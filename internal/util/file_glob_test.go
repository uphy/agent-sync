package util

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

func TestGlobWithExcludes(t *testing.T) {
	// Create a temporary test directory structure
	testDir, err := os.MkdirTemp("", "glob-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(testDir); err != nil {
			t.Logf("Failed to remove test directory: %v", err)
		}
	}()

	// Create test file structure:
	// - file1.md
	// - file2.md
	// - file_test.md
	// - subdir/
	//   - nested1.md
	//   - nested_test.md
	//   - deepdir/
	//     - deep1.md
	//     - deep_test.md
	// - temp.md

	// Create top-level files
	createTestFile(t, testDir, "file1.md")
	createTestFile(t, testDir, "file2.md")
	createTestFile(t, testDir, "file_test.md")
	createTestFile(t, testDir, "temp.md")

	// Create subdirectory and files
	subDir := filepath.Join(testDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}
	createTestFile(t, subDir, "nested1.md")
	createTestFile(t, subDir, "nested_test.md")

	// Create deep directory and files
	deepDir := filepath.Join(subDir, "deepdir")
	if err := os.MkdirAll(deepDir, 0755); err != nil {
		t.Fatalf("Failed to create deep directory: %v", err)
	}
	createTestFile(t, deepDir, "deep1.md")
	createTestFile(t, deepDir, "deep_test.md")

	// Test cases
	testCases := []struct {
		name     string
		patterns []string
		expected []string
	}{
		{
			name:     "Simple glob",
			patterns: []string{"*.md"},
			expected: []string{"file1.md", "file2.md", "file_test.md", "temp.md"},
		},
		{
			name:     "Simple exclude",
			patterns: []string{"*.md", "!*_test.md"},
			expected: []string{"file1.md", "file2.md", "temp.md"},
		},
		{
			name:     "Recursive glob",
			patterns: []string{"**/*.md"},
			expected: []string{
				"file1.md", "file2.md", "file_test.md", "temp.md",
				"subdir/nested1.md", "subdir/nested_test.md",
				"subdir/deepdir/deep1.md", "subdir/deepdir/deep_test.md",
			},
		},
		{
			name:     "Recursive with exclude",
			patterns: []string{"**/*.md", "!**/*_test.md"},
			expected: []string{
				"file1.md", "file2.md", "temp.md",
				"subdir/nested1.md",
				"subdir/deepdir/deep1.md",
			},
		},
		{
			name:     "Multiple excludes",
			patterns: []string{"**/*.md", "!**/*_test.md", "!temp.md"},
			expected: []string{
				"file1.md", "file2.md",
				"subdir/nested1.md",
				"subdir/deepdir/deep1.md",
			},
		},
		{
			name:     "Subdirectory only",
			patterns: []string{"subdir/**/*.md"},
			expected: []string{
				"subdir/nested1.md", "subdir/nested_test.md",
				"subdir/deepdir/deep1.md", "subdir/deepdir/deep_test.md",
			},
		},
		{
			name:     "Empty patterns",
			patterns: []string{},
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the function
			result, err := GlobWithExcludes(tc.patterns, testDir)
			if err != nil {
				t.Fatalf("Error in GlobWithExcludes: %v", err)
			}

			// Normalize expected paths
			expected := make([]string, len(tc.expected))
			for i, path := range tc.expected {
				expected[i] = filepath.FromSlash(path)
			}
			sort.Strings(expected)

			// Compare results
			if !reflect.DeepEqual(result, expected) {
				t.Errorf("Expected %v, got %v", expected, result)
			}
		})
	}
}

func createTestFile(t *testing.T, dir, name string) {
	path := filepath.Join(dir, name)
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("Failed to create test file %s: %v", path, err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("Failed to close test file: %v", err)
	}
}
