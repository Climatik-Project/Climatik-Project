package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/Climatik-Project/Climatik-Project/internal/webhook/runners"
)

type Alert struct {
	Status      string            `json:"status"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}

type AlertManagerPayload struct {
	Alerts []Alert `json:"alerts"`
}

type PrometheusAlertHandler struct {
	Runner runners.Runner
}

func (h *PrometheusAlertHandler) HandleAlert(payload []byte) error {
	var request AlertManagerPayload
	if err := json.Unmarshal(payload, &request); err != nil {
		return err
	}

	for _, alert := range request.Alerts {
		fmt.Printf("Processing alert: %v\n", alert)
		if err := h.Runner.Run(); err != nil {
			return err
		}
	}
	return nil
}
