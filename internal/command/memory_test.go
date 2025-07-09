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

// MockFileSystem implements util.FileSystem for testing
type MockFileSystem struct {
	files map[string][]byte
}

func NewMockFileSystem() *MockFileSystem {
	return &MockFileSystem{
		files: make(map[string][]byte),
	}
}

func (fs *MockFileSystem) ReadFile(path string) ([]byte, error) {
	if content, ok := fs.files[path]; ok {
		return content, nil
	}
	return nil, os.ErrNotExist
}

func (fs *MockFileSystem) WriteFile(path string, data []byte) error {
	fs.files[path] = data
	return nil
}

func (fs *MockFileSystem) FileExists(path string) bool {
	_, ok := fs.files[path]
	return ok
}

func (fs *MockFileSystem) AddFile(path string, content []byte) {
	fs.files[path] = content
}

func (fs *MockFileSystem) ListFiles(dir string, pattern string) ([]string, error) {
	var files []string
	for filePath := range fs.files {
		if strings.HasPrefix(filePath, dir) {
			// Check if it matches the pattern
			matched, err := filepath.Match(pattern, filepath.Base(filePath))
			if err != nil {
				return nil, err
			}
			if matched {
				files = append(files, filePath)
			}
		}
	}
	return files, nil
}

func (fs *MockFileSystem) ResolvePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join("/base", path)
}

type mockDirEntry struct {
	name  string
	isDir bool
}

func (e *mockDirEntry) Name() string               { return e.name }
func (e *mockDirEntry) IsDir() bool                { return e.isDir }
func (e *mockDirEntry) Type() os.FileMode          { return 0 }
func (e *mockDirEntry) Info() (os.FileInfo, error) { return nil, nil }

// TestRunMemory_ValidInput tests the memory command with valid input
func TestRunMemory_ValidInput(t *testing.T) {
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

	// Run the command using test data files
	err := RunMemory("testdata/memory_context.md", "claude", "testdata", registry)

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

	// Check if output contains the memory content (formatted by Claude agent)
	if !strings.Contains(output, "Test Memory Context") {
		t.Errorf("expected output to contain memory content, got: %s", output)
	}
}

// TestRunMemory_InvalidAgent tests the memory command with an invalid agent type
func TestRunMemory_InvalidAgent(t *testing.T) {
	// Setup
	registry := createTestRegistry()

	// Run the command with invalid agent type
	err := RunMemory("testdata/memory_context.md", "invalid-agent", "testdata", registry)

	// Verify
	if err == nil {
		t.Fatal("expected error for invalid agent type, got nil")
	}

	// Check the error type
	if _, ok := err.(*util.ErrInvalidAgent); !ok {
		t.Errorf("expected ErrInvalidAgent error, got %T: %v", err, err)
	}
}

// TestRunMemory_FileNotFound tests the memory command with a non-existent file
func TestRunMemory_FileNotFound(t *testing.T) {
	// Setup
	registry := createTestRegistry()

	// Run the command with non-existent file
	err := RunMemory("testdata/non-existent.md", "claude", "testdata", registry)

	// Verify
	if err == nil {
		t.Fatal("expected error for non-existent file, got nil")
	}
}

// TestMemoryFileResolver_Interface tests the MemoryFileResolver implementation
func TestMemoryFileResolver_Interface(t *testing.T) {
	// Setup
	fs := &util.RealFileSystem{}

	resolver := &MemoryFileResolver{
		FS:       fs,
		BasePath: "testdata",
	}

	// Test Exists for test files we've created
	if !resolver.Exists("testdata/memory_context.md") {
		t.Error("expected memory_context.md file to exist")
	}

	if resolver.Exists("testdata/non-existent.md") {
		t.Error("expected non-existent.md not to exist")
	}

	// Test ResolvePath
	resolved := resolver.ResolvePath("memory_context.md")
	if !strings.HasSuffix(resolved, "memory_context.md") {
		t.Errorf("expected resolved path to end with memory_context.md, got %q", resolved)
	}
}

// TestNewMemoryCommand_Flags tests the flags of the memory command
func TestNewMemoryCommand_Flags(t *testing.T) {
	registry := createTestRegistry()
	cmd := NewMemoryCommand(registry)

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

// TestMemoryCommand_RunE tests the RunE function of the memory command
func TestMemoryCommand_RunE(t *testing.T) {
	registry := createTestRegistry()
	cmd := NewMemoryCommand(registry)

	// The RunE function requires args, so we'll set them
	cmd.SetArgs([]string{"testdata/memory_context.md"})

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

	err = cmd.Args(cmd, []string{"testdata/memory_context.md", "extra"})
	if err == nil {
		t.Error("expected error for wrong number of args, got nil")
	}

	// Test Args validator with correct number of args
	err = cmd.Args(cmd, []string{"testdata/memory_context.md"})
	if err != nil {
		t.Errorf("expected no error for correct number of args, got %v", err)
	}
}
