package api

import (
	"fmt"

	"github.com/rohithmahesh3/plane-cli/pkg/plane"
)

func (c *Client) ListProjects() ([]plane.Project, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/", c.Workspace)

	var response struct {
		Results []plane.Project `json:"results"`
	}

	if err := c.Get(path, nil, &response); err != nil {
		return nil, err
	}

	return response.Results, nil
}

func (c *Client) GetProject(projectID string) (*plane.Project, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/", c.Workspace, projectID)

	var project plane.Project
	if err := c.Get(path, nil, &project); err != nil {
		return nil, err
	}

	return &project, nil
}

func (c *Client) CreateProject(req plane.CreateProjectRequest) (*plane.Project, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/", c.Workspace)

	var project plane.Project
	if err := c.Post(path, req, &project); err != nil {
		return nil, err
	}

	return &project, nil
}

func (c *Client) DeleteProject(projectID string) error {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/", c.Workspace, projectID)
	return c.Delete(path)
}

func (c *Client) GetProjectMembers(projectID string) ([]plane.User, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/members/", c.Workspace, projectID)

	var members []plane.User
	if err := c.Get(path, nil, &members); err != nil {
		return nil, err
	}

	return members, nil
}
