package cli

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"

	"github.com/uphy/agent-sync/internal/log"
	"github.com/urfave/cli/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

// executeCliCommand runs a CLI command with the given arguments and captures the output
func executeCliCommand(t *testing.T, cmd *cli.Command, args ...string) (string, error) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Prepare arguments (first arg should be the program name)
	testArgs := append([]string{"agent-sync"}, args...)

	// Run the command
	err := cmd.Run(context.Background(), testArgs)

	// Restore stdout
	if err := w.Close(); err != nil {
		t.Logf("Error closing pipe writer: %v", err)
	}
	os.Stdout = oldStdout

	// Capture output
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Logf("Error copying pipe output: %v", err)
	}

	return buf.String(), err
}

// TestContextCreation tests the creation of a CLI context with logger and output
func TestContextCreation(t *testing.T) {
	// Create a test logger that writes to a buffer
	var logBuf bytes.Buffer
	testCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(&logBuf),
		zap.DebugLevel,
	)
	testLogger := zap.New(testCore)

	// Create a mock output writer
	mockOutput := &mockOutputWriter{}

	// Create a context
	ctx := &Context{
		Logger: testLogger,
		Output: mockOutput,
	}

	// Verify the context has the expected components
	if ctx.Logger == nil {
		t.Error("Expected Logger to be set in Context")
	}
	if ctx.Output == nil {
		t.Error("Expected Output to be set in Context")
	}
}

// TestSetupCliV3Commands tests that commands are properly set up with context for urfave/cli
func TestSetupCliV3Commands(t *testing.T) {
	// Create a test logger
	testLogger := zaptest.NewLogger(t)

	// Create a mock output writer
	mockOutput := &mockOutputWriter{}

	// Create a context
	ctx := &Context{
		Logger: testLogger,
		Output: mockOutput,
	}

	// Create test commands
	cmd1 := &cli.Command{
		Name: "test1",
	}
	cmd2 := &cli.Command{
		Name: "test2",
	}
	commands := []*cli.Command{cmd1, cmd2}

	// Set up commands with context
	SetupCliV3Commands(commands, ctx)

	// Check that context was added to each command's metadata
	for _, cmd := range commands {
		if cmd.Metadata == nil {
			t.Errorf("Expected command %s to have metadata", cmd.Name)
			continue
		}

		ctxFromMeta, ok := cmd.Metadata["context"].(*Context)
		if !ok {
			t.Errorf("Expected command %s to have context in metadata", cmd.Name)
			continue
		}

		if ctxFromMeta != ctx {
			t.Errorf("Expected command %s to have the correct context in metadata", cmd.Name)
		}
	}
}

// TestGetCliContext tests retrieving the shared context from a CLI context
func TestGetCliContext(t *testing.T) {
	// Create a test logger
	testLogger := zaptest.NewLogger(t)

	// Create a mock output writer
	mockOutput := &mockOutputWriter{}

	// Create a context
	originalCtx := &Context{
		Logger: testLogger,
		Output: mockOutput,
	}

	// Create command with metadata containing the context
	cmd := &cli.Command{
		Name: "test",
		Metadata: map[string]interface{}{
			"context": originalCtx,
		},
	}

	// Helper function to get context from CLI command
	getContext := func(c *cli.Command) (*Context, bool) {
		if c == nil || c.Metadata == nil {
			return nil, false
		}

		ctx, ok := c.Metadata["context"].(*Context)
		return ctx, ok
	}

	// Retrieve the context
	retrievedCtx, ok := getContext(cmd)
	if !ok {
		t.Error("Expected to find context in CLI context metadata")
	}
	if retrievedCtx != originalCtx {
		t.Error("Expected retrieved context to match original context")
	}
}

