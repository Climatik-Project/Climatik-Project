package webhook

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Climatik-Project/Climatik-Project/internal/webhook/handlers"
)

var (
	signingSecret string
	slackBotToken string
)

func init() {
	signingSecret = os.Getenv("SLACK_SIGNING_SECRET")
	slackBotToken = os.Getenv("SLACK_BOT_TOKEN")
	if signingSecret == "" || slackBotToken == "" {
		log.Fatal("SLACK_SIGNING_SECRET and SLACK_BOT_TOKEN must be set")
	}
}

func CreateSlackWebhook(port int) {
	slackHandler, err := handlers.NewSlackHandler(signingSecret, slackBotToken)
	if err != nil {
		log.Fatalf("Failed to create SlackHandler: %v", err)
	}

	http.HandleFunc("/slack/command", slackHandler.CommandHandler)
	http.HandleFunc("/slack/interact", slackHandler.InteractionHandler)

	portStr := fmt.Sprintf(":%d", port)
	log.Printf("Starting server on port %d", port)
	if err := http.ListenAndServe(portStr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func main() {
	CreateSlackWebhook(8088)
}
