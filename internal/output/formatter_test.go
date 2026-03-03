package output

import (
	"reflect"
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

func TestTableFieldsRenameHTMLHeaders(t *testing.T) {
	val := reflect.ValueOf(struct {
		DescriptionHTML string `json:"description_html"`
	}{
		DescriptionHTML: "<p>Hello</p>",
	})

	fields := tableFields(val)
	require.Len(t, fields, 1)
	assert.Equal(t, "DescriptionMarkdown", fields[0].header)
	assert.True(t, fields[0].isHTML)
}

func TestFormatterGetRowConvertsHTMLValues(t *testing.T) {
	formatter := NewFormatter("table", false)
	val := reflect.ValueOf(struct {
		DescriptionHTML string `json:"description_html"`
	}{
		DescriptionHTML: "<p>Hello <strong>world</strong></p>",
	})

	row := formatter.getRow(tableFields(val))
	require.Len(t, row, 1)
	assert.Equal(t, "Hello **world**", row[0])
}
