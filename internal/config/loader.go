package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/user/agent-def/internal/util"
)

// LoadConfig reads and parses the agent-def YAML configuration from the given path.
// If the provided path is a directory, it looks for "agent-def.yml" or "agent-def.yaml" in that directory.
func LoadConfig(path string) (*Config, string, error) {
	cfgPath := path
	var cfgDir string

	info, err := os.Stat(path)
	if err != nil {
		return nil, "", fmt.Errorf("cannot stat path %q: %w", path, err)
	}

	if info.IsDir() {
		cfgDir = path
		// search for agent-def.yml or agent-def.yaml
		for _, name := range []string{"agent-def.yml", "agent-def.yaml"} {
			p := filepath.Join(path, name)
			if _, err := os.Stat(p); err == nil {
				cfgPath = p
				break
			}
		}
	} else {
		// If path points directly to a file, get its directory
		cfgDir = filepath.Dir(path)
	}

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read config file %q: %w", cfgPath, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, "", fmt.Errorf("failed to parse YAML in %q: %w", cfgPath, err)
	}

	// Normalize all paths in the config by expanding tildes (~) to home directory
	if err := expandTildeInConfig(&cfg); err != nil {
		return nil, "", fmt.Errorf("failed to expand tildes in paths: %w", err)
	}

	return &cfg, cfgDir, nil
}

// expandTildeInConfig normalizes all paths in the config by expanding tildes (~)
// to the user's home directory.
func expandTildeInConfig(cfg *Config) error {
	// Process all projects
	for projName, proj := range cfg.Projects {
		var err error

		// Expand tilde in Project.Root
		if proj.Root != "" {
			proj.Root, err = util.ExpandTilde(proj.Root)
			if err != nil {
				return fmt.Errorf("failed to expand tilde in project %s root: %w", projName, err)
			}
		}

		// Expand tilde in Project.Destinations
		for i, dest := range proj.Destinations {
			expanded, err := util.ExpandTilde(dest)
			if err != nil {
				return fmt.Errorf("failed to expand tilde in project %s destination: %w", projName, err)
			}
			proj.Destinations[i] = expanded
		}

		// Process Task.Sources and Target.Target in each Task
		for i, task := range proj.Tasks {
			// Expand tilde in Task.Sources
			for j, src := range task.Sources {
				expanded, err := util.ExpandTilde(src)
				if err != nil {
					return fmt.Errorf("failed to expand tilde in source for project %s task %s: %w",
						projName, task.Name, err)
				}
				task.Sources[j] = expanded
			}

			// Expand tilde in Target.Target for each Target
			for j, target := range task.Targets {
				if target.Target != "" {
					target.Target, err = util.ExpandTilde(target.Target)
					if err != nil {
						return fmt.Errorf("failed to expand tilde in target for project %s task %s: %w",
							projName, task.Name, err)
					}
					task.Targets[j] = target
				}
			}

			// Update the task in the Tasks slice
			proj.Tasks[i] = task
		}

		// Update the project in the Projects map
		cfg.Projects[projName] = proj
	}

	// Process user-level tasks
	for i, task := range cfg.User.Tasks {
		// Expand tilde in Task.Sources
		for j, src := range task.Sources {
			expanded, err := util.ExpandTilde(src)
			if err != nil {
				return fmt.Errorf("failed to expand tilde in source for user task %s: %w",
					task.Name, err)
			}
			task.Sources[j] = expanded
		}

		// Expand tilde in Target.Target for each Target
		for j, target := range task.Targets {
			if target.Target != "" {
				expanded, err := util.ExpandTilde(target.Target)
				if err != nil {
					return fmt.Errorf("failed to expand tilde in target for user task %s: %w",
						task.Name, err)
				}
				target.Target = expanded
				task.Targets[j] = target
			}
		}

		// Update the task in the User.Tasks slice
		cfg.User.Tasks[i] = task
	}

	return nil
}
