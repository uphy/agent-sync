package processor

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/agent-def/internal/config"
	"github.com/user/agent-def/internal/log"
	"github.com/user/agent-def/internal/util"
	"go.uber.org/zap"
)

// Manager orchestrates task execution based on loaded configuration.
type Manager struct {
	cfg       *config.Config
	configDir string
	force     bool
	logger    *zap.Logger
	output    log.OutputWriter
}

// NewManager creates a new Manager by loading configuration from the given path.
func NewManager(cfgPath string, logger *zap.Logger, output log.OutputWriter) (*Manager, error) {
	if logger == nil {
		logger = zap.NewNop()
	}

	logger.Debug("Loading configuration", zap.String("path", cfgPath))
	cfg, configDir, err := config.LoadConfig(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	if cfg.ConfigVersion == "" {
		logger.Error("Invalid configuration", zap.String("reason", "configVersion is required"))
		return nil, &util.ErrInvalidConfig{Reason: "configVersion is required"}
	}

	logger.Info("Configuration loaded successfully",
		zap.String("configVersion", cfg.ConfigVersion),
		zap.Int("projectCount", len(cfg.Projects)),
		zap.Int("userTaskCount", len(cfg.User.Tasks)))

	return &Manager{
		cfg:       cfg,
		configDir: configDir,
		force:     false,
		logger:    logger,
		output:    output,
	}, nil
}

// Build executes the build pipeline for the specified projects or user scope.
func (m *Manager) Build(projects []string, dryRun, force bool) error {
	m.force = force

	m.logger.Info("Starting build process",
		zap.Strings("projects", projects),
		zap.Bool("dryRun", dryRun),
		zap.Bool("force", force))
	// Process project-level tasks
	for name, proj := range m.cfg.Projects {
		if len(projects) > 0 && !contains(projects, name) {
			continue
		}
		// determine project root
		root := proj.Root
		if root == "" {
			// If no root is specified, use config directory
			root = m.configDir
			m.logger.Debug("Using config directory as project root", zap.String("project", name))
		} else if !filepath.IsAbs(root) {
			// If root is relative, make it relative to config directory
			root = filepath.Join(m.configDir, root)
			m.logger.Debug("Resolved relative project root",
				zap.String("project", name),
				zap.String("root", root))
		}
		m.logger.Info("Processing project",
			zap.String("name", name),
			zap.Int("taskCount", len(proj.Tasks)),
			zap.Strings("outputDirs", proj.OutputDirs))

		if m.output != nil {
			m.output.PrintProgress(fmt.Sprintf("Processing project: %s", name))
		}

		for i, task := range proj.Tasks {
			m.logger.Debug("Processing project task",
				zap.String("project", name),
				zap.Int("taskIndex", i),
				zap.String("taskName", task.Name),
				zap.String("taskType", string(task.Type)))

			pipeline := NewPipeline(task, root, proj.OutputDirs, false, dryRun, m.force, m.logger, m.output)
			if err := pipeline.Execute(); err != nil {
				m.logger.Error("Project task execution failed",
					zap.String("project", name),
					zap.Error(err))

				if m.output != nil {
					m.output.PrintError(err)
				}

				return fmt.Errorf("project %s task execution failed: %w", name, err)
			}
		}

		if m.output != nil {
			m.output.PrintSuccess(fmt.Sprintf("Project %s processed successfully", name))
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

		m.logger.Info("Processing user-level task",
			zap.String("taskName", task.Name),
			zap.String("taskType", string(task.Type)),
			zap.String("homeDir", home))

		if m.output != nil {
			m.output.PrintProgress(fmt.Sprintf("Processing user-level task: %s", task.Name))
		}

		// Use the config directory for resolving user task sources
		root := m.configDir
		pipeline := NewPipeline(task, root, []string{home}, true, dryRun, m.force, m.logger, m.output)
		if err := pipeline.Execute(); err != nil {
			m.logger.Error("User task execution failed", zap.Error(err))

			if m.output != nil {
				m.output.PrintError(err)
			}

			return fmt.Errorf("user task execution failed: %w", err)
		}

		if m.output != nil {
			m.output.PrintSuccess(fmt.Sprintf("User task %s processed successfully", task.Name))
		}
	}

	return nil
}

// Validate checks the configuration for correctness.
func (m *Manager) Validate() error {
	m.logger.Info("Validating configuration")

	// TODO: implement full schema and semantic validation
	if m.output != nil {
		m.output.Print("Configuration validation not fully implemented yet")
	}

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
