package config

import (
	"fmt"
)

// Config represents the root configuration for agent-def
type Config struct {
	// ConfigVersion defines the schema version of agent-def.yml
	ConfigVersion string `yaml:"configVersion"`
	// Projects holds named project configurations
	Projects map[string]Project `yaml:"projects"`
	// New fields for top-level project config
	OutputDirs []string `yaml:"outputDirs,omitempty"`
	Tasks      []Task   `yaml:"tasks,omitempty"`
	// User holds global user-level configuration
	User UserConfig `yaml:"user"`
}

// SetDefaultNames sets default names for all tasks across all projects and user config
func (c *Config) SetDefaultNames() {
	for projectName, project := range c.Projects {
		project.SetDefaultNames(projectName)
		c.Projects[projectName] = project // Update the map entry with the modified project
	}
	c.User.SetDefaultNames()

	// Set default names for top-level tasks
	for i := range c.Tasks {
		c.Tasks[i].SetDefaultName("default")
	}
}

// Project represents a project-specific configuration block
type Project struct {
	// OutputDirs are the output directories for generated files
	// Supports tilde (~) expansion for home directory.
	OutputDirs []string `yaml:"outputDirs"`
	// Tasks is the list of generation tasks for this project
	Tasks []Task `yaml:"tasks"`
}

// SetDefaultNames sets default names for all tasks in this project
func (p *Project) SetDefaultNames(projectName string) {
	for i := range p.Tasks {
		p.Tasks[i].SetDefaultName(projectName)
	}
}

// UserConfig represents global user-level configuration
type UserConfig struct {
	// Tasks is the list of generation tasks for user scope
	Tasks []Task `yaml:"tasks"`
	// Home is an optional override for the user's home directory, primarily for debugging purposes
	Home string `yaml:"home,omitempty"`
}

// SetDefaultNames sets default names for all tasks in user config
func (u *UserConfig) SetDefaultNames() {
	for i := range u.Tasks {
		u.Tasks[i].SetDefaultName("user")
	}
}

// Task represents a single generation task (command or memory)
type Task struct {
	// Name is an optional identifier for the task
	Name string `yaml:"name,omitempty"`
	// Type is either "command" or "memory"
	Type string `yaml:"type"`
	// Inputs are file or directory paths relative to config directory
	Inputs []string `yaml:"inputs"`
	// Outputs define the output agents and paths
	Outputs []Output `yaml:"outputs"`
}

// SetDefaultName sets a default name for a task if one is not provided
func (t *Task) SetDefaultName(prefix string) {
	if t.Name == "" {
		t.Name = fmt.Sprintf("%s-%s", prefix, t.Type)
	}
}

// Output represents an output destination for a task
type Output struct {
	// Agent is the target AI agent (e.g., "roo", "claude")
	Agent string `yaml:"agent"`
	// OutputPath is an optional custom output path
	// The path format determines concatenation behavior:
	// - If path ends with "/", it will be treated as a directory (non-concatenated outputs)
	// - If path doesn't end with "/", it will be treated as a file (concatenated outputs)
	OutputPath string `yaml:"outputPath,omitempty"`
}
