package interfaces

import "github.com/Climatik-Project/Climatik-Project/internal/webhook/runners"

type AlertHandler interface {
	HandleAlert(payload []byte) error
	UpdatePowerCappingConfig(param, value string) error
}

type AlertHandlerFactory interface {
	GetHandler(source string, runner runners.Runner) (AlertHandler, error)
}
