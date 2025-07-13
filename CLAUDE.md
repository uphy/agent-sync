# agent-def

## Development Commands

### Building and Running

```bash
# Build the project
go build ./cmd/agent-def

# Install the binary to GOPATH/bin
go install ./cmd/agent-def

# Run the example configuration
cd example/agent-def
go run ../../cmd/agent-def/main.go build -f

# Clean up example output
rm -rf ../dest
```

### Testing

```bash
# Run all tests
go test ./...

# Run integration tests
go test -timeout 30s -run ^TestAgentDef$ ./internal/cli

# Update integration test expectations
AGENT_DEF_REPLACE=true go test -timeout 30s -run ^TestAgentDef$ ./internal/cli
```

### Development Tasks with mise

This project uses [mise](https://mise.jdx.dev/) for task running:

```bash
# Run the example configuration
mise run example

# Run integration tests
mise run integration-test

# Update integration test expectations
mise run integration-test-replace

# Install the binary
mise run install
```

## Project Architecture

Agent-def is a tool for converting context and command definitions for various AI agents like Claude and Roo. The architecture is organized as follows:

### Core Components

1. **Command Line Interface** (`/cmd/agent-def/main.go`, `/internal/cli/`)
   - Uses Cobra for CLI implementation
   - Main commands: build, validate, list, init

2. **Configuration System** (`/internal/config/`)
   - Defines YAML configuration structure in `config.go`
   - Handles loading configurations from `agent-def.yml` files
   - Supports project and user-level configurations

3. **Processing Pipeline** (`/internal/processor/`)
   - `manager.go`: Manages the overall processing flow
   - `pipeline.go`: Handles processing tasks from configuration

4. **Agent Implementations** (`/internal/agent/`)
   - `registry.go`: Central registry for supported agents
   - `agent.go`: Common interface for all agents
   - `claude.go`: Claude-specific implementation
   - `roo.go`: Roo-specific implementation

5. **Template System** (`/internal/template/`)
   - `engine.go`: Template processing engine
   - `functions.go`: Template helper functions (file, include, reference, mcp)
   - Supports Go's `text/template` with custom functions

6. **Model Definitions** (`/internal/model/`)
   - `context.go`: Represents memory context
   - `command.go`: Represents command definitions

### Data Flow

1. Configuration is loaded from `agent-def.yml`
2. Tasks are processed according to their type (memory or command)
3. Source files are read and processed through the template engine
4. Results are transformed for the target agent format
5. Output is written to the specified destinations

### File Structure

- `/cmd/agent-def/` - Main application entry point
- `/internal/` - Internal packages
  - `/agent/` - Agent-specific implementations
  - `/cli/` - CLI command implementations
  - `/config/` - Configuration structures and loading
  - `/frontmatter/` - YAML frontmatter processing
  - `/model/` - Data models
  - `/parser/` - Markdown parsing utilities
  - `/processor/` - Core processing logic
  - `/template/` - Template engine and functions
  - `/util/` - Common utilities
- `/example/` - Example configurations and templates
- `/docs/` - Documentation

### Configuration Structure

The `agent-def.yml` file defines:

1. Project-specific configurations:
   - Source paths for memories and commands
   - Output destinations for generated files
   - Target agent formats

2. User-level configurations:
   - Global context and commands
   - User-specific settings

Each task can be of type "memory" (context) or "command" (commands), with multiple sources and targets.

For more details, refer to the @docs/config.md file.

# agent-def Configuration Reference

## Overview

agent-def is a tool for managing and converting context and command definitions for AI agents such as Claude and Roo. This reference documents all configuration options available in the `agent-def.yml` file.

## Configuration File Location

By default, agent-def looks for the following files in the current directory:
- `agent-def.yml`
- `agent-def.yaml`

A specific file path can also be directly specified.

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
| `name` | String | No | Optional identifier for the task |
| `type` | String | Yes | Type of task, either "command" or "memory" |
| `sources` | String Array | Yes | File or directory paths relative to Root. Supports tilde (~) expansion for home directory |
| `concat` | Boolean | No | When true, concatenates sources into one output |
| `targets` | Target Array | Yes | Defines the output agents and paths |

### Target Configuration

Options for each target:

| Setting | Type | Required | Description |
|---------|------|----------|-------------|
| `agent` | String | Yes | Target AI agent (e.g., "roo", "claude") |
| `target` | String | No | Optional custom output path. If not specified, the agent's default path is used. Supports tilde (~) expansion for home directory |

## Task Types

agent-def supports the following two task types:

### 1. Memory (`type: memory`)

Defines context information for AI agents. This serves as the agent's "memory" and can include general rules, architecture information, coding style, etc.

**Example output destinations:**
- Claude: Default is `~/.claude/CLAUDE.md` or `CLAUDE.md` within the project
- Roo: Default is `~/.roo/rules/rules.md` or `.roo/rules/context.md` within the project

### 2. Command (`type: command`)

Provides custom command definitions for AI agents. This creates shortcuts for performing specific tasks.

**Example output destinations:**
- Claude: Default is `~/.claude/commands/` or `.claude/commands/` within the project
- Roo: Default is `~/Library/Application Support/Code/User/globalStorage/rooveterinaryinc.roo-cline/settings/custom_modes.yaml` or `.roomodes` within the project

## Path Resolution

All paths are resolved according to these rules:

1. Tilde (`~`) is expanded to the user's home directory
2. `sources` paths within a project are interpreted as relative paths from that project's `root`
3. If `root` is not specified, paths are relative to the `agent-def.yml` file location

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

## Best Practices

1. **Organize Project Structure**: Store related memories and commands in an organized directory structure
2. **Share Common Settings**: Reuse common tasks across multiple projects
3. **Use Explicit Paths**: Use tilde (~) rather than absolute paths, especially for custom output paths
4. **Meaningful Naming**: Give tasks meaningful names that clearly indicate their purpose