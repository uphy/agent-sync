package agent

import (
	"fmt"
	"strings"

	"github.com/uphy/agent-sync/internal/model"
)

// Claude implements the Claude-specific conversion logic
type Claude struct{}

// ID returns the unique identifier for Claude agent
func (c *Claude) ID() string {
	return "claude"
}

// Name returns the display name for Claude agent
func (c *Claude) Name() string {
	return "Claude"
}

// FormatFile converts a path to Claude's file reference format
func (c *Claude) FormatFile(path string) string {
	return "@" + path
}

// FormatMCP formats an MCP command for Claude agent
func (c *Claude) FormatMCP(agent, command string, args ...string) string {
	return formatMCP(agent, command, args...)
}

// FormatMemory processes a memory context for Claude agent
func (c *Claude) FormatMemory(content string) (string, error) {
	// Process the memory context for Claude
	// For now, we're simply returning the content as is
	// Any Claude-specific formatting rules could be applied here
	return content, nil
}

// FormatCommand processes command definitions for Claude agent
func (c *Claude) FormatCommand(commands []model.Command) (string, error) {
	// Format commands according to Claude's requirements
	var formattedCommands []string

	for _, cmd := range commands {
		var formattedCmd string

		// Include frontmatter if Claude-specific attributes exist
		hasFrontmatter := cmd.Claude.Description != "" || cmd.Claude.AllowedTools != ""

		if hasFrontmatter {
			formattedCmd = "---\n"

			if cmd.Claude.Description != "" {
				formattedCmd += fmt.Sprintf("description: %s\n", cmd.Claude.Description)
			}

			if cmd.Claude.AllowedTools != "" {
				formattedCmd += fmt.Sprintf("allowed-tools: %s\n", cmd.Claude.AllowedTools)
			}

			formattedCmd += "---\n\n"
		}

		// Add the main content
		formattedCmd += cmd.Content

		formattedCommands = append(formattedCommands, formattedCmd)
	}

	return strings.Join(formattedCommands, "\n\n---\n\n"), nil
}

// MemoryPath returns the default path for Claude agent memory files
func (c *Claude) MemoryPath(userScope bool) string {
	if userScope {
		return ".claude/CLAUDE.md"
	}
	return "CLAUDE.md"
}

// CommandPath returns the default path for Claude agent command files
func (c *Claude) CommandPath(userScope bool) string {
	return ".claude/commands/"
}

// FormatMode processes mode definitions for Claude agent
// Claude does not support combining multiple modes.
func (c *Claude) FormatMode(modes []model.Mode) (string, error) {
	// No modes: return empty string without error
	if len(modes) == 0 {
		return "", nil
	}
	// Multiple modes are not supported
	if len(modes) > 1 {
		return "", fmt.Errorf("claude agent does not support multiple modes")
	}

	mode := modes[0]

	frontmatter := "---\n"
	if mode.Claude != nil {
		if mode.Claude.Name != "" {
			frontmatter += fmt.Sprintf("name: %q\n", mode.Claude.Name)
		}
		if mode.Claude.Description != "" {
			frontmatter += fmt.Sprintf("description: %q\n", mode.Claude.Description)
		}
		if len(mode.Claude.Tools) > 0 {
			frontmatter += "tools:\n"
			for _, tool := range mode.Claude.Tools {
				frontmatter += fmt.Sprintf("  - %q\n", tool)
			}
		}
	}
	frontmatter += "---\n\n"
	return frontmatter + mode.Content, nil
}

// ModePath returns the default path for Claude agent mode files
func (c *Claude) ModePath(userScope bool) string {
	return ".claude/agents/"
}
