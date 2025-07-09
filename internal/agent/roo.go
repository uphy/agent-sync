package agent

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/user/agent-def/internal/model"
	"gopkg.in/yaml.v3"
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
	var node yaml.Node
	err := node.Encode(config)
	if err != nil {
		return "", fmt.Errorf("failed to encode Roo commands to YAML node: %w", err)
	}

	// Find all customInstructions fields and set them to literal style
	// The document node contains a single mapping node at Content[0]
	for _, rootMap := range node.Content {
		// Find the customModes sequence node
		for i := 0; i < len(rootMap.Content); i += 2 {
			if rootMap.Content[i].Value == "customModes" {
				customModes := rootMap.Content[i+1]

				// For each mode in customModes
				for _, modeNode := range customModes.Content {
					// Each mode node is a mapping node with key-value pairs
					for j := 0; j < len(modeNode.Content); j += 2 {
						if modeNode.Content[j].Value == "customInstructions" {
							// Set literal style for customInstructions values
							valueNode := modeNode.Content[j+1]
							if strings.Contains(valueNode.Value, "\n") {
								valueNode.Style = yaml.LiteralStyle
							}
						}
					}
				}
				break
			}
		}
	}

	// Marshal the modified node to YAML
	var buf strings.Builder
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)

	if err := encoder.Encode(&node); err != nil {
		return "", fmt.Errorf("failed to marshal Roo commands to YAML: %w", err)
	}
	encoder.Close()

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
