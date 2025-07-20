package processor

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/uphy/agent-sync/internal/config"
	"github.com/uphy/agent-sync/internal/log"
	"github.com/uphy/agent-sync/internal/util"
	"go.uber.org/zap"
)

// Manager orchestrates task execution based on loaded configuration.
type Manager struct {
	cfg *config.Config
	// absConfigDir is the absolute directory path where the configuration file is located.
	absConfigDir string
	force        bool
	logger       *zap.Logger
	output       log.OutputWriter
}

// NewManager creates a new Manager by loading configuration from the given path.
func NewManager(cfgPath string, logger *zap.Logger, output log.OutputWriter) (*Manager, error) {
	if logger == nil {
		logger = zap.NewNop()
	}
	if !filepath.IsAbs(cfgPath) {
		return nil, fmt.Errorf("config path must be absolute: %s", cfgPath)
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
		cfg:          cfg,
		absConfigDir: configDir,
		force:        false,
		logger:       logger,
		output:       output,
	}, nil
}

// Apply executes the apply pipeline for the specified projects or user scope.
// If no projects are specified, all projects will be processed.
func (m *Manager) Apply(projects []string, dryRun, force bool) error {
	m.force = force

	m.logger.Info("Starting apply process",
		zap.Strings("projects", projects),
		zap.Bool("dryRun", dryRun),
		zap.Bool("force", force))
	// Process project-level tasks
	for name, proj := range m.cfg.Projects {
		// If `projects` is specified, only process those projects
		if len(projects) > 0 && !slices.Contains(projects, name) {
			continue
		}
		// Resolve absolute paths for input root and output directories
		absInputRoot := m.absConfigDir
		absOutputDirs := proj.OutputDirs
		for i, outputDir := range absOutputDirs {
			if filepath.IsAbs(outputDir) {
				continue
			}
			outputDir = filepath.Join(absInputRoot, outputDir)
			var err error
			outputDir, err = filepath.Abs(outputDir)
			if err != nil {
				return fmt.Errorf("failed to resolve absolute path for output directory %s: %w", outputDir, err)
			}
			absOutputDirs[i] = outputDir
		}

		// Always use config directory as the project root
		m.logger.Debug("Using config directory as project root", zap.String("project", name))
		m.logger.Info("Processing project",
			zap.String("name", name),
			zap.Int("taskCount", len(proj.Tasks)),
			zap.String("inputRoot", absInputRoot),
			zap.Strings("outputDirs", absOutputDirs))

		if m.output != nil {
			m.output.PrintProgress(fmt.Sprintf("Processing project: %s", name))
		}

		for i, task := range proj.Tasks {
			m.logger.Debug("Processing project task",
				zap.String("project", name),
				zap.Int("taskIndex", i),
				zap.String("taskName", task.Name),
				zap.String("taskType", string(task.Type)))

			pipeline, err := NewPipeline(task, absInputRoot, absOutputDirs, false, dryRun, m.force, m.logger, m.output)
			if err != nil {
				return fmt.Errorf("failed to create pipeline for project %s task %s: %w", name, task.Name, err)
			}
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

		absHome, err := filepath.Abs(home)
		if err != nil {
			return fmt.Errorf("failed to resolve user home directory: %w", err)
		}

		m.logger.Info("Processing user-level task",
			zap.String("taskName", task.Name),
			zap.String("taskType", string(task.Type)),
			zap.String("homeDir", absHome))

		if m.output != nil {
			m.output.PrintProgress(fmt.Sprintf("Processing user-level task: %s", task.Name))
		}

		// Use the config directory for resolving user task sources
		pipeline, err := NewPipeline(task, m.absConfigDir, []string{absHome}, true, dryRun, m.force, m.logger, m.output)
		if err != nil {
			return fmt.Errorf("failed to create pipeline for user task %s: %w", task.Name, err)
		}
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
