package parser

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// mdConverter είναι ο goldmark instance με όλα τα extensions ενεργά.
var mdConverter = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,           // GitHub Flavored Markdown (tables, strikethrough, task lists, autolink)
		extension.Footnote,      // Υποσημειώσεις
		extension.DefinitionList, // Definition lists
		extension.Typographer,   // Έξυπνα εισαγωγικά, em-dashes κλπ.
	),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(), // Αυτόματα IDs σε headings για anchor links
	),
	goldmark.WithRendererOptions(
		html.WithHardWraps(), // Σκληρά line breaks
		html.WithXHTML(),     // XHTML output
		html.WithUnsafe(),    // Επιτρέπει raw HTML μέσα στο markdown
	),
)

// MarkdownToHTML μετατρέπει markdown σε HTML string.
func MarkdownToHTML(markdown string) (string, error) {
	var buf bytes.Buffer
	if err := mdConverter.Convert([]byte(markdown), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
