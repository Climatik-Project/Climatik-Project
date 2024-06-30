package alert

import (
	"fmt"
	"os"
	"os/exec"
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
	alertMessage := fmt.Sprintf("GitOps Alert: Pod %s exceeded the power cap value of %d. Devices: %v", podName, powerCapValue, devices)
	alertFile := fmt.Sprintf("%s/alerts/%s-alert.txt", g.repoDir, podName)

	err := os.WriteFile(alertFile, []byte(alertMessage), 0644)
	if err != nil {
		return err
	}

	// Commit the alert file to the Git repository
	cmd := exec.Command("git", "-C", g.repoDir, "add", alertFile)
	err = cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("git", "-C", g.repoDir, "commit", "-m", fmt.Sprintf("Add alert for pod %s", podName))
	err = cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("git", "-C", g.repoDir, "push")
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
