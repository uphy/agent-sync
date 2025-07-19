# agent-def Configuration Reference

## Overview

agent-def is a tool for managing and converting context and command definitions for AI agents such as Claude, Roo, Cline, and Copilot. This reference documents all configuration options available in the `agent-def.yml` file.

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

## Configuration Reference

### Root Configuration

| Setting | Type | Required | Description |
|---------|------|----------|-------------|
| `configVersion` | String | Yes | Schema version of the configuration file (e.g., "1.0") |
| `projects` | Map | No* | Map of named project configurations where keys are project identifiers (arbitrary names) and values are project configuration objects. *Required in standard format, omitted in simplified format. |
| `user` | Object | No* | Global user-level configuration. *Required in standard format. |
| `outputDirs` | String Array | No* | Output directories where generated files will be placed. *Only used in simplified format. |
| `tasks` | Task Array | No* | List of generation tasks. *Only used in simplified format. |

### Project Configuration

Each project can be configured with the following settings:

| Setting | Type | Required | Description |
|---------|------|----------|-------------|
| `root` | String | No | Base path for inputs relative to the agent-def.yml location. Supports tilde (~) expansion for home directory |
| `outputDirs` | String Array | Yes | Output directories where generated files will be placed. Supports tilde (~) expansion for home directory. Multiple directories can be specified to support scenarios like git worktree |
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
| `inputs` | String Array | Yes | File or directory paths relative to Root. Supports tilde (~) expansion for home directory and glob patterns with exclusions |
| `outputs` | Output Array | Yes | Defines the output agents and their paths |

### Output Configuration

Options for each output:

| Setting | Type | Required | Description |
|---------|------|----------|-------------|
| `agent` | String | Yes | Target AI agent (e.g., "roo", "claude", "cline", "copilot") |
| `outputPath` | String | No | Optional custom output path. If not specified, the agent's default path is used. Supports tilde (~) expansion for home directory. The path format determines concatenation behavior: paths ending with "/" are treated as directories (non-concatenated outputs), while paths without a trailing "/" are treated as files (concatenated outputs) |

### Glob Pattern Support in Inputs

The `inputs` field supports glob patterns with doublestar support for recursive matching and exclusions:

| Pattern | Description |
|---------|-------------|
| `*.md` | All Markdown files in the current directory |
| `**/*.md` | All Markdown files in any subdirectory (recursive) |
| `!*_test.md` | Exclude all test Markdown files |

#### Examples

```yaml
inputs:
  # Include all Markdown files in commands directory and subdirectories
  - "commands/**/*.md"
  # Exclude test files
  - "!commands/**/*_test.md"
  # Exclude specific temporary file
  - "!commands/temp.md"
```

#### Pattern Order and Precedence

1. All include patterns (without `!` prefix) are processed first
2. Then exclude patterns (with `!` prefix) are applied to filter the results
3. The final list of files is sorted alphabetically by path

> **Note**: The final list of files will be sorted alphabetically by path. If you need files to be processed in a specific order, list them individually without glob patterns.

## Relationship Between outputDirs and outputs

The `outputDirs` setting (at the project level or root level in simplified format) defines the base directories where files will be generated. The `outputs` setting within each task defines which agents to generate files for and optionally the specific paths relative to the base directories.

- When multiple `outputDirs` are specified, the same output files will be generated in each directory.
- Each `outputs` entry defines a specific agent target and (optionally) a custom output path.
- If no custom `outputPath` is specified, the agent's default path is used (see "Default Output Path Behavior" section).

Example:
```yaml
outputDirs:
  - ~/projects/my-app
  - ~/projects/my-app-worktree
tasks:
  - type: memory
    outputs:
      - agent: roo
        # Will be generated in both:
        # ~/projects/my-app/.roo/rules/
        # ~/projects/my-app-worktree/.roo/rules/
```

## Output Path Interpretation

The format of the `outputPath` determines how files are processed:

- Paths ending with `/` are treated as **directories** where individual files will be placed
- Paths not ending with `/` are treated as **files** where all input content will be concatenated

