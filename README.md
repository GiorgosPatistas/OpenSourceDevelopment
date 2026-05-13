# my-ssg — Static Site Generator

A fast static site generator built with **TypeScript** (CLI) and **Go** (engine). Convert Markdown files into a complete HTML website using customizable templates.

---

## How It Works

The project is split into two parts that work together:

```
TypeScript CLI  →  finds the correct Go binary for your OS
Go Engine       →  reads Markdown, applies templates, outputs HTML
```

The TypeScript layer handles the user-facing command line interface, while the Go engine does the heavy lifting: parsing Markdown, processing YAML front matter, and rendering HTML templates.

---

## Project Structure

```
my-ssg/
├── src/
│   └── index.ts          # TypeScript CLI entry point
├── engine/
│   ├── main.go           # Go engine entry point
│   ├── parser/
│   │   ├── frontmatter.go  # YAML front matter parser
│   │   └── markdown.go     # Markdown → HTML (goldmark)
│   ├── renderer/
│   │   └── renderer.go   # HTML template renderer
│   ├── generator/
│   │   └── generator.go  # Site orchestrator
│   ├── build.sh          # Cross-compile script (Mac/Linux)
│   └── build.bat         # Cross-compile script (Windows)
├── bin/                  # Compiled Go binaries (git-ignored)
├── dist/                 # Compiled TypeScript (git-ignored)
├── example-site/         # Example site to test with
│   ├── content/          # Markdown source files
│   ├── templates/        # HTML templates
│   ├── static/           # CSS, images, etc.
│   └── ssg.config.json   # Site configuration
├── package.json
└── go.mod
```

---

## Prerequisites

- [Node.js](https://nodejs.org) v18+
- [pnpm](https://pnpm.io) (`npm install -g pnpm`)
- [Go](https://go.dev/dl) v1.22+

---

## Installation

**1. Clone the repository:**
```bash
git clone https://github.com/georgepatistas/my-ssg.git
cd my-ssg
```

**2. Install Node dependencies:**
```bash
pnpm install
```

**3. Compile the Go engine:**
```bash
# Windows
cd engine
.\build.bat

# Mac / Linux
cd engine
./build.sh
```

**4. Build the TypeScript CLI:**
```bash
pnpm build
```

---

## Usage

```bash
node dist/index.js build <path/to/your-site>
```

**Example:**
```bash
node dist/index.js build example-site
```

The generated site will be in `<your-site>/dist/`.

**To preview locally:**
```bash
cd example-site/dist
npx serve .
```

Then open `http://localhost:3000`.

---

## Site Structure

Your site folder must follow this structure:

```
your-site/
├── content/          # Markdown files (supports subdirectories)
│   ├── about.md
│   └── blog/
│       └── my-post.md
├── templates/
│   ├── layout.html   # Base layout (required)
│   ├── page.html     # Article template (defines "content" block)
│   └── index.html    # Home page template (defines "content" block)
├── static/           # Copied as-is to output (CSS, images, etc.)
└── ssg.config.json   # Optional site configuration
```

---

## Markdown Format

Each `.md` file should start with a YAML front matter block:

```markdown
---
title: "My Post Title"
date: 2024-06-15
description: "A short description"
draft: false
---

# Content starts here

Regular Markdown content...
```

| Field | Required | Description |
|---|---|---|
| `title` | No* | Page title |
| `date` | No | Publication date (affects sort order) |
| `description` | No | Short description for index and meta tags |
| `draft` | No | If `true`, page is excluded from build |
| `slug` | No | Custom URL slug |

*If omitted, the filename is used as the title.

**Supported Markdown features:** tables, task lists, strikethrough, footnotes, fenced code blocks, raw HTML, and more (via [goldmark](https://github.com/yuin/goldmark) with GFM extensions).

---

## Configuration

Create a `ssg.config.json` in your site's root folder:

```json
{
  "siteTitle": "My Blog",
  "siteUrl": "https://myblog.com",
  "description": "My personal blog"
}
```

| Field | Default | Description |
|---|---|---|
| `siteTitle` | `"My Site"` | Site name shown in header and title tag |
| `siteUrl` | `"/"` | Base URL of the site |
| `description` | `""` | Site description |
| `contentDir` | `"content"` | Directory for Markdown files |
| `templateDir` | `"templates"` | Directory for HTML templates |
| `staticDir` | `"static"` | Directory for static assets |
| `outputDir` | `"dist"` | Output directory |

---

## Templates

Templates use Go's `html/template` syntax. The `layout.html` defines the page skeleton and calls `{{template "content" .}}`. The `page.html` and `index.html` each define the `"content"` block that fills it.

**Available template variables:**

```
.Site.SiteTitle       → Site title from config
.Site.SiteURL         → Site base URL
.Page.Title           → Page title from front matter
.Page.Date            → Page date (time.Time)
.Page.Description     → Page description
.Page.URL             → Relative URL of the page
.Content              → Rendered HTML content
.Pages                → List of all published pages
```

**Built-in template functions:**

```
{{formatDate .Page.Date}}     → "15 June 2024"
{{formatDateISO .Page.Date}}  → "2024-06-15"
{{safeHTML .SomeVar}}         → Output HTML without escaping
```

---

## Tech Stack

| Layer | Technology | Purpose |
|---|---|---|
| CLI | TypeScript + Commander.js | User-facing command line interface |
| Engine | Go 1.22 | Markdown parsing, template rendering |
| Markdown | [goldmark](https://github.com/yuin/goldmark) | CommonMark + GFM extensions |
| YAML | [gopkg.in/yaml.v3](https://pkg.go.dev/gopkg.in/yaml.v3) | Front matter parsing |
| Templates | Go `html/template` | Safe HTML rendering |

---

## License

MIT — see [LICENSE](LICENSE) for details.
