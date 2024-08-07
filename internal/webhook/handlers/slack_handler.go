package handlers

import (
	"encoding/json"

	adapters "github.com/Climatik-Project/Climatik-Project/internal/alert/adapters"
	"github.com/Climatik-Project/Climatik-Project/internal/webhook/runners"
)

type SlackAlertHandler struct {
	Runner runners.Runner
}

func (h *SlackAlertHandler) HandleAlert(payload []byte) error {
	var request adapters.SlackAlert
	if err := json.Unmarshal(payload, &request); err != nil {
		return err
	}

	return h.Runner.Run()
}
