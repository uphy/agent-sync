# Implementation Design: Custom Agent-Def Configuration Path Support

## Overview

The current implementation only looks for the `agent-def.yml` file in the current directory. This design outlines changes to support specifying a custom path to the configuration file and ensuring that all relative paths in the configuration are interpreted relative to that file's directory.

## Changes Required

### 1. Add a configuration path flag to the build command (internal/cli/build.go)

```go
// New configuration path flag
var configPath string

// Add to flags
cmd.Flags().StringVarP(&configPath, "config", "c", ".", "Path to agent-def.yml file or directory containing it")

// Change the NewManager call from hardcoded "." to use the configPath
mgr, err := processor.NewManager(configPath)
```

### 2. Update config loader to track the configuration directory (internal/config/loader.go)

The `LoadConfig` function needs to return both the config and the directory where the config file was found:

```go
// Update signature to return the config directory
func LoadConfig(path string) (*Config, string, error) {
    cfgPath := path
    var cfgDir string
    
    info, err := os.Stat(path)
    if err != nil {
        return nil, "", fmt.Errorf("cannot stat path %q: %w", path, err)
    }

    if info.IsDir() {
        cfgDir = path
        // Search for agent-def.yml or agent-def.yaml in directory
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

    return &cfg, cfgDir, nil
}
```

### 3. Update Manager to store and use the config directory (internal/processor/manager.go)

```go
// Add configDir field to Manager struct
type Manager struct {
    cfg      *config.Config
    configDir string  // New field to store config directory
    force    bool
}

// Update NewManager to store the config directory
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

// Update Build method to use configDir for project root resolution
func (m *Manager) Build(projects []string, userOnly, watch, dryRun, force bool) error {
    m.force = force
    
    // Process project-level tasks unless userOnly is set
    if !userOnly {
        for name, proj := range m.cfg.Projects {
            if len(projects) > 0 && !contains(projects, name) {
                continue
            }
            
            // Determine project root relative to config directory
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

    // Process user-level tasks (user tasks part remains unchanged)
    // ...
}
```

## Implementation Steps

1. First, update the config loader to track and return the config file directory
2. Then, update the Manager to store and use the config directory
3. Finally, add the config flag to the build command
4. Update the command documentation to reflect the new flag option

## Testing

The implementation should be tested with various scenarios:
1. Default behavior (no flag) - should work as before
2. Specify directory containing agent-def.yml
3. Specify path directly to agent-def.yml file
4. Test with relative paths in the config that should be resolved from config file's directory

This change should be backward compatible with existing usage.