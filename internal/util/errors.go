package util

import "fmt"

// CustomError defines the interface for custom error types
// that can provide formatted error messages
type CustomError interface {
	error
	FormattedError() string
}

// ErrFileNotFound represents a file not found error
type ErrFileNotFound struct {
	Path string
}

func (e *ErrFileNotFound) Error() string {
	return fmt.Sprintf("file not found: %s", e.Path)
}

func (e *ErrFileNotFound) FormattedError() string {
	return fmt.Sprintf("Could not find file at path: %s", e.Path)
}

// ErrTemplateExecution represents a template execution error
type ErrTemplateExecution struct {
	Template string
	Cause    error
}

func (e *ErrTemplateExecution) Error() string {
	return fmt.Sprintf("template execution failed for '%s': %v", e.Template, e.Cause)
}

func (e *ErrTemplateExecution) FormattedError() string {
	return fmt.Sprintf("Failed to execute template '%s': %v", e.Template, e.Cause)
}

// ErrInvalidAgent represents an unknown agent type error
type ErrInvalidAgent struct {
	Type string
}

func (e *ErrInvalidAgent) Error() string {
	return fmt.Sprintf("unknown agent type: %s", e.Type)
}

func (e *ErrInvalidAgent) FormattedError() string {
	return fmt.Sprintf("Invalid agent type: '%s'. Use 'agent-def list' to see available agents.", e.Type)
}

// ErrParseFailure represents a parsing error
type ErrParseFailure struct {
	Path  string
	Cause error
}

func (e *ErrParseFailure) Error() string {
	return fmt.Sprintf("failed to parse '%s': %v", e.Path, e.Cause)
}

func (e *ErrParseFailure) FormattedError() string {
	return fmt.Sprintf("Failed to parse file '%s': %v", e.Path, e.Cause)
}

// ErrMalformedFrontmatter represents a malformed frontmatter error
type ErrMalformedFrontmatter struct {
	Path  string
	Cause error
}

func (e *ErrMalformedFrontmatter) Error() string {
	return fmt.Sprintf("malformed frontmatter in '%s': %v", e.Path, e.Cause)
}

func (e *ErrMalformedFrontmatter) FormattedError() string {
	return fmt.Sprintf("Malformed frontmatter in file '%s': %v", e.Path, e.Cause)
}

// ErrInvalidOutputFormat represents an invalid output format error
type ErrInvalidOutputFormat struct {
	Format string
}

func (e *ErrInvalidOutputFormat) Error() string {
	return fmt.Sprintf("invalid output format: %s", e.Format)
}

func (e *ErrInvalidOutputFormat) FormattedError() string {
	return fmt.Sprintf("Invalid output format: '%s'. Supported formats are 'json', 'yaml', and 'text'.", e.Format)
}

// ErrInvalidConfig represents an invalid configuration error
type ErrInvalidConfig struct {
	Reason string
}

func (e *ErrInvalidConfig) Error() string {
	return fmt.Sprintf("invalid configuration: %s", e.Reason)
}

func (e *ErrInvalidConfig) FormattedError() string {
	return fmt.Sprintf("Configuration error: %s", e.Reason)
}

// WrapError wraps an error with additional context
func WrapError(err error, message string) error {
	return fmt.Errorf("%s: %w", message, err)
}
