package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/user/agent-def/internal/log"
)

func TestLoadConfig_Logging(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "agent-def-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Test 1: Configuration with logging section
	configWithLogging := `
configVersion: "1.0"
projects:
  test-project:
    destinations:
      - ./output
    tasks:
      - name: test-task
        type: memory
        sources:
          - ./memories
        targets:
          - agent: claude
user:
  tasks: []
logging:
  enabled: true
  level: "debug"
  file: "~/logs/agent-def.log"
  max_size: 20
  max_age: 14
  max_files: 5
  format: "json"
  console_output: true
  color: false
  verbose: true
`
	configPath := filepath.Join(tmpDir, "config_with_logging.yml")
	err = os.WriteFile(configPath, []byte(configWithLogging), 0644)
	assert.NoError(t, err)

	cfg, _, err := LoadConfig(configPath)
	assert.NoError(t, err)
	assert.NotNil(t, cfg.Logging)
	assert.True(t, cfg.Logging.Enabled)
	assert.Equal(t, log.DebugLevel, cfg.Logging.Level)
	assert.Equal(t, "json", string(cfg.Logging.Format))
	assert.Equal(t, 20, cfg.Logging.MaxSize)
	assert.Equal(t, 14, cfg.Logging.MaxAge)
	assert.Equal(t, 5, cfg.Logging.MaxFiles)
	assert.True(t, cfg.Logging.ConsoleOutput)
	assert.False(t, cfg.Logging.Color)
	assert.True(t, cfg.Logging.Verbose)

	// Test 2: Configuration without logging section (backward compatibility)
	configWithoutLogging := `
configVersion: "1.0"
projects:
  test-project:
    destinations:
      - ./output
    tasks:
      - name: test-task
        type: memory
        sources:
          - ./memories
        targets:
          - agent: claude
user:
  tasks: []
`
	configPath = filepath.Join(tmpDir, "config_without_logging.yml")
	err = os.WriteFile(configPath, []byte(configWithoutLogging), 0644)
	assert.NoError(t, err)

	cfg, _, err = LoadConfig(configPath)
	assert.NoError(t, err)
	assert.NotNil(t, cfg.Logging)
	// Default values should be set
	assert.False(t, cfg.Logging.Enabled)
	assert.Equal(t, log.InfoLevel, cfg.Logging.Level)
	assert.Equal(t, "", cfg.Logging.File)
	assert.Equal(t, log.TextFormat, cfg.Logging.Format)
}

func TestExpandTildeInLoggingConfig(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "agent-def-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Mock home directory for testing
	originalHome := os.Getenv("HOME")
	mockHome := filepath.Join(tmpDir, "mock-home")
	err = os.Mkdir(mockHome, 0755)
	assert.NoError(t, err)
	os.Setenv("HOME", mockHome)
	defer os.Setenv("HOME", originalHome)

	configWithTilde := `
configVersion: "1.0"
projects:
  test-project:
    destinations:
      - ./output
    tasks: []
user:
  tasks: []
logging:
  enabled: true
  level: "info"
  file: "~/logs/agent-def.log"
`
	configPath := filepath.Join(tmpDir, "config_with_tilde.yml")
	err = os.WriteFile(configPath, []byte(configWithTilde), 0644)
	assert.NoError(t, err)

	cfg, _, err := LoadConfig(configPath)
	assert.NoError(t, err)
	assert.NotNil(t, cfg.Logging)

	// Check if tilde was expanded in the file path
	expectedPath := filepath.Join(mockHome, "logs", "agent-def.log")
	assert.Equal(t, expectedPath, cfg.Logging.File)
}
