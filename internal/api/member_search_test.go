package api

import (
	"testing"

	"github.com/rohithmahesh3/plane-cli/pkg/plane"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilterWorkspaceMembersSubstring(t *testing.T) {
	members := []plane.User{
		{ID: "1", DisplayName: "Rohith", Email: "rohith@example.com"},
		{ID: "2", DisplayName: "Alice", Email: "alice@example.com"},
	}

	filtered := FilterWorkspaceMembers(members, "roh", false, 0)
	require.Len(t, filtered, 1)
	assert.Equal(t, "1", filtered[0].ID)
}

func TestFilterWorkspaceMembersExact(t *testing.T) {
	members := []plane.User{
		{ID: "1", DisplayName: "Rohith", Email: "rohith@example.com"},
		{ID: "2", DisplayName: "Rohit", Email: "rohit@example.com"},
	}

	filtered := FilterWorkspaceMembers(members, "rohith", true, 0)
	require.Len(t, filtered, 1)
	assert.Equal(t, "1", filtered[0].ID)
}

func TestFilterWorkspaceMembersFullNameAndLimit(t *testing.T) {
	members := []plane.User{
		{ID: "1", FirstName: "Rohith", LastName: "Mahesh"},
		{ID: "2", FirstName: "Rohith", LastName: "Kumar"},
	}

	filtered := FilterWorkspaceMembers(members, "rohith", false, 1)
	require.Len(t, filtered, 1)
	assert.Equal(t, "1", filtered[0].ID)
}

func TestSuggestWorkspaceMembersRankingAndDedup(t *testing.T) {
	members := []plane.User{
		{ID: "2", DisplayName: "Rohan", Email: "rohan@example.com"},
		{ID: "1", DisplayName: "Rohith", Email: "rohith@example.com"},
		{ID: "1", DisplayName: "Rohith", Email: "rohith@example.com"},
	}

	suggestions := suggestWorkspaceMembers(members, "rohi", 10)
	require.Len(t, suggestions, 1)
	assert.Equal(t, "1", suggestions[0].ID)
}

func TestResolveAssigneesIncludesSuggestions(t *testing.T) {
	ClearResolverCache()
	membersMu.Lock()
	membersCache = map[string][]plane.User{
		"workspace-a": {
			{ID: "12345678-1234-1234-1234-123456789012", DisplayName: "Rohith", Email: "rohith@example.com"},
			{ID: "87654321-4321-4321-4321-210987654321", DisplayName: "Alice", Email: "alice@example.com"},
		},
	}
	membersMu.Unlock()
	t.Cleanup(ClearResolverCache)

	client := &Client{Workspace: "workspace-a"}
	_, err := client.ResolveAssignees("proj-1", []string{"rohi"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Closest matches:")
	assert.Contains(t, err.Error(), "Rohith")
	assert.Contains(t, err.Error(), "plane-cli workspace members --search rohi")
}

func TestWorkspaceMembersCacheIsScopedByWorkspace(t *testing.T) {
	ClearResolverCache()
	membersMu.Lock()
	membersCache = map[string][]plane.User{
		"workspace-a": {
			{ID: "member-a", DisplayName: "Alice"},
		},
		"workspace-b": {
			{ID: "member-b", DisplayName: "Bob"},
		},
	}
	membersMu.Unlock()
	t.Cleanup(ClearResolverCache)

	membersA, err := (&Client{Workspace: "workspace-a"}).getCachedWorkspaceMembers()
	require.NoError(t, err)
	require.Len(t, membersA, 1)
	assert.Equal(t, "member-a", membersA[0].ID)

	membersB, err := (&Client{Workspace: "workspace-b"}).getCachedWorkspaceMembers()
	require.NoError(t, err)
	require.Len(t, membersB, 1)
	assert.Equal(t, "member-b", membersB[0].ID)
}

func TestProjectLabelsCacheIsScopedByWorkspaceAndProject(t *testing.T) {
	ClearResolverCache()
	labelsMu.Lock()
	labelsCache = map[projectCacheKey][]plane.Label{
		{workspace: "workspace-a", projectID: "project-1"}: {
			{ID: "label-a", Name: "Alpha"},
		},
		{workspace: "workspace-a", projectID: "project-2"}: {
			{ID: "label-b", Name: "Beta"},
		},
	}
	labelsMu.Unlock()
	t.Cleanup(ClearResolverCache)

	labelsProject1, err := (&Client{Workspace: "workspace-a"}).getCachedProjectLabels("project-1")
	require.NoError(t, err)
	require.Len(t, labelsProject1, 1)
	assert.Equal(t, "label-a", labelsProject1[0].ID)

	labelsProject2, err := (&Client{Workspace: "workspace-a"}).getCachedProjectLabels("project-2")
	require.NoError(t, err)
	require.Len(t, labelsProject2, 1)
	assert.Equal(t, "label-b", labelsProject2[0].ID)
}

func TestProjectStatesCacheIsScopedByWorkspaceAndProject(t *testing.T) {
	ClearResolverCache()
	statesMu.Lock()
	statesCache = map[projectCacheKey][]plane.State{
		{workspace: "workspace-a", projectID: "project-1"}: {
			{ID: "state-a", Name: "Backlog"},
		},
		{workspace: "workspace-b", projectID: "project-1"}: {
			{ID: "state-b", Name: "Todo"},
		},
	}
	statesMu.Unlock()
	t.Cleanup(ClearResolverCache)

	statesWorkspaceA, err := (&Client{Workspace: "workspace-a"}).getCachedProjectStates("project-1")
	require.NoError(t, err)
	require.Len(t, statesWorkspaceA, 1)
	assert.Equal(t, "state-a", statesWorkspaceA[0].ID)

	statesWorkspaceB, err := (&Client{Workspace: "workspace-b"}).getCachedProjectStates("project-1")
	require.NoError(t, err)
	require.Len(t, statesWorkspaceB, 1)
	assert.Equal(t, "state-b", statesWorkspaceB[0].ID)
}

func TestClearResolverCacheClearsAllScopedCaches(t *testing.T) {
	membersMu.Lock()
	membersCache = map[string][]plane.User{"workspace-a": {{ID: "member-a"}}}
	membersMu.Unlock()

	labelsMu.Lock()
	labelsCache = map[projectCacheKey][]plane.Label{
		{workspace: "workspace-a", projectID: "project-1"}: {{ID: "label-a"}},
	}
	labelsMu.Unlock()

	statesMu.Lock()
	statesCache = map[projectCacheKey][]plane.State{
		{workspace: "workspace-a", projectID: "project-1"}: {{ID: "state-a"}},
	}
	statesMu.Unlock()

	ClearResolverCache()

	assert.Nil(t, membersCache)
	assert.Nil(t, labelsCache)
	assert.Nil(t, statesCache)
}
