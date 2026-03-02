//go:build integration
// +build integration

package api

import (
	"net/http"
	"os"
	"testing"

	"github.com/rohithmahesh3/plane-cli/internal/config"
	"github.com/rohithmahesh3/plane-cli/pkg/plane"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test configuration - set these environment variables to run integration tests
// PLANE_API_KEY - Your Plane API key
// PLANE_WORKSPACE - Test workspace slug (default: test-workspace)
// PLANE_API_HOST - API host URL (default: https://api.plane.so)

func setupIntegrationClient(t *testing.T) *Client {
	apiKey := os.Getenv("PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("PLANE_API_KEY not set, skipping integration test")
	}

	workspace := os.Getenv("PLANE_WORKSPACE")
	if workspace == "" {
		workspace = "test-workspace"
	}

	apiHost := os.Getenv("PLANE_API_HOST")
	if apiHost == "" {
		apiHost = "https://api.plane.so"
	}

	// Set up config
	config.Cfg.APIHost = apiHost
	config.Cfg.DefaultWorkspace = workspace
	err := config.SetAPIKey(apiKey)
	if err != nil {
		t.Logf("Warning: could not set API key in keyring: %v", err)
	}

	return &Client{
		HTTPClient: &http.Client{Timeout: DefaultTimeout},
		BaseURL:    apiHost,
		APIKey:     apiKey,
		Workspace:  workspace,
	}
}

func TestIntegration_ListWorkspaces(t *testing.T) {
	client := setupIntegrationClient(t)

	workspaces, err := client.ListWorkspaces()
	require.NoError(t, err)
	assert.NotEmpty(t, workspaces)

	t.Logf("Found %d workspaces:", len(workspaces))
	for _, ws := range workspaces {
		t.Logf("  - %s (%s): %s", ws.Name, ws.Slug, ws.ID)
	}
}

func TestIntegration_GetWorkspace(t *testing.T) {
	client := setupIntegrationClient(t)

	workspace, err := client.GetWorkspace(client.Workspace)
	require.NoError(t, err)
	assert.Equal(t, client.Workspace, workspace.Slug)
	assert.NotEmpty(t, workspace.ID)
	assert.NotEmpty(t, workspace.Name)

	t.Logf("Workspace: %s (%s)", workspace.Name, workspace.Slug)
}

func TestIntegration_ListProjects(t *testing.T) {
	client := setupIntegrationClient(t)

	projects, err := client.ListProjects()
	require.NoError(t, err)

	t.Logf("Found %d projects:", len(projects))
	for _, proj := range projects {
		t.Logf("  - %s (%s): %s", proj.Name, proj.Identifier, proj.ID)
	}
}

func TestIntegration_ProjectLifecycle(t *testing.T) {
	client := setupIntegrationClient(t)

	// Create a test project
	createReq := plane.CreateProjectRequest{
		Name:        "CLI Test Project",
		Identifier:  "CLITEST",
		Description: "Test project created by CLI integration tests",
	}

	project, err := client.CreateProject(createReq)
	require.NoError(t, err)
	assert.Equal(t, createReq.Name, project.Name)
	assert.Equal(t, createReq.Identifier, project.Identifier)

	t.Logf("Created project: %s (%s)", project.Name, project.ID)

	// Get project details
	retrievedProject, err := client.GetProject(project.ID)
	require.NoError(t, err)
	assert.Equal(t, project.ID, retrievedProject.ID)

	// List project members
	members, err := client.GetProjectMembers(project.ID)
	require.NoError(t, err)
	t.Logf("Project has %d members", len(members))

	// Clean up - delete project
	err = client.DeleteProject(project.ID)
	require.NoError(t, err)
	t.Log("Deleted test project")
}

func TestIntegration_IssueLifecycle(t *testing.T) {
	client := setupIntegrationClient(t)

	// First, get a project to work with
	projects, err := client.ListProjects()
	require.NoError(t, err)
	if len(projects) == 0 {
		t.Skip("No projects available for testing")
	}

	projectID := projects[0].ID
	t.Logf("Using project: %s (%s)", projects[0].Name, projectID)

	// Create an issue
	createReq := plane.CreateIssueRequest{
		Name:        "Test Issue from CLI",
		Description: "This is a test issue created by integration tests",
		Priority:    "medium",
	}

	issue, err := client.CreateIssue(projectID, createReq)
	require.NoError(t, err)
	assert.Equal(t, createReq.Name, issue.Name)
	assert.Equal(t, createReq.Priority, issue.Priority)

	t.Logf("Created issue: %s-%d (%s)", projects[0].Identifier, issue.SequenceID, issue.ID)

	// List issues
	opts := IssueListOptions{
		PerPage: 10,
	}
	issues, pagination, err := client.ListIssues(projectID, opts)
	require.NoError(t, err)
	assert.NotNil(t, pagination)
	t.Logf("Found %d issues in project (showing page with %d)", pagination.TotalResults, len(issues))

	// Get issue by ID
	retrievedIssue, err := client.GetIssue(projectID, issue.ID)
	require.NoError(t, err)
	assert.Equal(t, issue.ID, retrievedIssue.ID)

	// Get issue by sequence ID
	issueBySeq, err := client.GetIssueByIdentifier(projectID, issue.SequenceID)
	require.NoError(t, err)
	assert.Equal(t, issue.ID, issueBySeq.ID)

	// Update issue
	updateReq := plane.UpdateIssueRequest{
		Priority: "high",
	}
	updatedIssue, err := client.UpdateIssue(projectID, issue.ID, updateReq)
	require.NoError(t, err)
	assert.Equal(t, "high", updatedIssue.Priority)
	t.Log("Updated issue priority to high")

	// Search for issues
	searchResults, err := client.SearchIssues("Test Issue")
	require.NoError(t, err)
	t.Logf("Found %d issues matching search", len(searchResults))

	// Delete issue
	err = client.DeleteIssue(projectID, issue.ID)
	require.NoError(t, err)
	t.Log("Deleted test issue")
}

func TestIntegration_ListIssuesWithFilters(t *testing.T) {
	client := setupIntegrationClient(t)

	// Get a project
	projects, err := client.ListProjects()
	require.NoError(t, err)
	if len(projects) == 0 {
		t.Skip("No projects available for testing")
	}

	projectID := projects[0].ID

	// Test filtering by state
	opts := IssueListOptions{
		State:   "backlog",
		PerPage: 5,
	}
	issues, _, err := client.ListIssues(projectID, opts)
	require.NoError(t, err)
	t.Logf("Found %d issues in backlog state", len(issues))

	// Test filtering by priority
	opts = IssueListOptions{
		Priority: "high",
		PerPage:  5,
	}
	issues, _, err = client.ListIssues(projectID, opts)
	require.NoError(t, err)
	t.Logf("Found %d issues with high priority", len(issues))
}
