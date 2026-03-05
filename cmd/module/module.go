package module

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
	moduleName        string
	moduleDescription string
	moduleStatus      string
	showArchived      bool
)

var ModuleCmd = &cobra.Command{
	Use:     "module",
	Aliases: []string{"mod"},
	Short:   "Manage modules",
	Long:    `List, create, edit, and manage Plane modules.`,
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List modules",
	Long:    `List all modules in the current project.`,
	RunE:    runList,
}

var viewCmd = &cobra.Command{
	Use:   "view <id>",
	Short: "View module details",
	Long:  `Display detailed information about a specific module.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runView,
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new module",
	Long: `Create a new module in the current project.

Examples:
  plane module create --name "Authentication" --description "User auth features"
  plane module create -n "API Integration" -s "in-progress"`,
	RunE: runCreate,
}

var editCmd = &cobra.Command{
	Use:   "edit <id>",
	Short: "Edit a module",
	Long:  `Edit an existing module.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runEdit,
}

var deleteCmd = &cobra.Command{
	Use:     "delete <id>",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete a module",
	Long:    `Delete a module from the project.`,
	Args:    cobra.ExactArgs(1),
	RunE:    runDelete,
}

var archiveCmd = &cobra.Command{
	Use:   "archive <id>",
	Short: "Archive a module",
	Long:  `Archive a module to hide it from active modules.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runArchive,
}

var issuesCmd = &cobra.Command{
	Use:     "issues <id>",
	Aliases: []string{"work-items"},
	Short:   "List module issues",
	Long:    `List all work items in a specific module.`,
	Args:    cobra.ExactArgs(1),
	RunE:    runIssues,
}

var addIssuesCmd = &cobra.Command{
	Use:   "add-issues <module-id> <issue-ids...\u003e",
	Short: "Add issues to module",
	Long:  `Add work items to a module.`,
	Args:  cobra.MinimumNArgs(2),
	RunE:  runAddIssues,
}

var removeIssueCmd = &cobra.Command{
	Use:   "remove-issue <module-id> <issue-id>",
	Short: "Remove issue from module",
	Long:  `Remove a work item from a module.`,
	Args:  cobra.ExactArgs(2),
	RunE:  runRemoveIssue,
}

func init() {
	ModuleCmd.AddCommand(listCmd)
	ModuleCmd.AddCommand(viewCmd)
	ModuleCmd.AddCommand(createCmd)
	ModuleCmd.AddCommand(editCmd)
	ModuleCmd.AddCommand(deleteCmd)
	ModuleCmd.AddCommand(archiveCmd)
	ModuleCmd.AddCommand(issuesCmd)
	ModuleCmd.AddCommand(addIssuesCmd)
	ModuleCmd.AddCommand(removeIssueCmd)

	// List flags
	listCmd.Flags().BoolVarP(&showArchived, "archived", "a", false, "Show archived modules")

	// Create flags
	createCmd.Flags().StringVarP(&moduleName, "name", "n", "", "Module name")
	createCmd.Flags().StringVarP(&moduleDescription, "description", "d", "", "Module description")
	createCmd.Flags().StringVarP(&moduleStatus, "status", "s", "", "Module status")

	// Edit flags
	editCmd.Flags().StringVarP(&moduleName, "name", "n", "", "New module name")
	editCmd.Flags().StringVarP(&moduleDescription, "description", "d", "", "New module description")
	editCmd.Flags().StringVarP(&moduleStatus, "status", "s", "", "New module status")
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

	modules, err := client.ListModules(projectID, showArchived)
	if err != nil {
		return err
	}

	if len(modules) == 0 {
		output.Info("No modules found")
		return nil
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)

	type moduleOutput struct {
		ID          string `table:"ID" json:"id"`
		Name        string `table:"NAME" json:"name"`
		Status      string `table:"STATUS" json:"status,omitempty"`
		Description string `table:"DESCRIPTION" json:"description,omitempty"`
	}

	var outputs []moduleOutput
	for _, m := range modules {
		outputs = append(outputs, moduleOutput{
			ID:          m.ID,
			Name:        m.Name,
			Status:      m.Status,
			Description: m.Description,
		})
	}

	return formatter.Print(outputs)
}

func runView(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	moduleID := args[0]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	module, err := client.GetModule(projectID, moduleID)
	if err != nil {
		return err
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)
	return formatter.Print(module)
}

func runCreate(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified. Use --project flag or set default project")
	}

	// Interactive prompts if flags not provided
	if moduleName == "" {
		prompt := &survey.Input{
			Message: "Module name:",
			Help:    "e.g., Authentication, API Integration",
		}
		if err := survey.AskOne(prompt, &moduleName); err != nil {
			return err
		}
	}

	if moduleName == "" {
		return fmt.Errorf("module name is required")
	}

	if moduleDescription == "" {
		prompt := &survey.Input{
			Message: "Description (optional):",
		}
		_ = survey.AskOne(prompt, &moduleDescription)
	}

	if moduleStatus == "" {
		statusOptions := []string{"backlog", "planned", "in-progress", "paused", "completed", "cancelled"}
		prompt := &survey.Select{
			Message: "Status:",
			Options: statusOptions,
			Default: "backlog",
		}
		if err := survey.AskOne(prompt, &moduleStatus); err != nil {
			return err
		}
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	req := plane.CreateModuleRequest{
		Name:        moduleName,
		Description: moduleDescription,
		Status:      moduleStatus,
	}

	module, err := client.CreateModule(projectID, req)
	if err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Created module '%s' (%s)", module.Name, module.ID))
	return nil
}

func runEdit(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	moduleID := args[0]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	// Get current module
	module, err := client.GetModule(projectID, moduleID)
	if err != nil {
		return err
	}

	req := plane.UpdateModuleRequest{}

	// Interactive mode if no flags provided
	if moduleName == "" && moduleDescription == "" && moduleStatus == "" {
		output.Info(fmt.Sprintf("Editing module: %s", module.Name))

		prompt := &survey.Input{
			Message: "Name:",
			Default: module.Name,
		}
		if err := survey.AskOne(prompt, &req.Name); err != nil {
			return err
		}

		descPrompt := &survey.Input{
			Message: "Description:",
			Default: module.Description,
		}
		if err := survey.AskOne(descPrompt, &req.Description); err != nil {
			return err
		}

		statusOptions := []string{"backlog", "planned", "in-progress", "paused", "completed", "cancelled"}
		statusPrompt := &survey.Select{
			Message: "Status:",
			Options: statusOptions,
			Default: module.Status,
		}
		if err := survey.AskOne(statusPrompt, &req.Status); err != nil {
			return err
		}
	} else {
		// Use provided flags
		if moduleName != "" {
			req.Name = moduleName
		}
		if moduleDescription != "" {
			req.Description = moduleDescription
		}
		if moduleStatus != "" {
			req.Status = moduleStatus
		}
	}

	updatedModule, err := client.UpdateModule(projectID, moduleID, req)
	if err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Updated module '%s'", updatedModule.Name))
	return nil
}

func runDelete(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	moduleID := args[0]

	// Confirm deletion
	var confirm bool
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("Are you sure you want to delete module %s?", moduleID),
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

	if err := client.DeleteModule(projectID, moduleID); err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Deleted module %s", moduleID))
	return nil
}

func runArchive(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	moduleID := args[0]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	if err := client.ArchiveModule(projectID, moduleID); err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Archived module %s", moduleID))
	return nil
}

func runIssues(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	moduleID := args[0]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	issues, err := client.ListModuleIssues(projectID, moduleID)
	if err != nil {
		return err
	}

	if len(issues) == 0 {
		output.Info("No issues found in this module")
		return nil
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)

	type issueOutput struct {
		ID       string            `table:"ID" json:"id"`
		Sequence int               `table:"#" json:"sequence_id"`
		Title    string            `table:"TITLE" json:"title"`
		State    plane.StateOutput `table:"STATE" json:"state"`
		Priority string            `table:"PRIORITY" json:"priority"`
	}

	var outputs []issueOutput
	for _, issue := range issues {
		outputs = append(outputs, issueOutput{
			ID:       issue.ID,
			Sequence: issue.SequenceID,
			Title:    issue.Name,
			State: plane.StateOutput{
				ID:   issue.State.ID,
				Name: issue.State.Name,
			},
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

	moduleID := args[0]
	issueIDs := args[1:]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	if err := client.AddIssuesToModule(projectID, moduleID, issueIDs); err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Added %d issue(s) to module %s", len(issueIDs), moduleID))
	return nil
}

func runRemoveIssue(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	moduleID := args[0]
	issueID := args[1]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	if err := client.RemoveIssueFromModule(projectID, moduleID, issueID); err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Removed issue %s from module %s", issueID, moduleID))
	return nil
}
