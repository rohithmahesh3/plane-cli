package api

import (
	"fmt"

	"github.com/rohithmahesh3/plane-cli/pkg/plane"
)

// ListWorkspacePages retrieves all pages at workspace level
func (c *Client) ListWorkspacePages() ([]plane.Page, error) {
	path := fmt.Sprintf("/workspaces/%s/pages/", c.Workspace)

	var response struct {
		Results []plane.Page `json:"results"`
	}

	if err := c.Get(path, nil, &response); err != nil {
		return nil, err
	}

	return response.Results, nil
}

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

// ListProjectPages retrieves all pages for a project
func (c *Client) ListProjectPages(projectID string) ([]plane.Page, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/pages/", c.Workspace, projectID)

	var response struct {
		Results []plane.Page `json:"results"`
	}

	if err := c.Get(path, nil, &response); err != nil {
		return nil, err
	}

	return response.Results, nil
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

// DeletePage removes a page
func (c *Client) DeletePage(projectID, pageID string) error {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/pages/%s/", c.Workspace, projectID, pageID)
	return c.Delete(path)
}
