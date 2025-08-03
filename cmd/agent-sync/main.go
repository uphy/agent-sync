package main

import (
	"context"
	"fmt"
	"os"

	"github.com/uphy/agent-sync/internal/agent"
	internalcli "github.com/uphy/agent-sync/internal/cli"
	"github.com/uphy/agent-sync/internal/log"
	"github.com/uphy/agent-sync/internal/util"
	"github.com/urfave/cli/v3"
	"go.uber.org/zap"
)

var (
	// Version information - to be set during build
	version = "dev"
	commit  = "none"
	date    = "unknown"

	// Shared logger - will be initialized in Before hook
	logger *zap.Logger
)

func main() {
	// Create and configure agent registry
	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})
	registry.Register(&agent.Roo{})

	// Create shared context
	sharedContext := &internalcli.Context{}

	// Create root command with global flags and environment variable support
	rootCmd := &cli.Command{
		Name:  "agent-sync",
		Usage: "Convert context and command definitions for various AI agents",
		Description: `Agent Definition (agent-sync) is a tool for converting context and command
definitions for various AI agents like Claude and Roo.

It supports processing memory context files, command definitions,
and listing available agents.`,
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Output format (json, yaml, text)",
				Sources: cli.EnvVars("AGENT_SYNC_OUTPUT"),
				Action:  validateOutputFormat,
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Enable verbose output",
				Sources: cli.EnvVars("AGENT_SYNC_VERBOSE"),
			},
			&cli.StringFlag{
				Name:    "log-file",
				Usage:   "Log file path",
				Sources: cli.EnvVars("AGENT_SYNC_LOG_FILE"),
			},
			&cli.StringFlag{
				Name:    "log-level",
				Usage:   "Log level (debug, info, warn, error)",
				Value:   "info",
				Sources: cli.EnvVars("AGENT_SYNC_LOG_LEVEL"),
			},
			&cli.BoolFlag{
				Name:    "debug",
				Usage:   "Set log level to debug (shorthand for --log-level=debug)",
				Sources: cli.EnvVars("AGENT_SYNC_DEBUG"),
			},
		},
		Before: initializeLogging,
		Commands: []*cli.Command{
			internalcli.NewApplyCommand(),
			internalcli.NewInitCommand(),
		},
		Metadata: map[string]interface{}{
			"context": sharedContext,
		},
		ExitErrHandler: handleExitError,
	}
	// Propagate shared context to all subcommands so GetSharedContext works
	internalcli.SetupCliV3Commands(rootCmd.Commands, sharedContext)

	// Run the application
	if err := rootCmd.Run(context.Background(), os.Args); err != nil {
		// Error handling should be done by ExitErrHandler
		os.Exit(1)
	}
}

// validateOutputFormat validates the output format flag
func validateOutputFormat(ctx context.Context, cmd *cli.Command, format string) error {
	if format != "" && format != "json" && format != "yaml" && format != "text" {
		return &util.ErrInvalidOutputFormat{Format: format}
	}
	return nil
}

// initializeLogging sets up logging based on flags and environment variables
func initializeLogging(ctx context.Context, cmd *cli.Command) (context.Context, error) {
	// Initialize logging configuration
	logConfig := log.DefaultConfig()

	// Get shared context
	sharedContext, ok := cmd.Metadata["context"].(*internalcli.Context)
	if !ok {
		return ctx, fmt.Errorf("failed to get shared context")
	}

	// Extract flag values (both flags and env vars are handled automatically)
	verbose := cmd.Bool("verbose")
	logFile := cmd.String("log-file")
	logLevel := cmd.String("log-level")
	debugMode := cmd.Bool("debug")

	// Configure logging based on flags
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
	logger, err = log.NewZapLogger(logConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing logger: %v\n", err)
		os.Exit(1)
	}

	// Set up deferred logger sync via cmd's After hook
	registerLoggerCleanup(cmd, logger, logConfig.Level)

	// Initialize output writer
	output := &log.ConsoleOutput{
		Verbose: logConfig.Verbose,
		Color:   logConfig.Color,
	}

	// Populate shared context
	sharedContext.Logger = logger
	sharedContext.Output = output

	// Store logger in command metadata for global access
	cmd.Metadata["logger"] = logger

	return ctx, nil
}

// registerLoggerCleanup sets up a cleanup function for the logger
func registerLoggerCleanup(cmd *cli.Command, logger *zap.Logger, logLevel log.Level) {
	// Create cleanup function
	cleanup := func(ctx context.Context, cmd *cli.Command) error {
		if err := logger.Sync(); err != nil {
			if logLevel == log.DebugLevel {
				fmt.Fprintf(os.Stderr, "Failed to sync logger: %v\n", err)
			}
		}
		return nil
	}

	// Register cleanup with command's After hook
	if cmd.After == nil {
		cmd.After = cleanup
	} else {
		// Chain with existing After hook
		originalAfter := cmd.After
		cmd.After = func(ctx context.Context, cmd *cli.Command) error {
			err := cleanup(ctx, cmd)
			if err != nil {
				return err
			}
			return originalAfter(ctx, cmd)
		}
	}
}

// handleExitError provides custom error handling for the application
func handleExitError(ctx context.Context, cmd *cli.Command, err error) {
	if err == nil {
		return
	}

	// Get logger from context or use global logger
	var log *zap.Logger
	if cmd != nil && cmd.Metadata != nil {
		if l, ok := cmd.Metadata["logger"].(*zap.Logger); ok {
			log = l
		}
	}
	if log == nil {
		log = logger
	}

	// Handle custom errors
	if customErr, ok := err.(util.CustomError); ok {
		fmt.Fprintf(os.Stderr, "Error: %s\n", customErr.FormattedError())
		if log != nil {
			log.Error("Command execution failed",
				zap.String("error", customErr.FormattedError()))
		}
	} else {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		if log != nil {
			log.Error("Command execution failed",
				zap.Error(err))
		}
	}

	// Exit code is handled by the app.Run caller
}
