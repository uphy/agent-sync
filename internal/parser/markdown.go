package parser

import (
	"path/filepath"
	"strings"

	"github.com/user/agent-def/internal/model"
	"github.com/user/agent-def/internal/util"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

// ParseCommandFromContent parses a command definition from provided content
func ParseCommandFromContent(path string, contentBytes []byte) (model.Command, error) {
	// Parse frontmatter
	frontmatter, content, err := ParseFrontmatter(contentBytes)
	if err != nil {
		return model.Command{}, util.WrapError(err, "failed to parse frontmatter")
	}
	content = strings.Trim(content, "\n\t ")

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
		cmd.Groups = groups
	}

	// Validate required fields
	if cmd.Slug == "" {
		cmd.Slug = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	}
	if cmd.Name == "" {
		cmd.Name = cases.Title(language.Und).String(cmd.Slug)
	}

	return cmd, nil
}

// ParseCommand parses a command definition from a file
func ParseCommand(fs util.FileSystem, path string) (model.Command, error) {
	contentBytes, err := fs.ReadFile(path)
	if err != nil {
		return model.Command{}, util.WrapError(err, "failed to read command file")
	}

	return ParseCommandFromContent(path, contentBytes)
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