For example:
- `.roo/rules/` - Treated as a directory, each source file becomes a separate output file
- `.roo/rules/combined.md` - Treated as a file, all source files are concatenated into one file

If `outputPath` is not specified, the agent's default path is used (see "Default Output Path Behavior" sections below).

## Default Output Path Behavior

When no `outputPath` is specified, the following default paths are used. The path format (with or without trailing slash) determines whether files are concatenated or kept separate:

| Agent | Task Type | Default Path | Concatenation Behavior |
|-------|-----------|-------------|------------------------|
| Claude | memory | `CLAUDE.md` or `~/.claude/CLAUDE.md` | Concatenated (file path) |
| Claude | command | `.claude/commands/` | Non-concatenated (directory path) |
| Roo | memory | `.roo/rules/` | Non-concatenated (directory path) |
| Roo | command | `.roomodes` or global settings path | Concatenated (file path) |
| Cline | memory | `.clinerules/` or `~/Documents/Cline/Rules/` | Non-concatenated (directory path) |
| Cline | command | `.clinerules/workflows/` or `~/Documents/Cline/Workflows/` | Non-concatenated (directory path) |
| Copilot | memory | `.github/copilot-instructions.md` or `~/.vscode/copilot-instructions.md` | Concatenated (file path) |
| Copilot | command | `.github/prompts/` or `~/.vscode/prompts/` | Non-concatenated (directory path) |

## Task Processing Workflow

When the `agent-def build` command is executed, the following workflow occurs:

1. **Configuration Loading**: The configuration file is loaded and validated.
2. **Task Selection**: If specific projects are provided as command arguments, only those projects' tasks are processed. Otherwise, all projects and user tasks are processed.
3. **Input Processing**:
   - For each task, the input files are resolved based on glob patterns.
   - Each input file is processed through the template engine if applicable.
4. **Output Generation**:
   - For each output agent specified in the task, the content is formatted according to the agent's requirements.
   - The formatted content is written to the specified output paths in each output directory.
5. **Verification**: Optionally, with the `--dry-run` flag, the process will show what would be generated without actually writing files.

## Task Types

agent-def supports the following two task types:

### 1. Memory (`type: memory`)

Defines context information for AI agents. This serves as the agent's "memory" and can include general rules, architecture information, coding style, etc.

**Default output locations:**

| Agent | Scope | Default Path | Description |
|-------|-------|-------------|-------------|
| Claude | User | `~/.claude/CLAUDE.md` | User's global Claude memory file |
| Claude | Project | `CLAUDE.md` | Project-specific Claude memory file (in project root) |
| Roo | User | `~/.roo/rules/{filename}.md` | User's global Roo memory file |
| Roo | Project | `.roo/rules/{filename}.md` | Project-specific Roo memory file |
| Cline | User | `~/Documents/Cline/Rules/{filename}.md` | User's global Cline memory file |
| Cline | Project | `.clinerules/{filename}.md` | Project-specific Cline memory file |
| Copilot | User | `~/.vscode/copilot-instructions.md` | User's global Copilot instructions |
| Copilot | Project | `.github/copilot-instructions.md` | Project-specific Copilot instructions |

Note: For Roo, Cline, and similar agents, `{filename}` is derived from the input file's basename (the filename without its directory path). For example, an input file named `my-project/memories/coding-rules.md` would result in an output file named `coding-rules.md` in the appropriate output directory.

### 2. Command (`type: command`)

Provides custom command definitions for AI agents. This creates shortcuts for performing specific tasks.

**Default output locations:**

| Agent | Scope | Default Path | Description |
|-------|-------|-------------|-------------|
| Claude | User | `~/.claude/commands/{filename}.md` | User's global Claude command file |
| Claude | Project | `.claude/commands/{filename}.md` | Project-specific Claude command file |
| Roo | User | `~/Library/Application Support/Code/User/globalStorage/rooveterinaryinc.roo-cline/settings/custom_modes.yaml` | VSCode global settings for Roo custom modes (all commands are concatenated into this single file) |
| Roo | Project | `.roomodes` | Project-specific Roo custom modes file (all commands are concatenated into this single file) |
| Cline | User | `~/Documents/Cline/Workflows/{filename}.md` | User's global Cline workflow file |
| Cline | Project | `.clinerules/workflows/{filename}.md` | Project-specific Cline workflow file |
| Copilot | User | `~/.vscode/prompts/{filename}.prompt.md` | User's global Copilot prompt file |
| Copilot | Project | `.github/prompts/{filename}.prompt.md` | Project-specific Copilot prompt file |

