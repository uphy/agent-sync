---
name: test
description: Run the full test suite for the project
---

# Test

Executes unit and integration tests, producing verbose output on failures.

## Commands

```bash
# for JavaScript/TypeScript projects
if [ -f package.json ]; then
  npm install
  npm test -- --verbose
fi

# for Go projects
if [ -f go.mod ]; then
  go test ./... -v
fi
```

## Configuration

- JavaScript: tests under __tests__/ or *.test.js / *.spec.js  
- TypeScript: ts-jest in jest.config.js  
- Go: tests in *_test.go, coverage flags via -cover  

## Troubleshooting

- **Test failures**: review stack trace, rerun with `--runInBand` for Jest or `-timeout` flag for Go.  
- **Coverage reports**: add `--coverage` for JS or `-coverprofile=coverage.out` for Go.