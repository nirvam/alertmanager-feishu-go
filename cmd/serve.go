package cmd

import (
	"log"

	"github.com/nirvam/alertmanager-feishu-go/internal/config"
	"github.com/nirvam/alertmanager-feishu-go/internal/webhook"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Alertmanager Feishu service",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfig()
		if cfg.FeishuWebhookURL == "" {
			log.Fatal("FEISHU_WEBHOOK_URL is required")
		}
		if err := webhook.Start(cfg); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
