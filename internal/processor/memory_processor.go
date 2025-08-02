// Package processor provides functionality for processing agent-sync tasks
package processor

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/uphy/agent-sync/internal/agent"
	"github.com/uphy/agent-sync/internal/template"
)

// MemoryProcessor processes memory tasks
type MemoryProcessor struct {
	*BaseProcessor
}

// NewMemoryProcessor creates a new MemoryProcessor
func NewMemoryProcessor(base *BaseProcessor) *MemoryProcessor {
	return &MemoryProcessor{BaseProcessor: base}
}

// Process implements the task processing for memory task type
// Preserve external behavior (output shape, concatenation rules, default paths)
// while aligning the internal flow with command_processor.go / mode_processor.go.
func (p *MemoryProcessor) Process(inputs []string, cfg *OutputConfig) (*TaskResult, error) {
	concatParts := make([]string, 0, len(inputs))
	result := &TaskResult{Files: []ProcessedFile{}}

	for _, input := range inputs {
		absInputPath := filepath.Join(p.absInputRoot, input)

		// Read and parse the command
		inputContent, err := p.fs.ReadFile(absInputPath)
		if err != nil {
			return nil, fmt.Errorf("read input file %s: %w", absInputPath, err)
		}

		adapter := NewFSAdapter(p.fs)
		engine := template.NewEngine(adapter, cfg.AgentName, p.absInputRoot, p.registry)
		out, err := engine.Execute(absInputPath, string(inputContent), nil)
		if err != nil {
			return nil, fmt.Errorf("template execute %s: %w", input, err)
		}

		if cfg.IsDirectory {
			// Directory output: format and emit per file
			formatted, err := cfg.Agent.FormatMemory(out)
			if err != nil {
				return nil, fmt.Errorf("format memory for agent %s: %w", cfg.AgentName, err)
			}
			relPath := filepath.Join(cfg.RelPath, filepath.Base(input))
			result.Files = append(result.Files, ProcessedFile{
				relPath: relPath,
				Content: formatted,
			})
		} else {
			// File output: accumulate to concatenate later
			concatParts = append(concatParts, out)
		}
	}

	// If not directory, perform single formatting on concatenated content like command_processor
	if !cfg.IsDirectory {
		joined := strings.Join(concatParts, "\n\n")
		formatted, err := cfg.Agent.FormatMemory(joined)
		if err != nil {
			return nil, fmt.Errorf("format memory for agent %s: %w", cfg.AgentName, err)
		}
		result.Files = append(result.Files, ProcessedFile{
			relPath: cfg.RelPath,
			Content: formatted,
		})
	}

	return result, nil
}

// GetOutputPath returns the appropriate output path for memory tasks
func (p *MemoryProcessor) GetOutputPath(agent agent.Agent, outputPath string) string {
	if outputPath == "" {
		return agent.MemoryPath(p.userScope)
	}
	return outputPath
}
