// Package processor provides functionality for processing agent-sync tasks
package processor

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/uphy/agent-sync/internal/agent"
	"github.com/uphy/agent-sync/internal/config"
	"github.com/uphy/agent-sync/internal/log"
	"github.com/uphy/agent-sync/internal/util"
	"go.uber.org/zap"
)

// Pipeline represents the processing pipeline for a single task.
type Pipeline struct {
	// Task contains the configuration for the current task being processed,
	// including its name, type, inputs, targets, and concatenation settings.
	Task config.Task

	// AbsInputRoot is the base directory path for resolving input files and determining
	// output paths. For user-scoped tasks, it represents the user's home directory.
	AbsInputRoot string

	// AbsOutputDirs is a list of output directory paths where the processed
	// content will be written. Each output directory will receive a copy of the output.
	AbsOutputDirs []string

	// UserScope indicates whether the task operates in user scope (true) or project
	// scope (false). This affects output path resolution and file organization.
	UserScope bool

	// DryRun indicates whether the pipeline should actually write files to disk.
	// When true, the pipeline will simulate file operations but not actually write files.
	DryRun bool

	// Force indicates whether to overwrite files without confirmation.
	// When true, existing files will be overwritten without prompting.
	Force bool

	// fs is the file system interface used for all file operations,
	// such as reading source files and writing output files.
	fs util.FileSystem

	// registry contains all available agent implementations that can process
	// and format content for different target systems.
	registry *agent.Registry

	// logger is used for internal logging of pipeline operations.
	logger *zap.Logger

	// output is used for user-facing output of pipeline status and results.
	output log.OutputWriter
}

// NewPipeline creates a new Pipeline with context and registers built-in agents.
func NewPipeline(task config.Task, absInputRoot string, outputDirs []string, userScope bool, dryRun bool, force bool, logger *zap.Logger, output log.OutputWriter) (*Pipeline, error) {
	// Use no-op logger if none provided
	if logger == nil {
		logger = zap.NewNop()
	}

	if !filepath.IsAbs(absInputRoot) {
		return nil, fmt.Errorf("input root must be absolute: %s", absInputRoot)
	}
	if len(outputDirs) == 0 {
		return nil, fmt.Errorf("at least one output directory must be specified")
	}
	for _, dir := range outputDirs {
		if !filepath.IsAbs(dir) {
			return nil, fmt.Errorf("output directory must be absolute: %s", dir)
		}
	}

	p := &Pipeline{
		Task:          task,
		AbsInputRoot:  absInputRoot,
		AbsOutputDirs: outputDirs,
		UserScope:     userScope,
		DryRun:        dryRun,
		Force:         force,
		fs:            &util.RealFileSystem{},
		registry:      agent.NewRegistry(),
		logger:        logger,
		output:        output,
	}
	p.registry.Register(&agent.Roo{})
	p.registry.Register(&agent.Claude{})
	p.registry.Register(&agent.Cline{})
	return p, nil
}

// newTaskProcessor creates a TaskProcessor based on the task type
func (p *Pipeline) newTaskProcessor(taskType string) (TaskProcessor, error) {
	base := NewBaseProcessor(p.fs, p.logger, p.AbsInputRoot, p.registry, p.UserScope)

	switch taskType {
	case "memory":
		return NewMemoryProcessor(base), nil
	case "command":
		return NewCommandProcessor(base), nil
	case "mode":
		return NewModeProcessor(base), nil
	default:
		return nil, fmt.Errorf("unsupported task type %s", taskType)
	}
}

// Execute runs the full pipeline: resolve inputs, process based on task type, and write outputs.
// This implementation uses polymorphism through TaskProcessor interface instead of type switching.
func (p *Pipeline) Execute() error {
	// Log start of task execution
	p.logTaskStart()

	// Resolve and validate inputs
	inputs, err := p.resolveAndValidateInputs()
	if err != nil {
		return err // Error already logged in resolveAndValidateInputs
	}

	// Create the appropriate task processor based on task type
	processor, err := p.newTaskProcessor(p.Task.Type)
	if err != nil {
		return err
	}

	// Map to store files by agent
	filesByAgent := make(map[string][]ProcessedFile)

	// Process each output agent
	for _, output := range p.Task.Outputs {
		// Get output configuration
		cfg, err := p.getOutputConfig(output)
		if err != nil {
			return err
		}

		// Process the task using the appropriate processor
		result, err := processor.Process(inputs, cfg)
		if err != nil {
			return err
		}

		// Store files by agent
		for _, file := range result.Files {
			// Add agent name to each file
			file.AgentName = output.Agent
			filesByAgent[output.Agent] = append(filesByAgent[output.Agent], file)
		}
	}

	// Write all processed files organized by agent
	if err := p.writeOutputFilesByAgent(filesByAgent); err != nil {
		return err
	}

	// Log successful completion
	p.logTaskCompletion()

	return nil
}

