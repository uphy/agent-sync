package parser

import (
	"bytes"
	"strings"

	"gopkg.in/yaml.v3"
)

// ParseFrontmatter parses YAML frontmatter from Markdown content
func ParseFrontmatter(content []byte) (map[string]interface{}, string, error) {
	// Check for frontmatter delimiter
	if !bytes.HasPrefix(content, []byte("---\n")) {
		return nil, string(content), nil
	}

	// Find end of frontmatter
	end := bytes.Index(content[4:], []byte("\n---"))
	if end == -1 {
		return nil, string(content), nil
	}
	end += 4 // adjust for initial 4 bytes we skipped

	frontmatterContent := content[4:end]
	body := string(content[end+5:]) // skip the closing "---\n"

	// Parse frontmatter YAML
	var frontmatter map[string]interface{}
	if err := yaml.Unmarshal(frontmatterContent, &frontmatter); err != nil {
		return nil, body, err
	}

	return frontmatter, strings.TrimSpace(body), nil
}
