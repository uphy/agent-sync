package agent

import (
	"fmt"
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

// Tests for deprecated methods (ShouldConcatenate, DefaultMemoryPath, DefaultCommandPath) have been removed

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
