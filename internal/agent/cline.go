package agent

import (
	"path/filepath"
	"strings"

	"github.com/user/agent-def/internal/model"
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

// DefaultMemoryPath determines the output path for Cline agent memory files
func (c *Cline) DefaultMemoryPath(outputBaseDir string, userScope bool, fileName string) (string, error) {
	ext := ".md"
	name := strings.TrimSuffix(fileName, ext)

	if userScope {
		return filepath.Join(outputBaseDir, "Documents", "Cline", "Rules", name+ext), nil
	}

	// Project scope
	return filepath.Join(outputBaseDir, ".clinerules", name+ext), nil
}

// DefaultCommandPath determines the output path for Cline agent workflow files
func (c *Cline) DefaultCommandPath(outputBaseDir string, userScope bool, fileName string) (string, error) {
	// Keep the original file name as-is
	if userScope {
		return filepath.Join(outputBaseDir, "Documents", "Cline", "Workflows", fileName), nil
	}

	// Project scope
	return filepath.Join(outputBaseDir, ".clinerules", "workflows", fileName), nil
}

// ShouldConcatenate determines if content should be concatenated for Cline agent
func (c *Cline) ShouldConcatenate(taskType string) bool {
	return false
}
