// Package processor provides functionality for processing agent-def tasks
package processor

import (
	"github.com/uphy/agent-def/internal/agent"
)

// OutputConfig encapsulates output configuration settings
type OutputConfig struct {
	Agent       agent.Agent
	Path        string
	IsDirectory bool
	AgentName   string // Original agent name from config
}

// ProcessedFile represents a processed output file
type ProcessedFile struct {
	Path    string
	Content string
}

// TaskResult represents the result of processing a task
type TaskResult struct {
	Files []ProcessedFile
}
