package generator

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/georgepatistas/my-ssg/parser"
	"github.com/georgepatistas/my-ssg/renderer"
)

// Config ρυθμίζει τον Generator.
type Config struct {
	// Ο root φάκελος του site (αυτός που περνάει η TypeScript CLI)
	SiteDir string

	// Υποφάκελοι (σχετικά με τον SiteDir)
	ContentDir   string // default: "content"
	TemplatesDir string // default: "templates"
	StaticDir    string // default: "static"
	OutputDir    string // default: "dist"

	// Δεδομένα για το site
	SiteTitle   string
	SiteURL     string
	Description string
}

// DefaultConfig επιστρέφει έναν Config με τα default paths.
func DefaultConfig(siteDir string) Config {
	return Config{
		SiteDir:      siteDir,
		ContentDir:   "content",
		TemplatesDir: "templates",
		StaticDir:    "static",
		OutputDir:    "dist",
		SiteTitle:    "My Site",
		SiteURL:      "/",
		Description:  "",
	}
}

// Generator είναι ο κεντρικός orchestrator.
type Generator struct {
	cfg Config
}

// New δημιουργεί έναν Generator.
func New(cfg Config) *Generator {
	return &Generator{cfg: cfg}
}

// Build εκτελεί ολόκληρη τη διαδικασία παραγωγής του site.
func (g *Generator) Build() error {
	cfg := g.cfg

	contentDir := filepath.Join(cfg.SiteDir, cfg.ContentDir)
	templateDir := filepath.Join(cfg.SiteDir, cfg.TemplatesDir)
	staticDir := filepath.Join(cfg.SiteDir, cfg.StaticDir)
	outputDir := filepath.Join(cfg.SiteDir, cfg.OutputDir)

	// 1. Καθαρισμός και δημιουργία output dir
	fmt.Printf("📁 Δημιουργία output dir: %s\n", outputDir)
	if err := os.RemoveAll(outputDir); err != nil {
		return fmt.Errorf("σφάλμα καθαρισμού output dir: %w", err)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("σφάλμα δημιουργίας output dir: %w", err)
	}

	// 2. Parsing όλων των markdown αρχείων
	fmt.Printf("📖 Parsing markdown από: %s\n", contentDir)
	pages, err := g.parseContent(contentDir)
	if err != nil {
		return err
	}
	fmt.Printf("   ✓ Βρέθηκαν %d σελίδες\n", len(pages))

	// Φιλτράρουμε drafts
	publishedPages := filterDrafts(pages)
	fmt.Printf("   ✓ %d σελίδες για δημοσίευση (χωρίς drafts)\n", len(publishedPages))

	// Ταξινόμηση κατά ημερομηνία (νεότερες πρώτα)
	sort.Slice(publishedPages, func(i, j int) bool {
		return publishedPages[i].Date.After(publishedPages[j].Date)
	})

	// 3. Δημιουργία του Renderer
	siteData := renderer.SiteData{
		SiteTitle:   cfg.SiteTitle,
		SiteURL:     cfg.SiteURL,
		Description: cfg.Description,
	}
	r, err := renderer.New(templateDir, siteData)
	if err != nil {
		return fmt.Errorf("σφάλμα φόρτωσης templates: %w", err)
	}

	// 4. Render κάθε σελίδα
	fmt.Println("🔨 Rendering σελίδων...")
	for _, page := range publishedPages {
		html, err := r.RenderPage(page, publishedPages)
		if err != nil {
			fmt.Printf("   ⚠️  Παράλειψη '%s': %v\n", page.OutputFile, err)
			continue
		}

		outputPath := filepath.Join(outputDir, page.OutputFile)
		// Δημιουργία subdirectory αν χρειάζεται
		if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
			return fmt.Errorf("σφάλμα δημιουργίας dir για '%s': %w", page.OutputFile, err)
		}
		if err := os.WriteFile(outputPath, []byte(html), 0644); err != nil {
			return fmt.Errorf("σφάλμα αποθήκευσης '%s': %w", page.OutputFile, err)
		}
		fmt.Printf("   ✓ %s → %s\n", page.Title, page.OutputFile)
	}

	// 5. Render index.html
	fmt.Println("🏠 Rendering index.html...")
	indexHTML, err := r.RenderIndex(publishedPages)
	if err != nil {
		return fmt.Errorf("σφάλμα rendering index: %w", err)
	}
	if err := os.WriteFile(filepath.Join(outputDir, "index.html"), []byte(indexHTML), 0644); err != nil {
		return fmt.Errorf("σφάλμα αποθήκευσης index.html: %w", err)
	}
	fmt.Println("   ✓ index.html")

	// 6. Αντιγραφή static assets
	if _, err := os.Stat(staticDir); !os.IsNotExist(err) {
		fmt.Printf("📦 Αντιγραφή static assets από: %s\n", staticDir)
		if err := copyDir(staticDir, outputDir); err != nil {
			return fmt.Errorf("σφάλμα αντιγραφής static assets: %w", err)
		}
		fmt.Println("   ✓ Static assets αντιγράφηκαν")
	}

	fmt.Printf("\n✅ Site δημιουργήθηκε επιτυχώς στο: %s\n", outputDir)
	return nil
}

// parseContent διαβάζει όλα τα .md αρχεία από τον contentDir (αναδρομικά).
func (g *Generator) parseContent(contentDir string) ([]*parser.Page, error) {
	if _, err := os.Stat(contentDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("δεν βρέθηκε φάκελος content: %s", contentDir)
	}

	var pages []*parser.Page

	err := filepath.WalkDir(contentDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("σφάλμα ανάγνωσης '%s': %w", path, err)
		}

		page, err := parser.Parse(content)
		if err != nil {
			fmt.Printf("   ⚠️  Παράλειψη '%s': %v\n", path, err)
			return nil
		}

		// Μετατροπή markdown → HTML
		htmlContent, err := parser.MarkdownToHTML(page.RawMarkdown)
		if err != nil {
			return fmt.Errorf("σφάλμα μετατροπής markdown '%s': %w", path, err)
		}
		page.HTMLContent = htmlContent

		// Υπολογισμός output path
		relPath, _ := filepath.Rel(contentDir, path)
		page.OutputFile = strings.TrimSuffix(relPath, ".md") + ".html"
		page.OutputFile = filepath.ToSlash(page.OutputFile) // Windows-safe

		// Χρήση custom slug αν υπάρχει
		if page.Slug != "" {
			dir := filepath.Dir(page.OutputFile)
			if dir == "." {
				page.OutputFile = page.Slug + ".html"
			} else {
				page.OutputFile = dir + "/" + page.Slug + ".html"
			}
		}

		page.URL = "/" + page.OutputFile

		// Αν δεν έχει τίτλο, παίρνουμε από το filename
		if page.Title == "" {
			base := strings.TrimSuffix(filepath.Base(path), ".md")
			page.Title = strings.ReplaceAll(base, "-", " ")
			page.Title = strings.ReplaceAll(page.Title, "_", " ")
		}

		pages = append(pages, page)
		return nil
	})

	return pages, err
}

// filterDrafts επιστρέφει μόνο τις σελίδες που δεν είναι draft.
func filterDrafts(pages []*parser.Page) []*parser.Page {
	var result []*parser.Page
	for _, p := range pages {
		if !p.Draft {
			result = append(result, p)
		}
	}
	return result
}

// copyDir αντιγράφει αναδρομικά τα περιεχόμενα του src στο dst.
func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(src, path)
		targetPath := filepath.Join(dst, relPath)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		return copyFile(path, targetPath)
	})
}

// copyFile αντιγράφει ένα αρχείο.
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
