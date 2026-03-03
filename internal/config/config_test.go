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
	assert.Equal(t, "yaml", Cfg.OutputFormat)
	assert.Equal(t, DefaultAPIHost, Cfg.APIHost)
}

func TestSaveConfig(t *testing.T) {
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".config", AppName)
	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	// Change to temp directory to avoid loading local .plane/settings.yaml
	originalWd, _ := os.Getwd()
	defer func() {
		_ = os.Chdir(originalWd)
	}()
	err = os.Chdir(tempDir)
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

func TestInitConfigRejectsTableOutputFormat(t *testing.T) {
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".config", AppName)
	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	configPath := filepath.Join(configDir, ConfigFileName+".yaml")
	err = os.WriteFile(configPath, []byte("output_format: table\n"), 0644)
	require.NoError(t, err)

	SetConfigFile(configPath)

	err = InitConfig()
	require.Error(t, err)
	assert.EqualError(t, err, `invalid output format "table": table output has been removed; supported formats are json, yaml`)
}

func TestInitConfigAllowInvalidOutputAcceptsTable(t *testing.T) {
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".config", AppName)
	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	configPath := filepath.Join(configDir, ConfigFileName+".yaml")
	err = os.WriteFile(configPath, []byte("output_format: table\n"), 0644)
	require.NoError(t, err)

	SetConfigFile(configPath)

	err = InitConfigAllowInvalidOutput()
	require.NoError(t, err)
	assert.Equal(t, "table", Cfg.OutputFormat)
}

func TestLocalConfig(t *testing.T) {
	tempDir := t.TempDir()

	// Create global config
	configDir := filepath.Join(tempDir, ".config", AppName)
	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	SetConfigFile(filepath.Join(configDir, ConfigFileName+".yaml"))

	// Setup initial state
	Cfg = Config{
		DefaultWorkspace: "global-workspace",
		DefaultProject:   "global-project",
	}

	// Change to a temp directory to simulate a project
	projectDir := filepath.Join(tempDir, "my-project")
	err = os.MkdirAll(projectDir, 0755)
	require.NoError(t, err)

	// Save original working dir and restore later
	originalWd, _ := os.Getwd()
	defer func() {
		_ = os.Chdir(originalWd)
	}()

	err = os.Chdir(projectDir)
	require.NoError(t, err)

	// 1. Test without .plane/settings.yaml
	err = InitConfig()
	require.NoError(t, err)
	// It should retain global defaults if no file (actually, InitConfig will overwrite Cfg, but if we save it, it will read it. Let's just test loadLocalConfig directly or rely on InitConfig to fall back to nothing)

	// 2. Test with .plane/settings.yaml
	planeDir := filepath.Join(projectDir, ".plane")
	err = os.MkdirAll(planeDir, 0755)
	require.NoError(t, err)

	settingsContent := []byte(`
workspace: local-workspace
project: local-project
`)
	err = os.WriteFile(filepath.Join(planeDir, "settings.yaml"), settingsContent, 0644)
	require.NoError(t, err)

	// We can manually reset Cfg to defaults
	Cfg.DefaultWorkspace = "global-workspace"
	Cfg.DefaultProject = "global-project"

	// Call loadLocalConfig directly since InitConfig overrides Cfg entirely with Viper defaults
	loadLocalConfig()

	assert.Equal(t, "local-workspace", Cfg.DefaultWorkspace)
	assert.Equal(t, "local-project", Cfg.DefaultProject)
}