// Helper function for testing with environment variables
func withEnvironment(t *testing.T, envVars map[string]string, testFunc func()) {
	// Store original environment
	originalEnv := make(map[string]string)
	for _, key := range []string{
		"AGENT_SYNC_OUTPUT",
		"AGENT_SYNC_VERBOSE",
		"AGENT_SYNC_LOG_FILE",
		"AGENT_SYNC_LOG_LEVEL",
		"AGENT_SYNC_DEBUG",
		"AGENT_SYNC_DRY_RUN",
		"AGENT_SYNC_FORCE",
		"AGENT_SYNC_CONFIG",
		"AGENT_SYNC_INIT_FORCE",
	} {
		if val, exists := os.LookupEnv(key); exists {
			originalEnv[key] = val
		}
	}

	// Set test environment variables
	for key, value := range envVars {
		if err := os.Setenv(key, value); err != nil {
			t.Fatalf("Failed to set environment variable %s: %v", key, err)
		}
	}

	// Ensure we clean up after the test
	defer func() {
		// Restore original environment
		for key := range envVars {
			if err := os.Unsetenv(key); err != nil {
				t.Logf("Failed to unset environment variable %s: %v", key, err)
			}
		}
		for key, value := range originalEnv {
			if err := os.Setenv(key, value); err != nil {
				t.Logf("Failed to restore environment variable %s: %v", key, err)
			}
		}
	}()

	// Run the test
	testFunc()
}

// TestEnvironmentVariables tests the processing of environment variables for logging
func TestEnvironmentVariables(t *testing.T) {
	withEnvironment(t, map[string]string{
		"AGENT_SYNC_LOG_FILE":  "/tmp/test.log",
		"AGENT_SYNC_LOG_LEVEL": "debug",
	}, func() {
		// Create a test config to verify environment variable processing
		config := log.DefaultConfig()

		// Create command with environment-aware flags
		cmd := &cli.Command{
			Name: "test",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "log-file",
					Usage:    "Log file path",
					Sources:  cli.EnvVars("AGENT_SYNC_LOG_FILE"),
					Value:    "",
					Category: "Logging",
				},
				&cli.StringFlag{
					Name:     "log-level",
					Usage:    "Log level",
					Sources:  cli.EnvVars("AGENT_SYNC_LOG_LEVEL"),
					Value:    "info",
					Category: "Logging",
				},
			},
			Action: func(ctx context.Context, c *cli.Command) error {
				// Update config from flags (which includes env vars)
				if logFile := c.String("log-file"); logFile != "" {
					config.File = logFile
				}
				if logLevel := c.String("log-level"); logLevel != "" {
					config.Level = log.Level(logLevel)
				}
				return nil
			},
		}

		// Run the command with no arguments to test environment variables
		err := cmd.Run(context.Background(), []string{"agent-sync"})
		if err != nil {
			t.Fatalf("Failed to run app: %v", err)
		}

		// Verify that the environment variables were processed correctly
		if config.File != "/tmp/test.log" {
			t.Errorf("Expected log file to be /tmp/test.log, got %s", config.File)
		}
		if config.Level != log.Level("debug") {
			t.Errorf("Expected log level to be debug, got %s", config.Level)
		}
	})
}

// TestFlagPrecedence tests that command-line flags take precedence over environment variables
func TestFlagPrecedence(t *testing.T) {
	withEnvironment(t, map[string]string{
		"AGENT_SYNC_LOG_FILE":  "/tmp/env.log",
		"AGENT_SYNC_LOG_LEVEL": "error",
	}, func() {
		// Create a test config
		config := log.DefaultConfig()

		// Create command with environment-aware flags
		cmd := &cli.Command{
			Name: "test",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "log-file",
					Usage:    "Log file path",
					Sources:  cli.EnvVars("AGENT_SYNC_LOG_FILE"),
					Value:    "",
					Category: "Logging",
				},
				&cli.StringFlag{
					Name:     "log-level",
					Usage:    "Log level",
					Sources:  cli.EnvVars("AGENT_SYNC_LOG_LEVEL"),
					Value:    "info",
					Category: "Logging",
				},
			},
			Action: func(ctx context.Context, c *cli.Command) error {
				// Update config from flags (which includes env vars)
				if logFile := c.String("log-file"); logFile != "" {
					config.File = logFile
				}
				if logLevel := c.String("log-level"); logLevel != "" {
					config.Level = log.Level(logLevel)
				}
				return nil
			},
		}

		// Run the command with CLI arguments that should override env vars
		err := cmd.Run(context.Background(), []string{"agent-sync", "--log-file", "/tmp/flag.log", "--log-level", "debug"})
		if err != nil {
			t.Fatalf("Failed to run app: %v", err)
		}

		// Verify that the flags take precedence
		if config.File != "/tmp/flag.log" {
			t.Errorf("Expected log file to be /tmp/flag.log, got %s", config.File)
		}
		if config.Level != log.Level("debug") {
			t.Errorf("Expected log level to be debug, got %s", config.Level)
		}
	})
}

