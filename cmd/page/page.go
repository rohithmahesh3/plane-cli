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
	Long:    `List, create, and manage documentation pages at workspace or project level.`,
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List pages",
	Long:    `List all pages. Use --workspace to list workspace pages, otherwise lists project pages.`,
	RunE:    runList,
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

var deleteCmd = &cobra.Command{
	Use:     "delete <id>",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete a page",
	Long:    `Delete a page from the project or workspace.`,
	Args:    cobra.ExactArgs(1),
	RunE:    runDelete,
}

func init() {
	PageCmd.AddCommand(listCmd)
	PageCmd.AddCommand(viewCmd)
	PageCmd.AddCommand(createCmd)
	PageCmd.AddCommand(deleteCmd)

	// List flags
	listCmd.Flags().BoolVarP(&isWorkspacePage, "workspace", "w", false, "List workspace pages instead of project pages")

	// Create flags
	createCmd.Flags().StringVarP(&pageName, "name", "n", "", "Page name")
	createCmd.Flags().StringVarP(&pageDescription, "description", "d", "", "Page content/description")
	createCmd.Flags().BoolVarP(&isWorkspacePage, "workspace", "w", false, "Create as workspace page")
}

func runList(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient()
	if err != nil {
		return err
	}

	var pages []plane.Page

	if isWorkspacePage {
		pages, err = client.ListWorkspacePages()
	} else {
		projectID := config.Cfg.DefaultProject
		if projectID == "" {
			return fmt.Errorf("no project specified. Use --project flag or set default project")
		}
		pages, err = client.ListProjectPages(projectID)
	}

	if err != nil {
		return err
	}

	if len(pages) == 0 {
		output.Info("No pages found")
		return nil
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)

	type pageOutput struct {
		ID      string `table:"ID" json:"id"`
		Name    string `table:"NAME" json:"name"`
		Updated string `table:"UPDATED" json:"updated_at_formatted"`
	}

	var outputs []pageOutput
	for _, p := range pages {
		outputs = append(outputs, pageOutput{
			ID:      p.ID,
			Name:    p.Name,
			Updated: p.UpdatedAt.Format("2006-01-02"),
		})
	}

	return formatter.Print(outputs)
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

func runDelete(cmd *cobra.Command, args []string) error {
	pageID := args[0]

	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	if err := client.DeletePage(projectID, pageID); err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Deleted page %s", pageID))
	return nil
}
