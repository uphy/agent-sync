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

		// Write the processed files
		if err := p.writeOutputFiles(result.Files); err != nil {
			return err
		}
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

// writeOutputFiles writes the processed files to all output directories
func (p *Pipeline) writeOutputFiles(files []ProcessedFile) error {
	for _, absOutputDir := range p.AbsOutputDirs {
		for _, file := range files {
			absOutputFile := filepath.Join(absOutputDir, file.relPath)

			if isSub, err := util.IsSub(absOutputDir, absOutputFile); err != nil {
				return fmt.Errorf("check if output path is subdirectory: %w", err)
			} else if !isSub {
				return fmt.Errorf("output path %s is not a subdirectory of %s", absOutputFile, absOutputDir)
			}

			if p.DryRun {
				p.logger.Info("[DRY RUN] Would write file", zap.String("path", absOutputFile))
				if p.output != nil {
					p.output.PrintVerbose(fmt.Sprintf("[DRY RUN] Would write file %s", absOutputFile))
				}
			} else {
				if err := p.fs.WriteFile(absOutputFile, []byte(file.Content)); err != nil {
					return fmt.Errorf("write file %s: %w", absOutputFile, err)
				}
				p.logger.Info("Wrote file", zap.String("path", absOutputFile))
			}
		}
	}
	return nil
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
