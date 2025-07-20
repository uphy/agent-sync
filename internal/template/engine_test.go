package template

import (
	"sort"
	"testing"

	"path/filepath"
	"strings"

	"github.com/uphy/agent-sync/internal/agent"
)

func TestExecute_ReturnsRawContent(t *testing.T) {
	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})
	registry.Register(&agent.Roo{})

	engine := &Engine{
		AgentRegistry: registry,
	}
	input := "plain text content"
	output, err := engine.Execute("/test/path", input, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output != input {
		t.Errorf("expected output %q, got %q", input, output)
	}
}

func TestExecute_WithTemplateVariable(t *testing.T) {
	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})
	registry.Register(&agent.Roo{})

	engine := &Engine{
		AgentRegistry: registry,
	}
	templateContent := "hello {{.Name}}"
	data := map[string]interface{}{"Name": "World"}
	output, err := engine.Execute("/test/path", templateContent, data)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output != "hello World" {
		t.Errorf("expected output %q, got %q", "hello World", output)
	}
}

// MockFileResolver implements FileResolver for testing
type MockFileResolver struct {
	files      map[string]string
	base       string
	globExpect map[string][]string
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

func (r *MockFileResolver) ExpectGlob(pattern string, expected []string) {
	if r.globExpect == nil {
		r.globExpect = make(map[string][]string)
	}
	r.globExpect[pattern] = expected
}

func (r *MockFileResolver) Glob(patterns []string) ([]string, error) {
	if r.globExpect == nil {
		return nil, nil
	}
	results := make(map[string]struct{})
	for _, pattern := range patterns {
		if expected, ok := r.globExpect[pattern]; ok {
			for _, e := range expected {
				results[e] = struct{}{}
			}
		}
	}
	resultsSlice := make([]string, 0, len(results))
	for path := range results {
		resultsSlice = append(resultsSlice, path)
	}
	sort.Strings(resultsSlice)
	return resultsSlice, nil
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
	registry.Register(&agent.Roo{})

	engine := &Engine{
		FileResolver:       mockResolver,
		References:         make(map[string]string),
		AgentType:          "claude",
		absTemplateBaseDir: "/base",
		AgentRegistry:      registry,
		absCurrentFilePath: "/base/test.md",
	}

	// Setup expected glob results
	mockResolver.ExpectGlob("/base/testdata/content.md", []string{"/base/testdata/content.md"})

	// Test the include function
	fn := engine.IncludeFunc(true).(func(...string) (string, error))
	result, err := fn("@/testdata/content.md")

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
	registry.Register(&agent.Roo{})

	engine := &Engine{
		FileResolver:       mockResolver,
		References:         make(map[string]string),
		AgentType:          "claude",
		absTemplateBaseDir: "/base",
		AgentRegistry:      registry,
		absCurrentFilePath: "/base/test.md",
	}

	// Setup expected glob results
	mockResolver.ExpectGlob("/base/testdata/content.md", []string{"/base/testdata/content.md"})

	// Test the reference function
	fn := engine.ReferenceFunc(true).(func(...string) (string, error))
	result, err := fn("@/testdata/content.md")

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
	registry.Register(&agent.Roo{})
	registry.Register(&agent.Roo{})

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
	mockResolver.AddFile("/base/testdata/include_template.md", "# Template with Include\n\nThis template includes another file:\n\n{{include \"@/testdata/content.md\"}}\n\nEnd of template.")
	mockResolver.AddFile("/base/testdata/recursive_template.md", "# Recursive Template\n\nThis template includes another template that also has an include:\n\n{{include \"@/testdata/include_template.md\"}}\n\nEnd of recursive template.")

	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})
	registry.Register(&agent.Roo{})

	engine := &Engine{
		FileResolver:       mockResolver,
		References:         make(map[string]string),
		AgentType:          "claude",
		absTemplateBaseDir: "/base",
		AgentRegistry:      registry,
		absCurrentFilePath: "/base/test.md",
	}

	// Setup expected glob results
	mockResolver.ExpectGlob("/base/testdata/recursive_template.md", []string{"/base/testdata/recursive_template.md"})
	mockResolver.ExpectGlob("/base/testdata/include_template.md", []string{"/base/testdata/include_template.md"})
	mockResolver.ExpectGlob("/base/testdata/content.md", []string{"/base/testdata/content.md"})

	// Test the recursive include
	output, err := engine.Execute("/base/test.md", "{{include \"@/testdata/recursive_template.md\"}}", nil)

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
	registry.Register(&agent.Roo{})

	engine := &Engine{
		FileResolver:       mockResolver,
		References:         make(map[string]string),
		AgentType:          "claude",
		absTemplateBaseDir: "/base",
		AgentRegistry:      registry,
		absCurrentFilePath: "/base/test.md",
	}

	// Setup expected glob results
	mockResolver.ExpectGlob("/base/testdata/content1.md", []string{"/base/testdata/content1.md"})
	mockResolver.ExpectGlob("/base/testdata/content2.md", []string{"/base/testdata/content2.md"})

	// Process a template with multiple references
	output, err := engine.Execute("/base/test.md", "# Template with References\n\nFirst: {{reference \"@/testdata/content1.md\"}}\n\nSecond: {{reference \"@/testdata/content2.md\"}}", nil)

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

