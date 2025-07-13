package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/agent-def/internal/agent"
	"github.com/user/agent-def/internal/cli"
	"github.com/user/agent-def/internal/util"
)

var (
	// Version information - to be set during build
	version = "dev"
	commit  = "none"
	date    = "unknown"

	// Global flags
	outputFormat string
	verbose      bool
)

func main() {
	// Create and configure a centralized agent registry
	registry := agent.NewRegistry()
	registry.Register(&agent.Claude{})
	registry.Register(&agent.Roo{})

	rootCmd := &cobra.Command{
		Use:   "agent-def",
		Short: "Agent Definition - Convert context and command definitions for various AI agents",
		Long: `Agent Definition (agent-def) is a tool for converting context and command
definitions for various AI agents like Claude and Roo.

It supports processing memory context files, command definitions,
and listing available agents.`,
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Validate global flags
			if outputFormat != "" && outputFormat != "json" && outputFormat != "yaml" && outputFormat != "text" {
				return &util.ErrInvalidOutputFormat{Format: outputFormat}
			}
			return nil
		},
	}

	// Add global flags
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "", "Output format (json, yaml, text)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	// Register commands with the registry passed in
	rootCmd.AddCommand(
		cli.NewBuildCommand(),
		cli.NewInitCommand(),
	)

	// Execute the root command with improved error handling
	if err := rootCmd.Execute(); err != nil {
		// Check for custom errors that might have specific formatting
		if customErr, ok := err.(util.CustomError); ok {
			fmt.Fprintf(os.Stderr, "Error: %s\n", customErr.FormattedError())
		} else {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		}
		os.Exit(1)
	}
}
