package agent

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/uphy/agent-def/internal/model"
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

	// Build frontmatter
	frontmatterMap := map[string]interface{}{}

	// Set each field directly at the top level
	if cmd.Copilot.Mode != "" {
		frontmatterMap["mode"] = cmd.Copilot.Mode
	}

	if cmd.Copilot.Model != "" {
		frontmatterMap["model"] = cmd.Copilot.Model
	}

	if len(cmd.Copilot.Tools) > 0 {
		frontmatterMap["tools"] = cmd.Copilot.Tools
	}

	if cmd.Copilot.Description != "" {
		frontmatterMap["description"] = cmd.Copilot.Description
	}

	// Marshal to YAML
	yamlBytes, err := yaml.Marshal(frontmatterMap)
	if err != nil {
		return "", fmt.Errorf("failed to marshal frontmatter: %w", err)
	}

	frontmatter := fmt.Sprintf("---\n%s---\n\n", string(yamlBytes))

	// Add content
	return frontmatter + cmd.Content, nil
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
