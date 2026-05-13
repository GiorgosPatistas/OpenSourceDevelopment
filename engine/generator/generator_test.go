package generator

import (
	"testing"

	"github.com/georgepatistas/my-ssg/parser"
)

func TestFilterDrafts(t *testing.T) {
	t.Run("mixed: keeps only published pages", func(t *testing.T) {
		pages := []*parser.Page{
			{Title: "Published Post", Draft: false},
			{Title: "Draft Post", Draft: true},
			{Title: "Another Published", Draft: false},
		}
		result := filterDrafts(pages)
		if len(result) != 2 {
			t.Errorf("expected 2 published pages, got %d", len(result))
		}
		for _, p := range result {
			if p.Draft {
				t.Errorf("draft page '%s' should not be in the result", p.Title)
			}
		}
	})

	t.Run("all drafts → empty slice", func(t *testing.T) {
		pages := []*parser.Page{
			{Title: "Draft 1", Draft: true},
			{Title: "Draft 2", Draft: true},
		}
		result := filterDrafts(pages)
		if len(result) != 0 {
			t.Errorf("expected 0 pages, got %d", len(result))
		}
	})

	t.Run("no drafts → returns all pages", func(t *testing.T) {
		pages := []*parser.Page{
			{Title: "Post 1", Draft: false},
			{Title: "Post 2", Draft: false},
			{Title: "Post 3", Draft: false},
		}
		result := filterDrafts(pages)
		if len(result) != 3 {
			t.Errorf("expected 3 pages, got %d", len(result))
		}
	})

	t.Run("empty slice → nil or empty", func(t *testing.T) {
		result := filterDrafts([]*parser.Page{})
		if len(result) != 0 {
			t.Errorf("expected 0 pages for empty input, got %d", len(result))
		}
	})
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig("/path/to/site")

	tests := []struct {
		field    string
		got      string
		expected string
	}{
		{"SiteDir", cfg.SiteDir, "/path/to/site"},
		{"ContentDir", cfg.ContentDir, "content"},
		{"TemplatesDir", cfg.TemplatesDir, "templates"},
		{"StaticDir", cfg.StaticDir, "static"},
		{"OutputDir", cfg.OutputDir, "dist"},
		{"SiteTitle", cfg.SiteTitle, "My Site"},
		{"SiteURL", cfg.SiteURL, "/"},
	}

	for _, tt := range tests {
		if tt.got != tt.expected {
			t.Errorf("DefaultConfig.%s = %q, want %q", tt.field, tt.got, tt.expected)
		}
	}
}
