// Package processor provides functionality for processing agent-def tasks
package processor

import (
	"path/filepath"
	"strings"
)

// IsDirectory determines if a path represents a directory
// by checking if it ends with a path separator
func IsDirectory(path string) bool {
	return strings.HasSuffix(path, string(filepath.Separator)) ||
		strings.HasSuffix(path, "/")
}

// RemoveFileExtension returns the filename without its extension
func RemoveFileExtension(path string) string {
	ext := filepath.Ext(path)
	return strings.TrimSuffix(path, ext)
}
