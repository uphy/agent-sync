package agent

import (
	"errors"
	"path/filepath"

	"github.com/uphy/agent-sync/internal/model"
)

// Cline implements the Cline-specific conversion logic
type Cline struct{}

// ID returns the unique identifier for Cline agent
func (c *Cline) ID() string {
	return "cline"
}

// Name returns the display name for Cline agent
func (c *Cline) Name() string {
	return "Cline"
}

// FormatFile converts a path to Cline's file reference format
func (c *Cline) FormatFile(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return "@/" + path
}

// FormatMCP formats an MCP command for Cline agent
func (c *Cline) FormatMCP(agent, command string, args ...string) string {
	return formatMCP(agent, command, args...)
}

// FormatMemory processes a memory context for Cline agent
func (c *Cline) FormatMemory(content string) (string, error) {
	// Process the memory context for Cline
	// For now, we're simply returning the content as is
	// Any Cline-specific formatting rules could be applied here
	return content, nil
}

// FormatCommand processes command (workflow) definitions for Cline agent
func (c *Cline) FormatCommand(commands []model.Command) (string, error) {
	// For Cline, we simply use the original markdown content
	if len(commands) == 0 {
		return "", nil
	}

	// Just return the raw markdown content for commands
	// If there are multiple commands, we only use the first one
	// since each command is processed separately (ShouldConcatenate returns false)
	return commands[0].Content, nil
}

// MemoryPath returns the default path for Cline agent memory files
func (c *Cline) MemoryPath(userScope bool) string {
	if userScope {
		return "Documents/Cline/Rules/"
	}
	return ".clinerules/"
}

// CommandPath returns the default path for Cline agent command files
func (c *Cline) CommandPath(userScope bool) string {
	if userScope {
		return "Documents/Cline/Workflows/"
	}
	return ".clinerules/workflows/"
}

// FormatMode processes mode definitions for Cline agent
func (c *Cline) FormatMode(modes []model.Mode) (string, error) {
	// Cline does not support combining multiple modes
	if len(modes) > 1 {
		return "", errors.New("cline agent does not support multiple modes")
	}
	// No modes: return empty content without error
	if len(modes) == 0 {
		return "", nil
	}
	// Single mode: return its content as-is
	return modes[0].Content, nil
}

// ModePath returns the default path for Cline agent mode files
func (c *Cline) ModePath(userScope bool) string {
	return ""
}
