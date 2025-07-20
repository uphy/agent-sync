// Package processor provides functionality for processing agent-sync tasks
package processor

import (
	"github.com/uphy/agent-sync/internal/agent"
	"github.com/uphy/agent-sync/internal/util"
	"go.uber.org/zap"
)

// TaskProcessor defines the interface for processing different types of tasks
type TaskProcessor interface {
	// Process handles task-specific processing of inputs and returns the processed result
	Process(inputs []string, cfg *OutputConfig) (*TaskResult, error)

	// GetOutputPath returns the appropriate output path for the given agent and task type
	GetOutputPath(agent agent.Agent, outputPath string) string
}

// BaseProcessor contains common functionality for all task processors
type BaseProcessor struct {
	fs           util.FileSystem
	logger       *zap.Logger
	absInputRoot string
	registry     *agent.Registry
	userScope    bool
}

// NewBaseProcessor creates a new BaseProcessor with the given parameters
func NewBaseProcessor(fs util.FileSystem, logger *zap.Logger, absInputRoot string, registry *agent.Registry, userScope bool) *BaseProcessor {
	return &BaseProcessor{
		fs:           fs,
		logger:       logger,
		absInputRoot: absInputRoot,
		registry:     registry,
		userScope:    userScope,
	}
}
