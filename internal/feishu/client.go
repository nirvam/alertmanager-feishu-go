package feishu

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// GenSign generates the signature for Feishu webhook.
func GenSign(secret string, timestamp int64) (string, error) {
	// 拼接timestamp和secret
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, secret)

	h := hmac.New(sha256.New, []byte(stringToSign))
	// 飞书的签名逻辑实际上是：
	// stringToSign = timestamp + "\n" + secret
	// 然后以这个 stringToSign 为内容，计算 HMAC-SHA256？
	// 让我们参考 Python 版的逻辑：
	// string_to_sign = f"{timestamp}\n{secret}"
	// hmac_code = hmac.new(string_to_sign.encode("utf-8"), digestmod=hashlib.sha256).digest()
	// 注意：Python 的 hmac.new 如果不传 key，默认可能不太一样。
	// 回看 Python 代码：hmac.new(string_to_sign.encode("utf-8"), digestmod=hashlib.sha256).digest()
	// 这里第一个参数 string_to_sign 居然是作为 KEY 传入的，而 DATA 是空的！

	h.Write([]byte(""))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature, nil
}

// SendMessage sends a message to Feishu webhook.
func SendMessage(webhookURL string, secret string, msg Message) error {
	if secret != "" {
		timestamp := time.Now().Unix()
		sign, err := GenSign(secret, timestamp)
		if err != nil {
			return fmt.Errorf("failed to generate signature: %w", err)
		}
		msg.Timestamp = fmt.Sprintf("%d", timestamp)
		msg.Sign = sign
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("Feishu Response: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("feishu returned non-200 status: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err == nil {
		if code, ok := result["code"].(float64); ok && code != 0 {
			return fmt.Errorf("feishu API error: code=%v, msg=%v", code, result["msg"])
		}
	}

	return nil
}
