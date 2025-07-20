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

	// Helper function to collect files matching a pattern
	collectFiles := func(pattern string) (map[string]struct{}, error) {
		result := make(map[string]struct{})
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
				result[match] = struct{}{}
			}
		}
		return result, nil
	}

	// Process include patterns
	matchedFiles := make(map[string]struct{})
	for _, pattern := range includePatterns {
		includedFiles, err := collectFiles(pattern)
		if err != nil {
			return nil, fmt.Errorf("error processing include pattern '%s': %w", pattern, err)
		}

		// Add included files to the result set
		for path := range includedFiles {
			matchedFiles[path] = struct{}{}
		}
	}

	// Apply exclude patterns - use doublestar.Glob for exclude patterns too
	if len(excludePatterns) > 0 {
		// Collect files matching exclude patterns
		for _, pattern := range excludePatterns {
			excludedFiles, err := collectFiles(pattern)
			if err != nil {
				return nil, fmt.Errorf("error processing exclude pattern '%s': %w", pattern, err)
			}

			// Remove excluded files from included files
			for path := range excludedFiles {
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

func GlobWithExcludesNoBaseDir(patterns []string) ([]string, error) {
	baseToPatterns := make(map[string][]string)
	for _, path := range patterns {
		// Temporarily remove the '!' prefix for negative patterns
		var prefix string
		if after, found := strings.CutPrefix(path, "!"); found {
			path = after // Remove the '!' prefix for negative patterns
			prefix = "!"
		}

		base, pattern := doublestar.SplitPattern(path)
		baseToPatterns[base] = append(baseToPatterns[base], prefix+pattern)
	}

	result := make([]string, 0)
	for base, patterns := range baseToPatterns {
		paths, err := GlobWithExcludes(patterns, base)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve template paths for base %q: %w", base, err)
		}
		for i, p := range paths {
			paths[i] = filepath.Join(base, p)
		}
		result = append(result, paths...)
	}
	return result, nil
}
