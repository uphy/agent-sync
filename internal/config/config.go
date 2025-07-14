package config

import (
	"fmt"

	"github.com/user/agent-def/internal/log"
)

// Config represents the root configuration for agent-def
type Config struct {
	// ConfigVersion defines the schema version of agent-def.yml
	ConfigVersion string `yaml:"configVersion"`
	// Projects holds named project configurations
	Projects map[string]Project `yaml:"projects"`
	// User holds global user-level configuration
	User UserConfig `yaml:"user"`
	// Logging holds logging configuration
	Logging *log.LogConfig `yaml:"logging,omitempty"`
}

// SetDefaultNames sets default names for all tasks across all projects and user config
func (c *Config) SetDefaultNames() {
	for projectName, project := range c.Projects {
		project.SetDefaultNames(projectName)
		c.Projects[projectName] = project // Update the map entry with the modified project
	}
	c.User.SetDefaultNames()

	// Set default logging configuration if not provided
	if c.Logging == nil {
		defaultLogConfig := log.DefaultLogConfig()
		c.Logging = &defaultLogConfig
	}
}

// Project represents a project-specific configuration block
type Project struct {
	// Root is the base path for sources relative to agent-def.yml location
	// Supports tilde (~) expansion for home directory.
	Root string `yaml:"root,omitempty"`
	// Destinations are the output directories for generated files
	// Supports tilde (~) expansion for home directory.
	Destinations []string `yaml:"destinations"`
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
	// Sources are file or directory paths relative to Root
	// Supports tilde (~) expansion for home directory.
	Sources []string `yaml:"sources"`
	// Concat indicates whether to concatenate sources into one output
	Concat *bool `yaml:"concat,omitempty"`
	// Targets define the output agents and paths
	Targets []Target `yaml:"targets"`
}

// SetDefaultName sets a default name for a task if one is not provided
func (t *Task) SetDefaultName(prefix string) {
	if t.Name == "" {
		t.Name = fmt.Sprintf("%s-%s", prefix, t.Type)
	}
}

// Target represents an output destination for a task
type Target struct {
	// Agent is the target AI agent (e.g., "roo", "claude")
	Agent string `yaml:"agent"`
	// Target is an optional custom output path
	// Supports tilde (~) expansion for home directory.
	Target string `yaml:"target,omitempty"`
}
