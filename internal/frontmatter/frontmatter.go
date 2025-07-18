package frontmatter

import (
	"bytes"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/uphy/agent-def/internal/util"
)

// ExtractFrontmatter extracts the raw YAML frontmatter bytes from content
// It returns the raw frontmatter bytes, the remaining content bytes, and any error
func ExtractFrontmatter(content []byte) ([]byte, []byte, error) {
	// Check if content starts with '---'
	if !bytes.HasPrefix(content, []byte("---\n")) {
		// No frontmatter, return empty byte slice and original content
		return []byte{}, content, nil
	}

	// Find the end of frontmatter
	endIdx := bytes.Index(content[4:], []byte("\n---\n"))
	if endIdx == -1 {
		return nil, nil, &util.ErrMalformedFrontmatter{
			Path:  "content",
			Cause: fmt.Errorf("missing closing frontmatter delimiter '---'"),
		}
	}
	endIdx += 4 // Adjust for the offset in the slice

	// Extract frontmatter and remaining content
	frontmatterBytes := content[4:endIdx]
	remainingContent := content[endIdx+5:] // Skip past the closing '---\n'

	return frontmatterBytes, remainingContent, nil
}

// Parse parses YAML frontmatter from a markdown file
// It returns a map of frontmatter data, the remaining content, and any error
func Parse(content []byte) (map[string]interface{}, string, error) {
	// Extract the frontmatter bytes using the new function
	frontmatterBytes, remainingContent, err := ExtractFrontmatter(content)
	if err != nil {
		return nil, "", err
	}

	// If no frontmatter was found, return an empty map and the original content
	if len(frontmatterBytes) == 0 {
		return map[string]interface{}{}, string(content), nil
	}

	// Parse YAML frontmatter
	var data map[string]interface{}
	if err := yaml.Unmarshal(frontmatterBytes, &data); err != nil {
		return nil, "", &util.ErrMalformedFrontmatter{
			Path:  "content",
			Cause: err,
		}
	}

	return data, string(remainingContent), nil
}

// ParseFromFile reads a file and parses its frontmatter
func ParseFromFile(fs util.FileSystem, path string) (map[string]interface{}, string, error) {
	contentBytes, err := fs.ReadFile(path)
	if err != nil {
		return nil, "", util.WrapError(err, "failed to read file for frontmatter parsing")
	}

	frontmatter, content, err := Parse(contentBytes)
	if err != nil {
		return nil, "", &util.ErrMalformedFrontmatter{
			Path:  path,
			Cause: err,
		}
	}

	return frontmatter, content, nil
}
