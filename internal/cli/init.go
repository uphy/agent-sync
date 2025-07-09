package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func NewInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize agent-def.yml and directory structure",
		Long:  "Generate a sample agent-def.yml in the current directory along with context/ and commands/ folders.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("cannot get current directory: %w", err)
			}
			yml := filepath.Join(cwd, "agent-def.yml")
			if _, err := os.Stat(yml); err == nil {
				return fmt.Errorf("agent-def.yml already exists in %s", cwd)
			}
			template := `configVersion: "1.0"

projects:
  my-project:
    root: ./
    destinations: []
    tasks: []

user:
  tasks: []
`
			if err := os.WriteFile(yml, []byte(template), 0644); err != nil {
				return fmt.Errorf("failed to write agent-def.yml: %w", err)
			}
			if err := os.MkdirAll(filepath.Join(cwd, "context"), 0755); err != nil {
				return fmt.Errorf("failed to create context directory: %w", err)
			}
			if err := os.MkdirAll(filepath.Join(cwd, "commands"), 0755); err != nil {
				return fmt.Errorf("failed to create commands directory: %w", err)
			}
			fmt.Println("Generated agent-def.yml, context/, and commands/ in", cwd)
			return nil
		},
	}
	return cmd
}
