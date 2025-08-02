package agent

import (
	"fmt"
	"path/filepath"

	"github.com/uphy/agent-sync/internal/frontmatter"
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

	clean := func(s string) string {
		if s != "" && s[0] == '\n' {
			return s[1:]
		}
		return s
	}

	config := RooConfig{
		CustomModes: make([]CustomMode, 0, len(commands)),
	}

	for i, cmd := range commands {
		// Extract 'roo' section from generic frontmatter
		var section struct {
			Slug           string `yaml:"slug"`
			Name           string `yaml:"name"`
			Description    string `yaml:"description"`
			RoleDefinition string `yaml:"roleDefinition"`
			WhenToUse      string `yaml:"whenToUse"`
			Groups         any    `yaml:"groups"`
		}
		if cmd.Raw != nil {
			_ = cmd.UnmarshalSection("roo", &section)
		}

		// Validation for required fields
		if section.Slug == "" {
			return "", fmt.Errorf("command at index %d missing required field 'slug'", i)
		}
		if section.Name == "" {
			return "", fmt.Errorf("command at index %d missing required field 'name'", i)
		}
		if section.RoleDefinition == "" {
			return "", fmt.Errorf("command at index %d missing required field 'roleDefinition'", i)
		}

		// Fallback description to common one when empty
		if section.Description == "" && cmd.Description != "" {
			section.Description = cmd.Description
		}

		cm := CustomMode{
			Slug:               section.Slug,
			Name:               section.Name,
			Description:        clean(section.Description),
			RoleDefinition:     clean(section.RoleDefinition),
			WhenToUse:          clean(section.WhenToUse),
			Groups:             section.Groups,
			CustomInstructions: clean(cmd.Content),
		}
		if cm.Groups == nil {
			cm.Groups = []string{}
		}
		config.CustomModes = append(config.CustomModes, cm)
	}

	// Marshal the struct to YAML using shared helper (no fences for Roo)
	yml, err := frontmatter.RenderYAML(config)
	if err != nil {
		return "", fmt.Errorf("failed to marshal Roo commands to YAML: %w", err)
	}
	return yml, nil
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

	// Marshal the slice into a single YAML string using shared helper (no fences)
	yml, err := frontmatter.RenderYAML(out)
	if err != nil {
		return "", fmt.Errorf("failed to marshal Roo modes to YAML: %w", err)
	}
	return yml, nil
}

// ModePath returns the default path for Roo agent mode files
func (r *Roo) ModePath(userScope bool) string {
	return ""
}
