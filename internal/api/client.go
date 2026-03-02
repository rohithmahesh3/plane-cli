package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/rohithmahesh3/plane-cli/internal/config"
)

const (
	DefaultTimeout = 30 * time.Second
	APIVersion     = "v1"
)

type Client struct {
	HTTPClient *http.Client
	BaseURL    string
	APIKey     string
	Workspace  string
}

type Pagination struct {
	NextCursor      string `json:"next_cursor"`
	PrevCursor      string `json:"prev_cursor"`
	NextPageResults bool   `json:"next_page_results"`
	PrevPageResults bool   `json:"prev_page_results"`
	Count           int    `json:"count"`
	TotalPages      int    `json:"total_pages"`
	TotalResults    int    `json:"total_results"`
}

type Response struct {
	Pagination
	Results json.RawMessage `json:"results"`
}

func NewClient() (*Client, error) {
	apiKey, err := config.GetAPIKey()
	if err != nil {
		return nil, fmt.Errorf("not authenticated. Run 'plane auth login' first")
	}

	workspace := config.Cfg.DefaultWorkspace
	if workspace == "" {
		return nil, fmt.Errorf("no default workspace set. Use --workspace flag or set default workspace")
	}

	baseURL := config.Cfg.APIHost
	if baseURL == "" {
		baseURL = config.DefaultAPIHost
	}

	return &Client{
		HTTPClient: &http.Client{Timeout: DefaultTimeout},
		BaseURL:    baseURL,
		APIKey:     apiKey,
		Workspace:  workspace,
	}, nil
}

func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	urlStr := fmt.Sprintf("%s/api/%s%s", c.BaseURL, APIVersion, path)

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest(method, urlStr, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-Key", c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) error {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	if v != nil && len(body) > 0 {
		return json.Unmarshal(body, v)
	}

	return nil
}

func (c *Client) Get(path string, query url.Values, v interface{}) error {
	if query != nil {
		path = path + "?" + query.Encode()
	}

	req, err := c.NewRequest("GET", path, nil)
	if err != nil {
		return err
	}

	return c.Do(req, v)
}

func (c *Client) Post(path string, body interface{}, v interface{}) error {
	req, err := c.NewRequest("POST", path, body)
	if err != nil {
		return err
	}

	return c.Do(req, v)
}

func (c *Client) Patch(path string, body interface{}, v interface{}) error {
	req, err := c.NewRequest("PATCH", path, body)
	if err != nil {
		return err
	}

	return c.Do(req, v)
}

func (c *Client) Delete(path string) error {
	req, err := c.NewRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	return c.Do(req, nil)
}

func (c *Client) SetWorkspace(workspace string) {
	c.Workspace = workspace
}
