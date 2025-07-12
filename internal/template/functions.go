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
		fullPath := e.resolveTemplatePath(path)

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
		fullPath := e.resolveTemplatePath(path)

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

func (e *Engine) Agent() string {
	return e.AgentType
}

// resolveTemplatePath resolves paths for template includes and references
// using the special template path resolution rules:
// 1. Paths starting with "/" are relative to agent-def.yml's directory (BasePath) with the leading slash removed
// 2. Paths starting with "./" or "../" are relative to the including file's directory
// 3. Other paths (without "./" or "../" prefix) are relative to agent-def.yml's directory (BasePath)
// 4. OS-absolute paths (like C:\ on Windows) are preserved as-is
func (e *Engine) resolveTemplatePath(path string) string {
	// 1. Handle "/"-prefixed paths as relative to BasePath
	if strings.HasPrefix(path, "/") {
		trimmed := strings.TrimPrefix(path, "/")
		return filepath.Join(e.BasePath, trimmed)
	}

	// 2. Handle Windows-style absolute paths
	// This is specifically for Windows where C:\ style paths are absolute
	// On Unix, this should never trigger because / paths are handled above
	if filepath.IsAbs(path) && path[0] != '/' {
		return path
	}

	// 3. Handle explicitly relative paths (./ or ../) relative to current file
	if strings.HasPrefix(path, "./") || strings.HasPrefix(path, "../") {
		if e.CurrentFilePath != "" {
			return filepath.Join(filepath.Dir(e.CurrentFilePath), path)
		}
	}

	// 4. All other paths are relative to BasePath
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
