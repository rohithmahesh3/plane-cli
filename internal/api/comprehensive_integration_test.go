//go:build integration
// +build integration

package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rohithmahesh3/plane-cli/internal/config"
	"github.com/rohithmahesh3/plane-cli/internal/integrationtest"
	"github.com/rohithmahesh3/plane-cli/pkg/plane"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testClient *Client
var testProjectID string

func setupComprehensiveTest(t *testing.T) *Client {
	integrationtest.WaitForSlot(t)

	if testClient != nil {
		return testClient
	}

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

	config.Cfg.APIHost = apiHost
	config.Cfg.DefaultWorkspace = workspace

	originalService := config.KeyringService
	config.KeyringService = "plane-cli-integration-test"
	t.Cleanup(func() {
		config.KeyringService = originalService
		config.DeleteAPIKey()
	})

	err := config.SetAPIKey(apiKey)
	if err != nil {
		t.Skipf("Keyring not available in test environment: %v", err)
	}

	testClient = &Client{
		HTTPClient: &http.Client{Timeout: DefaultTimeout},
		BaseURL:    apiHost,
		APIKey:     apiKey,
		Workspace:  workspace,
	}

	projectID := os.Getenv("PLANE_PROJECT")
	if projectID == "" {
		projects, err := testClient.ListProjects()
		if err != nil {
			t.Skipf("failed to list projects for integration setup: %v", err)
		}
		if len(projects) == 0 {
			t.Skip("no projects available for integration tests")
		}
		projectID = projects[0].ID
	}

	config.Cfg.DefaultProject = projectID
	testProjectID = projectID

	return testClient
}

func createTestIssueForModule(t *testing.T, client *Client, name string) *plane.Issue {
	req := plane.CreateIssueRequest{
		Name:        fmt.Sprintf("%s %d", name, time.Now().UnixNano()),
		Description: "Test issue for comprehensive integration tests",
		Priority:    "low",
	}

	issue, err := client.CreateIssue(testProjectID, req)
	require.NoError(t, err)
	return issue
}

func cleanupTestIssue(t *testing.T, client *Client, issueID string) {
	_ = client.DeleteIssue(testProjectID, issueID)
}

func requireTestProject(t *testing.T, client *Client) *plane.Project {
	t.Helper()

	project, err := client.GetProject(testProjectID)
	require.NoError(t, err)
	return project
}

// ============================================
// LINKS TESTS
// ============================================

