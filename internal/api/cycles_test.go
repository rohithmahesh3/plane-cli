package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rohithmahesh3/plane-cli/pkg/plane"
)

func TestListCycles(t *testing.T) {
	mockCycles := []plane.Cycle{
		{ID: "cycle-1", Name: "Sprint 1", StartDate: "2024-01-01", EndDate: "2024-01-14"},
		{ID: "cycle-2", Name: "Sprint 2", StartDate: "2024-01-15", EndDate: "2024-01-28"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/workspaces/test-workspace/projects/test-project/cycles/" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
		}

		response := struct {
			Results []plane.Cycle `json:"results"`
		}{
			Results: mockCycles,
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

	cycles, err := client.ListCycles("test-project", false)
	if err != nil {
		t.Fatalf("ListCycles failed: %v", err)
	}

	if len(cycles) != 2 {
		t.Errorf("Expected 2 cycles, got %d", len(cycles))
	}
	if cycles[0].Name != "Sprint 1" {
		t.Errorf("Expected first cycle name 'Sprint 1', got '%s'", cycles[0].Name)
	}
}

func TestCreateCycle(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		var req plane.CreateCycleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}

		if req.Name != "Sprint 3" {
			t.Errorf("Expected name 'Sprint 3', got '%s'", req.Name)
		}

		cycle := plane.Cycle{
			ID:        "new-cycle-id",
			Name:      req.Name,
			StartDate: req.StartDate,
			EndDate:   req.EndDate,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(cycle)
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
		APIKey:     "test-key",
		Workspace:  "test-workspace",
	}

	req := plane.CreateCycleRequest{
		Name:      "Sprint 3",
		StartDate: "2024-02-01",
		EndDate:   "2024-02-14",
	}

	cycle, err := client.CreateCycle("test-project", req)
	if err != nil {
		t.Fatalf("CreateCycle failed: %v", err)
	}

	if cycle.Name != "Sprint 3" {
		t.Errorf("Expected cycle name 'Sprint 3', got '%s'", cycle.Name)
	}
}

func TestArchiveCycle(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/workspaces/test-workspace/projects/test-project/cycles/cycle-123/archive/" {
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

	err := client.ArchiveCycle("test-project", "cycle-123")
	if err != nil {
		t.Fatalf("ArchiveCycle failed: %v", err)
	}
}
