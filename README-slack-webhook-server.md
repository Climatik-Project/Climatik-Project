# Setting up a Slack App for Power Capping Config

## Workflow

1. When a user types the slash command `/modify-power-config`, `handleModifyPowerConfig` will be called, presenting the user with options.
2. When the user selects an option, `handleBlockActions` will be called, which will open a modal using `openParameterInputModal`.
3. When the user submits the modal, `handleViewSubmission` will be called, which will then call `handleParameterUpdate` to process the update.

## 1. Create a new Slack App

1. Go to https://api.slack.com/apps
2. Click "Create New App"
3. Choose "From scratch"
4. Name your app (e.g., "Power Capping Config") and select your workspace

## 2. Set up Slash Commands

1. In the left sidebar, click on "Slash Commands"
2. Click "Create New Command"
3. Set the command to "/modify-power-config"
4. Set the Request URL to "https://your-domain.com/slack/command"
5. Add a short description and usage hint
6. Save the command

## 3. Set up Interactivity

1. In the left sidebar, click on "Interactivity & Shortcuts"
2. Turn on Interactivity
3. Set the Request URL to "https://your-domain.com/slack/interact"
4. Save the changes

## 4. Set up OAuth & Permissions

1. In the left sidebar, click on "OAuth & Permissions"
2. Under "Scopes", add the following Bot Token Scopes:

- commands
- chat:write
- im:write

3. Scroll up and click "Install to Workspace"
4. Authorize the app

## 5. Get your app credentials

1. In the left sidebar, click on "Basic Information"
2. Under "App Credentials", you'll find your Signing Secret
3. Go back to "OAuth & Permissions" to find your Bot User OAuth Token

## 6. Set up your environment variables

In your server environment, set these variables:

```bash
export SLACK_SIGNING_SECRET=your_signing_secret_here
export SLACK_BOT_TOKEN=your_bot_token_here
```

## 7. Update your Go code

Ensure your `CreateWebhook` function is being called with the correct port:

```go
func main() {
    // ... other setup code ...
    CreateWebhook(8080) // or whatever port you're using
}
```

## 8. Deploy your webhook

Deploy your Go application to a server that's accessible via HTTPS. Slack requires HTTPS for all webhook URLs.

## 9. Test your app

1. In your Slack workspace, type "/modify-power-config"
2. You should see the options to modify efficiency level or power cap percentage
3. Select an option, enter a new value, and submit
4. Verify that your backend receives and processes the request correctly

## 10. Error Handling and Logging

Implement proper error handling and logging in your webhook handlers to help with debugging.

## 11. Security Considerations

- Always verify the request signature using the signing secret
- Use HTTPS for all communications
- Implement proper access controls to ensure only authorized users can modify configurations

**Note**: Remember to replace "https://your-domain.com" with your actual domain where your webhook is hosted.

## Important URLs

Make sure to update your Slack app configuration to point to these webhook URLs:
- Set the slash command URL to `https://your-domain.com/slack/command`
- Set the interactivity request URL to `https://your-domain.com/slack/interact`