// ErrorFileResolver is a specific mock for error testing
type ErrorFileResolver struct {
	*MockFileResolver
	errorPaths map[string]bool
}

func NewErrorFileResolver(base string) *ErrorFileResolver {
	return &ErrorFileResolver{
		MockFileResolver: NewMockFileResolver(base),
		errorPaths:       make(map[string]bool),
	}
}

func (r *ErrorFileResolver) AddErrorPath(path string) {
	r.errorPaths[path] = true
}

func (r *ErrorFileResolver) Exists(path string) bool {
	// Return false for explicitly marked error paths
	if r.errorPaths[path] {
		return false
	}
	return r.MockFileResolver.Exists(path)
}

func (r *ErrorFileResolver) Read(path string) ([]byte, error) {
	// Return error for explicitly marked error paths
	if r.errorPaths[path] {
		return nil, &ErrFileNotFound{Path: path}
	}
	return r.MockFileResolver.Read(path)
}

// Tests for the new includeRaw and referenceRaw functions

func TestIncludeRawFunc(t *testing.T) {
	// Setup
	mockResolver := NewMockFileResolver("/base")
	mockResolver.AddFile("/base/testdata/template_with_syntax.md",
		"# Template with Syntax\n\n"+
			"This file contains template syntax that should not be processed by raw functions:\n\n"+
			"{{.Value}} - This is a variable\n"+
			"{{include \"@/testdata/content.md\"}} - This is an include directive\n"+
			"{{agent}} - This is a function call\n\n"+
			"The raw functions should preserve these as-is.")

	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})
	registry.Register(&agent.Roo{})

	engine := &Engine{
		FileResolver:       mockResolver,
		References:         make(map[string]string),
		AgentType:          "claude",
		absTemplateBaseDir: "/base",
		AgentRegistry:      registry,
		absCurrentFilePath: "/base/test.md",
	}

	// Setup expected glob results
	mockResolver.ExpectGlob("/base/testdata/template_with_syntax.md", []string{"/base/testdata/template_with_syntax.md"})

	// Test the includeRaw function
	fn := engine.IncludeFunc(false).(func(...string) (string, error))
	result, err := fn("@/testdata/template_with_syntax.md")

	// Verify
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check that template syntax is preserved as-is (not processed)
	if !strings.Contains(result, "{{.Value}}") {
		t.Errorf("includeRaw did not preserve template variable syntax")
	}

	if !strings.Contains(result, "{{include \"@/testdata/content.md\"}}") {
		t.Errorf("includeRaw did not preserve include directive syntax")
	}

	if !strings.Contains(result, "{{agent}}") {
		t.Errorf("includeRaw did not preserve function call syntax")
	}
}

