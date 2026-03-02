package api

import (
	"fmt"

	"github.com/rohithmahesh3/plane-cli/pkg/plane"
)

// ListCycles retrieves all cycles for a project
func (c *Client) ListCycles(projectID string, archived bool) ([]plane.Cycle, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/cycles/", c.Workspace, projectID)

	var response struct {
		Results []plane.Cycle `json:"results"`
	}

	if err := c.Get(path, nil, &response); err != nil {
		return nil, err
	}

	// Filter archived cycles if needed
	if !archived {
		var activeCycles []plane.Cycle
		for _, cycle := range response.Results {
			if cycle.ArchivedAt == "" {
				activeCycles = append(activeCycles, cycle)
			}
		}
		return activeCycles, nil
	}

	return response.Results, nil
}

// GetCycle retrieves a specific cycle by ID
func (c *Client) GetCycle(projectID, cycleID string) (*plane.Cycle, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/cycles/%s/", c.Workspace, projectID, cycleID)

	var cycle plane.Cycle
	if err := c.Get(path, nil, &cycle); err != nil {
		return nil, err
	}

	return &cycle, nil
}

// CreateCycle creates a new cycle in a project
func (c *Client) CreateCycle(projectID string, req plane.CreateCycleRequest) (*plane.Cycle, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/cycles/", c.Workspace, projectID)

	var cycle plane.Cycle
	if err := c.Post(path, req, &cycle); err != nil {
		return nil, err
	}

	return &cycle, nil
}

// UpdateCycle updates an existing cycle
func (c *Client) UpdateCycle(projectID, cycleID string, req plane.UpdateCycleRequest) (*plane.Cycle, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/cycles/%s/", c.Workspace, projectID, cycleID)

	var cycle plane.Cycle
	if err := c.Patch(path, req, &cycle); err != nil {
		return nil, err
	}

	return &cycle, nil
}

// DeleteCycle removes a cycle from a project
func (c *Client) DeleteCycle(projectID, cycleID string) error {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/cycles/%s/", c.Workspace, projectID, cycleID)
	return c.Delete(path)
}

// ArchiveCycle archives a cycle
func (c *Client) ArchiveCycle(projectID, cycleID string) error {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/cycles/%s/archive/", c.Workspace, projectID, cycleID)
	return c.Post(path, nil, nil)
}

// UnarchiveCycle unarchives a cycle
func (c *Client) UnarchiveCycle(projectID, cycleID string) error {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/cycles/%s/unarchive/", c.Workspace, projectID, cycleID)
	return c.Post(path, nil, nil)
}

// ListCycleIssues retrieves all issues in a cycle
func (c *Client) ListCycleIssues(projectID, cycleID string) ([]plane.Issue, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/cycles/%s/cycle-issues/", c.Workspace, projectID, cycleID)

	var response struct {
		Results []plane.Issue `json:"results"`
	}

	if err := c.Get(path, nil, &response); err != nil {
		return nil, err
	}

	return response.Results, nil
}

// AddIssuesToCycle adds issues to a cycle
func (c *Client) AddIssuesToCycle(projectID, cycleID string, issueIDs []string) error {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/cycles/%s/cycle-issues/", c.Workspace, projectID, cycleID)

	req := struct {
		Issues []string `json:"issues"`
	}{
		Issues: issueIDs,
	}

	return c.Post(path, req, nil)
}

// RemoveIssueFromCycle removes an issue from a cycle
func (c *Client) RemoveIssueFromCycle(projectID, cycleID, issueID string) error {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/cycles/%s/cycle-issues/%s/", c.Workspace, projectID, cycleID, issueID)
	return c.Delete(path)
}
