# Configuration Reference

This document details all configuration parameters available in the `agent-def.yml` file, their types, requirements, and descriptions.

## Root Configuration

| Setting | Type | Required | Description |
|---------|------|----------|-------------|
| `configVersion` | String | Yes | Schema version of the configuration file (e.g., "1.0") |
| `projects` | Map | No* | Map of named project configurations where keys are project identifiers (arbitrary names) and values are project configuration objects. *Required in standard format, omitted in simplified format. |
| `user` | Object | No* | Global user-level configuration. *Required in standard format. |
| `outputDirs` | String Array | No* | Output directories where generated files will be placed. *Only used in simplified format. |
| `tasks` | Task Array | No* | List of generation tasks. *Only used in simplified format. |

## Project Configuration

Each project can be configured with the following settings:

| Setting | Type | Required | Description |
|---------|------|----------|-------------|
| `outputDirs` | String Array | Yes | Output directories where generated files will be placed. Supports tilde (~) expansion for home directory. Multiple directories can be specified to support scenarios like git worktree |
| `tasks` | Task Array | Yes | List of generation tasks for this project |

## UserConfig Settings

User-level configuration:

| Setting | Type | Required | Description |
|---------|------|----------|-------------|
| `tasks` | Task Array | Yes | List of generation tasks for user scope |
| `home` | String | No | Optional override for the user's home directory, primarily for debugging purposes |

## Task Configuration

Options for each task:

| Setting | Type | Required | Description |
|---------|------|----------|-------------|
| `name` | String | No | Optional identifier for the task. If not provided, a default name is automatically generated: for project tasks, "{project-name}-{type}" (e.g., "my-project-memory"); for user tasks, "user-{type}" (e.g., "user-command") |
| `type` | String | Yes | Type of task, either "command" or "memory" |
| `inputs` | String Array | Yes | File or directory paths relative to config directory. Supports glob patterns with exclusions |
| `outputs` | Output Array | Yes | Defines the output agents and their paths |

## Output Configuration

Options for each output:

| Setting | Type | Required | Description |
|---------|------|----------|-------------|
| `agent` | String | Yes | Target AI agent (e.g., "roo", "claude", "cline", "copilot") |
| `outputPath` | String | No | Optional custom output path. If not specified, the agent's default path is used. The path format determines concatenation behavior: paths ending with "/" are treated as directories (non-concatenated outputs), while paths without a trailing "/" are treated as files (concatenated outputs) |

## Navigation

- [Main Configuration Guide](config.md)
- [Input and Output Processing](input-output.md)
- [Task Types](task-types.md)
- [Template System](templates.md)
- [Command Line Interface](cli.md)
- [Troubleshooting](troubleshooting.md)
- [Examples and Best Practices](examples.md)