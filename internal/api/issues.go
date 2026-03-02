package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/rohithmahesh3/plane-cli/pkg/plane"
)

type IssueListOptions struct {
	State    string
	Priority string
	Assignee string
	Label    string
	Cycle    string
	Module   string
	Search   string
	PerPage  int
	Cursor   string
}

func (c *Client) ListIssues(projectID string, opts IssueListOptions) ([]plane.Issue, *Pagination, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/issues/", c.Workspace, projectID)

	query := url.Values{}
	if opts.State != "" {
		query.Set("state", opts.State)
	}
	if opts.Priority != "" {
		query.Set("priority", opts.Priority)
	}
	if opts.Assignee != "" {
		query.Set("assignee", opts.Assignee)
	}
	if opts.Label != "" {
		query.Set("label", opts.Label)
	}
	if opts.Cycle != "" {
		query.Set("cycle", opts.Cycle)
	}
	if opts.Module != "" {
		query.Set("module", opts.Module)
	}
	if opts.PerPage > 0 {
		query.Set("per_page", fmt.Sprintf("%d", opts.PerPage))
	}
	if opts.Cursor != "" {
		query.Set("cursor", opts.Cursor)
	}

	var response Response
	if err := c.Get(path, query, &response); err != nil {
		return nil, nil, err
	}

	var issues []plane.Issue
	if err := json.Unmarshal(response.Results, &issues); err != nil {
		return nil, nil, err
	}

	return issues, &response.Pagination, nil
}

func (c *Client) GetIssue(projectID, issueID string) (*plane.Issue, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/issues/%s/", c.Workspace, projectID, issueID)

	var issue plane.Issue
	if err := c.Get(path, nil, &issue); err != nil {
		return nil, err
	}

	return &issue, nil
}

func (c *Client) GetIssueByIdentifier(projectID string, sequenceID int) (*plane.Issue, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/issues/?sequence_id=%d", c.Workspace, projectID, sequenceID)

	var response Response
	if err := c.Get(path, nil, &response); err != nil {
		return nil, err
	}

	var issues []plane.Issue
	if err := json.Unmarshal(response.Results, &issues); err != nil {
		return nil, err
	}

	if len(issues) == 0 {
		return nil, fmt.Errorf("issue not found")
	}

	return &issues[0], nil
}

func (c *Client) CreateIssue(projectID string, req plane.CreateIssueRequest) (*plane.Issue, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/issues/", c.Workspace, projectID)

	var issue plane.Issue
	if err := c.Post(path, req, &issue); err != nil {
		return nil, err
	}

	return &issue, nil
}

func (c *Client) UpdateIssue(projectID, issueID string, req plane.UpdateIssueRequest) (*plane.Issue, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/issues/%s/", c.Workspace, projectID, issueID)

	var issue plane.Issue
	if err := c.Patch(path, req, &issue); err != nil {
		return nil, err
	}

	return &issue, nil
}

func (c *Client) DeleteIssue(projectID, issueID string) error {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/issues/%s/", c.Workspace, projectID, issueID)
	return c.Delete(path)
}

// SearchIssues searches for issues across the workspace
// Endpoint: GET /api/v1/workspaces/{workspace_slug}/issues/search/
// Returns: {"issues": [...]}
func (c *Client) SearchIssues(query string) ([]plane.Issue, error) {
	path := fmt.Sprintf("/workspaces/%s/issues/search/", c.Workspace)

	params := url.Values{}
	params.Set("search", query)

	// The search endpoint returns a different structure
	var response struct {
		Issues []plane.Issue `json:"issues"`
	}

	if err := c.Get(path, params, &response); err != nil {
		return nil, err
	}

	return response.Issues, nil
}
