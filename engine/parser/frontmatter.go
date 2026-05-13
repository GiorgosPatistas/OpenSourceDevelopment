package parser

import (
	"bytes"
	"errors"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Page αντιπροσωπεύει ένα parsed markdown αρχείο με metadata και HTML content.
type Page struct {
	// Metadata από το YAML front matter
	Title       string    `yaml:"title"`
	Date        time.Time `yaml:"date"`
	Description string    `yaml:"description"`
	Draft       bool      `yaml:"draft"`
	Slug        string    `yaml:"slug"` // προαιρετικό custom URL slug

	// Γεμίζουν μετά το parsing
	RawMarkdown string // το markdown χωρίς το front matter
	HTMLContent string // το τελικό HTML μετά τη μετατροπή
	OutputFile  string // το όνομα του output HTML αρχείου (π.χ. "about.html")
	URL         string // το relative URL (π.χ. "/about.html")
}

var (
	ErrNoFrontMatter    = errors.New("δεν βρέθηκε YAML front matter (--- ... ---)")
	ErrInvalidFrontMatter = errors.New("μη έγκυρο YAML front matter")
)

// Parse διαβάζει ένα .md αρχείο, εξάγει το YAML front matter και επιστρέφει ένα Page.
// Το front matter πρέπει να αρχίζει με "---" στην πρώτη γραμμή.
func Parse(content []byte) (*Page, error) {
	text := string(content)

	// Έλεγχος αν αρχίζει με ---
	if !strings.HasPrefix(strings.TrimLeft(text, "\r\n"), "---") {
		return nil, ErrNoFrontMatter
	}

	// Βρίσκουμε το κλείσιμο ---
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

// splitFrontMatter χωρίζει το content σε [yamlPart, markdownPart].
// Επιστρέφει nil αν δεν βρεθεί σωστό front matter.
func splitFrontMatter(text string) []string {
	// Αφαίρεση leading newline αν υπάρχει
	text = strings.TrimLeft(text, "\r\n")

	// Αφαίρεση του πρώτου ---
	if !strings.HasPrefix(text, "---") {
		return nil
	}
	text = text[3:]

	// Πάμε στην επόμενη γραμμή
	idx := strings.Index(text, "\n")
	if idx == -1 {
		return nil
	}
	text = text[idx+1:]

	// Βρίσκουμε το κλείσιμο ---
	closeIdx := strings.Index(text, "\n---")
	if closeIdx == -1 {
		// Δοκιμάζουμε αν τελειώνει αμέσως με ---
		if strings.HasPrefix(text, "---") {
			return []string{"", text[3:]}
		}
		return nil
	}

	yamlPart := text[:closeIdx]
	rest := text[closeIdx+4:] // +4 για να skip-άρουμε το "\n---"

	// Skip τυχόν newline μετά το κλείσιμο ---
	rest = strings.TrimLeft(rest, "\r\n")

	return []string{yamlPart, rest}
}
