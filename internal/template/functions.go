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
// that accepts multiple paths and combines their content with double newlines as separators
func (e *Engine) IncludeFunc(processTemplate bool) any {
	return func(paths ...string) (string, error) {
		if len(paths) == 0 {
			return "", fmt.Errorf("at least one path must be provided")
		}

		// For multiple paths, resolve all paths
		resolvedPaths, err := e.resolveTemplatePath(paths)
		if err != nil {
			return "", err
		}

		// Process each file and collect the results
		var results []string
		for i, fullPath := range resolvedPaths {
			// Check if file exists
			if !e.FileResolver.Exists(fullPath) {
				return "", &util.ErrFileNotFound{Path: fullPath}
			}

			// Process the include with or without template processing
			content, err := e.processInclude(fullPath, processTemplate)
			if err != nil {
				return "", fmt.Errorf("failed to process include for path %q: %w", paths[i], err)
			}

			results = append(results, content)
		}

		// Combine the results with double newlines as separators
		return strings.Join(results, "\n\n"), nil
	}
}

// ReferenceFunc generates a reference helper function with optional template processing
// that accepts multiple paths and combines their reference markers with commas as separators
func (e *Engine) ReferenceFunc(processTemplate bool) any {
	return func(paths ...string) (string, error) {
		if len(paths) == 0 {
			return "", fmt.Errorf("at least one path must be provided")
		}

		// For multiple paths, resolve all paths
		resolvedPaths, err := e.resolveTemplatePath(paths)
		if err != nil {
			return "", err
		}

		// Process each file and collect the reference markers
		var refMarkers []string
		for i, resolvedPath := range resolvedPaths {
			// Check if file exists
			if !e.FileResolver.Exists(resolvedPath) {
				return "", &util.ErrFileNotFound{Path: resolvedPath}
			}

			// Process the reference with or without template processing
			// This will store the content in e.References and return a reference marker
			refMarker, err := e.processReference(resolvedPath, processTemplate)
			if err != nil {
				return "", fmt.Errorf("failed to process reference for path %q: %w", paths[i], err)
			}
			refMarkers = append(refMarkers, refMarker)
		}

		// Return a combined reference marker with comma-separated paths
		return strings.Join(refMarkers, ", "), nil
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
// Returns an array of resolved absolute paths and an error if any path is invalid.
func (e *Engine) resolveTemplatePath(paths []string) ([]string, error) {
	resolved := make([]string, 0, len(paths))

	for _, path := range paths {
		if filepath.IsAbs(path) {
			return nil, fmt.Errorf("absolute paths are not allowed in template includes or references: %s", path)
		}

		// Handle "@/"-prefixed paths as relative to BasePath
		if trimmed, ok := strings.CutPrefix(path, "@/"); ok {
			resolved = append(resolved, filepath.Join(e.BasePath, trimmed))
			continue
		}

		// Handle explicitly relative paths (./ or ../) relative to current file
		resolved = append(resolved, filepath.Join(filepath.Dir(e.CurrentFilePath), path))
	}

	return resolved, nil
}

// resolveTemplatePathSingle is a helper method that resolves a single path
// using the same rules as resolveTemplatePath
func (e *Engine) resolveTemplatePathSingle(path string) (string, error) {
	resolved, err := e.resolveTemplatePath([]string{path})
	if err != nil {
		return "", err
	}
	return resolved[0], nil
}

// normalizePath normalizes a path for the current OS
func (e *Engine) normalizePath(path string) string {
	// Replace backslashes with forward slashes on Windows
	if filepath.Separator == '\\' {
		return strings.ReplaceAll(path, "\\", "/")
	}
	return path
}
