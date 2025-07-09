---
name: coding-style
description: Common coding style guidelines for all projects
---

# Indentation
- Use 2 spaces for JavaScript/TypeScript; 4 spaces for Go.

# Line Length
- Wrap lines at 80 characters.

# Naming Conventions
- camelCase for variables and functions in JS/TS.
- PascalCase for React components and types.
- snake_case for file names and constants in Go.

# Comments
- Use JSDoc or GoDoc style for public APIs.
- Keep comments concise and relevant.

# Imports & Organization
- Group imports by external, internal, and styles.
- Alphabetize within each group.

# Tests
- Place tests in `__tests__/` folders or files ending with `.test`.
- Write tests before refactoring code.

# Commit Messages
- Use present tense, imperative mood.
- Reference issue keys when available.
- Keep summary under 50 characters.

# Formatting
- Run `prettier --write` after each feature.
- Run `go fmt ./...` before each commit.

# Additional Rules
- Avoid `console.log` in production code.
- Provide meaningful error messages.