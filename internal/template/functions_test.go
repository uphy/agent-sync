package template

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/agent-def/internal/agent"
)

func TestFileFunc_NotNil(t *testing.T) {
	engine := &Engine{}
	fn := engine.FileFunc()
	if fn == nil {
		t.Fatal("expected FileFunc to return a function, got nil")
	}
}

func TestIncludeFunc_NotNil(t *testing.T) {
	engine := &Engine{}
	fn := engine.IncludeFunc()
	if fn == nil {
		t.Fatal("expected IncludeFunc to return a function, got nil")
	}
}

func TestReferenceFunc_NotNil(t *testing.T) {
	engine := &Engine{}
	fn := engine.ReferenceFunc()
	if fn == nil {
		t.Fatal("expected ReferenceFunc to return a function, got nil")
	}
}

func TestMCPFunc_NotNil(t *testing.T) {
	engine := &Engine{}
	fn := engine.MCPFunc()
	if fn == nil {
		t.Fatal("expected MCPFunc to return a function, got nil")
	}
}

func TestFileFunc_ErrorOnInvalidUsage(t *testing.T) {
	engine := &Engine{}
	fn := engine.FileFunc().(func(string) (string, error))
	_, err := fn("path/to/file")
	if err == nil {
		t.Fatal("expected error on invalid FileFunc usage, got nil")
	}
}

func TestIncludeFunc_ErrorOnInvalidUsage(t *testing.T) {
	engine := &Engine{}
	fn := engine.IncludeFunc().(func(string) (string, error))
	_, err := fn("path")
	if err == nil {
		t.Fatal("expected error on invalid IncludeFunc usage, got nil")
	}
}

func TestReferenceFunc_ErrorOnInvalidUsage(t *testing.T) {
	engine := &Engine{}
	fn := engine.ReferenceFunc().(func(string) (string, error))
	_, err := fn("path")
	if err == nil {
		t.Fatal("expected error on invalid ReferenceFunc usage, got nil")
	}
}

func TestMCPFunc_ErrorOnInvalidUsage(t *testing.T) {
	engine := &Engine{}
	fn := engine.MCPFunc().(func(string, string, ...string) (string, error))
	_, err := fn("agent", "command")
	if err == nil {
		t.Fatal("expected error on invalid MCPFunc usage, got nil")
	}
}

// Additional more comprehensive tests for helper functions

func TestFileFunc_ValidUsage(t *testing.T) {
	// Setup
	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})

	engine := &Engine{
		AgentType:     "claude",
		AgentRegistry: registry,
	}

	// Test with valid agent
	fn := engine.FileFunc().(func(string) (string, error))
	result, err := fn("example.md")

	// Verify
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// The actual format depends on Claude agent implementation
	if !strings.Contains(result, "example.md") {
		t.Errorf("expected result to contain file path, got %q", result)
	}
}

func TestIncludeFunc_WithMockFileResolver(t *testing.T) {
	// Setup mock file resolver
	mockResolver := &MockFileResolver{
		files: map[string]string{
			"/base/test.md": "Test content",
		},
		base: "/base",
	}

	engine := &Engine{
		FileResolver: mockResolver,
		BasePath:     "/base",
	}

	// Test with existing file
	fn := engine.IncludeFunc().(func(string) (string, error))
	result, err := fn("test.md")

	// Verify
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != "Test content" {
		t.Errorf("expected result %q, got %q", "Test content", result)
	}

	// Test with non-existent file
	_, err = fn("missing.md")
	if err == nil {
		t.Error("expected error for non-existent file, got nil")
	}
}

func TestReferenceFunc_WithMockFileResolver(t *testing.T) {
	// Setup mock file resolver
	mockResolver := &MockFileResolver{
		files: map[string]string{
			"/base/ref.md": "Reference content",
		},
		base: "/base",
	}

	engine := &Engine{
		FileResolver: mockResolver,
		References:   make(map[string]string),
		BasePath:     "/base",
	}

	// Test with existing file
	fn := engine.ReferenceFunc().(func(string) (string, error))
	result, err := fn("ref.md")

	// Verify
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check reference marker
	if result != "[参考: ref.md]" {
		t.Errorf("expected result %q, got %q", "[参考: ref.md]", result)
	}

	// Check stored reference
	if content, ok := engine.References["ref.md"]; !ok {
		t.Error("reference was not stored in engine")
	} else if content != "Reference content" {
		t.Errorf("expected stored content %q, got %q", "Reference content", content)
	}
}

func TestMCPFunc_WithValidAgent(t *testing.T) {
	// Setup registry with agents
	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})
	registry.Register(&agent.Roo{})

	// Test with Claude agent
	claudeEngine := &Engine{
		AgentType:     "claude",
		AgentRegistry: registry,
	}

	claudeFn := claudeEngine.MCPFunc().(func(string, string, ...string) (string, error))
	claudeResult, err := claudeFn("test-agent", "test-command", "arg1", "arg2")

	if err != nil {
		t.Fatalf("expected no error for Claude agent, got %v", err)
	}

	// Test with Roo agent
	rooEngine := &Engine{
		AgentType:     "roo",
		AgentRegistry: registry,
	}

	rooFn := rooEngine.MCPFunc().(func(string, string, ...string) (string, error))
	rooResult, err := rooFn("test-agent", "test-command", "arg1", "arg2")

	if err != nil {
		t.Fatalf("expected no error for Roo agent, got %v", err)
	}

	// Verify different implementations produce different results
	if claudeResult == rooResult {
		t.Errorf("expected different results for different agents, got same: %q", claudeResult)
	}
}

func TestResolveRelativePath(t *testing.T) {
	engine := &Engine{
		BasePath: "/base/path",
	}

	// Test with relative path
	result := engine.resolveRelativePath("relative/file.md")
	expected := filepath.Join("/base/path", "relative/file.md")
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}

	// Test with absolute path
	absPath := "/absolute/file.md"
	result = engine.resolveRelativePath(absPath)
	if result != absPath {
		t.Errorf("expected %q, got %q", absPath, result)
	}
}

func TestNormalizePath(t *testing.T) {
	engine := &Engine{}

	path := "path/to/file.md"
	normalized := engine.normalizePath(path)

	// On all systems, the result should use forward slashes
	if strings.Contains(normalized, "\\") {
		t.Errorf("expected normalized path without backslashes, got %q", normalized)
	}
}
