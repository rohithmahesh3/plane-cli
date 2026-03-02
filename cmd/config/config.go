package config

import (
	"fmt"

	"github.com/rohithmahesh3/plane-cli/internal/config"
	"github.com/rohithmahesh3/plane-cli/internal/output"
	"github.com/spf13/cobra"
)

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  `View and modify plane-cli configuration.`,
}

var getCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get configuration value",
	Long:  `Get a specific configuration value or all values.`,
	Args:  cobra.MaximumNArgs(1),
	RunE:  runGet,
}

var setCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set configuration value",
	Long:  `Set a configuration value.`,
	Args:  cobra.ExactArgs(2),
	RunE:  runSet,
}

func init() {
	ConfigCmd.AddCommand(getCmd)
	ConfigCmd.AddCommand(setCmd)
}

func runGet(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		// Show all config
		fmt.Printf("version: %s\n", config.Cfg.Version)
		fmt.Printf("default_workspace: %s\n", config.Cfg.DefaultWorkspace)
		fmt.Printf("default_project: %s\n", config.Cfg.DefaultProject)
		fmt.Printf("output_format: %s\n", config.Cfg.OutputFormat)
		fmt.Printf("api_host: %s\n", config.Cfg.APIHost)
		return nil
	}

	key := args[0]
	switch key {
	case "workspace":
		fmt.Println(config.Cfg.DefaultWorkspace)
	case "project":
		fmt.Println(config.Cfg.DefaultProject)
	case "output":
		fmt.Println(config.Cfg.OutputFormat)
	case "api_host":
		fmt.Println(config.Cfg.APIHost)
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}

	return nil
}

func runSet(cmd *cobra.Command, args []string) error {
	key := args[0]
	value := args[1]

	switch key {
	case "workspace":
		config.Cfg.DefaultWorkspace = value
	case "project":
		config.Cfg.DefaultProject = value
	case "output":
		config.Cfg.OutputFormat = value
	case "api_host":
		config.Cfg.APIHost = value
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}

	if err := config.SaveConfig(); err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Set %s to %s", key, value))
	return nil
}
