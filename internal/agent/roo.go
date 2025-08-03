package agent

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

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

// RooSlashMeta is a simple struct for Roo-specific metadata.
type RooSlashMeta struct {
	// Description: Optional Roo-specific description that overrides top-level Command.Description when provided
	Description  string `yaml:"description,omitempty"`
	ArgumentHint string `yaml:"argument-hint,omitempty"`
}

// FormatCommand processes command definitions for Roo agent (slash command markdown)
func (r *Roo) FormatCommand(commands []model.Command) (string, error) {
	// Concatenate multiple command files with a blank line between each output
	var outputs []string
	for _, cmd := range commands {
		var meta RooSlashMeta
		// Populate from roo section. Ignore error if section is missing.
		_ = cmd.UnmarshalSection("roo", &meta) // roo: argument-hint

		body := strings.TrimLeft(cmd.Content, "\n")

		// Determine description with Roo-specific override when provided.
		// Priority: roo.description > top-level cmd.Description
		if meta.Description == "" {
			meta.Description = cmd.Description
		}
		// Validation: require at least one description source
		if meta.Description == "" {
			return "", fmt.Errorf("Roo slash command requires a description: provide either top-level 'description' or 'roo.description'")
		}
		yamlWithFences, err := frontmatter.Wrap(meta)
		if err != nil {
			return "", fmt.Errorf("failed to render Roo slash command frontmatter: %w", err)
		}
		outputs = append(outputs, yamlWithFences+body)
	}
	return strings.Join(outputs, "\n\n"), nil
}

// MemoryPath returns the default path for Roo agent memory files
func (r *Roo) MemoryPath(userScope bool) string {
	return ".roo/rules/"
}

// CommandPath returns the default path for Roo agent command files (directory, non-concatenated)
func (r *Roo) CommandPath(userScope bool) string {
	return ".roo/commands/"
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
// Project scope aggregates into a single file ".roomodes"
// User scope aggregates into VS Code globalStorage custom_modes.yaml under the user's home directory
func (r *Roo) ModePath(userScope bool) string {
	if !userScope {
		// Project scope: single aggregated file in project root
		return ".roomodes"
	}
	// User scope: platform-specific VS Code globalStorage path (relative to home)
	// macOS: ~/Library/Application Support/Code/User/globalStorage/rooveterinaryinc.roo-cline/settings/custom_modes.yaml
	// Linux: ~/.config/Code/User/globalStorage/rooveterinaryinc.roo-cline/settings/custom_modes.yaml
	// Windows: %APPDATA%\Code\User\globalStorage\rooveterinaryinc.roo-cline\settings\custom_modes.yaml
	switch os := runtime.GOOS; os {
	case "darwin":
		return filepath.Join("Library", "Application Support", "Code", "User", "globalStorage", "rooveterinaryinc.roo-cline", "settings", "custom_modes.yaml")
	case "linux":
		return filepath.Join(".config", "Code", "User", "globalStorage", "rooveterinaryinc.roo-cline", "settings", "custom_modes.yaml")
	case "windows":
		// Return a relative path under APPDATA; pipeline will resolve against home for user scope
		return filepath.Join("AppData", "Roaming", "Code", "User", "globalStorage", "rooveterinaryinc.roo-cline", "settings", "custom_modes.yaml")
	default:
		// Fallback to Linux-like path
		return filepath.Join(".config", "Code", "User", "globalStorage", "rooveterinaryinc.roo-cline", "settings", "custom_modes.yaml")
	}
}

// Legacy compatibility for slash commands removed: only top-level 'description' and 'roo.argument-hint' are supported.
