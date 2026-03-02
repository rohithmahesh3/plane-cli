package workspace

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/rohithmahesh3/plane-cli/internal/api"
	"github.com/rohithmahesh3/plane-cli/internal/config"
	"github.com/rohithmahesh3/plane-cli/internal/output"
	"github.com/spf13/cobra"
)

var WorkspaceCmd = &cobra.Command{
	Use:     "workspace",
	Aliases: []string{"ws"},
	Short:   "Manage workspaces",
	Long: `Manage Plane workspaces.

Note: The Plane API does not support listing or retrieving workspace details.
You can only switch between configured workspaces.`,
}

var infoCmd = &cobra.Command{
	Use:   "info [slug]",
	Short: "Show workspace details",
	Long: `Display detailed information about a specific workspace.

Note: The Plane API does not have a workspace details endpoint.
This command will show the currently configured workspace.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runInfo,
}

var switchCmd = &cobra.Command{
	Use:   "switch [slug]",
	Short: "Switch default workspace",
	Long: `Set the default workspace for all future commands.

Note: Since the Plane API doesn't support workspace listing, you need to
provide the workspace slug manually or configure it interactively.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runSwitch,
}

var membersCmd = &cobra.Command{
	Use:     "members",
	Aliases: []string{"users", "people"},
	Short:   "List workspace members",
	Long:    `Display all members of the current workspace.`,
	RunE:    runMembers,
}

func init() {
	WorkspaceCmd.AddCommand(infoCmd)
	WorkspaceCmd.AddCommand(switchCmd)
	WorkspaceCmd.AddCommand(membersCmd)
}

func runInfo(cmd *cobra.Command, args []string) error {
	slug := config.Cfg.DefaultWorkspace
	if len(args) > 0 {
		slug = args[0]
	}

	if slug == "" {
		return fmt.Errorf("no workspace specified. Use --workspace flag or provide workspace slug")
	}

	// The Plane API doesn't have a workspace info endpoint
	// Show what we have configured
	fmt.Printf("Workspace: %s\n", slug)
	fmt.Printf("API Host: %s\n", config.Cfg.APIHost)
	fmt.Printf("Default Project: %s\n", config.Cfg.DefaultProject)

	// Try to list projects to verify access
	client, err := api.NewClient()
	if err != nil {
		return err
	}

	projects, err := client.ListProjects()
	if err != nil {
		output.Warning(fmt.Sprintf("Could not access workspace: %v", err))
		return nil
	}

	fmt.Printf("\nProjects in workspace: %d\n", len(projects))
	for _, p := range projects {
		fmt.Printf("  - %s (%s)\n", p.Name, p.Identifier)
	}

	return nil
}

func runSwitch(cmd *cobra.Command, args []string) error {
	var slug string
	if len(args) > 0 {
		slug = args[0]
	} else {
		// Interactive prompt
		prompt := &survey.Input{
			Message: "Enter workspace slug:",
			Help:    "This is the unique identifier for your workspace (found in the URL)",
		}
		if err := survey.AskOne(prompt, &slug); err != nil {
			return err
		}
	}

	if slug == "" {
		return fmt.Errorf("workspace slug is required")
	}

	// Verify the workspace by trying to list projects
	client, err := api.NewClient()
	if err != nil {
		return err
	}

	// Temporarily set the workspace to test it
	oldWorkspace := config.Cfg.DefaultWorkspace
	config.Cfg.DefaultWorkspace = slug
	client.Workspace = slug

	_, err = client.ListProjects()
	if err != nil {
		// Restore old workspace
		config.Cfg.DefaultWorkspace = oldWorkspace
		return fmt.Errorf("could not access workspace '%s': %w", slug, err)
	}

	// Save the new workspace
	config.Cfg.DefaultWorkspace = slug
	if err := config.SaveConfig(); err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Switched to workspace '%s'", slug))
	return nil
}

func runMembers(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient()
	if err != nil {
		return err
	}

	members, err := client.GetWorkspaceMembers()
	if err != nil {
		return err
	}

	if len(members) == 0 {
		output.Info("No members found")
		return nil
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)

	type memberOutput struct {
		ID          string `table:"ID" json:"id"`
		DisplayName string `table:"NAME" json:"display_name"`
		Email       string `table:"EMAIL" json:"email"`
		Role        int    `table:"ROLE" json:"role"`
	}

	var outputs []memberOutput
	for _, m := range members {
		outputs = append(outputs, memberOutput{
			ID:          m.ID,
			DisplayName: m.DisplayName,
			Email:       m.Email,
			Role:        0, // Role not provided in current User struct
		})
	}

	return formatter.Print(outputs)
}
