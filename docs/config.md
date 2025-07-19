# agent-def Configuration Reference

## Overview

agent-def is a tool for managing and converting context and command definitions for AI agents such as Claude, Roo, Cline, and Copilot. This reference documents all configuration options available in the `agent-def.yml` file.

## Documentation Structure

The configuration documentation has been organized into several focused documents:

- **[Configuration Reference](config-reference.md)** - Detailed configuration parameters and settings
- **[Input and Output Processing](input-output.md)** - Input handling, output paths, and file processing
- **[Task Types](task-types.md)** - Memory and command task types with agent-specific details
- **[Template System](templates.md)** - Template functions and path resolution
- **[Command Line Interface](cli.md)** - Command line options and usage
- **[Troubleshooting](troubleshooting.md)** - Error handling and debugging
- **[Examples and Best Practices](examples.md)** - Configuration examples and recommendations

## Configuration File Location

By default, agent-def looks for the following files in the current directory:
- `agent-def.yml`
- `agent-def.yaml`

A specific file path can also be directly specified using the `--config` command line option.

## Basic Structure

The configuration file can be organized in two formats: standard or simplified.

### Standard Format

```yaml
configVersion: "1.0"
projects:
  project-name:
    # Project-specific configuration
user:
  # Global user-level configuration
```

### Simplified Format (Single Project)

```yaml
configVersion: "1.0"
outputDirs:
  - ~/projects/my-app
tasks:
  # Project tasks defined directly at root level
```

## Navigation

- [Configuration Reference](config-reference.md)
- [Input and Output Processing](input-output.md)
- [Task Types](task-types.md)
- [Template System](templates.md)
- [Command Line Interface](cli.md)
- [Troubleshooting](troubleshooting.md)
- [Examples and Best Practices](examples.md)
- [Logging Guide](logging.md)