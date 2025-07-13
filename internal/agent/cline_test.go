package agent

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestCline_ID_Name(t *testing.T) {
	c := &Cline{}
	if c.ID() != "cline" {
		t.Errorf("expected ID 'cline', got %q", c.ID())
	}
	if c.Name() != "Cline" {
		t.Errorf("expected Name 'Cline', got %q", c.Name())
	}
}

func TestCline_FormatFile(t *testing.T) {
	c := &Cline{}

	tests := []struct {
		path     string
		expected string
	}{
		{"/absolute/path/file.txt", "/absolute/path/file.txt"},
		{"relative/path/file.txt", "@/relative/path/file.txt"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("FormatFile(%s)", tt.path), func(t *testing.T) {
			result := c.FormatFile(tt.path)
			if result != tt.expected {
				t.Errorf("Cline.FormatFile(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestCline_ShouldConcatenate(t *testing.T) {
	c := &Cline{}

	tests := []struct {
		taskType string
		expected bool
	}{
		{"memory", false},
		{"command", false},
		{"unknown", false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("TaskType=%s", tt.taskType), func(t *testing.T) {
			result := c.ShouldConcatenate(tt.taskType)
			if result != tt.expected {
				t.Errorf("Cline.ShouldConcatenate(%q) = %v, want %v", tt.taskType, result, tt.expected)
			}
		})
	}
}

func TestCline_DefaultMemoryPath(t *testing.T) {
	c := &Cline{}

	tests := []struct {
		name        string
		outputDir   string
		userScope   bool
		fileName    string
		expected    string
		expectError bool
	}{
		{
			name:      "project scope",
			outputDir: "/test/path",
			userScope: false,
			fileName:  "rules.md",
			expected:  "/test/path/.clinerules/rules.md",
		},
		{
			name:      "user scope",
			outputDir: "/home/user",
			userScope: true,
			fileName:  "context.md",
			expected:  filepath.Join("/home/user", "Documents", "Cline", "Rules", "context.md"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := c.DefaultMemoryPath(tt.outputDir, tt.userScope, tt.fileName)

			if tt.expectError && err == nil {
				t.Fatal("expected error, got nil")
			}

			if !tt.expectError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result != tt.expected {
				t.Errorf("DefaultMemoryPath(%q, %v, %q) = %q, want %q",
					tt.outputDir, tt.userScope, tt.fileName, result, tt.expected)
			}
		})
	}
}

func TestCline_DefaultCommandPath(t *testing.T) {
	c := &Cline{}

	tests := []struct {
		name        string
		outputDir   string
		userScope   bool
		fileName    string
		expected    string
		expectError bool
	}{
		{
			name:      "project scope",
			outputDir: "/test/path",
			userScope: false,
			fileName:  "workflow.md",
			expected:  "/test/path/.clinerules/workflows/workflow.md",
		},
		{
			name:      "user scope",
			outputDir: "/home/user",
			userScope: true,
			fileName:  "workflow.md",
			expected:  filepath.Join("/home/user", "Documents", "Cline", "Workflows", "workflow.md"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := c.DefaultCommandPath(tt.outputDir, tt.userScope, tt.fileName)

			if tt.expectError && err == nil {
				t.Fatal("expected error, got nil")
			}

			if !tt.expectError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result != tt.expected {
				t.Errorf("DefaultCommandPath(%q, %v, %q) = %q, want %q",
					tt.outputDir, tt.userScope, tt.fileName, result, tt.expected)
			}
		})
	}
}

func TestCline_FormatMCP(t *testing.T) {
	c := &Cline{}

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
