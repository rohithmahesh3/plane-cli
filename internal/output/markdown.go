package output

import (
	"strings"

	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/strikethrough"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/table"
)

func htmlToMarkdown(input string) (string, error) {
	if strings.TrimSpace(input) == "" {
		return "", nil
	}

	conv := converter.NewConverter(
		converter.WithPlugins(
			base.NewBasePlugin(),
			commonmark.NewCommonmarkPlugin(),
			table.NewTablePlugin(
				table.WithNewlineBehavior(table.NewlineBehaviorPreserve),
			),
			strikethrough.NewStrikethroughPlugin(),
		),
	)

	markdown, err := conv.ConvertString(input)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(markdown), nil
}

func convertHTMLValue(input string) string {
	markdown, err := htmlToMarkdown(input)
	if err == nil {
		return markdown
	}

	return input
}
