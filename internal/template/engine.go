package template

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/user/agent-def/internal/agent"
	"github.com/user/agent-def/internal/util"
)

// FileResolver resolves and reads files
type FileResolver interface {
	// Read reads a file at the given path
	Read(path string) ([]byte, error)

	// Exists checks if a file exists
	Exists(path string) bool

	// ResolvePath resolves a relative path to absolute
	ResolvePath(path string) string
}

// Engine represents the template engine wrapper
type Engine struct {
	// FileResolver resolves file paths
	FileResolver FileResolver

	// References stores reference content to be appended at the end
	References map[string]string

	// AgentType determines the output format
	AgentType string

	// BasePath is the base path for resolving relative paths
	BasePath string

	// AgentRegistry provides access to registered agents
	AgentRegistry *agent.Registry
}

// NewEngine creates a new template engine
func NewEngine(fileResolver FileResolver, agentType, basePath string, agentRegistry *agent.Registry) *Engine {
	return &Engine{
		FileResolver:  fileResolver,
		References:    make(map[string]string),
		AgentType:     agentType,
		BasePath:      basePath,
		AgentRegistry: agentRegistry,
	}
}

// Execute processes a template with the given data
func (e *Engine) Execute(content string, data interface{}) (string, error) {
	// Create a new template
	t := template.New("template").Delims("{{", "}}")

	// Register helper functions
	t = e.RegisterHelperFunctions(t)

	// Parse the template
	t, err := t.Parse(content)
	if err != nil {
		return "", &util.ErrTemplateExecution{
			Template: "content",
			Cause:    err,
		}
	}

	// Execute the template
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", &util.ErrTemplateExecution{
			Template: "content",
			Cause:    err,
		}
	}

	result := buf.String()

	// Append references if any
	if len(e.References) > 0 {
		result += "\n\n## References\n\n"
		for path, content := range e.References {
			result += fmt.Sprintf("### %s\n\n%s\n\n", path, content)
		}
	}

	return result, nil
}

// Process implements the model.TemplateProcessor interface
func (e *Engine) Process(content string) (string, error) {
	return e.Execute(content, nil)
}

// RegisterHelperFunctions registers all template helper functions
func (e *Engine) RegisterHelperFunctions(t *template.Template) *template.Template {
	funcMap := template.FuncMap{
		"file":      e.FileFunc(),
		"include":   e.IncludeFunc(),
		"reference": e.ReferenceFunc(),
		"mcp":       e.MCPFunc(),
	}
	return t.Funcs(funcMap)
}

// processInclude processes an included file recursively
func (e *Engine) processInclude(path string) (string, error) {
	// Read the file
	content, err := e.FileResolver.Read(path)
	if err != nil {
		return "", &util.ErrFileNotFound{Path: path}
	}

	// Process the content as a template recursively
	processed, err := e.Execute(string(content), nil)
	if err != nil {
		return "", &util.ErrTemplateExecution{
			Template: path,
			Cause:    err,
		}
	}

	return processed, nil
}

// processReference adds a reference to be appended at the end
func (e *Engine) processReference(path string) (string, error) {
	// Process the file content recursively
	content, err := e.processInclude(path)
	if err != nil {
		return "", err
	}

	// Store the content to be appended at the end
	e.References[path] = content

	// Return a reference marker
	return fmt.Sprintf("[参考: %s]", path), nil
}
