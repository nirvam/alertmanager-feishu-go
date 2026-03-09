package cmd

import (
	"log"
	"time"

	"github.com/nirvam/alertmanager-feishu-go/internal/config"
	"github.com/nirvam/alertmanager-feishu-go/internal/feishu"
	"github.com/nirvam/alertmanager-feishu-go/internal/formatter"
	"github.com/nirvam/alertmanager-feishu-go/pkg/alertmanager"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Send a test message to Feishu",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfig()
		if cfg.FeishuWebhookURL == "" {
			log.Fatal("FEISHU_WEBHOOK_URL is required")
		}

		testPayload := alertmanager.Payload{
			Status: "firing",
			Alerts: []alertmanager.Alert{
				{
					Status: "firing",
					Labels: map[string]string{"alertname": "TestAlert", "severity": "critical"},
					Annotations: map[string]string{
						"summary":     "This is a test alert",
						"description": "If you see this, the alertmanager-feishu-go service is working!",
					},
					StartsAt:     time.Now(),
					GeneratorURL: "http://prometheus:9090",
				},
			},
			CommonLabels: map[string]string{"alertname": "TestAlert"},
		}

		var msg feishu.Message
		if cfg.MessageType == "interactive" {
			msg = formatter.FormatCard(testPayload)
		} else {
			msg = formatter.FormatText(testPayload)
		}

		log.Println("Sending test message...")
		if err := feishu.SendMessage(cfg.FeishuWebhookURL, cfg.FeishuSecret, msg); err != nil {
			log.Fatalf("Error: %v", err)
		}
		log.Println("Test message sent successfully!")
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
