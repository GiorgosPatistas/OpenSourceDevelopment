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

// Config holds the Generator's configuration.
type Config struct {
	// Root directory of the site (passed by the TypeScript CLI)
	SiteDir string

	// Subdirectories (relative to SiteDir)
	ContentDir   string // default: "content"
	TemplatesDir string // default: "templates"
	StaticDir    string // default: "static"
	OutputDir    string // default: "dist"

	// Site metadata
	SiteTitle   string
	SiteURL     string
	Description string
}

// DefaultConfig returns a Config with default paths.
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

// Generator is the central orchestrator.
type Generator struct {
	cfg Config
}

// New creates a new Generator.
func New(cfg Config) *Generator {
	return &Generator{cfg: cfg}
}

// Build runs the full site generation pipeline.
func (g *Generator) Build() error {
	cfg := g.cfg

	contentDir := filepath.Join(cfg.SiteDir, cfg.ContentDir)
	templateDir := filepath.Join(cfg.SiteDir, cfg.TemplatesDir)
	staticDir := filepath.Join(cfg.SiteDir, cfg.StaticDir)
	outputDir := filepath.Join(cfg.SiteDir, cfg.OutputDir)

	// 1. Clean and create the output directory
	fmt.Printf("📁 Creating output directory: %s\n", outputDir)
	if err := os.RemoveAll(outputDir); err != nil {
		return fmt.Errorf("error cleaning output directory: %w", err)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("error creating output directory: %w", err)
	}

	// 2. Parse all markdown files
	fmt.Printf("📖 Parsing markdown from: %s\n", contentDir)
	pages, err := g.parseContent(contentDir)
	if err != nil {
		return err
	}
	fmt.Printf("   ✓ Found %d pages\n", len(pages))

	// Filter out drafts
	publishedPages := filterDrafts(pages)
	fmt.Printf("   ✓ %d pages to publish (drafts excluded)\n", len(publishedPages))

	// Sort by date (newest first)
	sort.Slice(publishedPages, func(i, j int) bool {
		return publishedPages[i].Date.After(publishedPages[j].Date)
	})

	// 3. Create the Renderer
	siteData := renderer.SiteData{
		SiteTitle:   cfg.SiteTitle,
		SiteURL:     cfg.SiteURL,
		Description: cfg.Description,
	}
	r, err := renderer.New(templateDir, siteData)
	if err != nil {
		return fmt.Errorf("error loading templates: %w", err)
	}

	// 4. Render each page
	fmt.Println("🔨 Rendering pages...")
	for _, page := range publishedPages {
		html, err := r.RenderPage(page, publishedPages)
		if err != nil {
			fmt.Printf("   ⚠️  Skipping '%s': %v\n", page.OutputFile, err)
			continue
		}

		outputPath := filepath.Join(outputDir, page.OutputFile)
		// Create subdirectory if needed
		if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
			return fmt.Errorf("error creating directory for '%s': %w", page.OutputFile, err)
		}
		if err := os.WriteFile(outputPath, []byte(html), 0644); err != nil {
			return fmt.Errorf("error writing '%s': %w", page.OutputFile, err)
		}
		fmt.Printf("   ✓ %s → %s\n", page.Title, page.OutputFile)
	}

	// 5. Render index.html
	fmt.Println("🏠 Rendering index.html...")
	indexHTML, err := r.RenderIndex(publishedPages)
	if err != nil {
		return fmt.Errorf("error rendering index: %w", err)
	}
	if err := os.WriteFile(filepath.Join(outputDir, "index.html"), []byte(indexHTML), 0644); err != nil {
		return fmt.Errorf("error writing index.html: %w", err)
	}
	fmt.Println("   ✓ index.html")

	// 6. Copy static assets
	if _, err := os.Stat(staticDir); !os.IsNotExist(err) {
		fmt.Printf("📦 Copying static assets from: %s\n", staticDir)
		if err := copyDir(staticDir, outputDir); err != nil {
			return fmt.Errorf("error copying static assets: %w", err)
		}
		fmt.Println("   ✓ Static assets copied")
	}

	fmt.Printf("\n✅ Site successfully generated at: %s\n", outputDir)
	return nil
}

// parseContent reads all .md files from contentDir recursively.
func (g *Generator) parseContent(contentDir string) ([]*parser.Page, error) {
	if _, err := os.Stat(contentDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("content directory not found: %s", contentDir)
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
			return fmt.Errorf("error reading '%s': %w", path, err)
		}

		page, err := parser.Parse(content)
		if err != nil {
			fmt.Printf("   ⚠️  Skipping '%s': %v\n", path, err)
			return nil
		}

		// Convert markdown to HTML
		htmlContent, err := parser.MarkdownToHTML(page.RawMarkdown)
		if err != nil {
			return fmt.Errorf("error converting markdown '%s': %w", path, err)
		}
		page.HTMLContent = htmlContent

		// Compute output path
		relPath, _ := filepath.Rel(contentDir, path)
		page.OutputFile = strings.TrimSuffix(relPath, ".md") + ".html"
		page.OutputFile = filepath.ToSlash(page.OutputFile) // Windows-safe

		// Use custom slug if provided
		if page.Slug != "" {
			dir := filepath.Dir(page.OutputFile)
			if dir == "." {
				page.OutputFile = page.Slug + ".html"
			} else {
				page.OutputFile = dir + "/" + page.Slug + ".html"
			}
		}

		page.URL = "/" + page.OutputFile

		// Fall back to filename as title if none provided
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

// filterDrafts returns only pages that are not marked as drafts.
func filterDrafts(pages []*parser.Page) []*parser.Page {
	var result []*parser.Page
	for _, p := range pages {
		if !p.Draft {
			result = append(result, p)
		}
	}
	return result
}

// copyDir recursively copies the contents of src into dst.
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

// copyFile copies a single file from src to dst.
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
