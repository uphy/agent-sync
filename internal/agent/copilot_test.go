package agent

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/agent-def/internal/model"
)

func TestCopilot_ID_Name(t *testing.T) {
	c := &Copilot{}
	if c.ID() != "copilot" {
		t.Errorf("expected ID 'copilot', got %q", c.ID())
	}
	if c.Name() != "Copilot" {
		t.Errorf("expected Name 'Copilot', got %q", c.Name())
	}
}

func TestCopilot_FormatFile(t *testing.T) {
	c := &Copilot{}
	tests := []struct {
		path     string
		expected string
	}{
		{"path/to/file.go", "`path/to/file.go`"},
		{"", "``"},
		{"file with spaces.txt", "`file with spaces.txt`"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := c.FormatFile(tt.path)
			if result != tt.expected {
				t.Errorf("FormatFile(%q) = %q, want %q", tt.path, result, tt.expected)
			}
		})
	}
}

func TestCopilot_FormatMemory(t *testing.T) {
	c := &Copilot{}
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "normal content",
			content:  "# Memory Content\nThis is a test memory.",
			expected: "# Memory Content\nThis is a test memory.",
		},
		{
			name:     "empty content",
			content:  "",
			expected: "",
		},
		{
			name:     "multiline content",
			content:  "Line 1\nLine 2\nLine 3",
			expected: "Line 1\nLine 2\nLine 3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := c.FormatMemory(tt.content)
			if err != nil {
				t.Fatalf("FormatMemory returned unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("FormatMemory(%q) = %q, want %q", tt.content, result, tt.expected)
			}
		})
	}
}

func TestCopilot_FormatCommand(t *testing.T) {
	c := &Copilot{}

	t.Run("empty commands", func(t *testing.T) {
		result, err := c.FormatCommand([]model.Command{})
		if err != nil {
			t.Fatalf("FormatCommand returned unexpected error: %v", err)
		}
		if result != "" {
			t.Errorf("FormatCommand(empty) = %q, want empty string", result)
		}
	})

	t.Run("multiple commands", func(t *testing.T) {
		commands := []model.Command{
			{Content: "Command 1"},
			{Content: "Command 2"},
		}
		_, err := c.FormatCommand(commands)
		if err == nil {
			t.Errorf("FormatCommand(multiple) should return error, got nil")
		}
		if !strings.Contains(err.Error(), "does not support multiple commands") {
			t.Errorf("Error message should mention multiple commands, got: %v", err)
		}
	})

	t.Run("complete frontmatter", func(t *testing.T) {
		cmd := model.Command{
			Content: "Command content",
			Copilot: model.Copilot{
				Mode:        "edit",
				Model:       "gpt-4",
				Tools:       []string{"tool1", "tool2"},
				Description: "Test description",
			},
		}
		result, err := c.FormatCommand([]model.Command{cmd})
		if err != nil {
			t.Fatalf("FormatCommand returned unexpected error: %v", err)
		}

		// Check that frontmatter contains all expected fields at the top level
		if !strings.Contains(result, "mode: edit") {
			t.Errorf("Frontmatter missing mode, got: %s", result)
		}
		if !strings.Contains(result, "model: gpt-4") {
			t.Errorf("Frontmatter missing model, got: %s", result)
		}
		if !strings.Contains(result, "- tool1") && !strings.Contains(result, "- tool2") {
			t.Errorf("Frontmatter missing tools, got: %s", result)
		}
		if !strings.Contains(result, "description: Test description") {
			t.Errorf("Frontmatter missing description, got: %s", result)
		}
		// Ensure the copilot nested structure is not present
		if strings.Contains(result, "copilot:") {
			t.Errorf("Frontmatter should not contain nested 'copilot:' structure, got: %s", result)
		}

		// Check that content is preserved
		if !strings.Contains(result, "Command content") {
			t.Errorf("Result missing command content, got: %s", result)
		}
	})

	t.Run("partial frontmatter", func(t *testing.T) {
		cmd := model.Command{
			Content: "Command content",
			Copilot: model.Copilot{
				Mode:  "edit",
				Tools: []string{"tool1"},
			},
		}
		result, err := c.FormatCommand([]model.Command{cmd})
		if err != nil {
			t.Fatalf("FormatCommand returned unexpected error: %v", err)
		}

		// Check that frontmatter contains expected fields at the top level
		if !strings.Contains(result, "mode: edit") {
			t.Errorf("Frontmatter missing mode, got: %s", result)
		}
		if !strings.Contains(result, "- tool1") {
			t.Errorf("Frontmatter missing tools, got: %s", result)
		}

		// Check that empty fields are not included
		if strings.Contains(result, "model:") {
			t.Errorf("Frontmatter should not include empty model, got: %s", result)
		}
		if strings.Contains(result, "description:") {
			t.Errorf("Frontmatter should not include empty description, got: %s", result)
		}

		// Ensure the copilot nested structure is not present
		if strings.Contains(result, "copilot:") {
			t.Errorf("Frontmatter should not contain nested 'copilot:' structure, got: %s", result)
		}
	})
}

