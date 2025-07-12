package agent

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/user/agent-def/internal/model"
)

// Roo implements the Roo-specific conversion logic
type Roo struct{}

// ID returns the unique identifier for Roo agent
func (r *Roo) ID() string {
	return "roo"
}

// Name returns the display name for Roo agent
func (r *Roo) Name() string {
	return "Roo"
}

// FormatFile converts a path to Roo's file reference format
func (r *Roo) FormatFile(path string) string {
	return "@/" + path
}

// FormatMCP formats an MCP command for Roo agent
func (r *Roo) FormatMCP(agent, command string, args ...string) string {
	return fmt.Sprintf("MCP tool (MCP Server: %s, Tool: %s, Arguments: %s)", agent, command, strings.Join(args, " "))
}

// FormatMemory processes a memory context for Roo agent
func (r *Roo) FormatMemory(content string) (string, error) {
	// Process the memory context for Roo
	// For now, we're simply returning the content as is
	// Any Roo-specific formatting rules could be applied here
	return content, nil
}

// FormatCommand processes command definitions for Roo agent
func (r *Roo) FormatCommand(commands []model.Command) (string, error) {
	// Define structs that match our desired YAML output structure
	type CustomMode struct {
		Slug               string   `yaml:"slug"`
		Name               string   `yaml:"name"`
		RoleDefinition     string   `yaml:"roleDefinition"`
		WhenToUse          string   `yaml:"whenToUse"`
		Groups             []string `yaml:"groups"`
		CustomInstructions string   `yaml:"customInstructions"`
	}

	type RooConfig struct {
		CustomModes []CustomMode `yaml:"customModes"`
	}

	// Clean leading newlines if present
	cleanContent := func(s string) string {
		if s != "" && s[0] == '\n' {
			return s[1:]
		}
		return s
	}

	// Create and populate the config struct
	config := RooConfig{
		CustomModes: make([]CustomMode, len(commands)),
	}

	for i, cmd := range commands {
		config.CustomModes[i] = CustomMode{
			Slug:               cmd.Slug,
			Name:               cmd.Name,
			RoleDefinition:     cleanContent(cmd.RoleDefinition),
			WhenToUse:          cleanContent(cmd.WhenToUse),
			Groups:             cmd.Groups,
			CustomInstructions: cleanContent(cmd.Content),
		}
	}

	// Encode to yaml.Node first
	// Marshal the struct to YAML
	yamlData, err := yaml.Marshal(config)
	if err != nil {
		return "", fmt.Errorf("failed to marshal Roo commands to YAML: %w", err)
	}

	// Since go-yaml doesn't have the same node manipulation capabilities as yaml.v3,
	// we'll use a different approach to preserve multi-line strings properly

	// Parse the YAML into a map for manipulation
	var parsedYaml map[string]interface{}
	if err := yaml.Unmarshal(yamlData, &parsedYaml); err != nil {
		return "", fmt.Errorf("failed to parse YAML for formatting: %w", err)
	}

	// Marshal the modified node to YAML
	// Re-marshal with proper formatting options
	var buf strings.Builder
	encoder := yaml.NewEncoder(
		&buf,
		yaml.Indent(2),
		yaml.UseSingleQuote(false),
		yaml.UseLiteralStyleIfMultiline(true), // This ensures multi-line strings use literal style (|)
	)

	if err := encoder.Encode(parsedYaml); err != nil {
		return "", fmt.Errorf("failed to marshal Roo commands to YAML: %w", err)
	}

	return buf.String(), nil
}

// DefaultMemoryPath determines the output path for Roo agent memory files
func (r *Roo) DefaultMemoryPath(outputBaseDir string, userScope bool, fileName string) (string, error) {
	ext := ".md"
	name := strings.TrimSuffix(fileName, ext)

	return filepath.Join(outputBaseDir, ".roo", "rules", name+ext), nil
}

// DefaultCommandPath determines the output path for Roo agent command files
func (r *Roo) DefaultCommandPath(outputBaseDir string, userScope bool, fileName string) (string, error) {
	if userScope {
		return filepath.Join(outputBaseDir, "Library", "Application Support", "Code", "User", "globalStorage", "rooveterinaryinc.roo-cline", "settings", "custom_modes.yaml"), nil
	}

	// Project scope
	return filepath.Join(outputBaseDir, ".roo", ".roomodes"), nil
}
