package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
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
