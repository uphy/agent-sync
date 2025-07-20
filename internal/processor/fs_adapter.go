// Package processor provides functionality for processing agent-sync tasks
package processor

import (
	"github.com/uphy/agent-sync/internal/util"
)

// FSAdapter bridges util.FileSystem to template.FileResolver.
type FSAdapter struct {
	fs util.FileSystem
}

// NewFSAdapter creates a new FSAdapter
func NewFSAdapter(fs util.FileSystem) *FSAdapter {
	return &FSAdapter{fs: fs}
}

func (a *FSAdapter) Read(path string) ([]byte, error) {
	return a.fs.ReadFile(path)
}

func (a *FSAdapter) Exists(path string) bool {
	return a.fs.FileExists(path)
}

func (a *FSAdapter) ResolvePath(path string) string {
	return a.fs.ResolvePath(path)
}

func (a *FSAdapter) Glob(patterns []string) ([]string, error) {
	return util.GlobWithExcludesNoBaseDir(patterns)
}
