package epic

import (
	"fmt"

	"github.com/rohithmahesh3/plane-cli/internal/api"
	"github.com/rohithmahesh3/plane-cli/internal/config"
	"github.com/rohithmahesh3/plane-cli/internal/output"
	"github.com/spf13/cobra"
)

var EpicCmd = &cobra.Command{
	Use:     "epic",
	Aliases: []string{"epics"},
	Short:   "Manage epics",
	Long:    `List and view epics for organizing large bodies of work.`,
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List epics",
	Long:    `List all epics in the current project.`,
	RunE:    runList,
}

var viewCmd = &cobra.Command{
	Use:   "view <id>",
	Short: "View epic details",
	Long:  `Display detailed information about a specific epic.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runView,
}

func init() {
	EpicCmd.AddCommand(listCmd)
	EpicCmd.AddCommand(viewCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified. Use --project flag or set default project")
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	epics, err := client.ListEpics(projectID)
	if err != nil {
		return err
	}

	if len(epics) == 0 {
		output.Info("No epics found")
		return nil
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)

	type epicOutput struct {
		ID       string `table:"ID" json:"id"`
		Sequence int    `table:"#" json:"sequence_id"`
		Name     string `table:"NAME" json:"name"`
		State    string `table:"STATE" json:"state"`
		Priority string `table:"PRIORITY" json:"priority"`
	}

	var outputs []epicOutput
	for _, e := range epics {
		outputs = append(outputs, epicOutput{
			ID:       e.ID,
			Sequence: e.SequenceID,
			Name:     e.Name,
			State:    e.State,
			Priority: e.Priority,
		})
	}

	return formatter.Print(outputs)
}

func runView(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	epicID := args[0]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	epic, err := client.GetEpic(projectID, epicID)
	if err != nil {
		return err
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)
	return formatter.Print(epic)
}
