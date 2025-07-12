package template

import (
	"testing"

	"path/filepath"
	"strings"

	"github.com/user/agent-def/internal/agent"
)

func TestExecute_ReturnsRawContent(t *testing.T) {
	engine := &Engine{}
	input := "plain text content"
	output, err := engine.Execute(input, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output != input {
		t.Errorf("expected output %q, got %q", input, output)
	}
}

func TestExecute_WithTemplateVariable(t *testing.T) {
	engine := &Engine{}
	templateContent := "hello {{.Name}}"
	data := map[string]interface{}{"Name": "World"}
	output, err := engine.Execute(templateContent, data)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output != "hello World" {
		t.Errorf("expected output %q, got %q", "hello World", output)
	}
}

// MockFileResolver implements FileResolver for testing
type MockFileResolver struct {
	files map[string]string
	base  string
}

func NewMockFileResolver(base string) *MockFileResolver {
	return &MockFileResolver{
		files: make(map[string]string),
		base:  base,
	}
}

func (r *MockFileResolver) AddFile(path, content string) {
	r.files[path] = content
}

func (r *MockFileResolver) Read(path string) ([]byte, error) {
	if content, ok := r.files[path]; ok {
		return []byte(content), nil
	}
	return nil, &ErrFileNotFound{Path: path}
}

func (r *MockFileResolver) Exists(path string) bool {
	_, ok := r.files[path]
	return ok
}

func (r *MockFileResolver) ResolvePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(r.base, path)
}

// ErrFileNotFound implements error for file not found
type ErrFileNotFound struct {
	Path string
}

func (e *ErrFileNotFound) Error() string {
	return "file not found: " + e.Path
}

func TestIncludeFunc(t *testing.T) {
	// Setup
	mockResolver := NewMockFileResolver("/base")
	mockResolver.AddFile("/base/testdata/content.md", "# Included Content\n\nThis is content from an included file.")

	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})

	engine := &Engine{
		FileResolver:  mockResolver,
		References:    make(map[string]string),
		AgentType:     "claude",
		BasePath:      "/base",
		AgentRegistry: registry,
	}

	// Test the include function
	fn := engine.IncludeFunc().(func(string) (string, error))
	result, err := fn("testdata/content.md")

	// Verify
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := "# Included Content\n\nThis is content from an included file."
	if result != expected {
		t.Errorf("expected result %q, got %q", expected, result)
	}
}

func TestReferenceFunc(t *testing.T) {
	// Setup
	mockResolver := NewMockFileResolver("/base")
	mockResolver.AddFile("/base/testdata/content.md", "# Referenced Content\n\nThis is content from a referenced file.")

	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})

	engine := &Engine{
		FileResolver:  mockResolver,
		References:    make(map[string]string),
		AgentType:     "claude",
		BasePath:      "/base",
		AgentRegistry: registry,
	}

	// Test the reference function
	fn := engine.ReferenceFunc().(func(string) (string, error))
	result, err := fn("testdata/content.md")

	// Verify
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check if the reference is stored
	if content, ok := engine.References["testdata/content.md"]; !ok {
		t.Error("reference was not stored in the engine")
	} else if content != "# Referenced Content\n\nThis is content from a referenced file." {
		t.Errorf("stored reference content doesn't match expected")
	}

	// Check the reference marker
	expected := "[参考: testdata/content.md]"
	if result != expected {
		t.Errorf("expected result %q, got %q", expected, result)
	}
}

func TestFileFunc(t *testing.T) {
	// Setup
	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})

	engine := &Engine{
		AgentType:     "claude",
		AgentRegistry: registry,
	}

	// Test the file function
	fn := engine.FileFunc().(func(string) (string, error))
	result, err := fn("path/to/file.md")

	// Verify
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Test that the result contains the path (exact format depends on agent implementation)
	if !strings.Contains(result, "file.md") {
		t.Errorf("expected result to contain file path, got %q", result)
	}
}

