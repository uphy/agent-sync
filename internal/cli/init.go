package cli

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/uphy/agent-sync/internal/log"
	"github.com/urfave/cli/v3"
	"go.uber.org/zap"
)

//go:embed templates templates/memories templates/commands
var templatesFS embed.FS

// NewInitCommand returns the 'init' command for urfave/cli.
// It initializes the default agent-sync project structure.
func NewInitCommand() *cli.Command {
	return &cli.Command{
		Name:        "init",
		Usage:       "Initialize agent-sync.yml and directory structure",
		Description: "Generate a sample agent-sync.yml in the current directory along with memories/ and commands/ folders with sample files.",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Force overwrite of existing files",
				Sources: cli.EnvVars("AGENT_SYNC_INIT_FORCE"),
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// Get the shared context from app metadata
			sharedContext := GetSharedContext(cmd)

			// Get command-specific flags
			force := cmd.Bool("force")

			var logger *zap.Logger
			var output log.OutputWriter

			// Get logger and output from shared context
			if sharedContext != nil {
				logger = sharedContext.Logger
				output = sharedContext.Output

				// Log command execution
				logger.Info("Executing init command",
					zap.Bool("force", force))
			}

			cwd, err := os.Getwd()
			if err != nil {
				if logger != nil {
					logger.Error("Failed to get current directory", zap.Error(err))
				}
				return fmt.Errorf("cannot get current directory: %w", err)
			}

			// Check if agent-sync.yml already exists
			yml := filepath.Join(cwd, "agent-sync.yml")
			if _, err := os.Stat(yml); err == nil && !force {
				if logger != nil {
					logger.Warn("agent-sync.yml already exists", zap.String("path", yml))
				}
				return fmt.Errorf("agent-sync.yml already exists in %s (use --force to overwrite)", cwd)
			}

			if logger != nil {
				logger.Debug("Copying template files", zap.String("destination", cwd))
			}

			// Copy all template files recursively
			if err := copyEmbeddedDirectory(templatesFS, "templates", cwd); err != nil {
				if logger != nil {
					logger.Error("Failed to copy template files", zap.Error(err))
				}
				return fmt.Errorf("failed to copy template files: %w", err)
			}

			if output != nil {
				output.PrintSuccess(fmt.Sprintf("Generated agent-sync.yml, memories/, and commands/ with sample files in %s", cwd))
			} else {
				fmt.Println("Generated agent-sync.yml, memories/, and commands/ with sample files in", cwd)
			}

			if logger != nil {
				logger.Info("Successfully initialized project",
					zap.String("directory", cwd),
					zap.String("config", yml))
			}

			return nil
		},
	}
}

// getAllFiles gets all files from a directory recursively
func getAllFiles(fsys fs.FS, dir string) ([]string, error) {
	var files []string

	err := fs.WalkDir(fsys, dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	return files, nil
}

// mapTemplatePath maps template path to filesystem path
func mapTemplatePath(templatePath string, cwd string) string {
	switch {
	case strings.HasPrefix(templatePath, "templates/memories/"):
		// templates/memories/* -> memories/*
		rel := strings.TrimPrefix(templatePath, "templates/memories/")
		return filepath.Join(cwd, "memories", rel)
	case strings.HasPrefix(templatePath, "templates/commands/"):
		// templates/commands/* -> commands/*
		rel := strings.TrimPrefix(templatePath, "templates/commands/")
		return filepath.Join(cwd, "commands", rel)
	case templatePath == "templates/agent-sync.yml":
		// templates/agent-sync.yml -> agent-sync.yml
		return filepath.Join(cwd, "agent-sync.yml")
	default:
		// Other files - remove templates prefix
		rel := strings.TrimPrefix(templatePath, "templates/")
		if rel == templatePath {
			// Prefix was not removed - use as is
			return filepath.Join(cwd, templatePath)
		}
		return filepath.Join(cwd, rel)
	}
}

// copyEmbeddedDirectory copies all files from an embedded directory
func copyEmbeddedDirectory(fsys fs.FS, dir string, cwd string) error {
	files, err := getAllFiles(fsys, dir)
	if err != nil {
		return fmt.Errorf("failed to get template files: %w", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no template files found in %s", dir)
	}

	var copyErrors []string
	for _, file := range files {
		destPath := mapTemplatePath(file, cwd)

		// Create directory
		destDir := filepath.Dir(destPath)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			copyErrors = append(copyErrors, fmt.Sprintf("failed to create directory %s: %v", destDir, err))
			continue
		}

		// Copy file
		data, err := fs.ReadFile(fsys, file)
		if err != nil {
			copyErrors = append(copyErrors, fmt.Sprintf("failed to read embedded file %s: %v", file, err))
			continue
		}

		if err := os.WriteFile(destPath, data, 0644); err != nil {
			copyErrors = append(copyErrors, fmt.Sprintf("failed to write file %s: %v", destPath, err))
			continue
		}

		fmt.Printf("Created: %s\n", destPath)
	}

	if len(copyErrors) > 0 {
		return fmt.Errorf("encountered %d errors while copying files:\n%s",
			len(copyErrors), strings.Join(copyErrors, "\n"))
	}

	return nil
}
