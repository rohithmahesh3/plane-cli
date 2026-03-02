package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rohithmahesh3/plane-cli/pkg/plane"
)

func TestListModules(t *testing.T) {
	mockModules := []plane.Module{
		{ID: "module-1", Name: "Authentication", Status: "in-progress"},
		{ID: "module-2", Name: "API Integration", Status: "backlog"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/workspaces/test-workspace/projects/test-project/modules/" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
		}

		response := struct {
			Results []plane.Module `json:"results"`
		}{
			Results: mockModules,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
		APIKey:     "test-key",
		Workspace:  "test-workspace",
	}

	modules, err := client.ListModules("test-project", false)
	if err != nil {
		t.Fatalf("ListModules failed: %v", err)
	}

	if len(modules) != 2 {
		t.Errorf("Expected 2 modules, got %d", len(modules))
	}
	if modules[0].Name != "Authentication" {
		t.Errorf("Expected first module name 'Authentication', got '%s'", modules[0].Name)
	}
}

func TestCreateModule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		var req plane.CreateModuleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}

		if req.Name != "User Dashboard" {
			t.Errorf("Expected name 'User Dashboard', got '%s'", req.Name)
		}

		module := plane.Module{
			ID:     "new-module-id",
			Name:   req.Name,
			Status: req.Status,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(module)
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
		APIKey:     "test-key",
		Workspace:  "test-workspace",
	}

	req := plane.CreateModuleRequest{
		Name:   "User Dashboard",
		Status: "backlog",
	}

	module, err := client.CreateModule("test-project", req)
	if err != nil {
		t.Fatalf("CreateModule failed: %v", err)
	}

	if module.Name != "User Dashboard" {
		t.Errorf("Expected module name 'User Dashboard', got '%s'", module.Name)
	}
}

func TestArchiveModule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/workspaces/test-workspace/projects/test-project/modules/module-123/archive/" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
		APIKey:     "test-key",
		Workspace:  "test-workspace",
	}

	err := client.ArchiveModule("test-project", "module-123")
	if err != nil {
		t.Fatalf("ArchiveModule failed: %v", err)
	}
}
