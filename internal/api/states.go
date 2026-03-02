package api

import (
	"fmt"

	"github.com/rohithmahesh3/plane-cli/pkg/plane"
)

// ListStates retrieves all states for a project
func (c *Client) ListStates(projectID string) ([]plane.State, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/states/", c.Workspace, projectID)

	var response struct {
		Results []plane.State `json:"results"`
	}

	if err := c.Get(path, nil, &response); err != nil {
		return nil, err
	}

	return response.Results, nil
}

// GetState retrieves a specific state by ID
func (c *Client) GetState(projectID, stateID string) (*plane.State, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/states/%s/", c.Workspace, projectID, stateID)

	var state plane.State
	if err := c.Get(path, nil, &state); err != nil {
		return nil, err
	}

	return &state, nil
}

// CreateState creates a new state in a project
func (c *Client) CreateState(projectID string, req plane.CreateStateRequest) (*plane.State, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/states/", c.Workspace, projectID)

	var state plane.State
	if err := c.Post(path, req, &state); err != nil {
		return nil, err
	}

	return &state, nil
}

// UpdateState updates an existing state
func (c *Client) UpdateState(projectID, stateID string, req plane.UpdateStateRequest) (*plane.State, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/states/%s/", c.Workspace, projectID, stateID)

	var state plane.State
	if err := c.Patch(path, req, &state); err != nil {
		return nil, err
	}

	return &state, nil
}

// DeleteState removes a state from a project
func (c *Client) DeleteState(projectID, stateID string) error {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/states/%s/", c.Workspace, projectID, stateID)
	return c.Delete(path)
}
