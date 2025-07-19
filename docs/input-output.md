# Input and Output Processing

This document details how agent-def processes input files and generates output files.

## Glob Pattern Support in Inputs

The `inputs` field supports glob patterns with doublestar support for recursive matching and exclusions:

| Pattern | Description |
|---------|-------------|
| `*.md` | All Markdown files in the current directory |
| `**/*.md` | All Markdown files in any subdirectory (recursive) |
| `!*_test.md` | Exclude all test Markdown files |

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

## Navigation

- [Main Configuration Guide](config.md)
- [Configuration Reference](config-reference.md)
- [Task Types](task-types.md)
- [Template System](templates.md)
- [Command Line Interface](cli.md)
- [Troubleshooting](troubleshooting.md)
- [Examples and Best Practices](examples.md)