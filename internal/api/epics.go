package api

import (
	"fmt"

	"github.com/rohithmahesh3/plane-cli/pkg/plane"
)

// ListEpics retrieves all epics for a project
func (c *Client) ListEpics(projectID string) ([]plane.Epic, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/epics/", c.Workspace, projectID)

	var response struct {
		Results []plane.Epic `json:"results"`
	}

	if err := c.Get(path, nil, &response); err != nil {
		return nil, err
	}

	return response.Results, nil
}

// GetEpic retrieves a specific epic by ID
func (c *Client) GetEpic(projectID, epicID string) (*plane.Epic, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/epics/%s/", c.Workspace, projectID, epicID)

	var epic plane.Epic
	if err := c.Get(path, nil, &epic); err != nil {
		return nil, err
	}

	return &epic, nil
}
