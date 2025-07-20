// Package processor provides functionality for processing agent-def tasks
package processor

import (
	"fmt"
	"path/filepath"

	"github.com/uphy/agent-def/internal/agent"
	"github.com/uphy/agent-def/internal/model"
	"github.com/uphy/agent-def/internal/parser"
	"github.com/uphy/agent-def/internal/template"
)

// CommandProcessor processes command tasks
type CommandProcessor struct {
	*BaseProcessor
}

// NewCommandProcessor creates a new CommandProcessor
func NewCommandProcessor(base *BaseProcessor) *CommandProcessor {
	return &CommandProcessor{BaseProcessor: base}
}

// Process implements the task processing for command task type
func (p *CommandProcessor) Process(inputs []string, cfg *OutputConfig) (*TaskResult, error) {
	cmds := make([]model.Command, 0, len(inputs))
	result := &TaskResult{
		Files: []ProcessedFile{},
	}

	// Process each input file
	for _, input := range inputs {
		absInputFilePath := filepath.Join(p.absInputRoot, input)

		// Read and parse the command
		inputContent, err := p.fs.ReadFile(absInputFilePath)
		if err != nil {
			return nil, fmt.Errorf("read input file %s: %w", absInputFilePath, err)
		}
		cmd, err := parser.ParseCommandFromContent(absInputFilePath, inputContent)
		if err != nil {
			return nil, fmt.Errorf("parse command from content %s: %w", absInputFilePath, err)
		}

		// Apply template processing
		adapter := NewFSAdapter(p.fs)
		engine := template.NewEngine(adapter, cfg.AgentName, p.absInputRoot, p.registry)
		out, err := engine.Execute(absInputFilePath, cmd.Content, nil)
		if err != nil {
			return nil, fmt.Errorf("template execute %s: %w", input, err)
		}
		cmd.Content = out

		// Handle directory mode output (individual files)
		if cfg.IsDirectory {
			relPath := filepath.Join(cfg.RelPath, filepath.Base(input))
			formattedCmd, err := cfg.Agent.FormatCommand([]model.Command{*cmd})
			if err != nil {
				return nil, fmt.Errorf("format command for agent %s: %w", cfg.AgentName, err)
			}
			result.Files = append(result.Files, ProcessedFile{
				relPath: relPath,
				Content: formattedCmd,
			})
		}

		// Store commands for potential concatenation
		cmds = append(cmds, *cmd)
	}

	// Handle file mode output (concatenated content)
	if !cfg.IsDirectory {
		formattedCmd, err := cfg.Agent.FormatCommand(cmds)
		if err != nil {
			return nil, fmt.Errorf("format command for agent %s: %w", cfg.AgentName, err)
		}
		result.Files = append(result.Files, ProcessedFile{
			relPath: cfg.RelPath,
			Content: formattedCmd,
		})
	}

	return result, nil
}

// GetOutputPath returns the appropriate output path for command tasks
func (p *CommandProcessor) GetOutputPath(agent agent.Agent, outputPath string) string {
	if outputPath == "" {
		return agent.CommandPath(p.userScope)
	}
	return outputPath
}
