package cli

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/uphy/agent-sync/internal/log"
	"go.uber.org/zap"
)

//go:embed templates templates/memories templates/commands
var templatesFS embed.FS

// initContext holds the CLI context for the init command
var initContext *Context

// NewInitCommand returns the 'init' command with a nil context.
// It initializes the default agent-sync project structure.
func NewInitCommand() *cobra.Command {
	return NewInitCommandWithContext(nil)
}

// NewInitCommandWithContext returns the 'init' command with the specified context.
func NewInitCommandWithContext(ctx *Context) *cobra.Command {
	// Store context for command execution
	initContext = ctx

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize agent-sync.yml and directory structure",
		Long:  "Generate a sample agent-sync.yml in the current directory along with memories/ and commands/ folders with sample files.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Use the context if available
			var logger *zap.Logger
			var output log.OutputWriter

			if initContext != nil {
				logger = initContext.Logger
				output = initContext.Output

				// Log command execution
				logger.Info("Executing init command")
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
			if _, err := os.Stat(yml); err == nil {
				if logger != nil {
					logger.Warn("agent-sync.yml already exists", zap.String("path", yml))
				}
				return fmt.Errorf("agent-sync.yml already exists in %s", cwd)
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
	return cmd
}

// 再帰的にディレクトリからすべてのファイルを取得する
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

// テンプレートパスから実際のファイルシステム上のパスへマッピングする
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
		// 他のファイルの場合はtemplatesプレフィックスを削除
		rel := strings.TrimPrefix(templatePath, "templates/")
		if rel == templatePath {
			// プレフィックスが削除されていない場合は、そのまま使用
			return filepath.Join(cwd, templatePath)
		}
		return filepath.Join(cwd, rel)
	}
}

// 再帰的にディレクトリ内のすべてのファイルをコピーする
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

		// ディレクトリを作成
		destDir := filepath.Dir(destPath)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			copyErrors = append(copyErrors, fmt.Sprintf("failed to create directory %s: %v", destDir, err))
			continue
		}

		// ファイルをコピー
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