func TestReferenceRawFunc(t *testing.T) {
	// Setup
	mockResolver := NewMockFileResolver("/base")
	mockResolver.AddFile("/base/testdata/template_with_syntax.md",
		"# Template with Syntax\n\n"+
			"This file contains template syntax that should not be processed by raw functions:\n\n"+
			"{{.Value}} - This is a variable\n"+
			"{{include \"testdata/content.md\"}} - This is an include directive\n"+
			"{{agent}} - This is a function call\n\n"+
			"The raw functions should preserve these as-is.")

	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})
	registry.Register(&agent.Roo{})

	engine := &Engine{
		FileResolver:       mockResolver,
		References:         make(map[string]string),
		AgentType:          "claude",
		absTemplateBaseDir: "/base",
		AgentRegistry:      registry,
		absCurrentFilePath: "/base/test.md",
	}

	// Setup expected glob results
	mockResolver.ExpectGlob("/base/testdata/template_with_syntax.md", []string{"/base/testdata/template_with_syntax.md"})

	// Test the referenceRaw function
	fn := engine.ReferenceFunc(false).(func(...string) (string, error))
	result, err := fn("@/testdata/template_with_syntax.md")

	// Verify
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check reference marker
	expected := "[参考: testdata/template_with_syntax.md]"
	if result != expected {
		t.Errorf("expected result %q, got %q", expected, result)
	}

	// Check that the stored reference contains the template syntax as-is
	storedContent, exists := engine.References["testdata/template_with_syntax.md"]
	if !exists {
		t.Error("reference was not stored in engine")
	}

	// Check that template syntax is preserved as-is in the stored reference
	if !strings.Contains(storedContent, "{{.Value}}") {
		t.Errorf("referenceRaw did not preserve template variable syntax")
	}

	if !strings.Contains(storedContent, "{{include \"testdata/content.md\"}}") {
		t.Errorf("referenceRaw did not preserve include directive syntax")
	}

	if !strings.Contains(storedContent, "{{agent}}") {
		t.Errorf("referenceRaw did not preserve function call syntax")
	}
}

func TestCompareIncludeAndIncludeRaw(t *testing.T) {
	// Setup with test files
	mockResolver := NewMockFileResolver("/base")

	// Add simple files without template syntax for direct comparison
	mockResolver.AddFile("/base/testdata/simple.md", "Simple content")
	mockResolver.AddFile("/base/testdata/with_template_syntax.md", "Content with {{.Value}} and {{agent}}")

	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})
	registry.Register(&agent.Roo{})

	engine := &Engine{
		FileResolver:       mockResolver,
		References:         make(map[string]string),
		AgentType:          "claude",
		absTemplateBaseDir: "/base",
		AgentRegistry:      registry,
		absCurrentFilePath: "/base/test.md",
	}

	// Setup expected glob results
	mockResolver.ExpectGlob("/base/testdata/with_template_syntax.md", []string{"/base/testdata/with_template_syntax.md"})

	// Test direct function calls
	includeFunc := engine.IncludeFunc(true).(func(...string) (string, error))
	includeRawFunc := engine.IncludeFunc(false).(func(...string) (string, error))

	// Get results from the template file with syntax
	includeResult, err := includeFunc("@/testdata/with_template_syntax.md")
	if err != nil {
		t.Fatalf("expected no error for include, got %v", err)
	}

	includeRawResult, err := includeRawFunc("@/testdata/with_template_syntax.md")
	if err != nil {
		t.Fatalf("expected no error for includeRaw, got %v", err)
	}

	// Verify - the important part is that includeRaw preserves template syntax
	// while include attempts to process it (which may fail or leave it as-is)
	if !strings.Contains(includeRawResult, "{{.Value}}") || !strings.Contains(includeRawResult, "{{agent}}") {
		t.Errorf("includeRaw should preserve template syntax as-is")
	}

	// Different from regular include in that it doesn't attempt to process templates
	if includeResult == includeRawResult && strings.Contains(includeResult, "{{") {
		t.Logf("Note: In this test, regular include didn't modify templates because they're processed at a different time")
	}
}

