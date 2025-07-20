# Agent Sync (agent-sync)

[![CI](https://github.com/uphy/agent-sync/actions/workflows/ci.yml/badge.svg)](https://github.com/uphy/agent-sync/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/uphy/agent-sync)](https://github.com/uphy/agent-sync/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/uphy/agent-sync)](https://goreportcard.com/report/github.com/uphy/agent-sync)
[![codecov](https://codecov.io/gh/uphy/agent-sync/branch/main/graph/badge.svg)](https://codecov.io/gh/uphy/agent-sync)
[![Go Reference](https://pkg.go.dev/badge/github.com/uphy/agent-sync.svg)](https://pkg.go.dev/github.com/uphy/agent-sync)

## What is Agent Sync?

Agent Sync (`agent-sync`) is a tool that streamlines working with multiple AI assistants by letting you write your context and commands once and automatically converting them to formats compatible with different AI agents.

### Key Features:

- **Write Once, Deploy Everywhere**: Maintain a single source of truth for AI instructions instead of separate files for each assistant
- **Multi-Agent Support**: Convert context (memory) and commands to formats compatible with Claude, Roo, Cline, Copilot, and more
- **Template System**: Use powerful templating to dynamically generate content with agent-specific formatting
- **Project and User Separation**: Maintain both project-specific and global user-level configurations
- **Flexible Output**: Control how content is organized (concatenated into single files or split into multiple files)

### How It Works:

1. You write context (memory) and commands in standard markdown files
2. You configure how these files should be processed in `agent-sync.yml`
3. You run `agent-sync apply` to generate agent-specific files in the correct formats and locations
4. Your AI assistants now have the same knowledge and capabilities across platforms

## Installation

### Homebrew Installation (macOS)

You can install `agent-sync` using Homebrew:

```bash
brew install uphy/tap/agent-sync
```

### Go Installation

Install via Go:

```bash
go install github.com/uphy/agent-sync/cmd/agent-sync@latest
```

### Binary Downloads

Or download a release from [GitHub Releases](https://github.com/uphy/agent-sync/releases).

## Getting Started

### Init Command

Initialize a new agent-sync project with sample configuration and directory structure:

```bash
mkdir my-ai-instructions
cd my-ai-instructions
agent-sync init
```

This command:
- Creates an `agent-sync.yml` configuration file
- Sets up the `memories/` and `commands/` directories
- Adds sample template files to help you get started

### Apply Command

The `apply` command processes your input files and generates agent-specific output files based on your configuration:

```bash
agent-sync apply [flags]
```

What it does:
1. Reads your input files (memories and commands)
2. Processes any templates in those files
3. Converts the content to formats compatible with each target agent
4. Writes the output files to their specified output paths

**Important Flags:**
- `-c, --config string`: Specify a custom path to your configuration file
- `--dry-run`: Preview what would be generated without actually writing any files (useful for testing)
- `-f, --force`: Skip confirmation prompts when overwriting existing files
- `--verbose`: Show detailed output about what's happening

## Common Usage Patterns

### Example 1: Project with Multiple Agents

For a project where you want to use the same instructions with Claude, Roo, and Cline:

```yaml
# agent-sync.yml
configVersion: "1.0"
projects:
  my-project:
    outputDirs:
      - ~/projects/my-app
    tasks:
      - type: memory
        inputs:
          - ./memories/*.md
        outputs:
          - agent: claude  # Will generate CLAUDE.md
          - agent: roo     # Will generate files in .roo/rules/
          - agent: cline   # Will generate files in .clinerules/
```

Run: `agent-sync apply` to generate all formats.

### Example 2: Shared Context and Project-Specific Commands

For maintaining shared context information but project-specific commands:

```yaml
# agent-sync.yml
configVersion: "1.0"
user:
  tasks:
    - type: memory
      name: "Global Rules"
      inputs:
        - ./global-rules/*.md
      outputs:
        - agent: claude
        - agent: roo
projects:
  project-a:
    outputDirs:
      - ~/projects/project-a
    tasks:
      - type: command
        inputs:
          - ./project-a/commands/*.md
        outputs:
          - agent: claude
          - agent: roo
```

Run: `agent-sync apply` to generate both global context and project-specific commands.

### Example 3: Testing Changes with Dry Run

Before applying changes, you can preview what would be generated:

```bash
agent-sync apply --dry-run --verbose
```

This shows what files would be generated without actually writing them.

## Configuration

agent-sync uses a YAML configuration file (`agent-sync.yml`) to define the inputs and outputs for your agent definitions.

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

agent-sync supports two task types:

1. **Memory** (`type: memory`) - Defines context information for AI agents (rules, architecture, guidelines, etc.)
2. **Command** (`type: command`) - Provides custom command definitions for AI agents (specialized modes, workflows, etc.)

For detailed configuration options, output destinations, concatenation behavior, template syntax, and best practices, please refer to the [Configuration Documentation](docs/config.md) which is organized into several focused guides.

## Template Syntax

Agent Sync uses Go's `text/template` with `{{` and `}}` delimiters. Some key template functions include:

| Function | Description |
|----------|-------------|
| `file "path/to/file"` | Formats a file reference according to the target agent |
| `include "path/to/file"` | Includes content from another file with template processing |
| `reference "path/to/file"` | References another file without template processing |
| `mcp "serverName" "toolName" "arg1"` | Formats an MCP command for the target agent |
| `agent` | Returns the current target agent identifier |

See the [Template System Documentation](docs/templates.md) for full details on template functions and path resolution.

## Release Process

agent-sync uses [GoReleaser](https://goreleaser.com/) for automated releases and distribution. This process builds binaries for multiple platforms and architectures, creates GitHub releases, and automatically updates the Homebrew formula.

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

The release process automatically updates the Homebrew formula in the [homebrew-tap](https://github.com/uphy/homebrew-tap) repository, allowing users to install the latest version of agent-sync using Homebrew:

```bash
brew install uphy/tap/agent-sync
```

## License

agent-sync is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
