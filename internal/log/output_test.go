package log

import (
	"errors"
	"testing"
)

func TestConsoleOutput(t *testing.T) {
	// This test doesn't assert anything, it just ensures that methods don't panic
	// A more thorough test would capture stdout/stderr and verify output

	// Create a console output with color and verbose mode enabled
	co := &ConsoleOutput{
		Verbose: true,
		Color:   true,
	}

	// Test standard output methods
	co.Print("Standard message")
	co.Printf("Formatted message: %s", "test")

	// Test progress output
	co.PrintProgress("In progress")

	// Test success output
	co.PrintSuccess("Operation completed")

	// Test error output
	co.PrintError(errors.New("something went wrong"))

	// Test verbose output
	co.PrintVerbose("Verbose details")

	// Create a non-colored, non-verbose console output
	co = &ConsoleOutput{
		Verbose: false,
		Color:   false,
	}

	// Test the same methods with different settings
	co.Print("Standard message")
	co.Printf("Formatted message: %s", "test")
	co.PrintProgress("In progress")
	co.PrintSuccess("Operation completed")
	co.PrintError(errors.New("something went wrong"))
	co.PrintVerbose("This should not appear")

	// Note: We can't easily test Confirm() method in an automated test
	// because it requires user input.
}
