package alert

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/Climatik-Project/Climatik-Project/api/v1alpha1"
	alert "github.com/Climatik-Project/Climatik-Project/internal/alert"
	adapters "github.com/Climatik-Project/Climatik-Project/internal/alert/adapters"
)

// MockSlackClient is a mock implementation of the SlackClient interface
type MockSlackClient struct {
	mock.Mock
}

func (m *MockSlackClient) PostMessage(channelID string, options ...slack.MsgOption) (string, string, error) {
	args := m.Called(channelID, options)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockSlackClient) UpdateMessage(channelID, timestamp string, options ...slack.MsgOption) (string, string, string, error) {
	args := m.Called(channelID, timestamp, options)
	return args.String(0), args.String(1), args.String(2), args.Error(3)
}

type MockPrometheusAPI struct {
	mock.Mock
}

// Implement all methods of v1.API interface. Here's an example for the Query method:
func (m *MockPrometheusAPI) Query(ctx context.Context, query string, ts time.Time) (model.Value, v1.Warnings, error) {
	args := m.Called(ctx, query, ts)
	return args.Get(0).(model.Value), args.Get(1).(v1.Warnings), args.Error(2)
}

// MockGitOperations is a mock implementation of the GitOperations interface
type MockGitOperations struct {
	mock.Mock
}

func (m *MockGitOperations) Add(repoDir string, files ...string) error {
	args := m.Called(repoDir, files)
	return args.Error(0)
}

func (m *MockGitOperations) Commit(repoDir, message string) error {
	args := m.Called(repoDir, message)
	return args.Error(0)
}

func (m *MockGitOperations) Push(repoDir, remote, branch string) error {
	args := m.Called(repoDir, remote, branch)
	return args.Error(0)
}

func TestSlackAlertManagerMocked(t *testing.T) {
	// Create a test server that mimics Slack's webhook endpoint
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	manager, err := adapters.NewSlackAlertManager(server.URL)
	assert.NoError(t, err)

	devices := map[string]string{
		"CPU":    "80% usage",
		"Memory": "70% usage",
	}

	mockConfig := NewMockPowerCappingConfig()

	err = manager.CreateAlert("test-pod", 100, devices, mockConfig)

	assert.NoError(t, err)
}

// TestPrometheusAlertManager tests the PrometheusAlertManager
func TestNewPrometheusAlertManager(t *testing.T) {
	manager, err := adapters.NewPrometheusAlertManager("http://prometheus:9090")
	assert.NoError(t, err)
	assert.NotNil(t, manager)
	// We can't directly access alertmanagerURL as it's private
	// Instead, we can test the functionality that uses it
}

func TestFormatPrometheusAlert(t *testing.T) {
	manager, _ := adapters.NewPrometheusAlertManager("http://prometheus:9090")
	alert := manager.FormatPrometheusAlert("test-pod", 100, map[string]string{"cpu": "high", "memory": "low"})

	assert.Equal(t, "PowerCappingAlert", alert.Labels["alertname"])
	assert.Equal(t, "critical", alert.Labels["severity"])
	assert.Equal(t, "test-pod", alert.Labels["pod"])
	assert.Contains(t, alert.Annotations["description"], "100 watts")
	assert.Contains(t, alert.Annotations["description"], "cpu:high")
	assert.Contains(t, alert.Annotations["description"], "memory:low")

}

func TestCreateAlert(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/alerts", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var alerts []adapters.PrometheusAlert
		err := json.NewDecoder(r.Body).Decode(&alerts)
		require.NoError(t, err)
		assert.Len(t, alerts, 1)

		alert := alerts[0]
		assert.Equal(t, "PowerCappingAlert", alert.Labels["alertname"])
		assert.Equal(t, "test-pod", alert.Labels["pod"])
		assert.Contains(t, alert.Annotations["description"], "100 watts")

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	manager := &adapters.PrometheusAlertManager{
		AlertmanagerURL: server.URL + "/api/v1/alerts",
	}

	// Create a mock PowerCappingConfig
	mockConfig := NewMockPowerCappingConfig()

	err := manager.CreateAlert("test-pod", 100, map[string]string{"cpu": "high"}, mockConfig)
	assert.NoError(t, err)
}

