package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/uphy/agent-sync/internal/agent"
	"github.com/uphy/agent-sync/internal/cli"
	"github.com/uphy/agent-sync/internal/log"
	"github.com/uphy/agent-sync/internal/util"
	"go.uber.org/zap"
)

var (
	// Version information - to be set during build
	version = "dev"
	commit  = "none"
	date    = "unknown"

	// Global flags
	outputFormat string
	verbose      bool
	logFile      string
	logLevel     string
	debugMode    bool
	logger       *zap.Logger
	context      *cli.Context
)

func main() {
	// Create and configure a centralized agent registry
	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})
	registry.Register(&agent.Roo{})

	rootCmd := &cobra.Command{
		Use:   "agent-sync",
		Short: "Agent Definition - Convert context and command definitions for various AI agents",
		Long: `Agent Definition (agent-sync) is a tool for converting context and command
definitions for various AI agents like Claude and Roo.

It supports processing memory context files, command definitions,
and listing available agents.`,
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
	}
	// Add global flags
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "", "Output format (json, yaml, text)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().StringVar(&logFile, "log-file", "", "Log file path")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().BoolVar(&debugMode, "debug", false, "Set log level to debug (shorthand for --log-level=debug)")

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {

		// Validate global flags
		if outputFormat != "" && outputFormat != "json" && outputFormat != "yaml" && outputFormat != "text" {
			return &util.ErrInvalidOutputFormat{Format: outputFormat}
		}
		// Initialize logging configuration based only on command line flags and environment variables
		logConfig := log.DefaultConfig()

		// Process command line flags first
		if logFile != "" {
			logConfig.File = logFile
			logConfig.Enabled = true
		}

		// Handle debug mode flag
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

		// Process environment variables (lower priority than command line flags)
		if envFile := os.Getenv("AGENT_SYNC_LOG_FILE"); envFile != "" && logFile == "" {
			logConfig.File = envFile
			logConfig.Enabled = true
		}

		if envLevel := os.Getenv("AGENT_SYNC_LOG_LEVEL"); envLevel != "" && logLevel == "" {
			logConfig.Level = log.Level(envLevel)
			logConfig.Enabled = true
		}

		// Initialize logger
		var err error
		logger, err = log.NewZapLogger(logConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing logger: %v\n", err)
			os.Exit(1)
		}
		defer func() {
			// Ignoring sync error during shutdown is common practice with zap
			// but we should at least log it if in debug mode
			if err := logger.Sync(); err != nil {
				if logConfig.Level == log.DebugLevel {
					fmt.Fprintf(os.Stderr, "Failed to sync logger: %v\n", err)
				}
			}
		}()

		// Initialize output writer
		output := &log.ConsoleOutput{
			Verbose: logConfig.Verbose,
			Color:   logConfig.Color,
		}

		// Test debug logging
		logger.Debug("MAIN DEBUG TEST: Logging configuration",
			zap.String("level", string(logConfig.Level)),
			zap.Bool("enabled", logConfig.Enabled),
			zap.Bool("verbose", logConfig.Verbose))

		// Create CLI context
		if context == nil {
			context = &cli.Context{}
		}
		context.Logger = logger
		context.Output = output

		return nil
	}
	// Register commands with context
	context = &cli.Context{}
	cli.SetupCommands(rootCmd, context)

	// Execute the root command with improved error handling
	if err := rootCmd.Execute(); err != nil {
		// Check for custom errors that might have specific formatting
		if customErr, ok := err.(util.CustomError); ok {
			fmt.Fprintf(os.Stderr, "Error: %s\n", customErr.FormattedError())
			if logger != nil {
				logger.Error("Command execution failed",
					zap.String("error", customErr.FormattedError()))
			}
		} else {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			if logger != nil {
				logger.Error("Command execution failed",
					zap.Error(err))
			}
		}
		os.Exit(1)
	}
}
