package util

import (
	"os"
	"path/filepath"
)

// FileSystem provides file system operations
type FileSystem interface {
	// ReadFile reads a file at the given path
	ReadFile(path string) ([]byte, error)

	// WriteFile writes content to a file
	WriteFile(path string, data []byte) error

	// FileExists checks if a file exists
	FileExists(path string) bool

	// IsFile checks if a path exists and is a regular file
	IsFile(path string) bool

	// IsDir checks if a path exists and is a directory
	IsDir(path string) bool

	// ListFiles lists files matching a pattern
	ListFiles(dir string, pattern string) ([]string, error)

	// ResolvePath resolves a relative path to absolute
	ResolvePath(path string) string

	// GlobWithExcludes expands glob patterns with support for exclusions
	GlobWithExcludes(patterns []string, baseDir string) ([]string, error)
}

// RealFileSystem implements FileSystem with the actual OS file system
type RealFileSystem struct{}

// ReadFile reads a file at the given path
func (fs *RealFileSystem) ReadFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, &ErrFileNotFound{Path: path}
		}
		return nil, WrapError(err, "failed to read file")
	}
	return data, nil
}

// WriteFile writes content to a file
func (fs *RealFileSystem) WriteFile(path string, data []byte) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return WrapError(err, "failed to create directory")
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return WrapError(err, "failed to write file")
	}

	return nil
}

// FileExists checks if a file exists
func (fs *RealFileSystem) FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// ListFiles lists files matching a pattern
func (fs *RealFileSystem) ListFiles(dir string, pattern string) ([]string, error) {
	matches, err := filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		return nil, WrapError(err, "failed to list files")
	}
	return matches, nil
}

// ResolvePath resolves a relative path to absolute
func (fs *RealFileSystem) ResolvePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	// Assuming the current working directory as the base
	abs, err := filepath.Abs(path)
	if err != nil {
		return path // Return original on error
	}
	return abs
}

// GlobWithExcludes expands glob patterns with support for exclusions
func (fs *RealFileSystem) GlobWithExcludes(patterns []string, baseDir string) ([]string, error) {
	return GlobWithExcludes(patterns, baseDir)
}

// IsFile checks if a path exists and is a regular file
func (fs *RealFileSystem) IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// IsDir checks if a path exists and is a directory
func (fs *RealFileSystem) IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
