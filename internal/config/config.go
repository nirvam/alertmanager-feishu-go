package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds the application configuration.
type Config struct {
	FeishuWebhookURL string
	FeishuSecret     string
	MessageType      string // "text" or "interactive"
	AppHost          string
	AppPort          int
}

// LoadConfig loads configuration from environment variables.
func LoadConfig() *Config {
	// Load .env if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	port, _ := strconv.Atoi(getEnv("APP_PORT", "8080"))
	return &Config{
		FeishuWebhookURL: os.Getenv("FEISHU_WEBHOOK_URL"),
		FeishuSecret:     os.Getenv("FEISHU_SECRET"),
		MessageType:      getEnv("MESSAGE_TYPE", "interactive"),
		AppHost:          getEnv("APP_HOST", "0.0.0.0"),
		AppPort:          port,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
