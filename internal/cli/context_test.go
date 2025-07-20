package cli

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/uphy/agent-sync/internal/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

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

// TestSetupCommands tests that commands are properly set up with context
func TestSetupCommands(t *testing.T) {
	// Create a test logger
	testLogger := zaptest.NewLogger(t)

	// Create a mock output writer
	mockOutput := &mockOutputWriter{}

	// Create a context
	ctx := &Context{
		Logger: testLogger,
		Output: mockOutput,
	}

	// Create a root command
	rootCmd := &cobra.Command{
		Use: "testcmd",
	}

	// Set up commands with context
	SetupCommands(rootCmd, ctx)

	// Check that commands were added to the root command
	if len(rootCmd.Commands()) == 0 {
		t.Error("Expected commands to be added to the root command")
	}

	// Verify that buildContext and initContext were set
	if buildContext != ctx {
		t.Error("Expected buildContext to be set to the provided context")
	}
	if initContext != ctx {
		t.Error("Expected initContext to be set to the provided context")
	}
}

// TestEnvironmentVariables tests the processing of environment variables for logging
func TestEnvironmentVariables(t *testing.T) {
	// Save original environment
	origLogFile := os.Getenv("AGENT_DEF_LOG_FILE")
	origLogLevel := os.Getenv("AGENT_DEF_LOG_LEVEL")
	defer func() {
		// Restore original environment
		if err := os.Setenv("AGENT_DEF_LOG_FILE", origLogFile); err != nil {
			t.Logf("Failed to restore AGENT_DEF_LOG_FILE: %v", err)
		}
		if err := os.Setenv("AGENT_DEF_LOG_LEVEL", origLogLevel); err != nil {
			t.Logf("Failed to restore AGENT_DEF_LOG_LEVEL: %v", err)
		}
	}()

	// Set test environment variables
	testLogFile := "/tmp/test.log"
	testLogLevel := "debug"
	if err := os.Setenv("AGENT_DEF_LOG_FILE", testLogFile); err != nil {
		t.Fatalf("Failed to set AGENT_DEF_LOG_FILE: %v", err)
	}
	if err := os.Setenv("AGENT_DEF_LOG_LEVEL", testLogLevel); err != nil {
		t.Fatalf("Failed to set AGENT_DEF_LOG_LEVEL: %v", err)
	}

	// Create a test config to verify environment variable processing
	config := log.DefaultConfig()

	// Process environment variables (similar to what main.go does)
	if envFile := os.Getenv("AGENT_DEF_LOG_FILE"); envFile != "" {
		config.File = envFile
	}
	if envLevel := os.Getenv("AGENT_DEF_LOG_LEVEL"); envLevel != "" {
		config.Level = log.Level(envLevel)
	}

	// Verify that the environment variables were processed correctly
	if config.File != testLogFile {
		t.Errorf("Expected log file to be %s, got %s", testLogFile, config.File)
	}
	if config.Level != log.Level(testLogLevel) {
		t.Errorf("Expected log level to be %s, got %s", testLogLevel, config.Level)
	}
}

// TestFlagPrecedence tests that command-line flags take precedence over environment variables
func TestFlagPrecedence(t *testing.T) {
	// Set environment variables
	if err := os.Setenv("AGENT_DEF_LOG_FILE", "/tmp/env.log"); err != nil {
		t.Fatalf("Failed to set AGENT_DEF_LOG_FILE: %v", err)
	}
	if err := os.Setenv("AGENT_DEF_LOG_LEVEL", "error"); err != nil {
		t.Fatalf("Failed to set AGENT_DEF_LOG_LEVEL: %v", err)
	}
	defer func() {
		if err := os.Unsetenv("AGENT_DEF_LOG_FILE"); err != nil {
			t.Logf("Failed to unset AGENT_DEF_LOG_FILE: %v", err)
		}
		if err := os.Unsetenv("AGENT_DEF_LOG_LEVEL"); err != nil {
			t.Logf("Failed to unset AGENT_DEF_LOG_LEVEL: %v", err)
		}
	}()

	// Define command line flags (similar to what main.go does)
	var logFile = "/tmp/flag.log" // Simulating CLI flag value
	var logLevel = "debug"        // Simulating CLI flag value

	// Create a test config
	config := log.DefaultConfig()

	// Process flags (similar to what main.go does)
	if logFile != "" {
		config.File = logFile
	}
	if logLevel != "" {
		config.Level = log.Level(logLevel)
	}

	// Process environment variables if flags not set
	// (in this test, flags are "set" so env vars should be ignored)
	if envFile := os.Getenv("AGENT_DEF_LOG_FILE"); envFile != "" && logFile == "" {
		config.File = envFile
	}
	if envLevel := os.Getenv("AGENT_DEF_LOG_LEVEL"); envLevel != "" && logLevel == "" {
		config.Level = log.Level(envLevel)
	}

	// Verify that the flags take precedence
	if config.File != "/tmp/flag.log" {
		t.Errorf("Expected log file to be /tmp/flag.log, got %s", config.File)
	}
	if config.Level != log.Level("debug") {
		t.Errorf("Expected log level to be debug, got %s", config.Level)
	}
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
