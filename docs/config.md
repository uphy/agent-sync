# agent-def Configuration Reference

## Overview

agent-def is a tool for managing and converting context and command definitions for AI agents such as Claude and Roo. This reference documents all configuration options available in the `agent-def.yml` file.

## Configuration File Location

By default, agent-def looks for the following files in the current directory:
- `agent-def.yml`
- `agent-def.yaml`

A specific file path can also be directly specified using command line options.

## Basic Structure

The configuration file is organized into the following main sections:

```yaml
configVersion: "1.0"
projects:
  project-name:
    # Project-specific configuration
user:
  # Global user-level configuration
```

## Configuration Reference

### Root Configuration

| Setting | Type | Required | Description |
|---------|------|----------|-------------|
| `configVersion` | String | Yes | Schema version of the configuration file (e.g., "1.0") |
| `projects` | Map | Yes | Map of named project configurations where keys are project identifiers (arbitrary names) and values are project configuration objects |
| `user` | Object | Yes | Global user-level configuration |

### Project Configuration

Each project can be configured with the following settings:

| Setting | Type | Required | Description |
|---------|------|----------|-------------|
| `root` | String | No | Base path for sources relative to the agent-def.yml location. Supports tilde (~) expansion for home directory |
| `destinations` | String Array | Yes | Output directories for generated files. Supports tilde (~) expansion for home directory. Multiple destinations can be specified to support scenarios like git worktree |
| `tasks` | Task Array | Yes | List of generation tasks for this project |

### UserConfig Settings

User-level configuration:

| Setting | Type | Required | Description |
|---------|------|----------|-------------|
| `tasks` | Task Array | Yes | List of generation tasks for user scope |
| `home` | String | No | Optional override for the user's home directory, primarily for debugging purposes |

### Task Configuration

Options for each task:

| Setting | Type | Required | Description |
|---------|------|----------|-------------|
| `name` | String | No | Optional identifier for the task. If not provided, a default name is automatically generated: for project tasks, "{project-name}-{type}" (e.g., "my-project-memory"); for user tasks, "user-{type}" (e.g., "user-command") |
| `type` | String | Yes | Type of task, either "command" or "memory" |
| `sources` | String Array | Yes | File or directory paths relative to Root. Supports tilde (~) expansion for home directory |
| `concat` | Boolean | No | When true, concatenates sources into one output |
| `targets` | Target Array | Yes | Defines the output agents and paths |

### Target Configuration

Options for each target:

| Setting | Type | Required | Description |
|---------|------|----------|-------------|
| `agent` | String | Yes | Target AI agent (e.g., "roo", "claude", "cline") |
| `target` | String | No | Optional custom output path. If not specified, the agent's default path is used. Supports tilde (~) expansion for home directory |

## Task Types

agent-def supports the following two task types:

### 1. Memory (`type: memory`)

Defines context information for AI agents. This serves as the agent's "memory" and can include general rules, architecture information, coding style, etc.

**Default output destinations:**

| Agent | Scope | Default Path | Description |
|-------|-------|-------------|-------------|
| Claude | User | `~/.claude/CLAUDE.md` | User's global Claude memory file |
| Claude | Project | `CLAUDE.md` | Project-specific Claude memory file (in project root) |
| Roo | User | `~/.roo/rules/{filename}.md` | User's global Roo memory file |
| Roo | Project | `.roo/rules/{filename}.md` | Project-specific Roo memory file |
| Cline | User | `~/Documents/Cline/Rules/{filename}.md` | User's global Cline memory file |
| Cline | Project | `.clinerules/{filename}.md` | Project-specific Cline memory file |

Note: For Roo and Cline, `{filename}` is derived from the source file name.

### 2. Command (`type: command`)

Provides custom command definitions for AI agents. This creates shortcuts for performing specific tasks.

**Default output destinations:**

| Agent | Scope | Default Path | Description |
|-------|-------|-------------|-------------|
| Claude | User | `~/.claude/commands/{filename}.md` | User's global Claude command file |
| Claude | Project | `.claude/commands/{filename}.md` | Project-specific Claude command file |
| Roo | User | `~/Library/Application Support/Code/User/globalStorage/rooveterinaryinc.roo-cline/settings/custom_modes.yaml` | VSCode global settings for Roo custom modes |
| Roo | Project | `.roomodes` | Project-specific Roo custom modes file (in project root) |
| Cline | User | `~/Documents/Cline/Workflows/{filename}.md` | User's global Cline workflow file |
| Cline | Project | `.clinerules/workflows/{filename}.md` | Project-specific Cline workflow file |

