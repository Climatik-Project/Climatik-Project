package runners

import "fmt"

type RunnerFactory struct{}

func (f *RunnerFactory) GetRunner(runnerType string, path string) (Runner, error) {
	switch runnerType {
	case "ansible":
		return &AnsibleRunner{PlaybookPath: path}, nil
	case "kubernetes":
		return &KubernetesRunner{JobManifestPath: path}, nil
	default:
		return nil, fmt.Errorf("unknown runner type: %s", runnerType)
	}
}
