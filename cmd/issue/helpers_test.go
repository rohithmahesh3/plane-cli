package issue

import (
	"strings"
	"testing"
)

func TestRenderDescriptionHTML(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty input",
			input: " \n\t ",
			want:  "",
		},
		{
			name:  "plain text paragraph",
			input: "hello world",
			want:  "<p>hello world</p>\n",
		},
		{
			name:  "multiple paragraphs and line breaks",
			input: "line 1\nline 2\n\nnext paragraph",
			want:  "<p>line 1<br />\nline 2</p>\n<p>next paragraph</p>\n",
		},
		{
			name:  "preserves existing html",
			input: "<script>alert(1)</script>\nplain",
			// Input starting with "<" is treated as existing HTML and returned unchanged
			want: "<script>alert(1)</script>\nplain",
		},
		{
			name:  "markdown with bold and list",
			input: "# Heading\n\n**bold text**\n\n- item 1\n- item 2",
			want:  "<h1>Heading</h1>\n<p><strong>bold text</strong></p>\n<ul>\n<li>item 1</li>\n<li>item 2</li>\n</ul>\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderDescriptionHTML(tt.input)
			// Normalize line endings for comparison
			got = strings.TrimSpace(got)
			want := strings.TrimSpace(tt.want)
			if got != want {
				t.Fatalf("renderDescriptionHTML() = %q, want %q", got, want)
			}
		})
	}
}
