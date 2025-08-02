package agent

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/uphy/agent-sync/internal/model"
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
	if filepath.IsAbs(path) {
		return path
	}
	return "@/" + path
}

// FormatMCP formats an MCP command for Roo agent
func (r *Roo) FormatMCP(agent, command string, args ...string) string {
	return formatMCP(agent, command, args...)
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
	// Validate required fields for each command
	for i, cmd := range commands {
		if cmd.Roo.Slug == "" {
			return "", fmt.Errorf("command at index %d missing required field 'slug'", i)
		}
		if cmd.Roo.Name == "" {
			return "", fmt.Errorf("command at index %d missing required field 'name'", i)
		}
		if cmd.Roo.RoleDefinition == "" {
			return "", fmt.Errorf("command at index %d missing required field 'roleDefinition'", i)
		}
		// Groups is validated separately during struct creation
	}

	// Define structs that match our desired YAML output structure
	type CustomMode struct {
		Slug               string `yaml:"slug"`
		Name               string `yaml:"name"`
		Description        string `yaml:"description"`
		RoleDefinition     string `yaml:"roleDefinition"`
		WhenToUse          string `yaml:"whenToUse"`
		Groups             any    `yaml:"groups"`
		CustomInstructions string `yaml:"customInstructions"`
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
		customMode := CustomMode{
			Slug:               cmd.Roo.Slug,
			Name:               cmd.Roo.Name,
			Description:        cleanContent(cmd.Roo.Description),
			RoleDefinition:     cleanContent(cmd.Roo.RoleDefinition),
			WhenToUse:          cleanContent(cmd.Roo.WhenToUse),
			Groups:             cmd.Roo.Groups,
			CustomInstructions: cleanContent(cmd.Content),
		}
		if customMode.Groups == nil {
			customMode.Groups = []string{}
		}
		config.CustomModes[i] = customMode
	}

	// Encode to yaml.Node first
	// Marshal the struct to YAML
	var buf strings.Builder
	encoder := yaml.NewEncoder(
		&buf,
		yaml.Indent(2),
		yaml.UseSingleQuote(false),
		yaml.IndentSequence(true),
		yaml.UseLiteralStyleIfMultiline(true), // This ensures multi-line strings use literal style (|)
	)
	if err := encoder.Encode(config); err != nil {
		return "", fmt.Errorf("failed to marshal Roo commands to YAML: %w", err)
	}
	return buf.String(), nil
}

// MemoryPath returns the default path for Roo agent memory files
func (r *Roo) MemoryPath(userScope bool) string {
	return ".roo/rules/"
}

// CommandPath returns the default path for Roo agent command files
func (r *Roo) CommandPath(userScope bool) string {
	if userScope {
		return "Library/Application Support/Code/User/globalStorage/rooveterinaryinc.roo-cline/settings/custom_modes.yaml"
	}
	return ".roomodes"
}

// FormatMode processes mode definitions for Roo agent
// Roo supports combining multiple modes into a single YAML output, similar to FormatCommand.
func (r *Roo) FormatMode(modes []model.Mode) (string, error) {
	// Define output struct to marshal as YAML (list of Roo modes directly)
	type RooMode struct {
		Slug               string `yaml:"slug"`
		Name               string `yaml:"name"`
		Description        string `yaml:"description"`
		RoleDefinition     string `yaml:"roleDefinition"`
		WhenToUse          string `yaml:"whenToUse"`
		Groups             any    `yaml:"groups"` // preserve user-provided structure as-is
		CustomInstructions string `yaml:"customInstructions"`
	}

	clean := func(s string) string {
		if s != "" && s[0] == '\n' {
			return s[1:]
		}
		return s
	}

	out := make([]RooMode, 0, len(modes))

	for i, m := range modes {
		if m.Raw == nil {
			return "", fmt.Errorf("mode at index %d missing frontmatter", i)
		}
		// Unmarshal the 'roo' section directly using model.Mode helper
		var rm RooMode
		if err := m.UnmarshalSection("roo", &rm); err != nil {
			return "", fmt.Errorf("mode at index %d: %w", i, err)
		}

		// Apply fallbacks
		if rm.Description == "" && m.Description != "" {
			rm.Description = m.Description
		}
		rm.CustomInstructions = clean(m.Content)

		// Validation for required fields
		if rm.Slug == "" {
			return "", fmt.Errorf("mode at index %d missing required field 'slug'", i)
		}
		if rm.Name == "" {
			return "", fmt.Errorf("mode at index %d missing required field 'name'", i)
		}
		if rm.RoleDefinition == "" {
			return "", fmt.Errorf("mode at index %d missing required field 'roleDefinition'", i)
		}
		// Ensure non-nil groups (preserve user structure)
		if rm.Groups == nil {
			rm.Groups = []string{}
		}

		out = append(out, RooMode{
			Slug:               rm.Slug,
			Name:               rm.Name,
			Description:        clean(rm.Description),
			RoleDefinition:     clean(rm.RoleDefinition),
			WhenToUse:          clean(rm.WhenToUse),
			Groups:             rm.Groups,
			CustomInstructions: rm.CustomInstructions,
		})
	}

	// Marshal the slice into a single YAML string
	var buf strings.Builder
	enc := yaml.NewEncoder(
		&buf,
		yaml.Indent(2),
		yaml.UseSingleQuote(false),
		yaml.IndentSequence(true),
		yaml.UseLiteralStyleIfMultiline(true),
	)
	if err := enc.Encode(out); err != nil {
		return "", fmt.Errorf("failed to marshal Roo modes to YAML: %w", err)
	}
	return buf.String(), nil
}

// ModePath returns the default path for Roo agent mode files
func (r *Roo) ModePath(userScope bool) string {
	return ""
}
