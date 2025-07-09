package command

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func newRootCommand() *cobra.Command {
	registry := createTestRegistry()
	root := &cobra.Command{Use: "agent-def"}
	root.AddCommand(
		NewMemoryCommand(registry),
		NewCommandCommand(registry),
		NewListCommand(registry),
	)
	return root
}

func TestCLI_ListCommand_Output(t *testing.T) {
	root := newRootCommand()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"list"})
	if err := root.Execute(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "List supported agents") {
		t.Errorf("expected output to contain %q, got %q", "List supported agents", out)
	}
}

func TestCLI_MemoryCommand_Help(t *testing.T) {
	root := newRootCommand()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"memory", "--help"})
	if err := root.Execute(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Process memory context") {
		t.Errorf("expected help to contain %q, got %q", "Process memory context", out)
	}
}

func TestCLI_CommandCommand_Help(t *testing.T) {
	root := newRootCommand()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"command", "--help"})
	if err := root.Execute(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Process command definitions") {
		t.Errorf("expected help to contain %q, got %q", "Process command definitions", out)
	}
}

// Integration tests for end-to-end workflow

func TestIntegration_MemoryCommand_Claude(t *testing.T) {
	root := newRootCommand()
	buf := &bytes.Buffer{}
	root.SetOut(buf)

	// Run the memory command with the Claude agent
	root.SetArgs([]string{"memory", "testdata/memory_context.md", "--agent", "claude"})
	if err := root.Execute(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	out := buf.String()

	// Check Claude-specific formatting
	if !strings.Contains(out, "Test Memory Context") {
		t.Errorf("expected output to contain memory content")
	}

	// Check that included content is processed
	if !strings.Contains(out, "included content") {
		t.Errorf("expected output to contain included content")
	}

	// Check that references are properly formatted for Claude
	if !strings.Contains(out, "[参考:") {
		t.Errorf("expected output to contain reference markers in Claude format")
	}

	// Check references section
	if !strings.Contains(out, "## References") {
		t.Errorf("expected output to contain references section")
	}
}

func TestIntegration_MemoryCommand_Roo(t *testing.T) {
	root := newRootCommand()
	buf := &bytes.Buffer{}
	root.SetOut(buf)

	// Run the memory command with the Roo agent
	root.SetArgs([]string{"memory", "testdata/memory_context.md", "--agent", "roo"})
	if err := root.Execute(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	out := buf.String()

	// Check that memory content is included
	if !strings.Contains(out, "Test Memory Context") {
		t.Errorf("expected output to contain memory content")
	}

	// Check that included content is processed
	if !strings.Contains(out, "included content") {
		t.Errorf("expected output to contain included content")
	}

	// Check that references section is included - format will be different from Claude
	if !strings.Contains(out, "## References") {
		t.Errorf("expected output to contain references section")
	}
}

func TestIntegration_CommandCommand_Claude(t *testing.T) {
	root := newRootCommand()
	buf := &bytes.Buffer{}
	root.SetOut(buf)

	// Run the command command with the Claude agent
	root.SetArgs([]string{"command", "testdata/test_command.md", "--agent", "claude"})
	if err := root.Execute(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	out := buf.String()

	// Check that command content is included
	if !strings.Contains(out, "test-command") {
		t.Errorf("expected output to contain command name")
	}

	// Check that parameters are included
	if !strings.Contains(out, "param1") {
		t.Errorf("expected output to contain parameter information")
	}
}

func TestIntegration_CommandCommand_Roo(t *testing.T) {
	root := newRootCommand()
	buf := &bytes.Buffer{}
	root.SetOut(buf)

	// Run the command command with the Roo agent
	root.SetArgs([]string{"command", "testdata/test_command.md", "--agent", "roo"})
	if err := root.Execute(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	out := buf.String()

	// Check that command content is included
	if !strings.Contains(out, "test-command") {
		t.Errorf("expected output to contain command name")
	}

	// Check that parameters are included
	if !strings.Contains(out, "param1") {
		t.Errorf("expected output to contain parameter information")
	}
}

func TestIntegration_CommandCommand_Directory(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "command-integration-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create multiple command files
	cmd1 := `---
name: cmd1
description: Command 1
parameters:
  - name: param1
    type: string
    required: true
---
# Command 1
This is command 1.
`
	cmd2 := `---
name: cmd2
description: Command 2
parameters:
  - name: param2
    type: string
    required: true
---
# Command 2
This is command 2.
`

	cmd1Path := filepath.Join(tempDir, "cmd1.md")
	cmd2Path := filepath.Join(tempDir, "cmd2.md")

	if err := os.WriteFile(cmd1Path, []byte(cmd1), 0644); err != nil {
		t.Fatalf("failed to write cmd1.md: %v", err)
	}
	if err := os.WriteFile(cmd2Path, []byte(cmd2), 0644); err != nil {
		t.Fatalf("failed to write cmd2.md: %v", err)
	}

	root := newRootCommand()
	buf := &bytes.Buffer{}
	root.SetOut(buf)

	// Run the command command with the directory
	root.SetArgs([]string{"command", tempDir, "--agent", "claude"})
	if err := root.Execute(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	out := buf.String()

	// Check that both commands are included
	if !strings.Contains(out, "cmd1") || !strings.Contains(out, "cmd2") {
		t.Errorf("expected output to contain both command names")
	}

	// Check that parameters are included
	if !strings.Contains(out, "param1") || !strings.Contains(out, "param2") {
		t.Errorf("expected output to contain both parameter information")
	}
}

func TestIntegration_InvalidAgent(t *testing.T) {
	root := newRootCommand()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)

	// Run the memory command with an invalid agent
	root.SetArgs([]string{"memory", "testdata/memory_context.md", "--agent", "invalid"})
	err := root.Execute()

	// This should return an error
	if err == nil {
		t.Fatalf("expected error for invalid agent, got nil")
	}

	// Error message should mention invalid agent
	errMsg := err.Error()
	if !strings.Contains(strings.ToLower(errMsg), "invalid agent") {
		t.Errorf("expected error message to mention invalid agent, got: %s", errMsg)
	}
}

func TestIntegration_RecursiveTemplateProcessing(t *testing.T) {
	// Create a temporary directory with recursive templates
	tempDir, err := os.MkdirTemp("", "recursive-template-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create template files
	base := `---
title: Base Template
---
# Base Template

{{include "level1.md"}}

End of base template.
`
	level1 := `# Level 1 Include

This is level 1.

{{include "level2.md"}}

End of level 1.
`
	level2 := `## Level 2 Include

This is level 2.
`

	basePath := filepath.Join(tempDir, "base.md")
	level1Path := filepath.Join(tempDir, "level1.md")
	level2Path := filepath.Join(tempDir, "level2.md")

	if err := os.WriteFile(basePath, []byte(base), 0644); err != nil {
		t.Fatalf("failed to write base.md: %v", err)
	}
	if err := os.WriteFile(level1Path, []byte(level1), 0644); err != nil {
		t.Fatalf("failed to write level1.md: %v", err)
	}
	if err := os.WriteFile(level2Path, []byte(level2), 0644); err != nil {
		t.Fatalf("failed to write level2.md: %v", err)
	}

	root := newRootCommand()
	buf := &bytes.Buffer{}
	root.SetOut(buf)

	// Run the memory command with the base template
	root.SetArgs([]string{"memory", basePath, "--agent", "claude", "--base", tempDir})
	if err := root.Execute(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	out := buf.String()

	// Check that all levels are included
	if !strings.Contains(out, "Base Template") {
		t.Errorf("expected output to contain base template content")
	}
	if !strings.Contains(out, "Level 1 Include") {
		t.Errorf("expected output to contain level 1 content")
	}
	if !strings.Contains(out, "Level 2 Include") {
		t.Errorf("expected output to contain level 2 content")
	}
}
