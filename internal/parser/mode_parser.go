package parser

import (
	"io"

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

	// Parse frontmatter and content (lenient)
	fm, body, err := frontmatter.Parse(content)
	if err != nil {
		return nil, util.WrapError(err, "failed to parse frontmatter")
	}

	mode := &model.Mode{
		Raw:     fm,
		Content: body,
	}

	// Extract common description from top-level frontmatter if present
	if fm != nil {
		if v, ok := fm["description"]; ok {
			if s, ok := v.(string); ok {
				mode.Description = s
			}
		}
	}

	return mode, nil
}

// ParseMode parses a mode definition from the given reader without a path.
// Deprecated: Prefer ParseModeWithPath so Mode.Path is populated.
func ParseMode(r io.Reader) (*model.Mode, error) {
	return ParseModeWithPath("", r)
}
