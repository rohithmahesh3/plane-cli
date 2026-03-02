//go:build integration
// +build integration

package issue

import (
	"fmt"
	"os"
	"testing"

	"github.com/rohithmahesh3/plane-cli/internal/api"
	"github.com/rohithmahesh3/plane-cli/internal/config"
	"github.com/rohithmahesh3/plane-cli/pkg/plane"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestEnvironment(t *testing.T) *api.Client {
	apiKey := os.Getenv("PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("PLANE_API_KEY not set, skipping integration test")
	}

	workspace := os.Getenv("PLANE_WORKSPACE")
	if workspace == "" {
		t.Skip("PLANE_WORKSPACE not set, skipping integration test")
	}

	projectID := os.Getenv("PLANE_PROJECT")
	if projectID == "" {
		t.Skip("PLANE_PROJECT not set, skipping integration test")
	}

	apiHost := os.Getenv("PLANE_API_HOST")
	if apiHost == "" {
		apiHost = "https://api.plane.so"
	}

	config.Cfg.APIHost = apiHost
	config.Cfg.DefaultWorkspace = workspace
	config.Cfg.DefaultProject = projectID
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
	return client
}

func createTestIssue(t *testing.T, client *api.Client) *plane.Issue {
	projectID := config.Cfg.DefaultProject
	req := plane.CreateIssueRequest{
		Name:        fmt.Sprintf("Test Issue %d", os.Getpid()),
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

	uuid := "a9bae425-06ff-4cf2-b1c5-2b6cc6b5d810"

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

func TestResolveStateByName(t *testing.T) {
	client := setupTestEnvironment(t)
	projectID := config.Cfg.DefaultProject

	resolved, err := client.ResolveState(projectID, "Todo")
	require.NoError(t, err)
	assert.NotEmpty(t, resolved)
	assert.Len(t, resolved, 36, "Expected UUID format")
}

func TestResolveStateCaseInsensitive(t *testing.T) {
	client := setupTestEnvironment(t)
	projectID := config.Cfg.DefaultProject

	resolved, err := client.ResolveState(projectID, "TODO")
	require.NoError(t, err)
	assert.NotEmpty(t, resolved)

	resolved2, err := client.ResolveState(projectID, "todo")
	require.NoError(t, err)
	assert.Equal(t, resolved, resolved2)
}

func TestResolveStateByUUID(t *testing.T) {
	client := setupTestEnvironment(t)
	projectID := config.Cfg.DefaultProject

	uuid := "1578ec52-0083-440d-b370-b3c9286bc091"

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

	resolvedState, err := client.ResolveState(projectID, "Done")
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

	_, err := client.ResolveState(projectID, "Todo")
	require.NoError(t, err)

	_, err = client.ResolveState(projectID, "Backlog")
	require.NoError(t, err)

	api.ClearResolverCache()
}
