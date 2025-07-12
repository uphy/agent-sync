package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// JoinPath joins path elements into a single path, handling both absolute and relative paths
func JoinPath(base string, paths ...string) string {
	if len(paths) == 0 {
		return base
	}

	// If the first path is absolute, ignore the base
	if filepath.IsAbs(paths[0]) {
		return filepath.Join(paths...)
	}

	// Join base with all paths
	allPaths := append([]string{base}, paths...)
	return filepath.Join(allPaths...)
}

// GetBasePath returns the base directory of a file path
func GetBasePath(path string) string {
	return filepath.Dir(path)
}

// GetRelativePath returns a path relative to a base directory
func GetRelativePath(basePath, path string) (string, error) {
	rel, err := filepath.Rel(basePath, path)
	if err != nil {
		return "", WrapError(err, "failed to get relative path")
	}
	return rel, nil
}

// NormalizePath normalizes a path for the current platform
func NormalizePath(path string) string {
	// Replace backslashes with forward slashes on Windows
	if os.PathSeparator == '\\' {
		path = strings.ReplaceAll(path, "\\", "/")
	}
	return path
}

// IsDirectory returns true if the path is a directory
func IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// EnsureDirectory ensures a directory exists, creating it if necessary
func EnsureDirectory(path string) error {
	if IsDirectory(path) {
		return nil
	}

	err := os.MkdirAll(path, 0755)
	if err != nil {
		return WrapError(err, "failed to create directory")
	}
	return nil
}

// ExpandTilde replaces a leading tilde (~) in a path with the user's home directory.
// If the path doesn't start with a tilde, it is returned unchanged.
// Returns an error if the home directory cannot be determined.
func ExpandTilde(path string) (string, error) {
	// Return unchanged if the path is empty or doesn't start with ~
	if path == "" || path[0] != '~' {
		return path, nil
	}

	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Replace ~ with the home directory
	if len(path) == 1 {
		// Just ~ by itself
		return homeDir, nil
	} else if path[1] == filepath.Separator {
		// ~/something - path with separator
		return filepath.Join(homeDir, path[2:]), nil
	}
	// ~something - path without separator (less common)
	return filepath.Join(homeDir, path[1:]), nil
}
