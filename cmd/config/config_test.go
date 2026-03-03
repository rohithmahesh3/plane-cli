package config

import (
	"os"
	"path/filepath"
	"testing"

	internalconfig "github.com/rohithmahesh3/plane-cli/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunSetOutputAcceptsStructuredFormats(t *testing.T) {
	for _, format := range []string{"yaml", "json"} {
		t.Run(format, func(t *testing.T) {
			setupConfigTest(t)

			err := runSet(nil, []string{"output", format})
			require.NoError(t, err)
			assert.Equal(t, format, internalconfig.Cfg.OutputFormat)
		})
	}
}

func TestRunSetOutputRejectsTable(t *testing.T) {
	setupConfigTest(t)
	initial := internalconfig.Cfg.OutputFormat

	err := runSet(nil, []string{"output", "table"})
	require.Error(t, err)
	assert.EqualError(t, err, `invalid output format "table": table output has been removed; supported formats are json, yaml`)
	assert.Equal(t, initial, internalconfig.Cfg.OutputFormat)
}

func setupConfigTest(t *testing.T) {
	t.Helper()

	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".config", internalconfig.AppName)
	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	originalWd, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = os.Chdir(originalWd)
	})

	err = os.Chdir(tempDir)
	require.NoError(t, err)

	internalconfig.SetConfigFile(filepath.Join(configDir, internalconfig.ConfigFileName+".yaml"))
	internalconfig.Cfg = internalconfig.Config{}

	err = internalconfig.InitConfig()
	require.NoError(t, err)
}
