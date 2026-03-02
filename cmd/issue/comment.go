package issue

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
	commentText   string
	commentAccess string
)

func init() {
	// Add comment subcommand
	commentCmd := &cobra.Command{
		Use:     "comment",
		Aliases: []string{"comments"},
		Short:   "Manage issue comments",
		Long:    `Add, list, and manage comments on issues.`,
	}

	commentListCmd := &cobra.Command{
		Use:     "list <issue-id>",
		Aliases: []string{"ls"},
		Short:   "List comments for an issue",
		Args:    cobra.ExactArgs(1),
		RunE:    runCommentList,
	}

	commentAddCmd := &cobra.Command{
		Use:   "add <issue-id>",
		Short: "Add a comment to an issue",
		Args:  cobra.ExactArgs(1),
		RunE:  runCommentAdd,
	}

	commentDeleteCmd := &cobra.Command{
		Use:     "delete <issue-id> <comment-id>",
		Aliases: []string{"rm", "remove"},
		Short:   "Delete a comment from an issue",
		Args:    cobra.ExactArgs(2),
		RunE:    runCommentDelete,
	}

	commentAddCmd.Flags().StringVarP(&commentText, "text", "t", "", "Comment text")
	commentAddCmd.Flags().StringVar(&commentAccess, "access", "INTERNAL", "Comment access (INTERNAL or EXTERNAL)")

	commentCmd.AddCommand(commentListCmd)
	commentCmd.AddCommand(commentAddCmd)
	commentCmd.AddCommand(commentDeleteCmd)

	IssueCmd.AddCommand(commentCmd)
}

func runCommentList(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	issueID := args[0]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	comments, err := client.ListComments(projectID, issueID)
	if err != nil {
		return err
	}

	if len(comments) == 0 {
		output.Info("No comments found for this issue")
		return nil
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)

	type commentOutput struct {
		ID      string `table:"ID" json:"id"`
		Author  string `table:"AUTHOR" json:"actor"`
		Access  string `table:"ACCESS" json:"access"`
		Comment string `table:"COMMENT" json:"comment_stripped"`
	}

	var outputs []commentOutput
	for _, c := range comments {
		// Strip HTML for display
		commentText := c.CommentStripped
		if commentText == "" {
			commentText = "(empty)"
		}
		if len(commentText) > 50 {
			commentText = commentText[:47] + "..."
		}

		access := c.Access
		if access == "" {
			access = "INTERNAL"
		}

		outputs = append(outputs, commentOutput{
			ID:      c.ID,
			Author:  c.Actor,
			Access:  access,
			Comment: commentText,
		})
	}

	return formatter.Print(outputs)
}

func runCommentAdd(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	issueID := args[0]

	// Interactive prompts if flags not provided
	if commentText == "" {
		prompt := &survey.Editor{
			Message:       "Comment:",
			FileName:      "*.md",
			HideDefault:   true,
			AppendDefault: true,
		}
		if err := survey.AskOne(prompt, &commentText); err != nil {
			return err
		}
	}

	if commentText == "" {
		return fmt.Errorf("comment text is required")
	}

	// Convert plain text to simple HTML
	commentHTML := "<p>" + commentText + "</p>"

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	req := plane.CreateCommentRequest{
		CommentHTML: commentHTML,
		Access:      commentAccess,
	}

	comment, err := client.CreateComment(projectID, issueID, req)
	if err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Added comment to issue (%s)", comment.ID))
	return nil
}

func runCommentDelete(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	issueID := args[0]
	commentID := args[1]

	// Confirm deletion
	var confirm bool
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("Are you sure you want to delete comment %s?", commentID),
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

	if err := client.DeleteComment(projectID, issueID, commentID); err != nil {
		return err
	}

	output.Success("Comment deleted")
	return nil
}
