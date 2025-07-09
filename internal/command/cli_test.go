package command

import (
	"testing"
)

func TestNewMemoryCommand_Use(t *testing.T) {
	registry := createTestRegistry()
	cmd := NewMemoryCommand(registry)
	if cmd == nil {
		t.Fatal("expected NewMemoryCommand to return a command, got nil")
	}
	if cmd.Use != "memory [file]" {
		t.Errorf("expected Use 'memory [file]', got %q", cmd.Use)
	}
}

func TestNewCommandCommand_Use(t *testing.T) {
	registry := createTestRegistry()
	cmd := NewCommandCommand(registry)
	if cmd == nil {
		t.Fatal("expected NewCommandCommand to return a command, got nil")
	}
	if cmd.Use != "command [file_or_directory]" {
		t.Errorf("expected Use 'command [file_or_directory]', got %q", cmd.Use)
	}
}

func TestNewListCommand_Use(t *testing.T) {
	registry := createTestRegistry()
	cmd := NewListCommand(registry)
	if cmd == nil {
		t.Fatal("expected NewListCommand to return a command, got nil")
	}
	if cmd.Use != "list" {
		t.Errorf("expected Use 'list', got %q", cmd.Use)
	}
}
