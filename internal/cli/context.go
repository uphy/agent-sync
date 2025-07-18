package cli

import (
	"github.com/spf13/cobra"
	"github.com/uphy/agent-def/internal/log"
	"go.uber.org/zap"
)

// Context holds shared resources and configuration for CLI commands.
// It provides consistent access to logging and output facilities.
type Context struct {
	// Logger is used for internal logging (debug, error tracking)
	Logger *zap.Logger

	// Output is used for user-facing console output
	Output log.OutputWriter
}

// SetupCommands registers all commands with the root command and provides them with context.
// This function wires up the command handlers with the shared CLI context.
func SetupCommands(rootCmd *cobra.Command, ctx *Context) {
	// Create command instances with context
	buildCmd := NewBuildCommandWithContext(ctx)
	initCmd := NewInitCommandWithContext(ctx)

	// Add commands to root
	rootCmd.AddCommand(buildCmd, initCmd)
}
