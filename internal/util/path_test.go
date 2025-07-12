package util

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExpandTilde(t *testing.T) {
	// Get the user's home directory for comparison in tests
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatal("Failed to get user home directory for tests:", err)
	}

	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{
			name:     "No tilde",
			input:    "/absolute/path",
			expected: "/absolute/path",
			hasError: false,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
			hasError: false,
		},
		{
			name:     "Just tilde",
			input:    "~",
			expected: homeDir,
			hasError: false,
		},
		{
			name:     "Tilde with separator",
			input:    "~/some/path",
			expected: filepath.Join(homeDir, "some/path"),
			hasError: false,
		},
		{
			name:     "Tilde without separator",
			input:    "~user",
			expected: filepath.Join(homeDir, "user"),
			hasError: false,
		},
		{
			name:     "Relative path",
			input:    "./relative/path",
			expected: "./relative/path",
			hasError: false,
		},
		{
			name:     "Tilde in the middle",
			input:    "/path/with/~/tilde",
			expected: "/path/with/~/tilde",
			hasError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ExpandTilde(tc.input)

			// Check error condition
			if tc.hasError && err == nil {
				t.Fatal("Expected error but got none")
			}

			if !tc.hasError && err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Check result value
			if result != tc.expected {
				t.Errorf("Expected %q but got %q", tc.expected, result)
			}
		})
	}
}
