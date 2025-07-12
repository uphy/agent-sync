package agent

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/user/agent-def/internal/model"
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
	return fmt.Sprintf("MCP tool (MCP Server: %s, Tool: %s, Arguments: %s)", agent, command, strings.Join(args, " "))
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
	commandContents := make([]string, len(commands))
	for i, cmd := range commands {
		commandContents[i] = cmd.Content
	}

	return strings.Join(commandContents, "\n\n---\n\n"), nil
}

// DefaultMemoryPath determines the output path for Claude agent memory files
func (c *Claude) DefaultMemoryPath(outputBaseDir string, userScope bool, fileName string) (string, error) {
	if userScope {
		return filepath.Join(outputBaseDir, ".claude", "CLAUDE.md"), nil
	} else {
		return filepath.Join(outputBaseDir, "CLAUDE.md"), nil
	}
}

// DefaultCommandPath determines the output path for Claude agent command files
func (c *Claude) DefaultCommandPath(outputBaseDir string, userScope bool, fileName string) (string, error) {
	ext := ".md"
	name := strings.TrimSuffix(fileName, ext)

	return filepath.Join(outputBaseDir, ".claude", "commands", name+ext), nil
}
