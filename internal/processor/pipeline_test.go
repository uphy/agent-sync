package processor

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/uphy/agent-sync/internal/config"
	"github.com/uphy/agent-sync/internal/util"
	"go.uber.org/zap"
)

// mockFileSystem mocks the util.FileSystem interface for testing
type mockFileSystem struct {
	existingFiles map[string]bool
	writtenFiles  map[string][]byte
	fileContents  map[string]string
}

func newMockFileSystem(existingFiles []string) *mockFileSystem {
	files := make(map[string]bool)
	contents := make(map[string]string)

	for _, f := range existingFiles {
		files[f] = true
		// Default content for existing files
		contents[f] = "mock content"
	}

	return &mockFileSystem{
		existingFiles: files,
		writtenFiles:  make(map[string][]byte),
		fileContents:  contents,
	}
}

// SetFileContent allows setting specific content for a file in tests
func (m *mockFileSystem) SetFileContent(path string, content string) {
	m.fileContents[path] = content
}

func (m *mockFileSystem) ReadFile(path string) ([]byte, error) {
	if !m.FileExists(path) {
		return nil, &util.ErrFileNotFound{Path: path}
	}

	// Return custom content if specified
	if content, ok := m.fileContents[path]; ok {
		return []byte(content), nil
	}

	return []byte("mock content"), nil
}

func (m *mockFileSystem) WriteFile(path string, data []byte) error {
	m.writtenFiles[path] = data
	return nil
}

func (m *mockFileSystem) FileExists(path string) bool {
	return m.existingFiles[path]
}

func (m *mockFileSystem) IsFile(path string) bool {
	return m.FileExists(path)
}

func (m *mockFileSystem) IsDir(path string) bool {
	return false
}

func (m *mockFileSystem) ListFiles(dir string, pattern string) ([]string, error) {
	return nil, nil
}

func (m *mockFileSystem) ResolvePath(path string) string {
	return path
}

func (m *mockFileSystem) GlobWithExcludes(patterns []string, baseDir string) ([]string, error) {
	// For testing, return fixed paths
	return []string{filepath.Join(baseDir, "test.md")}, nil
}

// mockOutputWriter mocks the log.OutputWriter interface for testing
type mockOutputWriter struct {
	messages []string
}

func newMockOutputWriter() *mockOutputWriter {
	return &mockOutputWriter{
		messages: make([]string, 0),
	}
}

func (m *mockOutputWriter) Print(msg string) {
	m.messages = append(m.messages, msg)
}

func (m *mockOutputWriter) Printf(format string, args ...interface{}) {
	m.messages = append(m.messages, fmt.Sprintf(format, args...))
}

func (m *mockOutputWriter) PrintProgress(msg string) {
	m.messages = append(m.messages, fmt.Sprintf("PROGRESS: %s", msg))
}

func (m *mockOutputWriter) PrintSuccess(msg string) {
	m.messages = append(m.messages, fmt.Sprintf("SUCCESS: %s", msg))
}

func (m *mockOutputWriter) PrintError(err error) {
	m.messages = append(m.messages, fmt.Sprintf("ERROR: %s", err))
}

func (m *mockOutputWriter) PrintVerbose(msg string) {
	m.messages = append(m.messages, fmt.Sprintf("VERBOSE: %s", msg))
}

func (m *mockOutputWriter) Confirm(prompt string) bool {
	return true
}

