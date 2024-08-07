// factory.go
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

func CreateAlertService(config map[string]map[string]string) (*PubSub, error) {
	pubsub := NewPubSub()

	for managerType, managerConfig := range config {
		manager, err := NewAlertManager(AlertManagerType(managerType), managerConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create alert manager %s: %w", managerType, err)
		}
		pubsub.Subscribe("alerts", manager)
	}

	return pubsub, nil
}
