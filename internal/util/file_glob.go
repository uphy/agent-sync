package util

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

// GlobWithExcludes expands a list of glob patterns, including negative patterns (prefixed with '!'),
// into a deterministically sorted list of file paths.
//
// Patterns like "**/*.md" match all markdown files recursively while "!**/*_test.md" excludes test files.
//
// The function returns paths relative to baseDir, sorted alphabetically.
func GlobWithExcludes(patterns []string, baseDir string) ([]string, error) {
	includePatterns := []string{}
	excludePatterns := []string{}

	// Separate include and exclude patterns
	for _, pattern := range patterns {
		if strings.HasPrefix(pattern, "!") {
			// This is an exclude pattern, remove the '!' prefix
			excludePatterns = append(excludePatterns, pattern[1:])
		} else {
			// This is an include pattern
			includePatterns = append(includePatterns, pattern)
		}
	}

	// If no include patterns are provided, return empty result
	if len(includePatterns) == 0 {
		return []string{}, nil
	}

	// Process include patterns
	matchedFiles := make(map[string]struct{})
	for _, pattern := range includePatterns {
		// Use doublestar.Glob directly for pattern matching
		matches, err := doublestar.Glob(os.DirFS(baseDir), pattern)
		if err != nil {
			return nil, fmt.Errorf("error expanding pattern '%s': %w", pattern, err)
		}

		for _, match := range matches {
			// Get absolute path for stat
			absPath := filepath.Join(baseDir, match)

			// Skip directories, only include files
			info, err := os.Stat(absPath)
			if err != nil {
				return nil, fmt.Errorf("error stating file '%s': %w", absPath, err)
			}

			if !info.IsDir() {
				// Store relative path
				matchedFiles[match] = struct{}{}
			}
		}
	}

	// Apply exclude patterns
	if len(excludePatterns) > 0 {
		for path := range matchedFiles {
			excluded := false
			for _, pattern := range excludePatterns {
				// Use doublestar.Match for each exclude pattern
				matched, err := doublestar.Match(pattern, path)
				if err == nil && matched {
					excluded = true
					break
				}
			}
			if excluded {
				delete(matchedFiles, path)
			}
		}
	}

	// Convert map to sorted slice
	result := make([]string, 0, len(matchedFiles))
	for path := range matchedFiles {
		result = append(result, path)
	}

	// Sort for deterministic output
	sort.Strings(result)

	return result, nil
}
