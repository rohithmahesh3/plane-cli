package cmd

import (
	"fmt"
	"os"

	"github.com/rohithmahesh3/plane-cli/cmd/auth"
	"github.com/rohithmahesh3/plane-cli/cmd/config"
	"github.com/rohithmahesh3/plane-cli/cmd/cycle"
	"github.com/rohithmahesh3/plane-cli/cmd/issue"
	"github.com/rohithmahesh3/plane-cli/cmd/module"
	"github.com/rohithmahesh3/plane-cli/cmd/project"
	issuetype "github.com/rohithmahesh3/plane-cli/cmd/type"
	"github.com/rohithmahesh3/plane-cli/cmd/workspace"
	cfg "github.com/rohithmahesh3/plane-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	version   = "dev"
	commit     = "none"
	date       = "unknown"
	
	workspaceSlug string
	projectID     string
	outputFmt     string
	noColor       bool
	configFile    string
)

var rootCmd = &cobra.Command{
	Use:   "plane",
	Short: "A CLI tool for managing Plane projects",
	Long: `plane-cli is a command-line interface for Plane project management.

It allows you to manage workspaces, projects, issues, cycles, and modules
from the comfort of your terminal.

Get started:
  plane auth login                    # Authenticate with Plane
  plane workspace list                # List your workspaces
  plane project list                  # List projects in current workspace
  plane issue list                    # List issues in current project`,
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip config initialization for certain commands
		if cmd.Name() == "login" || cmd.Name() == "version" || cmd.Name() == "completion" {
			return nil
		}
		
		if configFile != "" {
			cfg.SetConfigFile(configFile)
		}
		
		if err := cfg.InitConfig(); err != nil {
			return fmt.Errorf("failed to initialize config: %w", err)
		}
		
		// Override config with flags
		if workspaceSlug != "" {
			cfg.Cfg.DefaultWorkspace = workspaceSlug
		}
		if projectID != "" {
			cfg.Cfg.DefaultProject = projectID
		}
		if outputFmt != "" {
			cfg.Cfg.OutputFormat = outputFmt
		}
		
		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&workspaceSlug, "workspace", "", "Plane workspace slug (overrides config)")
	rootCmd.PersistentFlags().StringVar(&projectID, "project", "", "Plane project ID (overrides config)")
	rootCmd.PersistentFlags().StringVarP(&outputFmt, "output", "o", "", "Output format: table, json, yaml (overrides config)")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable colored output")
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "Config file path (default: ~/.config/plane-cli/config.yaml)")
	
	// Add subcommands
	rootCmd.AddCommand(auth.AuthCmd)
	rootCmd.AddCommand(workspace.WorkspaceCmd)
	rootCmd.AddCommand(project.ProjectCmd)
	rootCmd.AddCommand(issue.IssueCmd)
	rootCmd.AddCommand(cycle.CycleCmd)
	rootCmd.AddCommand(module.ModuleCmd)
	rootCmd.AddCommand(config.ConfigCmd)
	rootCmd.AddCommand(issuetype.TypeCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(completionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("plane-cli version %s (commit: %s, built: %s)\n", version, commit, date)
	},
}

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:
  $ source <(plane completion bash)
  # To load completions for each session, execute once:
  # Linux:
  $ plane completion bash > /etc/bash_completion.d/plane
  # macOS:
  $ plane completion bash > $(brew --prefix)/etc/bash_completion.d/plane

Zsh:
  $ source <(plane completion zsh)
  # To load completions for each session, execute once:
  $ plane completion zsh > "${fpath[1]}/_plane"

Fish:
  $ plane completion fish | source
  # To load completions for each session, execute once:
  $ plane completion fish > ~/.config/fish/completions/plane.fish

PowerShell:
  PS> plane completion powershell | Out-String | Invoke-Expression
  # To load completions for every new session, run:
  PS> plane completion powershell > plane.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		}
	},
}
