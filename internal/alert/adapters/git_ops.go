package adapters

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	powercappingv1alpha1 "github.com/Climatik-Project/Climatik-Project/api/v1alpha1"
)

type GitOpsAlertManager struct {
	repoDir string
	alerts  map[string]string // In-memory storage for alerts
}

func NewGitOpsAlertManager(repoURL, repoDir string) (*GitOpsAlertManager, error) {
	// We're not using repoURL in this simplified version, but keeping it for interface consistency
	return &GitOpsAlertManager{
		repoDir: repoDir,
		alerts:  make(map[string]string),
	}, nil
}

func (g *GitOpsAlertManager) CreateAlert(podName string, powerCapValue int, devices map[string]string, config *powercappingv1alpha1.PowerCappingConfig) error {
	alertMessage := g.formatAlertMessage(podName, powerCapValue, devices)
	alertFile := filepath.Join(g.repoDir, "alerts", podName+"-alert.yaml")
	g.alerts[alertFile] = alertMessage
	return nil
}

func (g *GitOpsAlertManager) formatAlertMessage(podName string, powerCapValue int, devices map[string]string) string {
	deviceStr := []string{}
	for device, value := range devices {
		deviceStr = append(deviceStr, fmt.Sprintf("%s: %s", device, value))
	}

	return fmt.Sprintf(`
		apiVersion: climatik.io/v1
		kind: PowerAlert
		metadata:
		name: %s-alert
		creationTimestamp: %s
		spec:
		podName: %s
		powerCapValue: %d
		devices: %s
	`, podName, time.Now().Format(time.RFC3339), podName, powerCapValue, strings.Join(deviceStr, "\n    "))
}

// GetAlerts returns the current alerts (for testing purposes)
func (g *GitOpsAlertManager) GetAlerts() map[string]string {
	return g.alerts
}
