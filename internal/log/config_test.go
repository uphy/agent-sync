package log

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.Enabled {
		t.Error("DefaultLogConfig: Enabled should be false")
	}
	if config.Level != InfoLevel {
		t.Errorf("DefaultLogConfig: Level should be info, got %s", config.Level)
	}
	if config.File != "" {
		t.Errorf("DefaultLogConfig: File should be empty, got %s", config.File)
	}
	if config.MaxSize != 10 {
		t.Errorf("DefaultLogConfig: MaxSize should be 10, got %d", config.MaxSize)
	}
	if config.MaxAge != 30 {
		t.Errorf("DefaultLogConfig: MaxAge should be 30, got %d", config.MaxAge)
	}
	if config.MaxFiles != 3 {
		t.Errorf("DefaultLogConfig: MaxFiles should be 3, got %d", config.MaxFiles)
	}
	if config.Format != TextFormat {
		t.Errorf("DefaultLogConfig: Format should be text, got %s", config.Format)
	}
	if config.ConsoleOutput {
		t.Error("DefaultLogConfig: ConsoleOutput should be false")
	}
	if !config.Color {
		t.Error("DefaultLogConfig: Color should be true")
	}
	if config.Verbose {
		t.Error("DefaultLogConfig: Verbose should be false")
	}
}

func TestLogLevelValidation(t *testing.T) {
	tests := []struct {
		level    Level
		expected bool
	}{
		{DebugLevel, true},
		{InfoLevel, true},
		{WarnLevel, true},
		{ErrorLevel, true},
		{Level("invalid"), false},
	}

	for _, test := range tests {
		config := DefaultConfig()
		config.Level = test.level
		if result := config.IsValidLevel(); result != test.expected {
			t.Errorf("IsValidLevel for %s: expected %v, got %v", test.level, test.expected, result)
		}
	}
}

func TestLogFormatValidation(t *testing.T) {
	tests := []struct {
		format   Format
		expected bool
	}{
		{TextFormat, true},
		{JSONFormat, true},
		{Format("invalid"), false},
	}

	for _, test := range tests {
		config := DefaultConfig()
		config.Format = test.format
		if result := config.IsValidFormat(); result != test.expected {
			t.Errorf("IsValidFormat for %s: expected %v, got %v", test.format, test.expected, result)
		}
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		modifyFunc  func(*Config)
		expectError bool
	}{
		{
			name:        "Valid default config",
			modifyFunc:  func(c *Config) {},
			expectError: false,
		},
		{
			name: "Invalid level",
			modifyFunc: func(c *Config) {
				c.Level = Level("invalid")
			},
			expectError: true,
		},
		{
			name: "Invalid format",
			modifyFunc: func(c *Config) {
				c.Format = Format("invalid")
			},
			expectError: true,
		},
		{
			name: "Invalid MaxSize",
			modifyFunc: func(c *Config) {
				c.MaxSize = 0
			},
			expectError: true,
		},
		{
			name: "Invalid MaxAge",
			modifyFunc: func(c *Config) {
				c.MaxAge = 0
			},
			expectError: true,
		},
		{
			name: "Invalid MaxFiles",
			modifyFunc: func(c *Config) {
				c.MaxFiles = -1
			},
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := DefaultConfig()
			test.modifyFunc(&config)

			err := config.Validate()
			if test.expectError && err == nil {
				t.Errorf("Expected validation error, but got nil")
			}
			if !test.expectError && err != nil {
				t.Errorf("Did not expect validation error, but got: %v", err)
			}
		})
	}
}

