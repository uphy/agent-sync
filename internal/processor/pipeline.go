package processor

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/uphy/agent-def/internal/agent"
	"github.com/uphy/agent-def/internal/config"
	"github.com/uphy/agent-def/internal/log"
	"github.com/uphy/agent-def/internal/model"
	"github.com/uphy/agent-def/internal/parser"
	"github.com/uphy/agent-def/internal/template"
	"github.com/uphy/agent-def/internal/util"
	"go.uber.org/zap"
)

// Pipeline represents the processing pipeline for a single task.
type Pipeline struct {
	// Task contains the configuration for the current task being processed,
	// including its name, type, inputs, targets, and concatenation settings.
	Task config.Task

	// InputRoot is the base directory path for resolving input files and determining
	// output paths. For user-scoped tasks, it represents the user's home directory.
	InputRoot string

	// OutputDirs is a list of output directory paths where the processed
	// content will be written. Each output directory will receive a copy of the output.
	OutputDirs []string

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

type outputFile struct {
	path    string
	content string
}

// fsAdapter bridges util.FileSystem to template.FileResolver.
type fsAdapter struct {
	fs util.FileSystem
}

func (a *fsAdapter) Read(path string) ([]byte, error) {
	return a.fs.ReadFile(path)
}

func (a *fsAdapter) Exists(path string) bool {
	return a.fs.FileExists(path)
}

func (a *fsAdapter) ResolvePath(path string) string {
	return a.fs.ResolvePath(path)
}

// NewPipeline creates a new Pipeline with context and registers built-in agents.
func NewPipeline(task config.Task, inputRoot string, outputDirs []string, userScope bool, dryRun bool, force bool, logger *zap.Logger, output log.OutputWriter) *Pipeline {
	// Use no-op logger if none provided
	if logger == nil {
		logger = zap.NewNop()
	}

	p := &Pipeline{
		Task:       task,
		InputRoot:  inputRoot,
		OutputDirs: outputDirs,
		UserScope:  userScope,
		DryRun:     dryRun,
		Force:      force,
		fs:         &util.RealFileSystem{},
		registry:   agent.NewRegistry(),
		logger:     logger,
		output:     output,
	}
	p.registry.Register(&agent.Roo{})
	p.registry.Register(&agent.Claude{})
	p.registry.Register(&agent.Cline{})
	return p
}

// Execute runs the full pipeline: resolve inputs, template, format, and write.
func (p *Pipeline) Execute() error {
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

	// Resolve input file paths.
	inputs, err := p.resolveInputs()
	if err != nil {
		p.logger.Error("Failed to resolve inputs",
			zap.String("error", err.Error()),
			zap.Strings("inputs", p.Task.Inputs))
		if p.output != nil {
			p.output.PrintError(fmt.Errorf("failed to resolve inputs: %w", err))
		}
		return fmt.Errorf("resolve inputs: %w", err)
	}

	if len(inputs) == 0 {
		p.logger.Error("No source files found", zap.String("task", t.Name))
		if p.output != nil {
			p.output.PrintError(fmt.Errorf("no source files found for task %s", t.Name))
		}
		return fmt.Errorf("no source files for task %s", t.Name)
	}

	p.logger.Debug("Inputs resolved successfully", zap.Strings("inputs", inputs))

	// For each output agent.
	for _, output := range t.Outputs {
		agent, ok := p.registry.Get(output.Agent)
		if !ok {
			return fmt.Errorf("unknown agent %s", output.Agent)
		}
		outputPath := ""
		switch t.Type {
		case "memory":
			outputPath = agent.MemoryPath(p.UserScope)
		case "command":
			outputPath = agent.CommandPath(p.UserScope)
		default:
			return fmt.Errorf("unsupported task type %s", t.Type)
		}
		if output.OutputPath != "" {
			outputPath = output.OutputPath
		}

		// Determine if we should treat this as a directory based on output path
		isDirectory := strings.HasSuffix(outputPath, "/")
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

		outputFiles := make([]outputFile, 0)
		switch t.Type {
		case "memory":
			var renderedOutputs []string
			for _, input := range inputs {
				inputPath := filepath.Join(p.InputRoot, input)
				engine := template.NewEngine(&fsAdapter{p.fs}, output.Agent, p.InputRoot, p.registry)
				out, err := engine.ExecuteFile(inputPath, nil)
				if err != nil {
					return fmt.Errorf("template execute %s: %w", input, err)
				}
				renderedOutputs = append(renderedOutputs, out)
			}
			memories := make([]string, 0, len(inputs))
			for i, input := range inputs {
				memory, err := agent.FormatMemory(renderedOutputs[i])
				if err != nil {
					return fmt.Errorf("format memory for agent %s: %w", output.Agent, err)
				}
				if isDirectory {
					out := filepath.Join(outputPath, filepath.Base(input))
					outputFiles = append(outputFiles, outputFile{path: out, content: memory})
				}
				memories = append(memories, memory)
			}
			if !isDirectory {
				memory, err := agent.FormatMemory(strings.Join(memories, "\n\n"))
				if err != nil {
					return fmt.Errorf("format memory for agent %s: %w", output.Agent, err)
				}
				outputFiles = append(outputFiles, outputFile{path: outputPath, content: memory})
			}
		case "command":
			cmds := make([]model.Command, 0, len(inputs))
			for _, input := range inputs {
				inputPath := filepath.Join(p.InputRoot, input)
				inputContent, err := p.fs.ReadFile(inputPath)
				if err != nil {
					return fmt.Errorf("read input file %s: %w", inputPath, err)
				}
				cmd, err := parser.ParseCommandFromContent(inputPath, inputContent)
				if err != nil {
					return fmt.Errorf("parse command from content %s: %w", inputPath, err)
				}
				engine := template.NewEngine(&fsAdapter{p.fs}, output.Agent, p.InputRoot, p.registry)
				out, err := engine.Execute(cmd.Content, nil)
				if err != nil {
					return fmt.Errorf("template execute %s: %w", input, err)
				}
				cmd.Content = out
				if isDirectory {
					out := filepath.Join(outputPath, filepath.Base(input))
					formattedCmd, err := agent.FormatCommand([]model.Command{*cmd})
					if err != nil {
						return fmt.Errorf("format command for agent %s: %w", output.Agent, err)
					}
					outputFiles = append(outputFiles, outputFile{path: out, content: formattedCmd})
				}
				cmds = append(cmds, *cmd)
			}
			if !isDirectory {
				formattedCmd, err := agent.FormatCommand(cmds)
				if err != nil {
					return fmt.Errorf("format command for agent %s: %w", output.Agent, err)
				}
				outputFiles = append(outputFiles, outputFile{path: outputPath, content: formattedCmd})
			}
		default:
			return fmt.Errorf("unsupported task type %s", t.Type)
		}

		for _, outputDir := range p.OutputDirs {
			for _, outFile := range outputFiles {
				out := filepath.Join(outputDir, outFile.path)
				content := outFile.content
				if p.DryRun {
					p.logger.Info("[DRY RUN] Would write file", zap.String("path", out))
					if p.output != nil {
						p.output.PrintVerbose(fmt.Sprintf("[DRY RUN] Would write file %s", out))
					}
				} else {
					if err := p.fs.WriteFile(out, []byte(content)); err != nil {
						return fmt.Errorf("write file %s: %w", out, err)
					}
					p.logger.Info("Wrote file", zap.String("path", out))
				}
			}
		}
	}

	p.logger.Info("Pipeline execution completed successfully",
		zap.String("task", t.Name),
		zap.String("type", string(t.Type)))

	return nil
}

// resolveInputs expands Task.Inputs with support for glob patterns and exclusions
func (p *Pipeline) resolveInputs() ([]string, error) {
	// Use the new GlobWithExcludes function from the filesystem interface
	paths, err := p.fs.GlobWithExcludes(p.Task.Inputs, p.InputRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to expand glob patterns: %w", err)
	}
	return paths, nil
}
