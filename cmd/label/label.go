package label

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
	labelName        string
	labelDescription string
	labelColor       string
)

var LabelCmd = &cobra.Command{
	Use:     "label",
	Aliases: []string{"labels", "tag"},
	Short:   "Manage project labels",
	Long:    `List, create, edit, and manage labels for categorizing work items in your project.`,
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List labels",
	Long:    `List all labels in the current project.`,
	RunE:    runList,
}

var viewCmd = &cobra.Command{
	Use:   "view <id>",
	Short: "View label details",
	Long:  `Display detailed information about a specific label.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runView,
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new label",
	Long: `Create a new label in the current project.

Examples:
  plane label create --name "Bug" --color "#EF4444"
  plane label create -n "Feature" -c "#3B82F6" -d "New features"`,
	RunE: runCreate,
}

var editCmd = &cobra.Command{
	Use:   "edit <id>",
	Short: "Edit a label",
	Long:  `Edit an existing label.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runEdit,
}

var deleteCmd = &cobra.Command{
	Use:     "delete <id>",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete a label",
	Long:    `Delete a label from the project.`,
	Args:    cobra.ExactArgs(1),
	RunE:    runDelete,
}

func init() {
	LabelCmd.AddCommand(listCmd)
	LabelCmd.AddCommand(viewCmd)
	LabelCmd.AddCommand(createCmd)
	LabelCmd.AddCommand(editCmd)
	LabelCmd.AddCommand(deleteCmd)

	// Create flags
	createCmd.Flags().StringVarP(&labelName, "name", "n", "", "Label name")
	createCmd.Flags().StringVarP(&labelDescription, "description", "d", "", "Label description")
	createCmd.Flags().StringVarP(&labelColor, "color", "c", "", "Label color (hex code, e.g., #EF4444)")

	// Edit flags
	editCmd.Flags().StringVarP(&labelName, "name", "n", "", "New label name")
	editCmd.Flags().StringVarP(&labelDescription, "description", "d", "", "New label description")
	editCmd.Flags().StringVarP(&labelColor, "color", "c", "", "New label color")
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

	labels, err := client.ListLabels(projectID)
	if err != nil {
		return err
	}

	if len(labels) == 0 {
		output.Info("No labels found")
		return nil
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)

	type labelOutput struct {
		ID          string `table:"ID" json:"id"`
		Name        string `table:"NAME" json:"name"`
		Color       string `table:"COLOR" json:"color"`
		Description string `table:"DESCRIPTION" json:"description,omitempty"`
	}

	var outputs []labelOutput
	for _, l := range labels {
		outputs = append(outputs, labelOutput{
			ID:          l.ID,
			Name:        l.Name,
			Color:       l.Color,
			Description: l.Description,
		})
	}

	return formatter.Print(outputs)
}

func runView(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	labelID := args[0]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	label, err := client.GetLabel(projectID, labelID)
	if err != nil {
		return err
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)
	return formatter.Print(label)
}

func runCreate(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified. Use --project flag or set default project")
	}

	// Interactive prompts if flags not provided
	if labelName == "" {
		prompt := &survey.Input{
			Message: "Label name:",
			Help:    "e.g., Bug, Feature, High Priority",
		}
		if err := survey.AskOne(prompt, &labelName); err != nil {
			return err
		}
	}

	if labelName == "" {
		return fmt.Errorf("label name is required")
	}

	if labelColor == "" {
		prompt := &survey.Input{
			Message: "Label color (hex code):",
			Default: "#EF4444",
			Help:    "Hex color code (e.g., #EF4444 for red, #3B82F6 for blue)",
		}
		if err := survey.AskOne(prompt, &labelColor); err != nil {
			return err
		}
	}

	if labelDescription == "" {
		prompt := &survey.Input{
			Message: "Description (optional):",
		}
		_ = survey.AskOne(prompt, &labelDescription)
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	req := plane.CreateLabelRequest{
		Name:        labelName,
		Description: labelDescription,
		Color:       labelColor,
	}

	label, err := client.CreateLabel(projectID, req)
	if err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Created label '%s' (%s)", label.Name, label.ID))
	return nil
}

func runEdit(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	labelID := args[0]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	// Get current label
	label, err := client.GetLabel(projectID, labelID)
	if err != nil {
		return err
	}

	req := plane.UpdateLabelRequest{}

	// Interactive mode if no flags provided
	if labelName == "" && labelDescription == "" && labelColor == "" {
		output.Info(fmt.Sprintf("Editing label: %s", label.Name))

		prompt := &survey.Input{
			Message: "Name:",
			Default: label.Name,
		}
		if err := survey.AskOne(prompt, &req.Name); err != nil {
			return err
		}

		descPrompt := &survey.Input{
			Message: "Description:",
			Default: label.Description,
		}
		if err := survey.AskOne(descPrompt, &req.Description); err != nil {
			return err
		}

		colorPrompt := &survey.Input{
			Message: "Color:",
			Default: label.Color,
		}
		if err := survey.AskOne(colorPrompt, &req.Color); err != nil {
			return err
		}
	} else {
		// Use provided flags
		if labelName != "" {
			req.Name = labelName
		}
		if labelDescription != "" {
			req.Description = labelDescription
		}
		if labelColor != "" {
			req.Color = labelColor
		}
	}

	updatedLabel, err := client.UpdateLabel(projectID, labelID, req)
	if err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Updated label '%s'", updatedLabel.Name))
	return nil
}

func runDelete(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	labelID := args[0]

	// Confirm deletion
	var confirm bool
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("Are you sure you want to delete label %s?", labelID),
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

	if err := client.DeleteLabel(projectID, labelID); err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Deleted label %s", labelID))
	return nil
}
