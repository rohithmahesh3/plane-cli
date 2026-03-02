package api

import (
	"fmt"

	"github.com/rohithmahesh3/plane-cli/pkg/plane"
)

// ListWorklogs retrieves all worklogs for an issue
func (c *Client) ListWorklogs(projectID, issueID string) ([]plane.Worklog, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/worklogs/", c.Workspace, projectID, issueID)

	var response struct {
		Results []plane.Worklog `json:"results"`
	}

	if err := c.Get(path, nil, &response); err != nil {
		return nil, err
	}

	return response.Results, nil
}

// GetWorklog retrieves a specific worklog
func (c *Client) GetWorklog(projectID, issueID, worklogID string) (*plane.Worklog, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/worklogs/%s/", c.Workspace, projectID, issueID, worklogID)

	var worklog plane.Worklog
	if err := c.Get(path, nil, &worklog); err != nil {
		return nil, err
	}

	return &worklog, nil
}

// CreateWorklog creates a new worklog for an issue
func (c *Client) CreateWorklog(projectID, issueID string, req plane.CreateWorklogRequest) (*plane.Worklog, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/worklogs/", c.Workspace, projectID, issueID)

	var worklog plane.Worklog
	if err := c.Post(path, req, &worklog); err != nil {
		return nil, err
	}

	return &worklog, nil
}

// UpdateWorklog updates an existing worklog
func (c *Client) UpdateWorklog(projectID, issueID, worklogID string, req plane.UpdateWorklogRequest) (*plane.Worklog, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/worklogs/%s/", c.Workspace, projectID, issueID, worklogID)

	var worklog plane.Worklog
	if err := c.Patch(path, req, &worklog); err != nil {
		return nil, err
	}

	return &worklog, nil
}

// DeleteWorklog removes a worklog from an issue
func (c *Client) DeleteWorklog(projectID, issueID, worklogID string) error {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/worklogs/%s/", c.Workspace, projectID, issueID, worklogID)
	return c.Delete(path)
}

// GetTotalWorklogTime sums a work item's worklogs.
// The current public Plane API exposes project-level total-worklogs, not a
// per-work-item total endpoint, so this is derived from the item worklogs list.
func (c *Client) GetTotalWorklogTime(projectID, issueID string) (int, error) {
	worklogs, err := c.ListWorklogs(projectID, issueID)
	if err != nil {
		return 0, err
	}

	total := 0
	for _, worklog := range worklogs {
		total += worklog.Duration
	}

	return total, nil
}
