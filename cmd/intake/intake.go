package intake

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/rohithmahesh3/plane-cli/internal/api"
	"github.com/rohithmahesh3/plane-cli/internal/config"
	"github.com/rohithmahesh3/plane-cli/internal/output"
	"github.com/rohithmahesh3/plane-cli/pkg/plane"
	"github.com/spf13/cobra"
)

var (
	intakeName     string
	intakePriority string
)

var IntakeCmd = &cobra.Command{
	Use:     "intake",
	Aliases: []string{"inbox", "requests"},
	Short:   "Manage intake issues",
	Long:    `List, create, and manage intake requests for your project.`,
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List intake issues",
	Long:    `List all intake issues in the current project.`,
	RunE:    runList,
}

var viewCmd = &cobra.Command{
	Use:   "view <id>",
	Short: "View intake issue details",
	Long:  `Display detailed information about a specific intake issue.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runView,
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new intake issue",
	Long: `Create a new intake issue in the current project.

Examples:
  plane intake create --name "Feature Request" --priority high
  plane intake create -n "Bug Report" -p urgent`,
	RunE: runCreate,
}

var deleteCmd = &cobra.Command{
	Use:     "delete <id>",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete an intake issue",
	Long:    `Delete an intake issue from the project.`,
	Args:    cobra.ExactArgs(1),
	RunE:    runDelete,
}

func init() {
	IntakeCmd.AddCommand(listCmd)
	IntakeCmd.AddCommand(viewCmd)
	IntakeCmd.AddCommand(createCmd)
	IntakeCmd.AddCommand(deleteCmd)

	// Create flags
	createCmd.Flags().StringVarP(&intakeName, "name", "n", "", "Issue name/title")
	createCmd.Flags().StringVarP(&intakePriority, "priority", "p", "medium", "Issue priority (low, medium, high, urgent)")
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

	intakeIssues, err := client.ListIntakeIssues(projectID)
	if err != nil {
		return err
	}

	if len(intakeIssues) == 0 {
		output.Info("No intake issues found")
		return nil
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)

	type intakeOutput struct {
		ID     string `table:"ID" json:"id"`
		Status string `table:"STATUS" json:"status_formatted"`
		Source string `table:"SOURCE" json:"source,omitempty"`
		Issue  string `table:"ISSUE" json:"issue,omitempty"`
	}

	var outputs []intakeOutput
	for _, i := range intakeIssues {
		outputs = append(outputs, intakeOutput{
			ID:     i.ID,
			Status: formatIntakeStatus(i.Status),
			Source: i.Source,
			Issue:  i.Issue,
		})
	}

	return formatter.Print(outputs)
}

func runView(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	intakeID := args[0]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	intake, err := client.GetIntakeIssue(projectID, intakeID)
	if err != nil {
		return err
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)
	return formatter.Print(intake)
}

func runCreate(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified. Use --project flag or set default project")
	}

	// Interactive prompts if flags not provided
	if intakeName == "" {
		prompt := &survey.Input{
			Message: "Issue name/title:",
			Help:    "Brief description of the request",
		}
		if err := survey.AskOne(prompt, &intakeName); err != nil {
			return err
		}
	}

	if intakeName == "" {
		return fmt.Errorf("issue name is required")
	}

	if intakePriority == "" {
		priorityOptions := []string{"low", "medium", "high", "urgent"}
		prompt := &survey.Select{
			Message: "Priority:",
			Options: priorityOptions,
			Default: "medium",
		}
		if err := survey.AskOne(prompt, &intakePriority); err != nil {
			return err
		}
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	req := plane.CreateIntakeIssueRequest{}
	req.Issue.Name = intakeName
	req.Issue.Priority = intakePriority

	intake, err := client.CreateIntakeIssue(projectID, req)
	if err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Created intake issue (%s)", intake.ID))
	return nil
}

func runDelete(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	intakeID := args[0]

	// Confirm deletion
	var confirm bool
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("Are you sure you want to delete intake issue %s?", intakeID),
		Default: false,
	}
	if err := survey.AskOne(prompt, &confirm); err != nil {
		return err
	}

	if !confirm {
		output.Info("Deletion cancelled")
		return nil
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	if err := client.DeleteIntakeIssue(projectID, intakeID); err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Deleted intake issue %s", intakeID))
	return nil
}

// formatIntakeStatus converts status code to human-readable string
func formatIntakeStatus(status int) string {
	switch status {
	case -2:
		return "Pending"
	case -1:
		return "Rejected"
	case 0:
		return "Snoozed"
	case 1:
		return "Accepted"
	case 2:
		return "Duplicate"
	default:
		return fmt.Sprintf("Unknown (%d)", status)
	}
}
