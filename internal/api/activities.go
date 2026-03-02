package api

import (
	"fmt"

	"github.com/rohithmahesh3/plane-cli/pkg/plane"
)

// ListActivities retrieves all activities for an issue
func (c *Client) ListActivities(projectID, issueID string) ([]plane.Activity, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/activities/", c.Workspace, projectID, issueID)

	var response struct {
		Results []plane.Activity `json:"results"`
	}

	if err := c.Get(path, nil, &response); err != nil {
		return nil, err
	}

	return response.Results, nil
}

// GetActivity retrieves a specific activity by ID
func (c *Client) GetActivity(projectID, issueID, activityID string) (*plane.Activity, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/work-items/%s/activities/%s/", c.Workspace, projectID, issueID, activityID)

	var activity plane.Activity
	if err := c.Get(path, nil, &activity); err != nil {
		return nil, err
	}

	return &activity, nil
}
