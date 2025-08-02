package model

import (
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/uphy/agent-sync/internal/frontmatter"
)

// Mode represents a subagent definition in an agent-agnostic way.
// Frontmatter keys are preserved in Raw.
// Description is a common, top-level field used as a fallback by agents.
type Mode struct {
	// Description is parsed from top-level frontmatter key "description".
	// Agent-specific descriptions override this during formatting if present.
	Description string `yaml:"-"`
	// Raw contains the full frontmatter as a generic map. Agent-specific
	// keys like "claude", "roo", "cline", "copilot" remain here.
	Raw map[string]any `yaml:"-"`
	// Content is the markdown body after frontmatter.
	Content string
}

// UnmarshalSection marshals a sub-map from m.Raw[key] to YAML and unmarshals into out.
// This centralizes section extraction for all agents (claude, roo, etc).
// Note: Go methods cannot be generic, so out is an interface{}.
func (m *Mode) UnmarshalSection(key string, out interface{}) error {
	if m.Raw == nil {
		return fmt.Errorf("frontmatter is nil")
	}
	sec, ok := m.Raw[key]
	if !ok {
		return fmt.Errorf("frontmatter missing section %q", key)
	}
	b, err := yaml.Marshal(sec)
	if err != nil {
		return fmt.Errorf("failed to marshal section %q: %w", key, err)
	}
	if err := yaml.Unmarshal(b, out); err != nil {
		return fmt.Errorf("failed to unmarshal section %q: %w", key, err)
	}
	return nil
}

// ParseMode parses a mode definition from file content into an agent-agnostic Mode.
// It mirrors ParseCommand behavior for frontmatter handling, but preserves body
// content exactly (no trimming) to maintain compatibility with existing tests.
func ParseMode(path string, content []byte) (*Mode, error) {
	// Parse frontmatter using frontmatter package (lenient)
	fm, body, err := frontmatter.Parse(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	mode := &Mode{
		Content: body,
		Raw:     map[string]any{},
	}

	// Preserve frontmatter generically by deep-converting into map[string]any
	if fm != nil {
		yamlBytes, err := yaml.Marshal(fm)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal frontmatter: %w", err)
		}
		tmp := map[string]any{}
		if err := yaml.Unmarshal(yamlBytes, &tmp); err != nil {
			return nil, fmt.Errorf("failed to unmarshal frontmatter to raw map: %w", err)
		}
		mode.Raw = tmp

		// Extract common description from top-level if present
		if v, ok := tmp["description"]; ok {
			if s, ok := v.(string); ok {
				mode.Description = s
			}
		}
	}

	// Minimal validation (agent-specific checks happen during formatting)
	if err := mode.validate(); err != nil {
		return nil, fmt.Errorf("mode validation failed: %w", err)
	}

	return mode, nil
}

// validate performs minimal, agent-agnostic validation.
func (m *Mode) validate() error {
	// Currently no required generic fields.
	return nil
}
