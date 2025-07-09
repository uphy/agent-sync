package util

import (
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
