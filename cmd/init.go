package cmd

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/rohithmahesh3/plane-cli/internal/api"
	"github.com/rohithmahesh3/plane-cli/internal/config"
	"github.com/rohithmahesh3/plane-cli/internal/output"
	"github.com/rohithmahesh3/plane-cli/pkg/plane"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	initWorkspace     string
	initProject       string
	initCreateProject bool
	initProjectName   string
	initProjectIdent  string
	initProjectDesc   string
	initSkipGitignore bool
	initUpgrade       bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Plane CLI for current directory",
	Long: `Initialize .plane/settings.yaml for project-local configuration.

This creates a .plane/settings.yaml file in the current directory with
workspace and project settings, allowing you to use plane commands without
specifying --workspace and --project flags.

Examples:
  plane init                                    # Interactive setup
  plane init --workspace my-ws                  # Use specific workspace
  plane init --workspace my-ws --project FRONT  # Use existing project
  plane init --workspace my-ws --create-project --project-name "New App"
  plane init --upgrade                          # Update existing settings`,
	RunE: runInit,
}

func init() {
	initCmd.Flags().StringVar(&initWorkspace, "workspace", "", "Workspace slug")
	initCmd.Flags().StringVar(&initProject, "project", "", "Project ID or identifier")
	initCmd.Flags().BoolVar(&initCreateProject, "create-project", false, "Create a new project")
	initCmd.Flags().StringVar(&initProjectName, "project-name", "", "New project name (with --create-project)")
	initCmd.Flags().StringVar(&initProjectIdent, "project-identifier", "", "New project identifier (with --create-project)")
	initCmd.Flags().StringVar(&initProjectDesc, "project-description", "", "New project description (with --create-project)")
	initCmd.Flags().BoolVar(&initSkipGitignore, "skip-gitignore", false, "Skip adding .plane/ to .gitignore")
	initCmd.Flags().BoolVar(&initUpgrade, "upgrade", false, "Upgrade existing .plane/settings.yaml")
}

func runInit(cmd *cobra.Command, args []string) error {
	apiKey, err := config.GetAPIKey()
	if err != nil || apiKey == "" {
		return fmt.Errorf("not authenticated. Run 'plane auth login' first")
	}

	if err := config.InitConfig(); err != nil {
		return fmt.Errorf("failed to initialize config: %w", err)
	}

	settingsPath := filepath.Join(".", ".plane", "settings.yaml")
	if _, err := os.Stat(settingsPath); err == nil {
		if !initUpgrade {
			var overwrite bool
			if err := survey.AskOne(&survey.Confirm{
				Message: ".plane/settings.yaml already exists. Overwrite?",
				Default: false,
			}, &overwrite); err != nil {
				return err
			}
			if !overwrite {
				output.Info("Cancelled. Use --upgrade to update existing settings")
				return nil
			}
		}
	}

	workspaceSlug, err := getWorkspaceSlug()
	if err != nil {
		return err
	}

	client := &api.Client{
		HTTPClient: &http.Client{Timeout: api.DefaultTimeout},
		BaseURL:    config.Cfg.APIHost,
		APIKey:     apiKey,
		Workspace:  workspaceSlug,
	}

	output.Info(fmt.Sprintf("Fetching projects from workspace '%s'...", workspaceSlug))
	projects, err := client.ListProjects()
	if err != nil {
		return fmt.Errorf("failed to fetch projects from workspace '%s': %w", workspaceSlug, err)
	}

	var projectID string
	if initCreateProject {
		projectID, err = createNewProjectNonInteractive(client, initProjectName, initProjectIdent, initProjectDesc)
		if err != nil {
			return err
		}
	} else if initProject != "" {
		projectID, err = resolveProjectID(initProject, projects)
		if err != nil {
			return err
		}
	} else {
		projectID, err = selectOrCreateProject(client, projects)
		if err != nil {
			return err
		}
	}

	project, err := client.GetProject(projectID)
	if err != nil {
		return fmt.Errorf("failed to get project details: %w", err)
	}

	fmt.Println()
	fmt.Println("  Preview settings:")
	fmt.Printf("    Workspace: %s\n", workspaceSlug)
	fmt.Printf("    Project: %s (%s)\n", project.Name, project.Identifier)
	fmt.Println()

	var confirm bool
	if err := survey.AskOne(&survey.Confirm{
		Message: "Create .plane/settings.yaml?",
		Default: true,
	}, &confirm); err != nil {
		return err
	}

	if !confirm {
		output.Info("Cancelled")
		return nil
	}

	if err := createSettingsFile(workspaceSlug, projectID, settingsPath); err != nil {
		return err
	}

	output.Success("Created .plane/settings.yaml")

	if !initSkipGitignore {
		if err := handleGitignore(); err != nil {
			output.Warning(fmt.Sprintf("Failed to update .gitignore: %v", err))
		}
	}

	printSuccessMessage()
	return nil
}

func getWorkspaceSlug() (string, error) {
	if initWorkspace != "" {
		return initWorkspace, nil
	}

	if config.Cfg.DefaultWorkspace != "" {
		var useExisting bool
		if err := survey.AskOne(&survey.Confirm{
			Message: fmt.Sprintf("Use workspace '%s' from global config?", config.Cfg.DefaultWorkspace),
			Default: true,
		}, &useExisting); err != nil {
			return "", err
		}
		if useExisting {
			return config.Cfg.DefaultWorkspace, nil
		}
	}

	var workspace string
	err := survey.AskOne(&survey.Input{
		Message: "Enter your workspace slug:",
		Help:    "Find your workspace slug in your Plane URL (e.g., app.plane.so/my-workspace → slug is 'my-workspace')",
	}, &workspace, survey.WithValidator(survey.Required))
	if err != nil {
		return "", err
	}

	return workspace, nil
}

