package alert

import "github.com/Climatik-Project/Climatik-Project/api/v1alpha1"

// AlertManager is the interface for creating power capping alerts
type AlertManager interface {
	CreateAlert(podName string, powerCapValue int, devices map[string]string, config *v1alpha1.PowerCappingConfig) error
}
