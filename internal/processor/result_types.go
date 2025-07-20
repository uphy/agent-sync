// Package processor provides functionality for processing agent-sync tasks
package processor

import (
	"github.com/uphy/agent-sync/internal/agent"
)

// OutputConfig encapsulates output configuration settings
type OutputConfig struct {
	Agent agent.Agent
	// RelPath is the relative path from the output directory or file where the processed content will be written
	// If the path ends with "/", it indicates a directory output
	// If it does not end with "/", it indicates a file output
	RelPath     string
	IsDirectory bool
	AgentName   string // Original agent name from config
}

// ProcessedFile represents a processed output file
type ProcessedFile struct {
	// relPath is the relative path of the file from the output directory
	relPath string
	Content string
}

// TaskResult represents the result of processing a task
type TaskResult struct {
	Files []ProcessedFile
}
