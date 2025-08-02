package parser

import (
	"io"

	"github.com/goccy/go-yaml"
	"github.com/uphy/agent-sync/internal/frontmatter"
	"github.com/uphy/agent-sync/internal/model"
	"github.com/uphy/agent-sync/internal/util"
)

// ParseModeWithPath parses a mode definition and records the original source path.
// The mode definition consists of YAML frontmatter and Markdown content.
func ParseModeWithPath(path string, r io.Reader) (*model.Mode, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// Parse frontmatter and content
	frontmatterData, body, err := frontmatter.Parse(content)
	if err != nil {
		return nil, util.WrapError(err, "failed to parse frontmatter")
	}

	mode := &model.Mode{
		Content: body,
	}

	// If we have frontmatter, try to parse the claude section
	if len(frontmatterData) > 0 {
		if claudeSection, ok := frontmatterData["claude"]; ok {
			// Convert to YAML bytes for unmarshaling
			yamlData, err := yaml.Marshal(claudeSection)
			if err != nil {
				return nil, util.WrapError(err, "failed to marshal claude frontmatter")
			}

			var claudeMode model.ClaudeCodeMode
			if err := yaml.Unmarshal(yamlData, &claudeMode); err != nil {
				return nil, util.WrapError(err, "failed to unmarshal claude frontmatter into ClaudeCodeMode")
			}
			mode.Claude = &claudeMode
		}
	}

	return mode, nil
}

// ParseMode parses a mode definition from the given reader without a path.
// Deprecated: Prefer ParseModeWithPath so Mode.Path is populated.
func ParseMode(r io.Reader) (*model.Mode, error) {
	return ParseModeWithPath("", r)
}
