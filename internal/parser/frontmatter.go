package parser

import (
	"bytes"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/user/agent-def/internal/util"
)

// ParseFrontmatter parses YAML frontmatter from a markdown file
// It returns a map of frontmatter data, the remaining content, and any error
func ParseFrontmatter(content []byte) (map[string]interface{}, string, error) {
	// Check if content starts with '---'
	if !bytes.HasPrefix(content, []byte("---\n")) {
		// No frontmatter, return empty map and original content
		return map[string]interface{}{}, string(content), nil
	}

	// Find the end of frontmatter
	endIdx := bytes.Index(content[4:], []byte("\n---\n"))
	if endIdx == -1 {
		return nil, "", &util.ErrMalformedFrontmatter{
			Path:  "content",
			Cause: fmt.Errorf("missing closing frontmatter delimiter '---'"),
		}
	}
	endIdx += 4 // Adjust for the offset in the slice

	// Extract frontmatter and remaining content
	frontmatter := content[4:endIdx]
	remainingContent := content[endIdx+5:] // Skip past the closing '---\n'

	// Parse YAML frontmatter
	var data map[string]interface{}
	if err := yaml.Unmarshal(frontmatter, &data); err != nil {
		return nil, "", &util.ErrMalformedFrontmatter{
			Path:  "content",
			Cause: err,
		}
	}

	return data, string(remainingContent), nil
}

// ParseFrontmatterFromFile reads a file and parses its frontmatter
func ParseFrontmatterFromFile(fs util.FileSystem, path string) (map[string]interface{}, string, error) {
	contentBytes, err := fs.ReadFile(path)
	if err != nil {
		return nil, "", util.WrapError(err, "failed to read file for frontmatter parsing")
	}

	frontmatter, content, err := ParseFrontmatter(contentBytes)
	if err != nil {
		return nil, "", &util.ErrMalformedFrontmatter{
			Path:  path,
			Cause: err,
		}
	}

	return frontmatter, content, nil
}
