// Package processor provides functionality for processing agent-sync tasks
package processor

import (
	"strings"

	"github.com/uphy/agent-sync/internal/agent"
)

// MemoryProcessor processes memory tasks
type MemoryProcessor struct {
	*BaseProcessor
}

// NewMemoryProcessor creates a new MemoryProcessor
func NewMemoryProcessor(base *BaseProcessor) *MemoryProcessor {
	return &MemoryProcessor{BaseProcessor: base}
}

// memoryStrategy provides parsing, content access, and formatting for memory (string-based)
type memoryStrategy struct {
	p         *BaseProcessor
	agentName string
}

func (s memoryStrategy) Parse(absPath string, raw []byte) (string, error) {
	return string(raw), nil
}

func (s memoryStrategy) GetContent(item string) string {
	return item
}

func (s memoryStrategy) SetContent(item string, content string) string {
	return content
}

func (s memoryStrategy) FormatOne(a agent.Agent, item string) (string, error) {
	return a.FormatMemory(item)
}

func (s memoryStrategy) FormatMany(a agent.Agent, items []string) (string, error) {
	joined := strings.Join(items, "\n\n")
	return a.FormatMemory(joined)
}

// Process implements the task processing for memory task type
func (p *MemoryProcessor) Process(inputs []string, cfg *OutputConfig) (*TaskResult, error) {
	strategy := memoryStrategy{p: p.BaseProcessor, agentName: cfg.AgentName}
	return processGeneric(p.BaseProcessor, inputs, cfg, strategy)
}

// GetOutputPath returns the appropriate output path for memory tasks
func (p *MemoryProcessor) GetOutputPath(agent agent.Agent, outputPath string) string {
	if outputPath == "" {
		return agent.MemoryPath(p.userScope)
	}
	return outputPath
}
