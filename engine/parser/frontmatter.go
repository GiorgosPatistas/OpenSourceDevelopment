package parser

import (
	"bytes"
	"errors"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Page represents a parsed markdown file with its metadata and HTML content.
type Page struct {
	// Metadata from the YAML front matter
	Title       string    `yaml:"title"`
	Date        time.Time `yaml:"date"`
	Description string    `yaml:"description"`
	Draft       bool      `yaml:"draft"`
	Slug        string    `yaml:"slug"` // optional custom URL slug

	// Populated after parsing
	RawMarkdown string // the markdown content without the front matter
	HTMLContent string // the final HTML after conversion
	OutputFile  string // the output HTML filename (e.g. "about.html")
	URL         string // the relative URL (e.g. "/about.html")
}

var (
	ErrNoFrontMatter      = errors.New("no YAML front matter found (--- ... ---)")
	ErrInvalidFrontMatter = errors.New("invalid YAML front matter")
)

// Parse reads a .md file, extracts the YAML front matter, and returns a Page.
// The front matter must begin with "---" on the first line.
func Parse(content []byte) (*Page, error) {
	text := string(content)

	// Check if the file starts with ---
	if !strings.HasPrefix(strings.TrimLeft(text, "\r\n"), "---") {
		return nil, ErrNoFrontMatter
	}

	// Find the closing ---
	parts := splitFrontMatter(text)
	if parts == nil {
		return nil, ErrInvalidFrontMatter
	}

	yamlStr := parts[0]
	markdownStr := parts[1]

	page := &Page{}
	decoder := yaml.NewDecoder(bytes.NewBufferString(yamlStr))
	if err := decoder.Decode(page); err != nil {
		return nil, ErrInvalidFrontMatter
	}

	page.RawMarkdown = strings.TrimSpace(markdownStr)
	return page, nil
}

// splitFrontMatter splits the content into [yamlPart, markdownPart].
// Returns nil if no valid front matter is found.
func splitFrontMatter(text string) []string {
	// Strip leading newlines if present
	text = strings.TrimLeft(text, "\r\n")

	// Remove the opening ---
	if !strings.HasPrefix(text, "---") {
		return nil
	}
	text = text[3:]

	// Advance to the next line
	idx := strings.Index(text, "\n")
	if idx == -1 {
		return nil
	}
	text = text[idx+1:]

	// Find the closing ---
	closeIdx := strings.Index(text, "\n---")
	if closeIdx == -1 {
		// Check if it closes immediately with ---
		if strings.HasPrefix(text, "---") {
			return []string{"", text[3:]}
		}
		return nil
	}

	yamlPart := text[:closeIdx]
	rest := text[closeIdx+4:] // +4 to skip "\n---"

	// Skip any newline after the closing ---
	rest = strings.TrimLeft(rest, "\r\n")

	return []string{yamlPart, rest}
}
