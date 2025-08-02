package processor

import (
	"fmt"
	"path/filepath"

	"github.com/uphy/agent-sync/internal/agent"
	"github.com/uphy/agent-sync/internal/model"
	"github.com/uphy/agent-sync/internal/template"
)

// ModeProcessor processes mode tasks
type ModeProcessor struct {
	*BaseProcessor
}

// NewModeProcessor creates a new ModeProcessor
func NewModeProcessor(base *BaseProcessor) *ModeProcessor {
	return &ModeProcessor{BaseProcessor: base}
}

// Process implements the task processing for mode task type
func (p *ModeProcessor) Process(inputs []string, cfg *OutputConfig) (*TaskResult, error) {
	// Collect modes when concatenated output is required
	modes := make([]model.Mode, 0, len(inputs))
	result := &TaskResult{
		Files: []ProcessedFile{},
	}

	// Process each input file
	for _, input := range inputs {
		absInputFilePath := filepath.Join(p.absInputRoot, input)

		// Read and parse the mode
		inputContent, err := p.fs.ReadFile(absInputFilePath)
		if err != nil {
			return nil, fmt.Errorf("read input file %s: %w", absInputFilePath, err)
		}
		mode, err := model.ParseMode(absInputFilePath, inputContent)
		if err != nil {
			return nil, fmt.Errorf("parse mode from content %s: %w", absInputFilePath, err)
		}

		// Apply template processing to content only
		adapter := NewFSAdapter(p.fs)
		engine := template.NewEngine(adapter, cfg.AgentName, p.absInputRoot, p.registry)
		out, err := engine.Execute(absInputFilePath, mode.Content, nil)
		if err != nil {
			return nil, fmt.Errorf("template execute %s: %w", input, err)
		}
		mode.Content = out

		// Handle directory mode output (individual files)
		if cfg.IsDirectory {
			relPath := filepath.Join(cfg.RelPath, filepath.Base(input))
			formatted, err := cfg.Agent.FormatMode([]model.Mode{*mode})
			if err != nil {
				return nil, fmt.Errorf("format mode for agent %s: %w", cfg.AgentName, err)
			}
			result.Files = append(result.Files, ProcessedFile{
				relPath: relPath,
				Content: formatted,
			})
		} else {
			// Store for potential concatenation
			modes = append(modes, *mode)
		}
	}

	// Handle file mode output (concatenated content via agent formatter)
	if !cfg.IsDirectory {
		// Concatenated output: format all modes at once via agent formatter
		formatted, err := cfg.Agent.FormatMode(modes)
		if err != nil {
			return nil, fmt.Errorf("format mode for agent %s: %w", cfg.AgentName, err)
		}
		result.Files = append(result.Files, ProcessedFile{
			relPath: cfg.RelPath,
			Content: formatted,
		})
	}

	return result, nil
}

// GetOutputPath returns the appropriate output path for mode tasks
func (p *ModeProcessor) GetOutputPath(agent agent.Agent, outputPath string) string {
	if outputPath == "" {
		return agent.ModePath(p.userScope)
	}
	return outputPath
}
