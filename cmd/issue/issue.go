package issue

import (
	"fmt"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/rohithmahesh3/plane-cli/internal/api"
	"github.com/rohithmahesh3/plane-cli/internal/config"
	"github.com/rohithmahesh3/plane-cli/internal/output"
	"github.com/rohithmahesh3/plane-cli/pkg/plane"
	"github.com/spf13/cobra"
)

var (
	stateFilter    string
	priorityFilter string
	assigneeFilter string
	labelFilter    string
	cycleFilter    string
	moduleFilter   string
	perPage        int
	allFlag        bool

	issueTitle       string
	issueDescription string
	issuePriority    string
	issueState       string
	issueAssignees   []string
	issueLabels      []string
)

var IssueCmd = &cobra.Command{
	Use:     "issue",
	Aliases: []string{"i", "issues", "ticket"},
	Short:   "Manage issues (work items)",
	Long:    `List, create, edit, and manage Plane issues/work items.`,
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List issues",
	Long: `List issues in the current project.

Examples:
  plane issue list
  plane issue list --state backlog
  plane issue list --priority high --assignee @alice
  plane issue list --label bug --cycle "Sprint 1"`,
	RunE: runList,
}

var viewCmd = &cobra.Command{
	Use:   "view <id>",
	Short: "View issue details",
	Long:  `Display detailed information about a specific issue.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runView,
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new issue",
	Long: `Create a new issue in the current project.

Examples:
  plane issue create --title "Bug fix" --priority high
  plane issue create -t "Feature request" -d "Description here"`,
	RunE: runCreate,
}

var editCmd = &cobra.Command{
	Use:   "edit <id>",
	Short: "Edit an issue",
	Long:  `Edit an existing issue.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runEdit,
}

