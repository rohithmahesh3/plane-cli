package api

import (
	"fmt"

	"github.com/rohithmahesh3/plane-cli/pkg/plane"
)

// ListModules retrieves all modules for a project
func (c *Client) ListModules(projectID string, archived bool) ([]plane.Module, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/modules/", c.Workspace, projectID)

	var response struct {
		Results []plane.Module `json:"results"`
	}

	if err := c.Get(path, nil, &response); err != nil {
		return nil, err
	}

	// Filter archived modules if needed
	if !archived {
		var activeModules []plane.Module
		for _, module := range response.Results {
			if module.ArchivedAt == "" {
				activeModules = append(activeModules, module)
			}
		}
		return activeModules, nil
	}

	return response.Results, nil
}

// GetModule retrieves a specific module by ID
func (c *Client) GetModule(projectID, moduleID string) (*plane.Module, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/modules/%s/", c.Workspace, projectID, moduleID)

	var module plane.Module
	if err := c.Get(path, nil, &module); err != nil {
		return nil, err
	}

	return &module, nil
}

// CreateModule creates a new module in a project
func (c *Client) CreateModule(projectID string, req plane.CreateModuleRequest) (*plane.Module, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/modules/", c.Workspace, projectID)

	var module plane.Module
	if err := c.Post(path, req, &module); err != nil {
		return nil, err
	}

	return &module, nil
}

// UpdateModule updates an existing module
func (c *Client) UpdateModule(projectID, moduleID string, req plane.UpdateModuleRequest) (*plane.Module, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/modules/%s/", c.Workspace, projectID, moduleID)

	var module plane.Module
	if err := c.Patch(path, req, &module); err != nil {
		return nil, err
	}

	return &module, nil
}

// DeleteModule removes a module from a project
func (c *Client) DeleteModule(projectID, moduleID string) error {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/modules/%s/", c.Workspace, projectID, moduleID)
	return c.Delete(path)
}

// ArchiveModule archives a module
func (c *Client) ArchiveModule(projectID, moduleID string) error {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/modules/%s/archive/", c.Workspace, projectID, moduleID)
	return c.Post(path, nil, nil)
}

// UnarchiveModule unarchives a module
func (c *Client) UnarchiveModule(projectID, moduleID string) error {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/modules/%s/unarchive/", c.Workspace, projectID, moduleID)
	return c.Post(path, nil, nil)
}

// ListModuleIssues retrieves all issues in a module
func (c *Client) ListModuleIssues(projectID, moduleID string) ([]plane.Issue, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/modules/%s/module-issues/", c.Workspace, projectID, moduleID)

	var response struct {
		Results []plane.Issue `json:"results"`
	}

	if err := c.Get(path, nil, &response); err != nil {
		return nil, err
	}

	return response.Results, nil
}

// AddIssuesToModule adds issues to a module
func (c *Client) AddIssuesToModule(projectID, moduleID string, issueIDs []string) error {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/modules/%s/module-issues/", c.Workspace, projectID, moduleID)

	req := struct {
		Issues []string `json:"issues"`
	}{
		Issues: issueIDs,
	}

	return c.Post(path, req, nil)
}

// RemoveIssueFromModule removes an issue from a module
func (c *Client) RemoveIssueFromModule(projectID, moduleID, issueID string) error {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/modules/%s/module-issues/%s/", c.Workspace, projectID, moduleID, issueID)
	return c.Delete(path)
}
