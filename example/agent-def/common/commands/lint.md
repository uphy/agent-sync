---
name: lint
description: Run project static analysis and code formatting checks
---

# Lint

Ensures that code meets style guidelines, catches common errors, and formats source files.

## Commands

```bash
# for JavaScript/TypeScript projects
if [ -f package.json ]; then
  npm install
  npm run lint
  npm run format:check
fi

# for Go projects
if [ -f go.mod ]; then
  go fmt ./...
  golangci-lint run ./...
fi
```

## Configuration

- JavaScript/TypeScript:
  - ESLint config in `.eslintrc.js` or `.eslintrc.json`
  - Prettier config in `.prettierrc`
- Go:
  - golangci-lint config in `.golangci.yml`

## Troubleshooting

- **Missing ESLint/Prettier**: ensure `npm install` completes without errors.
- **golangci-lint not found**: install via `brew install golangci-lint` or add to your PATH.