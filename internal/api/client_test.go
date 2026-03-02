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
			name:    "GET request without body",
			method:  "GET",
			path:    "/workspaces/",
			body:    nil,
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
				err := json.NewEncoder(w).Encode(tt.response)
				require.NoError(t, err)
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
			Results: mustMarshal(t, mockProjects),
		}
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(response)
		require.NoError(t, err)
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
			State: plane.FlexibleState{
				ID: "backlog-state-id",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/workspaces/test-workspace/projects/proj-1/work-items/", r.URL.Path)

		// Check query params
		query := r.URL.Query()
		assert.Equal(t, "backlog", query.Get("state"))
		assert.Equal(t, "25", query.Get("limit"))
		assert.Equal(t, "", query.Get("per_page"))
		assert.Equal(t, "", query.Get("priority"))
		assert.Equal(t, "", query.Get("cursor"))

		response := Response{
			Results: mustMarshal(t, mockIssues),
			Pagination: Pagination{
				NextCursor:      "20:1:0",
				NextPageResults: false,
				Count:           1,
			},
		}
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(response)
		require.NoError(t, err)
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: &http.Client{Timeout: DefaultTimeout},
		BaseURL:    server.URL,
		APIKey:     "test-api-key",
		Workspace:  "test-workspace",
	}

	opts := IssueListOptions{
		State: "backlog",
		Limit: 25,
	}

	issues, pagination, err := client.ListIssues("proj-1", opts)
	require.NoError(t, err)
	assert.Len(t, issues, 1)
	assert.Equal(t, "Test Issue", issues[0].Name)
	assert.NotNil(t, pagination)
	assert.Equal(t, "20:1:0", pagination.NextCursor)
}

func TestClient_GetIssueByIdentifier(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/api/v1/workspaces/test-workspace/work-items/TESTW-42/", r.URL.Path)
		assert.Equal(t, "assignees,state,labels", r.URL.Query().Get("expand"))

		response := plane.Issue{
			ID:         "issue-42",
			SequenceID: 42,
			Name:       "Issue by identifier",
		}

		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(response)
		require.NoError(t, err)
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: &http.Client{Timeout: DefaultTimeout},
		BaseURL:    server.URL,
		APIKey:     "test-api-key",
		Workspace:  "test-workspace",
	}

	issue, err := client.GetIssueByIdentifier("TESTW-42")
	require.NoError(t, err)
	assert.Equal(t, "issue-42", issue.ID)
	assert.Equal(t, 42, issue.SequenceID)
}

func TestClient_GetIssueBySequenceID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/workspaces/test-workspace/projects/proj-1/":
			err := json.NewEncoder(w).Encode(plane.Project{
				ID:         "proj-1",
				Identifier: "TESTW",
				Name:       "Test Workspace",
			})
			require.NoError(t, err)
		case "/api/v1/workspaces/test-workspace/work-items/TESTW-7/":
			assert.Equal(t, "assignees,state,labels", r.URL.Query().Get("expand"))
			err := json.NewEncoder(w).Encode(plane.Issue{
				ID:         "issue-7",
				SequenceID: 7,
				Name:       "Sequence issue",
			})
			require.NoError(t, err)
		default:
			t.Fatalf("unexpected request path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: &http.Client{Timeout: DefaultTimeout},
		BaseURL:    server.URL,
		APIKey:     "test-api-key",
		Workspace:  "test-workspace",
	}

	issue, err := client.GetIssueBySequenceID("proj-1", 7)
	require.NoError(t, err)
	assert.Equal(t, "issue-7", issue.ID)
}

func TestClient_CreateIssue(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/api/v1/workspaces/test-workspace/projects/proj-1/work-items/", r.URL.Path)

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
		err = json.NewEncoder(w).Encode(response)
		require.NoError(t, err)
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
		assert.Equal(t, "/api/v1/workspaces/test-workspace/projects/proj-1/work-items/issue-1/", r.URL.Path)
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

func mustMarshal(t *testing.T, v interface{}) []byte {
	data, err := json.Marshal(v)
	require.NoError(t, err)
	return data
}
