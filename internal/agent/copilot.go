package agent

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/uphy/agent-sync/internal/frontmatter"
	"github.com/uphy/agent-sync/internal/model"
)

// Copilot implements the Copilot-specific conversion logic
type Copilot struct{}

// Name returns the display name for Copilot agent
func (c *Copilot) Name() string {
	return "Copilot"
}

// ID returns the unique identifier for Copilot agent
func (c *Copilot) ID() string {
	return "copilot"
}

// FormatFile converts a path to Copilot's file reference format
func (c *Copilot) FormatFile(path string) string {
	// Copilot uses simple path format
	return fmt.Sprintf("`%s`", path)
}

// FormatMCP formats an MCP command for Copilot agent
func (c *Copilot) FormatMCP(agent, command string, args ...string) string {
	return formatMCP(agent, command, args...)
}

// FormatMemory processes a memory context for Copilot agent
func (c *Copilot) FormatMemory(content string) (string, error) {
	// For global instructions, simply return the content as is
	// No frontmatter is needed
	return content, nil
}

// FormatCommand processes command definitions for Copilot agent
func (c *Copilot) FormatCommand(commands []model.Command) (string, error) {
	if len(commands) == 0 {
		return "", nil
	}

	// Copilot does not support multiple commands in one file
	if len(commands) > 1 {
		return "", fmt.Errorf("copilot agent does not support multiple commands in one file")
	}

	cmd := commands[0]

	// Build frontmatter using UnmarshalSection("copilot") primarily,
	// but also accept top-level fields for backward compatibility.
	type copilotFm struct {
		Mode        string   `yaml:"mode,omitempty"`
		Model       string   `yaml:"model,omitempty"`
		Tools       []string `yaml:"tools,omitempty"`
		Description string   `yaml:"description,omitempty"`
	}
	var fm copilotFm

	// 1) Try to read nested "copilot" section via Command.UnmarshalSection
	if err := cmd.UnmarshalSection("copilot", &fm); err != nil && err.Error() != "frontmatter missing section \"copilot\"" {
		// Only ignore "missing section" errors. Propagate other marshal/unmarshal failures.
		return "", fmt.Errorf("copilot frontmatter parse error: %w", err)
	}

	// 2) No top-level fallback for mode/model/tools to keep a single source of truth.
	//    Only fallback Description from the common Command.Description field.

	// Final fallback: common top-level description field on Command
	if fm.Description == "" && cmd.Description != "" {
		fm.Description = cmd.Description
	}

	// Emit frontmatter only when at least one field is present.
	if fm.Mode == "" && fm.Model == "" && len(fm.Tools) == 0 && fm.Description == "" {
		return cmd.Content, nil
	}

	// Marshal frontmatter using shared helper
	yamlWithFences, err := frontmatter.Wrap(fm)
	if err != nil {
		return "", fmt.Errorf("failed to marshal copilot command frontmatter: %w", err)
	}

	// Add content
	return yamlWithFences + cmd.Content, nil
}

// DefaultMemoryPath determines the default output path for memory tasks
func (c *Copilot) DefaultMemoryPath(outputBaseDir string, userScope bool, fileName string) (string, error) {
	if userScope {
		// User scope goes in home directory
		return filepath.Join(homeDir(), ".vscode", "copilot-instructions.md"), nil
	}

	// Project scope uses .github/copilot-instructions.md
	return filepath.Join(outputBaseDir, ".github", "copilot-instructions.md"), nil
}

// DefaultCommandPath determines the default output path for command tasks
func (c *Copilot) DefaultCommandPath(outputBaseDir string, userScope bool, fileName string) (string, error) {
	if userScope {
		// User scope
		return filepath.Join(homeDir(), ".vscode", "prompts", fileName+".prompt.md"), nil
	}

	// Project scope
	if !strings.HasSuffix(fileName, ".prompt") {
		fileName += ".prompt"
	}
	return filepath.Join(outputBaseDir, ".github", "prompts", fileName+".md"), nil
}

// MemoryPath returns the default path for Copilot agent memory files
func (c *Copilot) MemoryPath(userScope bool) string {
	if userScope {
		return filepath.Join(".vscode", "copilot-instructions.md")
	}
	return filepath.Join(".github", "copilot-instructions.md")
}

// CommandPath returns the default path for Copilot agent command files
func (c *Copilot) CommandPath(userScope bool) string {
	if userScope {
		return filepath.Join(".vscode", "prompts") + "/"
	}
	return filepath.Join(".github", "prompts") + "/"
}

// ShouldConcatenate determines whether content should be concatenated
func (c *Copilot) ShouldConcatenate(taskType string) bool {
	if taskType == "memory" {
		// Memory type (custom instructions) should be concatenated
		return true
	}
	// Command type (prompts) should not be concatenated
	return false
}

// homeDir returns the user's home directory
func homeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return home
}

// FormatMode processes mode definitions for Copilot agent
func (c *Copilot) FormatMode(modes []model.Mode) (string, error) {
	// Copilot does not support multiple modes combined
	if len(modes) > 1 {
		return "", errors.New("copilot agent does not support multiple modes")
	}
	// No modes provided: return empty content with no error
	if len(modes) == 0 {
		return "", nil
	}
	// Single mode: return its content
	return modes[0].Content, nil
}

// ModePath returns the default path for Copilot agent mode files
func (c *Copilot) ModePath(userScope bool) string {
	return ""
}
