package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rohithmahesh3/plane-cli/pkg/plane"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name      string
		apiKey    string
		workspace string
		baseURL   string
		wantErr   bool
	}{
		{
			name:      "valid client",
			apiKey:    "test-api-key",
			workspace: "test-workspace",
			baseURL:   "https://api.plane.so",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				HTTPClient: &http.Client{Timeout: DefaultTimeout},
				BaseURL:    tt.baseURL,
				APIKey:     tt.apiKey,
				Workspace:  tt.workspace,
			}

			assert.NotNil(t, client)
			assert.Equal(t, tt.apiKey, client.APIKey)
			assert.Equal(t, tt.workspace, client.Workspace)
			assert.Equal(t, tt.baseURL, client.BaseURL)
		})
	}
}

func TestClient_NewRequest(t *testing.T) {
	client := &Client{
		HTTPClient: &http.Client{Timeout: DefaultTimeout},
		BaseURL:    "https://api.plane.so",
		APIKey:     "test-api-key",
		Workspace:  "test-workspace",
	}

	tests := []struct {
		name       string
		method     string
		path       string
		body       interface{}
		wantErr    bool
		wantHeader map[string]string
	}{
		{
			name:   "GET request without body",
			method: "GET",
			path:   "/workspaces/",
			body:   nil,
			wantErr: false,
			wantHeader: map[string]string{
				"X-API-Key":    "test-api-key",
				"Content-Type": "application/json",
				"Accept":       "application/json",
			},
		},
		{
			name:   "POST request with body",
			method: "POST",
			path:   "/workspaces/test-workspace/projects/",
			body: map[string]string{
				"name":       "Test Project",
				"identifier": "TEST",
			},
			wantErr: false,
			wantHeader: map[string]string{
				"X-API-Key":    "test-api-key",
				"Content-Type": "application/json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := client.NewRequest(tt.method, tt.path, tt.body)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.method, req.Method)

			for key, value := range tt.wantHeader {
				assert.Equal(t, value, req.Header.Get(key))
			}
		})
	}
}

func TestClient_Do(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		response   interface{}
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "successful response",
			statusCode: http.StatusOK,
			response:   map[string]string{"id": "123", "name": "Test"},
			wantErr:    false,
		},
		{
			name:       "error response",
			statusCode: http.StatusUnauthorized,
			response:   map[string]string{"error": "Unauthorized"},
			wantErr:    true,
			errMsg:     "API error (status 401)",
		},
		{
			name:       "bad request",
			statusCode: http.StatusBadRequest,
			response:   map[string]string{"error": "Bad Request"},
			wantErr:    true,
			errMsg:     "API error (status 400)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				json.NewEncoder(w).Encode(tt.response)
			}))
			defer server.Close()

			client := &Client{
				HTTPClient: &http.Client{Timeout: DefaultTimeout},
				BaseURL:    server.URL,
				APIKey:     "test-key",
			}

			req, err := client.NewRequest("GET", "/test", nil)
			require.NoError(t, err)

			var result map[string]string
			err = client.Do(req, &result)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.response, result)
			}
		})
	}
}

func TestClient_ListWorkspaces(t *testing.T) {
	mockWorkspaces := []plane.Workspace{
		{
			ID:   "ws-1",
			Name: "Workspace 1",
			Slug: "workspace-1",
		},
		{
			ID:   "ws-2",
			Name: "Workspace 2",
			Slug: "workspace-2",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/workspaces/", r.URL.Path)
		assert.Equal(t, "test-api-key", r.Header.Get("X-API-Key"))

		response := Response{
			Results: mustMarshal(mockWorkspaces),
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: &http.Client{Timeout: DefaultTimeout},
		BaseURL:    server.URL,
		APIKey:     "test-api-key",
	}

	workspaces, err := client.ListWorkspaces()
	require.NoError(t, err)
	assert.Len(t, workspaces, 2)
	assert.Equal(t, "Workspace 1", workspaces[0].Name)
}

func TestClient_ListProjects(t *testing.T) {
	mockProjects := []plane.Project{
		{
			ID:         "proj-1",
			Name:       "Project 1",
			Identifier: "PROJ1",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/workspaces/test-workspace/projects/", r.URL.Path)

		response := Response{
			Results: mustMarshal(mockProjects),
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: &http.Client{Timeout: DefaultTimeout},
		BaseURL:    server.URL,
		APIKey:     "test-api-key",
		Workspace:  "test-workspace",
	}

	projects, err := client.ListProjects()
	require.NoError(t, err)
	assert.Len(t, projects, 1)
	assert.Equal(t, "Project 1", projects[0].Name)
}

func TestClient_ListIssues(t *testing.T) {
	mockIssues := []plane.Issue{
		{
			ID:         "issue-1",
			SequenceID: 1,
			Name:       "Test Issue",
			Priority:   "high",
			State:      "backlog-state-id",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/workspaces/test-workspace/projects/proj-1/issues/", r.URL.Path)

		// Check query params
		query := r.URL.Query()
		assert.Equal(t, "backlog", query.Get("state"))
		assert.Equal(t, "high", query.Get("priority"))

		response := Response{
			Results:       mustMarshal(mockIssues),
			Pagination: Pagination{
				NextCursor:      "20:1:0",
				NextPageResults: false,
				Count:           1,
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: &http.Client{Timeout: DefaultTimeout},
		BaseURL:    server.URL,
		APIKey:     "test-api-key",
		Workspace:  "test-workspace",
	}

	opts := IssueListOptions{
		State:    "backlog",
		Priority: "high",
	}

	issues, pagination, err := client.ListIssues("proj-1", opts)
	require.NoError(t, err)
	assert.Len(t, issues, 1)
	assert.Equal(t, "Test Issue", issues[0].Name)
	assert.NotNil(t, pagination)
	assert.Equal(t, "20:1:0", pagination.NextCursor)
}

func TestClient_CreateIssue(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/api/v1/workspaces/test-workspace/projects/proj-1/issues/", r.URL.Path)

		var req plane.CreateIssueRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)

		assert.Equal(t, "New Issue", req.Name)
		assert.Equal(t, "high", req.Priority)

		response := plane.Issue{
			ID:         "new-issue-id",
			SequenceID: 42,
			Name:       req.Name,
			Priority:   req.Priority,
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: &http.Client{Timeout: DefaultTimeout},
		BaseURL:    server.URL,
		APIKey:     "test-api-key",
		Workspace:  "test-workspace",
	}

	req := plane.CreateIssueRequest{
		Name:     "New Issue",
		Priority: "high",
	}

	issue, err := client.CreateIssue("proj-1", req)
	require.NoError(t, err)
	assert.Equal(t, "new-issue-id", issue.ID)
	assert.Equal(t, 42, issue.SequenceID)
}

func TestClient_DeleteIssue(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "/api/v1/workspaces/test-workspace/projects/proj-1/issues/issue-1/", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: &http.Client{Timeout: DefaultTimeout},
		BaseURL:    server.URL,
		APIKey:     "test-api-key",
		Workspace:  "test-workspace",
	}

	err := client.DeleteIssue("proj-1", "issue-1")
	require.NoError(t, err)
}

func mustMarshal(v interface{}) []byte {
	data, _ := json.Marshal(v)
	return data
}
