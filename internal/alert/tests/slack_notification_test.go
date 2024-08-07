package alert

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	adapters "github.com/Climatik-Project/Climatik-Project/internal/alert/adapters"
)

func TestSlackAlertActualNotification(t *testing.T) {
	// Find the root directory
	rootDir, err := findRootDir()
	assert.NoError(t, err, "Failed to find root directory")

	// Load environment variables from .env file in the root directory
	err = godotenv.Load(filepath.Join(rootDir, ".env"))
	assert.NoError(t, err, "Error loading .env file")

	webhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	assert.NotEmpty(t, webhookURL, "SLACK_WEBHOOK_URL is not set")

	manager, err := adapters.NewSlackAlertManager(webhookURL)
	assert.NoError(t, err)

	devices := map[string]string{
		"CPU":    "80% usage",
		"Memory": "70% usage",
	}

	err = manager.CreateAlert("test-pod", 100, devices)
	assert.NoError(t, err)
}

// findRootDir tries to find the root directory of the project
func findRootDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", os.ErrNotExist
		}
		dir = parent
	}
}