func selectOrCreateProject(client *api.Client, projects []plane.Project) (string, error) {
	type projectOption struct {
		display  string
		id       string
		isCreate bool
	}

	options := []projectOption{}
	for _, p := range projects {
		options = append(options, projectOption{
			display: fmt.Sprintf("%s (%s)", p.Name, p.Identifier),
			id:      p.ID,
		})
	}
	options = append(options, projectOption{
		display:  "✨ Create new project",
		isCreate: true,
	})

	var selectedIdx int
	displays := make([]string, len(options))
	for i, opt := range options {
		displays[i] = opt.display
	}

	if err := survey.AskOne(&survey.Select{
		Message: "Select a project:",
		Options: displays,
	}, &selectedIdx); err != nil {
		return "", err
	}

	if options[selectedIdx].isCreate {
		return createNewProjectFlow(client)
	}

	return options[selectedIdx].id, nil
}

func createNewProjectFlow(client *api.Client) (string, error) {
	var name, identifier, description string

	err := survey.AskOne(&survey.Input{
		Message: "Project name:",
	}, &name, survey.WithValidator(survey.Required))
	if err != nil {
		return "", err
	}

	defaultIdentifier := generateIdentifier(name)
	err = survey.AskOne(&survey.Input{
		Message: "Project identifier:",
		Default: defaultIdentifier,
		Help:    "Short unique identifier (max 10 characters, uppercase letters, numbers, and hyphens)",
	}, &identifier)
	if err != nil {
		return "", err
	}

	if identifier == "" {
		identifier = defaultIdentifier
	}

	err = survey.AskOne(&survey.Input{
		Message: "Description (optional):",
	}, &description)
	if err != nil {
		return "", err
	}

	output.Info("Creating project...")
	project, err := client.CreateProject(plane.CreateProjectRequest{
		Name:        name,
		Identifier:  identifier,
		Description: description,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create project: %w", err)
	}

	output.Success(fmt.Sprintf("Created project '%s' (%s)", project.Name, project.Identifier))
	return project.ID, nil
}

func createNewProjectNonInteractive(client *api.Client, name, identifier, description string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("--project-name is required when using --create-project")
	}

	if identifier == "" {
		identifier = generateIdentifier(name)
	}

	output.Info("Creating project...")
	project, err := client.CreateProject(plane.CreateProjectRequest{
		Name:        name,
		Identifier:  identifier,
		Description: description,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create project: %w", err)
	}

	output.Success(fmt.Sprintf("Created project '%s' (%s)", project.Name, project.Identifier))
	return project.ID, nil
}

func generateIdentifier(name string) string {
	identifier := strings.ToUpper(strings.ReplaceAll(name, " ", "-"))
	identifier = strings.ReplaceAll(identifier, "_", "-")
	re := regexp.MustCompile("[^A-Z0-9-]")
	identifier = re.ReplaceAllString(identifier, "")
	if len(identifier) > 10 {
		identifier = identifier[:10]
	}
	identifier = strings.Trim(identifier, "-")
	if identifier == "" {
		identifier = "PROJ"
	}
	return identifier
}

func resolveProjectID(projectIDOrIdentifier string, projects []plane.Project) (string, error) {
	if len(projectIDOrIdentifier) == 36 && strings.Count(projectIDOrIdentifier, "-") == 4 {
		return projectIDOrIdentifier, nil
	}

	for _, p := range projects {
		if strings.EqualFold(p.Identifier, projectIDOrIdentifier) || p.ID == projectIDOrIdentifier {
			return p.ID, nil
		}
	}

	return "", fmt.Errorf("project '%s' not found. Use 'plane project list' to see available projects", projectIDOrIdentifier)
}

func createSettingsFile(workspace, projectID, settingsPath string) error {
	planeDir := filepath.Dir(settingsPath)
	if err := os.MkdirAll(planeDir, 0755); err != nil {
		return fmt.Errorf("failed to create .plane directory: %w", err)
	}

	settings := config.LocalConfig{
		Workspace: workspace,
		Project:   projectID,
	}

	var buf strings.Builder
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	if err := encoder.Encode(&settings); err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}
	if err := encoder.Close(); err != nil {
		return fmt.Errorf("failed to close encoder: %w", err)
	}

	header := "# Plane CLI project settings\n# This file overrides global config for this directory\n\n"
	fullContent := header + buf.String()

	if err := os.WriteFile(settingsPath, []byte(fullContent), 0644); err != nil {
		return fmt.Errorf("failed to write settings file: %w", err)
	}

	return nil
}

func handleGitignore() error {
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		return nil
	}

	if _, err := os.Stat(".gitignore"); os.IsNotExist(err) {
		return nil
	}

	content, err := os.ReadFile(".gitignore")
	if err != nil {
		return err
	}

	if strings.Contains(string(content), ".plane/") {
		return nil
	}

	var addToGitignore bool
	if err := survey.AskOne(&survey.Confirm{
		Message: "Add .plane/ to .gitignore?",
		Default: true,
	}, &addToGitignore); err != nil {
		return err
	}

	if !addToGitignore {
		return nil
	}

	f, err := os.OpenFile(".gitignore", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	if len(content) > 0 && content[len(content)-1] != '\n' {
		if _, err := f.WriteString("\n"); err != nil {
			return err
		}
	}

	_, err = f.WriteString("\n# Plane CLI\n.plane/\n")
	if err != nil {
		return err
	}

	output.Success("Added .plane/ to .gitignore")
	return nil
}

func printSuccessMessage() {
	fmt.Println()
	output.Info("You can now use plane commands without specifying workspace/project:")
	fmt.Println("  plane issue list")
	fmt.Println("  plane cycle list")
	fmt.Println("  plane module list")
}
