package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"text/template"

	"github.com/Climatik-Project/Climatik-Project/internal/webhook/runners"
)

type GitOpsAlertHandler struct {
	Runner runners.Runner
}

const gitOpsTemplate = `apiVersion: v1
kind: ConfigMap
metadata:
  name: alert-{{ .AlertName }}
data:
  alert: |
    status: {{ .Status }}
    labels: {{ .Labels }}
    annotations: {{ .Annotations }}
`

func createGitOpsFile(alert Alert) (string, error) {
	tmpl, err := template.New("gitOps").Parse(gitOpsTemplate)
	if err != nil {
		return "", err
	}

	var data bytes.Buffer
	if err := tmpl.Execute(&data, map[string]interface{}{
		"AlertName":   alert.Labels["alertname"],
		"Status":      alert.Status,
		"Labels":      alert.Labels,
		"Annotations": alert.Annotations,
	}); err != nil {
		return "", err
	}

	filePath := fmt.Sprintf("/tmp/alert-%s.yaml", alert.Labels["alertname"])
	if err := os.WriteFile(filePath, data.Bytes(), 0644); err != nil {
		return "", err
	}

	return filePath, nil
}

func (h *GitOpsAlertHandler) HandleAlert(payload []byte) error {
	var request AlertManagerPayload
	if err := json.Unmarshal(payload, &request); err != nil {
		return err
	}

	for _, alert := range request.Alerts {
		filePath, err := createGitOpsFile(alert)
		if err != nil {
			return err
		}

		// Replace the file path in the runner
		if runner, ok := h.Runner.(*runners.KubernetesRunner); ok {
			runner.JobManifestPath = filePath
		}

		if err := h.Runner.Run(); err != nil {
			return err
		}
	}
	return nil
}
