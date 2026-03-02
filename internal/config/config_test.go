package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitConfig(t *testing.T) {
	// Create temporary directory for test config
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".config", AppName)
	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	// Set config file to temp location
	SetConfigFile(filepath.Join(configDir, ConfigFileName+".yaml"))

	// Test initialization
	err = InitConfig()
	require.NoError(t, err)

	// Check defaults
	assert.Equal(t, "1.0", Cfg.Version)
	assert.Equal(t, "table", Cfg.OutputFormat)
	assert.Equal(t, DefaultAPIHost, Cfg.APIHost)
}

func TestSaveConfig(t *testing.T) {
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".config", AppName)
	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	SetConfigFile(filepath.Join(configDir, ConfigFileName+".yaml"))

	// Initialize
	err = InitConfig()
	require.NoError(t, err)

	// Modify config
	Cfg.DefaultWorkspace = "test-workspace"
	Cfg.DefaultProject = "test-project"
	Cfg.OutputFormat = "json"

	// Save
	err = SaveConfig()
	require.NoError(t, err)

	// Re-initialize and verify
	Cfg = Config{}
	err = InitConfig()
	require.NoError(t, err)

	assert.Equal(t, "test-workspace", Cfg.DefaultWorkspace)
	assert.Equal(t, "test-project", Cfg.DefaultProject)
	assert.Equal(t, "json", Cfg.OutputFormat)
}

func TestAPIKeyStorage(t *testing.T) {
	// Note: This test uses the actual keyring
	// In CI environments, this might fail without proper setup

	// Use a test service to avoid overwriting user credentials
	originalService := KeyringService
	KeyringService = "plane-cli-test"
	defer func() { KeyringService = originalService }()

	testKey := "test-api-key-12345"

	// Set API key
	err := SetAPIKey(testKey)
	if err != nil {
		t.Skip("Keyring not available in test environment:", err)
	}

	// Get API key
	retrievedKey, err := GetAPIKey()
	require.NoError(t, err)
	assert.Equal(t, testKey, retrievedKey)

	// Delete API key
	err = DeleteAPIKey()
	require.NoError(t, err)

	// Verify deletion
	_, err = GetAPIKey()
	assert.Error(t, err) // Should return error when key is deleted
}
