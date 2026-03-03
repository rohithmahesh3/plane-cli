package cmd

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/rohithmahesh3/plane-cli/cmd/auth"
	"github.com/rohithmahesh3/plane-cli/cmd/config"
	"github.com/rohithmahesh3/plane-cli/cmd/context"
	"github.com/rohithmahesh3/plane-cli/cmd/cycle"
	"github.com/rohithmahesh3/plane-cli/cmd/epic"
	"github.com/rohithmahesh3/plane-cli/cmd/inject"
	"github.com/rohithmahesh3/plane-cli/cmd/intake"
	"github.com/rohithmahesh3/plane-cli/cmd/issue"
	"github.com/rohithmahesh3/plane-cli/cmd/label"
	"github.com/rohithmahesh3/plane-cli/cmd/module"
	"github.com/rohithmahesh3/plane-cli/cmd/project"
	"github.com/rohithmahesh3/plane-cli/cmd/state"
	issuetype "github.com/rohithmahesh3/plane-cli/cmd/type"
	"github.com/rohithmahesh3/plane-cli/cmd/workspace"
	cfg "github.com/rohithmahesh3/plane-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"

	workspaceSlug string
	projectID     string
	outputFmt     string
	noColor       bool
	configFile    string
)

var rootCmd = &cobra.Command{
	Use:   "plane-cli",
	Short: "A CLI tool for managing Plane projects",
	Long: `plane-cli is a command-line interface for Plane project management.

It allows you to manage workspaces, projects, issues, cycles, and modules
from the comfort of your terminal.

Get started:
  plane-cli auth login                    # Authenticate with Plane
  plane-cli workspace info                # Show configured workspace access
  plane-cli workspace members --search alice # Find assignable workspace users by name/email
  plane-cli project list                  # List projects in current workspace
  plane-cli issue list                    # List issues in current project`,
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip config initialization for certain commands
		if cmd.Name() == "login" || cmd.Name() == "version" || cmd.Name() == "completion" || cmd.Name() == "init" {
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
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(auth.AuthCmd)
	rootCmd.AddCommand(workspace.WorkspaceCmd)
	rootCmd.AddCommand(project.ProjectCmd)
	rootCmd.AddCommand(issue.IssueCmd)
	rootCmd.AddCommand(state.StateCmd)
	rootCmd.AddCommand(label.LabelCmd)
	rootCmd.AddCommand(cycle.CycleCmd)
	rootCmd.AddCommand(module.ModuleCmd)
	rootCmd.AddCommand(epic.EpicCmd)
	rootCmd.AddCommand(intake.IntakeCmd)
	rootCmd.AddCommand(config.ConfigCmd)
	rootCmd.AddCommand(context.ContextCmd)
	rootCmd.AddCommand(issuetype.TypeCmd)
	rootCmd.AddCommand(inject.InjectCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(completionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		// Try to get version from Go module info (works with go install)
		// Version is set when installed via: go install github.com/rohithmahesh3/plane-cli@v1.0.1
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
			fmt.Printf("plane-cli version %s\n", info.Main.Version)
			return
		}

		// Fall back to ldflags for Makefile builds or local development
		fmt.Printf("plane-cli version %s (commit: %s, built: %s)\n", version, commit, date)
	},
}

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:
  $ source <(plane-cli completion bash)
  # To load completions for each session, execute once:
  # Linux:
  $ plane-cli completion bash > /etc/bash_completion.d/plane-cli
  # macOS:
  $ plane-cli completion bash > $(brew --prefix)/etc/bash_completion.d/plane-cli

Zsh:
  $ source <(plane-cli completion zsh)
  # To load completions for each session, execute once:
  $ plane-cli completion zsh > "${fpath[1]}/_plane-cli"

Fish:
  $ source <(plane-cli completion fish)
  # To load completions for each session, execute once:
  $ plane-cli completion fish > ~/.config/fish/completions/plane-cli.fish

PowerShell:
  PS> plane-cli completion powershell | Out-String | Invoke-Expression
  # To load completions for every new session, run:
  PS> plane-cli completion powershell > plane-cli.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		}
		return nil
	},
}