func TestMCPFunc(t *testing.T) {
	// Setup
	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})

	engine := &Engine{
		AgentType:     "claude",
		AgentRegistry: registry,
	}

	// Test the MCP function
	fn := engine.MCPFunc().(func(string, string, ...string) (string, error))
	result, err := fn("test-agent", "test-command", "arg1", "arg2")

	// Verify
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Test that the result contains all parts (exact format depends on agent implementation)
	if !strings.Contains(result, "test-agent") ||
		!strings.Contains(result, "test-command") ||
		!strings.Contains(result, "arg1") ||
		!strings.Contains(result, "arg2") {
		t.Errorf("expected result to contain all MCP parts, got %q", result)
	}
}

func TestRecursiveTemplateProcessing(t *testing.T) {
	// Setup
	mockResolver := NewMockFileResolver("/base")
	mockResolver.AddFile("/base/testdata/content.md", "# Included Content\n\nThis is content from an included file.")
	mockResolver.AddFile("/base/testdata/include_template.md", "# Template with Include\n\nThis template includes another file:\n\n{{include \"testdata/content.md\"}}\n\nEnd of template.")
	mockResolver.AddFile("/base/testdata/recursive_template.md", "# Recursive Template\n\nThis template includes another template that also has an include:\n\n{{include \"testdata/include_template.md\"}}\n\nEnd of recursive template.")

	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})

	engine := &Engine{
		FileResolver:  mockResolver,
		References:    make(map[string]string),
		AgentType:     "claude",
		BasePath:      "/base",
		AgentRegistry: registry,
	}

	// Test the recursive include
	output, err := engine.Execute("{{include \"testdata/recursive_template.md\"}}", nil)

	// Verify
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check that all content is included
	if !strings.Contains(output, "Recursive Template") ||
		!strings.Contains(output, "Template with Include") ||
		!strings.Contains(output, "Included Content") ||
		!strings.Contains(output, "End of recursive template") {
		t.Errorf("expected output to contain all nested content, got %q", output)
	}
}

func TestReferenceCollection(t *testing.T) {
	// Setup
	mockResolver := NewMockFileResolver("/base")
	mockResolver.AddFile("/base/testdata/content1.md", "# Referenced Content 1\n\nThis is content from first referenced file.")
	mockResolver.AddFile("/base/testdata/content2.md", "# Referenced Content 2\n\nThis is content from second referenced file.")

	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})

	engine := &Engine{
		FileResolver:  mockResolver,
		References:    make(map[string]string),
		AgentType:     "claude",
		BasePath:      "/base",
		AgentRegistry: registry,
	}

	// Process a template with multiple references
	output, err := engine.Execute("# Template with References\n\nFirst: {{reference \"testdata/content1.md\"}}\n\nSecond: {{reference \"testdata/content2.md\"}}", nil)

	// Verify
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check that the references section is appended
	if !strings.Contains(output, "## References") {
		t.Error("expected output to contain References section")
	}

	// Check that both references are included
	if !strings.Contains(output, "### testdata/content1.md") ||
		!strings.Contains(output, "### testdata/content2.md") {
		t.Error("expected output to contain both references")
	}

	if !strings.Contains(output, "Referenced Content 1") ||
		!strings.Contains(output, "Referenced Content 2") {
		t.Error("expected output to contain content from both references")
	}
}

func TestExecute_ErrorHandling(t *testing.T) {
	// Setup
	mockResolver := NewMockFileResolver("/base")
	// Don't add any files, to test error handling

	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})

	engine := &Engine{
		FileResolver:  mockResolver,
		References:    make(map[string]string),
		AgentType:     "claude",
		BasePath:      "/base",
		AgentRegistry: registry,
	}

	// Test include with non-existent file
	_, err := engine.Execute("{{include \"non-existent-file.md\"}}", nil)

	// Verify
	if err == nil {
		t.Error("expected error for non-existent include file, got nil")
	}
}