Note: For Claude, Cline, Copilot, and similar agents, `{filename}` is derived from the input file's basename (the filename without its directory path). For example, an input file named `my-project/commands/deploy.md` would result in an output file named `deploy.md` in the appropriate output directory.

#### Agent-specific Command Frontmatter

Each agent may support specific frontmatter attributes for commands:

**Claude Command Frontmatter:**
- `claude.description`: Brief description of the command shown in the help menu
- `claude.allowed-tools`: List of tools the command is permitted to use

**Roo Frontmatter:**
- `roo.slug`: Command identifier
- `roo.name`: Display name of the command
- `roo.roleDefinition`: Description of what the command does
- `roo.whenToUse`: Description of when to use the command
- `roo.groups`: Permission groups for the command

**Cline Frontmatter:**
- Cline workflows use standard markdown without special frontmatter requirements

**Copilot Command Frontmatter:**
- `mode`: The operational mode for Copilot (e.g., "chat", "inline")
- `model`: The specific model to use
- `tools`: List of available tools
- `description`: Brief description of the prompt's purpose

## Path Resolution

All paths are resolved according to these rules:

1. Tilde (`~`) is expanded to the user's home directory
2. `inputs` paths within a project are interpreted as relative paths from that project's `root`
3. If `root` is not specified, paths are relative to the `agent-def.yml` file location

### Template Path Resolution

Within templates, paths are resolved using special rules:

1. Paths starting with `/` are relative to the configuration file directory with the leading slash removed
   - Example: `/shared/template.md` → `<config-dir>/shared/template.md`
   
2. Paths starting with `./` or `../` are relative to the including file's directory
   - Example from file `/templates/main.md`: `./partial.md` → `/templates/partial.md`
   - Example from file `/templates/main.md`: `../shared/partial.md` → `/shared/partial.md`
   
3. Other paths without `./` or `../` prefix are relative to the configuration file directory
   - Example: `shared/template.md` → `<config-dir>/shared/template.md`
   
4. OS-absolute paths (like C:\ on Windows) are preserved as-is
   - Example: `C:\templates\file.md` remains unchanged

Examples assuming config file is at `/project/agent-def.yml`:

| Template Path | Including File | Resolved Path |
|---------------|----------------|---------------|
| `/shared/template.md` | Any file | `/project/shared/template.md` |
| `./partial.md` | `/project/templates/main.md` | `/project/templates/partial.md` |
| `../shared/common.md` | `/project/templates/main.md` | `/project/shared/common.md` |
| `utils/helper.md` | Any file | `/project/utils/helper.md` |
| `C:\absolute\path.md` | Any file (Windows) | `C:\absolute\path.md` |

## Simplified Configuration Format

agent-def supports a simplified configuration format designed for single project use cases. This format allows you to place `outputDirs` and `tasks` directly at the top level, omitting the `projects` section:

```yaml
configVersion: "1.0"
outputDirs:
  - ..
tasks:
  - name: memories
    type: memory
    inputs:
      - memories/*.md
    outputs:
      - agent: roo
        outputPath: ".roo/rules/"  # Directory path (with trailing slash) for separate files
```

**Note**: You cannot mix this simplified format with the traditional `projects` section. Attempting to do so will result in an error.

## Template Functions

agent-def supports the following template functions in source files:

