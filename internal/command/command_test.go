package command

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/agent-def/internal/util"
)

// TestRunCommand_SingleFile tests processing a single command file
func TestRunCommand_SingleFile(t *testing.T) {
	// Setup
	registry := createTestRegistry()

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Restore when test is done
	defer func() {
		os.Stdout = oldStdout
	}()

	// Run the command with a single file
	err := RunCommand("testdata/test_command.md", "claude", "testdata", registry)

	// Close writer to get all output
	w.Close()

	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verify
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check if output contains the command content (formatted by Claude agent)
	if !strings.Contains(output, "test-command") {
		t.Errorf("expected output to contain command name, got: %s", output)
	}
}

// TestRunCommand_Directory tests processing a directory of command files
func TestRunCommand_Directory(t *testing.T) {
	// Setup
	registry := createTestRegistry()

	// Create a temporary directory for multiple command files
	tempDir, err := os.MkdirTemp("", "command-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test command files
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

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Restore when test is done
	defer func() {
		os.Stdout = oldStdout
	}()

	// Run the command with the directory
	err = RunCommand(tempDir, "claude", tempDir, registry)

	// Close writer to get all output
	w.Close()

	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verify
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check if output contains both commands (formatted by Claude agent)
	if !strings.Contains(output, "cmd1") || !strings.Contains(output, "cmd2") {
		t.Errorf("expected output to contain both command names, got: %s", output)
	}
}

// TestRunCommand_InvalidAgent tests the command command with an invalid agent type
func TestRunCommand_InvalidAgent(t *testing.T) {
	// Setup
	registry := createTestRegistry()

	// Run the command with invalid agent type
	err := RunCommand("testdata/test_command.md", "invalid-agent", "testdata", registry)

	// Verify
	if err == nil {
		t.Fatal("expected error for invalid agent type, got nil")
	}

	// Check the error type
	if _, ok := err.(*util.ErrInvalidAgent); !ok {
		t.Errorf("expected ErrInvalidAgent error, got %T: %v", err, err)
	}
}

// TestRunCommand_FileNotFound tests the command command with a non-existent file
func TestRunCommand_FileNotFound(t *testing.T) {
	// Setup
	registry := createTestRegistry()

	// Run the command with non-existent file
	err := RunCommand("testdata/non-existent.md", "claude", "testdata", registry)

	// Verify
	if err == nil {
		t.Fatal("expected error for non-existent file, got nil")
	}
}

// TestNewCommandCommand_Flags tests the flags of the command command
func TestNewCommandCommand_Flags(t *testing.T) {
	registry := createTestRegistry()
	cmd := NewCommandCommand(registry)

	// Check agent flag
	agentFlag := cmd.Flag("agent")
	if agentFlag == nil {
		t.Fatal("expected agent flag to exist")
	}
	if agentFlag.Shorthand != "a" {
		t.Errorf("expected agent flag shorthand 'a', got %q", agentFlag.Shorthand)
	}
	if agentFlag.DefValue != "claude" {
		t.Errorf("expected agent flag default value 'claude', got %q", agentFlag.DefValue)
	}

	// Check base flag
	baseFlag := cmd.Flag("base")
	if baseFlag == nil {
		t.Fatal("expected base flag to exist")
	}
	if baseFlag.Shorthand != "b" {
		t.Errorf("expected base flag shorthand 'b', got %q", baseFlag.Shorthand)
	}
	if baseFlag.DefValue != "" {
		t.Errorf("expected base flag default value '', got %q", baseFlag.DefValue)
	}
}

// TestCommandCommand_RunE tests the RunE function of the command command
func TestCommandCommand_RunE(t *testing.T) {
	registry := createTestRegistry()
	cmd := NewCommandCommand(registry)

	// The RunE function requires args, so we'll set them
	cmd.SetArgs([]string{"testdata/test_command.md"})

	// We're not going to actually execute it, just check that it's set correctly
	if cmd.RunE == nil {
		t.Fatal("expected RunE function to be set")
	}

	// Check args validation
	if cmd.Args == nil {
		t.Fatal("expected Args validator to be set")
	}

	// Test Args validator with wrong number of args
	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("expected error for wrong number of args, got nil")
	}

	err = cmd.Args(cmd, []string{"testdata/test_command.md", "extra"})
	if err == nil {
		t.Error("expected error for wrong number of args, got nil")
	}

	// Test Args validator with correct number of args
	err = cmd.Args(cmd, []string{"testdata/test_command.md"})
	if err != nil {
		t.Errorf("expected no error for correct number of args, got %v", err)
	}
}