func TestCompareReferenceAndReferenceRaw(t *testing.T) {
	// Setup with test files
	mockResolver := NewMockFileResolver("/base")

	// Add a file with template syntax
	mockResolver.AddFile("/base/testdata/with_template_syntax.md",
		"Content with {{.Value}} and {{agent}}")

	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})
	registry.Register(&agent.Roo{})

	// Setup expected glob results - used by both engines
	mockResolver.ExpectGlob("/base/testdata/with_template_syntax.md", []string{"/base/testdata/with_template_syntax.md"})

	// Test both functions separately

	// First test with reference function
	engine1 := &Engine{
		FileResolver:       mockResolver,
		References:         make(map[string]string),
		AgentType:          "claude",
		absTemplateBaseDir: "/base",
		AgentRegistry:      registry,
		absCurrentFilePath: "/base/test.md",
	}

	refFunc := engine1.ReferenceFunc(true).(func(...string) (string, error))
	refResult, err := refFunc("@/testdata/with_template_syntax.md")
	if err != nil {
		t.Fatalf("expected no error for reference, got %v", err)
	}

	// Verify reference marker format
	if refResult != "[参考: testdata/with_template_syntax.md]" {
		t.Errorf("expected reference marker format, got %q", refResult)
	}

	// Check stored content for the regular reference
	refContent, exists := engine1.References["testdata/with_template_syntax.md"]
	if !exists {
		t.Error("reference content not stored for regular reference")
	}

	// Then test with referenceRaw function
	engine2 := &Engine{
		FileResolver:       mockResolver,
		References:         make(map[string]string),
		AgentType:          "claude",
		absTemplateBaseDir: "/base",
		AgentRegistry:      registry,
		absCurrentFilePath: "/base/test.md",
	}

	refRawFunc := engine2.ReferenceFunc(false).(func(...string) (string, error))
	refRawResult, err := refRawFunc("@/testdata/with_template_syntax.md")
	if err != nil {
		t.Fatalf("expected no error for referenceRaw, got %v", err)
	}

	// Verify reference marker format (should be the same)
	if refRawResult != "[参考: testdata/with_template_syntax.md]" {
		t.Errorf("expected reference marker format, got %q", refRawResult)
	}

	// Check stored content for the raw reference
	refRawContent, exists := engine2.References["testdata/with_template_syntax.md"]
	if !exists {
		t.Error("reference content not stored for raw reference")
	}

	// Now verify the key difference: raw reference preserves template syntax
	if !strings.Contains(refRawContent, "{{.Value}}") || !strings.Contains(refRawContent, "{{agent}}") {
		t.Error("raw reference did not preserve template syntax")
	}

	// Regular reference might preserve template syntax too depending on implementation
	// The important thing is that raw references MUST preserve it
	if refContent == refRawContent && strings.Contains(refContent, "{{") {
		t.Logf("Note: In this test, regular reference didn't modify templates because they're processed at a different time")
	}
}

func TestExecuteWithPath(t *testing.T) {
	// Setup
	mockResolver := NewMockFileResolver("/base")
	mockResolver.AddFile("/base/testdata/referenced.md", "Content from referenced file")
	mockResolver.AddFile("/base/testdata/main.md", "Main content with reference: {{reference \"@/testdata/referenced.md\"}}")

	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})
	registry.Register(&agent.Roo{})

	engine := &Engine{
		FileResolver:       mockResolver,
		References:         make(map[string]string),
		AgentType:          "claude",
		absTemplateBaseDir: "/base",
		AgentRegistry:      registry,
	}

	// Setup expected glob results
	mockResolver.ExpectGlob("/base/testdata/referenced.md", []string{"/base/testdata/referenced.md"})
	mockResolver.ExpectGlob("/base/testdata/main.md", []string{"/base/testdata/main.md"})

	// Test ExecuteWithPath directly
	result, err := engine.Execute("/base/testdata/main.md", "Content from {{agent}} with path: {{.Path}}", map[string]string{
		"Path": "/base/testdata/main.md",
	})

	// Verify
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check content was processed correctly
	if !strings.Contains(result, "Content from claude with path: /base/testdata/main.md") {
		t.Errorf("expected template to be processed with agent and path, got %q", result)
	}

	// Test path management by nesting templates
	complexContent := "{{include \"@/testdata/main.md\"}}"
	complexResult, err := engine.Execute("/base/nested.md", complexContent, nil)

	if err != nil {
		t.Fatalf("expected no error for nested template, got %v", err)
	}

	// Check that nested content was processed and paths were maintained properly
	if !strings.Contains(complexResult, "Main content with reference") {
		t.Errorf("expected nested template content to be included, got %q", complexResult)
	}

	// Check that references section is included
	if !strings.Contains(complexResult, "## References") {
		t.Errorf("expected references section, got %q", complexResult)
	}

	if !strings.Contains(complexResult, "testdata/referenced.md") {
		t.Errorf("expected referenced content path to be included, got %q", complexResult)
	}
}

