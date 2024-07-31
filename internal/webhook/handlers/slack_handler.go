package handlers

import (
	"encoding/json"

	"github.com/Climatik-Project/Climatik-Project/internal/alert"
	"github.com/Climatik-Project/Climatik-Project/internal/webhook/runners"
)

type SlackAlertHandler struct {
	Runner runners.Runner
}

func (h *SlackAlertHandler) HandleAlert(payload []byte) error {
	var request alert.SlackAlert
	if err := json.Unmarshal(payload, &request); err != nil {
		return err
	}

	return h.Runner.Run()
}
