---
layout: default
title: Input and Output Processing
---

[Back to index](index.md)

# Input and Output Processing

This document details how agent-sync processes input files and generates output files.

## Glob Pattern Support in Inputs

The `inputs` field supports glob patterns with doublestar support for recursive matching and exclusions. Input paths are always resolved relative to the configuration file directory (where your `agent-sync.yml` is located):

| Pattern | Description |
|---------|-------------|
| `*.md` | All Markdown files in the current directory |
| `**/*.md` | All Markdown files in any subdirectory (recursive) |
| `!*_test.md` | Exclude all test Markdown files |

For detailed information about glob pattern syntax and behavior, see the [Glob Patterns](glob-patterns.md) documentation.

### Examples

```yaml
inputs:
  # Include all Markdown files in commands directory and subdirectories
  - "commands/**/*.md"
  # Exclude test files
  - "!commands/**/*_test.md"
  # Exclude specific temporary file
  - "!commands/temp.md"
```

### Pattern Order and Precedence

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

## Handling of Relative Output Directories

When specifying output directories in the `outputDirs` setting:

- Absolute paths (starting with `/` on Unix or drive letters on Windows) are used as-is
- Relative paths (not starting with `/` or drive letters) are:
  - Treated as relative to the agent-sync.yml file location
  - Always converted to absolute paths internally for clarity and precision

This approach allows you to specify output directories relative to your configuration file, making configurations more portable across different environments, while ensuring that all paths used internally are unambiguous.

Example:
```yaml
# If agent-sync.yml is located at /path/to/project/agent-sync.yml
outputDirs:
  - dist  # Will be resolved to absolute path: /path/to/project/dist
```

## Output Path Interpretation

The format of the `outputPath` determines how files are processed:

- Paths ending with `/` are treated as **directories** where individual files will be placed
- Paths not ending with `/` are treated as **files** where all input content will be concatenated

The `outputPath` can be specified as either relative or absolute:
- Absolute paths (starting with `/` on Unix or drive letters on Windows) are used as-is
- Relative paths are resolved relative to each output directory

For example:
- `.roo/rules/` - Treated as a directory, each source file becomes a separate output file
- `.roo/rules/combined.md` - Treated as a file, all source files are concatenated into one file
- `/absolute/path/rules/` - Absolute path to a directory for separate output files

If `outputPath` is not specified, the agent's default path is used (see "Default Output Path Behavior" sections below).

## Default Output Path Behavior

When no `outputPath` is specified, the following default paths are used. The path format (with or without trailing slash) determines whether files are concatenated or kept separate:

| Agent | Task Type | Default Path | Concatenation Behavior |
|-------|-----------|-------------|------------------------|
| Claude | memory | `CLAUDE.md` or `~/.claude/CLAUDE.md` | Concatenated (file path) |
| Claude | command | `.claude/commands/` | Non-concatenated (directory path) |
| Roo | memory | `.roo/rules/` | Non-concatenated (directory path) |
| Roo | command | `.roo/commands/` or `~/.roo/commands/` | Non-concatenated (directory path) |
| Cline | memory | `.clinerules/` or `~/Documents/Cline/Rules/` | Non-concatenated (directory path) |
| Cline | command | `.clinerules/workflows/` or `~/Documents/Cline/Workflows/` | Non-concatenated (directory path) |
| Copilot | memory | `.github/copilot-instructions.md` or `~/.vscode/copilot-instructions.md` | Concatenated (file path) |
| Copilot | command | `.github/prompts/` or `~/.vscode/prompts/` | Non-concatenated (directory path) |

## Task Processing Workflow
When the `agent-sync apply` command is executed, the following workflow occurs:

1. **Configuration Loading**: The configuration file is loaded and validated.
2. **Task Selection**: If specific projects are provided as command arguments, only those projects' tasks are processed. Otherwise, all projects and user tasks are processed.
3. **Input Processing**:
   - For each task, the input files are resolved based on glob patterns relative to the configuration file directory.
   - Absolute paths are constructed internally for all input files to ensure precise handling.
   - Each input file is processed through the template engine if applicable.
4. **Output Generation**:
   - For each output agent specified in the task, the content is formatted according to the agent's requirements.
   - The formatted content is written to the specified output paths in each output directory.
5. **Verification**: Optionally, with the `--dry-run` flag, the process will show what would be generated without actually writing files. The enhanced output displays content organized by agent with file status ([CREATE], [MODIFY], [UNCHANGED]), file paths, sizes, and per-agent summaries.

## Navigation

- [Main Configuration Guide](config.md)
- [Configuration Reference](config-reference.md)
- [Task Types](task-types.md)
- [Template System](templates.md)
- [Glob Patterns](glob-patterns.md)
- [Command Line Interface](cli.md)
- [Troubleshooting](troubleshooting.md)
- [Examples and Best Practices](examples.md)

---

| Previous | Next |
|----------|------|
| [Template System](templates.md) | [Glob Patterns](glob-patterns.md) |

[Back to index](index.md)