package issue

import (
	"html"
	"regexp"
	"strconv"
	"strings"

	"github.com/rohithmahesh3/plane-cli/internal/api"
	"github.com/rohithmahesh3/plane-cli/pkg/plane"
)

var paragraphBreakPattern = regexp.MustCompile(`\n\s*\n+`)

func resolveIssue(client *api.Client, projectID, ref string) (*plane.Issue, error) {
	if seqID, err := strconv.Atoi(strings.TrimSpace(ref)); err == nil {
		return client.GetIssueBySequenceID(projectID, seqID)
	}

	if looksLikeUUID(ref) {
		return client.GetIssue(projectID, ref)
	}

	issue, err := client.GetIssue(projectID, ref)
	if err == nil {
		return issue, nil
	}

	return client.GetIssueByIdentifier(ref)
}

func resolveIssueID(client *api.Client, projectID, ref string) (string, error) {
	issue, err := resolveIssue(client, projectID, ref)
	if err != nil {
		return "", err
	}

	return issue.ID, nil
}

func looksLikeUUID(s string) bool {
	if len(s) != 36 {
		return false
	}

	for i, r := range s {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if r != '-' {
				return false
			}
			continue
		}

		if (r < '0' || r > '9') && (r < 'a' || r > 'f') && (r < 'A' || r > 'F') {
			return false
		}
	}

	return true
}

func renderDescriptionHTML(input string) string {
	normalized := strings.TrimSpace(strings.ReplaceAll(input, "\r\n", "\n"))
	if normalized == "" {
		return ""
	}

	if strings.HasPrefix(normalized, "<") {
		return normalized
	}

	paragraphs := paragraphBreakPattern.Split(normalized, -1)
	var builder strings.Builder

	for _, paragraph := range paragraphs {
		paragraph = strings.TrimSpace(paragraph)
		if paragraph == "" {
			continue
		}

		if builder.Len() > 0 {
			builder.WriteByte('\n')
		}

		builder.WriteString("<p>")
		builder.WriteString(strings.ReplaceAll(html.EscapeString(paragraph), "\n", "<br>"))
		builder.WriteString("</p>")
	}

	return builder.String()
}
