package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rohithmahesh3/plane-cli/pkg/plane"
)

func TestListStates(t *testing.T) {
	mockStates := []plane.State{
		{ID: "state-1", Name: "Todo", Color: "#3B82F6", Group: "unstarted"},
		{ID: "state-2", Name: "In Progress", Color: "#F59E0B", Group: "started"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/workspaces/test-workspace/projects/test-project/states/" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
		}

		response := struct {
			Results []plane.State `json:"results"`
		}{
			Results: mockStates,
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

	states, err := client.ListStates("test-project")
	if err != nil {
		t.Fatalf("ListStates failed: %v", err)
	}

	if len(states) != 2 {
		t.Errorf("Expected 2 states, got %d", len(states))
	}
	if states[0].Name != "Todo" {
		t.Errorf("Expected first state name 'Todo', got '%s'", states[0].Name)
	}
}

func TestCreateState(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		var req plane.CreateStateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}

		if req.Name != "Review" {
			t.Errorf("Expected name 'Review', got '%s'", req.Name)
		}

		state := plane.State{
			ID:    "new-state-id",
			Name:  req.Name,
			Color: req.Color,
			Group: req.Group,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(state)
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
		APIKey:     "test-key",
		Workspace:  "test-workspace",
	}

	req := plane.CreateStateRequest{
		Name:  "Review",
		Color: "#8B5CF6",
		Group: "started",
	}

	state, err := client.CreateState("test-project", req)
	if err != nil {
		t.Fatalf("CreateState failed: %v", err)
	}

	if state.Name != "Review" {
		t.Errorf("Expected state name 'Review', got '%s'", state.Name)
	}
}

func TestDeleteState(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/workspaces/test-workspace/projects/test-project/states/state-123/" {
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

	err := client.DeleteState("test-project", "state-123")
	if err != nil {
		t.Fatalf("DeleteState failed: %v", err)
	}
}