func TestPathManagement(t *testing.T) {
	// Setup a more complex scenario to test path management
	mockResolver := NewMockFileResolver("/base")
	mockResolver.AddFile("/base/file1.md", "Content from file1 referencing {{include \"file2.md\"}}")
	mockResolver.AddFile("/base/file2.md", "Content from file2 with CurrentPath: {{include \"subdir/file3.md\"}}")
	mockResolver.AddFile("/base/subdir/file3.md", "Content from file3")

	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})
	registry.Register(&agent.Roo{})

	engine := &Engine{
		FileResolver:       mockResolver,
		References:         make(map[string]string),
		AgentType:          "claude",
		absTemplateBaseDir: "/base",
		AgentRegistry:      registry,
	}

	// Setup expected glob results for nested includes
	mockResolver.ExpectGlob("/base/file2.md", []string{"/base/file2.md"})
	mockResolver.ExpectGlob("/base/subdir/file3.md", []string{"/base/subdir/file3.md"})

	// Execute file1 which includes file2 which includes file3
	result, err := engine.ExecuteFile("/base/file1.md", nil)

	if err != nil {
		t.Fatalf("expected no error for nested includes, got %v", err)
	}

	// Check that all content is included properly
	expectedContent := "Content from file1 referencing Content from file2 with CurrentPath: Content from file3"
	if !strings.Contains(result, expectedContent) {
		t.Errorf("expected all files to be included properly, got %q", result)
	}
}

func TestIncludeWithGlobPattern(t *testing.T) {
	// Setup
	mockResolver := NewMockFileResolver("/base")
	mockResolver.AddFile("/base/testdata/glob1.md", "# Content from glob1\n\nThis is content from glob1.")
	mockResolver.AddFile("/base/testdata/glob2.md", "# Content from glob2\n\nThis is content from glob2.")
	mockResolver.AddFile("/base/testdata/glob3.md", "# Content from glob3\n\nThis is content from glob3.")
	mockResolver.AddFile("/base/testdata/other.txt", "Other content that should not be included")

	// Setup expected glob results
	mockResolver.ExpectGlob("/base/testdata/glob*.md", []string{
		"/base/testdata/glob1.md",
		"/base/testdata/glob2.md",
		"/base/testdata/glob3.md",
	})

	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})
	registry.Register(&agent.Roo{})

	engine := &Engine{
		FileResolver:       mockResolver,
		References:         make(map[string]string),
		AgentType:          "claude",
		absTemplateBaseDir: "/base",
		AgentRegistry:      registry,
		absCurrentFilePath: "/base/current.md",
	}

	// Test the include function with glob pattern
	fn := engine.IncludeFunc(true).(func(...string) (string, error))
	result, err := fn("testdata/glob*.md")

	// Verify
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// All three files should be included and separated by double newlines
	if !strings.Contains(result, "Content from glob1") {
		t.Errorf("expected result to contain content from glob1")
	}
	if !strings.Contains(result, "Content from glob2") {
		t.Errorf("expected result to contain content from glob2")
	}
	if !strings.Contains(result, "Content from glob3") {
		t.Errorf("expected result to contain content from glob3")
	}
	if strings.Contains(result, "Other content") {
		t.Errorf("result should not contain content from files not matching the glob pattern")
	}
}

