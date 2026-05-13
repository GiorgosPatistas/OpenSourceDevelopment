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

// SiteData holds global site information available to all templates.
type SiteData struct {
	SiteTitle   string
	SiteURL     string
	Description string
	BuildDate   time.Time
}

// PageData holds the data passed to the template for each page.
type PageData struct {
	Site    SiteData
	Page    *parser.Page
	Content template.HTML  // HTML content (marked safe to prevent escaping)
	Pages   []*parser.Page // all pages (for navigation/index)
}

// Renderer loads and executes HTML templates.
type Renderer struct {
	templateDir string
	siteData    SiteData
	funcMap     template.FuncMap
	layoutPath  string
	pagePath    string
	indexPath   string
}

// New creates a Renderer, verifying that the required templates exist.
func New(templateDir string, siteData SiteData) (*Renderer, error) {
	layoutPath := filepath.Join(templateDir, "layout.html")
	if _, err := os.Stat(layoutPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("templates/layout.html not found")
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

	// Register optional template paths
	pagePath := filepath.Join(templateDir, "page.html")
	if _, err := os.Stat(pagePath); !os.IsNotExist(err) {
		r.pagePath = pagePath
	}

	indexPath := filepath.Join(templateDir, "index.html")
	if _, err := os.Stat(indexPath); !os.IsNotExist(err) {
		r.indexPath = indexPath
	}

	// Validation: do a trial parse of the layout to catch syntax errors early
	if err := r.validateLayout(); err != nil {
		return nil, err
	}

	return r, nil
}

// validateLayout does a trial parse to catch syntax errors.
func (r *Renderer) validateLayout() error {
	files := []string{r.layoutPath}
	if r.pagePath != "" {
		files = append(files, r.pagePath)
	}
	_, err := template.New("layout.html").Funcs(r.funcMap).ParseFiles(files...)
	if err != nil {
		return fmt.Errorf("error parsing templates: %w", err)
	}
	return nil
}

// RenderPage applies the layout + page template to a single page.
// Result: a complete HTML file.
func (r *Renderer) RenderPage(page *parser.Page, allPages []*parser.Page) (string, error) {
	// Load a fresh template set each time (layout.html + page.html)
	files := []string{r.layoutPath}
	if r.pagePath != "" {
		files = append(files, r.pagePath)
	}

	tmpl, err := template.New("layout.html").Funcs(r.funcMap).ParseFiles(files...)
	if err != nil {
		return "", fmt.Errorf("error parsing page templates: %w", err)
	}

	data := PageData{
		Site:    r.siteData,
		Page:    page,
		Content: template.HTML(page.HTMLContent),
		Pages:   allPages,
	}

	var buf bytes.Buffer
	// Execute "layout.html" which calls {{template "content" .}}
	// defined in page.html
	if err := tmpl.ExecuteTemplate(&buf, "layout.html", data); err != nil {
		return "", fmt.Errorf("error rendering page '%s': %w", page.Title, err)
	}
	return buf.String(), nil
}

// RenderIndex generates the index.html.
// If templates/index.html exists, it is used as the "content" block inside the layout.
// Otherwise, a simple list of page links is generated automatically.
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
		// Load layout.html + index.html together into the same template set.
		// index.html defines {{define "content"}} called by the layout.
		tmpl, err = template.New("layout.html").Funcs(r.funcMap).ParseFiles(r.layoutPath, r.indexPath)
		if err != nil {
			return "", fmt.Errorf("error parsing index templates: %w", err)
		}
	} else {
		// Fallback: auto-generate content as a list of pages
		tmpl, err = template.New("layout.html").Funcs(r.funcMap).ParseFiles(r.layoutPath)
		if err != nil {
			return "", fmt.Errorf("error parsing layout for index: %w", err)
		}
		// Define the "content" template inline
		defaultContent := `{{define "content"}}` + buildDefaultIndexContent(allPages) + `{{end}}`
		tmpl, err = tmpl.Parse(defaultContent)
		if err != nil {
			return "", fmt.Errorf("error parsing default index content: %w", err)
		}
		data.Content = ""
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "layout.html", data); err != nil {
		return "", fmt.Errorf("error rendering index: %w", err)
	}
	return buf.String(), nil
}

// buildDefaultIndexContent generates a simple HTML list linking to all pages.
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