var deleteCmd = &cobra.Command{
	Use:     "delete <id>",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete an issue",
	Long:    `Delete an issue from the project.`,
	Args:    cobra.ExactArgs(1),
	RunE:    runDelete,
}

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for issues",
	Long:  `Search for issues across the workspace.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runSearch,
}

func init() {
	IssueCmd.AddCommand(listCmd)
	IssueCmd.AddCommand(viewCmd)
	IssueCmd.AddCommand(createCmd)
	IssueCmd.AddCommand(editCmd)
	IssueCmd.AddCommand(deleteCmd)
	IssueCmd.AddCommand(searchCmd)

	// List flags
	listCmd.Flags().StringVarP(&stateFilter, "state", "s", "", "Filter by state (backlog, todo, in-progress, done)")
	listCmd.Flags().StringVar(&priorityFilter, "priority", "", "Filter by priority (low, medium, high, urgent)")
	listCmd.Flags().StringVar(&assigneeFilter, "assignee", "", "Filter by assignee (@username or 'me')")
	listCmd.Flags().StringVar(&labelFilter, "label", "", "Filter by label")
	listCmd.Flags().StringVar(&cycleFilter, "cycle", "", "Filter by cycle name")
	listCmd.Flags().StringVar(&moduleFilter, "module", "", "Filter by module name")
	listCmd.Flags().IntVarP(&perPage, "limit", "l", 20, "Number of issues to show per page")
	listCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "Show all issues (fetch all pages)")

	// Create flags
	createCmd.Flags().StringVarP(&issueTitle, "title", "t", "", "Issue title")
	createCmd.Flags().StringVarP(&issueDescription, "description", "d", "", "Issue description")
	createCmd.Flags().StringVarP(&issuePriority, "priority", "p", "medium", "Issue priority (low, medium, high, urgent)")
	createCmd.Flags().StringSliceVarP(&issueAssignees, "assignee", "a", nil, "Assignee(s) (@username)")
	createCmd.Flags().StringSliceVar(&issueLabels, "label", nil, "Label(s)")

	// Edit flags
	editCmd.Flags().StringVarP(&issueTitle, "title", "t", "", "New title")
	editCmd.Flags().StringVarP(&issueDescription, "description", "d", "", "New description")
	editCmd.Flags().StringVar(&issuePriority, "priority", "", "New priority (low, medium, high, urgent)")
	editCmd.Flags().StringVar(&issueState, "state", "", "New state (backlog, todo, in-progress, done)")
	editCmd.Flags().StringSliceVarP(&issueAssignees, "assignee", "a", nil, "New assignee(s) (@username)")
	editCmd.Flags().StringSliceVar(&issueLabels, "label", nil, "New label(s)")
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

	opts := api.IssueListOptions{
		State:    stateFilter,
		Priority: priorityFilter,
		Assignee: assigneeFilter,
		Label:    labelFilter,
		Cycle:    cycleFilter,
		Module:   moduleFilter,
		PerPage:  perPage,
	}

	issues, _, err := client.ListIssues(projectID, opts)
	if err != nil {
		return err
	}

	if len(issues) == 0 {
		output.Info("No issues found")
		return nil
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)

	type issueOutput struct {
		ID       string `table:"ID" json:"id"`
		Sequence int    `table:"#" json:"sequence_id"`
		Title    string `table:"TITLE" json:"title"`
		State    string `table:"STATE" json:"state_name"`
		Priority string `table:"PRIORITY" json:"priority"`
		Assignee string `table:"ASSIGNEE" json:"assignee"`
	}

	var outputs []issueOutput
	for _, issue := range issues {
		assignee := "-"
		if len(issue.Assignees) > 0 {
			u := issue.Assignees[0]
			if u.DisplayName != "" {
				assignee = "@" + u.DisplayName
			} else if u.Email != "" {
				assignee = u.Email
			}
		}

		stateName := issue.State.Name
		if stateName == "" {
			stateName = "-"
		}

		outputs = append(outputs, issueOutput{
			ID:       issue.ID,
			Sequence: issue.SequenceID,
			Title:    issue.Name,
			State:    stateName,
			Priority: issue.Priority,
			Assignee: assignee,
		})
	}

	return formatter.Print(outputs)
}

func runView(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	issueID := args[0]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	// Try to parse as sequence ID (number)
	var issue *plane.Issue
	if seqID, err := strconv.Atoi(issueID); err == nil {
		issue, err = client.GetIssueByIdentifier(projectID, seqID)
		if err != nil {
			return err
		}
	} else {
		issue, err = client.GetIssue(projectID, issueID)
		if err != nil {
			return err
		}
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)
	return formatter.Print(issue)
}

func runCreate(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified. Use --project flag or set default project")
	}

	// Interactive prompts if flags not provided
	if issueTitle == "" {
		prompt := &survey.Input{
			Message: "Issue title:",
		}
		if err := survey.AskOne(prompt, &issueTitle); err != nil {
			return err
		}
	}

	if issueTitle == "" {
		return fmt.Errorf("issue title is required")
	}

	if issueDescription == "" {
		prompt := &survey.Editor{
			Message:       "Issue description:",
			FileName:      "*.md",
			HideDefault:   true,
			AppendDefault: true,
		}
		if err := survey.AskOne(prompt, &issueDescription); err != nil {
			return err
		}
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	resolvedAssignees, err := client.ResolveAssignees(projectID, issueAssignees)
	if err != nil {
		return fmt.Errorf("failed to resolve assignees: %w", err)
	}

	resolvedLabels, err := client.ResolveLabels(projectID, issueLabels)
	if err != nil {
		return fmt.Errorf("failed to resolve labels: %w", err)
	}

	req := plane.CreateIssueRequest{
		Name:        issueTitle,
		Description: issueDescription,
		Priority:    issuePriority,
		Assignees:   resolvedAssignees,
		Labels:      resolvedLabels,
	}

	issue, err := client.CreateIssue(projectID, req)
	if err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Created issue %s-%d", config.Cfg.DefaultProject[:4], issue.SequenceID))
	return nil
}

func runEdit(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	issueID := args[0]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	// Get current issue
	issue, err := client.GetIssue(projectID, issueID)
	if err != nil {
		return err
	}

	req := plane.UpdateIssueRequest{}

	hasFlags := issueTitle != "" || issueDescription != "" || issuePriority != "" || issueState != "" || len(issueAssignees) > 0 || len(issueLabels) > 0

	// Interactive mode if no flags provided
	if !hasFlags {
		// Show current values and prompt for changes
		output.Info(fmt.Sprintf("Editing issue %d: %s", issue.SequenceID, issue.Name))

		prompt := &survey.Input{
			Message: "Title:",
			Default: issue.Name,
		}
		if err := survey.AskOne(prompt, &req.Name); err != nil {
			return err
		}

		priorityOptions := []string{"low", "medium", "high", "urgent"}
		priorityPrompt := &survey.Select{
			Message: "Priority:",
			Options: priorityOptions,
			Default: issue.Priority,
		}
		if err := survey.AskOne(priorityPrompt, &req.Priority); err != nil {
			return err
		}
	} else {
		// Use provided flags
		if issueTitle != "" {
			req.Name = issueTitle
		}
		if issueDescription != "" {
			req.Description = issueDescription
		}
		if issuePriority != "" {
			req.Priority = issuePriority
		}
		if issueState != "" {
			resolvedState, err := client.ResolveState(projectID, issueState)
			if err != nil {
				return fmt.Errorf("failed to resolve state: %w", err)
			}
			req.State = resolvedState
		}
		if len(issueAssignees) > 0 {
			resolvedAssignees, err := client.ResolveAssignees(projectID, issueAssignees)
			if err != nil {
				return fmt.Errorf("failed to resolve assignees: %w", err)
			}
			req.Assignees = resolvedAssignees
		}
		if len(issueLabels) > 0 {
			resolvedLabels, err := client.ResolveLabels(projectID, issueLabels)
			if err != nil {
				return fmt.Errorf("failed to resolve labels: %w", err)
			}
			req.Labels = resolvedLabels
		}
	}

	updatedIssue, err := client.UpdateIssue(projectID, issueID, req)
	if err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Updated issue %d", updatedIssue.SequenceID))
	return nil
}

func runDelete(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	issueID := args[0]

	// Confirm deletion
	var confirm bool
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("Are you sure you want to delete issue %s?", issueID),
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

	if err := client.DeleteIssue(projectID, issueID); err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Deleted issue %s", issueID))
	return nil
}

func runSearch(cmd *cobra.Command, args []string) error {
	query := args[0]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	issues, err := client.SearchIssues(query)
	if err != nil {
		return err
	}

	if len(issues) == 0 {
		output.Info("No issues found")
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
		stateName := issue.State.Name
		if stateName == "" {
			stateName = "-"
		}

		outputs = append(outputs, issueOutput{
			ID:       issue.ID,
			Sequence: issue.SequenceID,
			Title:    issue.Name,
			State:    stateName,
			Priority: issue.Priority,
		})
	}

	return formatter.Print(outputs)
}
