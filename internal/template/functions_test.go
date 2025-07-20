package template

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/uphy/agent-def/internal/agent"
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
	fn := engine.IncludeFunc(true)
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
	fn := engine.ReferenceFunc(true)
	if fn == nil {
		t.Fatal("expected ReferenceFunc to return a function, got nil")
	}
}

func TestIncludeRawFunc_NotNil(t *testing.T) {
	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})
	registry.Register(&agent.Roo{})

	engine := &Engine{
		AgentRegistry: registry,
	}
	fn := engine.IncludeFunc(false)
	if fn == nil {
		t.Fatal("expected IncludeRawFunc to return a function, got nil")
	}
}

func TestReferenceRawFunc_NotNil(t *testing.T) {
	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})
	registry.Register(&agent.Roo{})

	engine := &Engine{
		AgentRegistry: registry,
	}
	fn := engine.ReferenceFunc(false)
	if fn == nil {
		t.Fatal("expected ReferenceRawFunc to return a function, got nil")
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
	mockResolver := NewMockFileResolver("/base")

	engine := &Engine{
		AgentRegistry:      registry,
		FileResolver:       mockResolver,
		absTemplateBaseDir: "/base",
		absCurrentFilePath: "/base/current.md",
	}

	// Setup expected glob results - to simulate resolving the path
	mockResolver.ExpectGlob("/base/path", []string{"/base/path"})
	fn := engine.IncludeFunc(true).(func(...string) (string, error))
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
	mockResolver := NewMockFileResolver("/base")

	engine := &Engine{
		AgentRegistry:      registry,
		FileResolver:       mockResolver,
		absTemplateBaseDir: "/base",
		References:         make(map[string]string),
		absCurrentFilePath: "/base/current.md",
	}

	// Setup expected glob results - to simulate resolving the path
	mockResolver.ExpectGlob("/base/path", []string{"/base/path"})
	fn := engine.ReferenceFunc(true).(func(...string) (string, error))
	_, err := fn("path")
	if err == nil {
		t.Fatal("expected error on invalid ReferenceFunc usage, got nil")
	}
}

func TestIncludeRawFunc_ErrorOnInvalidUsage(t *testing.T) {
	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})
	registry.Register(&agent.Roo{})

	// Create mock file resolver for non-existent file testing
	mockResolver := NewMockFileResolver("/base")

	engine := &Engine{
		AgentRegistry:      registry,
		FileResolver:       mockResolver,
		absTemplateBaseDir: "/base",
		absCurrentFilePath: "/base/current.md",
	}

	// Setup expected glob results - to simulate resolving the path
	mockResolver.ExpectGlob("/base/path", []string{"/base/path"})
	fn := engine.IncludeFunc(false).(func(...string) (string, error))
	_, err := fn("path")
	if err == nil {
		t.Fatal("expected error on invalid IncludeFunc(false) usage, got nil")
	}
}

func TestReferenceRawFunc_ErrorOnInvalidUsage(t *testing.T) {
	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})
	registry.Register(&agent.Roo{})

	// Create mock file resolver for non-existent file testing
	mockResolver := NewMockFileResolver("/base")

	engine := &Engine{
		AgentRegistry:      registry,
		FileResolver:       mockResolver,
		absTemplateBaseDir: "/base",
		References:         make(map[string]string),
		absCurrentFilePath: "/base/current.md",
	}

	// Setup expected glob results - to simulate resolving the path
	mockResolver.ExpectGlob("/base/path", []string{"/base/path"})
	fn := engine.ReferenceFunc(false).(func(...string) (string, error))
	_, err := fn("path")
	if err == nil {
		t.Fatal("expected error on invalid ReferenceRawFunc usage, got nil")
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
	// Setup specialized error mock file resolver
	mockResolver := NewErrorFileResolver("/base")
	mockResolver.AddFile("/base/test.md", "Test content")
	mockResolver.AddFile("/base/test2.md", "Test content2")

	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})
	registry.Register(&agent.Roo{})

	engine := &Engine{
		FileResolver:       mockResolver,
		absTemplateBaseDir: "/base",
		AgentRegistry:      registry,
		absCurrentFilePath: "/base/current.md",
	}

	// Setup expected glob results
	mockResolver.ExpectGlob("/base/test.md", []string{"/base/test.md"})
	mockResolver.ExpectGlob("/base/test2.md", []string{"/base/test2.md"})

	// Test with existing file
	fn := engine.IncludeFunc(true).(func(...string) (string, error))
	result, err := fn("@/test.md", "@/test2.md")

	// Verify
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != "Test content\n\nTest content2" {
		t.Errorf("expected result %q, got %q", "Test content\n\nTest content2", result)
	}

	// Test with non-existent file
	// Add the path to the glob results but mark it as a non-existent file
	nonExistentPath := "/base/current.md/missing.md"
	mockResolver.ExpectGlob(nonExistentPath, []string{nonExistentPath})
	mockResolver.AddErrorPath(nonExistentPath)

	result, err = fn("missing.md")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result for non-existent file, got %q", result)
	}
}

