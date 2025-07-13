package processor

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/user/agent-def/internal/agent"
	"github.com/user/agent-def/internal/config"
	"github.com/user/agent-def/internal/model"
	"github.com/user/agent-def/internal/parser"
	"github.com/user/agent-def/internal/template"
	"github.com/user/agent-def/internal/util"
)

// Pipeline represents the processing pipeline for a single task.
type Pipeline struct {
	// Task contains the configuration for the current task being processed,
	// including its name, type, sources, targets, and concatenation settings.
	Task config.Task

	// SourceRoot is the base directory path for resolving source files and determining
	// output paths. For user-scoped tasks, it represents the user's home directory.
	SourceRoot string

	// Destinations is a list of output directory paths where the processed
	// content will be written. Each destination will receive a copy of the output.
	Destinations []string

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
func NewPipeline(task config.Task, sourceRoot string, destinations []string, userScope bool, dryRun bool, force bool) *Pipeline {
	p := &Pipeline{
		Task:         task,
		SourceRoot:   sourceRoot,
		Destinations: destinations,
		UserScope:    userScope,
		DryRun:       dryRun,
		Force:        force,
		fs:           &util.RealFileSystem{},
		registry:     agent.NewRegistry(),
	}
	p.registry.Register(&agent.Roo{})
	p.registry.Register(&agent.Claude{})
	return p
}

// Execute runs the full pipeline: resolve sources, template, format, and write.
func (p *Pipeline) Execute() error {
	t := p.Task
	fmt.Printf("Executing task '%s' of type '%s'\n", t.Name, t.Type)

	// Resolve source file paths.
	sources, err := p.resolveSources()
	if err != nil {
		return fmt.Errorf("resolve sources: %w", err)
	}
	if len(sources) == 0 {
		return fmt.Errorf("no source files for task %s", t.Name)
	}

	// For each target agent.
	for _, tgt := range t.Targets {
		agent, ok := p.registry.Get(tgt.Agent)
		if !ok {
			return fmt.Errorf("unknown agent %s", tgt.Agent)
		}
		// Determine if sources should be concatenated.
		concat := false
		if t.Concat != nil {
			// Based on user-defined concat setting.
			concat = *t.Concat
		} else {
			// Default concat behavior based on agent type.
			agent, ok := p.registry.Get(tgt.Agent)
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
				engine := template.NewEngine(&fsAdapter{p.fs}, tgt.Agent, p.SourceRoot, p.registry)
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
				return fmt.Errorf("format for agent %s: %w", tgt.Agent, err)
			}

			// Write outputs to destinations.
			for _, dest := range p.Destinations {
				outPath, err := p.defaultOutputPath(dest, tgt, grp, concat)
				if err != nil {
					return fmt.Errorf("determine output path: %w", err)
				}

				if p.DryRun {
					message := fmt.Sprintf("[DRY RUN] Would write output to %s", outPath)
					if p.fs.FileExists(outPath) {
						message += " (Will overwrite)"
					}
					fmt.Println(message)
				} else {
					// Check if file exists and handle confirmation if needed
					fileExists := p.fs.FileExists(outPath)
					shouldWrite := true

					if fileExists && !p.Force {
						// Prompt for confirmation
						fmt.Printf("File already exists: %s. Overwrite? [y/N]: ", outPath)
						reader := bufio.NewReader(os.Stdin)
						response, _ := reader.ReadString('\n')
						response = strings.TrimSpace(strings.ToLower(response))
						shouldWrite = response == "y" || response == "yes"
					}

					if shouldWrite {
						if err := p.fs.WriteFile(outPath, []byte(formatted)); err != nil {
							return fmt.Errorf("write output %s: %w", outPath, err)
						}
						fmt.Printf("Wrote output to %s\n", outPath)
					} else {
						fmt.Printf("Skipped writing to %s\n", outPath)
					}
				}
			}
		}
	}

	return nil
}

// resolveSources expands Task.Sources (dirs, globs, files) to relative .md files.
func (p *Pipeline) resolveSources() ([]string, error) {
	var result []string
	for _, src := range p.Task.Sources {
		abs := filepath.Join(p.SourceRoot, src)
		info, err := os.Stat(abs)
		if err == nil && info.IsDir() {
			matches, err := p.fs.ListFiles(abs, "*.md")
			if err != nil {
				return nil, err
			}
			for _, m := range matches {
				rel, _ := filepath.Rel(p.SourceRoot, m)
				result = append(result, rel)
			}
		} else if strings.Contains(src, "*") {
			matches, err := filepath.Glob(abs)
			if err != nil {
				return nil, err
			}
			for _, m := range matches {
				rel, _ := filepath.Rel(p.SourceRoot, m)
				result = append(result, rel)
			}
		} else {
			result = append(result, src)
		}
	}
	return result, nil
}

// defaultOutputPath computes output path based on scope, type, and agent conventions.
func (p *Pipeline) defaultOutputPath(outputBaseDir string, tgt config.Target, sourceFiles []string, concat bool) (string, error) {
	// Override if target path provided.
	if tgt.Target != "" {
		return filepath.Join(outputBaseDir, tgt.Target), nil
	}

	name := filepath.Base(sourceFiles[len(sourceFiles)-1])
	if concat {
		name = p.Task.Name
	}

	// Use the agent's appropriate output path method based on task type
	agent, ok := p.registry.Get(tgt.Agent)
	if !ok {
		return "", fmt.Errorf("unknown agent %s", tgt.Agent)
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
