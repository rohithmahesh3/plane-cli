package api

import (
	"fmt"

	"github.com/rohithmahesh3/plane-cli/pkg/plane"
)

// ListIntakeIssues retrieves all intake issues for a project
func (c *Client) ListIntakeIssues(projectID string) ([]plane.IntakeIssue, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/intake-issues/", c.Workspace, projectID)

	var response struct {
		Results []plane.IntakeIssue `json:"results"`
	}

	if err := c.Get(path, nil, &response); err != nil {
		return nil, err
	}

	return response.Results, nil
}

// GetIntakeIssue retrieves a specific intake issue
func (c *Client) GetIntakeIssue(projectID, intakeID string) (*plane.IntakeIssue, error) {
	intake, err := c.getIntakeIssue(projectID, intakeID)
	if err == nil {
		return intake, nil
	}

	if !isAPIStatusError(err, 404) {
		return nil, err
	}

	intakeIssues, listErr := c.ListIntakeIssues(projectID)
	if listErr != nil {
		return nil, err
	}

	for _, item := range intakeIssues {
		if item.ID == intakeID && item.Issue != "" {
			return c.getIntakeIssue(projectID, item.Issue)
		}
	}

	return nil, err
}

func (c *Client) getIntakeIssue(projectID, intakeID string) (*plane.IntakeIssue, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/intake-issues/%s/", c.Workspace, projectID, intakeID)

	var intake plane.IntakeIssue
	if err := c.Get(path, nil, &intake); err != nil {
		return nil, err
	}

	return &intake, nil
}

// CreateIntakeIssue creates a new intake issue
func (c *Client) CreateIntakeIssue(projectID string, req plane.CreateIntakeIssueRequest) (*plane.IntakeIssue, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/intake-issues/", c.Workspace, projectID)

	var intake plane.IntakeIssue
	if err := c.Post(path, req, &intake); err != nil {
		return nil, err
	}

	return &intake, nil
}

// DeleteIntakeIssue removes an intake issue
func (c *Client) DeleteIntakeIssue(projectID, intakeID string) error {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/intake-issues/%s/", c.Workspace, projectID, intakeID)
	return c.Delete(path)
}
