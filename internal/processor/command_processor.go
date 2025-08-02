// Package processor provides functionality for processing agent-sync tasks
package processor

import (
	"fmt"

	"github.com/uphy/agent-sync/internal/agent"
	"github.com/uphy/agent-sync/internal/model"
)

// CommandProcessor processes command tasks
type CommandProcessor struct {
	*BaseProcessor
}

// NewCommandProcessor creates a new CommandProcessor
func NewCommandProcessor(base *BaseProcessor) *CommandProcessor {
	return &CommandProcessor{BaseProcessor: base}
}

// commandStrategy provides parsing, content access, and formatting for commands
type commandStrategy struct {
	p         *BaseProcessor
	agentName string
}

func (s commandStrategy) Parse(absPath string, raw []byte) (model.Command, error) {
	cmd, err := model.ParseCommand(absPath, raw)
	if err != nil {
		return model.Command{}, fmt.Errorf("parse command from content %s: %w", absPath, err)
	}
	return *cmd, nil
}

func (s commandStrategy) GetContent(item model.Command) string {
	return item.Content
}

func (s commandStrategy) SetContent(item model.Command, content string) model.Command {
	item.Content = content
	return item
}

func (s commandStrategy) FormatOne(a agent.Agent, item model.Command) (string, error) {
	return a.FormatCommand([]model.Command{item})
}

func (s commandStrategy) FormatMany(a agent.Agent, items []model.Command) (string, error) {
	return a.FormatCommand(items)
}

// Process implements the task processing for command task type
func (p *CommandProcessor) Process(inputs []string, cfg *OutputConfig) (*TaskResult, error) {
	strategy := commandStrategy{p: p.BaseProcessor, agentName: cfg.AgentName}
	return processGeneric(p.BaseProcessor, inputs, cfg, strategy)
}

// GetOutputPath returns the appropriate output path for command tasks
func (p *CommandProcessor) GetOutputPath(agent agent.Agent, outputPath string) string {
	if outputPath == "" {
		return agent.CommandPath(p.userScope)
	}
	return outputPath
}
