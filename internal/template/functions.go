package template

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/uphy/agent-def/internal/util"
)

// FileFunc generates a file reference helper function
func (e *Engine) FileFunc() any {
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

// IncludeFunc generates an include helper function with optional template processing
func (e *Engine) IncludeFunc(processTemplate bool) any {
	return func(path string) (string, error) {
		// Resolve path relative to the base path
		fullPath, err := e.resolveTemplatePath(path)
		if err != nil {
			return "", err
		}

		// Check if file exists
		if !e.FileResolver.Exists(fullPath) {
			return "", &util.ErrFileNotFound{Path: fullPath}
		}

		// Process the include with or without template processing
		return e.processInclude(fullPath, processTemplate)
	}
}

// ReferenceFunc generates a reference helper function with optional template processing
func (e *Engine) ReferenceFunc(processTemplate bool) any {
	return func(path string) (string, error) {
		// Resolve path relative to the base path
		fullPath, err := e.resolveTemplatePath(path)
		if err != nil {
			return "", err
		}

		// Check if file exists
		if !e.FileResolver.Exists(fullPath) {
			return "", &util.ErrFileNotFound{Path: fullPath}
		}

		// Process the reference with or without template processing
		return e.processReference(path, fullPath, processTemplate)
	}
}

// MCPFunc generates an MCP command helper function
func (e *Engine) MCPFunc() any {
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

// Agent returns the current agent type being used for template processing.
func (e *Engine) Agent() string {
	return e.AgentType
}

// resolveTemplatePath resolves paths for template includes and references
// using the special template path resolution rules:
// 1. Absolute paths (like C:\ on Windows or /root on Unix) are not allowed and will return an error
// 2. Paths starting with "@/" are relative to agent-def.yml's directory (BasePath) with the "@/" prefix removed
// 3. All other paths (including those with "./" or "../" prefix) are relative to the including file's directory (CurrentFilePath)
//
// Returns the resolved absolute path and an error if the path is invalid.
func (e *Engine) resolveTemplatePath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return "", fmt.Errorf("absolute paths are not allowed in template includes or references: %s", path)
	}

	// Handle "@/"-prefixed paths as relative to BasePath
	if trimmed, ok := strings.CutPrefix(path, "@/"); ok {
		return filepath.Join(e.BasePath, trimmed), nil
	}

	// Handle explicitly relative paths (./ or ../) relative to current file
	return filepath.Join(filepath.Dir(e.CurrentFilePath), path), nil
}

// normalizePath normalizes a path for the current OS
func (e *Engine) normalizePath(path string) string {
	// Replace backslashes with forward slashes on Windows
	if filepath.Separator == '\\' {
		return strings.ReplaceAll(path, "\\", "/")
	}
	return path
}
