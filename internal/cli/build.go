package cli

import (
	"github.com/spf13/cobra"
	"github.com/user/agent-def/internal/processor"
)

// NewBuildCommand returns the 'build' command for the CLI.
func NewBuildCommand() *cobra.Command {
	var userOnly bool
	var watch bool
	var dryRun bool
	var force bool

	cmd := &cobra.Command{
		Use:   "build [project...]",
		Short: "Generate files based on agent-def.yml",
		Long: `Generate files for one or more projects or user-level tasks as defined in agent-def.yml.
If no project names are provided, all projects will be processed.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr, err := processor.NewManager(".")
			if err != nil {
				return err
			}
			return mgr.Build(args, userOnly, watch, dryRun, force)
		},
	}

	cmd.Flags().BoolVar(&userOnly, "user", false, "Build only user-level tasks")
	cmd.Flags().BoolVar(&watch, "watch", false, "Watch for changes and rebuild")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be generated without writing files")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite without prompting for confirmation")

	return cmd
}
