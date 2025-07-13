package processor

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/agent-def/internal/config"
	"github.com/user/agent-def/internal/util"
)

// Manager orchestrates task execution based on loaded configuration.
type Manager struct {
	cfg       *config.Config
	configDir string
	force     bool
}

// NewManager creates a new Manager by loading configuration from the given path.
func NewManager(cfgPath string) (*Manager, error) {
	cfg, configDir, err := config.LoadConfig(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	if cfg.ConfigVersion == "" {
		return nil, &util.ErrInvalidConfig{Reason: "configVersion is required"}
	}
	return &Manager{cfg: cfg, configDir: configDir, force: false}, nil
}

// Build executes the build pipeline for the specified projects or user scope.
func (m *Manager) Build(projects []string, userOnly, dryRun, force bool) error {
	m.force = force
	// Process project-level tasks unless userOnly is set
	if !userOnly {
		for name, proj := range m.cfg.Projects {
			if len(projects) > 0 && !contains(projects, name) {
				continue
			}
			// determine project root
			root := proj.Root
			if root == "" {
				// If no root is specified, use config directory
				root = m.configDir
			} else if !filepath.IsAbs(root) {
				// If root is relative, make it relative to config directory
				root = filepath.Join(m.configDir, root)
			}
			for _, task := range proj.Tasks {
				pipeline := NewPipeline(task, root, proj.Destinations, false, dryRun, m.force)
				if err := pipeline.Execute(); err != nil {
					return fmt.Errorf("project %s task execution failed: %w", name, err)
				}
			}
		}
	}

	// Process user-level tasks
	for _, task := range m.cfg.User.Tasks {
		// Use User.Home setting when specified, otherwise fallback to system home directory
		var home string
		if m.cfg.User.Home != "" {
			home = m.cfg.User.Home
		} else {
			var err error
			home, err = os.UserHomeDir()
			if err != nil {
				home = ""
			}
		}
		// Use the config directory for resolving user task sources
		root := m.configDir
		pipeline := NewPipeline(task, root, []string{home}, true, dryRun, m.force)
		if err := pipeline.Execute(); err != nil {
			return fmt.Errorf("user task execution failed: %w", err)
		}
	}

	return nil
}

// Validate checks the configuration for correctness.
func (m *Manager) Validate() error {
	// TODO: implement full schema and semantic validation
	fmt.Println("Manager.Validate invoked")
	return nil
}

// contains returns true if str is in list.
func contains(list []string, str string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}