// resolveAndValidateInputs expands Task.Inputs with support for glob patterns and exclusions,
// then validates that at least one input file was found.
func (p *Pipeline) resolveAndValidateInputs() ([]string, error) {
	// Expand globs and apply exclusions
	paths, err := p.fs.GlobWithExcludes(p.Task.Inputs, p.AbsInputRoot)
	if err != nil {
		p.logError("Failed to resolve inputs", err,
			zap.Strings("inputs", p.Task.Inputs))
		return nil, fmt.Errorf("resolve inputs: %w", err)
	}

	// Validate that we have at least one input file
	if len(paths) == 0 {
		err := fmt.Errorf("no source files for task %s", p.Task.Name)
		p.logError("No source files found", err,
			zap.String("task", p.Task.Name))
		return nil, err
	}

	p.logger.Debug("Inputs resolved successfully", zap.Strings("inputs", paths))
	return paths, nil
}

// getOutputConfig resolves the output path for an agent based on task type and user configuration
func (p *Pipeline) getOutputConfig(output config.Output) (*OutputConfig, error) {
	// Look up the agent implementation
	agent, ok := p.registry.Get(output.Agent)
	if !ok {
		return nil, fmt.Errorf("unknown agent %s", output.Agent)
	}

	// Create processor based on task type
	processor, err := p.newTaskProcessor(p.Task.Type)
	if err != nil {
		return nil, err
	}

	// Get output path using processor
	relOutputPath := processor.GetOutputPath(agent, output.OutputPath)

	// Determine if we should treat this as a directory based on output path
	isDirectory := strings.HasSuffix(relOutputPath, "/")

	// Log the output mode
	if isDirectory {
		p.logger.Debug("Using directory mode for output (path ends with '/')",
			zap.String("agent", output.Agent),
			zap.String("outputPath", output.OutputPath),
			zap.Bool("isDirectory", isDirectory))
	} else {
		p.logger.Debug("Using file mode for output (path does not end with '/')",
			zap.String("agent", output.Agent),
			zap.String("outputPath", output.OutputPath),
			zap.Bool("isDirectory", isDirectory))
	}

	return &OutputConfig{
		Agent:       agent,
		RelPath:     relOutputPath,
		IsDirectory: isDirectory,
		AgentName:   output.Agent,
	}, nil
}

// logError logs an error to both the zap logger and output writer
func (p *Pipeline) logError(msg string, err error, fields ...zap.Field) {
	// Log to zap logger
	p.logger.Error(msg, append(fields, zap.Error(err))...)

	// Log to output writer if available
	if p.output != nil {
		p.output.PrintError(fmt.Errorf("%s: %w", msg, err))
	}
}

