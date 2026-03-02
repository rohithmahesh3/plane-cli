package api

import (
	"fmt"

	"github.com/rohithmahesh3/plane-cli/pkg/plane"
)

// GetWorkspacePage retrieves a specific workspace page
func (c *Client) GetWorkspacePage(pageID string) (*plane.Page, error) {
	path := fmt.Sprintf("/workspaces/%s/pages/%s/", c.Workspace, pageID)

	var page plane.Page
	if err := c.Get(path, nil, &page); err != nil {
		return nil, err
	}

	return &page, nil
}

// CreateWorkspacePage creates a new workspace page
func (c *Client) CreateWorkspacePage(req plane.CreatePageRequest) (*plane.Page, error) {
	path := fmt.Sprintf("/workspaces/%s/pages/", c.Workspace)

	var page plane.Page
	if err := c.Post(path, req, &page); err != nil {
		return nil, err
	}

	return &page, nil
}

// GetProjectPage retrieves a specific project page
func (c *Client) GetProjectPage(projectID, pageID string) (*plane.Page, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/pages/%s/", c.Workspace, projectID, pageID)

	var page plane.Page
	if err := c.Get(path, nil, &page); err != nil {
		return nil, err
	}

	return &page, nil
}

// CreateProjectPage creates a new page in a project
func (c *Client) CreateProjectPage(projectID string, req plane.CreatePageRequest) (*plane.Page, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/pages/", c.Workspace, projectID)

	var page plane.Page
	if err := c.Post(path, req, &page); err != nil {
		return nil, err
	}

	return &page, nil
}
