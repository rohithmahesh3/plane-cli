package auth

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/rohithmahesh3/plane-cli/internal/api"
	"github.com/rohithmahesh3/plane-cli/internal/config"
	"github.com/rohithmahesh3/plane-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	token     string
	apiHost   string
	workspace string
)

var AuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication commands",
	Long:  `Manage authentication with Plane.`,
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Plane",
	Long: `Authenticate with Plane using an API key.

You can generate an API key from your Plane workspace settings:
1. Go to Profile Settings → Personal Access Tokens
2. Click "Add personal access token"
3. Copy the generated token

Example:
  plane auth login
  plane auth login --token YOUR_API_KEY --workspace my-workspace`,
	RunE: runLogin,
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from Plane",
	Long:  `Remove stored credentials and configuration.`,
	RunE:  runLogout,
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check authentication status",
	Long:  `Check if you're authenticated and show current configuration.`,
	RunE:  runStatus,
}

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show current user information",
	Long:  `Display information about the currently authenticated user.`,
	RunE:  runWhoami,
}

func init() {
	AuthCmd.AddCommand(loginCmd)
	AuthCmd.AddCommand(logoutCmd)
	AuthCmd.AddCommand(statusCmd)
	AuthCmd.AddCommand(whoamiCmd)
	
	loginCmd.Flags().StringVar(&token, "token", "", "API key (will prompt if not provided)")
	loginCmd.Flags().StringVar(&apiHost, "api-host", config.DefaultAPIHost, "Plane API host URL")
	loginCmd.Flags().StringVar(&workspace, "workspace", "", "Default workspace slug")
}

func runLogin(cmd *cobra.Command, args []string) error {
	// Interactive prompts if flags not provided
	if token == "" {
		prompt := &survey.Password{
			Message: "Enter your Plane API key:",
			Help:    "You can generate an API key from Profile Settings → Personal Access Tokens",
		}
		if err := survey.AskOne(prompt, &token); err != nil {
			return err
		}
	}
	
	if token == "" {
		return fmt.Errorf("API key is required")
	}
	
	if apiHost == "" {
		apiHost = config.DefaultAPIHost
	}
	
	if workspace == "" {
		prompt := &survey.Input{
			Message: "Enter your default workspace slug:",
			Help:    "This is the unique identifier for your workspace (found in the URL)",
		}
		if err := survey.AskOne(prompt, &workspace); err != nil {
			return err
		}
	}
	
	if workspace == "" {
		return fmt.Errorf("workspace is required")
	}
	
	// Test the credentials
	if err := config.SetAPIKey(token); err != nil {
		return fmt.Errorf("failed to save API key: %w", err)
	}
	
	// Initialize config
	if err := config.InitConfig(); err != nil {
		return fmt.Errorf("failed to initialize config: %w", err)
	}
	
	config.Cfg.APIHost = apiHost
	config.Cfg.DefaultWorkspace = workspace
	
	// Test authentication
	client, err := api.NewClient()
	if err != nil {
		_ = config.DeleteAPIKey()
		return fmt.Errorf("failed to create API client: %w", err)
	}
	
	workspaces, err := client.ListWorkspaces()
	if err != nil {
		_ = config.DeleteAPIKey()
		return fmt.Errorf("authentication failed: %w", err)
	}
	
	// Save config
	if err := config.SaveConfig(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}
	
	output.Success(fmt.Sprintf("Successfully authenticated with workspace '%s'", workspace))
	output.Info(fmt.Sprintf("Found %d workspace(s)", len(workspaces)))
	
	return nil
}

func runLogout(cmd *cobra.Command, args []string) error {
	if err := config.DeleteAPIKey(); err != nil {
		return fmt.Errorf("failed to remove API key: %w", err)
	}
	
	output.Success("Successfully logged out")
	return nil
}

func runStatus(cmd *cobra.Command, args []string) error {
	apiKey, err := config.GetAPIKey()
	if err != nil || apiKey == "" {
		output.Error("Not authenticated")
		output.Info("Run 'plane auth login' to authenticate")
		return nil
	}
	
	if err := config.InitConfig(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	
	fmt.Println("Authentication Status: ✓ Authenticated")
	fmt.Printf("Workspace: %s\n", config.Cfg.DefaultWorkspace)
	fmt.Printf("API Host: %s\n", config.Cfg.APIHost)
	fmt.Printf("Default Project: %s\n", config.Cfg.DefaultProject)
	fmt.Printf("Output Format: %s\n", config.Cfg.OutputFormat)
	
	return nil
}

func runWhoami(cmd *cobra.Command, args []string) error {
	if err := config.InitConfig(); err != nil {
		return fmt.Errorf("failed to initialize config: %w", err)
	}
	
	client, err := api.NewClient()
	if err != nil {
		return err
	}
	
	// TODO: Add GetMe endpoint to API
	output.Info("User information endpoint to be implemented")
	fmt.Printf("Workspace: %s\n", client.Workspace)
	
	return nil
}
