package parser

import (
	"path/filepath"
	"strings"

	"github.com/user/agent-def/internal/model"
	"github.com/user/agent-def/internal/util"
)

// ExtractMarkdownContent extracts the content of a markdown file without frontmatter
func ExtractMarkdownContent(content []byte) (string, error) {
	_, remainingContent, err := ParseFrontmatter(content)
	if err != nil {
		return "", err
	}
	return remainingContent, nil
}

// ParseMarkdownFile parses a markdown file and returns its content without frontmatter
func ParseMarkdownFile(fs util.FileSystem, path string) (string, error) {
	content, err := fs.ReadFile(path)
	if err != nil {
		return "", err
	}

	return ExtractMarkdownContent(content)
}

// ParseCommand parses a command definition from a file
func ParseCommand(fs util.FileSystem, path string) (model.Command, error) {
	contentBytes, err := fs.ReadFile(path)
	if err != nil {
		return model.Command{}, util.WrapError(err, "failed to read command file")
	}

	// Parse frontmatter
	frontmatter, content, err := ParseFrontmatter(contentBytes)
	if err != nil {
		return model.Command{}, util.WrapError(err, "failed to parse frontmatter")
	}

	// Create command model
	cmd := model.Command{
		Path:    path,
		Content: content,
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
		for _, g := range groups {
			if group, ok := g.(string); ok {
				cmd.Groups = append(cmd.Groups, group)
			}
		}
	}

	// Validate required fields
	if cmd.Slug == "" {
		cmd.Slug = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	}
	if cmd.Name == "" {
		cmd.Name = strings.Title(cmd.Slug)
	}

	return cmd, nil
}

// ParseCommands parses all command files in a directory
func ParseCommands(fs util.FileSystem, dir string) ([]model.Command, error) {
	files, err := fs.ListFiles(dir, "*.md")
	if err != nil {
		return nil, util.WrapError(err, "failed to list markdown files")
	}

	var commands []model.Command
	for _, file := range files {
		cmd, err := ParseCommand(fs, file)
		if err != nil {
			return nil, util.WrapError(err, "failed to parse command")
		}
		commands = append(commands, cmd)
	}

	return commands, nil
}