// TestApplyCommand tests the apply command with both flags and environment variables
func TestApplyCommand(t *testing.T) {
	// Create a test apply command
	testContext := &Context{
		Logger: zaptest.NewLogger(t),
		Output: &mockOutputWriter{},
	}

	// Create a test flag to capture the dry-run value
	var dryRunCaptured bool

	// Create a test apply command
	rootCmd := &cli.Command{
		Name: "agent-sync",
		Commands: []*cli.Command{
			{
				Name: "apply",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:     "dry-run",
						Usage:    "Show what would be generated without writing files",
						Sources:  cli.EnvVars("AGENT_SYNC_DRY_RUN"),
						Value:    false,
						Category: "Output",
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					// Capture the flag value for testing
					dryRunCaptured = c.Bool("dry-run")
					return nil
				},
				Metadata: map[string]interface{}{
					"context": testContext,
				},
			},
		},
	}

	// Test with flag
	_, err := executeCliCommand(t, rootCmd, "apply", "--dry-run")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !dryRunCaptured {
		t.Errorf("Expected dry-run flag to be true when passed as CLI arg")
	}

	// Reset the captured value
	dryRunCaptured = false

	// Test with environment variable
	withEnvironment(t, map[string]string{
		"AGENT_SYNC_DRY_RUN": "true",
	}, func() {
		// Create a new fresh command for this test to avoid caching of flag values
		envRootCmd := &cli.Command{
			Name: "agent-sync",
			Commands: []*cli.Command{
				{
					Name: "apply",
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:     "dry-run",
							Usage:    "Show what would be generated without writing files",
							Sources:  cli.EnvVars("AGENT_SYNC_DRY_RUN"),
							Value:    false,
							Category: "Output",
						},
					},
					Action: func(ctx context.Context, c *cli.Command) error {
						// Capture the flag value for testing
						dryRunCaptured = c.Bool("dry-run")
						return nil
					},
					Metadata: map[string]interface{}{
						"context": testContext,
					},
				},
			},
		}

		_, err := executeCliCommand(t, envRootCmd, "apply")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if !dryRunCaptured {
			t.Errorf("Expected dry-run flag to be true when set via environment variable")
		}
	})
}

// mockOutputWriter is a mock implementation of the OutputWriter interface for testing
type mockOutputWriter struct {
	PrintCallCount         int
	PrintProgressCallCount int
	PrintSuccessCallCount  int
	PrintErrorCallCount    int
	PrintVerboseCallCount  int
	LastPrintMessage       string
}

func (m *mockOutputWriter) Print(msg string) {
	m.PrintCallCount++
	m.LastPrintMessage = msg
}

func (m *mockOutputWriter) Printf(format string, args ...interface{}) {
	m.PrintCallCount++
}

func (m *mockOutputWriter) PrintProgress(msg string) {
	m.PrintProgressCallCount++
	m.LastPrintMessage = msg
}

func (m *mockOutputWriter) PrintSuccess(msg string) {
	m.PrintSuccessCallCount++
	m.LastPrintMessage = msg
}

func (m *mockOutputWriter) PrintError(err error) {
	m.PrintErrorCallCount++
}

func (m *mockOutputWriter) PrintVerbose(msg string) {
	m.PrintVerboseCallCount++
	m.LastPrintMessage = msg
}

func (m *mockOutputWriter) Confirm(prompt string) bool {
	return true
}
