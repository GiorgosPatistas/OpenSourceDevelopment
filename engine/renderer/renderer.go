package renderer

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/georgepatistas/my-ssg/parser"
)

// SiteData περιέχει global πληροφορίες για το site που είναι διαθέσιμες σε όλα τα templates.
type SiteData struct {
	SiteTitle   string
	SiteURL     string
	Description string
	BuildDate   time.Time
}

// PageData είναι τα δεδομένα που περνάμε στο template για κάθε σελίδα.
type PageData struct {
	Site    SiteData
	Page    *parser.Page
	Content template.HTML // το HTML content (marked safe για να μην γίνει escape)
	Pages   []*parser.Page // όλες οι σελίδες (για navigation/index)
}

// Renderer φορτώνει και εκτελεί τα HTML templates.
type Renderer struct {
	templateDir string
	siteData    SiteData
	funcMap     template.FuncMap
	layoutPath  string
	pagePath    string
	indexPath   string
}

// New δημιουργεί έναν Renderer ελέγχοντας ότι υπάρχουν τα απαραίτητα templates.
func New(templateDir string, siteData SiteData) (*Renderer, error) {
	layoutPath := filepath.Join(templateDir, "layout.html")
	if _, err := os.Stat(layoutPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("δεν βρέθηκε το templates/layout.html")
	}

	funcMap := template.FuncMap{
		"formatDate": func(t time.Time) string {
			if t.IsZero() {
				return ""
			}
			return t.Format("02 January 2006")
		},
		"formatDateISO": func(t time.Time) string {
			if t.IsZero() {
				return ""
			}
			return t.Format("2006-01-02")
		},
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
	}

	r := &Renderer{
		templateDir: templateDir,
		siteData:    siteData,
		funcMap:     funcMap,
		layoutPath:  layoutPath,
	}

	// Καταχώρηση προαιρετικών template paths
	pagePath := filepath.Join(templateDir, "page.html")
	if _, err := os.Stat(pagePath); !os.IsNotExist(err) {
		r.pagePath = pagePath
	}

	indexPath := filepath.Join(templateDir, "index.html")
	if _, err := os.Stat(indexPath); !os.IsNotExist(err) {
		r.indexPath = indexPath
	}

	// Validation: δοκιμαστικό parse του layout για να πιάσουμε syntax errors νωρίς
	if err := r.validateLayout(); err != nil {
		return nil, err
	}

	return r, nil
}

// validateLayout κάνει ένα δοκιμαστικό parse για να πιάσει syntax errors.
func (r *Renderer) validateLayout() error {
	files := []string{r.layoutPath}
	if r.pagePath != "" {
		files = append(files, r.pagePath)
	}
	_, err := template.New("layout.html").Funcs(r.funcMap).ParseFiles(files...)
	if err != nil {
		return fmt.Errorf("σφάλμα parsing templates: %w", err)
	}
	return nil
}

// RenderPage εφαρμόζει το layout + page template σε μια σελίδα.
// Αποτέλεσμα: πλήρες HTML αρχείο.
func (r *Renderer) RenderPage(page *parser.Page, allPages []*parser.Page) (string, error) {
	// Κάθε φορά φορτώνουμε φρέσκο template set (layout.html + page.html)
	files := []string{r.layoutPath}
	if r.pagePath != "" {
		files = append(files, r.pagePath)
	}

	tmpl, err := template.New("layout.html").Funcs(r.funcMap).ParseFiles(files...)
	if err != nil {
		return "", fmt.Errorf("σφάλμα parsing page templates: %w", err)
	}

	data := PageData{
		Site:    r.siteData,
		Page:    page,
		Content: template.HTML(page.HTMLContent),
		Pages:   allPages,
	}

	var buf bytes.Buffer
	// Εκτελούμε το "layout.html" που καλεί {{template "content" .}}
	// το οποίο είναι defined στο page.html
	if err := tmpl.ExecuteTemplate(&buf, "layout.html", data); err != nil {
		return "", fmt.Errorf("σφάλμα rendering σελίδας '%s': %w", page.Title, err)
	}
	return buf.String(), nil
}

// RenderIndex δημιουργεί το index.html.
// Αν υπάρχει templates/index.html, το χρησιμοποιεί ως "content" block μέσα στο layout.
// Αλλιώς, δημιουργεί αυτόματα μια λίστα με links.
func (r *Renderer) RenderIndex(allPages []*parser.Page) (string, error) {
	data := PageData{
		Site:  r.siteData,
		Pages: allPages,
		Page: &parser.Page{
			Title:       r.siteData.SiteTitle,
			Description: r.siteData.Description,
		},
	}

	var tmpl *template.Template
	var err error

	if r.indexPath != "" {
		// Φορτώνουμε layout.html + index.html μαζί στο ίδιο template set
		// Το index.html ορίζει {{define "content"}} που καλείται από το layout
		tmpl, err = template.New("layout.html").Funcs(r.funcMap).ParseFiles(r.layoutPath, r.indexPath)
		if err != nil {
			return "", fmt.Errorf("σφάλμα parsing index templates: %w", err)
		}
	} else {
		// Fallback: φτιάχνουμε αυτόματα το content με λίστα σελίδων
		tmpl, err = template.New("layout.html").Funcs(r.funcMap).ParseFiles(r.layoutPath)
		if err != nil {
			return "", fmt.Errorf("σφάλμα parsing layout για index: %w", err)
		}
		// Ορίζουμε inline το "content" template
		defaultContent := `{{define "content"}}` + buildDefaultIndexContent(allPages) + `{{end}}`
		tmpl, err = tmpl.Parse(defaultContent)
		if err != nil {
			return "", fmt.Errorf("σφάλμα parsing default index content: %w", err)
		}
		data.Content = ""
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "layout.html", data); err != nil {
		return "", fmt.Errorf("σφάλμα rendering index: %w", err)
	}
	return buf.String(), nil
}

// buildDefaultIndexContent δημιουργεί ένα απλό HTML list με links σε όλες τις σελίδες.
func buildDefaultIndexContent(pages []*parser.Page) string {
	var b bytes.Buffer
	b.WriteString("<ul class=\"page-list\">\n")
	for _, p := range pages {
		dateStr := ""
		if !p.Date.IsZero() {
			dateStr = fmt.Sprintf(` <time datetime="%s">%s</time>`,
				p.Date.Format("2006-01-02"),
				p.Date.Format("02 Jan 2006"))
		}
		b.WriteString(fmt.Sprintf(
			"  <li><a href=\"%s\">%s</a>%s</li>\n",
			p.URL, p.Title, dateStr,
		))
	}
	b.WriteString("</ul>")
	return b.String()
}
