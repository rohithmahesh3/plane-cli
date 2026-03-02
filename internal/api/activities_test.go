package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rohithmahesh3/plane-cli/pkg/plane"
)

func TestListActivities(t *testing.T) {
	mockActivities := []plane.Activity{
		{
			ID:      "activity-1",
			Verb:    "created",
			Comment: "created the work item",
			Actor:   "user-1",
		},
		{
			ID:       "activity-2",
			Verb:     "updated",
			Field:    "state",
			OldValue: "Todo",
			NewValue: "In Progress",
			Actor:    "user-1",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/workspaces/test-workspace/projects/test-project/work-items/issue-123/activities/" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
		}

		response := struct {
			Results []plane.Activity `json:"results"`
		}{
			Results: mockActivities,
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

	activities, err := client.ListActivities("test-project", "issue-123")
	if err != nil {
		t.Fatalf("ListActivities failed: %v", err)
	}

	if len(activities) != 2 {
		t.Errorf("Expected 2 activities, got %d", len(activities))
	}
	if activities[0].Verb != "created" {
		t.Errorf("Expected first activity verb 'created', got '%s'", activities[0].Verb)
	}
}

func TestGetActivity(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/workspaces/test-workspace/projects/test-project/work-items/issue-123/activities/activity-456/" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
		}

		activity := plane.Activity{
			ID:    "activity-456",
			Verb:  "updated",
			Field: "priority",
			Actor: "user-1",
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(activity)
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
		APIKey:     "test-key",
		Workspace:  "test-workspace",
	}

	activity, err := client.GetActivity("test-project", "issue-123", "activity-456")
	if err != nil {
		t.Fatalf("GetActivity failed: %v", err)
	}

	if activity.ID != "activity-456" {
		t.Errorf("Expected activity ID 'activity-456', got '%s'", activity.ID)
	}
}
