package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rohithmahesh3/plane-cli/pkg/plane"
)

func TestListWorklogs(t *testing.T) {
	mockWorklogs := []plane.Worklog{
		{ID: "worklog-1", Description: "Initial work", Duration: 60},
		{ID: "worklog-2", Description: "More work", Duration: 90},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/workspaces/test-workspace/projects/test-project/work-items/issue-123/worklogs/" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
		}

		response := struct {
			Results []plane.Worklog `json:"results"`
		}{
			Results: mockWorklogs,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
		APIKey:     "test-key",
		Workspace:  "test-workspace",
	}

	worklogs, err := client.ListWorklogs("test-project", "issue-123")
	if err != nil {
		t.Fatalf("ListWorklogs failed: %v", err)
	}

	if len(worklogs) != 2 {
		t.Errorf("Expected 2 worklogs, got %d", len(worklogs))
	}
	if worklogs[0].Duration != 60 {
		t.Errorf("Expected first worklog duration 60, got %d", worklogs[0].Duration)
	}
}

func TestCreateWorklog(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		var req plane.CreateWorklogRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}

		if req.Duration != 120 {
			t.Errorf("Expected duration 120, got %d", req.Duration)
		}
		if req.Description != "Worked on feature" {
			t.Errorf("Expected description 'Worked on feature', got '%s'", req.Description)
		}

		worklog := plane.Worklog{
			ID:          "new-worklog-id",
			Description: req.Description,
			Duration:    req.Duration,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(worklog)
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
		APIKey:     "test-key",
		Workspace:  "test-workspace",
	}

	req := plane.CreateWorklogRequest{
		Description: "Worked on feature",
		Duration:    120,
	}

	worklog, err := client.CreateWorklog("test-project", "issue-123", req)
	if err != nil {
		t.Fatalf("CreateWorklog failed: %v", err)
	}

	if worklog.Duration != 120 {
		t.Errorf("Expected worklog duration 120, got %d", worklog.Duration)
	}
}

func TestGetTotalWorklogTime(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/workspaces/test-workspace/projects/test-project/work-items/issue-123/worklogs/" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
		}

		response := struct {
			Results []plane.Worklog `json:"results"`
		}{
			Results: []plane.Worklog{
				{ID: "worklog-1", Duration: 60},
				{ID: "worklog-2", Duration: 90},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
		APIKey:     "test-key",
		Workspace:  "test-workspace",
	}

	total, err := client.GetTotalWorklogTime("test-project", "issue-123")
	if err != nil {
		t.Fatalf("GetTotalWorklogTime failed: %v", err)
	}

	if total != 150 {
		t.Errorf("Expected total 150, got %d", total)
	}
}

func TestDeleteWorklog(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/workspaces/test-workspace/projects/test-project/work-items/issue-123/worklogs/worklog-456/" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
		APIKey:     "test-key",
		Workspace:  "test-workspace",
	}

	err := client.DeleteWorklog("test-project", "issue-123", "worklog-456")
	if err != nil {
		t.Fatalf("DeleteWorklog failed: %v", err)
	}
}
