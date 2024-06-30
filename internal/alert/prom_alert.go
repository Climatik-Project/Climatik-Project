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

var (
	AlertmanagerURL = getEnv("ALERTMANAGER_URL", "http://alertmanager:9093/api/v1/alerts")
)

type PrometheusAlert struct {
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	StartsAt    time.Time         `json:"startsAt"`
}

func sendAlertToPrometheus(alert string) error {
	alertStruct := PrometheusAlert{
		Labels: map[string]string{
			"alertname": "PowerCappingAlert",
			"severity":  "critical",
		},
		Annotations: map[string]string{
			"summary":     "Power capping alert",
			"description": alert,
		},
		StartsAt: time.Now(),
	}

	alertBody, err := json.Marshal([]PrometheusAlert{alertStruct})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", AlertmanagerURL, bytes.NewBuffer(alertBody))
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
		return fmt.Errorf("failed to send alert to Prometheus Alertmanager: %s", resp.Status)
	}

	return nil
}

type PrometheusAlertManager struct {
	client v1.API
}

func NewPrometheusAlertManager(prometheusAddress string) (*PrometheusAlertManager, error) {
	client, err := api.NewClient(api.Config{Address: prometheusAddress})
	if err != nil {
		return nil, err
	}
	return &PrometheusAlertManager{
		client: v1.NewAPI(client),
	}, nil
}

func (p *PrometheusAlertManager) CreateAlert(podName string, powerCap int, devices map[string]string) error {
	device := ""
	for k, label := range devices {
		device += k + ":" + label + ","
	}
	alert := fmt.Sprintf(`ALERT PowerCappingAlert
    IF rate(kepler_container_joules_total{pod='%s'}) > %d
    FOR 5m
    LABELS { severity="critical" }
    ANNOTATIONS {
        summary = "Power capping alert for pod %s",
        description = "The pod is exceeding the power cap of %d watts."
		device = "%s"
    }`, podName, powerCap, podName, powerCap, device)

	err := sendAlertToPrometheus(alert)

	return err
}
