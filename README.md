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

### Init Command

Initialize a new agent-def project:

```
agent-def init
```

Example:

```
mkdir my-new-project
cd my-new-project
agent-def init
```

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

### init

Initialize a new agent-def project with sample configuration and directory structure:

```
agent-def init
```

This command:
- Creates an `agent-def.yml` configuration file
- Sets up the `context/` and `commands/` directories
- Adds sample template files to help you get started

No flags required.

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

## Release Process

agent-def uses [GoReleaser](https://goreleaser.com/) for automated releases and distribution. This process builds binaries for multiple platforms and architectures, creates GitHub releases, and automatically updates the Homebrew formula.

### Creating a New Release

To create a new release, follow these steps:

1. Ensure all changes are committed and pushed to the main branch
2. Create and push a new tag with the version number:

   ```bash
   # Tag with a new version (following Semantic Versioning)
   git tag -a v1.2.3 -m "Release v1.2.3"
   
   # Push the tag to the remote repository
   git push origin v1.2.3
   ```

3. The GitHub Actions workflow will automatically:
   - Build binaries for Linux, macOS, and Windows (amd64 and arm64)
   - Create a new GitHub release with the binaries and checksums
   - Update the Homebrew formula in the tap repository

### GitHub Actions Automation

The release process is automated through GitHub Actions, which uses the [GoReleaser Action](https://github.com/goreleaser/goreleaser-action). The workflow is triggered when a new version tag (matching `v*.*.*`) is pushed to the repository.

The `.github/workflows/release.yml` file defines the release workflow, which:
- Checks out the repository with full history
- Sets up Go
- Runs GoReleaser to build, package, and publish the release

### Required Secrets and Environment Variables

For the automated release process to work, you need to set up the following GitHub secret:

- `HOMEBREW_TAP_TOKEN`: A GitHub Personal Access Token with write access to the Homebrew tap repository. This token needs the following permissions:
  - `contents:write` - To push the formula to the tap repository
  - `workflows` - To trigger other workflows if necessary

To set up this secret:
1. Go to your GitHub repository settings
2. Navigate to "Secrets and variables" > "Actions"
3. Click on "New repository secret"
4. Add `HOMEBREW_TAP_TOKEN` with your PAT value

### Homebrew Formula

The release process automatically updates the Homebrew formula in the [homebrew-tap](https://github.com/uphy/homebrew-tap) repository, allowing users to install the latest version of agent-def using Homebrew:

```bash
# Add the tap (first time only)
brew tap uphy/tap

# Install agent-def
brew install agent-def
```