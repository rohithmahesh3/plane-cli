package api

import (
	"fmt"

	"github.com/rohithmahesh3/plane-cli/pkg/plane"
)

func (c *Client) ListWorkspaces() ([]plane.Workspace, error) {
	path := "/workspaces/"
	
	var response struct {
		Results []plane.Workspace `json:"results"`
	}
	
	if err := c.Get(path, nil, &response); err != nil {
		return nil, err
	}
	
	return response.Results, nil
}

func (c *Client) GetWorkspace(slug string) (*plane.Workspace, error) {
	path := fmt.Sprintf("/workspaces/%s/", slug)
	
	var workspace plane.Workspace
	if err := c.Get(path, nil, &workspace); err != nil {
		return nil, err
	}
	
	return &workspace, nil
}
