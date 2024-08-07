package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AlertLevel string

const (
	AlertLevelInfo     AlertLevel = "INFO"
	AlertLevelWarning  AlertLevel = "WARNING"
	AlertLevelCritical AlertLevel = "CRITICAL"
)

type SlackAlert struct {
	PodName       string
	PowerCapValue int
	CurrentPower  float64
	Devices       map[string]string
	Level         AlertLevel
	Timestamp     time.Time
}

type SlackAlertManager struct {
	webhookURL string
}

func NewSlackAlertManager(webhookURL string) (*SlackAlertManager, error) {
	return &SlackAlertManager{
		webhookURL: webhookURL,
	}, nil
}

func (s *SlackAlertManager) CreateAlert(podName string, powerCapValue int, devices map[string]string) error {
	alert := SlackAlert{
		PodName:       podName,
		PowerCapValue: powerCapValue,
		CurrentPower:  float64(powerCapValue), // You might want to get the actual current power from somewhere
		Devices:       devices,
		Level:         AlertLevelWarning,
		Timestamp:     time.Now(),
	}
	return s.sendWebhookAlert(alert)
}

func (s *SlackAlertManager) sendWebhookAlert(alert SlackAlert) error {
	message := s.createWebhookMessage(alert)
	return s.sendSlackWebhook(message)
}

func (s *SlackAlertManager) createWebhookMessage(alert SlackAlert) map[string]interface{} {
	attachment := s.createAttachment(alert)
	message := map[string]interface{}{
		"attachments": []interface{}{attachment},
	}
	return message
}

func (s *SlackAlertManager) createAttachment(alert SlackAlert) map[string]interface{} {
	return map[string]interface{}{
		"color":  getColorForAlertLevel(alert.Level),
		"title":  fmt.Sprintf("Power Cap Alert for Pod %s", alert.PodName),
		"text":   fmt.Sprintf("Power cap value of %d exceeded. Current power: %.2f", alert.PowerCapValue, alert.CurrentPower),
		"fields": getFieldsForDevices(alert.Devices),
		"footer": "Climatik Power Management",
		"ts":     alert.Timestamp.Unix(),
	}
}

func getColorForAlertLevel(level AlertLevel) string {
	switch level {
	case AlertLevelInfo:
		return "#36a64f" // Green
	case AlertLevelWarning:
		return "#ffcc00" // Yellow
	case AlertLevelCritical:
		return "#ff0000" // Red
	default:
		return "#808080" // Gray
	}
}

func getFieldsForDevices(devices map[string]string) []map[string]interface{} {
	var fields []map[string]interface{}
	for device, value := range devices {
		fields = append(fields, map[string]interface{}{
			"title": device,
			"value": value,
			"short": true,
		})
	}
	return fields
}

func (s *SlackAlertManager) sendSlackWebhook(message map[string]interface{}) error {
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	resp, err := http.Post(s.webhookURL, "application/json", bytes.NewBuffer(jsonMessage))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error sending message to Slack: %s", resp.Status)
	}

	return nil
}
