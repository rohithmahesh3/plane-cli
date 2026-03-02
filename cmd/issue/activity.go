package issue

import (
	"fmt"

	"github.com/rohithmahesh3/plane-cli/internal/api"
	"github.com/rohithmahesh3/plane-cli/internal/config"
	"github.com/rohithmahesh3/plane-cli/internal/output"
	"github.com/spf13/cobra"
)

func init() {
	// Add activity subcommand
	activityCmd := &cobra.Command{
		Use:     "activity",
		Aliases: []string{"history", "log"},
		Short:   "View issue activity history",
		Long:    `View the activity history and audit trail for a work item.`,
	}

	activityListCmd := &cobra.Command{
		Use:     "list <issue-id>",
		Aliases: []string{"ls"},
		Short:   "List activity history",
		Long:    `Display the activity history for a specific issue.`,
		Args:    cobra.ExactArgs(1),
		RunE:    runActivityList,
	}

	activityViewCmd := &cobra.Command{
		Use:   "view <issue-id> <activity-id>",
		Short: "View activity details",
		Long:  `Display detailed information about a specific activity.`,
		Args:  cobra.ExactArgs(2),
		RunE:  runActivityView,
	}

	activityCmd.AddCommand(activityListCmd)
	activityCmd.AddCommand(activityViewCmd)

	IssueCmd.AddCommand(activityCmd)
}

func runActivityList(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	issueID := args[0]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	issueID, err = resolveIssueID(client, projectID, issueID)
	if err != nil {
		return err
	}

	activities, err := client.ListActivities(projectID, issueID)
	if err != nil {
		return err
	}

	if len(activities) == 0 {
		output.Info("No activity found for this issue")
		return nil
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)

	type activityOutput struct {
		ID        string `table:"ID" json:"id"`
		Verb      string `table:"ACTION" json:"verb"`
		Field     string `table:"FIELD" json:"field,omitempty"`
		OldValue  string `table:"OLD" json:"old_value,omitempty"`
		NewValue  string `table:"NEW" json:"new_value,omitempty"`
		Comment   string `table:"COMMENT" json:"comment,omitempty"`
		Actor     string `table:"ACTOR" json:"actor,omitempty"`
		CreatedAt string `table:"WHEN" json:"created_at_formatted"`
	}

	var outputs []activityOutput
	for _, a := range activities {
		outputs = append(outputs, activityOutput{
			ID:        a.ID,
			Verb:      a.Verb,
			Field:     a.Field,
			OldValue:  truncateString(a.OldValue, 20),
			NewValue:  truncateString(a.NewValue, 20),
			Comment:   truncateString(a.Comment, 30),
			Actor:     a.Actor,
			CreatedAt: a.CreatedAt.Format("2006-01-02 15:04"),
		})
	}

	return formatter.Print(outputs)
}

func runActivityView(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	issueID := args[0]
	activityID := args[1]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	issueID, err = resolveIssueID(client, projectID, issueID)
	if err != nil {
		return err
	}

	activity, err := client.GetActivity(projectID, issueID, activityID)
	if err != nil {
		return err
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)
	return formatter.Print(activity)
}
