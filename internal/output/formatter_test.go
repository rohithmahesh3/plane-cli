package output

import (
	"io"
	"os"
	"testing"

	"github.com/rohithmahesh3/plane-cli/pkg/plane"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeStructuredOutputRenamesHTMLFields(t *testing.T) {
	normalized, err := normalizeStructuredOutput(plane.Issue{
		Name:                "Probe",
		DescriptionHTML:     "<h1>Title</h1><p>Hello <strong>world</strong></p>",
		DescriptionStripped: "Title Hello world",
	})
	require.NoError(t, err)

	record, ok := normalized.(map[string]interface{})
	require.True(t, ok)

	assert.Equal(t, "Probe", record["name"])
	assert.Equal(t, "Title Hello world", record["description_stripped"])
	assert.NotContains(t, record, "description_html")
	assert.Equal(t, "# Title\n\nHello **world**", record["description_markdown"])
}

func TestNormalizeStructuredOutputRenamesCommentHTMLInSlices(t *testing.T) {
	normalized, err := normalizeStructuredOutput([]plane.Comment{
		{
			ID:          "1",
			CommentHTML: "<p>Hi <del>there</del></p>",
			CommentJSON: `{"x":1}`,
		},
	})
	require.NoError(t, err)

	items, ok := normalized.([]interface{})
	require.True(t, ok)
	require.Len(t, items, 1)

	record, ok := items[0].(map[string]interface{})
	require.True(t, ok)

	assert.NotContains(t, record, "comment_html")
	assert.Equal(t, "Hi ~~there~~", record["comment_markdown"])
	assert.Equal(t, `{"x":1}`, record["comment_json"])
}

func TestHTMLToMarkdownConvertsTables(t *testing.T) {
	markdown, err := htmlToMarkdown("<table><tr><th>A</th><th>B</th></tr><tr><td>1</td><td>2</td></tr></table>")
	require.NoError(t, err)

	assert.Contains(t, markdown, "| A | B |")
	assert.Contains(t, markdown, "| 1 | 2 |")
}

func TestFormatterPrintDefaultsToYAML(t *testing.T) {
	formatter := NewFormatter("", false)

	output, err := captureStdout(t, func() error {
		return formatter.Print(plane.Issue{
			Name:            "Probe",
			DescriptionHTML: "<p>Hello <strong>world</strong></p>",
		})
	})
	require.NoError(t, err)

	assert.Contains(t, output, "name: Probe")
	assert.Contains(t, output, "description_markdown: Hello **world**")
	assert.NotContains(t, output, "description_html")
}

func TestFormatterPrintRejectsTable(t *testing.T) {
	formatter := NewFormatter("table", false)

	err := formatter.Print(plane.Issue{Name: "Probe"})
	require.Error(t, err)
	assert.EqualError(t, err, `invalid output format "table": table output has been removed; supported formats are json, yaml`)
}

func TestValidateFormat(t *testing.T) {
	require.NoError(t, ValidateFormat(""))
	require.NoError(t, ValidateFormat("json"))
	require.NoError(t, ValidateFormat("yaml"))

	err := ValidateFormat("bogus")
	require.Error(t, err)
	assert.EqualError(t, err, `invalid output format "bogus": supported formats are json, yaml`)
}

func captureStdout(t *testing.T, fn func() error) (string, error) {
	t.Helper()

	originalStdout := os.Stdout
	reader, writer, err := os.Pipe()
	require.NoError(t, err)

	os.Stdout = writer
	runErr := fn()
	require.NoError(t, writer.Close())
	os.Stdout = originalStdout

	out, readErr := io.ReadAll(reader)
	require.NoError(t, readErr)
	require.NoError(t, reader.Close())

	return string(out), runErr
}
