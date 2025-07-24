package cli

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/uphy/agent-sync/internal/log"
	"github.com/urfave/cli/v3"
	"go.uber.org/zap"
)

// setupTestApp creates a test CLI app with all commands for integration testing
func setupTestApp() *cli.Command {
	// Create a test logger
	logConfig := log.DefaultConfig()
	logConfig.Enabled = true
	logConfig.Level = log.DebugLevel
	logger, _ := zap.NewDevelopment()

	// Create output writer
	output := &log.ConsoleOutput{
		Verbose: true,
		Color:   false,
	}

	// Create shared context
	sharedContext := &Context{
		Logger: logger,
		Output: output,
	}

	// Create root command with global flags and environment variable support
	rootCmd := &cli.Command{
		Name:  "agent-sync",
		Usage: "Convert context and command definitions for various AI agents",
		Description: `Agent Sync (agent-sync) is a tool for converting context and command
definitions for various AI agents like Claude and Roo.

It supports processing memory context files, command definitions,
and listing available agents.`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "output",
				Aliases:  []string{"o"},
				Usage:    "Output format (json, yaml, text)",
				Sources:  cli.EnvVars("AGENT_SYNC_OUTPUT"),
				Category: "Output",
			},
			&cli.BoolFlag{
				Name:     "verbose",
				Aliases:  []string{"v"},
				Usage:    "Enable verbose output",
				Sources:  cli.EnvVars("AGENT_SYNC_VERBOSE"),
				Category: "Logging",
			},
			&cli.StringFlag{
				Name:     "log-file",
				Usage:    "Log file path",
				Sources:  cli.EnvVars("AGENT_SYNC_LOG_FILE"),
				Category: "Logging",
			},
			&cli.StringFlag{
				Name:     "log-level",
				Usage:    "Log level (debug, info, warn, error)",
				Value:    "info",
				Sources:  cli.EnvVars("AGENT_SYNC_LOG_LEVEL"),
				Category: "Logging",
			},
			&cli.BoolFlag{
				Name:     "debug",
				Usage:    "Set log level to debug (shorthand for --log-level=debug)",
				Sources:  cli.EnvVars("AGENT_SYNC_DEBUG"),
				Category: "Logging",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "apply",
				Usage: "Generate files based on configuration",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "config",
						Aliases:  []string{"c"},
						Usage:    "Path to agent-sync.yml file or directory containing it",
						Value:    ".",
						Sources:  cli.EnvVars("AGENT_SYNC_CONFIG"),
						Category: "Config",
					},
					&cli.BoolFlag{
						Name:     "dry-run",
						Usage:    "Show what would be generated without writing files",
						Sources:  cli.EnvVars("AGENT_SYNC_DRY_RUN"),
						Category: "Output",
					},
					&cli.BoolFlag{
						Name:     "force",
						Aliases:  []string{"f"},
						Usage:    "Force overwrite without prompting for confirmation",
						Sources:  cli.EnvVars("AGENT_SYNC_FORCE"),
						Category: "Output",
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					return nil // Mock implementation for testing
				},
				Metadata: map[string]interface{}{
					"context": sharedContext,
				},
			},
			{
				Name:  "init",
				Usage: "Initialize a new agent-sync.yml configuration",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:     "force",
						Aliases:  []string{"f"},
						Usage:    "Force overwrite of existing files",
						Sources:  cli.EnvVars("AGENT_SYNC_INIT_FORCE"),
						Category: "Output",
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					return nil // Mock implementation for testing
				},
				Metadata: map[string]interface{}{
					"context": sharedContext,
				},
			},
		},
		Metadata: map[string]interface{}{
			"context": sharedContext,
		},
	}

	return rootCmd
}

