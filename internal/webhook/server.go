package webhook

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/nirvam/alertmanager-feishu-go/internal/config"
	"github.com/nirvam/alertmanager-feishu-go/internal/feishu"
	"github.com/nirvam/alertmanager-feishu-go/internal/formatter"
	"github.com/nirvam/alertmanager-feishu-go/pkg/alertmanager"
)

// Start starts the webhook receiver server.
func Start(cfg *config.Config) error {
	http.HandleFunc("/webhook", handleWebhook(cfg))
	addr := fmt.Sprintf("%s:%d", cfg.AppHost, cfg.AppPort)
	log.Printf("Starting server on %s", addr)
	return http.ListenAndServe(addr, nil)
}

func handleWebhook(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var payload alertmanager.Payload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		var msg feishu.Message
		if cfg.MessageType == "interactive" {
			msg = formatter.FormatCard(payload)
		} else {
			msg = formatter.FormatText(payload)
		}

		if err := feishu.SendMessage(cfg.FeishuWebhookURL, cfg.FeishuSecret, msg); err != nil {
			log.Printf("Error sending message to Feishu: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
}
