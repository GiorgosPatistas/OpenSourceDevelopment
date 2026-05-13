package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/georgepatistas/my-ssg/generator"
)

// siteConfig mirrors the optional ssg.config.json file
type siteConfig struct {
	SiteTitle   string `json:"siteTitle"`
	SiteURL     string `json:"siteUrl"`
	Description string `json:"description"`
	ContentDir  string `json:"contentDir"`
	TemplateDir string `json:"templateDir"`
	StaticDir   string `json:"staticDir"`
	OutputDir   string `json:"outputDir"`
}

func main() {
	// The TypeScript CLI passes the site directory as the first argument
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "❌ Usage: engine <path/to/site>")
		os.Exit(1)
	}

	siteDir := os.Args[1]

	// Check if the directory exists
	if _, err := os.Stat(siteDir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "❌ Directory not found: %s\n", siteDir)
		os.Exit(1)
	}

	// Load default config
	cfg := generator.DefaultConfig(siteDir)

	// Try to load ssg.config.json if it exists
	configPath := filepath.Join(siteDir, "ssg.config.json")
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		if sc, err := loadSiteConfig(configPath); err == nil {
			applyConfig(&cfg, sc)
			fmt.Printf("⚙️  Loading config from: %s\n", configPath)
		}
	}

	fmt.Printf("🚀 Starting build for: %s\n", siteDir)
	fmt.Println("──────────────────────────────────────────")

	gen := generator.New(cfg)
	if err := gen.Build(); err != nil {
		fmt.Fprintf(os.Stderr, "\n❌ Build failed: %v\n", err)
		os.Exit(1)
	}
}

func loadSiteConfig(path string) (*siteConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg siteConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func applyConfig(cfg *generator.Config, sc *siteConfig) {
	if sc.SiteTitle != "" {
		cfg.SiteTitle = sc.SiteTitle
	}
	if sc.SiteURL != "" {
		cfg.SiteURL = sc.SiteURL
	}
	if sc.Description != "" {
		cfg.Description = sc.Description
	}
	if sc.ContentDir != "" {
		cfg.ContentDir = sc.ContentDir
	}
	if sc.TemplateDir != "" {
		cfg.TemplatesDir = sc.TemplateDir
	}
	if sc.StaticDir != "" {
		cfg.StaticDir = sc.StaticDir
	}
	if sc.OutputDir != "" {
		cfg.OutputDir = sc.OutputDir
	}
}