func TestLoadFromYAMLContent(t *testing.T) {
	yamlContent := `
enabled: false
level: debug
file: /tmp/test.log
max_size: 20
max_age: 60
max_files: 5
format: json
console_output: true
color: false
verbose: true
`

	config, err := LoadConfigFromYAMLContent(yamlContent)
	if err != nil {
		t.Fatalf("Failed to load config from YAML content: %v", err)
	}

	if config.Enabled {
		t.Error("Config.Enabled should be false")
	}
	if config.Level != DebugLevel {
		t.Errorf("Config.Level should be debug, got %s", config.Level)
	}
	if config.File != "/tmp/test.log" {
		t.Errorf("Config.File should be /tmp/test.log, got %s", config.File)
	}
	if config.MaxSize != 20 {
		t.Errorf("Config.MaxSize should be 20, got %d", config.MaxSize)
	}
	if config.MaxAge != 60 {
		t.Errorf("Config.MaxAge should be 60, got %d", config.MaxAge)
	}
	if config.MaxFiles != 5 {
		t.Errorf("Config.MaxFiles should be 5, got %d", config.MaxFiles)
	}
	if config.Format != JSONFormat {
		t.Errorf("Config.Format should be json, got %s", config.Format)
	}
	if !config.ConsoleOutput {
		t.Error("Config.ConsoleOutput should be true")
	}
	if config.Color {
		t.Error("Config.Color should be false")
	}
	if !config.Verbose {
		t.Error("Config.Verbose should be true")
	}
}

func TestLoadFromYAMLFile(t *testing.T) {
	// 一時ディレクトリとファイルを作成
	tmpDir, err := os.MkdirTemp("", "log-config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("Failed to remove temporary directory %s: %v", tmpDir, err)
		}
	}()

	configPath := filepath.Join(tmpDir, "test-config.yml")
	yamlContent := []byte(`
enabled: false
level: warn
max_size: 15
max_age: 45
max_files: 4
format: json
console_output: true
color: false
`)

	err = os.WriteFile(configPath, yamlContent, 0644)
	if err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	config, err := LoadConfigFromYAML(configPath)
	if err != nil {
		t.Fatalf("Failed to load config from YAML file: %v", err)
	}

	if config.Enabled {
		t.Error("Config.Enabled should be false")
	}
	if config.Level != WarnLevel {
		t.Errorf("Config.Level should be warn, got %s", config.Level)
	}
	if config.MaxSize != 15 {
		t.Errorf("Config.MaxSize should be 15, got %d", config.MaxSize)
	}
	if config.MaxAge != 45 {
		t.Errorf("Config.MaxAge should be 45, got %d", config.MaxAge)
	}
	if config.MaxFiles != 4 {
		t.Errorf("Config.MaxFiles should be 4, got %d", config.MaxFiles)
	}
	if config.Format != JSONFormat {
		t.Errorf("Config.Format should be json, got %s", config.Format)
	}
	if !config.ConsoleOutput {
		t.Error("Config.ConsoleOutput should be true")
	}
	if config.Color {
		t.Error("Config.Color should be false")
	}
}

func TestLoadFromNonExistentFile(t *testing.T) {
	_, err := LoadConfigFromYAML("/non/existent/path.yml")
	if err == nil {
		t.Error("Expected error when loading from non-existent file, but got nil")
	}
}

func TestLoadInvalidYAMLContent(t *testing.T) {
	_, err := LoadConfigFromYAMLContent("invalid: : yaml: content")
	if err == nil {
		t.Error("Expected error when loading invalid YAML content, but got nil")
	}
}

func TestConfigString(t *testing.T) {
	config := Config{
		Enabled:       true,
		Level:         InfoLevel,
		File:          "/var/log/app.log",
		MaxSize:       10,
		MaxAge:        30,
		MaxFiles:      3,
		Format:        TextFormat,
		ConsoleOutput: true,
		Color:         false,
		Verbose:       true,
	}

	str := config.String()

	// 基本的な内容が含まれているか確認
	expectedSubstrings := []string{
		"Config{",
		"Enabled: true",
		"Level: info",
		"File: /var/log/app.log",
		"MaxSize: 10MB",
		"MaxAge: 30 days",
		"MaxFiles: 3",
		"Format: text",
		"ConsoleOutput: true",
		"Color: false",
		"Verbose: true",
	}

	for _, substr := range expectedSubstrings {
		if !contains(str, substr) {
			t.Errorf("String() output does not contain expected substring: %s", substr)
		}
	}
}

// テストヘルパー関数
func contains(str, substr string) bool {
	return strings.Contains(str, substr)
}
