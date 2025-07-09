package model

import (
	"github.com/user/agent-def/internal/util"
)

// Context represents the context definition
type Context struct {
	// Content is the raw content of the context file
	Content string

	// Path is the original file path
	Path string

	// Processed is the processed content after template execution
	Processed string

	// Metadata contains any extracted metadata from the context
	Metadata map[string]interface{}
}

// NewContext creates a new Context instance
func NewContext(path string, content string) *Context {
	return &Context{
		Path:     path,
		Content:  content,
		Metadata: make(map[string]interface{}),
	}
}

// LoadContext loads a context from a file
func LoadContext(fs util.FileSystem, path string) (*Context, error) {
	content, err := fs.ReadFile(path)
	if err != nil {
		return nil, util.WrapError(err, "failed to read context file")
	}

	return NewContext(path, string(content)), nil
}

// ProcessContext processes a context using the template engine
// This maintains backward compatibility with existing tests
func ProcessContext(engine interface{}, content string) (Context, error) {
	// For backward compatibility with tests using nil engine
	if engine == nil {
		return Context{
			Content:   content,
			Processed: content,
		}, nil
	}

	// When a real engine is provided
	ctx := NewContext("", content)
	err := ProcessContextWithEngine(engine, ctx)
	return *ctx, err
}

// ProcessContextWithEngine processes a context using the template engine
// The engine parameter should be *template.Engine but using interface{} to avoid import cycles
func ProcessContextWithEngine(engine interface{}, ctx *Context) error {
	if engine == nil {
		// No processing needed for nil engine, just copy content to processed
		ctx.Processed = ctx.Content
		return nil
	}

	templateEngine, ok := engine.(TemplateProcessor)
	if !ok {
		return util.WrapError(nil, "invalid template engine provided")
	}

	processed, err := templateEngine.Process(ctx.Content)
	if err != nil {
		return util.WrapError(err, "failed to process context template")
	}

	ctx.Processed = processed
	return nil
}

// TemplateProcessor is an interface for template processing
type TemplateProcessor interface {
	Process(content string) (string, error)
}
