//go:build integration
// +build integration

package workspace

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/rohithmahesh3/plane-cli/internal/api"
	"github.com/rohithmahesh3/plane-cli/internal/config"
	"github.com/rohithmahesh3/plane-cli/internal/integrationtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestEnvironment(t *testing.T) {
	integrationtest.WaitForSlot(t)

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

func TestWorkspaceMembersSearchIntegration(t *testing.T) {
	setupTestEnvironment(t)

	client, err := api.NewClient()
	require.NoError(t, err)

	members, err := client.GetWorkspaceMembers()
	require.NoError(t, err)
	require.NotEmpty(t, members)

	target := members[0]
	query := target.ID
	if query == "" {
		query = target.Email
	}
	if query == "" {
		query = target.DisplayName
	}
	require.NotEmpty(t, query)

	memberSearch = query
	memberExact = true
	memberLimit = 0
	t.Cleanup(func() {
		memberSearch = ""
		memberExact = false
		memberLimit = 0
	})

	output := captureStdout(t, func() {
		err = runMembers(membersCmd, []string{})
	})

	require.NoError(t, err)

	var results []struct {
		ID string `json:"id"`
	}
	require.NoError(t, json.Unmarshal([]byte(output), &results))
	require.Len(t, results, 1)
	assert.Equal(t, target.ID, results[0].ID)
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	oldStdout := os.Stdout
	reader, writer, err := os.Pipe()
	require.NoError(t, err)

	os.Stdout = writer
	t.Cleanup(func() {
		os.Stdout = oldStdout
	})

	fn()

	require.NoError(t, writer.Close())
	output, err := io.ReadAll(reader)
	require.NoError(t, err)
	require.NoError(t, reader.Close())

	os.Stdout = oldStdout
	return string(output)
}
