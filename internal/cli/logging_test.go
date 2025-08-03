package cli

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uphy/agent-sync/internal/log"
	"github.com/urfave/cli/v3"
)

func TestInitializeLogging(t *testing.T) {
	testCases := []struct {
		name            string
		args            []string
		expectedLevel   string
		expectedVerbose bool
		expectError     bool
	}{
		{
			name:            "Default settings",
			args:            []string{"app"},
			expectedLevel:   "info",
			expectedVerbose: false,
			expectError:     false,
		},
		{
			name:            "With verbose flag",
			args:            []string{"app", "--verbose"},
			expectedLevel:   "debug",
			expectedVerbose: true,
			expectError:     false,
		},
		{
			name:            "With debug flag (shorthand forces verbose console)",
			args:            []string{"app", "--debug"},
			expectedLevel:   "debug",
			expectedVerbose: true,
			expectError:     false,
		},
		{
			name:            "With log level flag",
			args:            []string{"app", "--log-level", "warn"},
			expectedLevel:   "warn",
			expectedVerbose: false,
			expectError:     false,
		},
		{
			name:            "With invalid output format",
			args:            []string{"app", "--output", "invalid"},
			expectedLevel:   "info",
			expectedVerbose: false,
			expectError:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a test command with the required flags
			cmd := &cli.Command{
				Name: "test-app",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "output",
						Usage: "Output format",
					},
					&cli.BoolFlag{
						Name:  "verbose",
						Usage: "Enable verbose output",
					},
					&cli.StringFlag{
						Name:  "log-file",
						Usage: "Log file path",
					},
					&cli.StringFlag{
						Name:  "log-level",
						Usage: "Log level",
						Value: "info",
					},
					&cli.BoolFlag{
						Name:  "debug",
						Usage: "Set log level to debug",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					// We don't need to do anything here for the test
					return nil
				},
				Before: InitializeLogging,
				Metadata: map[string]interface{}{
					"context": &Context{},
				},
			}

			// Run the command with test arguments
			err := cmd.Run(context.Background(), tc.args)

			// Check error state
			if tc.expectError {
				assert.Error(t, err, "Expected error but got none")
			} else {
				assert.NoError(t, err, "Expected no error but got: %v", err)

				// If we don't expect an error, check the logger configuration
				// Get the shared context
				sharedContext := GetSharedContext(cmd)
				assert.NotNil(t, sharedContext, "Shared context should not be nil")

				// Check verbose setting
				output, ok := sharedContext.Output.(*log.ConsoleOutput)
				assert.True(t, ok, "Output should be of type ConsoleOutput")
				assert.Equal(t, tc.expectedVerbose, output.Verbose, "Verbose setting mismatch")

				// Get the logger
				logger := GetLogger(cmd)
				assert.NotNil(t, logger, "Logger should not be nil")
			}
		})
	}
}
