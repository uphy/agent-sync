package agent

import (
	"strings"
	"testing"

	"github.com/uphy/agent-def/internal/model"
)

func TestClaude_FormatCommand(t *testing.T) {
	c := &Claude{}

	tests := []struct {
		name     string
		commands []model.Command
		want     string
	}{
		{
			name: "No frontmatter",
			commands: []model.Command{
				{
					Content: "# Test Command\n\nThis is a test command with no frontmatter.",
					Path:    "test.md",
				},
			},
			want: "# Test Command\n\nThis is a test command with no frontmatter.",
		},
		{
			name: "With description",
			commands: []model.Command{
				{
					Claude: model.Claude{
						Description: "Test description",
					},
					Content: "# Test Command\n\nThis is a test command with description.",
					Path:    "test.md",
				},
			},
			want: "---\ndescription: Test description\n---\n\n# Test Command\n\nThis is a test command with description.",
		},
		{
			name: "With allowed-tools",
			commands: []model.Command{
				{
					Claude: model.Claude{
						AllowedTools: "Bash(git status:*)",
					},
					Content: "# Test Command\n\nThis is a test command with allowed-tools.",
					Path:    "test.md",
				},
			},
			want: "---\nallowed-tools: Bash(git status:*)\n---\n\n# Test Command\n\nThis is a test command with allowed-tools.",
		},
		{
			name: "With both attributes",
			commands: []model.Command{
				{
					Claude: model.Claude{
						Description:  "Test description",
						AllowedTools: "Bash(git status:*), Bash(git commit:*)",
					},
					Content: "# Test Command\n\nThis is a test command with both description and allowed-tools.",
					Path:    "test.md",
				},
			},
			want: "---\ndescription: Test description\nallowed-tools: Bash(git status:*), Bash(git commit:*)\n---\n\n# Test Command\n\nThis is a test command with both description and allowed-tools.",
		},
		{
			name: "Multiple commands",
			commands: []model.Command{
				{
					Claude: model.Claude{
						Description: "First command",
					},
					Content: "# First Command\n\nThis is the first command.",
					Path:    "first.md",
				},
				{
					Claude: model.Claude{
						AllowedTools: "Bash(ls:*)",
					},
					Content: "# Second Command\n\nThis is the second command.",
					Path:    "second.md",
				},
			},
			want: "---\ndescription: First command\n---\n\n# First Command\n\nThis is the first command.\n\n---\n\n---\nallowed-tools: Bash(ls:*)\n---\n\n# Second Command\n\nThis is the second command.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := c.FormatCommand(tt.commands)
			if err != nil {
				t.Errorf("Claude.FormatCommand() error = %v", err)
				return
			}

			// Normalize line endings for comparison
			got := strings.ReplaceAll(result, "\r\n", "\n")
			want := strings.ReplaceAll(tt.want, "\r\n", "\n")

			if got != want {
				t.Errorf("Claude.FormatCommand() = %q, want %q", got, want)
			}
		})
	}
}
