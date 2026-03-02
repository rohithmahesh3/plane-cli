package issue

import (
	"fmt"
	"os"

	"github.com/rohithmahesh3/plane-cli/internal/api"
	"github.com/rohithmahesh3/plane-cli/internal/config"
	"github.com/rohithmahesh3/plane-cli/internal/output"
	"github.com/rohithmahesh3/plane-cli/pkg/plane"
	"github.com/spf13/cobra"
)

var (
	attachmentName     string
	attachmentArchived bool
)

func init() {
	// Add attachment subcommand
	attachmentCmd := &cobra.Command{
		Use:     "attachment",
		Aliases: []string{"attach", "file"},
		Short:   "Manage issue attachments",
		Long:    `Upload, list, and manage file attachments on issues.`,
	}

	attachmentListCmd := &cobra.Command{
		Use:     "list <issue-id>",
		Aliases: []string{"ls"},
		Short:   "List attachments for an issue",
		Args:    cobra.ExactArgs(1),
		RunE:    runAttachmentList,
	}

	attachmentUploadCmd := &cobra.Command{
		Use:   "upload <issue-id> <file-path>",
		Short: "Upload a file attachment",
		Long:  `Upload a file as an attachment to an issue.`,
		Args:  cobra.ExactArgs(2),
		RunE:  runAttachmentUpload,
	}

	attachmentEditCmd := &cobra.Command{
		Use:   "edit <issue-id> <attachment-id>",
		Short: "Edit attachment metadata",
		Long:  `Update attachment properties like name and archive status.`,
		Args:  cobra.ExactArgs(2),
		RunE:  runAttachmentEdit,
	}

	attachmentDeleteCmd := &cobra.Command{
		Use:     "delete <issue-id> <attachment-id>",
		Aliases: []string{"rm", "remove"},
		Short:   "Delete an attachment",
		Args:    cobra.ExactArgs(2),
		RunE:    runAttachmentDelete,
	}

	attachmentEditCmd.Flags().StringVarP(&attachmentName, "name", "n", "", "New filename")
	attachmentEditCmd.Flags().BoolVarP(&attachmentArchived, "archive", "a", false, "Archive the attachment")
	attachmentEditCmd.Flags().BoolVarP(&attachmentArchived, "unarchive", "u", false, "Unarchive the attachment")

	attachmentCmd.AddCommand(attachmentListCmd)
	attachmentCmd.AddCommand(attachmentUploadCmd)
	attachmentCmd.AddCommand(attachmentEditCmd)
	attachmentCmd.AddCommand(attachmentDeleteCmd)

	IssueCmd.AddCommand(attachmentCmd)
}

func runAttachmentList(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	issueID := args[0]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	issueID, err = resolveIssueID(client, projectID, issueID)
	if err != nil {
		return err
	}

	attachments, err := client.ListAttachments(projectID, issueID)
	if err != nil {
		return err
	}

	if len(attachments) == 0 {
		output.Info("No attachments found for this issue")
		return nil
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)

	type attachmentOutput struct {
		ID       string `table:"ID" json:"id"`
		Name     string `table:"NAME" json:"name"`
		Size     string `table:"SIZE" json:"size_formatted"`
		Type     string `table:"TYPE" json:"type"`
		Uploaded string `table:"UPLOADED" json:"uploaded_at"`
	}

	var outputs []attachmentOutput
	for _, a := range attachments {
		outputs = append(outputs, attachmentOutput{
			ID:       a.ID,
			Name:     a.Attributes.Name,
			Size:     formatBytes(a.Attributes.Size),
			Type:     a.Attributes.Type,
			Uploaded: a.CreatedAt.Format("2006-01-02"),
		})
	}

	return formatter.Print(outputs)
}

func runAttachmentUpload(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	issueID := args[0]
	filePath := args[1]

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", filePath)
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	issueID, err = resolveIssueID(client, projectID, issueID)
	if err != nil {
		return err
	}

	attachment, err := client.UploadAttachment(projectID, issueID, filePath)
	if err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Uploaded '%s' (%s)", attachment.Attributes.Name, formatBytes(attachment.Attributes.Size)))
	return nil
}

func runAttachmentEdit(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	issueID := args[0]
	attachmentID := args[1]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	issueID, err = resolveIssueID(client, projectID, issueID)
	if err != nil {
		return err
	}

	req := plane.UpdateAttachmentRequest{}

	if attachmentName != "" {
		req.Attributes.Name = attachmentName
	}

	// Handle archive/unarchive flags
	if cmd.Flags().Changed("archive") {
		req.IsArchived = true
	} else if cmd.Flags().Changed("unarchive") {
		req.IsArchived = false
	}

	attachment, err := client.UpdateAttachment(projectID, issueID, attachmentID, req)
	if err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Updated attachment '%s'", attachment.Attributes.Name))
	return nil
}

func runAttachmentDelete(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	issueID := args[0]
	attachmentID := args[1]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	issueID, err = resolveIssueID(client, projectID, issueID)
	if err != nil {
		return err
	}

	if err := client.DeleteAttachment(projectID, issueID, attachmentID); err != nil {
		return err
	}

	output.Success("Attachment deleted")
	return nil
}

// formatBytes converts bytes to human-readable format
func formatBytes(bytes int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
