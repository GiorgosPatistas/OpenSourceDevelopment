package parser

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// mdConverter is the goldmark instance with all extensions enabled.
var mdConverter = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,            // GitHub Flavored Markdown (tables, strikethrough, task lists, autolink)
		extension.Footnote,       // Footnotes
		extension.DefinitionList, // Definition lists
		extension.Typographer,    // Smart quotes, em-dashes, etc.
	),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(), // Auto-generate IDs on headings for anchor links
	),
	goldmark.WithRendererOptions(
		html.WithHardWraps(), // Hard line breaks
		html.WithXHTML(),     // XHTML output
		html.WithUnsafe(),    // Allow raw HTML inside markdown
	),
)

// MarkdownToHTML converts a markdown string to an HTML string.
func MarkdownToHTML(markdown string) (string, error) {
	var buf bytes.Buffer
	if err := mdConverter.Convert([]byte(markdown), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
