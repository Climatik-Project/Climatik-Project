# Setting Up and Using Slack Webhooks

This guide will walk you through the process of creating a Slack webhook and using it to send messages to a specific channel.

## 1. Creating a Slack App and Webhook

1. Go to the [Slack API website](https://api.slack.com/apps).
2. Click on "Create New App" and choose "From scratch".
3. Give your app a name and select the workspace where you want to use it.
4. Once your app is created, go to "Incoming Webhooks" in the left sidebar.
5. Toggle "Activate Incoming Webhooks" to On.
6. Scroll down and click "Add New Webhook to Workspace".
7. Choose the channel where you want the messages to appear.
8. Click "Allow" to authorize the webhook.
9. You'll see a new webhook URL. Copy this URL, as you'll need it later.

## 2. Setting Up Environment Variables

To keep your webhook URL secure, we'll use environment variables:

1. Create a `.env` file in your project root if it doesn't exist already.
2. Add the following line to your `.env` file:

    ```bash
    SLACK_WEBHOOK_URL=https://hooks.slack.com/services/YOUR/WEBHOOK/URL
    ```

3. Replace `https://hooks.slack.com/services/YOUR/WEBHOOK/URL` with the actual webhook URL you copied earlier.
4. Make sure your `.gitignore` file includes `.env` to avoid accidentally committing sensitive information.

## 3. Using the Webhook in Your Code

Here's how you can use the webhook in your Go code:

```go
package main

import (
 "os"
 "log"
 "github.com/joho/godotenv"
 "github.com/Climatik-Project/Climatik-Project/internal/alert/adapters"
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
 slackManager, err := adapters.NewSlackAlertManager(webhookURL, nil)
 if err != nil {
     log.Fatalf("Failed to create Slack alert manager: %v", err)
 }

 // Send a test alert
 err = slackManager.CreateAlert("test-pod", 100, map[string]string{"cpu": "high"})
 if err != nil {
     log.Fatalf("Failed to send alert: %v", err)
 }

 log.Println("Alert sent successfully!")
}
```

### Troubleshooting

If you're not seeing messages in Slack, double-check that your webhook URL is correct and that you've selected the right channel when setting up the webhook.
Ensure that your .env file is in the correct location and that the environment variable is being loaded properly.
Check the Slack App settings to make sure the app has the necessary permissions for the channel you're posting to.
