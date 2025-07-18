# Agent Definition (agent-def)
[![CI](https://github.com/uphy/agent-def/actions/workflows/ci.yml/badge.svg)](https://github.com/uphy/agent-def/actions/workflows/ci.yml)
[![Build](https://github.com/uphy/agent-def/actions/workflows/build.yml/badge.svg)](https://github.com/uphy/agent-def/actions/workflows/build.yml)
[![Release](https://img.shields.io/github/v/release/uphy/agent-def)](https://github.com/uphy/agent-def/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/uphy/agent-def)](https://goreportcard.com/report/github.com/uphy/agent-def)
[![codecov](https://codecov.io/gh/uphy/agent-def/branch/main/graph/badge.svg)](https://codecov.io/gh/uphy/agent-def)
[![Go Reference](https://pkg.go.dev/badge/github.com/uphy/agent-def.svg)](https://pkg.go.dev/github.com/uphy/agent-def)

Agent Definition (`agent-def`) is a tool for converting context and command definitions for various AI agents like Claude, Roo, and Cline.

## Installation

Install via Go:

```
go install github.com/uphy/agent-def/cmd/agent-def@latest
```

Or download a release from GitHub Releases.

## Usage

### Init Command

Initialize a new agent-def project with sample configuration and directory structure:

```
agent-def init
```

This command:
- Creates an `agent-def.yml` configuration file
- Sets up the `memories/` and `commands/` directories
- Adds sample template files to help you get started

Example:

```
mkdir my-new-project
cd my-new-project
agent-def init
```

### Build Command

Process the agent definitions according to the configuration:

```
agent-def build [flags]
```

Flags:
- `-c, --config string` Path to agent-def.yml file or directory containing it (default ".")
- `--dry-run` Show what would be generated without writing files
- `-f, --force` Force overwrite without prompting for confirmation

Example:

```
agent-def build --config ./configs/agent-def.yml my-app
```

The build command processes both project-specific and user-level tasks by default. If you specify project names as arguments, only those projects will be processed, but user-level tasks will still be processed.

## Configuration

agent-def uses a YAML configuration file (`agent-def.yml`) to define the sources and destinations for your agent definitions.

### Basic Structure

```yaml
configVersion: "1.0"
projects:
  project-name:
    # Project-specific configuration
user:
  # Global user-level configuration
```

### Task Types

agent-def supports two task types:

1. **Memory** (`type: memory`) - Defines context information for AI agents
2. **Command** (`type: command`) - Provides custom command definitions for AI agents

For detailed configuration options, output destinations, concatenation behavior, template syntax, and best practices, please refer to the [Configuration Reference](docs/config.md).

## Template Syntax

Agent Definition uses Go's `text/template` with `{{` and `}}` delimiters. Some key template functions include:

| Function | Description |
|----------|-------------|
| `file "path/to/file"` | Formats a file reference according to the target agent |
| `include "path/to/file"` | Includes content from another file with template processing |
| `reference "path/to/file"` | References another file without template processing |
| `mcp "serverName" "toolName" "arg1"` | Formats an MCP command for the target agent |
| `agent` | Returns the current target agent identifier |

See the [Configuration Reference](docs/config.md) for full details on template functions and path resolution.

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
brew install uphy/tap/agent-def
```