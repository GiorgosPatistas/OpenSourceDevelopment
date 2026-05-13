# Contributing to my-ssg

Thank you for your interest in contributing! This guide will help you get started.

---

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [How to Contribute](#how-to-contribute)
- [Development Setup](#development-setup)
- [Commit Conventions](#commit-conventions)
- [Pull Request Process](#pull-request-process)

---

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating, you agree to uphold it.

---

## Getting Started

1. **Fork** the repository on GitHub
2. **Clone** your fork locally:
   ```bash
   git clone https://github.com/YOUR-USERNAME/OpenSourceDevelopment.git
   cd OpenSourceDevelopment
   ```
3. **Add the upstream remote** so you can sync changes:
   ```bash
   git remote add upstream https://github.com/GiorgosPatistas/OpenSourceDevelopment.git
   ```

---

## How to Contribute

### Reporting Bugs

- Search [existing issues](https://github.com/GiorgosPatistas/OpenSourceDevelopment/issues) first to avoid duplicates
- Use the **Bug Report** issue template
- Include steps to reproduce, expected vs. actual behavior, and your OS/Go/Node versions

### Suggesting Features

- Open a **Feature Request** issue with a clear description of the problem it solves
- Discuss the idea before starting implementation to avoid wasted effort

### Submitting Code

- Look for issues labeled `good first issue` or `help wanted`
- Comment on the issue to let others know you're working on it
- Keep PRs focused — one fix or feature per PR

---

## Development Setup

### Prerequisites

- [Go](https://go.dev/dl) v1.22+
- [Node.js](https://nodejs.org) v18+
- [pnpm](https://pnpm.io) (`npm install -g pnpm`)

### Setup

```bash
# Install Node dependencies
pnpm install

# Build the TypeScript CLI
pnpm build

# Compile the Go engine (pick your OS)
cd engine
./build.sh        # macOS / Linux
.\build.bat       # Windows

# Test with the example site
node dist/index.js build example-site
```

### Running Tests

```bash
# Go tests
cd engine && go test ./...

# TypeScript tests
pnpm test
```

---

## Commit Conventions

We use [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <short description>
```

| Type | When to use |
|---|---|
| `feat` | A new feature |
| `fix` | A bug fix |
| `docs` | Documentation changes only |
| `test` | Adding or fixing tests |
| `refactor` | Code change that is neither a fix nor a feature |
| `chore` | Build process, dependency updates, etc. |

**Examples:**
```
feat(engine): add syntax highlighting for code blocks
fix(parser): handle missing front matter gracefully
docs: update installation instructions for Windows
```

---

## Pull Request Process

1. **Create a branch** from `main` with a descriptive name:
   ```bash
   git checkout -b feat/syntax-highlighting
   ```

2. **Make your changes**, keeping commits small and focused

3. **Sync with upstream** before opening a PR:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

4. **Open a Pull Request** against `main` and fill in the PR template

5. **Address review feedback** — maintainers may request changes

6. Once approved, your PR will be **squash-merged** into `main`

---

Thank you for contributing!
