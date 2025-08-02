package model

import (
	"reflect"
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

func TestParseCommand_PreservesFrontmatterGenerically(t *testing.T) {
	yml := []byte(`---
description: top desc
roo:
  slug: deploy
  name: Deploy
  roleDefinition: Do deploy
  whenToUse: When needed
  groups: [dev, ops]
claude:
  allowed-tools: Bash(ls:*)
---
# Title
Body`)
	cmd, err := ParseCommand("commands/deploy.md", yml)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cmd.Description != "top desc" {
		t.Errorf("expected Description 'top desc', got %q", cmd.Description)
	}
	if cmd.Raw == nil {
		t.Fatalf("expected Raw to be set")
	}
	// spot-check nested keys exist
	if _, ok := cmd.Raw["roo"]; !ok {
		t.Errorf("expected roo section present in Raw")
	}
	if _, ok := cmd.Raw["claude"]; !ok {
		t.Errorf("expected claude section present in Raw")
	}
	if got, want := cmd.Path, "commands/deploy.md"; got != want {
		t.Errorf("Path mismatch: got %q want %q", got, want)
	}
	if want := "# Title\nBody"; cmd.Content != want {
		t.Errorf("content mismatch: got %q want %q", cmd.Content, want)
	}
}

func TestCommand_UnmarshalSection_ExtractsAgentSection(t *testing.T) {
	yml := []byte(`---
roo:
  slug: review
  name: Code Reviewer
  roleDefinition: Review code
---
content`)
	cmd, err := ParseCommand("x.md", yml)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	var roo struct {
		Slug           string `yaml:"slug"`
		Name           string `yaml:"name"`
		RoleDefinition string `yaml:"roleDefinition"`
	}
	if err := cmd.UnmarshalSection("roo", &roo); err != nil {
		t.Fatalf("unmarshal section error: %v", err)
	}
	if !reflect.DeepEqual(roo, struct {
		Slug           string `yaml:"slug"`
		Name           string `yaml:"name"`
		RoleDefinition string `yaml:"roleDefinition"`
	}{"review", "Code Reviewer", "Review code"}) {
		t.Errorf("unexpected roo section: %+v", roo)
	}
}