// TestPipelineDryRun tests the dry run functionality
func TestPipelineDryRun(t *testing.T) {
	// Test cases
	tests := []struct {
		name          string
		existingFiles []string
		fileContents  map[string]string
		wantCreated   int
		wantModified  int
		wantUnchanged int
	}{
		{
			name:          "All files new",
			existingFiles: []string{},
			fileContents:  map[string]string{},
			wantCreated:   2,
			wantModified:  0,
			wantUnchanged: 0,
		},
		{
			name:          "All files existing and different",
			existingFiles: []string{"/output/output1.md", "/output/output2.md"},
			fileContents:  map[string]string{},
			wantCreated:   0,
			wantModified:  2,
			wantUnchanged: 0,
		},
		{
			name:          "Mixed new and existing",
			existingFiles: []string{"/output/output1.md"},
			fileContents:  map[string]string{},
			wantCreated:   1,
			wantModified:  1,
			wantUnchanged: 0,
		},
		{
			name:          "All files unchanged",
			existingFiles: []string{"/output/output1.md", "/output/output2.md"},
			fileContents: map[string]string{
				"/output/output1.md": "Test content 1",
				"/output/output2.md": "Test content 2 with more characters",
			},
			wantCreated:   0,
			wantModified:  0,
			wantUnchanged: 2,
		},
		{
			name:          "Mixed unchanged, modified and new",
			existingFiles: []string{"/output/output1.md", "/output/output2.md"},
			fileContents: map[string]string{
				"/output/output1.md": "Test content 1",
				"/output/output2.md": "Old content that will change",
			},
			wantCreated:   0,
			wantModified:  1,
			wantUnchanged: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test objects
			mockFS := newMockFileSystem(tt.existingFiles)

			// Set file contents for comparison
			for path, content := range tt.fileContents {
				mockFS.SetFileContent(path, content)
			}
			mockOutput := newMockOutputWriter()
			logger := zap.NewNop()

			task := config.Task{
				Name:   "test-task",
				Type:   "memory",
				Inputs: []string{"*.md"},
				Outputs: []config.Output{
					{Agent: "roo", OutputPath: "output/"},
				},
			}

			// Create pipeline with dry run enabled
			pipeline := &Pipeline{
				Task:          task,
				AbsInputRoot:  "/input",
				AbsOutputDirs: []string{"/output"},
				UserScope:     false,
				DryRun:        true,
				Force:         false,
				fs:            mockFS,
				registry:      nil, // Not used in this test
				logger:        logger,
				output:        mockOutput,
			}

			// Mock files to process
			files := []ProcessedFile{
				{relPath: "output1.md", Content: "Test content 1"},
				{relPath: "output2.md", Content: "Test content 2 with more characters"},
			}

			// Execute writeOutputFiles method
			err := pipeline.writeOutputFiles(files)

			// Check for errors
			if err != nil {
				t.Fatalf("writeOutputFiles returned error: %v", err)
			}

			// No files should be actually written in dry run mode
			if len(mockFS.writtenFiles) > 0 {
				t.Errorf("Expected no files to be written in dry run mode, got %d", len(mockFS.writtenFiles))
			}

			// Verify output contains correct information
			createCount := 0
			modifyCount := 0
			unchangedCount := 0

			for _, msg := range mockOutput.messages {
				if strings.Contains(msg, "[CREATE]") {
					createCount++
				} else if strings.Contains(msg, "[MODIFY]") {
					modifyCount++
				} else if strings.Contains(msg, "[UNCHANGED]") {
					unchangedCount++
				}
			}

			// Check counts match expected values
			if createCount != tt.wantCreated {
				t.Errorf("Expected %d CREATE messages, got %d", tt.wantCreated, createCount)
			}

			if modifyCount != tt.wantModified {
				t.Errorf("Expected %d MODIFY messages, got %d", tt.wantModified, modifyCount)
			}

			if unchangedCount != tt.wantUnchanged {
				t.Errorf("Expected %d UNCHANGED messages, got %d", tt.wantUnchanged, unchangedCount)
			}

			// Check for size information
			sizeFound := false
			for _, msg := range mockOutput.messages {
				if strings.Contains(msg, "B)") || strings.Contains(msg, "KB)") || strings.Contains(msg, "MB)") {
					sizeFound = true
					break
				}
			}

			if !sizeFound {
				t.Error("Expected size information in output, none found")
			}

			// Check for summary message
			summaryFound := false
			for _, msg := range mockOutput.messages {
				if strings.Contains(msg, "Summary:") &&
					strings.Contains(msg, fmt.Sprintf("%d files would be created", tt.wantCreated)) &&
					strings.Contains(msg, fmt.Sprintf("%d files would be modified", tt.wantModified)) &&
					strings.Contains(msg, fmt.Sprintf("%d files would remain unchanged", tt.wantUnchanged)) {
					summaryFound = true
					break
				}
			}

			if !summaryFound {
				t.Error("Expected summary message in output, none found")
			}
		})
	}
}

// TestFormatDryRunFileStatus tests the formatDryRunFileStatus helper function
func TestFormatDryRunFileStatus(t *testing.T) {
	tests := []struct {
		name          string
		fileExists    bool
		unchanged     bool
		contentLength int
		want          string
		wantStatus    string
		wantSizeUnit  string
	}{
		{
			name:          "New small file",
			fileExists:    false,
			unchanged:     false,
			contentLength: 500,
			wantStatus:    "CREATE",
			wantSizeUnit:  "B",
		},
		{
			name:          "Modified KB file",
			fileExists:    true,
			unchanged:     false,
			contentLength: 2048,
			wantStatus:    "MODIFY",
			wantSizeUnit:  "KB",
		},
		{
			name:          "New MB file",
			fileExists:    false,
			unchanged:     false,
			contentLength: 3 * 1024 * 1024,
			wantStatus:    "CREATE",
			wantSizeUnit:  "MB",
		},
		{
			name:          "Unchanged small file",
			fileExists:    true,
			unchanged:     true,
			contentLength: 500,
			wantStatus:    "UNCHANGED",
			wantSizeUnit:  "B",
		},
		{
			name:          "Unchanged large file",
			fileExists:    true,
			unchanged:     true,
			contentLength: 1024 * 1024,
			wantStatus:    "UNCHANGED",
			wantSizeUnit:  "MB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock filesystem with or without the test file
			mockFS := &mockFileSystem{
				existingFiles: map[string]bool{
					"/test/file.txt": tt.fileExists,
				},
			}

			// Create pipeline with mock filesystem
			pipeline := &Pipeline{
				fs: mockFS,
			}

			// Execute formatDryRunFileStatus
			result := pipeline.formatDryRunFileStatus("/test/file.txt", tt.contentLength, tt.unchanged)

			// Verify status (CREATE/OVERWRITE)
			if !strings.Contains(result, "["+tt.wantStatus+"]") {
				t.Errorf("Expected status %s, got %s", tt.wantStatus, result)
			}

			// Verify file path is included
			if !strings.Contains(result, "/test/file.txt") {
				t.Errorf("Expected file path in result, got %s", result)
			}

			// Verify size unit is included
			if !strings.Contains(result, tt.wantSizeUnit) {
				t.Errorf("Expected size unit %s in result, got %s", tt.wantSizeUnit, result)
			}
		})
	}
}