func TestIncludeRawWithGlobPattern(t *testing.T) {
	// Setup
	mockResolver := NewMockFileResolver("/base")
	mockResolver.AddFile("/base/testdata/raw1.md", "# Raw Content 1\n\n{{.Value}}")
	mockResolver.AddFile("/base/testdata/raw2.md", "# Raw Content 2\n\n{{include \"another.md\"}}")

	// Setup expected glob results
	mockResolver.ExpectGlob("/base/testdata/raw*.md", []string{
		"/base/testdata/raw1.md",
		"/base/testdata/raw2.md",
	})

	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})
	registry.Register(&agent.Roo{})

	engine := &Engine{
		FileResolver:       mockResolver,
		References:         make(map[string]string),
		AgentType:          "claude",
		absTemplateBaseDir: "/base",
		AgentRegistry:      registry,
		absCurrentFilePath: "/base/current.md",
	}

	// Test the includeRaw function with glob pattern
	fn := engine.IncludeFunc(false).(func(...string) (string, error))
	result, err := fn("testdata/raw*.md")

	// Verify
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check that template syntax is preserved as-is (not processed)
	if !strings.Contains(result, "{{.Value}}") {
		t.Errorf("includeRaw should preserve template variable syntax")
	}

	if !strings.Contains(result, "{{include \"another.md\"}}") {
		t.Errorf("includeRaw should preserve include directive syntax")
	}
}

func TestReferenceWithGlobPattern(t *testing.T) {
	// Setup
	mockResolver := NewMockFileResolver("/base")
	mockResolver.AddFile("/base/testdata/ref1.md", "# Reference Content 1")
	mockResolver.AddFile("/base/testdata/ref2.md", "# Reference Content 2")

	// Setup expected glob results
	mockResolver.ExpectGlob("/base/testdata/ref*.md", []string{
		"/base/testdata/ref1.md",
		"/base/testdata/ref2.md",
	})

	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})
	registry.Register(&agent.Roo{})

	engine := &Engine{
		FileResolver:       mockResolver,
		References:         make(map[string]string),
		AgentType:          "claude",
		absTemplateBaseDir: "/base",
		AgentRegistry:      registry,
		absCurrentFilePath: "/base/current.md",
	}

	// Test the reference function with glob pattern
	fn := engine.ReferenceFunc(true).(func(...string) (string, error))
	result, err := fn("testdata/ref*.md")

	// Verify
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check that all reference markers are included
	if !strings.Contains(result, "testdata/ref1.md") {
		t.Errorf("expected result to contain reference marker for ref1.md")
	}
	if !strings.Contains(result, "testdata/ref2.md") {
		t.Errorf("expected result to contain reference marker for ref2.md")
	}

	// Check that both references are stored
	content1, exists := engine.References["testdata/ref1.md"]
	if !exists {
		t.Errorf("reference content for ref1.md was not stored")
	} else if !strings.Contains(content1, "Reference Content 1") {
		t.Errorf("stored reference content doesn't match expected for ref1.md")
	}

	content2, exists := engine.References["testdata/ref2.md"]
	if !exists {
		t.Errorf("reference content for ref2.md was not stored")
	} else if !strings.Contains(content2, "Reference Content 2") {
		t.Errorf("stored reference content doesn't match expected for ref2.md")
	}
}