| Function | Description | Example |
|----------|-------------|---------|
| `file "path/to/file"` | Formats a file reference according to the output agent | `{{ file "src/main.go" }}` → `` `src/main.go` `` (Copilot) or `@/src/main.go` (Cline) |
| `include "path/to/file"` | Includes content from another file with template processing | `{{ include "common/header.md" }}` |
| `includeRaw "path/to/file"` | Includes content from another file without template processing | `{{ includeRaw "common/header.md" }}` |
| `reference "path/to/file"` | References another file's content with template processing | `{{ reference "data/config.json" }}` |
| `referenceRaw "path/to/file"` | References another file's content without template processing | `{{ referenceRaw "data/config.json" }}` |
| `mcp "agent" "command" "arg1" "arg2"` | Formats an MCP command for the output agent | `{{ mcp "github" "get-issue" "owner" "repo" "123" }}` |
| `agent` | Returns the current output agent identifier | `{{ if eq agent "claude" }}Claude-specific content{{ end }}` |
| `ifAGENT "content"` | Conditionally includes content only for the specified agent | `{{ ifRoo "This will only appear in Roo output" }}` |

### Template Function Examples

**File references:**
```
For Claude: {{ file "src/main.go" }} → @src/main.go
For Roo: {{ file "src/main.go" }} → @/src/main.go
For Cline: {{ file "src/main.go" }} → @/src/main.go
For Copilot: {{ file "src/main.go" }} → `src/main.go`
```

**Including templates:**
```
{{ include "header.md" }}
```
This will include and process the content of `header.md`, executing any template directives it contains.

**Including content without template processing:**
```
{{ includeRaw "header.md" }}
```
This will include the content of `header.md` exactly as it is, without executing any template directives it contains.

**Referencing files with template processing:**
```
{{ reference "data.json" }}
```
This will include the content of `data.json` with template processing.

**Referencing files without template processing:**
```
{{ referenceRaw "data.json" }}
```
This will include the raw content of `data.json` without any template processing.

**MCP commands:**
```
{{ mcp "mcp_server_name" "tool_name" "arguments" }}
```
This formats an MCP command for the output agent.

**Agent-specific content:**
```
{{ if eq agent "claude" }}
Claude-specific content here
{{ else if eq agent "roo" }}
Roo-specific content here
{{ end }}
```

**Simplified agent-specific content:**
```
{{ ifRoo "This content only appears in Roo output" }}
{{ ifClaude "This content only appears in Claude output" }}
```

## Command Line Arguments

agent-def supports the following commands and flags:

### Global Flags

These flags can be used with any command:

| Flag | Description |
|------|-------------|
| `--output, -o` | Output format (json, yaml, text) |
| `--verbose, -v` | Enable verbose output |
| `--log-file` | Log file path |
| `--log-level` | Log level (debug, info, warn, error) |
| `--debug` | Set log level to debug (shorthand for --log-level=debug) |

### Commands

#### `build`

Generates files based on agent-def.yml configuration.

Usage: `agent-def build [project...]`

If no project names are provided, all projects will be processed.

Flags:
- `--config, -c`: Path to agent-def.yml file or directory containing it (default: ".")
- `--dry-run`: Show what would be generated without writing files
- `--force, -f`: Force overwrite without prompting for confirmation

#### `validate`

Validates the agent-def.yml configuration without generating files.

Usage: `agent-def validate`

Flags:
- `--config, -c`: Path to agent-def.yml file or directory containing it (default: ".")

#### `list`

Lists available agents or project configurations.

Usage: `agent-def list [agents|projects]`

#### `init`

Initializes a new agent-def.yml configuration and sample files.

Usage: `agent-def init`

Flags:
- `--force, -f`: Force overwrite of existing files

## Error Handling and Troubleshooting

### Common Errors

| Error | Description | Solution |
|-------|-------------|----------|
| `configuration file not found` | The agent-def.yml file could not be found | Ensure you're in the correct directory or specify the path with `--config` |
| `invalid configuration format` | The configuration file has syntax errors | Check the YAML syntax for errors |
| `mixed configuration format` | Both simplified and standard formats were detected | Use either the simplified format or the standard format, not both |
| `agent not found` | The specified agent is not supported | Check for typos or use `agent-def list agents` to see supported agents |
| `file access denied` | Permission issues when reading/writing files | Check file permissions |

### Environment Variables

agent-def recognizes the following environment variables:

| Variable | Description |
|----------|-------------|
| `AGENT_DEF_LOG_FILE` | Log file path (overridden by `--log-file` flag) |
| `AGENT_DEF_LOG_LEVEL` | Log level (overridden by `--log-level` flag) |

