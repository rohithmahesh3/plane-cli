package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/rohithmahesh3/plane-cli/pkg/plane"
	"github.com/stretchr/testify/require"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestGetIntakeIssueFallsBackFromIntakeIDToIssueID(t *testing.T) {
	client := &Client{
		HTTPClient: &http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				var (
					status = http.StatusOK
					body   string
				)

				switch r.URL.Path {
				case "/api/v1/workspaces/test-workspace/projects/test-project/intake-issues/intake-1/":
					status = http.StatusNotFound
					body = `{"error":"not found"}`
				case "/api/v1/workspaces/test-workspace/projects/test-project/intake-issues/":
					response := struct {
						Results []plane.IntakeIssue `json:"results"`
					}{
						Results: []plane.IntakeIssue{
							{ID: "intake-1", Issue: "issue-1"},
						},
					}
					payload, err := json.Marshal(response)
					require.NoError(t, err)
					body = string(payload)
				case "/api/v1/workspaces/test-workspace/projects/test-project/intake-issues/issue-1/":
					intake := plane.IntakeIssue{ID: "intake-1", Issue: "issue-1"}
					payload, err := json.Marshal(intake)
					require.NoError(t, err)
					body = string(payload)
				default:
					t.Fatalf("unexpected request path %s", r.URL.Path)
				}

				return &http.Response{
					StatusCode: status,
					Header:     make(http.Header),
					Body:       io.NopCloser(strings.NewReader(body)),
				}, nil
			}),
		},
		BaseURL:   "https://example.test",
		APIKey:    "test-key",
		Workspace: "test-workspace",
	}

	intake, err := client.GetIntakeIssue("test-project", "intake-1")
	require.NoError(t, err)
	require.Equal(t, "intake-1", intake.ID)
	require.Equal(t, "issue-1", intake.Issue)
}
