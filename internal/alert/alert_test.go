// alert_test.go
package alert

import (
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAlertManager is a mock implementation of the AlertManager interface
type MockAlertManager struct {
	mock.Mock
}

var execCommand = exec.Command

type MockSlackClient struct {
	mock.Mock
}

func (m *MockSlackClient) PostMessage(channelID string, options ...slack.MsgOption) (string, string, error) {
	args := m.Called(channelID, options)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockAlertManager) CreateAlert(podName string, powerCapValue int, devices map[string]string) error {
	args := m.Called(podName, powerCapValue, devices)
	return args.Error(0)
}

// TestNewAlertService tests the creation of a new AlertService
func TestNewAlertService(t *testing.T) {
	config := map[string]map[string]string{
		"prometheus": {"prometheusAddress": "http://prometheus:9090"},
		"gitops":     {"repoURL": "https://github.com/test/repo.git", "repoDir": "/tmp/repo"},
		"slack":      {"webhookURL": "http://slack.com/webhook", "token": "test-token", "channel": "#test"},
	}

	service, err := NewAlertService(config)
	assert.NoError(t, err)
	assert.NotNil(t, service)
	assert.NotNil(t, service.pubsub)
}

// TestAlertServiceSendAlert tests the SendAlert method of AlertService
func TestAlertServiceSendAlert(t *testing.T) {
	pubsub := NewPubSub()
	service := &AlertService{pubsub: pubsub}

	mockManager := new(MockAlertManager)
	mockManager.On("CreateAlert", "test-pod", 100, map[string]string{"cpu": "high"}).Return(nil)

	pubsub.Subscribe("alerts", mockManager)

	err := service.SendAlert("test-pod", 100, map[string]string{"cpu": "high"})
	assert.NoError(t, err)

	// Wait for goroutine to finish
	time.Sleep(100 * time.Millisecond)

	mockManager.AssertExpectations(t)
}

// TestPrometheusAlertManager tests the PrometheusAlertManager
func TestPrometheusAlertManager(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/alerts", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	manager, err := NewPrometheusAlertManager(server.URL)
	assert.NoError(t, err)

	err = manager.CreateAlert("test-pod", 100, map[string]string{"cpu": "high"})
	assert.NoError(t, err)
}

// TestGitOpsAlertManager tests the GitOpsAlertManager
func TestGitOpsAlertManager(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gitops-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	manager, err := NewGitOpsAlertManager("https://github.com/test/repo.git", tempDir)
	assert.NoError(t, err)

	// Mock git commands
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	err = manager.CreateAlert("test-pod", 100, map[string]string{"cpu": "high"})
	assert.NoError(t, err)

	// Check if file was created
	_, err = os.Stat(tempDir + "/alerts/test-pod-alert.yaml")
	assert.NoError(t, err)
}

// TestSlackAlertManager tests the SlackAlertManager
func TestSlackAlertManager(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	manager, err := NewSlackAlertManager(server.URL, "test-token", "#test-channel")
	assert.NoError(t, err)

	// Mock Slack client
	mockClient := &MockSlackClient{}
	manager.client = mockClient

	mockClient.On("PostMessage", "#test-channel", mock.Anything).Return("", "", nil)

	err = manager.CreateAlert("test-pod", 100, map[string]string{"cpu": "high"})
	assert.NoError(t, err)

	mockClient.AssertExpectations(t)
}

// mockExecCommand is a helper function to mock exec.Command
func mockExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

// TestHelperProcess isn't a real test. It's used to mock exec.Command
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// mocked behavior
	os.Exit(0)
}
