package api

import (
	"strings"

	"github.com/rohithmahesh3/plane-cli/pkg/plane"
)

// FilterWorkspaceMembers performs client-side filtering for workspace members.
// It does not imply any server-side search support in the Plane API.
func FilterWorkspaceMembers(members []plane.User, query string, exact bool, limit int) []plane.User {
	if len(members) == 0 {
		return nil
	}

	if strings.TrimSpace(query) == "" {
		return limitWorkspaceMembers(members, limit)
	}

	normalizedQuery := normalizeMemberSearchValue(query)
	filtered := make([]plane.User, 0, len(members))

	for _, member := range members {
		if memberMatchesQuery(member, normalizedQuery, exact) {
			filtered = append(filtered, member)
		}
	}

	return limitWorkspaceMembers(filtered, limit)
}

func memberMatchesQuery(member plane.User, normalizedQuery string, exact bool) bool {
	for _, candidate := range memberSearchFields(member) {
		if candidate == "" {
			continue
		}

		if exact {
			if candidate == normalizedQuery {
				return true
			}
			continue
		}

		if strings.Contains(candidate, normalizedQuery) {
			return true
		}
	}

	return false
}

func limitWorkspaceMembers(members []plane.User, limit int) []plane.User {
	if limit <= 0 || len(members) <= limit {
		return members
	}

	return members[:limit]
}

func memberFullName(member plane.User) string {
	parts := make([]string, 0, 2)
	if first := strings.TrimSpace(member.FirstName); first != "" {
		parts = append(parts, first)
	}
	if last := strings.TrimSpace(member.LastName); last != "" {
		parts = append(parts, last)
	}

	return strings.Join(parts, " ")
}

func memberSearchFields(member plane.User) []string {
	fullName := memberFullName(member)

	return []string{
		normalizeMemberSearchValue(member.DisplayName),
		normalizeMemberSearchValue(member.Email),
		normalizeMemberSearchValue(member.FirstName),
		normalizeMemberSearchValue(member.LastName),
		normalizeMemberSearchValue(fullName),
		normalizeMemberSearchValue(member.ID),
	}
}

func normalizeMemberSearchValue(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}