// TestAgentDef runs integration tests for the agent-sync command.
// It discovers test directories in testdata/ and runs each test case.
func TestAgentDef(t *testing.T) {
	replace := os.Getenv("AGENT_SYNC_REPLACE") == "true"
	// Build the agent-sync binary first
	binaryPath, err := buildAgentDefBinary()
	if err != nil {
		t.Fatalf("Failed to build agent-sync binary: %v", err)
	}
	defer func() {
		if err := os.Remove(binaryPath); err != nil {
			t.Logf("Failed to remove binary %s: %v", binaryPath, err)
		}
	}() // Clean up the binary after tests

	// Find all test directories under testdata
	testDirs, err := findTestDirectories()
	if err != nil {
		t.Fatalf("Failed to find test directories: %v", err)
	}
	if len(testDirs) == 0 {
		t.Log("No test directories found in testdata/")
		return
	}
	// Run each test case
	for _, testDir := range testDirs {
		testName := filepath.Base(testDir)
		fmt.Println(testName)

		t.Run(testName, func(t *testing.T) {
			// Determine if this is an expected error case
			expectError := testName == "edge-cases" || testName == "mixed-format-test"

			// Create a temporary directory for the test
			tempDir, err := os.MkdirTemp("", "agent-sync-test-")
			if err != nil {
				t.Fatalf("Failed to create temporary directory: %v", err)
			}
			defer func() {
				if err := os.RemoveAll(tempDir); err != nil {
					t.Logf("Failed to remove temporary directory %s: %v", tempDir, err)
				}
			}()
			buildDir := filepath.Join(tempDir, "build")

			// Path to test directory containing agent-sync.yml
			testPath := filepath.Join(testDir, "test")

			// Copy test files to temporary directory
			if err := copyDir(testPath, tempDir); err != nil {
				t.Fatalf("Failed to copy test files to temporary directory: %v", err)
			}

			// Run the agent-sync binary with config flag pointing to the local directory
			cmd := exec.Command(binaryPath, "apply", "-f", "--config", ".")
			cmd.Dir = tempDir
			var stderr bytes.Buffer
			cmd.Stderr = &stderr
			output, err := cmd.Output()

			// Check if command result matches expected error state
			if expectError {
				if err == nil {
					t.Errorf("Expected command to fail for error test case, but it succeeded")
				}
				// For error cases, we don't compare directories
				return
			} else if err != nil {
				t.Fatalf("agent-sync command failed: %v\n%s\n%s", err, output, stderr.String())
			}

			// Compare output with expected output
			expectedDir := filepath.Join(testDir, "expected")
			if replace {
				if err := os.RemoveAll(expectedDir); err != nil {
					t.Fatalf("Failed to remove expected directory: %v", err)
				}
				if err := copyDir(tempDir, expectedDir); err != nil {
					t.Fatalf("Failed to copy output to expected directory: %v", err)
				}
				return
			}
			diffReport, err := compareDirectories(expectedDir, buildDir)
			if err != nil {
				t.Fatalf("Failed to compare directories: %v", err)
			}

			if diffReport != "" {
				t.Errorf("Differences found between expected and actual output:\n%s", diffReport)
			}
		})
	}
}

// findTestDirectories finds all test directories under testdata.
// Each test directory should have a "test" and "expected" subdirectory.
func findTestDirectories() ([]string, error) {
	var dirs []string

	// Get current working directory for debugging
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current working directory: %v\n", err)
	} else {
		fmt.Printf("Current working directory: %s\n", cwd)
	}

	// Now testdata is directly in the internal/cli directory
	testdataDir := "testdata"
	absTestdataDir, err := filepath.Abs(testdataDir)
	if err != nil {
		fmt.Printf("Error getting absolute path for testdata: %v\n", err)
	} else {
		fmt.Printf("Looking for testdata at: %s\n", absTestdataDir)
	}

	// Check if testdata directory exists
	if _, err := os.Stat(testdataDir); os.IsNotExist(err) {
		fmt.Printf("testdata directory not found\n")
		return nil, nil // No testdata directory, return empty list
	}

	// Find all directories under testdata
	entries, err := os.ReadDir(testdataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read testdata directory: %w", err)
	}

	fmt.Printf("Found %d entries in testdata directory\n", len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			fmt.Printf("Skipping non-directory entry: %s\n", entry.Name())
			continue
		}

		dirPath := filepath.Join(testdataDir, entry.Name())
		fmt.Printf("Checking directory: %s\n", dirPath)

		// Check if test directory has both "test" and "expected" subdirectories
		testExists := false
		expectedExists := false

		subEntries, err := os.ReadDir(dirPath)
		if err != nil {
			fmt.Printf("Error reading directory %s: %v\n", dirPath, err)
			continue
		}

		fmt.Printf("Found %d entries in %s\n", len(subEntries), dirPath)
		for _, subEntry := range subEntries {
			if !subEntry.IsDir() {
				continue
			}

			fmt.Printf("- Found subdirectory: %s\n", subEntry.Name())
			if subEntry.Name() == "test" {
				testExists = true
				fmt.Printf("  - Found 'test' directory\n")
			} else if subEntry.Name() == "expected" {
				expectedExists = true
				fmt.Printf("  - Found 'expected' directory\n")
			}
		}

		if testExists && expectedExists {
			fmt.Printf("Adding valid test directory: %s\n", dirPath)
			dirs = append(dirs, dirPath)
		} else {
			fmt.Printf("Directory %s does not have both 'test' and 'expected' subdirectories\n", dirPath)
		}
	}

	return dirs, nil
}

// copyDir recursively copies a directory to a destination.
func copyDir(src, dst string) error {
	// Get properties of source directory
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("error getting stats for source directory: %w", err)
	}

	// Create destination directory with same permissions
	if err = os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return fmt.Errorf("error creating destination directory: %w", err)
	}

	// Read source directory entries
	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("error reading source directory: %w", err)
	}

	// Copy each entry
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// Recursive copy for directories
			if err = copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Copy file
			if err = copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile copies a single file from src to dst.
