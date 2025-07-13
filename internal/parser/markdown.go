package parser

import (
	"github.com/user/agent-def/internal/frontmatter"
	"github.com/user/agent-def/internal/model"
	"github.com/user/agent-def/internal/util"
)

// ExtractMarkdownContent extracts the content of a markdown file without frontmatter
func ExtractMarkdownContent(content []byte) (string, error) {
	_, remainingContent, err := frontmatter.Parse(content)
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
func ParseCommandFromContent(path string, contentBytes []byte) (*model.Command, error) {
	return model.ParseCommand(path, contentBytes)
}

// ParseCommand parses a command definition from a file
func ParseCommand(fs util.FileSystem, path string) (*model.Command, error) {
	contentBytes, err := fs.ReadFile(path)
	if err != nil {
		return nil, util.WrapError(err, "failed to read command file")
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
		commands = append(commands, *cmd)
	}

	return commands, nil
}
