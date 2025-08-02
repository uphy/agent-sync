package processor

import (
	"fmt"

	"github.com/uphy/agent-sync/internal/agent"
	"github.com/uphy/agent-sync/internal/model"
)

// ModeProcessor processes mode tasks
type ModeProcessor struct {
	*BaseProcessor
}

// NewModeProcessor creates a new ModeProcessor
func NewModeProcessor(base *BaseProcessor) *ModeProcessor {
	return &ModeProcessor{BaseProcessor: base}
}

// modeStrategy provides parsing, content access, and formatting for modes
type modeStrategy struct {
	p         *BaseProcessor
	agentName string
}

func (s modeStrategy) Parse(absPath string, raw []byte) (model.Mode, error) {
	m, err := model.ParseMode(absPath, raw)
	if err != nil {
		return model.Mode{}, fmt.Errorf("parse mode from content %s: %w", absPath, err)
	}
	return *m, nil
}

func (s modeStrategy) GetContent(item model.Mode) string {
	return item.Content
}

func (s modeStrategy) SetContent(item model.Mode, content string) model.Mode {
	item.Content = content
	return item
}

func (s modeStrategy) FormatOne(a agent.Agent, item model.Mode) (string, error) {
	return a.FormatMode([]model.Mode{item})
}

func (s modeStrategy) FormatMany(a agent.Agent, items []model.Mode) (string, error) {
	return a.FormatMode(items)
}

// Process implements the task processing for mode task type
func (p *ModeProcessor) Process(inputs []string, cfg *OutputConfig) (*TaskResult, error) {
	strategy := modeStrategy{p: p.BaseProcessor, agentName: cfg.AgentName}
	return processGeneric(p.BaseProcessor, inputs, cfg, strategy)
}

// GetOutputPath returns the appropriate output path for mode tasks
func (p *ModeProcessor) GetOutputPath(agent agent.Agent, outputPath string) string {
	if outputPath == "" {
		return agent.ModePath(p.userScope)
	}
	return outputPath
}
