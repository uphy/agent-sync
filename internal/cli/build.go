package cli

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/uphy/agent-def/internal/log"
	"github.com/uphy/agent-def/internal/processor"
	"go.uber.org/zap"
)

// buildContext holds the CLI context for the build command
var buildContext *Context

// NewBuildCommand returns the 'build' command for the CLI.
func NewBuildCommand() *cobra.Command {
	return NewBuildCommandWithContext(nil)
}

// NewBuildCommandWithContext returns the 'build' command with the specified context.
func NewBuildCommandWithContext(ctx *Context) *cobra.Command {
	// Store context for command execution
	buildContext = ctx
	var dryRun bool
	var force bool
	var configPath string

	cmd := &cobra.Command{
		Use:   "build [project...]",
		Short: "Generate files based on agent-def.yml",
		Long: `Generate files for one or more projects or user-level tasks as defined in agent-def.yml.
If no project names are provided, all projects will be processed.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Use the context if available
			var logger *zap.Logger
			var output log.OutputWriter

			if buildContext != nil {
				logger = buildContext.Logger
				output = buildContext.Output

				// Log command execution
				logger.Info("Executing build command",
					zap.String("configPath", configPath),
					zap.Strings("projects", args),
					zap.Bool("dryRun", dryRun),
					zap.Bool("force", force))
			}

			// Initialize manager with context components
			absConfigPath, err := filepath.Abs(configPath)
			if err != nil {
				return err
			}

			mgr, err := processor.NewManager(absConfigPath, logger, output)
			if err != nil {
				return err
			}

			// Execute build
			return mgr.Build(args, dryRun, force)
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be generated without writing files")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite without prompting for confirmation")
	cmd.Flags().StringVarP(&configPath, "config", "c", ".", "Path to agent-def.yml file or directory containing it")

	return cmd
}
