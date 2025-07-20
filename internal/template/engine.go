package template

import (
	"bytes"
	"fmt"
	"path/filepath"
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

	// absTemplateBaseDir is the base path for resolving relative paths
	absTemplateBaseDir string

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
func NewEngine(fileResolver FileResolver, agentType, absTemplateBaseDir string, agentRegistry *agent.Registry) *Engine {
	return &Engine{
		FileResolver:       fileResolver,
		References:         make(map[string]string),
		AgentType:          agentType,
		absTemplateBaseDir: absTemplateBaseDir,
		AgentRegistry:      agentRegistry,
	}
}

// Execute processes a template with the given data while managing CurrentFilePath
func (e *Engine) Execute(path string, content string, data any) (string, error) {
	// Process the template content with proper path management
	result, err := e.executeWithoutReferencesWithPath(path, content, data)
	if err != nil {
		return "", err
	}

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
func (e *Engine) ExecuteFile(absFilePath string, data any) (string, error) {
	if !filepath.IsAbs(absFilePath) {
		return "", fmt.Errorf("absolute file path is required: %s", absFilePath)
	}

	// Read the file
	content, err := e.FileResolver.Read(absFilePath)
	if err != nil {
		return "", &util.ErrFileNotFound{Path: absFilePath}
	}

	// Use the new ExecuteWithPath method which handles path management
	return e.Execute(absFilePath, string(content), data)
}

// Execute processes a template with the given data
func (e *Engine) executeWithoutReferences(content string, data any) (string, error) {
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

	return buf.String(), nil
}

// executeWithoutReferencesWithPath processes a template with the given data, ensuring CurrentFilePath
// is properly managed during execution. This method uses defer to guarantee path restoration.
func (e *Engine) executeWithoutReferencesWithPath(path string, content string, data any) (string, error) {
	// Save the current path
	previousPath := e.CurrentFilePath

	// Set the current path to the path being processed
	e.CurrentFilePath = path

	// Use defer to ensure path restoration happens regardless of execution result
	defer func() {
		e.CurrentFilePath = previousPath
	}()

	// Execute the template
	return e.executeWithoutReferences(content, data)
}

// RegisterHelperFunctions registers all template helper functions
func (e *Engine) RegisterHelperFunctions(t *template.Template) *template.Template {
	funcMap := template.FuncMap{
		"file":         e.FileFunc(),
		"include":      e.IncludeFunc(true),
		"reference":    e.ReferenceFunc(true),
		"includeRaw":   e.IncludeFunc(false),
		"referenceRaw": e.ReferenceFunc(false),
		"mcp":          e.MCPFunc(),
		"agent":        e.Agent,
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

// processInclude processes an included file with optional template processing
func (e *Engine) processInclude(path string, processTemplate bool) (string, error) {
	// Read the file
	content, err := e.FileResolver.Read(path)
	if err != nil {
		return "", &util.ErrFileNotFound{Path: path}
	}

	if processTemplate {
		// Process the content as a template recursively using path management
		processed, err := e.executeWithoutReferencesWithPath(path, string(content), nil)
		if err != nil {
			return "", &util.ErrTemplateExecution{
				Template: path,
				Cause:    err,
			}
		}
		return processed, nil
	} else {
		// Return the raw content without template processing
		// No need for path management as we're not processing templates
		return string(content), nil
	}
}

// processReference adds a reference to be appended at the end
func (e *Engine) processReference(fullPath string, processTemplate bool) (string, error) {
	originalPath, err := filepath.Rel(e.absTemplateBaseDir, fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to get relative path for %q: %w", fullPath, err)
	}

	// Process the file content with or without template processing
	// The path management will be handled by processInclude
	content, err := e.processInclude(fullPath, processTemplate)
	if err != nil {
		return "", err
	}

	// Store the content to be appended at the end
	e.References[originalPath] = content

	// Return a reference marker with original relative path
	return fmt.Sprintf("[参考: %s]", originalPath), nil
}
