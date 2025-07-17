package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_Basic(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "agent-def-test")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("Failed to remove temporary directory: %v", err)
		}
	}()

	// Basic configuration without logging section
	configBasic := `
configVersion: "1.0"
projects:
  test-project:
    outputDirs:
      - ./output
    tasks:
      - name: test-task
        type: memory
        inputs:
          - ./memories
        outputs:
          - agent: claude
user:
  tasks: []
`
	configPath := filepath.Join(tmpDir, "config_basic.yml")
	err = os.WriteFile(configPath, []byte(configBasic), 0644)
	assert.NoError(t, err)

	cfg, _, err := LoadConfig(configPath)
	assert.NoError(t, err)
	assert.Equal(t, "1.0", cfg.ConfigVersion)
	assert.Contains(t, cfg.Projects, "test-project")
	assert.Equal(t, 1, len(cfg.Projects["test-project"].Tasks))
	assert.Equal(t, "test-task", cfg.Projects["test-project"].Tasks[0].Name)
}

func TestExpandTildeInConfig(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "agent-def-test")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("Failed to remove temporary directory: %v", err)
		}
	}()

	// Mock home directory for testing
	originalHome := os.Getenv("HOME")
	mockHome := filepath.Join(tmpDir, "mock-home")
	err = os.Mkdir(mockHome, 0755)
	assert.NoError(t, err)
	if err := os.Setenv("HOME", mockHome); err != nil {
		t.Fatalf("Failed to set HOME environment variable: %v", err)
	}
	defer func() {
		if err := os.Setenv("HOME", originalHome); err != nil {
			t.Logf("Failed to restore HOME environment variable: %v", err)
		}
	}()

	configWithTilde := `
configVersion: "1.0"
projects:
  test-project:
    outputDirs:
      - ~/output
    tasks:
      - name: test-task
        type: memory
        inputs:
          - ~/memories
        outputs:
          - agent: claude
            outputPath: ~/custom/path
user:
  tasks: []
`
	configPath := filepath.Join(tmpDir, "config_with_tilde.yml")
	err = os.WriteFile(configPath, []byte(configWithTilde), 0644)
	assert.NoError(t, err)

	cfg, _, err := LoadConfig(configPath)
	assert.NoError(t, err)

	// Check if tildes were expanded in various paths
	expectedOutputDir := filepath.Join(mockHome, "output")
	assert.Equal(t, expectedOutputDir, cfg.Projects["test-project"].OutputDirs[0])

	expectedInputPath := filepath.Join(mockHome, "memories")
	assert.Equal(t, expectedInputPath, cfg.Projects["test-project"].Tasks[0].Inputs[0])

	expectedOutputPath := filepath.Join(mockHome, "custom", "path")
	assert.Equal(t, expectedOutputPath, cfg.Projects["test-project"].Tasks[0].Outputs[0].OutputPath)
}
