package model

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/user/agent-def/internal/frontmatter"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Command represents a command definition
type Command struct {
	// Roo contains the command properties from frontmatter
	Roo Roo `yaml:"roo"`

	// Claude contains Claude-specific properties from frontmatter
	Claude Claude `yaml:"claude"`

	// Copilot contains Copilot-specific properties from frontmatter
	Copilot Copilot `yaml:"copilot"`

	// Content is the main content of the command (not in frontmatter)
	Content string `yaml:"-"`

	// Path is the original file path (not in frontmatter)
	Path string `yaml:"-"`
}

// Roo contains the command properties from frontmatter
type Roo struct {
	// Slug is the command identifier
	Slug string `yaml:"slug"`

	// Name is the display name of the command
	Name string `yaml:"name"`

	// RoleDefinition describes what the command does
	RoleDefinition string `yaml:"roleDefinition"`

	// WhenToUse describes when to use the command
	WhenToUse string `yaml:"whenToUse"`

	// Groups are permission groups for the command
	Groups any `yaml:"groups"`
}

// Claude contains the Claude-specific properties from frontmatter
type Claude struct {
	// Description is a brief description of the command
	Description string `yaml:"description"`

	// AllowedTools lists tools the command can use
	AllowedTools string `yaml:"allowed-tools"`
}

// Copilot contains the Copilot-specific properties from frontmatter
type Copilot struct {
	// Mode is the chat mode (ask, edit, agent)
	Mode string `yaml:"mode,omitempty"`

	// Model is the AI model to use
	Model string `yaml:"model,omitempty"`

	// Tools are the available tools
	Tools []string `yaml:"tools,omitempty"`

	// Description is the prompt description
	Description string `yaml:"description,omitempty"`
}

// ParseCommand parses a command definition from file content
func ParseCommand(path string, content []byte) (*Command, error) {
	// Parse frontmatter using frontmatter package
	fm, textContent, err := frontmatter.Parse(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	// Standardize with parser package by trimming whitespace
	textContent = strings.Trim(textContent, "\n\t ")

	// Create command model
	cmd := &Command{
		Path:    path,
		Content: textContent,
	}

	// Parse YAML directly into Command struct
	if len(fm) > 0 {
		yamlBytes, err := yaml.Marshal(fm)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal frontmatter: %w", err)
		}

		if err := yaml.Unmarshal(yamlBytes, cmd); err != nil {
			return nil, fmt.Errorf("failed to unmarshal command: %w", err)
		}
	}

	// Set defaults and validate
	cmd.setDefaults()
	if err := cmd.validate(); err != nil {
		return nil, fmt.Errorf("command validation failed: %w", err)
	}

	return cmd, nil
}

// setDefaults sets default values for missing Command fields
func (c *Command) setDefaults() {
	if c.Roo.Slug == "" {
		c.Roo.Slug = strings.TrimSuffix(filepath.Base(c.Path), filepath.Ext(c.Path))
	}
	if c.Roo.Name == "" {
		c.Roo.Name = cases.Title(language.Und).String(c.Roo.Slug)
	}
}

// validate validates that the command has all required fields
func (c *Command) validate() error {
	// Minimum validation - ensure required fields have values
	if c.Roo.Slug == "" {
		return fmt.Errorf("command slug is required")
	}

	if c.Roo.Name == "" {
		return fmt.Errorf("command name is required")
	}

	return nil
}
