// git_ops.go
package alert

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type GitOpsAlertManager struct {
	repoURL string
	repoDir string
}

func NewGitOpsAlertManager(repoURL, repoDir string) (*GitOpsAlertManager, error) {
	return &GitOpsAlertManager{
		repoURL: repoURL,
		repoDir: repoDir,
	}, nil
}

func (g *GitOpsAlertManager) CreateAlert(podName string, powerCapValue int, devices map[string]string) error {
	alertMessage := g.formatAlertMessage(podName, powerCapValue, devices)
	alertFile := fmt.Sprintf("%s/alerts/%s-alert.yaml", g.repoDir, podName)

	if err := g.writeAlertFile(alertFile, alertMessage); err != nil {
		return fmt.Errorf("failed to write alert file: %w", err)
	}

	if err := g.commitAndPushChanges(podName); err != nil {
		return fmt.Errorf("failed to commit and push changes: %w", err)
	}

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

func (g *GitOpsAlertManager) writeAlertFile(alertFile, alertMessage string) error {
	return os.WriteFile(alertFile, []byte(alertMessage), 0644)
}

func (g *GitOpsAlertManager) commitAndPushChanges(podName string) error {
	commands := [][]string{
		{"git", "-C", g.repoDir, "add", "."},
		{"git", "-C", g.repoDir, "commit", "-m", fmt.Sprintf("Add alert for pod %s", podName)},
		{"git", "-C", g.repoDir, "push", "origin", "main"},
	}

	for _, cmd := range commands {
		if err := exec.Command(cmd[0], cmd[1:]...).Run(); err != nil {
			return fmt.Errorf("failed to execute command %v: %w", cmd, err)
		}
	}

	return nil
}
