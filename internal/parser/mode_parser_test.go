package parser

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/uphy/agent-sync/internal/model"
)

func TestParseMode_ValidFrontmatter(t *testing.T) {
	input := `---
claude:
  name: Test Mode
  description: Test description
  tools: [tool1, tool2]
---
# Test Mode
This is a test mode.
`
	r := bytes.NewBufferString(input)
	mode, err := ParseMode(r)
	require.NoError(t, err)

	// Raw should contain the "claude" map with the expected keys
	claude, ok := mode.Raw["claude"].(map[string]interface{})
	require.True(t, ok, "claude section should be a map[string]interface{}")
	assert.Equal(t, "Test Mode", claude["name"])
	assert.Equal(t, "Test description", claude["description"])
	assert.ElementsMatch(t, []interface{}{"tool1", "tool2"}, claude["tools"].([]interface{}))

	assert.Equal(t, "# Test Mode\nThis is a test mode.\n", mode.Content)
}

func TestParseMode_NoFrontmatter(t *testing.T) {
	input := `# Test Mode
This is a test mode without frontmatter.
`
	r := bytes.NewBufferString(input)
	mode, err := ParseMode(r)
	require.NoError(t, err)
	assert.Equal(t, &model.Mode{
		Raw:     map[string]any{},
		Content: input,
	}, mode)
}

func TestParseMode_InvalidFrontmatter(t *testing.T) {
	input := `---
invalid: yaml: :
---
# Test Mode
`
	r := bytes.NewBufferString(input)
	_, err := ParseMode(r)
	require.Error(t, err)
}

func TestParseMode_Empty(t *testing.T) {
	input := ``
	r := bytes.NewBufferString(input)
	mode, err := ParseMode(r)
	require.NoError(t, err)
	assert.Equal(t, &model.Mode{
		Raw:     map[string]any{},
		Content: "",
	}, mode)
}
