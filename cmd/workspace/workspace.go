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
	Long:    `List, view, and switch between Plane workspaces.`,
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List workspaces",
	Long:    `List all workspaces you have access to.`,
	RunE:    runList,
}

var infoCmd = &cobra.Command{
	Use:   "info [slug]",
	Short: "Show workspace details",
	Long:  `Display detailed information about a specific workspace.`,
	Args:  cobra.MaximumNArgs(1),
	RunE:  runInfo,
}

var switchCmd = &cobra.Command{
	Use:   "switch [slug]",
	Short: "Switch default workspace",
	Long:  `Set the default workspace for all future commands.`,
	Args:  cobra.MaximumNArgs(1),
	RunE:  runSwitch,
}

func init() {
	WorkspaceCmd.AddCommand(listCmd)
	WorkspaceCmd.AddCommand(infoCmd)
	WorkspaceCmd.AddCommand(switchCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient()
	if err != nil {
		return err
	}
	
	workspaces, err := client.ListWorkspaces()
	if err != nil {
		return err
	}
	
	if len(workspaces) == 0 {
		output.Info("No workspaces found")
		return nil
	}
	
	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)
	
	// Add a marker for default workspace
	type workspaceOutput struct {
		ID          string `table:"ID" json:"id"`
		Name        string `table:"NAME" json:"name"`
		Slug        string `table:"SLUG" json:"slug"`
		Description string `table:"DESCRIPTION" json:"description,omitempty"`
		Default     string `table:"DEFAULT" json:"default,omitempty"`
	}
	
	var outputs []workspaceOutput
	for _, ws := range workspaces {
		isDefault := ""
		if ws.Slug == config.Cfg.DefaultWorkspace {
			isDefault = "✓"
		}
		outputs = append(outputs, workspaceOutput{
			ID:          ws.ID,
			Name:        ws.Name,
			Slug:        ws.Slug,
			Description: ws.Description,
			Default:     isDefault,
		})
	}
	
	return formatter.Print(outputs)
}

func runInfo(cmd *cobra.Command, args []string) error {
	slug := config.Cfg.DefaultWorkspace
	if len(args) > 0 {
		slug = args[0]
	}
	
	if slug == "" {
		return fmt.Errorf("no workspace specified. Use --workspace flag or provide workspace slug")
	}
	
	client, err := api.NewClient()
	if err != nil {
		return err
	}
	
	workspace, err := client.GetWorkspace(slug)
	if err != nil {
		return err
	}
	
	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)
	return formatter.Print(workspace)
}

func runSwitch(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient()
	if err != nil {
		return err
	}
	
	workspaces, err := client.ListWorkspaces()
	if err != nil {
		return err
	}
	
	var slug string
	if len(args) > 0 {
		slug = args[0]
	} else {
		// Interactive selection
		var options []string
		for _, ws := range workspaces {
			options = append(options, fmt.Sprintf("%s (%s)", ws.Name, ws.Slug))
		}
		
		var selected string
		prompt := &survey.Select{
			Message: "Select workspace:",
			Options: options,
		}
		if err := survey.AskOne(prompt, &selected); err != nil {
			return err
		}
		
		// Extract slug from selection
		for _, ws := range workspaces {
			if fmt.Sprintf("%s (%s)", ws.Name, ws.Slug) == selected {
				slug = ws.Slug
				break
			}
		}
	}
	
	if slug == "" {
		return fmt.Errorf("workspace slug is required")
	}
	
	// Verify workspace exists
	_, err = client.GetWorkspace(slug)
	if err != nil {
		return fmt.Errorf("workspace not found: %s", slug)
	}
	
	config.Cfg.DefaultWorkspace = slug
	if err := config.SaveConfig(); err != nil {
		return err
	}
	
	output.Success(fmt.Sprintf("Switched to workspace '%s'", slug))
	return nil
}
