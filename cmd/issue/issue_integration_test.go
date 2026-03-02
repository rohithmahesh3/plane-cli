//go:build integration
// +build integration

package issue

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/rohithmahesh3/plane-cli/internal/api"
	"github.com/rohithmahesh3/plane-cli/internal/config"
	"github.com/rohithmahesh3/plane-cli/internal/integrationtest"
	"github.com/rohithmahesh3/plane-cli/pkg/plane"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestEnvironment(t *testing.T) *api.Client {
	integrationtest.WaitForSlot(t)

	apiKey := os.Getenv("PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("PLANE_API_KEY not set, skipping integration test")
	}

	workspace := os.Getenv("PLANE_WORKSPACE")
	if workspace == "" {
		t.Skip("PLANE_WORKSPACE not set, skipping integration test")
	}

	apiHost := os.Getenv("PLANE_API_HOST")
	if apiHost == "" {
		apiHost = "https://api.plane.so"
	}

	config.Cfg.APIHost = apiHost
	config.Cfg.DefaultWorkspace = workspace
	config.Cfg.OutputFormat = "json"

	originalService := config.KeyringService
	config.KeyringService = "plane-cli-cmd-test"
	t.Cleanup(func() {
		config.KeyringService = originalService
		config.DeleteAPIKey()
		api.ClearResolverCache()
	})

	err := config.SetAPIKey(apiKey)
	if err != nil {
		t.Skipf("Keyring not available in test environment: %v", err)
	}

	client, err := api.NewClient()
	require.NoError(t, err)

	projectID := os.Getenv("PLANE_PROJECT")
	if projectID == "" {
		projects, err := client.ListProjects()
		require.NoError(t, err)
		require.NotEmpty(t, projects, "expected at least one project for integration tests")
		projectID = projects[0].ID
	}
	config.Cfg.DefaultProject = projectID

	return client
}

func createTestIssue(t *testing.T, client *api.Client) *plane.Issue {
	projectID := config.Cfg.DefaultProject
	req := plane.CreateIssueRequest{
		Name:        fmt.Sprintf("Test Issue %d", time.Now().UnixNano()),
		Description: "Integration test issue",
		Priority:    "low",
	}

	issue, err := client.CreateIssue(projectID, req)
	require.NoError(t, err)
	return issue
}

func cleanupTestIssue(t *testing.T, client *api.Client, issueID string) {
	projectID := config.Cfg.DefaultProject
	_ = client.DeleteIssue(projectID, issueID)
}

func TestResolveAssigneeByUsername(t *testing.T) {
	client := setupTestEnvironment(t)
	projectID := config.Cfg.DefaultProject

	assigneeUsername := os.Getenv("PLANE_TEST_ASSIGNEE")
	if assigneeUsername == "" {
		t.Skip("PLANE_TEST_ASSIGNEE not set, skipping assignee test")
	}

	resolved, err := client.ResolveAssignees(projectID, []string{assigneeUsername})
	require.NoError(t, err)
	assert.Len(t, resolved, 1)
	assert.Len(t, resolved[0], 36, "Expected UUID format")
}

func TestResolveAssigneeByUUID(t *testing.T) {
	client := setupTestEnvironment(t)
	projectID := config.Cfg.DefaultProject

	assigneeUsername := os.Getenv("PLANE_TEST_ASSIGNEE")
	if assigneeUsername == "" {
		t.Skip("PLANE_TEST_ASSIGNEE not set, skipping assignee test")
	}

	resolvedByName, err := client.ResolveAssignees(projectID, []string{assigneeUsername})
	require.NoError(t, err)
	require.Len(t, resolvedByName, 1)

	uuid := resolvedByName[0]

	resolved, err := client.ResolveAssignees(projectID, []string{uuid})
	require.NoError(t, err)
	assert.Len(t, resolved, 1)
	assert.Equal(t, uuid, resolved[0])
}

