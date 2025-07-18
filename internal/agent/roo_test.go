package agent

import (
	"testing"
)

func TestRoo_ID_Name(t *testing.T) {
	r := &Roo{}
	if r.ID() != "roo" {
		t.Errorf("expected ID 'roo', got %q", r.ID())
	}
	if r.Name() != "Roo" {
		t.Errorf("expected Name 'Roo', got %q", r.Name())
	}
}

// ShouldConcatenate method has been removed as part of transition to path-based concatenation behavior

func TestMCPFunc_WithValidAgent(t *testing.T) {
	// Create a Roo agent
	roo := &Roo{}

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
			result := roo.FormatMCP(tc.agent, tc.command, tc.args...)
			if result != tc.expected {
				t.Errorf("FormatMCP(%q, %q, %v) = %q, expected %q",
					tc.agent, tc.command, tc.args, result, tc.expected)
			}
		})
	}
}
