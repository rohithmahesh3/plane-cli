//go:build integration
// +build integration

package workspace

import (
	"bytes"
	"os"
	"testing"

	"github.com/rohithmahesh3/plane-cli/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestEnvironment(t *testing.T) {
	apiKey := os.Getenv("PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("PLANE_API_KEY not set, skipping integration test")
	}

	workspace := os.Getenv("PLANE_WORKSPACE")
	if workspace == "" {
		workspace = "test-workspace"
	}

	apiHost := os.Getenv("PLANE_API_HOST")
	if apiHost == "" {
		apiHost = "https://api.plane.so"
	}

	config.Cfg.APIHost = apiHost
	config.Cfg.DefaultWorkspace = workspace
	config.Cfg.OutputFormat = "json"

	originalService := config.KeyringService
	config.KeyringService = "plane-cli-cmd-test"
	t.Cleanup(func() { config.KeyringService = originalService; config.DeleteAPIKey() })

	err := config.SetAPIKey(apiKey)
	require.NoError(t, err)
}

func TestWorkspaceInfoIntegration(t *testing.T) {
	setupTestEnvironment(t)

	cmd := infoCmd
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err := runInfo(cmd, []string{})
	assert.NoError(t, err)
}
