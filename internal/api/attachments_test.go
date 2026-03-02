package api

import (
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/rohithmahesh3/plane-cli/pkg/plane"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUploadCredentials(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/v1/workspaces/test-workspace/projects/test-project/work-items/issue-123/attachments/", r.URL.Path)

		var req struct {
			Name string `json:"name"`
			Type string `json:"type"`
			Size int64  `json:"size"`
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)

		assert.Equal(t, "test-file.txt", req.Name)
		assert.Equal(t, "text/plain; charset=utf-8", req.Type)
		assert.EqualValues(t, 1024, req.Size)

		response := plane.UploadCredentials{
			UploadData: plane.UploadData{
				URL:    "https://uploads.example.com",
				Fields: map[string]string{"key": "attachments/test-file.txt"},
			},
			AssetID: "attachment-123",
			Attachment: plane.Attachment{
				ID: "attachment-123",
			},
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(response)
		require.NoError(t, err)
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
		APIKey:     "test-key",
		Workspace:  "test-workspace",
	}

	creds, err := client.GetUploadCredentials("test-project", "issue-123", "test-file.txt", 1024)
	require.NoError(t, err)
	assert.Equal(t, "https://uploads.example.com", creds.UploadTarget().URL)
	assert.Equal(t, "attachment-123", creds.Attachment.ID)
}

func TestUploadAttachment(t *testing.T) {
	uploadServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		err := r.ParseMultipartForm(1024)
		require.NoError(t, err)

		assert.Equal(t, "attachments/upload.txt", r.FormValue("key"))

		file, _, err := r.FormFile("file")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, file.Close())
		}()

		data, err := io.ReadAll(file)
		require.NoError(t, err)
		assert.Equal(t, "hello upload", string(data))

		w.WriteHeader(http.StatusNoContent)
	}))
	defer uploadServer.Close()

	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/workspaces/test-workspace/projects/test-project/work-items/issue-123/attachments/":
			var req struct {
				Name string `json:"name"`
				Type string `json:"type"`
				Size int64  `json:"size"`
			}
			err := json.NewDecoder(r.Body).Decode(&req)
			require.NoError(t, err)
			assert.Equal(t, "upload.txt", req.Name)
			assert.Equal(t, mime.TypeByExtension(".txt"), req.Type)

			err = json.NewEncoder(w).Encode(plane.UploadCredentials{
				UploadData: plane.UploadData{
					URL:    uploadServer.URL,
					Fields: map[string]string{"key": "attachments/upload.txt"},
				},
				AssetID: "attachment-456",
				Attachment: plane.Attachment{
					ID: "attachment-456",
				},
			})
			require.NoError(t, err)
		case r.Method == http.MethodPatch && r.URL.Path == "/api/v1/workspaces/test-workspace/projects/test-project/work-items/issue-123/attachments/attachment-456/":
			var req map[string]bool
			err := json.NewDecoder(r.Body).Decode(&req)
			require.NoError(t, err)
			assert.True(t, req["is_uploaded"])

			err = json.NewEncoder(w).Encode(plane.Attachment{
				ID:         "attachment-456",
				IsUploaded: true,
				Attributes: plane.AttachmentAttributes{
					Name: "upload.txt",
					Size: 12,
					Type: "text/plain",
				},
			})
			require.NoError(t, err)
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	}))
	defer apiServer.Close()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "upload.txt")
	err := os.WriteFile(filePath, []byte("hello upload"), 0o600)
	require.NoError(t, err)

	client := &Client{
		HTTPClient: apiServer.Client(),
		BaseURL:    apiServer.URL,
		APIKey:     "test-key",
		Workspace:  "test-workspace",
	}

	attachment, err := client.UploadAttachment("test-project", "issue-123", filePath)
	require.NoError(t, err)
	assert.Equal(t, "attachment-456", attachment.ID)
	assert.True(t, attachment.IsUploaded)
	assert.Equal(t, "upload.txt", attachment.Attributes.Name)
}

func TestGetAttachment(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/v1/workspaces/test-workspace/projects/test-project/work-items/issue-123/attachments/attachment-456/", r.URL.Path)

		err := json.NewEncoder(w).Encode(plane.Attachment{
			ID:         "attachment-456",
			IsUploaded: true,
		})
		require.NoError(t, err)
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
		APIKey:     "test-key",
		Workspace:  "test-workspace",
	}

	attachment, err := client.GetAttachment("test-project", "issue-123", "attachment-456")
	require.NoError(t, err)
	assert.Equal(t, "attachment-456", attachment.ID)
	assert.True(t, attachment.IsUploaded)
}
