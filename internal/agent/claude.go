package agent

import (
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/uphy/agent-sync/internal/frontmatter"
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
		var formattedCmd strings.Builder

		// Extract Claude-specific frontmatter from generic Command.Raw
		type claudeFm struct {
			Description  string `yaml:"description,omitempty"`
			AllowedTools string `yaml:"allowed-tools,omitempty"`
		}
		var fm claudeFm
		if cmd.Raw != nil {
			_ = cmd.UnmarshalSection("claude", &fm)
		}
		// Fallback to common description if not provided under claude
		if fm.Description == "" && cmd.Description != "" {
			fm.Description = cmd.Description
		}

		// Include frontmatter if Claude-specific attributes exist
		if fm.Description != "" || fm.AllowedTools != "" {
			yamlWithFences, err := frontmatter.Wrap(fm)
			if err != nil {
				return "", fmt.Errorf("failed to marshal claude command frontmatter: %w", err)
			}
			formattedCmd.WriteString(yamlWithFences)
		}

		// Add the main content
		formattedCmd.WriteString(cmd.Content)

		formattedCommands = append(formattedCommands, formattedCmd.String())
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

	// Define ClaudeMode; tools are plain []string (no special normalization)
	type ClaudeMode struct {
		Name        string   `yaml:"name"`
		Description string   `yaml:"description"`
		Tools       []string `yaml:"tools,omitempty"`
	}

	var cm ClaudeMode
	// Unmarshal section if present
	if mode.Raw != nil {
		_ = mode.UnmarshalSection("claude", &cm)
	}

	// Fallback to common description when empty
	if cm.Description == "" && mode.Description != "" {
		cm.Description = mode.Description
	}

	// Encode with goccy/go-yaml, then post-process to match expected quoting style
	var sb strings.Builder
	enc := yaml.NewEncoder(
		&sb,
		yaml.Indent(2),
		yaml.UseSingleQuote(false),
		yaml.IndentSequence(true),
		yaml.UseLiteralStyleIfMultiline(true),
	)
	if err := enc.Encode(cm); err != nil {
		return "", fmt.Errorf("failed to marshal claude mode frontmatter: %w", err)
	}
	yamlOut := strings.TrimRight(sb.String(), "\n") // avoid trailing blank line before closing ---

	// Tests expect quoted scalars and quoted list items. goccy/yaml encoder
	// doesn't force quotes for plain scalars, so patch them deterministically.
	// Only touch the specific fields we emit: name, description, tools items.
	var out strings.Builder
	lines := strings.Split(yamlOut, "\n")
	for i, line := range lines {
		switch {
		case strings.HasPrefix(line, "name: "):
			val := strings.TrimSpace(strings.TrimPrefix(line, "name: "))
			if val != "" && !strings.HasPrefix(val, "\"") && !strings.HasSuffix(val, "\"") {
				line = "name: " + fmt.Sprintf("%q", val)
			}
		case strings.HasPrefix(line, "description: "):
			val := strings.TrimSpace(strings.TrimPrefix(line, "description: "))
			if val != "" && !strings.HasPrefix(val, "\"") && !strings.HasSuffix(val, "\"") {
				line = "description: " + fmt.Sprintf("%q", val)
			}
		case strings.HasPrefix(strings.TrimLeft(line, " "), "- "):
			// Quote list items under tools
			trimmed := strings.TrimSpace(strings.TrimPrefix(strings.TrimLeft(line, " "), "- "))
			if trimmed != "" && !strings.HasPrefix(trimmed, "\"") && !strings.HasSuffix(trimmed, "\"") {
				indent := line[:len(line)-len(strings.TrimLeft(line, " "))]
				line = indent + "- " + fmt.Sprintf("%q", trimmed)
			}
		}
		out.WriteString(line)
		if i < len(lines)-1 {
			out.WriteByte('\n')
		}
	}

	frontmatter := "---\n" + out.String() + "\n---\n\n"
	return frontmatter + mode.Content, nil
}

// ModePath returns the default path for Claude agent mode files
func (c *Claude) ModePath(userScope bool) string {
	return ".claude/agents/"
}
