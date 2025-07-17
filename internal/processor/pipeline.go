package processor

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/user/agent-def/internal/agent"
	"github.com/user/agent-def/internal/config"
	"github.com/user/agent-def/internal/log"
	"github.com/user/agent-def/internal/model"
	"github.com/user/agent-def/internal/parser"
	"github.com/user/agent-def/internal/template"
	"github.com/user/agent-def/internal/util"
	"go.uber.org/zap"
)

// Pipeline represents the processing pipeline for a single task.
type Pipeline struct {
	// Task contains the configuration for the current task being processed,
	// including its name, type, inputs, targets, and concatenation settings.
	Task config.Task

	// SourceRoot is the base directory path for resolving source files and determining
	// output paths. For user-scoped tasks, it represents the user's home directory.
	SourceRoot string

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
func NewPipeline(task config.Task, sourceRoot string, destinations []string, userScope bool, dryRun bool, force bool, logger *zap.Logger, output log.OutputWriter) *Pipeline {
	// Use no-op logger if none provided
	if logger == nil {
		logger = zap.NewNop()
	}

	p := &Pipeline{
		Task:       task,
		SourceRoot: sourceRoot,
		OutputDirs: destinations,
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
	sources, err := p.resolveInputs()
	if err != nil {
		p.logger.Error("Failed to resolve inputs",
			zap.String("error", err.Error()),
			zap.Strings("inputs", p.Task.Inputs))
		if p.output != nil {
			p.output.PrintError(fmt.Errorf("failed to resolve inputs: %w", err))
		}
		return fmt.Errorf("resolve inputs: %w", err)
	}

	if len(sources) == 0 {
		p.logger.Error("No source files found", zap.String("task", t.Name))
		if p.output != nil {
			p.output.PrintError(fmt.Errorf("no source files found for task %s", t.Name))
		}
		return fmt.Errorf("no source files for task %s", t.Name)
	}

	p.logger.Debug("Inputs resolved successfully", zap.Strings("sources", sources))

	// For each output agent.
	for _, output := range t.Outputs {
		agent, ok := p.registry.Get(output.Agent)
		if !ok {
			return fmt.Errorf("unknown agent %s", output.Agent)
		}
		// Determine if inputs should be concatenated.
		concat := false
		if t.Concat != nil {
			// Based on user-defined concat setting.
			concat = *t.Concat
		} else {
			// Default concat behavior based on agent type.
			agent, ok := p.registry.Get(output.Agent)
			if ok {
				concat = agent.ShouldConcatenate(p.Task.Type)
			}
		}

		// Group files based on concat setting.
		var sourceGroups [][]string
		if concat {
			// all files are placed in a single group for combined processing
			sourceGroups = [][]string{sources}
		} else {
			// each file gets its own group and is processed separately
			for _, f := range sources {
				sourceGroups = append(sourceGroups, []string{f})
			}
		}

		// Process each group for each target.
		for _, grp := range sourceGroups {
			// Template processing.
			var rendered []string
			for _, rel := range grp {
				full := filepath.Join(p.SourceRoot, rel)
				engine := template.NewEngine(&fsAdapter{p.fs}, output.Agent, p.SourceRoot, p.registry)
				out, err := engine.ExecuteFile(full, nil)
				if err != nil {
					return fmt.Errorf("template execute %s: %w", rel, err)
				}
				rendered = append(rendered, out)
			}
			content := strings.Join(rendered, "\n\n")

			// Agent-specific formatting.
			var formatted string
			switch t.Type {
			case "memory":
				formatted, err = agent.FormatMemory(content)
			case "command":
				var cmds []model.Command
				for i, rel := range grp {
					full := filepath.Join(p.SourceRoot, rel)

					// Use the already templated content instead of re-reading the file
					templatedContent := rendered[i]

					// Parse the command from the templated content
					cmd, err2 := parser.ParseCommandFromContent(full, []byte(templatedContent))
					if err2 != nil {
						return fmt.Errorf("parse command %s: %w", rel, err2)
					}
					cmds = append(cmds, *cmd)
				}
				formatted, err = agent.FormatCommand(cmds)
			default:
				return fmt.Errorf("unsupported task type %s", t.Type)
			}
			if err != nil {
				return fmt.Errorf("format for agent %s: %w", output.Agent, err)
			}

			// Write outputs to output directories.
			for _, dest := range p.OutputDirs {
				outPath, err := p.defaultOutputPath(dest, output, grp, concat)
				if err != nil {
					return fmt.Errorf("determine output path: %w", err)
				}

				if p.DryRun {
					message := fmt.Sprintf("[DRY RUN] Would write output to %s", outPath)
					if p.fs.FileExists(outPath) {
						message += " (Will overwrite)"
					}

					p.logger.Info("Dry run - would write file",
						zap.String("path", outPath),
						zap.Bool("fileExists", p.fs.FileExists(outPath)))

					if p.output != nil {
						p.output.PrintVerbose(message)
					}
				} else {
					// Check if file exists and handle confirmation if needed
					fileExists := p.fs.FileExists(outPath)
					shouldWrite := true

					p.logger.Debug("Preparing to write file",
						zap.String("path", outPath),
						zap.Bool("fileExists", fileExists),
						zap.Bool("force", p.Force))

					if fileExists && !p.Force {
						// Prompt for confirmation
						promptMsg := fmt.Sprintf("File already exists: %s. Overwrite? [y/N]: ", outPath)
						if p.output != nil {
							p.output.Print(promptMsg)
						} else {
							fmt.Print(promptMsg) // Fallback if no output writer
						}

						reader := bufio.NewReader(os.Stdin)
						response, _ := reader.ReadString('\n')
						response = strings.TrimSpace(strings.ToLower(response))
						shouldWrite = response == "y" || response == "yes"

						p.logger.Debug("User response to overwrite prompt",
							zap.String("response", response),
							zap.Bool("shouldWrite", shouldWrite))
					}

					if shouldWrite {
						if err := p.fs.WriteFile(outPath, []byte(formatted)); err != nil {
							p.logger.Error("Failed to write output file",
								zap.String("path", outPath),
								zap.Error(err))
							if p.output != nil {
								p.output.PrintError(fmt.Errorf("failed to write output %s: %w", outPath, err))
							}
							return fmt.Errorf("write output %s: %w", outPath, err)
						}

						p.logger.Info("Successfully wrote output file",
							zap.String("path", outPath))

						if p.output != nil {
							p.output.PrintSuccess(fmt.Sprintf("Wrote output to %s", outPath))
						}
					} else {
						p.logger.Info("Skipped writing file", zap.String("path", outPath))

						if p.output != nil {
							p.output.Print(fmt.Sprintf("Skipped writing to %s", outPath))
						}
					}
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
	paths, err := p.fs.GlobWithExcludes(p.Task.Inputs, p.SourceRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to expand glob patterns: %w", err)
	}
	return paths, nil
}

// defaultOutputPath computes output path based on scope, type, and agent conventions.
func (p *Pipeline) defaultOutputPath(outputBaseDir string, output config.Output, sourceFiles []string, concat bool) (string, error) {
	// Override if output path provided.
	if output.OutputPath != "" {
		return filepath.Join(outputBaseDir, output.OutputPath), nil
	}

	name := filepath.Base(sourceFiles[len(sourceFiles)-1])
	if concat {
		name = p.Task.Name
	}

	// Use the agent's appropriate output path method based on task type
	agent, ok := p.registry.Get(output.Agent)
	if !ok {
		return "", fmt.Errorf("unknown agent %s", output.Agent)
	}

	// Call the appropriate method based on task type
	switch p.Task.Type {
	case "memory":
		return agent.DefaultMemoryPath(outputBaseDir, p.UserScope, name)
	case "command":
		return agent.DefaultCommandPath(outputBaseDir, p.UserScope, name)
	default:
		return "", fmt.Errorf("unsupported task type %s", p.Task.Type)
	}
}
