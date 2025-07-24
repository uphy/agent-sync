package cli

import (
	"github.com/uphy/agent-sync/internal/log"
	"github.com/urfave/cli/v3"
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

// SetupCliV3Commands registers all commands for the urfave/cli/v3 version
// and sets up the context in each command's metadata.
func SetupCliV3Commands(commands []*cli.Command, ctx *Context) {
	// Register context with all commands
	for _, cmd := range commands {
		if cmd.Metadata == nil {
			cmd.Metadata = make(map[string]interface{})
		}
		cmd.Metadata["context"] = ctx
	}
}
