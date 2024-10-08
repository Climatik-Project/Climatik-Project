package factory

import (
	"fmt"

	"github.com/Climatik-Project/Climatik-Project/internal/webhook/handlers"
)

type AlertHandler interface {
	HandleAlert(body []byte) error
}

type AlertHandlerFactory struct{}

func (f *AlertHandlerFactory) GetHandler(source string) (AlertHandler, error) {
	switch source {
	case "slack":
		return &handlers.SlackHandler{}, nil
	case "prometheus":
		// Implement PrometheusAlertHandler creation
		return nil, fmt.Errorf("PrometheusAlertHandler not implemented")
	case "gitops":
		// Implement GitOpsAlertHandler creation
		return nil, fmt.Errorf("GitOpsAlertHandler not implemented")
	default:
		return nil, fmt.Errorf("unknown alert source: %s", source)
	}
}
