package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rohithmahesh3/plane-cli/pkg/plane"
)

func TestListLabels(t *testing.T) {
	mockLabels := []plane.Label{
		{ID: "label-1", Name: "Bug", Color: "#EF4444"},
		{ID: "label-2", Name: "Feature", Color: "#3B82F6"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/workspaces/test-workspace/projects/test-project/labels/" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
		}

		response := struct {
			Results []plane.Label `json:"results"`
		}{
			Results: mockLabels,
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

	labels, err := client.ListLabels("test-project")
	if err != nil {
		t.Fatalf("ListLabels failed: %v", err)
	}

	if len(labels) != 2 {
		t.Errorf("Expected 2 labels, got %d", len(labels))
	}
	if labels[0].Name != "Bug" {
		t.Errorf("Expected first label name 'Bug', got '%s'", labels[0].Name)
	}
}

func TestCreateLabel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		var req plane.CreateLabelRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}

		if req.Name != "High Priority" {
			t.Errorf("Expected name 'High Priority', got '%s'", req.Name)
		}

		label := plane.Label{
			ID:    "new-label-id",
			Name:  req.Name,
			Color: req.Color,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(label)
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
		APIKey:     "test-key",
		Workspace:  "test-workspace",
	}

	req := plane.CreateLabelRequest{
		Name:  "High Priority",
		Color: "#DC2626",
	}

	label, err := client.CreateLabel("test-project", req)
	if err != nil {
		t.Fatalf("CreateLabel failed: %v", err)
	}

	if label.Name != "High Priority" {
		t.Errorf("Expected label name 'High Priority', got '%s'", label.Name)
	}
}

func TestDeleteLabel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/workspaces/test-workspace/projects/test-project/labels/label-123/" {
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

	err := client.DeleteLabel("test-project", "label-123")
	if err != nil {
		t.Fatalf("DeleteLabel failed: %v", err)
	}
}