func TestReferenceFunc_WithMockFileResolver(t *testing.T) {
	// Setup mock file resolver
	mockResolver := NewMockFileResolver("/base")
	mockResolver.AddFile("/base/ref.md", "Reference content")
	mockResolver.AddFile("/base/ref2.md", "Reference content2")

	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})
	registry.Register(&agent.Roo{})

	engine := &Engine{
		FileResolver:       mockResolver,
		References:         make(map[string]string),
		absTemplateBaseDir: "/base",
		AgentRegistry:      registry,
		absCurrentFilePath: "/base/current.md",
	}

	// Setup expected glob results
	mockResolver.ExpectGlob("/base/ref.md", []string{"/base/ref.md"})
	mockResolver.ExpectGlob("/base/ref2.md", []string{"/base/ref2.md"})

	// Test with existing file
	fn := engine.ReferenceFunc(true).(func(...string) (string, error))
	result, err := fn("@/ref.md", "@/ref2.md")

	// Verify
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check reference marker
	if result != "[参考: ref.md], [参考: ref2.md]" {
		t.Errorf("expected result %q, got %q", "[参考: ref.md], [参考: ref2.md]", result)
	}

	// Check stored reference
	if content, ok := engine.References["ref.md"]; !ok {
		t.Error("reference was not stored in engine")
	} else if content != "Reference content" {
		t.Errorf("expected stored content %q, got %q", "Reference content", content)
	}
	if content, ok := engine.References["ref2.md"]; !ok {
		t.Error("reference was not stored in engine")
	} else if content != "Reference content2" {
		fmt.Println(content)
		t.Errorf("expected stored content %q, got %q", "Reference content2", content)
	}
}

func TestIncludeRawFunc_WithMockFileResolver(t *testing.T) {
	// Setup specialized error mock file resolver
	mockResolver := NewErrorFileResolver("/base")
	mockResolver.AddFile("/base/template_file.md", "Content with {{template}} syntax that should remain as-is")

	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})
	registry.Register(&agent.Roo{})

	engine := &Engine{
		FileResolver:       mockResolver,
		absTemplateBaseDir: "/base",
		AgentRegistry:      registry,
		absCurrentFilePath: "/base/current.md",
	}

	// Setup expected glob results
	mockResolver.ExpectGlob("/base/template_file.md", []string{"/base/template_file.md"})

	// Test with existing file
	fn := engine.IncludeFunc(false).(func(...string) (string, error))
	result, err := fn("@/template_file.md")

	// Verify
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check that template syntax is preserved as-is (not processed)
	expected := "Content with {{template}} syntax that should remain as-is"
	if result != expected {
		t.Errorf("expected raw content %q, got %q", expected, result)
	}

	// Test with non-existent file
	// Add the path to the glob results but mark it as a non-existent file
	nonExistentPath := "/base/current.md/missing.md"
	mockResolver.ExpectGlob(nonExistentPath, []string{nonExistentPath})
	mockResolver.AddErrorPath(nonExistentPath)

	result, err = fn("missing.md")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result for non-existent file, got %q", result)
	}
}

