package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// LoadConfig reads and parses the agent-def YAML configuration from the given path.
// If the provided path is a directory, it looks for "agent-def.yml" or "agent-def.yaml" in that directory.
func LoadConfig(path string) (*Config, error) {
	cfgPath := path
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("cannot stat path %q: %w", path, err)
	}

	if info.IsDir() {
		// search for agent-def.yml or agent-def.yaml
		for _, name := range []string{"agent-def.yml", "agent-def.yaml"} {
			p := filepath.Join(path, name)
			if _, err := os.Stat(p); err == nil {
				cfgPath = p
				break
			}
		}
	}

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %q: %w", cfgPath, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse YAML in %q: %w", cfgPath, err)
	}

	return &cfg, nil
}
