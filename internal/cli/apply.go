package cli

import (
	"context"
	"path/filepath"

	"github.com/uphy/agent-sync/internal/config"
	"github.com/uphy/agent-sync/internal/log"
	"github.com/uphy/agent-sync/internal/processor"
	"github.com/urfave/cli/v3"
	"go.uber.org/zap"
)

// NewApplyCommand returns the 'apply' command for urfave/cli.
func NewApplyCommand() *cli.Command {
	return &cli.Command{
		Name:        "apply",
		Usage:       "Generate files based on agent-sync.yml",
		Description: "Generate files for projects and user-level tasks as defined in agent-sync.yml.",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "dry-run",
				Usage:   "Preview files that would be generated (showing paths, sizes, and whether files would be created or overwritten)",
				Sources: cli.EnvVars("AGENT_SYNC_DRY_RUN"),
			},
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Force overwrite without prompting for confirmation",
				Sources: cli.EnvVars("AGENT_SYNC_FORCE"),
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Path to agent-sync.yml file or directory containing it",
				Value:   ".",
				Sources: cli.EnvVars("AGENT_SYNC_CONFIG"),
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// Access the shared context from metadata
			sharedContext := GetSharedContext(cmd)

			// Get command-specific flags
			dryRun := cmd.Bool("dry-run")
			force := cmd.Bool("force")
			configPath := cmd.String("config")

			var logger *zap.Logger
			var output log.OutputWriter

			// Get logger and output from shared context
			if sharedContext != nil {
				logger = sharedContext.Logger
				output = sharedContext.Output

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

			// Always validate config against embedded JSON Schema before applying
			if err := config.ValidateConfigFile(absConfigPath); err != nil {
				// Let the global handler format error output
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
}
