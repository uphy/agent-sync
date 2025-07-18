package agent

import (
	"testing"
)

func TestClaude_ID_Name(t *testing.T) {
	c := &Claude{}
	if c.ID() != "claude" {
		t.Errorf("expected ID 'claude', got %q", c.ID())
	}
	if c.Name() != "Claude" {
		t.Errorf("expected Name 'Claude', got %q", c.Name())
	}
}

// ShouldConcatenate method has been removed as part of transition to path-based concatenation behavior