func TestResolveAssigneeNotFound(t *testing.T) {
	client := setupTestEnvironment(t)
	projectID := config.Cfg.DefaultProject

	_, err := client.ResolveAssignees(projectID, []string{"nonexistent-user-12345"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestResolveAssigneeNotFoundIncludesSuggestions(t *testing.T) {
	client := setupTestEnvironment(t)
	projectID := config.Cfg.DefaultProject

	assigneeUsername := os.Getenv("PLANE_TEST_ASSIGNEE")
	if len(assigneeUsername) < 2 {
		t.Skip("PLANE_TEST_ASSIGNEE must be set to at least 2 characters")
	}

	query := assigneeUsername[:len(assigneeUsername)-1]

	_, err := client.ResolveAssignees(projectID, []string{query})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Closest matches:")
	assert.Contains(t, err.Error(), "plane workspace members --search")
	assert.Contains(t, strings.ToLower(err.Error()), strings.ToLower(assigneeUsername))
}

func TestResolveStateByName(t *testing.T) {
	client := setupTestEnvironment(t)
	projectID := config.Cfg.DefaultProject

	states, err := client.ListStates(projectID)
	require.NoError(t, err)
	require.NotEmpty(t, states)

	resolved, err := client.ResolveState(projectID, states[0].Name)
	require.NoError(t, err)
	assert.NotEmpty(t, resolved)
	assert.Len(t, resolved, 36, "Expected UUID format")
}

func TestResolveStateCaseInsensitive(t *testing.T) {
	client := setupTestEnvironment(t)
	projectID := config.Cfg.DefaultProject

	states, err := client.ListStates(projectID)
	require.NoError(t, err)
	require.NotEmpty(t, states)

	stateName := states[0].Name

	resolved, err := client.ResolveState(projectID, strings.ToUpper(stateName))
	require.NoError(t, err)
	assert.NotEmpty(t, resolved)

	resolved2, err := client.ResolveState(projectID, strings.ToLower(stateName))
	require.NoError(t, err)
	assert.Equal(t, resolved, resolved2)
}

func TestResolveStateByUUID(t *testing.T) {
	client := setupTestEnvironment(t)
	projectID := config.Cfg.DefaultProject

	states, err := client.ListStates(projectID)
	require.NoError(t, err)
	require.NotEmpty(t, states)

	uuid := states[0].ID

	resolved, err := client.ResolveState(projectID, uuid)
	require.NoError(t, err)
	assert.Equal(t, uuid, resolved)
}

func TestResolveStateNotFound(t *testing.T) {
	client := setupTestEnvironment(t)
	projectID := config.Cfg.DefaultProject

	_, err := client.ResolveState(projectID, "NonExistentState")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestIssueCreateWithAssigneeIntegration(t *testing.T) {
	client := setupTestEnvironment(t)
	projectID := config.Cfg.DefaultProject

	assigneeUsername := os.Getenv("PLANE_TEST_ASSIGNEE")
	if assigneeUsername == "" {
		t.Skip("PLANE_TEST_ASSIGNEE not set, skipping assignee test")
	}

	resolvedAssignees, err := client.ResolveAssignees(projectID, []string{assigneeUsername})
	require.NoError(t, err, "Failed to resolve assignee username to UUID")

	req := plane.CreateIssueRequest{
		Name:        fmt.Sprintf("Test Issue with Assignee %d", os.Getpid()),
		Description: "Integration test for assignee on create",
		Priority:    "medium",
		Assignees:   resolvedAssignees,
	}

	issue, err := client.CreateIssue(projectID, req)
	require.NoError(t, err)
	require.NotNil(t, issue)
	defer cleanupTestIssue(t, client, issue.ID)

	assert.NotEmpty(t, issue.ID)
	assert.Equal(t, req.Name, issue.Name)

	fetchedIssue, err := client.GetIssue(projectID, issue.ID)
	require.NoError(t, err)

	assert.Len(t, fetchedIssue.Assignees, 1, "Expected one assignee")
	if len(fetchedIssue.Assignees) > 0 {
		assert.Equal(t, assigneeUsername, fetchedIssue.Assignees[0].DisplayName)
	}
}

func TestIssueEditWithAssigneeIntegration(t *testing.T) {
	client := setupTestEnvironment(t)
	projectID := config.Cfg.DefaultProject

	assigneeUsername := os.Getenv("PLANE_TEST_ASSIGNEE")
	if assigneeUsername == "" {
		t.Skip("PLANE_TEST_ASSIGNEE not set, skipping assignee test")
	}

	issue := createTestIssue(t, client)
	defer cleanupTestIssue(t, client, issue.ID)

	resolvedAssignees, err := client.ResolveAssignees(projectID, []string{assigneeUsername})
	require.NoError(t, err, "Failed to resolve assignee username to UUID")

	updateReq := plane.UpdateIssueRequest{
		Assignees: resolvedAssignees,
	}

	updatedIssue, err := client.UpdateIssue(projectID, issue.ID, updateReq)
	require.NoError(t, err)
	require.NotNil(t, updatedIssue)

	fetchedIssue, err := client.GetIssue(projectID, issue.ID)
	require.NoError(t, err)

	assert.Len(t, fetchedIssue.Assignees, 1, "Expected one assignee after update")
	if len(fetchedIssue.Assignees) > 0 {
		assert.Equal(t, assigneeUsername, fetchedIssue.Assignees[0].DisplayName)
	}
}

func TestIssueEditClearAssigneeIntegration(t *testing.T) {
	client := setupTestEnvironment(t)
	projectID := config.Cfg.DefaultProject

	assigneeUsername := os.Getenv("PLANE_TEST_ASSIGNEE")
	if assigneeUsername == "" {
		t.Skip("PLANE_TEST_ASSIGNEE not set, skipping assignee test")
	}

	issue := createTestIssue(t, client)
	defer cleanupTestIssue(t, client, issue.ID)

	resolvedAssignees, err := client.ResolveAssignees(projectID, []string{assigneeUsername})
	require.NoError(t, err, "Failed to resolve assignee username to UUID")

	updateReq := plane.UpdateIssueRequest{
		Assignees: resolvedAssignees,
	}
	_, err = client.UpdateIssue(projectID, issue.ID, updateReq)
	require.NoError(t, err)

	fetchedIssue, err := client.GetIssue(projectID, issue.ID)
	require.NoError(t, err)
	assert.Len(t, fetchedIssue.Assignees, 1, "Expected one assignee after adding")
}

func TestIssueCreateWithLabelsIntegration(t *testing.T) {
	client := setupTestEnvironment(t)
	projectID := config.Cfg.DefaultProject

	label1, err := client.CreateLabel(projectID, plane.CreateLabelRequest{
		Name:  fmt.Sprintf("test-label-1-%d", os.Getpid()),
		Color: "#EF4444",
	})
	require.NoError(t, err, "Failed to create test label 1")
	defer client.DeleteLabel(projectID, label1.ID)

	label2, err := client.CreateLabel(projectID, plane.CreateLabelRequest{
		Name:  fmt.Sprintf("test-label-2-%d", os.Getpid()),
		Color: "#3B82F6",
	})
	require.NoError(t, err, "Failed to create test label 2")
	defer client.DeleteLabel(projectID, label2.ID)

	api.ClearResolverCache()

	resolvedLabels, err := client.ResolveLabels(projectID, []string{label1.Name, label2.Name})
	require.NoError(t, err, "Failed to resolve label names to UUIDs")

	req := plane.CreateIssueRequest{
		Name:        fmt.Sprintf("Test Issue with Labels %d", os.Getpid()),
		Description: "Integration test for labels on create",
		Priority:    "medium",
		Labels:      resolvedLabels,
	}

	issue, err := client.CreateIssue(projectID, req)
	require.NoError(t, err)
	require.NotNil(t, issue)
	defer cleanupTestIssue(t, client, issue.ID)

	assert.NotEmpty(t, issue.ID)
	assert.Equal(t, req.Name, issue.Name)

	fetchedIssue, err := client.GetIssue(projectID, issue.ID)
	require.NoError(t, err)

	assert.Len(t, fetchedIssue.Labels, 2, "Expected two labels to be set")
}

func TestIssueEditWithStateIntegration(t *testing.T) {
	client := setupTestEnvironment(t)
	projectID := config.Cfg.DefaultProject

	issue := createTestIssue(t, client)
	defer cleanupTestIssue(t, client, issue.ID)

	states, err := client.ListStates(projectID)
	require.NoError(t, err)
	require.NotEmpty(t, states)

	resolvedState, err := client.ResolveState(projectID, states[0].Name)
	require.NoError(t, err, "Failed to resolve state name to UUID")

	updateReq := plane.UpdateIssueRequest{
		State: resolvedState,
	}

	updatedIssue, err := client.UpdateIssue(projectID, issue.ID, updateReq)
	require.NoError(t, err)
	require.NotNil(t, updatedIssue)

	fetchedIssue, err := client.GetIssue(projectID, issue.ID)
	require.NoError(t, err)

	assert.NotEmpty(t, fetchedIssue.State, "Expected state to be set")
}

func TestResolverCaching(t *testing.T) {
	client := setupTestEnvironment(t)
	projectID := config.Cfg.DefaultProject

	states, err := client.ListStates(projectID)
	require.NoError(t, err)
	require.NotEmpty(t, states)

	_, err = client.ResolveState(projectID, states[0].Name)
	require.NoError(t, err)

	secondState := states[0].Name
	if len(states) > 1 {
		secondState = states[1].Name
	}

	_, err = client.ResolveState(projectID, secondState)
	require.NoError(t, err)

	api.ClearResolverCache()
}
