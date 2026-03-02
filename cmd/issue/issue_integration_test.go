//go:build integration
// +build integration

package issue

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

	// Make sure we have a default project
	// For testing, we just need a string. It will be evaluated by API.
	// But actually we need a real project ID for "issue list"
	// We'll set a dummy project ID, and if it fails, it means the API works but project is wrong.
	// A better way is to pass a valid project ID, but user didn't give one.
	// Actually, issue list requires a project ID. We will test search which is workspace-level.

	config.Cfg.APIHost = apiHost
	config.Cfg.DefaultWorkspace = workspace
	config.Cfg.OutputFormat = "json"

	// Mock Keyring for test by setting it directly if needed,
	// but the API client in cmds uses config.GetAPIKey().
	// For this test, we must inject the API key to the keyring, or modify API client to read from env.
	// We will write it to a test keyring if possible, or assume it's in the environment.

	// Wait, cmd commands will call api.NewClient() which reads from GetAPIKey().
	// Since we made KeyringService a variable, we can change it for the test and inject the key.
	originalService := config.KeyringService
	config.KeyringService = "plane-cli-cmd-test"
	t.Cleanup(func() { config.KeyringService = originalService; config.DeleteAPIKey() })

	err := config.SetAPIKey(apiKey)
	require.NoError(t, err)
}

func TestIssueSearchIntegration(t *testing.T) {
	setupTestEnvironment(t)

	// Test 'plane issue search "Test"'
	cmd := searchCmd

	// Redirect output to buffer (optional, but good for testing)
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	// execute the function directly instead of cobra execution to avoid os.Exit
	err := runSearch(cmd, []string{"Test"})
	assert.NoError(t, err)
}
