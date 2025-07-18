package template

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/uphy/agent-def/internal/agent"
	"github.com/uphy/agent-def/internal/util"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

	// OriginalPaths maps absolute paths to their original relative paths
	OriginalPaths map[string]string

	// AgentType determines the output format
	AgentType string

	// BasePath is the base path for resolving relative paths
	BasePath string

	// CurrentFilePath is the path of the file currently being processed
	CurrentFilePath string

	// AgentRegistry provides access to registered agents
	AgentRegistry *agent.Registry
}

// Context holds the context information for template processing,
// including information about which agent is being used for output.
type Context struct {
	Agent string
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

// NewEngineWithFilePath creates a new template engine with the current file path
func NewEngineWithFilePath(fileResolver FileResolver, agentType, basePath, currentFilePath string, agentRegistry *agent.Registry) *Engine {
	return &Engine{
		FileResolver:    fileResolver,
		References:      make(map[string]string),
		AgentType:       agentType,
		BasePath:        basePath,
		CurrentFilePath: currentFilePath,
		AgentRegistry:   agentRegistry,
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

// ExecuteFile processes a template file with the given data
func (e *Engine) ExecuteFile(filePath string, data interface{}) (string, error) {
	// Save previous current file path
	previousPath := e.CurrentFilePath

	// Set current file path to the file being processed
	e.CurrentFilePath = filePath

	// Read the file
	content, err := e.FileResolver.Read(filePath)
	if err != nil {
		// Restore previous path before returning error
		e.CurrentFilePath = previousPath
		return "", &util.ErrFileNotFound{Path: filePath}
	}

	// Execute the template
	result, err := e.Execute(string(content), data)

	// Restore previous current file path before returning
	e.CurrentFilePath = previousPath

	return result, err
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
		"agent":     e.Agent,
	}
	caser := cases.Title(language.English)
	for _, agent := range e.AgentRegistry.List() {
		name := caser.String(agent.ID())
		funcMap["is"+name] = func() bool {
			return agent.ID() == e.AgentType
		}
		funcMap["if"+name] = func(args ...string) string {
			ifStr := args[0]
			elseStr := ""
			switch len(args) {
			case 1:
				// No else part
			case 2:
				elseStr = args[1]
			default:
				return fmt.Sprintf("if%s: too many arguments", name)
			}

			if agent.ID() == e.AgentType {
				return ifStr
			}
			return elseStr
		}
	}
	return t.Funcs(funcMap)
}

// processInclude processes an included file recursively
func (e *Engine) processInclude(path string) (string, error) {
	// Save previous current file path
	previousPath := e.CurrentFilePath

	// Set current file path to the file being included
	e.CurrentFilePath = path

	// Read the file
	content, err := e.FileResolver.Read(path)
	if err != nil {
		// Restore previous path before returning error
		e.CurrentFilePath = previousPath
		return "", &util.ErrFileNotFound{Path: path}
	}

	// Process the content as a template recursively
	processed, err := e.Execute(string(content), nil)
	if err != nil {
		// Restore previous path before returning error
		e.CurrentFilePath = previousPath
		return "", &util.ErrTemplateExecution{
			Template: path,
			Cause:    err,
		}
	}

	// Restore previous current file path before returning
	e.CurrentFilePath = previousPath

	return processed, nil
}

// processReference adds a reference to be appended at the end
func (e *Engine) processReference(originalPath, fullPath string) (string, error) {
	// Save previous current file path
	previousPath := e.CurrentFilePath

	// Set current file path to the file being referenced
	e.CurrentFilePath = fullPath

	// Process the file content recursively
	content, err := e.processInclude(fullPath)
	if err != nil {
		// Restore previous path before returning error
		e.CurrentFilePath = previousPath
		return "", err
	}

	// Store the content to be appended at the end
	e.References[originalPath] = content

	// Restore previous current file path before returning
	e.CurrentFilePath = previousPath

	// Return a reference marker with original relative path
	return fmt.Sprintf("[参考: %s]", originalPath), nil
}
