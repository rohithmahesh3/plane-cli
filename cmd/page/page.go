package page

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
	pageName        string
	pageDescription string
	isWorkspacePage bool
)

var PageCmd = &cobra.Command{
	Use:     "page",
	Aliases: []string{"pages", "doc", "wiki"},
	Short:   "Manage pages (documentation)",
	Long:    `Create and view documentation pages at workspace or project level.`,
}

var viewCmd = &cobra.Command{
	Use:   "view <id>",
	Short: "View page details",
	Long:  `Display detailed information about a specific page.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runView,
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new page",
	Long: `Create a new documentation page.

Examples:
  plane page create --name "API Documentation"
  plane page create -n "Getting Started" --workspace`,
	RunE: runCreate,
}

func init() {
	PageCmd.AddCommand(viewCmd)
	PageCmd.AddCommand(createCmd)

	// Create flags
	createCmd.Flags().StringVarP(&pageName, "name", "n", "", "Page name")
	createCmd.Flags().StringVarP(&pageDescription, "description", "d", "", "Page content/description")
	createCmd.Flags().BoolVarP(&isWorkspacePage, "workspace", "w", false, "Create as workspace page")
}

func runView(cmd *cobra.Command, args []string) error {
	pageID := args[0]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	var page *plane.Page

	// Try workspace page first, then project page
	page, err = client.GetWorkspacePage(pageID)
	if err != nil {
		// Try project page
		projectID := config.Cfg.DefaultProject
		if projectID == "" {
			return fmt.Errorf("no project specified")
		}
		page, err = client.GetProjectPage(projectID, pageID)
		if err != nil {
			return err
		}
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)
	return formatter.Print(page)
}

func runCreate(cmd *cobra.Command, args []string) error {
	// Interactive prompts if flags not provided
	if pageName == "" {
		prompt := &survey.Input{
			Message: "Page name:",
			Help:    "e.g., API Documentation, Getting Started Guide",
		}
		if err := survey.AskOne(prompt, &pageName); err != nil {
			return err
		}
	}

	if pageName == "" {
		return fmt.Errorf("page name is required")
	}

	if pageDescription == "" {
		prompt := &survey.Editor{
			Message:       "Page content:",
			FileName:      "*.md",
			HideDefault:   true,
			AppendDefault: true,
		}
		if err := survey.AskOne(prompt, &pageDescription); err != nil {
			return err
		}
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	req := plane.CreatePageRequest{
		Name:            pageName,
		DescriptionHTML: "<p>" + pageDescription + "</p>",
	}

	var page *plane.Page
	if isWorkspacePage {
		page, err = client.CreateWorkspacePage(req)
	} else {
		projectID := config.Cfg.DefaultProject
		if projectID == "" {
			return fmt.Errorf("no project specified")
		}
		page, err = client.CreateProjectPage(projectID, req)
	}

	if err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Created page '%s' (%s)", page.Name, page.ID))
	return nil
}
