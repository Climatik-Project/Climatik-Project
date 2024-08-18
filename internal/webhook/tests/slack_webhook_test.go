package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Climatik-Project/Climatik-Project/internal/webhook"
	"github.com/Climatik-Project/Climatik-Project/internal/webhook/factory"
	"github.com/Climatik-Project/Climatik-Project/internal/webhook/handlers"
	"github.com/slack-go/slack"
)

func TestCreateSlackWebhook(t *testing.T) {
	// Set up environment variables for testing
	os.Setenv("SLACK_SIGNING_SECRET", "test_secret")
	os.Setenv("SLACK_BOT_TOKEN", "test_token")

	// Create a channel to signal when the server has started
	started := make(chan bool)

	// Start the server in a goroutine
	go func() {
		webhook.CreateSlackWebhook(8088)
		started <- true
	}()

	// Wait for the server to start
	<-started

	// Test that the server is running and endpoints are registered
	resp, err := http.Get("http://localhost:8088/slack/command")
	if err != nil {
		t.Fatalf("Failed to reach /slack/command endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", resp.StatusCode)
	}

	// Clean up environment variables
	os.Unsetenv("SLACK_SIGNING_SECRET")
	os.Unsetenv("SLACK_BOT_TOKEN")
}

func TestSlackHandlerCommandHandler(t *testing.T) {
	handler, _ := handlers.NewSlackHandler("test_secret", "test_token")

	// Create a test request
	slashCommand := &slack.SlashCommand{
		Command: "/modify-power-config",
	}
	payload, _ := json.Marshal(slashCommand)
	req := httptest.NewRequest("POST", "/slack/command", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler.CommandHandler(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := `{"blocks":[{"type":"section","text":{"type":"mrkdwn","text":"Choose a parameter to modify:"}},{"type":"actions","elements":[{"type":"static_select","placeholder":{"type":"plain_text","text":"Select a parameter"},"action_id":"select_parameter","options":[{"text":{"type":"plain_text","text":"Efficiency Level"},"value":"efficiency_level","description":{"type":"plain_text","text":"Set the efficiency level"}},{"text":{"type":"plain_text","text":"Power Cap Percentage"},"value":"power_cap_percentage","description":{"type":"plain_text","text":"Set the power cap percentage"}}]}]}]}`
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestSlackHandlerInteractionHandler(t *testing.T) {
	handler, _ := handlers.NewSlackHandler("test_secret", "test_token")

	// Create a test request
	interaction := &slack.InteractionCallback{
		Type: slack.InteractionTypeBlockActions,
		ActionCallback: slack.ActionCallbacks{
			BlockActions: []*slack.BlockAction{
				{
					ActionID: "select_parameter",
					SelectedOption: slack.OptionBlockObject{
						Value: "efficiency_level",
					},
				},
			},
		},
	}
	payload, _ := json.Marshal(interaction)
	req := httptest.NewRequest("POST", "/slack/interact", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler.InteractionHandler(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Note: The actual response body would depend on the implementation of openParameterInputModal
	// You may want to mock this function or check for a specific part of the response
}

func TestAlertHandlerFactory(t *testing.T) {
	factory := &factory.AlertHandlerFactory{}

	// Test Slack handler
	slackHandler, err := factory.GetHandler("slack")
	if err != nil {
		t.Errorf("Failed to get Slack handler: %v", err)
	}
	if _, ok := slackHandler.(*handlers.SlackHandler); !ok {
		t.Errorf("Expected SlackHandler, got %T", slackHandler)
	}

	// Test unsupported handler
	_, err = factory.GetHandler("unsupported")
	if err == nil {
		t.Error("Expected error for unsupported handler, got nil")
	}

	// Test unimplemented handlers
	_, err = factory.GetHandler("prometheus")
	if err == nil {
		t.Error("Expected error for unimplemented PrometheusAlertHandler, got nil")
	}

	_, err = factory.GetHandler("gitops")
	if err == nil {
		t.Error("Expected error for unimplemented GitOpsAlertHandler, got nil")
	}
}
