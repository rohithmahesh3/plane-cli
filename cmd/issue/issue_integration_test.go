//go:build integration
// +build integration

package issue

import (
	"bytes"
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
	t.Cleanup(func() { config.KeyringService = originalService; config.DeleteAPIKey() })

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

func TestIssueSearchIntegration(t *testing.T) {
	setupTestEnvironment(t)

	cmd := searchCmd
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err := runSearch(cmd, []string{"Test"})
	assert.NoError(t, err)
}

func TestIssueCreateWithAssigneeIntegration(t *testing.T) {
	client := setupTestEnvironment(t)
	projectID := config.Cfg.DefaultProject

	assigneeUsername := os.Getenv("PLANE_TEST_ASSIGNEE")
	if assigneeUsername == "" {
		t.Skip("PLANE_TEST_ASSIGNEE not set, skipping assignee test")
	}

	req := plane.CreateIssueRequest{
		Name:        fmt.Sprintf("Test Issue with Assignee %d", os.Getpid()),
		Description: "Integration test for assignee on create",
		Priority:    "medium",
		Assignees:   []string{assigneeUsername},
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
		assert.Equal(t, assigneeUsername, fetchedIssue.Assignees[0].Username)
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

	updateReq := plane.UpdateIssueRequest{
		Assignees: []string{assigneeUsername},
	}

	updatedIssue, err := client.UpdateIssue(projectID, issue.ID, updateReq)
	require.NoError(t, err)
	require.NotNil(t, updatedIssue)

	fetchedIssue, err := client.GetIssue(projectID, issue.ID)
	require.NoError(t, err)

	assert.Len(t, fetchedIssue.Assignees, 1, "Expected one assignee after update")
	if len(fetchedIssue.Assignees) > 0 {
		assert.Equal(t, assigneeUsername, fetchedIssue.Assignees[0].Username)
	}
}

func TestIssueEditReplaceAssigneeIntegration(t *testing.T) {
	client := setupTestEnvironment(t)
	projectID := config.Cfg.DefaultProject

	assigneeUsername := os.Getenv("PLANE_TEST_ASSIGNEE")
	if assigneeUsername == "" {
		t.Skip("PLANE_TEST_ASSIGNEE not set, skipping assignee test")
	}

	issue := createTestIssue(t, client)
	defer cleanupTestIssue(t, client, issue.ID)

	updateReq := plane.UpdateIssueRequest{
		Assignees: []string{assigneeUsername},
	}
	_, err := client.UpdateIssue(projectID, issue.ID, updateReq)
	require.NoError(t, err)

	clearReq := plane.UpdateIssueRequest{
		Assignees: []string{},
	}
	_, err = client.UpdateIssue(projectID, issue.ID, clearReq)
	require.NoError(t, err)

	fetchedIssue, err := client.GetIssue(projectID, issue.ID)
	require.NoError(t, err)

	assert.Len(t, fetchedIssue.Assignees, 0, "Expected no assignees after clearing")
}

func TestIssueCreateWithLabelsIntegration(t *testing.T) {
	client := setupTestEnvironment(t)
	projectID := config.Cfg.DefaultProject

	req := plane.CreateIssueRequest{
		Name:        fmt.Sprintf("Test Issue with Labels %d", os.Getpid()),
		Description: "Integration test for labels on create",
		Priority:    "medium",
		Labels:      []string{"bug", "test"},
	}

	issue, err := client.CreateIssue(projectID, req)
	require.NoError(t, err)
	require.NotNil(t, issue)
	defer cleanupTestIssue(t, client, issue.ID)

	assert.NotEmpty(t, issue.ID)
	assert.Equal(t, req.Name, issue.Name)

	fetchedIssue, err := client.GetIssue(projectID, issue.ID)
	require.NoError(t, err)

	assert.NotEmpty(t, fetchedIssue.Labels, "Expected labels to be set")
}
