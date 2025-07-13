package cli

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

//go:embed templates templates/memories templates/commands
var templatesFS embed.FS

func NewInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize agent-def.yml and directory structure",
		Long:  "Generate a sample agent-def.yml in the current directory along with memories/ and commands/ folders with sample files.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("cannot get current directory: %w", err)
			}

			// Check if agent-def.yml already exists
			yml := filepath.Join(cwd, "agent-def.yml")
			if _, err := os.Stat(yml); err == nil {
				return fmt.Errorf("agent-def.yml already exists in %s", cwd)
			}

			// 再帰的にすべてのテンプレートファイルをコピー
			if err := copyEmbeddedDirectory(templatesFS, "templates", cwd); err != nil {
				return fmt.Errorf("failed to copy template files: %w", err)
			}

			fmt.Println("Generated agent-def.yml, memories/, and commands/ with sample files in", cwd)
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
	case templatePath == "templates/agent-def.yml":
		// templates/agent-def.yml -> agent-def.yml
		return filepath.Join(cwd, "agent-def.yml")
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
