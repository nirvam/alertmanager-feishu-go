package formatter

import (
	"testing"
	"time"

	"github.com/nirvam/alertmanager-feishu-go/pkg/alertmanager"
	"github.com/stretchr/testify/assert"
)

func TestFormatText(t *testing.T) {
	payload := alertmanager.Payload{
		Status: "firing",
		CommonLabels: map[string]string{
			"alertname": "TestAlert",
		},
		Alerts: []alertmanager.Alert{
			{
				Status: "firing",
				Annotations: map[string]string{
					"summary":     "This is a test alert",
					"description": "If you see this, the service is working!",
				},
				StartsAt: time.Now(),
			},
		},
	}

	msg := FormatText(payload)
	assert.Equal(t, "text", msg.MsgType)
	assert.Contains(t, msg.Content["text"].(string), "Status: FIRING")
	assert.Contains(t, msg.Content["text"].(string), "Summary: This is a test alert")
}

func TestFormatCard(t *testing.T) {
	payload := alertmanager.Payload{
		Status: "firing",
		CommonLabels: map[string]string{
			"alertname": "TestAlert",
		},
		Alerts: []alertmanager.Alert{
			{
				Status: "firing",
				Labels: map[string]string{
					"severity": "critical",
				},
				Annotations: map[string]string{
					"summary":     "This is a test alert",
					"description": "If you see this, the service is working!",
				},
				StartsAt:     time.Now(),
				GeneratorURL: "http://prometheus:9090",
			},
		},
	}

	msg := FormatCard(payload)
	assert.Equal(t, "interactive", msg.MsgType)
	assert.NotNil(t, msg.Card)
	assert.Equal(t, "2.0", msg.Card.Schema)
	assert.Equal(t, "red", msg.Card.Header.Template)
	assert.True(t, len(msg.Card.Body.Elements) > 0)
}
