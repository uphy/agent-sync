package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/uphy/agent-sync/internal/log"
	"github.com/uphy/agent-sync/internal/util"
	"github.com/urfave/cli/v3"
	"go.uber.org/zap"
)

// InitializeLogging sets up logging based on command-line flags and environment variables.
// This function is designed to be used as the Before hook in urfave/cli applications.
func InitializeLogging(ctx context.Context, cmd *cli.Command) (context.Context, error) {
	// Initialize logging configuration based on command line flags and environment variables
	logConfig := log.DefaultConfig()

	// Create a shared context for all commands
	sharedContext := &Context{}

	// Store the context in app metadata for command access
	if cmd.Metadata == nil {
		cmd.Metadata = make(map[string]interface{})
	}
	cmd.Metadata["context"] = sharedContext

	// Extract flag values
	// With urfave/cli, both flags and env vars are handled automatically
	// according to the flag definitions and EnvVars settings
	outputFormat := cmd.String("output")
	verbose := cmd.Bool("verbose")
	logFile := cmd.String("log-file")
	logLevel := cmd.String("log-level")
	debugMode := cmd.Bool("debug")

	// Validate output format if provided
	if outputFormat != "" && outputFormat != "json" && outputFormat != "yaml" && outputFormat != "text" {
		return ctx, &util.ErrInvalidOutputFormat{Format: outputFormat}
	}

	// Configure logging based on flags
	// Note: We don't need to manually check environment variables because
	// urfave/cli handles that automatically based on the flag definitions
	if logFile != "" {
		logConfig.File = logFile
		logConfig.Enabled = true
	}

	// Handle debug mode flag (takes precedence over log-level)
	if debugMode {
		logLevel = "debug"
		logConfig.Enabled = true
	}

	// Process log level
	if logLevel != "" {
		logConfig.Level = log.Level(logLevel)
		logConfig.Enabled = true
	}

	// Process verbose flag
	if verbose {
		if logLevel == "" || logLevel == "info" {
			// Set to debug level if verbose is specified without specific level
			logConfig.Level = log.DebugLevel
		}
		logConfig.Verbose = true
		logConfig.Enabled = true
		logConfig.ConsoleOutput = true
	}

	// Ensure console output is enabled for debug logs even when not in verbose mode
	if logConfig.Level == log.DebugLevel || debugMode {
		logConfig.ConsoleOutput = true
	}

	// Initialize logger
	var err error
	logger, err := log.NewZapLogger(logConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing logger: %v\n", err)
		os.Exit(1)
	}

	// Set up deferred logger sync
	// Note: urfave/cli has no direct support for defer in Before,
	// so we need to register a cleanup function with the app
	cmd.Metadata["loggerCleanup"] = func() {
		if err := logger.Sync(); err != nil {
			if logConfig.Level == log.DebugLevel {
				fmt.Fprintf(os.Stderr, "Failed to sync logger: %v\n", err)
			}
		}
	}

	// Register app shutdown hook for logger cleanup
	// This uses urfave/cli/v3's After hook for cleanup
	if cmd.After == nil {
		cmd.After = func(cleanupCtx context.Context, c *cli.Command) error {
			if cleanup, ok := c.Metadata["loggerCleanup"].(func()); ok {
				cleanup()
			}
			return nil
		}
	} else {
		// There's already an After hook, so we'll chain them
		originalAfter := cmd.After
		cmd.After = func(cleanupCtx context.Context, c *cli.Command) error {
			if cleanup, ok := c.Metadata["loggerCleanup"].(func()); ok {
				cleanup()
			}
			return originalAfter(cleanupCtx, c)
		}
	}

	// Initialize output writer
	output := &log.ConsoleOutput{
		Verbose: logConfig.Verbose,
		Color:   logConfig.Color,
	}

	// Populate shared context
	sharedContext.Logger = logger
	sharedContext.Output = output

	// Store logger in app metadata for global access if needed
	cmd.Metadata["logger"] = logger

	return ctx, nil
}

// GetSharedContext retrieves the shared context from command metadata
func GetSharedContext(cmd *cli.Command) *Context {
	if ctx, ok := cmd.Metadata["context"].(*Context); ok {
		return ctx
	}
	return nil
}

// GetLogger retrieves the logger directly from command metadata
func GetLogger(cmd *cli.Command) *zap.Logger {
	if logger, ok := cmd.Metadata["logger"].(*zap.Logger); ok {
		return logger
	}
	return nil
}