func TestLinksLifecycle(t *testing.T) {
	client := setupComprehensiveTest(t)

	issue := createTestIssueForModule(t, client, "Links Test Issue")
	defer cleanupTestIssue(t, client, issue.ID)

	link, err := client.CreateLink(testProjectID, issue.ID, plane.CreateLinkRequest{
		URL:   "https://example.com/test-link",
		Title: "Test Link",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, link.ID)
	assert.Equal(t, "https://example.com/test-link", link.URL)
	t.Logf("Created link: %s", link.ID)

	fetchedLink, err := client.GetLink(testProjectID, issue.ID, link.ID)
	require.NoError(t, err)
	assert.Equal(t, link.ID, fetchedLink.ID)

	links, err := client.ListLinks(testProjectID, issue.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(links), 1)

	updatedLink, err := client.UpdateLink(testProjectID, issue.ID, link.ID, plane.UpdateLinkRequest{
		Title: "Updated Test Link",
	})
	require.NoError(t, err)
	assert.Equal(t, "Updated Test Link", updatedLink.Title)

	err = client.DeleteLink(testProjectID, issue.ID, link.ID)
	require.NoError(t, err)
}

// ============================================
// COMMENTS TESTS
// ============================================

func TestCommentsLifecycle(t *testing.T) {
	client := setupComprehensiveTest(t)

	issue := createTestIssueForModule(t, client, "Comments Test Issue")
	defer cleanupTestIssue(t, client, issue.ID)

	comment, err := client.CreateComment(testProjectID, issue.ID, plane.CreateCommentRequest{
		CommentHTML: "<p>Test comment from integration test</p>",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, comment.ID)
	t.Logf("Created comment: %s", comment.ID)

	fetchedComment, err := client.GetComment(testProjectID, issue.ID, comment.ID)
	require.NoError(t, err)
	assert.Equal(t, comment.ID, fetchedComment.ID)

	comments, err := client.ListComments(testProjectID, issue.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(comments), 1)

	updatedComment, err := client.UpdateComment(testProjectID, issue.ID, comment.ID, plane.UpdateCommentRequest{
		CommentHTML: "<p>Updated test comment</p>",
	})
	require.NoError(t, err)
	assert.Contains(t, updatedComment.CommentHTML, "Updated")

	err = client.DeleteComment(testProjectID, issue.ID, comment.ID)
	require.NoError(t, err)
}

// ============================================
// ACTIVITIES TESTS
// ============================================

func TestActivitiesList(t *testing.T) {
	client := setupComprehensiveTest(t)

	issue := createTestIssueForModule(t, client, "Activities Test Issue")
	defer cleanupTestIssue(t, client, issue.ID)

	_, err := client.UpdateIssue(testProjectID, issue.ID, plane.UpdateIssueRequest{
		Priority: "high",
	})
	require.NoError(t, err)

	activities, err := client.ListActivities(testProjectID, issue.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(activities), 1)
	t.Logf("Found %d activities", len(activities))

	for _, a := range activities {
		t.Logf("  - Activity: %s (verb: %s, field: %s)", a.ID, a.Verb, a.Field)
	}
}

// ============================================
// WORKLOGS TESTS
// ============================================

func TestWorklogsLifecycle(t *testing.T) {
	client := setupComprehensiveTest(t)
	project := requireTestProject(t, client)
	if !project.IsTimeTrackingEnabled {
		t.Skip("time tracking is disabled for this project")
	}

	issue := createTestIssueForModule(t, client, "Worklogs Test Issue")
	defer cleanupTestIssue(t, client, issue.ID)

	worklog, err := client.CreateWorklog(testProjectID, issue.ID, plane.CreateWorklogRequest{
		Description: "Test worklog entry",
		Duration:    3600,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, worklog.ID)
	assert.Equal(t, 3600, worklog.Duration)
	t.Logf("Created worklog: %s (duration: %d seconds)", worklog.ID, worklog.Duration)

	fetchedWorklog, err := client.GetWorklog(testProjectID, issue.ID, worklog.ID)
	require.NoError(t, err)
	assert.Equal(t, worklog.ID, fetchedWorklog.ID)

	worklogs, err := client.ListWorklogs(testProjectID, issue.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(worklogs), 1)

	updatedWorklog, err := client.UpdateWorklog(testProjectID, issue.ID, worklog.ID, plane.UpdateWorklogRequest{
		Duration:    7200,
		Description: "Updated worklog entry",
	})
	require.NoError(t, err)
	assert.Equal(t, 7200, updatedWorklog.Duration)

	totalTime, err := client.GetTotalWorklogTime(testProjectID, issue.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, totalTime, 7200)
	t.Logf("Total time logged: %d seconds", totalTime)

	err = client.DeleteWorklog(testProjectID, issue.ID, worklog.ID)
	require.NoError(t, err)
}

// ============================================
// CYCLES TESTS
// ============================================

func TestCyclesLifecycle(t *testing.T) {
	client := setupComprehensiveTest(t)

	cycleName := fmt.Sprintf("Test Cycle %d", time.Now().UnixNano())
	cycle, err := client.CreateCycle(testProjectID, plane.CreateCycleRequest{
		Name:        cycleName,
		Description: "Test cycle from integration tests",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, cycle.ID)
	t.Logf("Created cycle: %s (%s)", cycle.Name, cycle.ID)

	defer func() {
		err := client.DeleteCycle(testProjectID, cycle.ID)
		if err != nil {
			t.Logf("Warning: failed to delete cycle: %v", err)
		}
	}()

	fetchedCycle, err := client.GetCycle(testProjectID, cycle.ID)
	require.NoError(t, err)
	assert.Equal(t, cycle.ID, fetchedCycle.ID)

	cycles, err := client.ListCycles(testProjectID, false)
	require.NoError(t, err)
	found := false
	for _, c := range cycles {
		if c.ID == cycle.ID {
			found = true
			break
		}
	}
	assert.True(t, found, "Created cycle should be in list")

	issue := createTestIssueForModule(t, client, "Cycle Test Issue")
	defer cleanupTestIssue(t, client, issue.ID)

	err = client.AddIssuesToCycle(testProjectID, cycle.ID, []string{issue.ID})
	require.NoError(t, err)
	t.Log("Added issue to cycle")

	cycleIssues, err := client.ListCycleIssues(testProjectID, cycle.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(cycleIssues), 1)
	t.Logf("Cycle has %d issues", len(cycleIssues))

	updatedCycle, err := client.UpdateCycle(testProjectID, cycle.ID, plane.UpdateCycleRequest{
		Name:        cycleName + " Updated",
		Description: "Updated description",
	})
	require.NoError(t, err)
	assert.Contains(t, updatedCycle.Name, "Updated")

	err = client.RemoveIssueFromCycle(testProjectID, cycle.ID, issue.ID)
	require.NoError(t, err)
	t.Log("Removed issue from cycle")
}

// ============================================
// MODULES TESTS
// ============================================

func TestModulesLifecycle(t *testing.T) {
	client := setupComprehensiveTest(t)

	moduleName := fmt.Sprintf("Test Module %d", time.Now().UnixNano())
	module, err := client.CreateModule(testProjectID, plane.CreateModuleRequest{
		Name:        moduleName,
		Description: "Test module from integration tests",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, module.ID)
	t.Logf("Created module: %s (%s)", module.Name, module.ID)

	defer func() {
		err := client.DeleteModule(testProjectID, module.ID)
		if err != nil {
			t.Logf("Warning: failed to delete module: %v", err)
		}
	}()

	fetchedModule, err := client.GetModule(testProjectID, module.ID)
	require.NoError(t, err)
	assert.Equal(t, module.ID, fetchedModule.ID)

	modules, err := client.ListModules(testProjectID, false)
	require.NoError(t, err)
	found := false
	for _, m := range modules {
		if m.ID == module.ID {
			found = true
			break
		}
	}
	assert.True(t, found, "Created module should be in list")

	issue := createTestIssueForModule(t, client, "Module Test Issue")
	defer cleanupTestIssue(t, client, issue.ID)

	err = client.AddIssuesToModule(testProjectID, module.ID, []string{issue.ID})
	require.NoError(t, err)
	t.Log("Added issue to module")

	moduleIssues, err := client.ListModuleIssues(testProjectID, module.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(moduleIssues), 1)
	t.Logf("Module has %d issues", len(moduleIssues))

	updatedModule, err := client.UpdateModule(testProjectID, module.ID, plane.UpdateModuleRequest{
		Name:        moduleName + " Updated",
		Description: "Updated description",
	})
	require.NoError(t, err)
	assert.Contains(t, updatedModule.Name, "Updated")

	err = client.RemoveIssueFromModule(testProjectID, module.ID, issue.ID)
	require.NoError(t, err)
	t.Log("Removed issue from module")
}

// ============================================
// LABELS TESTS
// ============================================

func TestLabelsLifecycle(t *testing.T) {
	client := setupComprehensiveTest(t)

	labelName := fmt.Sprintf("test-label-%d", time.Now().UnixNano())
	label, err := client.CreateLabel(testProjectID, plane.CreateLabelRequest{
		Name:        labelName,
		Description: "Test label from integration tests",
		Color:       "#FF5733",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, label.ID)
	t.Logf("Created label: %s (%s)", label.Name, label.ID)

	defer func() {
		err := client.DeleteLabel(testProjectID, label.ID)
		if err != nil {
			t.Logf("Warning: failed to delete label: %v", err)
		}
	}()

	fetchedLabel, err := client.GetLabel(testProjectID, label.ID)
	require.NoError(t, err)
	assert.Equal(t, label.ID, fetchedLabel.ID)

	labels, err := client.ListLabels(testProjectID)
	require.NoError(t, err)
	found := false
	for _, l := range labels {
		if l.ID == label.ID {
			found = true
			break
		}
	}
	assert.True(t, found, "Created label should be in list")

	updatedLabel, err := client.UpdateLabel(testProjectID, label.ID, plane.UpdateLabelRequest{
		Name:        labelName + "-updated",
		Description: "Updated description",
		Color:       "#3366FF",
	})
	require.NoError(t, err)
	assert.Contains(t, updatedLabel.Name, "updated")
	assert.Equal(t, "#3366FF", updatedLabel.Color)
}

// ============================================
// STATES TESTS
// ============================================

func TestStatesLifecycle(t *testing.T) {
	client := setupComprehensiveTest(t)

	stateName := fmt.Sprintf("Test State %d", time.Now().UnixNano())
	state, err := client.CreateState(testProjectID, plane.CreateStateRequest{
		Name:        stateName,
		Description: "Test state from integration tests",
		Color:       "#00FF00",
		Group:       "backlog",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, state.ID)
	t.Logf("Created state: %s (%s)", state.Name, state.ID)

	defer func() {
		err := client.DeleteState(testProjectID, state.ID)
		if err != nil {
			t.Logf("Warning: failed to delete state: %v", err)
		}
	}()

	fetchedState, err := client.GetState(testProjectID, state.ID)
	require.NoError(t, err)
	assert.Equal(t, state.ID, fetchedState.ID)

	states, err := client.ListStates(testProjectID)
	require.NoError(t, err)
	found := false
	for _, s := range states {
		if s.ID == state.ID {
			found = true
			break
		}
	}
	assert.True(t, found, "Created state should be in list")

	updatedState, err := client.UpdateState(testProjectID, state.ID, plane.UpdateStateRequest{
		Name:        stateName + " Updated",
		Description: "Updated description",
		Color:       "#FF0000",
	})
	require.NoError(t, err)
	assert.Contains(t, updatedState.Name, "Updated")
}

// ============================================
// ISSUE TYPES TESTS
// ============================================

func TestIssueTypesLifecycle(t *testing.T) {
	client := setupComprehensiveTest(t)
	project := requireTestProject(t, client)
	if !project.IsIssueTypeEnabled {
		t.Skip("issue types are disabled for this project")
	}

	typeName := fmt.Sprintf("Test Type %d", time.Now().UnixNano())
	issueType, err := client.CreateIssueType(testProjectID, plane.CreateIssueTypeRequest{
		Name:        typeName,
		Description: "Test issue type from integration tests",
		IsActive:    true,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, issueType.ID)
	t.Logf("Created issue type: %s (%s)", issueType.Name, issueType.ID)

	defer func() {
		err := client.DeleteIssueType(testProjectID, issueType.ID)
		if err != nil {
			t.Logf("Warning: failed to delete issue type: %v", err)
		}
	}()

	fetchedType, err := client.GetIssueType(testProjectID, issueType.ID)
	require.NoError(t, err)
	assert.Equal(t, issueType.ID, fetchedType.ID)

	types, err := client.ListIssueTypes(testProjectID)
	require.NoError(t, err)
	found := false
	for _, it := range types {
		if it.ID == issueType.ID {
			found = true
			break
		}
	}
	assert.True(t, found, "Created issue type should be in list")

	updatedType, err := client.UpdateIssueType(testProjectID, issueType.ID, plane.UpdateIssueTypeRequest{
		Name:        typeName + " Updated",
		Description: "Updated description",
	})
	require.NoError(t, err)
	assert.Contains(t, updatedType.Name, "Updated")
}

// ============================================
// PROJECTS TESTS
// ============================================

func TestProjectsListAndGet(t *testing.T) {
	client := setupComprehensiveTest(t)

	projects, err := client.ListProjects()
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(projects), 1)
	t.Logf("Found %d projects", len(projects))

	project, err := client.GetProject(testProjectID)
	require.NoError(t, err)
	assert.Equal(t, testProjectID, project.ID)
	t.Logf("Project: %s (%s)", project.Name, project.Identifier)
}

func TestProjectMembers(t *testing.T) {
	client := setupComprehensiveTest(t)

	members, err := client.GetProjectMembers(testProjectID)
	require.NoError(t, err)
	t.Logf("Project has %d members", len(members))

	for _, m := range members {
		t.Logf("  - %s (%s)", m.DisplayName, m.ID)
	}
}

// ============================================
// ISSUE ADVANCED TESTS
// ============================================

func TestIssueWithParent(t *testing.T) {
	client := setupComprehensiveTest(t)

	parent := createTestIssueForModule(t, client, "Parent Issue")
	defer cleanupTestIssue(t, client, parent.ID)

	child := createTestIssueForModule(t, client, "Child Issue")
	defer cleanupTestIssue(t, client, child.ID)

	updatedChild, err := client.UpdateIssue(testProjectID, child.ID, plane.UpdateIssueRequest{
		Parent: parent.ID,
	})
	require.NoError(t, err)
	assert.Equal(t, parent.ID, updatedChild.Parent)
	t.Logf("Set parent relationship: %s -> %s", child.ID, parent.ID)
}

func TestIssueWithCycle(t *testing.T) {
	client := setupComprehensiveTest(t)

	cycle, err := client.CreateCycle(testProjectID, plane.CreateCycleRequest{
		Name:      fmt.Sprintf("Issue Cycle Test %d", time.Now().UnixNano()),
		StartDate: time.Now().Format("2006-01-02"),
		EndDate:   time.Now().Add(7 * 24 * time.Hour).Format("2006-01-02"),
	})
	require.NoError(t, err)
	defer client.DeleteCycle(testProjectID, cycle.ID)

	issue := createTestIssueForModule(t, client, "Cycle Issue Test")
	defer cleanupTestIssue(t, client, issue.ID)

	err = client.AddIssuesToCycle(testProjectID, cycle.ID, []string{issue.ID})
	require.NoError(t, err)

	cycleIssues, err := client.ListCycleIssues(testProjectID, cycle.ID)
	require.NoError(t, err)
	found := false
	for _, item := range cycleIssues {
		if item.ID == issue.ID {
			found = true
			break
		}
	}
	assert.True(t, found, "expected issue to be associated with the cycle")
	t.Logf("Issue %s is now in cycle %s", issue.ID, cycle.ID)
}

func TestIssueWithModule(t *testing.T) {
	client := setupComprehensiveTest(t)

	module, err := client.CreateModule(testProjectID, plane.CreateModuleRequest{
		Name: fmt.Sprintf("Issue Module Test %d", time.Now().UnixNano()),
	})
	require.NoError(t, err)
	defer client.DeleteModule(testProjectID, module.ID)

	issue := createTestIssueForModule(t, client, "Module Issue Test")
	defer cleanupTestIssue(t, client, issue.ID)

	err = client.AddIssuesToModule(testProjectID, module.ID, []string{issue.ID})
	require.NoError(t, err)

	moduleIssues, err := client.ListModuleIssues(testProjectID, module.ID)
	require.NoError(t, err)
	found := false
	for _, item := range moduleIssues {
		if item.ID == issue.ID {
			found = true
			break
		}
	}
	assert.True(t, found, "expected issue to be associated with the module")
	t.Logf("Issue %s is now in module %s", issue.ID, module.ID)
}

func TestIssueSearch(t *testing.T) {
	client := setupComprehensiveTest(t)

	uniqueName := fmt.Sprintf("Searchable Issue %d", time.Now().UnixNano())
	issue := createTestIssueForModule(t, client, uniqueName)
	defer cleanupTestIssue(t, client, issue.ID)

	results, err := client.SearchIssues(uniqueName)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(results), 1)
	t.Logf("Search found %d results", len(results))
}

func TestIssueListWithPagination(t *testing.T) {
	client := setupComprehensiveTest(t)

	opts := IssueListOptions{
		Limit: 5,
	}
	issues, pagination, err := client.ListIssues(testProjectID, opts)
	require.NoError(t, err)
	assert.NotNil(t, pagination)
	t.Logf("Found %d issues (page shows %d, total: %d)", len(issues), len(issues), pagination.TotalResults)
}

func TestIssueGetByIdentifier(t *testing.T) {
	client := setupComprehensiveTest(t)

	issue := createTestIssueForModule(t, client, "Identifier Test Issue")
	defer cleanupTestIssue(t, client, issue.ID)

	fetched, err := client.GetIssueBySequenceID(testProjectID, issue.SequenceID)
	require.NoError(t, err)
	assert.Equal(t, issue.ID, fetched.ID)
	t.Logf("Fetched issue by sequence_id %d: %s", issue.SequenceID, fetched.ID)
}

// ============================================
// ATTACHMENTS TESTS
// ============================================

func TestAttachmentsList(t *testing.T) {
	client := setupComprehensiveTest(t)

	issue := createTestIssueForModule(t, client, "Attachments Test Issue")
	defer cleanupTestIssue(t, client, issue.ID)

	attachments, err := client.ListAttachments(testProjectID, issue.ID)
	require.NoError(t, err)
	t.Logf("Issue has %d attachments", len(attachments))
}

func TestAttachmentUploadCredentials(t *testing.T) {
	client := setupComprehensiveTest(t)

	issue := createTestIssueForModule(t, client, "Upload Creds Test Issue")
	defer cleanupTestIssue(t, client, issue.ID)

	creds, err := client.GetUploadCredentials(testProjectID, issue.ID, "test-file.png", 1024)
	require.NoError(t, err)
	assert.NotEmpty(t, creds.UploadTarget().URL)
	assert.NotEmpty(t, creds.Attachment.ID)
	t.Logf("Got upload credentials for attachment %s", creds.Attachment.ID)
}

func TestAttachmentUploadLifecycle(t *testing.T) {
	client := setupComprehensiveTest(t)

	issue := createTestIssueForModule(t, client, "Attachment Upload Test Issue")
	defer cleanupTestIssue(t, client, issue.ID)

	filePath := filepath.Join(t.TempDir(), "upload.png")
	err := os.WriteFile(filePath, []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4,
		0x89, 0x00, 0x00, 0x00, 0x0D, 0x49, 0x44, 0x41,
		0x54, 0x78, 0x9C, 0x63, 0xF8, 0xCF, 0xC0, 0x00,
		0x00, 0x03, 0x01, 0x01, 0x00, 0x18, 0xDD, 0x8D,
		0x18, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E,
		0x44, 0xAE, 0x42, 0x60, 0x82,
	}, 0o600)
	require.NoError(t, err)

	attachment, err := client.UploadAttachment(testProjectID, issue.ID, filePath)
	require.NoError(t, err)
	assert.NotEmpty(t, attachment.ID)
	assert.True(t, attachment.IsUploaded)

	fetched, err := client.GetAttachment(testProjectID, issue.ID, attachment.ID)
	require.NoError(t, err)
	assert.Equal(t, attachment.ID, fetched.ID)
	assert.True(t, fetched.IsUploaded)
}

// ============================================
// ARCHIVE/UNARCHIVE TESTS
// ============================================

func TestCycleArchive(t *testing.T) {
	client := setupComprehensiveTest(t)

	cycle, err := client.CreateCycle(testProjectID, plane.CreateCycleRequest{
		Name:      fmt.Sprintf("Archive Test Cycle %d", time.Now().UnixNano()),
		StartDate: "2024-01-01",
		EndDate:   "2024-01-07",
	})
	require.NoError(t, err)
	t.Logf("Created cycle: %s", cycle.ID)

	err = client.ArchiveCycle(testProjectID, cycle.ID)
	require.NoError(t, err)
	t.Log("Archived cycle")
	defer client.DeleteCycle(testProjectID, cycle.ID)
}

func TestModuleArchive(t *testing.T) {
	client := setupComprehensiveTest(t)

	module, err := client.CreateModule(testProjectID, plane.CreateModuleRequest{
		Name:   fmt.Sprintf("Archive Test Module %d", time.Now().UnixNano()),
		Status: "completed",
	})
	require.NoError(t, err)
	t.Logf("Created module: %s", module.ID)

	err = client.ArchiveModule(testProjectID, module.ID)
	require.NoError(t, err)
	t.Log("Archived module")

	modules, err := client.ListModules(testProjectID, false)
	require.NoError(t, err)
	for _, m := range modules {
		assert.NotEqual(t, module.ID, m.ID, "Archived module should not appear in active list")
	}

	modulesWithArchived, err := client.ListModules(testProjectID, true)
	require.NoError(t, err)
	found := false
	for _, m := range modulesWithArchived {
		if m.ID == module.ID {
			found = true
			break
		}
	}
	assert.True(t, found, "Archived module should appear when archived=true")

	defer client.DeleteModule(testProjectID, module.ID)
}

// ============================================
// ISSUE WITH DATES TESTS
// ============================================

func TestIssueWithDates(t *testing.T) {
	client := setupComprehensiveTest(t)

	issue := createTestIssueForModule(t, client, "Dates Test Issue")
	defer cleanupTestIssue(t, client, issue.ID)

	startDate := "2024-01-01"
	targetDate := "2024-12-31"

	updated, err := client.UpdateIssue(testProjectID, issue.ID, plane.UpdateIssueRequest{
		StartDate:  startDate,
		TargetDate: targetDate,
	})
	require.NoError(t, err)
	assert.Equal(t, startDate, updated.StartDate)
	assert.Equal(t, targetDate, updated.TargetDate)
	t.Logf("Set dates: start=%s, target=%s", updated.StartDate, updated.TargetDate)
}

func TestIssueWithEstimatePoint(t *testing.T) {
	client := setupComprehensiveTest(t)
	project := requireTestProject(t, client)
	if project.Estimate == nil {
		t.Skip("estimate points are disabled for this project")
	}

	issue := createTestIssueForModule(t, client, "Estimate Test Issue")
	defer cleanupTestIssue(t, client, issue.ID)

	updated, err := client.UpdateIssue(testProjectID, issue.ID, plane.UpdateIssueRequest{
		EstimatePoint: 5,
	})
	require.NoError(t, err)
	assert.Equal(t, 5, updated.EstimatePoint)
	t.Logf("Set estimate point: %d", updated.EstimatePoint)
}

// ============================================
// ISSUE WITH HTML DESCRIPTION TESTS
// ============================================

func TestIssueWithHTMLDescription(t *testing.T) {
	client := setupComprehensiveTest(t)

	issue, err := client.CreateIssue(testProjectID, plane.CreateIssueRequest{
		Name:            fmt.Sprintf("HTML Desc Test %d", time.Now().UnixNano()),
		DescriptionHTML: "<h1>Test</h1><p>This is <strong>HTML</strong> description</p>",
		Priority:        "medium",
	})
	require.NoError(t, err)
	defer cleanupTestIssue(t, client, issue.ID)

	fetched, err := client.GetIssue(testProjectID, issue.ID)
	require.NoError(t, err)
	assert.NotEmpty(t, fetched.DescriptionHTML)
	t.Logf("Issue created with HTML description")
}

// ============================================
// BULK OPERATIONS TESTS
// ============================================

func TestBulkAddIssuesToCycle(t *testing.T) {
	client := setupComprehensiveTest(t)

	cycle, err := client.CreateCycle(testProjectID, plane.CreateCycleRequest{
		Name: fmt.Sprintf("Bulk Cycle Test %d", time.Now().UnixNano()),
	})
	require.NoError(t, err)
	defer client.DeleteCycle(testProjectID, cycle.ID)

	issue1 := createTestIssueForModule(t, client, "Bulk Issue 1")
	defer cleanupTestIssue(t, client, issue1.ID)

	issue2 := createTestIssueForModule(t, client, "Bulk Issue 2")
	defer cleanupTestIssue(t, client, issue2.ID)

	err = client.AddIssuesToCycle(testProjectID, cycle.ID, []string{issue1.ID, issue2.ID})
	require.NoError(t, err)

	cycleIssues, err := client.ListCycleIssues(testProjectID, cycle.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(cycleIssues), 2)
	t.Logf("Cycle has %d issues after bulk add", len(cycleIssues))
}

func TestBulkAddIssuesToModule(t *testing.T) {
	client := setupComprehensiveTest(t)

	module, err := client.CreateModule(testProjectID, plane.CreateModuleRequest{
		Name: fmt.Sprintf("Bulk Module Test %d", time.Now().UnixNano()),
	})
	require.NoError(t, err)
	defer client.DeleteModule(testProjectID, module.ID)

	issue1 := createTestIssueForModule(t, client, "Bulk Module Issue 1")
	defer cleanupTestIssue(t, client, issue1.ID)

	issue2 := createTestIssueForModule(t, client, "Bulk Module Issue 2")
	defer cleanupTestIssue(t, client, issue2.ID)

	err = client.AddIssuesToModule(testProjectID, module.ID, []string{issue1.ID, issue2.ID})
	require.NoError(t, err)

	moduleIssues, err := client.ListModuleIssues(testProjectID, module.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(moduleIssues), 2)
	t.Logf("Module has %d issues after bulk add", len(moduleIssues))
}
