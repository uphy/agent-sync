package template

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/agent-def/internal/agent"
)

func TestFileFunc_NotNil(t *testing.T) {
	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})
	registry.Register(&agent.Roo{})

	engine := &Engine{
		AgentRegistry: registry,
	}
	fn := engine.FileFunc()
	if fn == nil {
		t.Fatal("expected FileFunc to return a function, got nil")
	}
}

func TestIncludeFunc_NotNil(t *testing.T) {
	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})
	registry.Register(&agent.Roo{})

	engine := &Engine{
		AgentRegistry: registry,
	}
	fn := engine.IncludeFunc()
	if fn == nil {
		t.Fatal("expected IncludeFunc to return a function, got nil")
	}
}

func TestReferenceFunc_NotNil(t *testing.T) {
	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})
	registry.Register(&agent.Roo{})

	engine := &Engine{
		AgentRegistry: registry,
	}
	fn := engine.ReferenceFunc()
	if fn == nil {
		t.Fatal("expected ReferenceFunc to return a function, got nil")
	}
}

func TestMCPFunc_NotNil(t *testing.T) {
	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})
	registry.Register(&agent.Roo{})

	engine := &Engine{
		AgentRegistry: registry,
	}
	fn := engine.MCPFunc()
	if fn == nil {
		t.Fatal("expected MCPFunc to return a function, got nil")
	}
}

func TestFileFunc_ErrorOnInvalidUsage(t *testing.T) {
	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})
	registry.Register(&agent.Roo{})

	engine := &Engine{
		AgentRegistry: registry,
	}
	fn := engine.FileFunc().(func(string) (string, error))
	_, err := fn("path/to/file")
	if err == nil {
		t.Fatal("expected error on invalid FileFunc usage, got nil")
	}
}

func TestIncludeFunc_ErrorOnInvalidUsage(t *testing.T) {
	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})
	registry.Register(&agent.Roo{})

	// モックファイルリゾルバーを作成 - 存在しないファイルのテスト用
	mockResolver := &MockFileResolver{
		files: map[string]string{},
		base:  "/base",
	}

	engine := &Engine{
		AgentRegistry: registry,
		FileResolver:  mockResolver,
		BasePath:      "/base",
	}
	fn := engine.IncludeFunc().(func(string) (string, error))
	_, err := fn("path")
	if err == nil {
		t.Fatal("expected error on invalid IncludeFunc usage, got nil")
	}
}

func TestReferenceFunc_ErrorOnInvalidUsage(t *testing.T) {
	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})
	registry.Register(&agent.Roo{})

	// モックファイルリゾルバーを作成 - 存在しないファイルのテスト用
	mockResolver := &MockFileResolver{
		files: map[string]string{},
		base:  "/base",
	}

	engine := &Engine{
		AgentRegistry: registry,
		FileResolver:  mockResolver,
		BasePath:      "/base",
		References:    make(map[string]string),
	}
	fn := engine.ReferenceFunc().(func(string) (string, error))
	_, err := fn("path")
	if err == nil {
		t.Fatal("expected error on invalid ReferenceFunc usage, got nil")
	}
}

func TestMCPFunc_ErrorOnInvalidUsage(t *testing.T) {
	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})
	registry.Register(&agent.Roo{})

	engine := &Engine{
		AgentRegistry: registry,
	}
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
	registry.Register(&agent.Roo{})

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

	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})
	registry.Register(&agent.Roo{})

	engine := &Engine{
		FileResolver:  mockResolver,
		BasePath:      "/base",
		AgentRegistry: registry,
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

	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})
	registry.Register(&agent.Roo{})

	engine := &Engine{
		FileResolver:  mockResolver,
		References:    make(map[string]string),
		BasePath:      "/base",
		AgentRegistry: registry,
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

	testCases := []struct {
		agentType string
		agent     string
		command   string
		args      []string
		expected  string
	}{
		{
			agentType: "claude",
			agent:     "github",
			command:   "get-issue",
			args:      []string{},
			expected:  "MCP tool `github.get-issue`",
		},
		{
			agentType: "roo",
			agent:     "github",
			command:   "get-issue",
			args:      []string{},
			expected:  "MCP tool `github.get-issue`",
		},
		{
			agentType: "claude",
			agent:     "jira",
			command:   "get-ticket",
			args:      []string{"PROJ-123"},
			expected:  "MCP tool `jira.get-ticket(PROJ-123)`",
		},
		{
			agentType: "roo",
			agent:     "jira",
			command:   "get-ticket",
			args:      []string{"PROJ-123"},
			expected:  "MCP tool `jira.get-ticket(PROJ-123)`",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.agentType+"_"+tc.command+"_"+strings.Join(tc.args, "_"), func(t *testing.T) {
			engine := &Engine{
				AgentType:     tc.agentType,
				AgentRegistry: registry,
			}

			fn := engine.MCPFunc().(func(string, string, ...string) (string, error))
			result, err := fn(tc.agent, tc.command, tc.args...)

			if err != nil {
				t.Fatalf("expected no error for %s agent, got %v", tc.agentType, err)
			}

			if result != tc.expected {
				t.Errorf("FormatMCP for %s: expected %q, got %q",
					tc.agentType, tc.expected, result)
			}
		})
	}
}

func TestResolveTemplatePath(t *testing.T) {
	// Create a test engine
	basePath := "/base/path"
	engine := &Engine{
		BasePath: basePath,
	}

	// Test cases
	testCases := []struct {
		name            string
		path            string
		currentFilePath string
		expected        string
	}{
		{
			name:            "Windows-style absolute path",
			path:            "C:\\absolute\\path", // Only works on Windows but test is valid everywhere
			currentFilePath: "/current/file/path",
			expected:        "C:\\absolute\\path", // Should be preserved as-is
		},
		{
			name:            "Slash-prefixed path",
			path:            "/relative/to/base",
			currentFilePath: "/current/file/path",
			expected:        filepath.Join(basePath, "relative/to/base"),
		},
		{
			name:            "Dot-prefixed path (with current file)",
			path:            "./relative/to/current",
			currentFilePath: "/current/file/path",
			expected:        filepath.Join("/current/file", "relative/to/current"),
		},
		{
			name:            "Parent path (with current file)",
			path:            "../sibling/directory",
			currentFilePath: "/current/file/path",
			expected:        filepath.Join("/current", "sibling/directory"),
		},
		{
			name:            "Plain relative path (no ./ or ../)",
			path:            "relative/to/base",
			currentFilePath: "/current/file/path",
			expected:        filepath.Join(basePath, "relative/to/base"), // Should be relative to BasePath
		},
		{
			name:            "Plain relative path (without current file)",
			path:            "relative/to/base",
			currentFilePath: "",
			expected:        filepath.Join(basePath, "relative/to/base"),
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			engine.CurrentFilePath = tc.currentFilePath
			result := engine.resolveTemplatePath(tc.path)
			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}
		})
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