func TestReferenceRawWithGlobPattern(t *testing.T) {
	// Setup
	mockResolver := NewMockFileResolver("/base")
	mockResolver.AddFile("/base/testdata/refraw1.md", "# Raw Reference 1\n\n{{.Value}}")
	mockResolver.AddFile("/base/testdata/refraw2.md", "# Raw Reference 2\n\n{{include \"another.md\"}}")

	// Setup expected glob results
	mockResolver.ExpectGlob("/base/testdata/refraw*.md", []string{
		"/base/testdata/refraw1.md",
		"/base/testdata/refraw2.md",
	})

	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})
	registry.Register(&agent.Roo{})

	engine := &Engine{
		FileResolver:       mockResolver,
		References:         make(map[string]string),
		AgentType:          "claude",
		absTemplateBaseDir: "/base",
		AgentRegistry:      registry,
		absCurrentFilePath: "/base/current.md",
	}

	// Test the referenceRaw function with glob pattern
	fn := engine.ReferenceFunc(false).(func(...string) (string, error))
	result, err := fn("testdata/refraw*.md")

	// Verify
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check that both reference markers are included
	if !strings.Contains(result, "testdata/refraw1.md") {
		t.Errorf("expected result to contain reference marker for refraw1.md")
	}
	if !strings.Contains(result, "testdata/refraw2.md") {
		t.Errorf("expected result to contain reference marker for refraw2.md")
	}

	// Check that template syntax is preserved as-is in stored references
	content1 := engine.References["testdata/refraw1.md"]
	if !strings.Contains(content1, "{{.Value}}") {
		t.Errorf("referenceRaw should preserve template variable syntax")
	}

	content2 := engine.References["testdata/refraw2.md"]
	if !strings.Contains(content2, "{{include \"another.md\"}}") {
		t.Errorf("referenceRaw should preserve include directive syntax")
	}
}

func TestGlobWithBasePathPrefix(t *testing.T) {
	// Setup
	mockResolver := NewMockFileResolver("/base")
	mockResolver.AddFile("/base/testdata/base1.md", "Content from base1")
	mockResolver.AddFile("/base/testdata/base2.md", "Content from base2")

	// Setup expected glob results for @/ prefixed path
	mockResolver.ExpectGlob("/base/testdata/base*.md", []string{
		"/base/testdata/base1.md",
		"/base/testdata/base2.md",
	})

	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})
	registry.Register(&agent.Roo{})

	engine := &Engine{
		FileResolver:       mockResolver,
		References:         make(map[string]string),
		AgentType:          "claude",
		absTemplateBaseDir: "/base",
		AgentRegistry:      registry,
		absCurrentFilePath: "/base/subdir/current.md",
	}

	// Test include with @/ prefix
	fn := engine.IncludeFunc(true).(func(...string) (string, error))
	result, err := fn("@/testdata/base*.md")

	// Verify
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Both files should be included
	if !strings.Contains(result, "Content from base1") {
		t.Errorf("expected result to contain content from base1.md")
	}
	if !strings.Contains(result, "Content from base2") {
		t.Errorf("expected result to contain content from base2.md")
	}
}

func TestGlobWithNegativePatterns(t *testing.T) {
	// Setup
	mockResolver := NewMockFileResolver("/base")
	mockResolver.AddFile("/base/testdata/include1.md", "Include 1")
	mockResolver.AddFile("/base/testdata/include2.md", "Include 2")
	mockResolver.AddFile("/base/testdata/exclude.md", "Should be excluded")

	// Setup expected glob results with negative pattern
	// The resolution of negative patterns happens in the Glob implementation,
	// but we can simulate the behavior by not including the excluded file in the results
	mockResolver.ExpectGlob("/base/testdata/include*.md", []string{
		"/base/testdata/include1.md",
		"/base/testdata/include2.md",
	})
	mockResolver.ExpectGlob("!/base/testdata/exclude.md", []string{})

	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})
	registry.Register(&agent.Roo{})

	engine := &Engine{
		FileResolver:       mockResolver,
		References:         make(map[string]string),
		AgentType:          "claude",
		absTemplateBaseDir: "/base",
		AgentRegistry:      registry,
		absCurrentFilePath: "/base/current.md",
	}

	// Test include with positive and negative patterns
	fn := engine.IncludeFunc(true).(func(...string) (string, error))
	result, err := fn("testdata/include*.md", "!testdata/exclude.md")

	// Verify
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Only non-excluded files should be included
	if !strings.Contains(result, "Include 1") {
		t.Errorf("expected result to contain content from include1.md")
	}
	if !strings.Contains(result, "Include 2") {
		t.Errorf("expected result to contain content from include2.md")
	}
	if strings.Contains(result, "Should be excluded") {
		t.Errorf("result should not contain content from excluded file")
	}
}
