---
title: "How I Built a Static Site Generator with Go"
description: "The story behind this project"
date: 2024-07-15
---

# How I Built a Static Site Generator with Go

I wanted to understand how an SSG works under the hood, so I decided to build one from scratch.

## Architecture

The project is made up of two parts:

1. **TypeScript CLI** — the interface the user interacts with
2. **Go Engine** — the core that does the heavy lifting

The reason for this split architecture is that Go delivers excellent performance for file processing, while the TypeScript ecosystem is ideal for building CLIs.

## Technical Details

The Go engine uses:
- **goldmark** for markdown parsing (CommonMark compliant)
- **html/template** for safe HTML rendering
- **gopkg.in/yaml.v3** for front matter parsing

Total dependencies: just 2! That's the beauty of Go.
