package agent

import (
	"fmt"
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

func TestClaude_ShouldConcatenate(t *testing.T) {
	c := &Claude{}

	tests := []struct {
		taskType string
		expected bool
	}{
		{"memory", true},
		{"command", false},
		{"unknown", false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("TaskType=%s", tt.taskType), func(t *testing.T) {
			result := c.ShouldConcatenate(tt.taskType)
			if result != tt.expected {
				t.Errorf("Claude.ShouldConcatenate(%q) = %v, want %v", tt.taskType, result, tt.expected)
			}
		})
	}
}
