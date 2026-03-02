package cycle

import (
	"github.com/rohithmahesh3/plane-cli/internal/output"
	"github.com/spf13/cobra"
)

var CycleCmd = &cobra.Command{
	Use:     "cycle",
	Aliases: []string{"sprint"},
	Short:   "Manage cycles (sprints)",
	Long:    `List and manage Plane cycles (sprints).`,
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List cycles",
	RunE: func(cmd *cobra.Command, args []string) error {
		output.Info("Cycle listing - to be implemented")
		return nil
	},
}

func init() {
	CycleCmd.AddCommand(listCmd)
}
