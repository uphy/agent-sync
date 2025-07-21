package cli

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/uphy/agent-sync/internal/log"
	"github.com/uphy/agent-sync/internal/processor"
	"go.uber.org/zap"
)

// applyContext holds the CLI context for the apply command
var applyContext *Context

// NewApplyCommand returns the 'apply' command for the CLI.
func NewApplyCommand() *cobra.Command {
	return NewApplyCommandWithContext(nil)
}

// NewApplyCommandWithContext returns the 'apply' command with the specified context.
func NewApplyCommandWithContext(ctx *Context) *cobra.Command {
	// Store context for command execution
	applyContext = ctx
	var dryRun bool
	var force bool
	var configPath string

	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Generate files based on agent-sync.yml",
		Long:  `Generate files for projects and user-level tasks as defined in agent-sync.yml.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Use the context if available
			var logger *zap.Logger
			var output log.OutputWriter

			if applyContext != nil {
				logger = applyContext.Logger
				output = applyContext.Output

				// Log command execution
				logger.Info("Executing apply command",
					zap.String("configPath", configPath),
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

			// Execute apply
			return mgr.Apply(dryRun, force)
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview files that would be generated (showing paths, sizes, and whether files would be created or overwritten)")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite without prompting for confirmation")
	cmd.Flags().StringVarP(&configPath, "config", "c", ".", "Path to agent-sync.yml file or directory containing it")

	return cmd
}
