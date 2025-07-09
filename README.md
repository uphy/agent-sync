# Agent Definition (agent-def)
[![CI](https://github.com/user/agent-def/actions/workflows/ci.yml/badge.svg)](https://github.com/user/agent-def/actions/workflows/ci.yml)  
[![Build](https://github.com/user/agent-def/actions/workflows/build.yml/badge.svg)](https://github.com/user/agent-def/actions/workflows/build.yml)  
[![Release](https://img.shields.io/github/v/release/user/agent-def)](https://github.com/user/agent-def/releases/latest)  
[![Go Report Card](https://goreportcard.com/badge/github.com/user/agent-def)](https://goreportcard.com/report/github.com/user/agent-def)  
[![codecov](https://codecov.io/gh/user/agent-def/branch/main/graph/badge.svg)](https://codecov.io/gh/user/agent-def)  
[![Go Reference](https://pkg.go.dev/badge/github.com/user/agent-def.svg)](https://pkg.go.dev/github.com/user/agent-def)  

Agent Definition (`agent-def`) is a tool for converting context and command definitions for various AI agents like Claude and Roo.

## Installation

Install via Go:

```
go install github.com/user/agent-def/cmd/agent-def@latest
```

Or download a release from GitHub Releases.

## Usage

### Memory Command

Process a memory context file:

```
agent-def memory <file> [flags]
```

Example:

```
agent-def memory examples/memory_context.md --agent claude
```

### Command Command

Process command definition(s):

```
agent-def command <file_or_directory> [flags]
```

Example:

```
agent-def command examples/command_definition.md --agent roo --base .
```

### List Command

List supported agents:

```
agent-def list [flags]
```

Example:

```
agent-def list --verbose
```

## Global Flags

- `-o, --output string` Output format (`json`, `yaml`, `text`)
- `-v, --verbose` Enable verbose output

## Commands

### memory

Flags:

- `-a, --agent string` Target agent type (default `claude`)
- `-b, --base string` Base path for resolving relative paths (defaults to input file directory)

### command

Flags:

- `-a, --agent string` Target agent type (default `claude`)
- `-b, --base string` Base path for resolving relative paths (defaults to input file/directory)

### list

Flags:

- `-v, --verbose` Show detailed capabilities for each agent

## Template Syntax

Agent Definition uses Go's `text/template` with `{{` and `}}` delimiters. The following helper functions are available:

- **file**: Insert an agent-specific file reference.

  ```
  {{file "path/to/file"}}
  ```

- **include**: Include and process another markdown file.

  ```
  {{include "path/to/file.md"}}
  ```

- **reference**: Include a reference to another file and append its content in the References section.

  ```
  {{reference "path/to/file.md"}}
  ```

- **mcp**: Generate an agent-specific MCP command invocation.

  ```
  {{mcp "serverName" "toolName" "arg1" "arg2"}}
  ```

## Examples

See the `examples/` directory for sample memory context and command definition files, and example outputs for different agents.