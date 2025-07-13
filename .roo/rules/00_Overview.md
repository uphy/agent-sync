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

For more details, refer to the @/docs/config.md file.