package agent

import (
	"strings"
	"testing"

	"github.com/uphy/agent-sync/internal/model"
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

func TestRoo_FormatCommand_FrontmatterCases(t *testing.T) {
	r := &Roo{}

	makeCmd := func(raw map[string]any, body string) model.Command {
		// Simulate ParseCommand behavior: top-level description also exists in Raw.
		// Our FormatCommand now reads via UnmarshalSection("", &meta) which pulls from Raw[""].
		// For tests, we emulate this by ensuring Raw has the keys used by UnmarshalSection.
		return model.Command{
			Raw:     raw,
			Content: body,
		}
	}

	body := "Command body content"

	t.Run("only top-level description", func(t *testing.T) {
		// Simulate ParseCommand setting Command.Description, while Raw may not be used for top-level.
		cmd := makeCmd(map[string]any{}, body)
		cmd.Description = "Top description"

		out, err := r.FormatCommand([]model.Command{cmd})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !strings.HasPrefix(out, "---\n") {
			t.Fatalf("expected frontmatter, got: %q", out)
		}
		if !strings.Contains(out, "description: Top description") {
			t.Errorf("expected description in frontmatter, got: %q", out)
		}
		if strings.Contains(out, "roo:") {
			t.Errorf("did not expect roo section")
		}
		if !strings.HasSuffix(out, "\n\n"+body) && !strings.HasSuffix(out, "\n"+body) {
			t.Errorf("expected body appended after frontmatter, got: %q", out)
		}
	})

	t.Run("only roo.argument-hint (should error without any description)", func(t *testing.T) {
		cmd := makeCmd(map[string]any{
			"roo": map[string]any{
				"argument-hint": "ARG HINT",
			},
		}, body)

		_, err := r.FormatCommand([]model.Command{cmd})
		if err == nil {
			t.Fatalf("expected error due to missing description, got nil")
		}
		if !strings.Contains(err.Error(), "requires a description") {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("both description and roo.argument-hint", func(t *testing.T) {
		cmd := makeCmd(map[string]any{
			"roo": map[string]any{
				"argument-hint": "HINT",
			},
		}, body)
		cmd.Description = "Desc"

		out, err := r.FormatCommand([]model.Command{cmd})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !strings.Contains(out, "description: Desc") {
			t.Errorf("expected description in frontmatter, got: %q", out)
		}
		if !strings.Contains(out, "argument-hint: HINT") {
			t.Errorf("expected roo.argument-hint, got: %q", out)
		}
	})

	t.Run("neither description nor roo.argument-hint present -> should error", func(t *testing.T) {
		cmd := makeCmd(map[string]any{}, body)

		_, err := r.FormatCommand([]model.Command{cmd})
		if err == nil {
			t.Fatalf("expected error due to missing description, got nil")
		}
		if !strings.Contains(err.Error(), "requires a description") {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
