package cycle

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
	cycleName        string
	cycleDescription string
	cycleStartDate   string
	cycleEndDate     string
	showArchived     bool
)

var CycleCmd = &cobra.Command{
	Use:     "cycle",
	Aliases: []string{"sprint"},
	Short:   "Manage cycles (sprints)",
	Long:    `List, create, edit, and manage Plane cycles (sprints).`,
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List cycles",
	Long:    `List all cycles in the current project.`,
	RunE:    runList,
}

var viewCmd = &cobra.Command{
	Use:   "view <id>",
	Short: "View cycle details",
	Long:  `Display detailed information about a specific cycle.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runView,
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new cycle",
	Long: `Create a new cycle (sprint) in the current project.

Examples:
  plane cycle create --name "Sprint 1" --start-date 2024-01-01 --end-date 2024-01-14
  plane cycle create -n "Q1 Planning" -d "First quarter planning cycle"`,
	RunE: runCreate,
}

var editCmd = &cobra.Command{
	Use:   "edit <id>",
	Short: "Edit a cycle",
	Long:  `Edit an existing cycle.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runEdit,
}

var deleteCmd = &cobra.Command{
	Use:     "delete <id>",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete a cycle",
	Long:    `Delete a cycle from the project.`,
	Args:    cobra.ExactArgs(1),
	RunE:    runDelete,
}

var archiveCmd = &cobra.Command{
	Use:   "archive <id>",
	Short: "Archive a cycle",
	Long:  `Archive a cycle to hide it from active cycles.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runArchive,
}

var unarchiveCmd = &cobra.Command{
	Use:   "unarchive <id>",
	Short: "Unarchive a cycle",
	Long:  `Restore an archived cycle to active status.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runUnarchive,
}

var issuesCmd = &cobra.Command{
	Use:     "issues <id>",
	Aliases: []string{"work-items"},
	Short:   "List cycle issues",
	Long:    `List all work items in a specific cycle.`,
	Args:    cobra.ExactArgs(1),
	RunE:    runIssues,
}

var addIssuesCmd = &cobra.Command{
	Use:   "add-issues <cycle-id> <issue-ids...\u003e",
	Short: "Add issues to cycle",
	Long:  `Add work items to a cycle.`,
	Args:  cobra.MinimumNArgs(2),
	RunE:  runAddIssues,
}

var removeIssueCmd = &cobra.Command{
	Use:   "remove-issue <cycle-id> <issue-id>",
	Short: "Remove issue from cycle",
	Long:  `Remove a work item from a cycle.`,
	Args:  cobra.ExactArgs(2),
	RunE:  runRemoveIssue,
}

func init() {
	CycleCmd.AddCommand(listCmd)
	CycleCmd.AddCommand(viewCmd)
	CycleCmd.AddCommand(createCmd)
	CycleCmd.AddCommand(editCmd)
	CycleCmd.AddCommand(deleteCmd)
	CycleCmd.AddCommand(archiveCmd)
	CycleCmd.AddCommand(unarchiveCmd)
	CycleCmd.AddCommand(issuesCmd)
	CycleCmd.AddCommand(addIssuesCmd)
	CycleCmd.AddCommand(removeIssueCmd)

	// List flags
	listCmd.Flags().BoolVarP(&showArchived, "archived", "a", false, "Show archived cycles")

	// Create flags
	createCmd.Flags().StringVarP(&cycleName, "name", "n", "", "Cycle name")
	createCmd.Flags().StringVarP(&cycleDescription, "description", "d", "", "Cycle description")
	createCmd.Flags().StringVarP(&cycleStartDate, "start-date", "s", "", "Start date (YYYY-MM-DD)")
	createCmd.Flags().StringVarP(&cycleEndDate, "end-date", "e", "", "End date (YYYY-MM-DD)")

	// Edit flags
	editCmd.Flags().StringVarP(&cycleName, "name", "n", "", "New cycle name")
	editCmd.Flags().StringVarP(&cycleDescription, "description", "d", "", "New cycle description")
	editCmd.Flags().StringVarP(&cycleStartDate, "start-date", "s", "", "New start date (YYYY-MM-DD)")
	editCmd.Flags().StringVarP(&cycleEndDate, "end-date", "e", "", "New end date (YYYY-MM-DD)")
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

	cycles, err := client.ListCycles(projectID, showArchived)
	if err != nil {
		return err
	}

	if len(cycles) == 0 {
		output.Info("No cycles found")
		return nil
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)

	type cycleOutput struct {
		ID          string `table:"ID" json:"id"`
		Name        string `table:"NAME" json:"name"`
		StartDate   string `table:"START" json:"start_date,omitempty"`
		EndDate     string `table:"END" json:"end_date,omitempty"`
		Status      string `table:"STATUS" json:"status,omitempty"`
		Description string `table:"DESCRIPTION" json:"description,omitempty"`
	}

	var outputs []cycleOutput
	for _, c := range cycles {
		outputs = append(outputs, cycleOutput{
			ID:          c.ID,
			Name:        c.Name,
			StartDate:   c.StartDate,
			EndDate:     c.EndDate,
			Status:      c.Status,
			Description: c.Description,
		})
	}

	return formatter.Print(outputs)
}

func runView(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	cycleID := args[0]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	cycle, err := client.GetCycle(projectID, cycleID)
	if err != nil {
		return err
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)
	return formatter.Print(cycle)
}

