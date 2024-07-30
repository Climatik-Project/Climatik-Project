// prom_alert.go
package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

type PrometheusAlertManager struct {
	client          v1.API
	alertmanagerURL string
}

func NewPrometheusAlertManager(prometheusAddress string) (*PrometheusAlertManager, error) {
	client, err := api.NewClient(api.Config{Address: prometheusAddress})
	if err != nil {
		return nil, fmt.Errorf("failed to create Prometheus client: %w", err)
	}
	return &PrometheusAlertManager{
		client:          v1.NewAPI(client),
		alertmanagerURL: prometheusAddress + "/api/v1/alerts",
	}, nil
}

func (p *PrometheusAlertManager) CreateAlert(podName string, powerCapValue int, devices map[string]string) error {
	alert := p.formatPrometheusAlert(podName, powerCapValue, devices)
	if err := p.sendAlertToPrometheus(alert); err != nil {
		return fmt.Errorf("failed to send alert to Prometheus: %w", err)
	}
	return nil
}

func (p *PrometheusAlertManager) formatPrometheusAlert(podName string, powerCapValue int, devices map[string]string) PrometheusAlert {
	deviceStr := ""
	for device, value := range devices {
		deviceStr += fmt.Sprintf("%s:%s,", device, value)
	}
	deviceStr = deviceStr[:len(deviceStr)-1] // Remove trailing comma

	return PrometheusAlert{
		Labels: map[string]string{
			"alertname": "PowerCappingAlert",
			"severity":  "critical",
			"pod":       podName,
		},
		Annotations: map[string]string{
			"summary":     fmt.Sprintf("Power capping alert for pod %s", podName),
			"description": fmt.Sprintf("The pod is exceeding the power cap of %d watts. Devices: %s", powerCapValue, deviceStr),
		},
		StartsAt: time.Now(),
	}
}

func (p *PrometheusAlertManager) sendAlertToPrometheus(alert PrometheusAlert) error {
	alertBody, err := json.Marshal([]PrometheusAlert{alert})
	if err != nil {
		return fmt.Errorf("failed to marshal alert: %w", err)
	}

	req, err := http.NewRequest("POST", p.alertmanagerURL, bytes.NewBuffer(alertBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send alert to Prometheus Alertmanager: %s", resp.Status)
	}

	return nil
}

type PrometheusAlert struct {
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	StartsAt    time.Time         `json:"startsAt"`
}
