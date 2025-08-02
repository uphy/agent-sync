package model

import (
	"fmt"

	"github.com/goccy/go-yaml"
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
