package template

import (
	"path/filepath"
	"strings"

	"github.com/user/agent-def/internal/util"
)

// FileFunc generates a file reference helper function
func (e *Engine) FileFunc() interface{} {
	return func(path string) (string, error) {
		// Get the agent implementation
		agent, found := e.AgentRegistry.Get(e.AgentType)
		if !found {
			return "", &util.ErrInvalidAgent{Type: e.AgentType}
		}

		// Format the path according to the agent
		return agent.FormatFile(path), nil
	}
}

// IncludeFunc generates an include helper function
func (e *Engine) IncludeFunc() interface{} {
	return func(path string) (string, error) {
		// Resolve path relative to the base path
		fullPath := e.resolveRelativePath(path)

		// Check if file exists
		if !e.FileResolver.Exists(fullPath) {
			return "", &util.ErrFileNotFound{Path: fullPath}
		}

		// Process the include recursively
		return e.processInclude(fullPath)
	}
}

// ReferenceFunc generates a reference helper function
func (e *Engine) ReferenceFunc() interface{} {
	return func(path string) (string, error) {
		// Resolve path relative to the base path
		fullPath := e.resolveRelativePath(path)

		// Check if file exists
		if !e.FileResolver.Exists(fullPath) {
			return "", &util.ErrFileNotFound{Path: fullPath}
		}

		// Process the reference
		return e.processReference(fullPath)
	}
}

// MCPFunc generates an MCP command helper function
func (e *Engine) MCPFunc() interface{} {
	return func(agentName, command string, args ...string) (string, error) {
		// Get the agent implementation
		agent, found := e.AgentRegistry.Get(e.AgentType)
		if !found {
			return "", &util.ErrInvalidAgent{Type: e.AgentType}
		}

		// Format the MCP command according to the agent
		return agent.FormatMCP(agentName, command, args...), nil
	}
}

// resolveRelativePath resolves a path relative to the base path
func (e *Engine) resolveRelativePath(path string) string {
	// If path is absolute, return as is
	if filepath.IsAbs(path) {
		return path
	}

	// Otherwise, join with base path
	return filepath.Join(e.BasePath, path)
}

// normalizePath normalizes a path for the current OS
func (e *Engine) normalizePath(path string) string {
	// Replace backslashes with forward slashes on Windows
	if filepath.Separator == '\\' {
		return strings.ReplaceAll(path, "\\", "/")
	}
	return path
}