### Debugging Tips

1. **Use `--dry-run`**: Preview what would be generated without writing files
2. **Enable debug logging**: Use `--debug` or `--log-level=debug` for detailed logs
3. **Check output path formatting**: Remember that paths ending with `/` are treated as directories, while paths without are treated as files
4. **Validate configuration**: Run `agent-def validate` to check your configuration before building
5. **Check file permissions**: Ensure you have read permissions for input files and write permissions for output directories

## Configuration Example

```yaml
configVersion: "1.0"

# Project-specific configuration
projects:
  my-app:
    # Output directories where generated files will be placed
    outputDirs:
      - ~/projects/my-app
      
    # Generation tasks for this project
    tasks:
      # Project-specific context information
      - name: "Project Memory"
        type: memory
        inputs:
          - ./projects/my-app/memories/
        outputs:
          # File path (no trailing slash) = concatenated output
          - agent: claude  # Outputs to single file CLAUDE.md
          - agent: roo     # Outputs to single file
            outputPath: ".roo/rules/context.md"
          
      # Project-specific commands
      - name: "Project Commands"
        type: command
        inputs:
          - ./projects/my-app/commands/
        outputs:
          # Directory path (trailing slash) = separate files
          - agent: claude
            outputPath: ".claude/commands/"  # Outputs to directory with separate files
          - agent: roo  # Outputs to single file .roomodes (default)

# User-level global configuration
user:
  tasks:
    # Global context information
    - name: "User Global Memory"
      type: memory
      inputs:
        - ./user/memories/
      outputs:
        - agent: claude
          outputPath: "~/.claude_global_context.md"  # Custom file output path (concatenated)
        - agent: roo
          outputPath: "~/.roo/rules/global.md"  # Custom file output path (concatenated)
          
    # Global commands
    - name: "User Global Commands"
      type: command
      inputs:
        - ./user/commands/
      outputs:
        - agent: claude
          outputPath: "~/.claude/commands/"  # Directory path (trailing slash) = separate files
        - agent: roo  # Default path for user commands (concatenated)
```

## OS-Specific Considerations

### Path Handling

- **Windows**: Both forward slashes (`/`) and backslashes (`\`) are supported in paths. Home directory expansion (`~`) works with the user's profile directory.
- **macOS/Linux**: Use forward slashes (`/`) for paths. Home directory expansion (`~`) works with the user's home directory.

### Default Output Paths

Some default paths differ between operating systems:

- **Windows**:
  - Roo user commands: `%APPDATA%\Code\User\globalStorage\rooveterinaryinc.roo-cline\settings\custom_modes.yaml`
  - Cline user memory: `%USERPROFILE%\Documents\Cline\Rules\`

- **macOS**:
  - Roo user commands: `~/Library/Application Support/Code/User/globalStorage/rooveterinaryinc.roo-cline/settings/custom_modes.yaml`
  - Cline user memory: `~/Documents/Cline/Rules/`

- **Linux**:
  - Roo user commands: `~/.config/Code/User/globalStorage/rooveterinaryinc.roo-cline/settings/custom_modes.yaml`
  - Cline user memory: `~/Documents/Cline/Rules/`

## Best Practices

1. **Organize Project Structure**: Store related memories and commands in an organized directory structure
2. **Share Common Settings**: Reuse common tasks across multiple projects
3. **Use Explicit Paths**: Use tilde (~) rather than absolute paths, especially for custom outputPath values
4. **Meaningful Naming**: Give tasks meaningful names that clearly indicate their purpose
5. **Use Templates**: Leverage template functions to maintain DRY (Don't Repeat Yourself) principles
6. **Agent-Specific Content**: Use the `ifAGENT` function to include content specific to certain agents
7. **Version Control**: Keep your agent-def configuration and template files under version control
8. **Modular Content**: Break down large memory files into smaller, focused files for easier maintenance
9. **Test Changes**: Use `--dry-run` to preview changes before applying them
10. **Document Configuration**: Add comments to your configuration file to explain complex setups