func copyFile(src, dst string) error {
	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("error opening source file: %w", err)
	}
	defer func() {
		if err := srcFile.Close(); err != nil {
			// Just log the error since we're in a defer and can't return it
			fmt.Printf("Error closing source file: %v\n", err)
		}
	}()

	// Get source file info for permissions
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("error getting source file stats: %w", err)
	}

	// Create destination file
	dstFile, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("error creating destination file: %w", err)
	}
	defer func() {
		if err := dstFile.Close(); err != nil {
			// Just log the error since we're in a defer and can't return it
			fmt.Printf("Error closing destination file: %v\n", err)
		}
	}()

	// Copy file contents
	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("error copying file contents: %w", err)
	}

	return nil
}

// compareDirectories compares two directories recursively using the `diff -r` command.
func compareDirectories(expected, actual string) (string, error) {
	// Use diff command for directory comparison
	cmd := exec.Command("diff", "-r", expected, actual)
	output, err := cmd.CombinedOutput()

	// If diff command executed successfully with no differences
	if err == nil {
		return "", nil
	}

	// Exit status 1 from diff means differences were found (not an error for our purposes)
	exitErr, ok := err.(*exec.ExitError)
	if ok && exitErr.ExitCode() == 1 {
		return string(output), nil
	}

	// Any other error means the diff command failed to run properly
	return "", fmt.Errorf("error comparing directories with diff command: %w\nOutput: %s", err, string(output))
}

// buildAgentDefBinary builds the agent-sync binary to a temporary location and returns the path.
func buildAgentDefBinary() (string, error) {
	// Create a temporary file to store the binary
	tempFile, err := os.CreateTemp("", "agent-sync-test")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	if err := tempFile.Close(); err != nil {
		return "", fmt.Errorf("failed to close temporary file: %w", err)
	}

	binaryPath := tempFile.Name()
	fmt.Printf("Building agent-sync binary to: %s\n", binaryPath)

	// Get project root to find main.go
	projectRoot, err := findProjectRoot()
	if err != nil {
		return "", fmt.Errorf("failed to find project root: %w", err)
	}

	// Build the binary
	cmd := exec.Command("go", "build", "-o", binaryPath, filepath.Join(projectRoot, "cmd", "agent-sync", "main.go"))
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if output, err := cmd.Output(); err != nil {
		return "", fmt.Errorf("failed to build binary: %v\n%s\n%s", err, output, stderr.String())
	}

	return binaryPath, nil
}

// findProjectRoot walks up the directory tree from the current directory
// until it finds a go.mod file, which indicates the project root.
func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current working directory: %w", err)
	}

	for {
		// Check if go.mod exists in this directory
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		// Move up one directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root without finding go.mod
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("could not find project root")
}

// TestIntegrationWithCliV3 tests urfave/cli/v3 integration with in-memory commands
func TestIntegrationWithCliV3(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Define test cases
	testCases := []struct {
		name         string
		args         []string
		envVars      map[string]string
		expectOutput string
		expectError  bool
	}{
		{
			name:        "apply with dry-run flag",
			args:        []string{"apply", "--dry-run"},
			expectError: false,
		},
		{
			name: "apply with environment variables",
			args: []string{"apply"},
			envVars: map[string]string{
				"AGENT_SYNC_DRY_RUN": "true",
				"AGENT_SYNC_CONFIG":  "./testdata/default-project-test/test",
			},
			expectError: false,
		},
		{
			name:        "apply with config flag",
			args:        []string{"apply", "--config", "./testdata/default-project-test/test"},
			expectError: false,
		},
		{
			name:        "init command",
			args:        []string{"init"},
			expectError: false,
		},
		{
			name:        "init with force flag",
			args:        []string{"init", "--force"},
			expectError: false,
		},
	}

	// Run each test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app := setupTestApp()

			withEnvironment(t, tc.envVars, func() {
				_, err := executeCliCommand(t, app, tc.args...)

				if tc.expectError && err == nil {
					t.Errorf("Expected error, got none")
				}
				if !tc.expectError && err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			})
		})
	}
}

// TestFlagEnvironmentPrecedence tests flag vs environment variable precedence
func TestFlagEnvironmentPrecedence(t *testing.T) {
	// Create a command with config flag and environment variable
	var configPath string

	cmd := &cli.Command{
		Name: "test",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Config path",
				Value:   "default-path",
				Sources: cli.EnvVars("AGENT_SYNC_CONFIG"),
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			configPath = c.String("config")
			return nil
		},
	}

	// Test environment variable without flag
	withEnvironment(t, map[string]string{
		"AGENT_SYNC_CONFIG": "env-path",
	}, func() {
		err := cmd.Run(context.Background(), []string{"test"})
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if configPath != "env-path" {
			t.Errorf("Expected config to be 'env-path', got '%s'", configPath)
		}
	})

	// Test flag overrides environment variable
	withEnvironment(t, map[string]string{
		"AGENT_SYNC_CONFIG": "env-path",
	}, func() {
		err := cmd.Run(context.Background(), []string{"test", "--config", "flag-path"})
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if configPath != "flag-path" {
			t.Errorf("Expected config to be 'flag-path', got '%s'", configPath)
		}
	})
}
