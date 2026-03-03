package issue

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"

	"github.com/rohithmahesh3/plane-cli/internal/api"
	"github.com/rohithmahesh3/plane-cli/pkg/plane"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

var codeBlockPattern = regexp.MustCompile(`(?s)<pre><code([^>]*)>(.*?)</code></pre>`)

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

// renderDescriptionHTML converts Markdown text to HTML.
// If the input already appears to be HTML (starts with "<"), it returns it unchanged.
func renderDescriptionHTML(input string) string {
	normalized := strings.TrimSpace(strings.ReplaceAll(input, "\r\n", "\n"))
	if normalized == "" {
		return ""
	}

	// If input already looks like HTML, return as-is
	if strings.HasPrefix(normalized, "<") {
		return normalized
	}

	// Configure goldmark with common extensions for GitHub-flavored Markdown
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.Table,
			extension.Strikethrough,
			extension.Linkify,
			extension.TaskList,
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(normalized), &buf); err != nil {
		// Fallback: return escaped text wrapped in paragraph if parsing fails
		return "<p>" + normalized + "</p>"
	}

	return normalizeCodeBlockBlankLines(buf.String())
}

func normalizeCodeBlockBlankLines(rendered string) string {
	return codeBlockPattern.ReplaceAllStringFunc(rendered, func(block string) string {
		matches := codeBlockPattern.FindStringSubmatch(block)
		if len(matches) != 3 {
			return block
		}

		lines := strings.Split(matches[2], "\n")
		for i := 0; i < len(lines)-1; i++ {
			if lines[i] == "" {
				lines[i] = "<br />"
			}
		}

		return "<pre><code" + matches[1] + ">" + strings.Join(lines, "\n") + "</code></pre>"
	})
}
