package issue

import "testing"

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
			want:  "<p>hello world</p>",
		},
		{
			name:  "multiple paragraphs and line breaks",
			input: "line 1\nline 2\n\nnext paragraph",
			want:  "<p>line 1<br>line 2</p>\n<p>next paragraph</p>",
		},
		{
			name:  "escapes html-like text",
			input: "<script>alert(1)</script>\nplain",
			want:  "<script>alert(1)</script>\nplain",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderDescriptionHTML(tt.input)
			if got != tt.want {
				t.Fatalf("renderDescriptionHTML() = %q, want %q", got, tt.want)
			}
		})
	}
}
