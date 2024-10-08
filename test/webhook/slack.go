package main

import (
	"log"
	"os"

	alert "github.com/Climatik-Project/Climatik-Project/internal/alert/adapters"
	mockConfig "github.com/Climatik-Project/Climatik-Project/internal/alert/tests"
	"github.com/joho/godotenv"
)

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get the webhook URL from environment variables
	webhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	if webhookURL == "" {
		log.Fatal("SLACK_WEBHOOK_URL is not set")
	}

	// Create a new SlackAlertManager
	slackManager, err := alert.NewSlackAlertManager(webhookURL)
	if err != nil {
		log.Fatalf("Failed to create Slack alert manager: %v", err)
	}

	mockConfig := mockConfig.NewMockPowerCappingConfig()

	// Send a test alert
	err = slackManager.CreateAlert("test-pod", 100, map[string]string{"cpu": "high"}, mockConfig)
	if err != nil {
		log.Fatalf("Failed to send alert: %v", err)
	}

	log.Println("Alert sent successfully!")
}
