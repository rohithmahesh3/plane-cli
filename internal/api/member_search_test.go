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
	membersCache = []plane.User{
		{ID: "12345678-1234-1234-1234-123456789012", DisplayName: "Rohith", Email: "rohith@example.com"},
		{ID: "87654321-4321-4321-4321-210987654321", DisplayName: "Alice", Email: "alice@example.com"},
	}
	membersMu.Unlock()
	t.Cleanup(ClearResolverCache)

	client := &Client{}
	_, err := client.ResolveAssignees("proj-1", []string{"rohi"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Closest matches:")
	assert.Contains(t, err.Error(), "Rohith")
	assert.Contains(t, err.Error(), "plane workspace members --search rohi")
}
