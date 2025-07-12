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

func TestParseCommand_DefaultFieldsEmpty(t *testing.T) {
	cmd, _ := ParseCommand("p", []byte{})
	if cmd.Slug != "" {
		t.Errorf("expected Slug empty, got %q", cmd.Slug)
	}
	if cmd.Name != "" {
		t.Errorf("expected Name empty, got %q", cmd.Name)
	}
	if cmd.RoleDefinition != "" {
		t.Errorf("expected RoleDefinition empty, got %q", cmd.RoleDefinition)
	}
	if cmd.WhenToUse != "" {
		t.Errorf("expected WhenToUse empty, got %q", cmd.WhenToUse)
	}
	if cmd.Groups != nil {
		t.Errorf("expected Groups empty, got %v", cmd.Groups)
	}
}
