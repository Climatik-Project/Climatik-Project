package alert

// AlertManager is the interface for creating power capping alerts
type AlertManager interface {
	CreateAlert(podName string, powerCapValue int, devices map[string]string) error
}