func runCreate(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified. Use --project flag or set default project")
	}

	// Interactive prompts if flags not provided
	if cycleName == "" {
		prompt := &survey.Input{
			Message: "Cycle name:",
			Help:    "e.g., Sprint 1, Q1 Planning",
		}
		if err := survey.AskOne(prompt, &cycleName); err != nil {
			return err
		}
	}

	if cycleName == "" {
		return fmt.Errorf("cycle name is required")
	}

	if cycleDescription == "" {
		prompt := &survey.Input{
			Message: "Description (optional):",
		}
		_ = survey.AskOne(prompt, &cycleDescription)
	}

	if cycleStartDate == "" {
		prompt := &survey.Input{
			Message: "Start date (YYYY-MM-DD):",
			Help:    "When does this cycle start?",
		}
		_ = survey.AskOne(prompt, &cycleStartDate)
	}

	if cycleEndDate == "" {
		prompt := &survey.Input{
			Message: "End date (YYYY-MM-DD):",
			Help:    "When does this cycle end?",
		}
		_ = survey.AskOne(prompt, &cycleEndDate)
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	req := plane.CreateCycleRequest{
		Name:        cycleName,
		Description: cycleDescription,
		StartDate:   cycleStartDate,
		EndDate:     cycleEndDate,
	}

	cycle, err := client.CreateCycle(projectID, req)
	if err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Created cycle '%s' (%s)", cycle.Name, cycle.ID))
	return nil
}

func runEdit(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	cycleID := args[0]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	// Get current cycle
	cycle, err := client.GetCycle(projectID, cycleID)
	if err != nil {
		return err
	}

	req := plane.UpdateCycleRequest{}

	// Interactive mode if no flags provided
	if cycleName == "" && cycleDescription == "" && cycleStartDate == "" && cycleEndDate == "" {
		output.Info(fmt.Sprintf("Editing cycle: %s", cycle.Name))

		prompt := &survey.Input{
			Message: "Name:",
			Default: cycle.Name,
		}
		if err := survey.AskOne(prompt, &req.Name); err != nil {
			return err
		}

		descPrompt := &survey.Input{
			Message: "Description:",
			Default: cycle.Description,
		}
		if err := survey.AskOne(descPrompt, &req.Description); err != nil {
			return err
		}

		startPrompt := &survey.Input{
			Message: "Start date (YYYY-MM-DD):",
			Default: cycle.StartDate,
		}
		if err := survey.AskOne(startPrompt, &req.StartDate); err != nil {
			return err
		}

		endPrompt := &survey.Input{
			Message: "End date (YYYY-MM-DD):",
			Default: cycle.EndDate,
		}
		if err := survey.AskOne(endPrompt, &req.EndDate); err != nil {
			return err
		}
	} else {
		// Use provided flags
		if cycleName != "" {
			req.Name = cycleName
		}
		if cycleDescription != "" {
			req.Description = cycleDescription
		}
		if cycleStartDate != "" {
			req.StartDate = cycleStartDate
		}
		if cycleEndDate != "" {
			req.EndDate = cycleEndDate
		}
	}

	updatedCycle, err := client.UpdateCycle(projectID, cycleID, req)
	if err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Updated cycle '%s'", updatedCycle.Name))
	return nil
}

func runDelete(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	cycleID := args[0]

	// Confirm deletion
	var confirm bool
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("Are you sure you want to delete cycle %s?", cycleID),
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

	if err := client.DeleteCycle(projectID, cycleID); err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Deleted cycle %s", cycleID))
	return nil
}

func runArchive(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	cycleID := args[0]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	if err := client.ArchiveCycle(projectID, cycleID); err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Archived cycle %s", cycleID))
	return nil
}

func runUnarchive(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	cycleID := args[0]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	if err := client.UnarchiveCycle(projectID, cycleID); err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Unarchived cycle %s", cycleID))
	return nil
}

func runIssues(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	cycleID := args[0]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	issues, err := client.ListCycleIssues(projectID, cycleID)
	if err != nil {
		return err
	}

	if len(issues) == 0 {
		output.Info("No issues found in this cycle")
		return nil
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)

	type issueOutput struct {
		ID       string `table:"ID" json:"id"`
		Sequence int    `table:"#" json:"sequence_id"`
		Title    string `table:"TITLE" json:"title"`
		State    string `table:"STATE" json:"state_name"`
		Priority string `table:"PRIORITY" json:"priority"`
	}

	var outputs []issueOutput
	for _, issue := range issues {
		outputs = append(outputs, issueOutput{
			ID:       issue.ID,
			Sequence: issue.SequenceID,
			Title:    issue.Name,
			State:    issue.State,
			Priority: issue.Priority,
		})
	}

	return formatter.Print(outputs)
}

func runAddIssues(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	cycleID := args[0]
	issueIDs := args[1:]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	if err := client.AddIssuesToCycle(projectID, cycleID, issueIDs); err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Added %d issue(s) to cycle %s", len(issueIDs), cycleID))
	return nil
}

func runRemoveIssue(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	cycleID := args[0]
	issueID := args[1]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	if err := client.RemoveIssueFromCycle(projectID, cycleID, issueID); err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Removed issue %s from cycle %s", issueID, cycleID))
	return nil
}
