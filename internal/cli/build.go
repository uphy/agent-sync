package cli

import (
	"github.com/spf13/cobra"
	"github.com/user/agent-def/internal/processor"
)

// NewBuildCommand returns the 'build' command for the CLI.
func NewBuildCommand() *cobra.Command {
	var userOnly bool
	var dryRun bool
	var force bool
	var configPath string

	cmd := &cobra.Command{
		Use:   "build [project...]",
		Short: "Generate files based on agent-def.yml",
		Long: `Generate files for one or more projects or user-level tasks as defined in agent-def.yml.
If no project names are provided, all projects will be processed.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr, err := processor.NewManager(configPath)
			if err != nil {
				return err
			}
			return mgr.Build(args, userOnly, dryRun, force)
		},
	}

	cmd.Flags().BoolVar(&userOnly, "user", false, "Build only user-level tasks")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be generated without writing files")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite without prompting for confirmation")
	cmd.Flags().StringVarP(&configPath, "config", "c", ".", "Path to agent-def.yml file or directory containing it")

	return cmd
}
