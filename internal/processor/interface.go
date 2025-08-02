// Package processor provides functionality for processing agent-sync tasks
package processor

import (
	"fmt"
	"path/filepath"

	"github.com/uphy/agent-sync/internal/agent"
	"github.com/uphy/agent-sync/internal/template"
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

// Template engine factory kept internal to avoid repetition
func (p *BaseProcessor) templateEngine(agentName string) *template.Engine {
	fsAdapter := NewFSAdapter(p.fs)
	return template.NewEngine(fsAdapter, agentName, p.absInputRoot, p.registry)
}

// resolveOutputRelPath builds the per-input relative output path under cfg.RelPath
func resolveOutputRelPath(baseRel string, input string) string {
	return filepath.Join(baseRel, filepath.Base(input))
}

// ProcessorStrategy provides the per-task-type behavior plugged into the generic driver.
// T is the parsed item type (e.g., model.Command, model.Mode, or string for memory).
// The driver performs templating centrally using hooks to read/write the "content" portion.
type ProcessorStrategy[T any] interface {
	Parse(absPath string, raw []byte) (T, error)
	GetContent(item T) string
	SetContent(item T, content string) T
	FormatOne(a agent.Agent, item T) (string, error)
	FormatMany(a agent.Agent, items []T) (string, error)
}

// processGeneric is a shared driver that handles common processing flow for all task types.
// We avoid method type parameters by using a top-level generic function that accepts BaseProcessor explicitly.
func processGeneric[T any](
	p *BaseProcessor,
	inputs []string,
	cfg *OutputConfig,
	strategy ProcessorStrategy[T],
) (*TaskResult, error) {
	result := &TaskResult{Files: []ProcessedFile{}}

	items := make([]T, 0, len(inputs))
	for _, input := range inputs {
		absInputPath := filepath.Join(p.absInputRoot, input)

		// Read file
		raw, err := p.fs.ReadFile(absInputPath)
		if err != nil {
			return nil, fmt.Errorf("read input file %s: %w", absInputPath, err)
		}

		// Parse to typed item
		item, err := strategy.Parse(absInputPath, raw)
		if err != nil {
			return nil, fmt.Errorf("parse item from content %s: %w", absInputPath, err)
		}

		// Apply templating centrally using strategy-provided content accessors
		engine := p.templateEngine(cfg.AgentName)
		content := strategy.GetContent(item)
		out, err := engine.Execute(absInputPath, content, nil)
		if err != nil {
			return nil, fmt.Errorf("template execute %s: %w", input, err)
		}
		item = strategy.SetContent(item, out)

		// Directory output → format per item immediately
		if cfg.IsDirectory {
			content, err := strategy.FormatOne(cfg.Agent, item)
			if err != nil {
				return nil, fmt.Errorf("format item for agent %s: %w", cfg.AgentName, err)
			}
			result.Files = append(result.Files, ProcessedFile{
				relPath:   resolveOutputRelPath(cfg.RelPath, input),
				Content:   content,
				AgentName: cfg.AgentName,
			})
			continue
		}

		// File output → accumulate for single formatting later
		items = append(items, item)
	}

	// File output: single formatting on all items
	if !cfg.IsDirectory {
		content, err := strategy.FormatMany(cfg.Agent, items)
		if err != nil {
			return nil, fmt.Errorf("format items for agent %s: %w", cfg.AgentName, err)
		}
		result.Files = append(result.Files, ProcessedFile{
			relPath:   cfg.RelPath,
			Content:   content,
			AgentName: cfg.AgentName,
		})
	}

	return result, nil
}
