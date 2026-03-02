package issuetype

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
	typeName        string
	typeDescription string
)

var TypeCmd = &cobra.Command{
	Use:     "type",
	Aliases: []string{"issue-type"},
	Short:   "Manage issue types",
	Long:    `List, create, and manage custom issue types for your project.`,
}

var typeListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List issue types",
	RunE:    runTypeList,
}

var typeCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new issue type",
	RunE:  runTypeCreate,
}

var typeDeleteCmd = &cobra.Command{
	Use:     "delete <type-id>",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete an issue type",
	Args:    cobra.ExactArgs(1),
	RunE:    runTypeDelete,
}

func init() {
	typeCreateCmd.Flags().StringVarP(&typeName, "name", "n", "", "Issue type name")
	typeCreateCmd.Flags().StringVarP(&typeDescription, "description", "d", "", "Issue type description")

	TypeCmd.AddCommand(typeListCmd)
	TypeCmd.AddCommand(typeCreateCmd)
	TypeCmd.AddCommand(typeDeleteCmd)
}

func runTypeList(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified. Use --project flag or set default project")
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	project, err := client.GetProject(projectID)
	if err != nil {
		return err
	}
	if !project.IsIssueTypeEnabled {
		return fmt.Errorf("issue types are disabled for project %s", projectID)
	}

	types, err := client.ListIssueTypes(projectID)
	if err != nil {
		return err
	}

	if len(types) == 0 {
		output.Info("No issue types found")
		return nil
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)

	type typeOutput struct {
		ID          string `table:"ID" json:"id"`
		Name        string `table:"NAME" json:"name"`
		Description string `table:"DESCRIPTION" json:"description"`
		IsDefault   string `table:"DEFAULT" json:"is_default"`
	}

	var outputs []typeOutput
	for _, t := range types {
		isDefault := ""
		if t.IsDefault {
			isDefault = "✓"
		}
		outputs = append(outputs, typeOutput{
			ID:          t.ID,
			Name:        t.Name,
			Description: t.Description,
			IsDefault:   isDefault,
		})
	}

	return formatter.Print(outputs)
}

func runTypeCreate(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	// Interactive prompts if flags not provided
	if typeName == "" {
		prompt := &survey.Input{
			Message: "Issue type name:",
			Help:    "e.g., Bug, Feature, Task",
		}
		if err := survey.AskOne(prompt, &typeName); err != nil {
			return err
		}
	}

	if typeName == "" {
		return fmt.Errorf("name is required")
	}

	if typeDescription == "" {
		prompt := &survey.Input{
			Message: "Description (optional):",
		}
		_ = survey.AskOne(prompt, &typeDescription)
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	project, err := client.GetProject(projectID)
	if err != nil {
		return err
	}
	if !project.IsIssueTypeEnabled {
		return fmt.Errorf("issue types are disabled for project %s", projectID)
	}

	req := plane.CreateIssueTypeRequest{
		Name:        typeName,
		Description: typeDescription,
		IsActive:    true,
	}

	issueType, err := client.CreateIssueType(projectID, req)
	if err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Created issue type '%s' (%s)", issueType.Name, issueType.ID))
	return nil
}

func runTypeDelete(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	typeID := args[0]

	// Confirm deletion
	var confirm bool
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("Are you sure you want to delete issue type %s?", typeID),
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

	project, err := client.GetProject(projectID)
	if err != nil {
		return err
	}
	if !project.IsIssueTypeEnabled {
		return fmt.Errorf("issue types are disabled for project %s", projectID)
	}

	if err := client.DeleteIssueType(projectID, typeID); err != nil {
		return err
	}

	output.Success("Issue type deleted")
	return nil
}
