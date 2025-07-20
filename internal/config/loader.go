package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/uphy/agent-def/internal/util"
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

	// Check for mixed format (both top-level fields and projects section)
	if (len(cfg.OutputDirs) > 0 || len(cfg.Tasks) > 0) && len(cfg.Projects) > 0 {
		return nil, "", fmt.Errorf("configuration cannot mix top-level outputDirs/tasks with projects section")
	}

	// Convert top-level settings to a "default" project
	if len(cfg.OutputDirs) > 0 || len(cfg.Tasks) > 0 {
		if cfg.Projects == nil {
			cfg.Projects = make(map[string]Project)
		}

		// Create "default" project
		defaultProject := Project{
			OutputDirs: cfg.OutputDirs,
			Tasks:      cfg.Tasks,
		}

		cfg.Projects["default"] = defaultProject

		// Clear top-level fields after conversion
		cfg.OutputDirs = nil
		cfg.Tasks = nil
	}

	// Normalize all paths in the config by expanding tildes (~) to home directory
	if err := expandTildeInConfig(&cfg); err != nil {
		return nil, "", fmt.Errorf("failed to expand tildes in paths: %w", err)
	}

	// Set default names for tasks that don't have names specified
	cfg.SetDefaultNames()

	return &cfg, cfgDir, nil
}

// expandTildeInConfig normalizes all paths in the config by expanding tildes (~)
// to the user's home directory.
func expandTildeInConfig(cfg *Config) error {
	// Process all projects
	for projName, proj := range cfg.Projects {
		var err error

		// Expand tilde in Project.OutputDirs
		for i, dest := range proj.OutputDirs {
			expanded, err := util.ExpandTilde(dest)
			if err != nil {
				return fmt.Errorf("failed to expand tilde in project %s output directory: %w", projName, err)
			}
			proj.OutputDirs[i] = expanded
		}

		// Process Task.Inputs and Output.OutputPath in each Task
		for i, task := range proj.Tasks {
			// Expand tilde in Task.Inputs
			for j, src := range task.Inputs {
				expanded, err := util.ExpandTilde(src)
				if err != nil {
					return fmt.Errorf("failed to expand tilde in input for project %s task %s: %w",
						projName, task.Name, err)
				}
				task.Inputs[j] = expanded
			}

			// Expand tilde in Output.OutputPath for each Output
			for j, output := range task.Outputs {
				if output.OutputPath != "" {
					output.OutputPath, err = util.ExpandTilde(output.OutputPath)
					if err != nil {
						return fmt.Errorf("failed to expand tilde in output path for project %s task %s: %w",
							projName, task.Name, err)
					}
					task.Outputs[j] = output
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
		// Expand tilde in Task.Inputs
		for j, src := range task.Inputs {
			expanded, err := util.ExpandTilde(src)
			if err != nil {
				return fmt.Errorf("failed to expand tilde in input for user task %s: %w",
					task.Name, err)
			}
			task.Inputs[j] = expanded
		}

		// Expand tilde in Output.OutputPath for each Output
		for j, output := range task.Outputs {
			if output.OutputPath != "" {
				expanded, err := util.ExpandTilde(output.OutputPath)
				if err != nil {
					return fmt.Errorf("failed to expand tilde in output path for user task %s: %w",
						task.Name, err)
				}
				output.OutputPath = expanded
				task.Outputs[j] = output
			}
		}

		// Update the task in the User.Tasks slice
		cfg.User.Tasks[i] = task
	}

	return nil
}
