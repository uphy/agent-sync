package agent

import (
	"fmt"
	"testing"
)

func TestRoo_ID_Name(t *testing.T) {
	r := &Roo{}
	if r.ID() != "roo" {
		t.Errorf("expected ID 'roo', got %q", r.ID())
	}
	if r.Name() != "Roo" {
		t.Errorf("expected Name 'Roo', got %q", r.Name())
	}
}

func TestRoo_ShouldConcatenate(t *testing.T) {
	r := &Roo{}

	tests := []struct {
		taskType string
		expected bool
	}{
		{"memory", false},
		{"command", true},
		{"unknown", false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("TaskType=%s", tt.taskType), func(t *testing.T) {
			result := r.ShouldConcatenate(tt.taskType)
			if result != tt.expected {
				t.Errorf("Roo.ShouldConcatenate(%q) = %v, want %v", tt.taskType, result, tt.expected)
			}
		})
	}
}
