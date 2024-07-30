package factory

import (
	"fmt"

	"github.com/Climatik-Project/Climatik-Project/internal/webhook/handlers"
	"github.com/Climatik-Project/Climatik-Project/internal/webhook/runners"
)

type AlertHandlerFactory struct{}

func (f *AlertHandlerFactory) GetHandler(source string, runner runners.Runner) (handlers.AlertHandler, error) {
	switch source {
	case "slack":
		return &handlers.SlackAlertHandler{Runner: runner}, nil
	case "prometheus":
		return &handlers.PrometheusAlertHandler{Runner: runner}, nil
	case "gitops":
		return &handlers.GitOpsAlertHandler{Runner: runner}, nil
	default:
		return nil, fmt.Errorf("unknown alert source: %s", source)
	}
}