func TestReferenceRawFunc_WithMockFileResolver(t *testing.T) {
	// Setup mock file resolver with a file containing template syntax
	mockResolver := NewMockFileResolver("/base")
	mockResolver.AddFile("/base/template_ref.md", "Reference with {{template}} syntax that should remain as-is")

	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})
	registry.Register(&agent.Roo{})

	engine := &Engine{
		FileResolver:       mockResolver,
		References:         make(map[string]string),
		absTemplateBaseDir: "/base",
		AgentRegistry:      registry,
		absCurrentFilePath: "/base/current.md",
	}

	// Setup expected glob results
	mockResolver.ExpectGlob("/base/template_ref.md", []string{"/base/template_ref.md"})

	// Test with existing file
	fn := engine.ReferenceFunc(false).(func(...string) (string, error))
	result, err := fn("@/template_ref.md")

	// Verify
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check reference marker
	if result != "[参考: template_ref.md]" {
		t.Errorf("expected result %q, got %q", "[参考: template_ref.md]", result)
	}

	// Check stored raw reference with template syntax preserved
	expected := "Reference with {{template}} syntax that should remain as-is"
	if content, ok := engine.References["template_ref.md"]; !ok {
		t.Error("reference was not stored in engine")
	} else if content != expected {
		t.Errorf("expected raw content %q, got %q", expected, content)
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
	// Create a test engine with mock FileResolver
	basePath := "/base/path"
	mockResolver := NewMockFileResolver(basePath)
	engine := &Engine{
		absTemplateBaseDir: basePath,
		FileResolver:       mockResolver,
	}

	// Test cases
	testCases := []struct {
		name            string
		paths           []string
		currentFilePath string
		expected        []string
		expectError     bool
	}{
		{
			name:            "@/-prefixed path",
			paths:           []string{"@/relative/to/base"},
			currentFilePath: "/current/file/path",
			expected:        []string{filepath.Join(basePath, "relative/to/base")},
			expectError:     false,
		},
		{
			name:            "Dot-prefixed path (with current file)",
			paths:           []string{"./relative/to/current"},
			currentFilePath: "/current/file/path",
			expected:        []string{filepath.Join("/current/file", "relative/to/current")},
			expectError:     false,
		},
		{
			name:            "Parent path (with current file)",
			paths:           []string{"../sibling/directory"},
			currentFilePath: "/current/file/path",
			expected:        []string{filepath.Join("/current", "sibling/directory")},
			expectError:     false,
		},
		{
			name:            "Parent-parent path (with current file)",
			paths:           []string{"../../grandparent/directory"},
			currentFilePath: "/current/file/path",
			expected:        []string{filepath.Join("/", "grandparent/directory")},
			expectError:     false,
		},
		{
			name:            "Parent-parent with additional directories",
			paths:           []string{"../../x/y/z"},
			currentFilePath: "/a/b/c/d/file.md",
			expected:        []string{filepath.Join("/a/b", "x/y/z")},
			expectError:     false,
		},
		{
			name:            "Plain relative path (no prefix)",
			paths:           []string{"relative/to/current"},
			currentFilePath: "/current/file/path",
			expected:        []string{filepath.Join("/current/file", "relative/to/current")},
			expectError:     false,
		},
		{
			name:            "Plain relative path (without current file)",
			paths:           []string{"relative/to/current"},
			currentFilePath: "",
			expected:        []string{filepath.Join("", "relative/to/current")},
			expectError:     false,
		},
		{
			name:            "Multiple paths mixed",
			paths:           []string{"@/file1.md", "./file2.md", "../file3.md"},
			currentFilePath: "/current/file/path",
			expected: []string{
				filepath.Join(basePath, "file1.md"),
				filepath.Join("/current/file", "file2.md"),
				filepath.Join("/current", "file3.md"),
			},
			expectError: false,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			engine.absCurrentFilePath = tc.currentFilePath

			// Setup mock resolver to return the expected resolved paths for this test case
			if !tc.expectError {
				for i, path := range tc.paths {
					var resolvedPath string
					if strings.HasPrefix(path, "@/") {
						resolvedPath = filepath.Join(basePath, strings.TrimPrefix(path, "@/"))
					} else {
						resolvedPath = filepath.Join(filepath.Dir(tc.currentFilePath), path)
					}
					mockResolver.ExpectGlob(resolvedPath, []string{tc.expected[i]})
				}
			}

			result, err := engine.resolveTemplatePath(tc.paths)

			// Check error expectation
			if tc.expectError && err == nil {
				t.Errorf("Expected error but got none")
			} else if !tc.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			// Check result only if no error is expected
			if !tc.expectError {
				if len(result) != len(tc.expected) {
					t.Errorf("Expected %d results, got %d", len(tc.expected), len(result))
				} else {
					for i, exp := range tc.expected {
						if result[i] != exp {
							t.Errorf("Expected result[%d] to be %q, got %q", i, exp, result[i])
						}
					}
				}
			}
		})
	}
}

func TestResolveTemplatePathSingle(t *testing.T) {
	// Create a test engine with mock FileResolver
	basePath := "/base/path"
	mockResolver := NewMockFileResolver(basePath)
	engine := &Engine{
		absTemplateBaseDir: basePath,
		FileResolver:       mockResolver,
	}

	// Test cases
	testCases := []struct {
		name            string
		path            string
		currentFilePath string
		expected        string
		expectError     bool
	}{
		{
			name:            "@/-prefixed path",
			path:            "@/relative/to/base",
			currentFilePath: "/current/file/path",
			expected:        filepath.Join(basePath, "relative/to/base"),
			expectError:     false,
		},
		{
			name:            "Dot-prefixed path (with current file)",
			path:            "./relative/to/current",
			currentFilePath: "/current/file/path",
			expected:        filepath.Join("/current/file", "relative/to/current"),
			expectError:     false,
		},
		{
			name:            "Parent path (with current file)",
			path:            "../sibling/directory",
			currentFilePath: "/current/file/path",
			expected:        filepath.Join("/current", "sibling/directory"),
			expectError:     false,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			engine.absCurrentFilePath = tc.currentFilePath

			// Setup mock resolver for resolveTemplatePathSingle
			if !tc.expectError {
				var resolvedPath string
				if strings.HasPrefix(tc.path, "@/") {
					resolvedPath = filepath.Join(basePath, strings.TrimPrefix(tc.path, "@/"))
				} else {
					resolvedPath = filepath.Join(filepath.Dir(tc.currentFilePath), tc.path)
				}
				mockResolver.ExpectGlob(resolvedPath, []string{tc.expected})
			}

			result, err := engine.resolveTemplatePath([]string{tc.path})

			// Check error expectation
			if tc.expectError && err == nil {
				t.Errorf("Expected error but got none")
			} else if !tc.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			// Check result only if no error is expected
			if !tc.expectError && result[0] != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result[0])
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
