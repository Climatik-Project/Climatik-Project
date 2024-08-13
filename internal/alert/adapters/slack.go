package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	powercappingv1alpha1 "github.com/Climatik-Project/Climatik-Project/api/v1alpha1"
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
	Config        *powercappingv1alpha1.PowerCappingConfig
	Message       string
}

type SlackAlertManager struct {
	webhookURL string
}

func NewSlackAlertManager(webhookURL string) (*SlackAlertManager, error) {
	return &SlackAlertManager{
		webhookURL: webhookURL,
	}, nil
}

func (s *SlackAlertManager) CreateAlert(podName string, powerCapValue int, devices map[string]string, config *powercappingv1alpha1.PowerCappingConfig) error {
	currentPower := float64(powerCapValue) // You might want to get the actual current power from somewhere

	alert := SlackAlert{
		PodName:       podName,
		PowerCapValue: powerCapValue,
		CurrentPower:  currentPower,
		Devices:       devices,
		Level:         AlertLevelWarning,
		Timestamp:     time.Now(),
		Config:        config,
	}

	// Create a more informative message
	message := fmt.Sprintf("*Power Capping Alert for pod %s*\n", podName)
	message += fmt.Sprintf("Current power: %.2f watts\n", currentPower)
	message += fmt.Sprintf("Power cap: %d watts\n", powerCapValue)
	message += fmt.Sprintf("Devices: %v\n", devices)
	message += "\n*Configuration Details:*\n"
	message += fmt.Sprintf("Workload Type: %s\n", config.Spec.WorkloadType)
	message += fmt.Sprintf("Efficiency Level: %s\n", config.Spec.EfficiencyLevel)
	message += fmt.Sprintf("Power Cap Kind: %s\n", config.Spec.PowerCappingSpec.Kind)

	if config.Spec.PowerCappingSpec.Kind == powercappingv1alpha1.RelativePowerCapOfPeakPowerConsumptionInPercentage {
		message += fmt.Sprintf("Power Cap Percentage: %d%%\n", config.Spec.PowerCappingSpec.RelativePowerCapInPercentageSpec.PowerCapPercentage)
		message += fmt.Sprintf("Sample Window: %d seconds\n", config.Spec.PowerCappingSpec.RelativePowerCapInPercentageSpec.SampleWindow)
	}

	message += "\n*Actions:*\n"
	message += "To modify the power capping configuration, click here: <your_app_url>/modify-config"

	// Update the alert struct with the new message
	alert.Message = message

	return s.sendWebhookAlert(alert)
}

func (s *SlackAlertManager) sendWebhookAlert(alert SlackAlert) error {
	payload := map[string]interface{}{
		"text": alert.Message,
		"attachments": []map[string]interface{}{
			{
				"color": getColorForAlertLevel(alert.Level),
				"fields": []map[string]interface{}{
					{"title": "Pod", "value": alert.PodName, "short": true},
					{"title": "Power Cap", "value": fmt.Sprintf("%d W", alert.PowerCapValue), "short": true},
					{"title": "Current Power", "value": fmt.Sprintf("%.2f W", alert.CurrentPower), "short": true},
					{"title": "Timestamp", "value": alert.Timestamp.Format(time.RFC3339), "short": true},
				},
			},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	resp, err := http.Post(s.webhookURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to send Slack alert: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack API returned non-OK status: %s", resp.Status)
	}

	return nil
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
