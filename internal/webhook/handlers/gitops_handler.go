package handlers

import (
	"encoding/json"

	"github.com/Climatik-Project/Climatik-Project/internal/alert"
	"github.com/Climatik-Project/Climatik-Project/internal/webhook/runners"
)

type GitOpsAlertHandler struct {
	Runner runners.Runner
}

func (h *GitOpsAlertHandler) HandleAlert(payload []byte) error {
	var request alert.GitOpsAlertManager
	if err := json.Unmarshal(payload, &request); err != nil {
		return err
	}

	return h.Runner.Run()
}
