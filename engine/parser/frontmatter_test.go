package parser

import (
	"errors"
	"testing"
)

func TestParse(t *testing.T) {
	t.Run("valid front matter with all fields", func(t *testing.T) {
		input := []byte(`---
title: "Hello World"
date: 2024-01-15
description: "A test page"
draft: false
---
# Content here
`)
		page, err := Parse(input)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if page.Title != "Hello World" {
			t.Errorf("expected title 'Hello World', got %q", page.Title)
		}
		if page.Description != "A test page" {
			t.Errorf("expected description 'A test page', got %q", page.Description)
		}
		if page.Draft != false {
			t.Error("expected Draft=false")
		}
		if page.RawMarkdown != "# Content here" {
			t.Errorf("unexpected RawMarkdown: %q", page.RawMarkdown)
		}
	})

	t.Run("missing front matter returns ErrNoFrontMatter", func(t *testing.T) {
		input := []byte("# Just markdown, no front matter")
		_, err := Parse(input)
		if !errors.Is(err, ErrNoFrontMatter) {
			t.Errorf("expected ErrNoFrontMatter, got: %v", err)
		}
	})

	t.Run("invalid YAML returns ErrInvalidFrontMatter", func(t *testing.T) {
		input := []byte("---\ntitle: [invalid yaml\n---\n")
		_, err := Parse(input)
		if !errors.Is(err, ErrInvalidFrontMatter) {
			t.Errorf("expected ErrInvalidFrontMatter, got: %v", err)
		}
	})

	t.Run("draft: true", func(t *testing.T) {
		input := []byte("---\ntitle: Draft Post\ndraft: true\n---\nContent\n")
		page, err := Parse(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !page.Draft {
			t.Error("expected Draft=true")
		}
	})

	t.Run("custom slug", func(t *testing.T) {
		input := []byte("---\ntitle: My Page\nslug: custom-slug\n---\nContent\n")
		page, err := Parse(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if page.Slug != "custom-slug" {
			t.Errorf("expected slug 'custom-slug', got %q", page.Slug)
		}
	})

	t.Run("empty body after front matter", func(t *testing.T) {
		input := []byte("---\ntitle: Empty Body\n---\n")
		page, err := Parse(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if page.RawMarkdown != "" {
			t.Errorf("expected empty RawMarkdown, got %q", page.RawMarkdown)
		}
	})

	t.Run("completely empty input", func(t *testing.T) {
		_, err := Parse([]byte(""))
		if !errors.Is(err, ErrNoFrontMatter) {
			t.Errorf("expected ErrNoFrontMatter for empty input, got: %v", err)
		}
	})

	t.Run("multiline markdown body", func(t *testing.T) {
		input := []byte("---\ntitle: Multi\n---\n# H1\n\nParagraph text.\n\n- item one\n- item two\n")
		page, err := Parse(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if page.Title != "Multi" {
			t.Errorf("unexpected title: %q", page.Title)
		}
		if page.RawMarkdown == "" {
			t.Error("expected non-empty RawMarkdown")
		}
	})
}
