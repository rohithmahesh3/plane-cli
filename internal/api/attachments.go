package api

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/rohithmahesh3/plane-cli/pkg/plane"
)

// ListAttachments retrieves all attachments for an issue
func (c *Client) ListAttachments(projectID, issueID string) ([]plane.Attachment, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/attachments/", c.Workspace, projectID, issueID)

	var response struct {
		Results []plane.Attachment `json:"results"`
	}

	if err := c.Get(path, nil, &response); err != nil {
		return nil, err
	}

	return response.Results, nil
}

// GetAttachment retrieves a specific attachment
func (c *Client) GetAttachment(projectID, issueID, attachmentID string) (*plane.Attachment, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/attachments/%s/", c.Workspace, projectID, issueID, attachmentID)

	var attachment plane.Attachment
	if err := c.Get(path, nil, &attachment); err != nil {
		return nil, err
	}

	return &attachment, nil
}

// GetUploadCredentials gets credentials for uploading an attachment
func (c *Client) GetUploadCredentials(projectID, issueID, filename string, size int64) (*plane.UploadCredentials, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/attachments/upload-credentials/", c.Workspace, projectID, issueID)

	req := struct {
		Name string `json:"name"`
		Size int64  `json:"size"`
	}{
		Name: filename,
		Size: size,
	}

	var credentials plane.UploadCredentials
	if err := c.Post(path, req, &credentials); err != nil {
		return nil, err
	}

	return &credentials, nil
}

// UploadAttachment uploads a file attachment to an issue
func (c *Client) UploadAttachment(projectID, issueID, filePath string) (*plane.Attachment, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	filename := filepath.Base(filePath)

	// Get upload credentials
	credentials, err := c.GetUploadCredentials(projectID, issueID, filename, fileInfo.Size())
	if err != nil {
		return nil, fmt.Errorf("failed to get upload credentials: %w", err)
	}

	// Upload file to the provided URL (simplified - assumes direct upload)
	// In a real implementation, this would use the credentials to upload to S3 or similar
	_ = credentials

	// For now, we'll use a multipart upload to the Plane API directly
	// This is a simplified implementation
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/attachments/", c.Workspace, projectID, issueID)

	// Create multipart form
	var requestBody io.Reader
	var contentType string

	// Create a pipe to stream the multipart form
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer pw.Close()
		part, err := writer.CreateFormFile("file", filename)
		if err != nil {
			pw.CloseWithError(err)
			return
		}
		if _, err := io.Copy(part, file); err != nil {
			pw.CloseWithError(err)
			return
		}
		writer.Close()
	}()

	requestBody = pr
	contentType = writer.FormDataContentType()

	// Create request
	urlStr := fmt.Sprintf("%s/api/%s%s", c.BaseURL, APIVersion, path)
	req, err := http.NewRequest("POST", urlStr, requestBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-Key", c.APIKey)
	req.Header.Set("Content-Type", contentType)

	// Execute request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("upload failed (status %d): %s", resp.StatusCode, string(body))
	}

	var attachment plane.Attachment
	if err := c.unmarshalResponse(resp.Body, &attachment); err != nil {
		return nil, err
	}

	return &attachment, nil
}

// DeleteAttachment removes an attachment from an issue
func (c *Client) DeleteAttachment(projectID, issueID, attachmentID string) error {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/attachments/%s/", c.Workspace, projectID, issueID, attachmentID)
	return c.Delete(path)
}

// UpdateAttachment updates an existing attachment
func (c *Client) UpdateAttachment(projectID, issueID, attachmentID string, req plane.UpdateAttachmentRequest) (*plane.Attachment, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/attachments/%s/", c.Workspace, projectID, issueID, attachmentID)

	var attachment plane.Attachment
	if err := c.Patch(path, req, &attachment); err != nil {
		return nil, err
	}

	return &attachment, nil
}

// Helper function to unmarshal response
func (c *Client) unmarshalResponse(body io.Reader, v interface{}) error {
	data, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, v)
}