// writeOutputFilesByAgent writes the processed files to all output directories, grouped by agent
func (p *Pipeline) writeOutputFilesByAgent(filesByAgent map[string][]ProcessedFile) error {
	if p.DryRun && p.output != nil {
		// Process and print files grouped by agent
		for agentName, files := range filesByAgent {
			// Count for per-agent summary
			createCount := 0
			modifyCount := 0
			unchangedCount := 0

			// Print agent header
			p.output.Print(fmt.Sprintf("\nAgent: %s", agentName))

			// Track status messages for this agent
			var statusMessages []string

			// Process all files for this agent across all output directories
			for _, absOutputDir := range p.AbsOutputDirs {
				for _, file := range files {
					absOutputFile := filepath.Join(absOutputDir, file.relPath)
					contentLength := len(file.Content)

					if isSub, err := util.IsSub(absOutputDir, absOutputFile); err != nil {
						return fmt.Errorf("check if output path is subdirectory: %w", err)
					} else if !isSub {
						return fmt.Errorf("output path %s is not a subdirectory of %s", absOutputFile, absOutputDir)
					}

					fileExists := p.fs.FileExists(absOutputFile)
					unchanged := false

					// Check if file content would change
					if fileExists {
						existingContent, err := p.fs.ReadFile(absOutputFile)
						if err == nil {
							// Compare content
							if string(existingContent) == file.Content {
								unchanged = true
								unchangedCount++
							} else {
								modifyCount++
							}
						} else {
							// If we can't read the file, assume it will be modified
							p.logger.Debug("Could not read existing file for comparison",
								zap.String("path", absOutputFile),
								zap.Error(err))
							modifyCount++
						}
					} else {
						createCount++
					}

					statusMsg := p.formatDryRunFileStatus(absOutputFile, contentLength, unchanged)
					p.logger.Info("[DRY RUN] " + statusMsg)

					// Collect status messages without the redundant [DRY RUN] prefix
					statusMessages = append(statusMessages, fmt.Sprintf("  %s", statusMsg))
				}
			}

			// Print all collected file statuses for this agent
			for _, msg := range statusMessages {
				p.output.Print(msg)
			}

			// Print per-agent summary
			p.output.Print(fmt.Sprintf("\n  Summary: %d files would be created, %d files would be modified, %d files would remain unchanged",
				createCount, modifyCount, unchangedCount))
		}
	} else {
		// Non-dry-run mode: actually write the files
		for _, files := range filesByAgent {
			for _, absOutputDir := range p.AbsOutputDirs {
				for _, file := range files {
					absOutputFile := filepath.Join(absOutputDir, file.relPath)
					contentLength := len(file.Content)

					if isSub, err := util.IsSub(absOutputDir, absOutputFile); err != nil {
						return fmt.Errorf("check if output path is subdirectory: %w", err)
					} else if !isSub {
						return fmt.Errorf("output path %s is not a subdirectory of %s", absOutputFile, absOutputDir)
					}

					if err := p.fs.WriteFile(absOutputFile, []byte(file.Content)); err != nil {
						return fmt.Errorf("write file %s: %w", absOutputFile, err)
					}
					p.logger.Info("Wrote file", zap.String("path", absOutputFile), zap.Int("bytes", contentLength))
				}
			}
		}
	}

	return nil
}

// writeOutputFiles writes the processed files to all output directories
// Kept for backward compatibility
func (p *Pipeline) writeOutputFiles(files []ProcessedFile) error {
	// Group files by a dummy agent name to use the new method
	filesByAgent := map[string][]ProcessedFile{"": files}
	return p.writeOutputFilesByAgent(filesByAgent)
}

// formatDryRunFileStatus returns a formatted status string for dry run output
func (p *Pipeline) formatDryRunFileStatus(absOutputFile string, contentLength int, unchanged bool) string {
	fileExists := p.fs.FileExists(absOutputFile)
	status := "CREATE"
	if fileExists {
		if unchanged {
			status = "UNCHANGED"
		} else {
			status = "MODIFY"
		}
	}

	// Format size in a human-readable way
	var sizeStr string
	if contentLength < 1024 {
		sizeStr = fmt.Sprintf("%d B", contentLength)
	} else if contentLength < 1024*1024 {
		sizeStr = fmt.Sprintf("%.1f KB", float64(contentLength)/1024)
	} else {
		sizeStr = fmt.Sprintf("%.1f MB", float64(contentLength)/(1024*1024))
	}

	return fmt.Sprintf("[%s] %s (%s)", status, absOutputFile, sizeStr)
}

// logTaskStart logs the start of task execution
func (p *Pipeline) logTaskStart() {
	t := p.Task

	// User-facing output
	if p.output != nil {
		p.output.PrintProgress(fmt.Sprintf("Executing task '%s' of type '%s'", t.Name, t.Type))
	}

	// Internal logging
	p.logger.Info("Pipeline execution started",
		zap.String("task", t.Name),
		zap.String("type", string(t.Type)),
		zap.Bool("dry-run", p.DryRun),
		zap.Bool("force", p.Force),
		zap.Bool("userScope", p.UserScope))
}

// logTaskCompletion logs the successful completion of a task
func (p *Pipeline) logTaskCompletion() {
	t := p.Task

	p.logger.Info("Pipeline execution completed successfully",
		zap.String("task", t.Name),
		zap.String("type", string(t.Type)))
}
