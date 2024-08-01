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

// TestSlackAlertManager tests the SlackAlertManager
func TestSlackAlertManager(t *testing.T) {
	mockClient := new(MockSlackClient)
	manager := &SlackAlertManager{
		webhookURL: "http://slack.com/webhook",
		channel:    "#test-channel",
		client:     mockClient,
	}

	mockClient.On("PostMessage", "#test-channel", mock.Anything).Return("channel", "timestamp", nil)

	err := manager.CreateAlert("test-pod", 100, map[string]string{"cpu": "high"})
	assert.NoError(t, err)

	mockClient.AssertExpectations(t)
}

// TestPrometheusAlertManager tests the PrometheusAlertManager
func TestNewPrometheusAlertManager(t *testing.T) {
	manager, err := NewPrometheusAlertManager("http://prometheus:9090")
	assert.NoError(t, err)
	assert.NotNil(t, manager)
	assert.Equal(t, "http://prometheus:9090/api/v1/alerts", manager.alertmanagerURL)
}

func TestFormatPrometheusAlert(t *testing.T) {
	manager := &PrometheusAlertManager{}
	alert := manager.formatPrometheusAlert("test-pod", 100, map[string]string{"cpu": "high", "memory": "low"})

	assert.Equal(t, "PowerCappingAlert", alert.Labels["alertname"])
	assert.Equal(t, "critical", alert.Labels["severity"])
	assert.Equal(t, "test-pod", alert.Labels["pod"])
	assert.Contains(t, alert.Annotations["description"], "100 watts")
	assert.Contains(t, alert.Annotations["description"], "cpu:high,memory:low")
}

func TestCreateAlert(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/alerts", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var alerts []PrometheusAlert
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

	manager := &PrometheusAlertManager{
		alertmanagerURL: server.URL + "/api/v1/alerts",
	}

	err := manager.CreateAlert("test-pod", 100, map[string]string{"cpu": "high"})
	assert.NoError(t, err)
}

func TestCreateAlertError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	manager := &PrometheusAlertManager{
		alertmanagerURL: server.URL + "/api/v1/alerts",
	}

	err := manager.CreateAlert("test-pod", 100, map[string]string{"cpu": "high"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send alert to Prometheus Alertmanager")
}

func TestSendAlertToPrometheus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	manager := &PrometheusAlertManager{
		alertmanagerURL: server.URL,
	}

	alert := PrometheusAlert{
		Labels:      map[string]string{"alertname": "TestAlert"},
		Annotations: map[string]string{"description": "Test description"},
		StartsAt:    time.Now(),
	}

	err := manager.sendAlertToPrometheus(alert)
	assert.NoError(t, err)
}

// TestGitOpsAlertManager tests the GitOpsAlertManager
func TestGitOpsAlertManager(t *testing.T) {
	manager, err := NewGitOpsAlertManager("https://github.com/test/repo.git", "/tmp/test-repo")
	assert.NoError(t, err)
	assert.NotNil(t, manager)

	gitOpsManager, ok := manager.(*GitOpsAlertManager)
	assert.True(t, ok, "manager should be of type *GitOpsAlertManager")

	err = manager.CreateAlert("test-pod", 100, map[string]string{"cpu": "high"})
	assert.NoError(t, err)

	alerts := gitOpsManager.GetAlerts()
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

	pubsub := NewPubSub()
	pubsub.Subscribe("alerts", mockSlackManager)
	pubsub.Subscribe("alerts", mockPrometheusManager)
	pubsub.Subscribe("alerts", mockGitOpsManager)

	service := &AlertService{pubsub: pubsub}

	mockSlackManager.On("CreateAlert", "test-pod", 100, map[string]string{"cpu": "high"}).Return(nil)
	mockPrometheusManager.On("CreateAlert", "test-pod", 100, map[string]string{"cpu": "high"}).Return(nil)
	mockGitOpsManager.On("CreateAlert", "test-pod", 100, map[string]string{"cpu": "high"}).Return(nil)

	err := service.SendAlert("test-pod", 100, map[string]string{"cpu": "high"})
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

func (m *MockAlertManager) CreateAlert(podName string, powerCapValue int, devices map[string]string) error {
	args := m.Called(podName, powerCapValue, devices)
	return args.Error(0)
}
