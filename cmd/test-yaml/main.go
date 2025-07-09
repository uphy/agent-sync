package main

import (
	"fmt"
	"log"

	"github.com/user/agent-def/internal/agent"
	"github.com/user/agent-def/internal/model"
)

func main() {
	// Create a Roo agent
	roo := &agent.Roo{}

	// Create a test command with multiline content
	testCmd := model.Command{
		Slug:           "test",
		Name:           "test",
		RoleDefinition: "Test role\nwith multiple\nlines",
		WhenToUse:      "Use this when\nyou need to test",
		Groups:         []string{"test"},
		Content:        "# Test\n\nExecutes unit and integration tests, producing verbose output on failures.\n\n## Commands\n\n```bash\n# for JavaScript/TypeScript projects\nif [ -f package.json ]; then\n  npm install\n  npm test -- --verbose\nfi\n\n# for Go projects\nif [ -f go.mod ]; then\n  go test ./... -v\nfi\n```\n\n## Configuration\n\n- JavaScript: tests under __tests__/ or *.test.js / *.spec.js  \n- TypeScript: ts-jest in jest.config.js  \n- Go: tests in *_test.go, coverage flags via -cover  \n\n## Troubleshooting\n\n- **Test failures**: review stack trace, rerun with `--runInBand` for Jest or `-timeout` flag for Go.  \n- **Coverage reports**: add `--coverage` for JS or `-coverprofile=coverage.out` for Go.",
	}

	// Format the command
	formatted, err := roo.FormatCommand([]model.Command{testCmd})
	if err != nil {
		log.Fatalf("Failed to format command: %v", err)
	}

	// Print the formatted output
	fmt.Println(formatted)
}
