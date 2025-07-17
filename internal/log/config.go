package log

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
)

// Level represents the logging level
type Level string

// ログレベル定数
const (
	DebugLevel Level = "debug"
	InfoLevel  Level = "info"
	WarnLevel  Level = "warn"
	ErrorLevel Level = "error"
)

// Format はログ出力フォーマットを表す型
type Format string

// ログフォーマット定数
const (
	TextFormat Format = "text"
	JSONFormat Format = "json"
)

// Config represents logging configuration
type Config struct {
	// 基本設定
	Enabled bool  `yaml:"enabled"` // ロギングを有効にするか
	Level   Level `yaml:"level"`   // "debug", "info", "warn", "error"

	// ファイル出力設定
	File     string `yaml:"file"`      // ログファイルパス、空の場合はファイル出力なし
	MaxSize  int    `yaml:"max_size"`  // ログファイルの最大サイズ (MB)
	MaxAge   int    `yaml:"max_age"`   // ログファイルの保持日数
	MaxFiles int    `yaml:"max_files"` // 保持する古いログファイル数

	// 出力フォーマット
	Format Format `yaml:"format"` // "json", "text"など

	// コンソール出力設定
	ConsoleOutput bool `yaml:"console_output"` // コンソールにも出力するか
	Color         bool `yaml:"color"`          // カラー出力を使用するか
	Verbose       bool `yaml:"verbose"`        // 詳細出力を有効にするか
}

// DefaultConfig はデフォルトのログ設定を返す
func DefaultConfig() Config {
	return Config{
		Enabled:       false, // デフォルトで無効に変更
		Level:         InfoLevel,
		File:          "",
		MaxSize:       10, // 10MB
		MaxAge:        30, // 30日
		MaxFiles:      3,
		Format:        TextFormat,
		ConsoleOutput: false,
		Color:         true,
		Verbose:       false,
	}
}

// IsValidLevel はログレベルが有効かどうかを検証する
func (c *Config) IsValidLevel() bool {
	switch c.Level {
	case DebugLevel, InfoLevel, WarnLevel, ErrorLevel:
		return true
	default:
		return false
	}
}

// IsValidFormat はフォーマットが有効かどうかを検証する
func (c *Config) IsValidFormat() bool {
	switch c.Format {
	case TextFormat, JSONFormat:
		return true
	default:
		return false
	}
}

// Validate は設定が有効かどうかを検証する
func (c *Config) Validate() error {
	if !c.IsValidLevel() {
		return fmt.Errorf("invalid log level: %s", c.Level)
	}

	if !c.IsValidFormat() {
		return fmt.Errorf("invalid log format: %s", c.Format)
	}

	if c.File != "" {
		// ファイルパスのディレクトリ部分が存在するか確認
		dir := filepath.Dir(c.File)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return fmt.Errorf("log file directory does not exist: %s", dir)
		}
	}

	if c.MaxSize <= 0 {
		return fmt.Errorf("max_size must be positive, got: %d", c.MaxSize)
	}

	if c.MaxAge <= 0 {
		return fmt.Errorf("max_age must be positive, got: %d", c.MaxAge)
	}

	if c.MaxFiles < 0 {
		return fmt.Errorf("max_files cannot be negative, got: %d", c.MaxFiles)
	}

	return nil
}

// LoadConfigFromYAML loads logging configuration from the specified YAML file.
func LoadConfigFromYAML(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := DefaultConfig()
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &config, nil
}

// LoadConfigFromYAMLContent loads logging configuration from YAML content string.
func LoadConfigFromYAMLContent(content string) (*Config, error) {
	config := DefaultConfig()
	err := yaml.Unmarshal([]byte(content), &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config content: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &config, nil
}

// String はConfigの文字列表現を返す
func (c *Config) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Config{Enabled: %t, ", c.Enabled))
	sb.WriteString(fmt.Sprintf("Level: %s, ", c.Level))
	if c.File != "" {
		sb.WriteString(fmt.Sprintf("File: %s, ", c.File))
	}
	sb.WriteString(fmt.Sprintf("MaxSize: %dMB, ", c.MaxSize))
	sb.WriteString(fmt.Sprintf("MaxAge: %d days, ", c.MaxAge))
	sb.WriteString(fmt.Sprintf("MaxFiles: %d, ", c.MaxFiles))
	sb.WriteString(fmt.Sprintf("Format: %s, ", c.Format))
	sb.WriteString(fmt.Sprintf("ConsoleOutput: %t, ", c.ConsoleOutput))
	sb.WriteString(fmt.Sprintf("Color: %t, ", c.Color))
	sb.WriteString(fmt.Sprintf("Verbose: %t}", c.Verbose))
	return sb.String()
}
