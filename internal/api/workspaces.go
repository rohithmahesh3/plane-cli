package api

import (
	"fmt"

	"github.com/rohithmahesh3/plane-cli/pkg/plane"
)

// GetUserInfo retrieves the current authenticated user info
// This can be used to verify authentication and get user details
func (c *Client) GetUserInfo() (*plane.User, error) {
	path := "/users/me/"

	var user plane.User
	if err := c.Get(path, nil, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// ListWorkspaces is NOT available in the Plane API
// The API does not have a /workspaces/ endpoint
// Use the workspace slug from configuration instead
func (c *Client) ListWorkspaces() ([]plane.Workspace, error) {
	return nil, fmt.Errorf("workspace listing is not available in the Plane API. " +
		"Please configure your workspace slug using 'plane config set workspace <slug>'")
}

// GetWorkspace is NOT available in the Plane API
// The API does not have a /workspaces/{slug}/ endpoint
func (c *Client) GetWorkspace(slug string) (*plane.Workspace, error) {
	return nil, fmt.Errorf("workspace details endpoint is not available in the Plane API. " +
		"Use 'plane project list' to see projects in the configured workspace")
}

// GetWorkspaceMembers retrieves all members of the workspace
func (c *Client) GetWorkspaceMembers() ([]plane.User, error) {
	path := fmt.Sprintf("/workspaces/%s/members/", c.Workspace)

	var members []plane.User
	if err := c.Get(path, nil, &members); err != nil {
		return nil, err
	}

	return members, nil
}
