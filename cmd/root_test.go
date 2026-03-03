package cmd

import (
	"os"
	"path/filepath"
	"testing"

	cfg "github.com/rohithmahesh3/plane-cli/internal/config"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldAllowInvalidOutputConfig(t *testing.T) {
	root := &cobra.Command{Use: "plane-cli"}
	configCmd := &cobra.Command{Use: "config"}
	setCmd := &cobra.Command{Use: "set"}
	configCmd.AddCommand(setCmd)
	root.AddCommand(configCmd)

	assert.True(t, shouldAllowInvalidOutputConfig(setCmd, []string{"output", "yaml"}))
	assert.False(t, shouldAllowInvalidOutputConfig(setCmd, []string{"workspace", "foo"}))

	otherCmd := &cobra.Command{Use: "issue"}
	root.AddCommand(otherCmd)
	assert.False(t, shouldAllowInvalidOutputConfig(otherCmd, []string{"output", "yaml"}))
}

func TestPersistentPreRunEAllowsOutputRecovery(t *testing.T) {
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".config", cfg.AppName)
	err := os.MkdirAll(configDir, 0o755)
	require.NoError(t, err)

	configPath := filepath.Join(configDir, cfg.ConfigFileName+".yaml")
	err = os.WriteFile(configPath, []byte("output_format: table\n"), 0o644)
	require.NoError(t, err)

	originalConfigFile := configFile
	originalOutputFmt := outputFmt
	configFile = configPath
	outputFmt = ""
	t.Cleanup(func() {
		configFile = originalConfigFile
		outputFmt = originalOutputFmt
	})

	cmd, _, err := rootCmd.Find([]string{"config", "set"})
	require.NoError(t, err)
	require.NotNil(t, cmd)

	err = rootCmd.PersistentPreRunE(cmd, []string{"output", "yaml"})
	require.NoError(t, err)
	assert.Equal(t, "table", cfg.Cfg.OutputFormat)
}
