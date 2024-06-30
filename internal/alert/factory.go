package alert

import (
	"fmt"
)

type AlertManagerType string

const (
	Prometheus AlertManagerType = "prometheus"
	GitOps     AlertManagerType = "gitops"
	Slack      AlertManagerType = "slack"
)

func NewAlertManager(managerType AlertManagerType, config map[string]string) (AlertManager, error) {
	switch managerType {
	case Prometheus:
		return NewPrometheusAlertManager(config["prometheusAddress"])
	case GitOps:
		return NewGitOpsAlertManager(config["repoURL"], config["repoDir"])
	case Slack:
		return NewSlackAlertManager(config["webhookURL"])
	default:
		return nil, fmt.Errorf("unsupported alert manager type: %s", managerType)
	}
}
