package api

import (
	"fmt"
	"sort"
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

func suggestWorkspaceMembers(members []plane.User, query string, limit int) []plane.User {
	normalizedQuery := normalizeMemberSearchValue(query)
	if normalizedQuery == "" || len(members) == 0 {
		return nil
	}

	type scoredMember struct {
		member plane.User
		score  int
	}

	seen := make(map[string]struct{}, len(members))
	scored := make([]scoredMember, 0, len(members))
	for _, member := range members {
		key := member.ID
		if key == "" {
			key = member.Email + ":" + member.DisplayName
		}
		if _, ok := seen[key]; ok {
			continue
		}

		score := scoreWorkspaceMember(member, normalizedQuery)
		if score < 0 {
			continue
		}

		seen[key] = struct{}{}
		scored = append(scored, scoredMember{member: member, score: score})
	}

	sort.SliceStable(scored, func(i, j int) bool {
		if scored[i].score != scored[j].score {
			return scored[i].score < scored[j].score
		}

		left := normalizeMemberSearchValue(preferredMemberLabel(scored[i].member))
		right := normalizeMemberSearchValue(preferredMemberLabel(scored[j].member))
		if left != right {
			return left < right
		}

		return scored[i].member.ID < scored[j].member.ID
	})

	suggestions := make([]plane.User, 0, len(scored))
	for _, item := range scored {
		suggestions = append(suggestions, item.member)
	}

	return limitWorkspaceMembers(suggestions, limit)
}

func FormatWorkspaceMemberSuggestion(member plane.User) string {
	label := preferredMemberLabel(member)
	email := strings.TrimSpace(member.Email)
	if email != "" && !strings.EqualFold(email, label) {
		return fmt.Sprintf("%s (%s) [%s]", label, email, member.ID)
	}

	return fmt.Sprintf("%s [%s]", label, member.ID)
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

func scoreWorkspaceMember(member plane.User, normalizedQuery string) int {
	displayName := normalizeMemberSearchValue(member.DisplayName)
	email := normalizeMemberSearchValue(member.Email)
	fullName := normalizeMemberSearchValue(memberFullName(member))
	firstName := normalizeMemberSearchValue(member.FirstName)
	lastName := normalizeMemberSearchValue(member.LastName)
	id := normalizeMemberSearchValue(member.ID)

	switch {
	case displayName != "" && displayName == normalizedQuery:
		return 0
	case email != "" && email == normalizedQuery:
		return 1
	case fullName != "" && fullName == normalizedQuery:
		return 2
	case displayName != "" && strings.HasPrefix(displayName, normalizedQuery):
		return 3
	case email != "" && strings.HasPrefix(email, normalizedQuery):
		return 4
	case displayName != "" && strings.Contains(displayName, normalizedQuery):
		return 5
	case email != "" && strings.Contains(email, normalizedQuery):
		return 6
	case fullName != "" && strings.Contains(fullName, normalizedQuery):
		return 7
	case firstName != "" && strings.Contains(firstName, normalizedQuery):
		return 8
	case lastName != "" && strings.Contains(lastName, normalizedQuery):
		return 8
	case id != "" && strings.Contains(id, normalizedQuery):
		return 9
	default:
		return -1
	}
}

func limitWorkspaceMembers(members []plane.User, limit int) []plane.User {
	if limit <= 0 || len(members) <= limit {
		return members
	}

	return members[:limit]
}

func preferredMemberLabel(member plane.User) string {
	if label := strings.TrimSpace(member.DisplayName); label != "" {
		return label
	}

	if label := strings.TrimSpace(memberFullName(member)); label != "" {
		return label
	}

	if label := strings.TrimSpace(member.Email); label != "" {
		return label
	}

	return member.ID
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
