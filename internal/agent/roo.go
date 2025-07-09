package agent

import (
	"encoding/json"
	"fmt"
	"strings"

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
	// Convert commands to JSON format for .roo/custom_modes.json
	type RooMode struct {
		Slug           string   `json:"slug"`
		Name           string   `json:"name"`
		RoleDefinition string   `json:"roleDefinition"`
		WhenToUse      string   `json:"whenToUse"`
		Groups         []string `json:"groups"`
		Content        string   `json:"content"`
	}

	var rooModes []RooMode
	for _, cmd := range commands {
		rooModes = append(rooModes, RooMode{
			Slug:           cmd.Slug,
			Name:           cmd.Name,
			RoleDefinition: cmd.RoleDefinition,
			WhenToUse:      cmd.WhenToUse,
			Groups:         cmd.Groups,
			Content:        cmd.Content,
		})
	}

	// Convert to JSON
	jsonBytes, err := json.MarshalIndent(rooModes, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal Roo commands to JSON: %w", err)
	}

	return string(jsonBytes), nil
}
