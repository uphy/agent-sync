package model

import (
	"testing"
)

func TestParseCommand_SetsPathAndContent(t *testing.T) {
	inputPath := "path/to/file.md"
	inputContent := []byte("example content")
	cmd, err := ParseCommand(inputPath, inputContent)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cmd.Path != inputPath {
		t.Errorf("expected Path %q, got %q", inputPath, cmd.Path)
	}
	if cmd.Content != string(inputContent) {
		t.Errorf("expected Content %q, got %q", string(inputContent), cmd.Content)
	}
}

func TestParseCommand_DefaultFieldsFromPath(t *testing.T) {
	cmd, _ := ParseCommand("p", []byte{})

	// Now Slug is set from path and Name from Slug
	if cmd.Roo.Slug != "p" {
		t.Errorf("expected Slug to be 'p', got %q", cmd.Roo.Slug)
	}
	if cmd.Roo.Name != "P" {
		t.Errorf("expected Name to be 'P', got %q", cmd.Roo.Name)
	}
	if cmd.Roo.RoleDefinition != "" {
		t.Errorf("expected RoleDefinition empty, got %q", cmd.Roo.RoleDefinition)
	}
	if cmd.Roo.WhenToUse != "" {
		t.Errorf("expected WhenToUse empty, got %q", cmd.Roo.WhenToUse)
	}
	if cmd.Roo.Groups != nil {
		t.Errorf("expected Groups empty, got %v", cmd.Roo.Groups)
	}
}
