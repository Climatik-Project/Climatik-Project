package handlers

import (
	"encoding/json"

	adapters "github.com/Climatik-Project/Climatik-Project/internal/alert/adapters"
	"github.com/Climatik-Project/Climatik-Project/internal/webhook/runners"
)

type PrometheusAlertHandler struct {
	Runner runners.Runner
}

func (h *PrometheusAlertHandler) HandleAlert(payload []byte) error {
	var request adapters.PrometheusAlert
	if err := json.Unmarshal(payload, &request); err != nil {
		return err
	}
	return h.Runner.Run()
}
