package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type SlackAlertManager struct {
	webhookURL string
}

type slackPayload struct {
	Text string `json:"text"`
}

func NewSlackAlertManager(webhookURL string) (*SlackAlertManager, error) {
	return &SlackAlertManager{
		webhookURL: webhookURL,
	}, nil
}

func (s *SlackAlertManager) CreateAlert(podName string, powerCapValue int, devices map[string]string) error {
	alertMessage := fmt.Sprintf("Slack Alert: Pod %s exceeded the power cap value of %d. Devices: %v", podName, powerCapValue, devices)
	payload := slackPayload{Text: alertMessage}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", s.webhookURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send alert to Slack, status code: %d", resp.StatusCode)
	}

	return nil
}
