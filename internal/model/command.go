package model

import (
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/uphy/agent-sync/internal/frontmatter"
)

// Command represents a command definition in an agent-agnostic way.
// Frontmatter keys are preserved in Raw.
// Description is a common, top-level field used as a fallback by agents.
type Command struct {
	// Description is parsed from top-level frontmatter key "description".
	// Agent-specific descriptions override this during formatting if present.
	Description string `yaml:"-"`
	// Raw contains the full frontmatter as a generic map. Agent-specific
	// keys like "claude", "roo", "cline", "copilot" remain here.
	Raw map[string]any `yaml:"-"`
	// Content is the markdown body after frontmatter.
	Content string `yaml:"-"`
	// Path is the original file path (not in frontmatter)
	Path string `yaml:"-"`
}

// UnmarshalSection marshals a sub-map from c.Raw[key] to YAML and unmarshals into out.
// This mirrors model.Mode.UnmarshalSection and centralizes section extraction for all agents.
func (c *Command) UnmarshalSection(key string, out interface{}) error {
	if c.Raw == nil {
		return fmt.Errorf("frontmatter is nil")
	}
	sec, ok := c.Raw[key]
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

// ParseCommand parses a command definition from file content into an agent-agnostic Command.
func ParseCommand(path string, content []byte) (*Command, error) {
	// Parse frontmatter using frontmatter package
	fm, textContent, err := frontmatter.Parse(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	// Normalize content by trimming surrounding whitespace
	textContent = strings.Trim(textContent, "\n\t ")

	cmd := &Command{
		Path:    path,
		Content: textContent,
		Raw:     map[string]any{},
	}

	// Preserve frontmatter generically.
	// Extract a top-level description if present.
	if len(fm) > 0 {
		// marshal/unmarshal through YAML to deep-convert into map[string]any
		yamlBytes, err := yaml.Marshal(fm)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal frontmatter: %w", err)
		}
		tmp := map[string]any{}
		if err := yaml.Unmarshal(yamlBytes, &tmp); err != nil {
			return nil, fmt.Errorf("failed to unmarshal frontmatter to raw map: %w", err)
		}
		cmd.Raw = tmp

		// pick common description from top-level if present
		if v, ok := tmp["description"]; ok {
			if s, ok := v.(string); ok {
				cmd.Description = s
			}
		}
	}

	// Only generic validation here; agent-specific checks happen in agent implementations.
	if err := cmd.validate(); err != nil {
		return nil, fmt.Errorf("command validation failed: %w", err)
	}
	return cmd, nil
}

// validate performs minimal, agent-agnostic validation.
// Agent-specific required fields are validated in each agent's formatter.
func (c *Command) validate() error {
	if c.Content == "" {
		// Commands are allowed to have empty body; keep this non-fatal.
		// Return nil to avoid breaking existing flows.
		return nil
	}
	return nil
}