func TestCreateAlertError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	manager := &adapters.PrometheusAlertManager{
		AlertmanagerURL: server.URL + "/api/v1/alerts",
	}

	mockConfig := NewMockPowerCappingConfig()

	err := manager.CreateAlert("test-pod", 100, map[string]string{"cpu": "high"}, mockConfig)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send alert to Prometheus Alertmanager")
}

func TestSendAlertToPrometheus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	manager := &adapters.PrometheusAlertManager{
		AlertmanagerURL: server.URL,
	}

	alert := adapters.PrometheusAlert{
		Labels:      map[string]string{"alertname": "TestAlert"},
		Annotations: map[string]string{"description": "Test description"},
		StartsAt:    time.Now(),
	}

	err := manager.SendAlertToPrometheus(alert)
	assert.NoError(t, err)
}

// TestGitOpsAlertManager tests the GitOpsAlertManager
func TestGitOpsAlertManager(t *testing.T) {
	manager, err := adapters.NewGitOpsAlertManager("https://github.com/test/repo.git", "/tmp/test-repo")
	assert.NoError(t, err)
	assert.NotNil(t, manager)

	mockConfig := NewMockPowerCappingConfig()

	err = manager.CreateAlert("test-pod", 100, map[string]string{"cpu": "high"}, mockConfig)
	assert.NoError(t, err)

	alerts := manager.GetAlerts()
	assert.Len(t, alerts, 1)

	alertFile := "/tmp/test-repo/alerts/test-pod-alert.yaml"
	alertContent, exists := alerts[alertFile]
	assert.True(t, exists)
	assert.Contains(t, alertContent, "test-pod")
	assert.Contains(t, alertContent, "powerCapValue: 100")
	assert.Contains(t, alertContent, "cpu: high")
}

// TestAlertService tests the AlertService
func TestAlertService(t *testing.T) {
	mockSlackManager := new(MockAlertManager)
	mockPrometheusManager := new(MockAlertManager)
	mockGitOpsManager := new(MockAlertManager)

	pubsub := alert.NewPubSub()
	pubsub.Subscribe("alerts", mockSlackManager)
	pubsub.Subscribe("alerts", mockPrometheusManager)
	pubsub.Subscribe("alerts", mockGitOpsManager)

	service := &alert.AlertService{Pubsub: pubsub}

	// Create a mock PowerCappingConfig
	mockConfig := NewMockPowerCappingConfig()

	mockSlackManager.On("CreateAlert", "test-pod", 100, map[string]string{"cpu": "high"}, mockConfig).Return(nil)
	mockPrometheusManager.On("CreateAlert", "test-pod", 100, map[string]string{"cpu": "high"}, mockConfig).Return(nil)
	mockGitOpsManager.On("CreateAlert", "test-pod", 100, map[string]string{"cpu": "high"}, mockConfig).Return(nil)

	err := service.SendAlert("test-pod", 100, map[string]string{"cpu": "high"}, mockConfig)
	assert.NoError(t, err)

	// Wait for goroutines to finish
	time.Sleep(100 * time.Millisecond)

	mockSlackManager.AssertExpectations(t)
	mockPrometheusManager.AssertExpectations(t)
	mockGitOpsManager.AssertExpectations(t)
}

// MockAlertManager is a mock implementation of the AlertManager interface
type MockAlertManager struct {
	mock.Mock
}

func (m *MockAlertManager) CreateAlert(podName string, powerCapValue int, devices map[string]string, config *v1alpha1.PowerCappingConfig) error {
	args := m.Called(podName, powerCapValue, devices, config)
	return args.Error(0)
}
