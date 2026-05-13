package parser

import (
	"strings"
	"testing"
)

func TestMarkdownToHTML(t *testing.T) {
	t.Run("heading παράγει <h1>", func(t *testing.T) {
		html, err := MarkdownToHTML("# Hello World")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(html, "<h1") {
			t.Errorf("expected <h1> tag in output, got: %s", html)
		}
		if !strings.Contains(html, "Hello World") {
			t.Errorf("expected heading text in output, got: %s", html)
		}
	})

	t.Run("bold κείμενο παράγει <strong>", func(t *testing.T) {
		html, err := MarkdownToHTML("**bold text**")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(html, "<strong>bold text</strong>") {
			t.Errorf("expected <strong>bold text</strong>, got: %s", html)
		}
	})

	t.Run("link παράγει <a href>", func(t *testing.T) {
		html, err := MarkdownToHTML("[click here](https://example.com)")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(html, `href="https://example.com"`) {
			t.Errorf("expected href attribute, got: %s", html)
		}
		if !strings.Contains(html, "click here") {
			t.Errorf("expected link text, got: %s", html)
		}
	})

	t.Run("GFM table παράγει <table>", func(t *testing.T) {
		md := "| Name  | Age |\n|-------|-----|\n| Alice | 30  |"
		html, err := MarkdownToHTML(md)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(html, "<table>") {
			t.Errorf("expected <table>, got: %s", html)
		}
		if !strings.Contains(html, "Alice") {
			t.Errorf("expected table content, got: %s", html)
		}
	})

	t.Run("unordered list παράγει <ul>", func(t *testing.T) {
		html, err := MarkdownToHTML("- item one\n- item two\n- item three")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(html, "<ul>") {
			t.Errorf("expected <ul>, got: %s", html)
		}
		if !strings.Contains(html, "<li>") {
			t.Errorf("expected <li> items, got: %s", html)
		}
	})

	t.Run("άδεια είσοδος επιστρέφει άδειο string", func(t *testing.T) {
		html, err := MarkdownToHTML("")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if strings.TrimSpace(html) != "" {
			t.Errorf("expected empty output for empty input, got: %q", html)
		}
	})

	t.Run("strikethrough (GFM)", func(t *testing.T) {
		html, err := MarkdownToHTML("~~deleted~~")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(html, "<del>deleted</del>") {
			t.Errorf("expected <del>deleted</del>, got: %s", html)
		}
	})
}
