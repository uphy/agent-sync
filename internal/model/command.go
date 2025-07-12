package model

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/user/agent-def/internal/util"
)

// Command represents a command definition
type Command struct {
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

	// Content is the main content of the command
	Content string

	// Path is the original file path
	Path string
}

// Validate validates that the command has all required fields
func (c *Command) Validate() error {
	// Minimum validation - ensure required fields have values
	if c.Slug == "" {
		return fmt.Errorf("command slug is required")
	}

	if c.Name == "" {
		return fmt.Errorf("command name is required")
	}

	return nil
}

// parseFrontmatter extracts YAML frontmatter from markdown content
func parseFrontmatter(content []byte) (map[string]interface{}, string, error) {
	// Check if content starts with '---'
	if len(content) < 4 || string(content[:4]) != "---\n" {
		// No frontmatter, return empty map and original content
		return map[string]interface{}{}, string(content), nil
	}

	// Find the end of frontmatter
	endIdx := strings.Index(string(content[4:]), "\n---\n")
	if endIdx == -1 {
		return nil, "", fmt.Errorf("malformed frontmatter: missing closing '---'")
	}
	endIdx += 4 // Adjust for the offset in the slice

	// Extract frontmatter and remaining content
	frontmatter := content[4:endIdx]
	remainingContent := content[endIdx+5:] // Skip past the closing '---\n'

	// Parse YAML frontmatter
	var data map[string]interface{}
	if err := yaml.Unmarshal(frontmatter, &data); err != nil {
		return nil, "", fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	return data, string(remainingContent), nil
}

// ParseCommand parses a command definition from file content
func ParseCommand(path string, content []byte) (Command, error) {
	// Parse frontmatter
	frontmatter, textContent, err := parseFrontmatter(content)
	if err != nil {
		return Command{}, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	// Create command model
	cmd := Command{
		Path:    path,
		Content: textContent,
	}

	// Map frontmatter to command fields
	if slug, ok := frontmatter["slug"].(string); ok {
		cmd.Slug = slug
	}
	if name, ok := frontmatter["name"].(string); ok {
		cmd.Name = name
	}
	if role, ok := frontmatter["roleDefinition"].(string); ok {
		cmd.RoleDefinition = role
	}
	if when, ok := frontmatter["whenToUse"].(string); ok {
		cmd.WhenToUse = when
	}
	if groups, ok := frontmatter["groups"].([]interface{}); ok {
		cmd.Groups = groups
	}

	// Fill in defaults for missing fields
	if cmd.Slug == "" {
		cmd.Slug = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	}
	if cmd.Name == "" {
		cmd.Name = strings.Title(cmd.Slug)
	}

	return cmd, nil
}

// LoadCommandsFromDirectory loads all commands from a directory
func LoadCommandsFromDirectory(fs util.FileSystem, dir string) ([]Command, error) {
	files, err := fs.ListFiles(dir, "*.md")
	if err != nil {
		return nil, util.WrapError(err, "failed to list command files")
	}

	var commands []Command
	for _, file := range files {
		content, err := fs.ReadFile(file)
		if err != nil {
			return nil, util.WrapError(err, fmt.Sprintf("failed to read command file: %s", file))
		}

		cmd, err := ParseCommand(file, content)
		if err != nil {
			return nil, util.WrapError(err, fmt.Sprintf("failed to parse command: %s", file))
		}

		if err := cmd.Validate(); err != nil {
			return nil, util.WrapError(err, fmt.Sprintf("invalid command in %s", file))
		}

		commands = append(commands, cmd)
	}

	return commands, nil
}
