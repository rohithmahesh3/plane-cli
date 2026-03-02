package project

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/rohithmahesh3/plane-cli/internal/api"
	"github.com/rohithmahesh3/plane-cli/internal/config"
	"github.com/rohithmahesh3/plane-cli/internal/output"
	"github.com/rohithmahesh3/plane-cli/pkg/plane"
	"github.com/spf13/cobra"
)

var (
	allFlag bool
)

var ProjectCmd = &cobra.Command{
	Use:     "project",
	Aliases: []string{"proj"},
	Short:   "Manage projects",
	Long:    `List, create, and manage Plane projects.`,
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List projects",
	Long:    `List all projects in the current workspace.`,
	RunE:    runList,
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new project",
	Long:  `Create a new project in the current workspace.`,
	RunE:  runCreate,
}

var infoCmd = &cobra.Command{
	Use:   "info [id]",
	Short: "Show project details",
	Long:  `Display detailed information about a specific project.`,
	Args:  cobra.MaximumNArgs(1),
	RunE:  runInfo,
}

var deleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a project",
	Long:  `Delete a project from the workspace.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runDelete,
}

var membersCmd = &cobra.Command{
	Use:   "members [id]",
	Short: "List project members",
	Long:  `List all members of a specific project.`,
	Args:  cobra.MaximumNArgs(1),
	RunE:  runMembers,
}

func init() {
	ProjectCmd.AddCommand(listCmd)
	ProjectCmd.AddCommand(createCmd)
	ProjectCmd.AddCommand(infoCmd)
	ProjectCmd.AddCommand(deleteCmd)
	ProjectCmd.AddCommand(membersCmd)

	listCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "List all accessible projects")
}

func runList(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient()
	if err != nil {
		return err
	}

	projects, err := client.ListProjects()
	if err != nil {
		return err
	}

	if len(projects) == 0 {
		output.Info("No projects found")
		return nil
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)

	type projectOutput struct {
		ID         string `table:"ID" json:"id"`
		Identifier string `table:"IDENTIFIER" json:"identifier"`
		Name       string `table:"NAME" json:"name"`
		Default    string `table:"DEFAULT" json:"default,omitempty"`
	}

	var outputs []projectOutput
	for _, p := range projects {
		isDefault := ""
		if p.ID == config.Cfg.DefaultProject {
			isDefault = "✓"
		}
		outputs = append(outputs, projectOutput{
			ID:         p.ID,
			Identifier: p.Identifier,
			Name:       p.Name,
			Default:    isDefault,
		})
	}

	return formatter.Print(outputs)
}

func runCreate(cmd *cobra.Command, args []string) error {
	var name, identifier string

	if len(args) >= 1 {
		name = args[0]
	} else {
		prompt := &survey.Input{
			Message: "Project name:",
		}
		if err := survey.AskOne(prompt, &name); err != nil {
			return err
		}
	}

	if name == "" {
		return fmt.Errorf("project name is required")
	}

	// Generate default identifier from name
	defaultIdentifier := strings.ToUpper(strings.ReplaceAll(name, " ", "-"))
	defaultIdentifier = strings.ReplaceAll(defaultIdentifier, "_", "-")
	if len(defaultIdentifier) > 10 {
		defaultIdentifier = defaultIdentifier[:10]
	}

	prompt := &survey.Input{
		Message: "Project identifier:",
		Default: defaultIdentifier,
		Help:    "Short unique identifier for the project (e.g., PROJ, WEB)",
	}
	if err := survey.AskOne(prompt, &identifier); err != nil {
		return err
	}

	if identifier == "" {
		identifier = defaultIdentifier
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	project, err := client.CreateProject(plane.CreateProjectRequest{
		Name:       name,
		Identifier: identifier,
	})
	if err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Created project '%s' (%s)", project.Name, project.Identifier))

	// Ask if user wants to set as default
	var setDefault bool
	confirmPrompt := &survey.Confirm{
		Message: "Set as default project?",
		Default: true,
	}
	if err := survey.AskOne(confirmPrompt, &setDefault); err != nil {
		return err
	}

	if setDefault {
		config.Cfg.DefaultProject = project.ID
		if err := config.SaveConfig(); err != nil {
			return err
		}
		output.Info("Set as default project")
	}

	return nil
}

func runInfo(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if len(args) > 0 {
		projectID = args[0]
	}

	if projectID == "" {
		return fmt.Errorf("no project specified. Use --project flag or provide project ID")
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	project, err := client.GetProject(projectID)
	if err != nil {
		return err
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)
	return formatter.Print(project)
}

func runDelete(cmd *cobra.Command, args []string) error {
	projectID := args[0]

	// Confirm deletion
	var confirm bool
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("Are you sure you want to delete project %s?", projectID),
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

	if err := client.DeleteProject(projectID); err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Deleted project %s", projectID))
	return nil
}

func runMembers(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if len(args) > 0 {
		projectID = args[0]
	}

	if projectID == "" {
		return fmt.Errorf("no project specified. Use --project flag or provide project ID")
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	members, err := client.GetProjectMembers(projectID)
	if err != nil {
		return err
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)
	return formatter.Print(members)
}