func TestCopilot_DefaultMemoryPath(t *testing.T) {
	c := &Copilot{}
	// Store original home directory
	origHomeDir := os.Getenv("HOME")
	os.Setenv("HOME", "/mock/home")
	defer os.Setenv("HOME", origHomeDir)

	t.Run("user scope", func(t *testing.T) {

		path, err := c.DefaultMemoryPath("output", true, "test")
		if err != nil {
			t.Fatalf("DefaultMemoryPath returned unexpected error: %v", err)
		}
		expected := filepath.Join("/mock/home", ".vscode", "copilot-instructions.md")
		if path != expected {
			t.Errorf("DefaultMemoryPath(user) = %q, want %q", path, expected)
		}
	})

	t.Run("project scope", func(t *testing.T) {
		path, err := c.DefaultMemoryPath("output", false, "test")
		if err != nil {
			t.Fatalf("DefaultMemoryPath returned unexpected error: %v", err)
		}
		expected := filepath.Join("output", ".github", "copilot-instructions.md")
		if path != expected {
			t.Errorf("DefaultMemoryPath(project) = %q, want %q", path, expected)
		}
	})
}

func TestCopilot_DefaultCommandPath(t *testing.T) {
	c := &Copilot{}
	// Store original home directory
	origHomeDir := os.Getenv("HOME")
	os.Setenv("HOME", "/mock/home")
	defer os.Setenv("HOME", origHomeDir)

	t.Run("user scope", func(t *testing.T) {

		path, err := c.DefaultCommandPath("output", true, "deploy")
		if err != nil {
			t.Fatalf("DefaultCommandPath returned unexpected error: %v", err)
		}
		expected := filepath.Join("/mock/home", ".vscode", "prompts", "deploy.prompt.md")
		if path != expected {
			t.Errorf("DefaultCommandPath(user) = %q, want %q", path, expected)
		}
	})

	t.Run("project scope without .prompt suffix", func(t *testing.T) {
		path, err := c.DefaultCommandPath("output", false, "deploy")
		if err != nil {
			t.Fatalf("DefaultCommandPath returned unexpected error: %v", err)
		}
		expected := filepath.Join("output", ".github", "prompts", "deploy.prompt.md")
		if path != expected {
			t.Errorf("DefaultCommandPath(project) = %q, want %q", path, expected)
		}
	})

	t.Run("project scope with .prompt suffix", func(t *testing.T) {
		path, err := c.DefaultCommandPath("output", false, "deploy.prompt")
		if err != nil {
			t.Fatalf("DefaultCommandPath returned unexpected error: %v", err)
		}
		expected := filepath.Join("output", ".github", "prompts", "deploy.prompt.md")
		if path != expected {
			t.Errorf("DefaultCommandPath(project) = %q, want %q", path, expected)
		}
	})
}

func TestCopilot_ShouldConcatenate(t *testing.T) {
	c := &Copilot{}

	tests := []struct {
		taskType string
		expected bool
	}{
		{"memory", true},
		{"command", false},
		{"unknown", false},
	}

	for _, tt := range tests {
		t.Run(tt.taskType, func(t *testing.T) {
			result := c.ShouldConcatenate(tt.taskType)
			if result != tt.expected {
				t.Errorf("Copilot.ShouldConcatenate(%q) = %v, want %v", tt.taskType, result, tt.expected)
			}
		})
	}
}

func TestCopilot_FormatMCP(t *testing.T) {
	c := &Copilot{}

	testCases := []struct {
		name     string
		agent    string
		command  string
		args     []string
		expected string
	}{
		{
			name:     "command without args",
			agent:    "github",
			command:  "get-issue",
			args:     []string{},
			expected: "MCP tool `github.get-issue`",
		},
		{
			name:     "command with single arg",
			agent:    "jira",
			command:  "get-ticket",
			args:     []string{"PROJ-123"},
			expected: "MCP tool `jira.get-ticket(PROJ-123)`",
		},
		{
			name:     "command with multiple args",
			agent:    "github",
			command:  "get-issue",
			args:     []string{"owner", "repo", "123"},
			expected: "MCP tool `github.get-issue(owner, repo, 123)`",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := c.FormatMCP(tc.agent, tc.command, tc.args...)
			if result != tc.expected {
				t.Errorf("FormatMCP(%q, %q, %v) = %q, expected %q",
					tc.agent, tc.command, tc.args, result, tc.expected)
			}
		})
	}
}
