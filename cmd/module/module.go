package module

import (
	"github.com/rohithmahesh3/plane-cli/internal/output"
	"github.com/spf13/cobra"
)

var ModuleCmd = &cobra.Command{
	Use:     "module",
	Aliases: []string{"mod"},
	Short:   "Manage modules",
	Long:    `List and manage Plane modules.`,
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List modules",
	RunE: func(cmd *cobra.Command, args []string) error {
		output.Info("Module listing - to be implemented")
		return nil
	},
}

func init() {
	ModuleCmd.AddCommand(listCmd)
}