Note: For Claude, `{filename}` is derived from the source file name.

#### Agent-specific Command Frontmatter

Each agent may support specific frontmatter attributes for commands:

**Claude Code Frontmatter:**
- `claude.description`: Brief description of the command shown in the help menu
- `claude.allowed-tools`: List of tools the command is permitted to use

**Roo Frontmatter:**
- `roo.slug`: Command identifier
- `roo.name`: Display name of the command
- `roo.roleDefinition`: Description of what the command does
- `roo.whenToUse`: Description of when to use the command
- `roo.groups`: Permission groups for the command

### Default Concatenation Behavior

If not explicitly set via `concat` option, the following default concatenation behaviors apply:

| Agent | Task Type | Default Concatenation |
|-------|-----------|----------------------|
| Claude | memory | Yes (concat = true) |
| Claude | command | No (concat = false) |
| Roo | memory | No (concat = false) |
| Roo | command | Yes (concat = true) |
| Cline | memory | No (concat = false) |
| Cline | command | No (concat = false) |

## Path Resolution

All paths are resolved according to these rules:

1. Tilde (`~`) is expanded to the user's home directory
2. `sources` paths within a project are interpreted as relative paths from that project's `root`
3. If `root` is not specified, paths are relative to the `agent-def.yml` file location

### Template Path Resolution

Within templates, paths are resolved using special rules:
1. Paths starting with `/` are relative to the configuration file directory with the leading slash removed
2. Paths starting with `./` or `../` are relative to the including file's directory
3. Other paths without `./` or `../` prefix are relative to the configuration file directory
4. OS-absolute paths (like C:\ on Windows) are preserved as-is

## Configuration Example

```yaml
configVersion: "1.0"

# Project-specific configuration
projects:
  my-app:
    # Output directories (can specify multiple)
    destinations:
      - ~/projects/my-app
      
    # Generation tasks for this project
    tasks:
      # Project-specific context information
      - name: "Project Memory"
        type: memory
        sources:
          - ./projects/my-app/memories/
        targets:
          - agent: claude  # Outputs to CLAUDE.md
          - agent: roo     # Outputs to .roo/rules/context.md
          
      # Project-specific commands
      - name: "Project Commands"
        type: command
        sources:
          - ./projects/my-app/commands/
        targets:
          - agent: claude  # Outputs to .claude/commands/
          - agent: roo     # Outputs to .roomodes file

# User-level global configuration
user:
  tasks:
    # Global context information
    - name: "User Global Memory"
      type: memory
      sources:
        - ./user/memories/
      targets:
        - agent: claude
          target: "~/.claude_global_context.md"  # Custom output path
        - agent: roo
          
    # Global commands
    - name: "User Global Commands"
      type: command
      sources:
        - ./user/commands/
      targets:
        - agent: claude
        - agent: roo
```

## Template Functions

agent-def supports the following template functions in source files:

| Function | Description |
|----------|-------------|
| `file "path/to/file"` | Formats a file reference according to the target agent |
| `include "path/to/file"` | Includes content from another file with template processing |
| `reference "path/to/file"` | References another file without template processing |
| `mcp "agent" "command" "arg1" "arg2"` | Formats an MCP command for the target agent |
| `agent` | Returns the current target agent identifier |
| `ifAGENT "agent-name" "content"` | Conditionally includes content only for the specified agent |

## Best Practices

1. **Organize Project Structure**: Store related memories and commands in an organized directory structure
2. **Share Common Settings**: Reuse common tasks across multiple projects
3. **Use Explicit Paths**: Use tilde (~) rather than absolute paths, especially for custom output paths
4. **Meaningful Naming**: Give tasks meaningful names that clearly indicate their purpose
5. **Use Templates**: Leverage template functions to maintain DRY (Don't Repeat Yourself) principles
6. **Agent-Specific Content**: Use the `ifAGENT` function to include content specific to certain agents