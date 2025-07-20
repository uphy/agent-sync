// Package processor provides functionality for processing agent-def tasks
package processor

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/uphy/agent-def/internal/agent"
	"github.com/uphy/agent-def/internal/template"
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
func (p *MemoryProcessor) Process(inputs []string, cfg *OutputConfig) (*TaskResult, error) {
	// Render templates for all input files
	var renderedOutputs []string
	for _, input := range inputs {
		absInputPath := filepath.Join(p.absInputRoot, input)
		adapter := NewFSAdapter(p.fs)
		engine := template.NewEngine(adapter, cfg.AgentName, p.absInputRoot, p.registry)
		out, err := engine.ExecuteFile(absInputPath, nil)
		if err != nil {
			return nil, fmt.Errorf("template execute %s: %w", input, err)
		}
		renderedOutputs = append(renderedOutputs, out)
	}

	result := &TaskResult{
		Files: []ProcessedFile{},
	}

	// Process outputs for directory mode (individual files)
	memories := make([]string, 0, len(inputs))
	for i, input := range inputs {
		memory, err := cfg.Agent.FormatMemory(renderedOutputs[i])
		if err != nil {
			return nil, fmt.Errorf("format memory for agent %s: %w", cfg.AgentName, err)
		}

		// If directory mode, add each input as a separate file
		if cfg.IsDirectory {
			relPath := filepath.Join(cfg.RelPath, filepath.Base(input))
			result.Files = append(result.Files, ProcessedFile{
				relPath: relPath,
				Content: memory,
			})
		}

		// Save all memories for potential concatenation
		memories = append(memories, memory)
	}

	// Process output for file mode (concatenated content)
	if !cfg.IsDirectory {
		memory, err := cfg.Agent.FormatMemory(strings.Join(memories, "\n\n"))
		if err != nil {
			return nil, fmt.Errorf("format memory for agent %s: %w", cfg.AgentName, err)
		}
		result.Files = append(result.Files, ProcessedFile{
			relPath: cfg.RelPath,
			Content: memory,
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
