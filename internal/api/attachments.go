package api

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/rohithmahesh3/plane-cli/pkg/plane"
)

// ListAttachments retrieves all attachments for an issue
func (c *Client) ListAttachments(projectID, issueID string) ([]plane.Attachment, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/attachments/", c.Workspace, projectID, issueID)

	body, err := c.GetRaw(path, nil)
	if err != nil {
		return nil, err
	}

	return unmarshalListResponse[plane.Attachment](body)
}

// GetAttachment retrieves a specific attachment
func (c *Client) GetAttachment(projectID, issueID, attachmentID string) (*plane.Attachment, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/attachments/%s/", c.Workspace, projectID, issueID, attachmentID)

	var attachment plane.Attachment
	if err := c.Get(path, nil, &attachment); err == nil {
		return &attachment, nil
	}

	attachments, err := c.ListAttachments(projectID, issueID)
	if err != nil {
		return nil, err
	}

	for _, item := range attachments {
		if item.ID == attachmentID {
			return &item, nil
		}
	}

	return nil, fmt.Errorf("attachment %s not found", attachmentID)
}

// GetUploadCredentials gets credentials for uploading an attachment
func (c *Client) GetUploadCredentials(projectID, issueID, filename string, size int64) (*plane.UploadCredentials, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/attachments/", c.Workspace, projectID, issueID)

	contentType := mime.TypeByExtension(filepath.Ext(filename))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	req := struct {
		Name string `json:"name"`
		Type string `json:"type"`
		Size int64  `json:"size"`
	}{
		Name: filename,
		Type: contentType,
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
	defer func() {
		_ = file.Close()
	}()

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

	uploadTarget := credentials.UploadTarget()
	if uploadTarget.URL == "" {
		return nil, fmt.Errorf("upload target missing from credentials response")
	}

	requestBody, contentType, err := buildMultipartPayload(uploadTarget.Fields, "file", filename, file)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uploadTarget.URL, requestBody)
	if err != nil {
		return nil, err
	}

	if seeker, ok := requestBody.(io.Seeker); ok {
		size, err := seeker.Seek(0, io.SeekEnd)
		if err != nil {
			return nil, err
		}
		if _, err := seeker.Seek(0, io.SeekStart); err != nil {
			return nil, err
		}
		req.ContentLength = size
	}
	req.Header.Set("Content-Type", contentType)

	// Execute request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("upload failed (status %d): %s", resp.StatusCode, string(body))
	}

	attachmentID := credentials.AssetID
	if attachmentID == "" {
		attachmentID = credentials.Attachment.ID
	}
	if attachmentID == "" {
		return nil, fmt.Errorf("attachment ID missing from credentials response")
	}

	attachment, err := c.completeAttachmentUpload(projectID, issueID, attachmentID)
	if err != nil {
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

func (c *Client) completeAttachmentUpload(projectID, issueID, attachmentID string) (plane.Attachment, error) {
	attachmentPath := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/attachments/%s/", c.Workspace, projectID, issueID, attachmentID)

	// Plane documents upload completion as a PATCH to the attachment resource.
	attempts := []struct {
		method string
		path   string
		body   interface{}
	}{
		{
			method: http.MethodPatch,
			path:   attachmentPath,
			body:   map[string]bool{"is_uploaded": true},
		},
		{
			method: http.MethodPost,
			path:   attachmentPath + "complete-upload/",
			body:   map[string]string{"asset_id": attachmentID},
		},
		{
			method: http.MethodPost,
			path:   fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/attachments/complete-upload/", c.Workspace, projectID, issueID),
			body:   map[string]string{"asset_id": attachmentID},
		},
	}

	var lastErr error
	for _, attempt := range attempts {
		var attachment plane.Attachment

		switch attempt.method {
		case http.MethodPost:
			lastErr = c.Post(attempt.path, attempt.body, &attachment)
		case http.MethodPatch:
			lastErr = c.Patch(attempt.path, attempt.body, &attachment)
		default:
			lastErr = fmt.Errorf("unsupported method %s", attempt.method)
		}

		if lastErr == nil {
			if attachment.ID != "" {
				return attachment, nil
			}

			fetched, err := c.GetAttachment(projectID, issueID, attachmentID)
			if err == nil {
				return *fetched, nil
			}
			lastErr = err
			continue
		}
	}

	for i := 0; i < 5; i++ {
		attachment, err := c.GetAttachment(projectID, issueID, attachmentID)
		if err == nil && attachment.ID != "" {
			return *attachment, nil
		}
		time.Sleep(1 * time.Second)
	}

	return plane.Attachment{}, lastErr
}

func buildMultipartPayload(fields map[string]string, fileField, filename string, file io.Reader) (io.Reader, string, error) {
	var payload bytes.Buffer
	writer := multipart.NewWriter(&payload)

	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			return nil, "", err
		}
	}

	part, err := writer.CreateFormFile(fileField, filename)
	if err != nil {
		return nil, "", err
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, "", err
	}

	if err := writer.Close(); err != nil {
		return nil, "", err
	}

	return bytes.NewReader(payload.Bytes()), writer.FormDataContentType(), nil
}
