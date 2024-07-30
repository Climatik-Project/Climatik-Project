package runners

import (
	"fmt"
	"os/exec"
)

type AnsibleRunner struct {
	PlaybookPath string
}

func (r *AnsibleRunner) Run() error {
	cmd := exec.Command("ansible-playbook", r.PlaybookPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run ansible-playbook: %v, output: %s", err, string(output))
	}
	fmt.Println(string(output))
	return nil
}
