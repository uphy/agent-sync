# agent-sync

## Development Commands

This project uses [mise](https://mise.jdx.dev/) for task running:

```bash
# Run all tests (including integration tests)
mise run test

# Run only integration tests
mise run integration-test

# Run linter
mise run lint

# Run pre-commit checks
mise run pre-commit
```

## Project Architecture

agent-sync is a tool for converting context and command definitions for various AI agents like Claude, Roo, Cline, and Copilot. The architecture is organized as follows:

### Core Components

1. **Command Line Interface** (`/cmd/agent-sync/main.go`, `/internal/cli/`)
   - Uses Cobra for CLI implementation
   - Main commands: apply, init
   - Supports global flags for logging, output format, etc.

2. **Configuration System** (`/internal/config/`)
   - Defines YAML configuration structure in `config.go`
   - Handles loading configurations from `agent-sync.yml` files
   - Supports project and user-level configurations
   - Supports standard and simplified formats

3. **Processing Pipeline** (`/internal/processor/`)
   - `manager.go`: Manages the overall processing flow
   - `pipeline.go`: Handles processing tasks from configuration
   - `memory_processor.go`: Processes memory/context tasks
   - `command_processor.go`: Processes command tasks
   - `fs_adapter.go`: File system adapter for I/O operations

4. **Agent Implementations** (`/internal/agent/`)
   - `registry.go`: Central registry for supported agents
   - `agent.go`: Common interface for all agents
   - `claude.go`: Claude-specific implementation
   - `roo.go`: Roo-specific implementation
   - `cline.go`: Cline-specific implementation
   - `copilot.go`: Copilot-specific implementation

5. **Template System** (`/internal/template/`)
   - `engine.go`: Template processing engine
   - `functions.go`: Template helper functions (file, include, reference, mcp)
   - Supports Go's `text/template` with custom functions
   - Path resolution for template includes and references

6. **Logging System** (`/internal/log/`)
   - `config.go`: Logging configuration options
   - `logger.go`: Core logging implementation
   - `output.go`: Log output handling
   - `test.go`: Testing utilities for logs

7. **Model Definitions** (`/internal/model/`)
   - `command.go`: Represents command definitions with agent-specific properties
     - Roo: slug, name, description, roleDefinition, whenToUse, groups
     - Claude: description, allowed-tools
     - Copilot: mode, model, tools, description

8. **Utility Functions** (`/internal/util/`)
   - `file_glob.go`: Glob pattern matching for file selection
   - `path.go`: Path manipulation utilities
   - `errors.go`: Error handling utilities
   - `file.go`: File operation utilities

### Data Flow

1. Configuration is loaded from `agent-sync.yml`
2. Tasks are processed according to their type (memory or command)
3. Source files are read and processed through the template engine
4. Results are transformed for the target agent format
5. Output is written to the specified destinations

### File Structure

- `/cmd/agent-sync/` - Main application entry point
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
  - `/docs/cli.md` - Command line interface reference
  - `/docs/config.md` - Main configuration guide
  - `/docs/config-reference.md` - Detailed configuration reference
  - `/docs/examples.md` - Configuration examples and best practices
  - `/docs/glob-patterns.md` - Glob pattern syntax and behavior
  - `/docs/input-output.md` - Input handling and output paths
  - `/docs/logging.md` - Logging system documentation
  - `/docs/task-types.md` - Memory and command task types
  - `/docs/templates.md` - Template functions and path resolution
  - `/docs/troubleshooting.md` - Error handling and debugging

### Configuration Structure

The `agent-sync.yml` file defines:

1. Project-specific configurations:
   - Source paths for memories and commands
   - Output destinations for generated files
   - Target agent formats

2. User-level configurations:
   - Global context and commands
   - User-specific settings

Each task can be of type "memory" (context) or "command" (commands), with multiple sources and targets.

The configuration can be organized in two formats: standard or simplified.

For more details, refer to the documentation files in the `/docs` directory.
