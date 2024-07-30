package handlers

import (
	"encoding/json"

	"github.com/Climatik-Project/Climatik-Project/internal/webhook/runners"
)

type SlackEvent struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

type SlackRequest struct {
	Token       string     `json:"token"`
	TeamID      string     `json:"team_id"`
	APIAppID    string     `json:"api_app_id"`
	Event       SlackEvent `json:"event"`
	Type        string     `json:"type"`
	Challenge   string     `json:"challenge"`
	EventID     string     `json:"event_id"`
	EventTime   int64      `json:"event_time"`
	AuthedUsers []string   `json:"authed_users"`
}

type SlackAlertHandler struct {
	Runner runners.Runner
}

func (h *SlackAlertHandler) HandleAlert(payload []byte) error {
	var request SlackRequest
	if err := json.Unmarshal(payload, &request); err != nil {
		return err
	}

	if request.Event.Type == "message" {
		return h.Runner.Run()
	}
	return nil
}
