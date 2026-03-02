package api

import (
	"fmt"
	"strings"
	"sync"

	"github.com/rohithmahesh3/plane-cli/pkg/plane"
)

var (
	membersCache []plane.User
	labelsCache  []plane.Label
	statesCache  []plane.State
	membersMu    sync.RWMutex
	labelsMu     sync.RWMutex
	statesMu     sync.RWMutex
)

// ResolveAssignees converts display names/emails to user UUIDs
// Accepts formats: @display_name, display_name, email, or UUID
func (c *Client) ResolveAssignees(projectID string, assignees []string) ([]string, error) {
	if len(assignees) == 0 {
		return nil, nil
	}

	members, err := c.getCachedWorkspaceMembers()
	if err != nil {
		return nil, fmt.Errorf("failed to get workspace members: %w", err)
	}

	resolved := make([]string, 0, len(assignees))
	for _, a := range assignees {
		username := strings.TrimPrefix(a, "@")

		if isValidUUID(username) {
			resolved = append(resolved, username)
			continue
		}

		found := false
		for _, m := range members {
			if m.DisplayName == username || m.Email == username || m.ID == username {
				resolved = append(resolved, m.ID)
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("assignee '%s' not found in workspace members", username)
		}
	}

	return resolved, nil
}

// ResolveLabels converts label names to label UUIDs
// Accepts label names or UUIDs
func (c *Client) ResolveLabels(projectID string, labels []string) ([]string, error) {
	if len(labels) == 0 {
		return nil, nil
	}

	allLabels, err := c.getCachedProjectLabels(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project labels: %w", err)
	}

	resolved := make([]string, 0, len(labels))
	for _, l := range labels {
		if isValidUUID(l) {
			resolved = append(resolved, l)
			continue
		}

		found := false
		for _, label := range allLabels {
			if strings.EqualFold(label.Name, l) {
				resolved = append(resolved, label.ID)
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("label '%s' not found in project. Create it first with 'plane label create'", l)
		}
	}

	return resolved, nil
}

// ResolveState converts state name to state UUID
// Accepts state name or UUID
func (c *Client) ResolveState(projectID string, stateName string) (string, error) {
	if stateName == "" {
		return "", nil
	}

	if isValidUUID(stateName) {
		return stateName, nil
	}

	allStates, err := c.getCachedProjectStates(projectID)
	if err != nil {
		return "", fmt.Errorf("failed to get project states: %w", err)
	}

	for _, s := range allStates {
		if strings.EqualFold(s.Name, stateName) {
			return s.ID, nil
		}
	}

	return "", fmt.Errorf("state '%s' not found in project", stateName)
}

// getCachedWorkspaceMembers returns cached workspace members or fetches them
func (c *Client) getCachedWorkspaceMembers() ([]plane.User, error) {
	membersMu.RLock()
	if membersCache != nil {
		defer membersMu.RUnlock()
		return membersCache, nil
	}
	membersMu.RUnlock()

	membersMu.Lock()
	defer membersMu.Unlock()

	if membersCache != nil {
		return membersCache, nil
	}

	members, err := c.GetWorkspaceMembers()
	if err != nil {
		return nil, err
	}

	membersCache = members
	return membersCache, nil
}

// getCachedProjectLabels returns cached project labels or fetches them
func (c *Client) getCachedProjectLabels(projectID string) ([]plane.Label, error) {
	labelsMu.RLock()
	if labelsCache != nil {
		defer labelsMu.RUnlock()
		return labelsCache, nil
	}
	labelsMu.RUnlock()

	labelsMu.Lock()
	defer labelsMu.Unlock()

	if labelsCache != nil {
		return labelsCache, nil
	}

	labels, err := c.ListLabels(projectID)
	if err != nil {
		return nil, err
	}

	labelsCache = labels
	return labelsCache, nil
}

// getCachedProjectStates returns cached project states or fetches them
func (c *Client) getCachedProjectStates(projectID string) ([]plane.State, error) {
	statesMu.RLock()
	if statesCache != nil {
		defer statesMu.RUnlock()
		return statesCache, nil
	}
	statesMu.RUnlock()

	statesMu.Lock()
	defer statesMu.Unlock()

	if statesCache != nil {
		return statesCache, nil
	}

	states, err := c.ListStates(projectID)
	if err != nil {
		return nil, err
	}

	statesCache = states
	return statesCache, nil
}

// ClearResolverCache clears all cached resolver data
func ClearResolverCache() {
	membersMu.Lock()
	membersCache = nil
	membersMu.Unlock()

	labelsMu.Lock()
	labelsCache = nil
	labelsMu.Unlock()

	statesMu.Lock()
	statesCache = nil
	statesMu.Unlock()
}

// isValidUUID checks if a string looks like a UUID
func isValidUUID(s string) bool {
	if len(s) != 36 {
		return false
	}
	for i, r := range s {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if r != '-' {
				return false
			}
		} else {
			if (r < '0' || r > '9') && (r < 'a' || r > 'f') && (r < 'A' || r > 'F') {
				return false
			}
		}
	}
	return true
}
