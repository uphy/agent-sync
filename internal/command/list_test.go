package command

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/user/agent-def/internal/agent"
)

// TestRunList tests the list command output
func TestRunList(t *testing.T) {
	// Setup
	registry := createTestRegistry()

	// Create a buffer to capture output
	buf := &bytes.Buffer{}

	// Run the list command with the buffer as output
	err := RunList(buf, false, registry)

	// Verify
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	output := buf.String()

	// Check if output contains the agent names
	if !strings.Contains(output, "claude") || !strings.Contains(output, "roo") {
		t.Errorf("expected output to contain agent names, got: %s", output)
	}

	// Check the formatting of the output
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) < 3 { // header + 2 agents
		t.Errorf("expected at least 3 lines in output, got %d", len(lines))
	}

	// Check if it has a header line
	if !strings.Contains(lines[0], "NAME") || !strings.Contains(lines[0], "CAPABILITIES") {
		t.Errorf("expected first line to be header, got: %s", lines[0])
	}
}

// TestRunList_Verbose tests the list command with verbose output
func TestRunList_Verbose(t *testing.T) {
	// Setup
	registry := createTestRegistry()

	// Create buffers for both verbose and non-verbose output
	verboseBuf := &bytes.Buffer{}
	normalBuf := &bytes.Buffer{}

	// Run the list command with verbose and non-verbose flags
	err1 := RunList(verboseBuf, true, registry)
	err2 := RunList(normalBuf, false, registry)

	// Verify no errors
	if err1 != nil {
		t.Fatalf("expected no error for verbose mode, got %v", err1)
	}
	if err2 != nil {
		t.Fatalf("expected no error for non-verbose mode, got %v", err2)
	}

	verboseOutput := verboseBuf.String()
	normalOutput := normalBuf.String()

	// Check if the verbose output is the same or longer than non-verbose
	// When descriptions are long enough to be truncated in non-verbose mode
	if len(verboseOutput) < len(normalOutput) {
		t.Errorf("expected verbose output to be at least as long as normal output")
	}
}

// TestRunList_EmptyRegistry tests the list command with an empty registry
func TestRunList_EmptyRegistry(t *testing.T) {
	// Create an empty registry
	registry := agent.NewRegistry()

	// Create a buffer for output
	buf := &bytes.Buffer{}

	// Run the list command
	err := RunList(buf, false, registry)

	// Verify
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	output := buf.String()

	// Check if the output indicates no agents
	if !strings.Contains(output, "No agents registered") {
		t.Errorf("expected output to indicate no agents, got: %s", output)
	}
}

// TestNewListCommand tests the creation of the list command
func TestNewListCommand(t *testing.T) {
	registry := createTestRegistry()
	cmd := NewListCommand(registry)

	// Check the command properties
	if cmd.Use != "list" {
		t.Errorf("expected Use 'list', got %q", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("expected non-empty Short description")
	}

	if cmd.Long == "" {
		t.Error("expected non-empty Long description")
	}

	// Check the RunE function
	if cmd.RunE == nil {
		t.Fatal("expected RunE function to be set")
	}

	// Check flags
	verboseFlag := cmd.Flag("verbose")
	if verboseFlag == nil {
		t.Fatal("expected verbose flag to exist")
	}
	if verboseFlag.Shorthand != "v" {
		t.Errorf("expected verbose flag shorthand 'v', got %q", verboseFlag.Shorthand)
	}
}

// TestListCommand_Execute tests the execution of the list command
func TestListCommand_Execute(t *testing.T) {
	registry := createTestRegistry()
	cmd := NewListCommand(registry)

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Restore when test is done
	defer func() {
		os.Stdout = oldStdout
	}()

	// Set args (none needed for list command)
	cmd.SetArgs([]string{})

	// Execute the command
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Close writer to get all output
	w.Close()

	// Read captured output
	var outBuf bytes.Buffer
	io.Copy(&outBuf, r)
	output := outBuf.String()

	// Verify output
	if output == "" {
		t.Error("expected non-empty output")
	}
}
