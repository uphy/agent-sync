---
layout: default
title: Examples and Best Practices
---

[Back to index](index.md)

# Examples and Best Practices

This document provides configuration examples and best practices for using agent-sync effectively.

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

## Simplified Configuration Example

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

## Best Practices

1. **Organize Project Structure**: Store related memories and commands in an organized directory structure
   ```
   project/
   ├── memories/
   │   ├── architecture.md
   │   ├── coding-standards.md
   │   └── deployment.md
   └── commands/
       ├── deploy.md
       ├── test.md
       └── debug.md
   ```

2. **Share Common Settings**: Reuse common tasks across multiple projects by using templates and includes

3. **Use Explicit Paths**: Use tilde (~) rather than absolute paths, especially for custom outputPath values
   ```yaml
   outputPath: "~/.roo/rules/"  # Good: Portable across users
   outputPath: "/Users/alice/rules/"  # Bad: Hard-coded user path
   ```

4. **Meaningful Naming**: Give tasks meaningful names that clearly indicate their purpose
   ```yaml
   name: "Coding Standards"  # Good: Clear purpose
   name: "Task1"  # Bad: Generic name
   ```

5. **Use Templates**: Leverage template functions to maintain DRY (Don't Repeat Yourself) principles
   ```markdown
{% raw %}
   {{ include "common/header.md" }}
{% endraw %}
   ```

6. **Agent-Specific Content**: Use the `ifAGENT` function to include content specific to certain agents
   ```markdown
{% raw %}
   {{ ifRoo "This appears only in Roo output" }}
{% endraw %}
   ```

7. **Version Control**: Keep your agent-sync configuration and template files under version control

8. **Modular Content**: Break down large memory files into smaller, focused files for easier maintenance

9. **Test Changes**: Use `--dry-run` to preview changes before applying them
   ```bash
   agent-sync apply --dry-run
   ```
   
   Example output from enhanced dry-run mode:
   ```
   DRY RUN MODE: No files will actually be written
   Executing task 'Project Memory' of type 'memory'
   
   Agent: claude
     [CREATE] /Users/alice/projects/my-app/CLAUDE.md (2.3 KB)
   
     Summary: 1 files would be created, 0 files would be modified, 0 files would remain unchanged
   
   Agent: roo
     [MODIFY] /Users/alice/projects/my-app/.roo/rules/architecture.md (1.1 KB)
     [CREATE] /Users/alice/projects/my-app/.roo/rules/coding-standards.md (856 B)
     [UNCHANGED] /Users/alice/projects/my-app/.roo/rules/deployment.md (512 B)
   
     Summary: 1 files would be created, 1 files would be modified, 1 files would remain unchanged
   
   Executing task 'Project Commands' of type 'command'
   
   Agent: claude
     [CREATE] /Users/alice/projects/my-app/.claude/commands/deploy.md (734 B)
     [CREATE] /Users/alice/projects/my-app/.claude/commands/test.md (512 B)
   
     Summary: 2 files would be created, 0 files would be modified, 0 files would remain unchanged
   
   Agent: roo
     [MODIFY] /Users/alice/projects/my-app/.roomodes (1.2 KB)
   
     Summary: 0 files would be created, 1 files would be modified, 0 files would remain unchanged
   ```

10. **Document Configuration**: Add comments to your configuration file to explain complex setups
    ```yaml
    # This outputs to separate files in the .roo/rules/ directory
    # Each input file becomes its own output file
    outputs:
      - agent: roo
        outputPath: ".roo/rules/"
    ```

## Common Patterns

### Multiple Output Destinations

```yaml
outputDirs:
  - ~/projects/main
  - ~/projects/worktree
```

### Agent-Specific Output Paths

```yaml
outputs:
  - agent: claude
    outputPath: "CLAUDE.md"
  - agent: roo
    outputPath: ".roo/rules/"
  - agent: cline
    outputPath: ".clinerules/"
```

### Selective File Processing

```yaml
inputs:
  - "core/**/*.md"
  - "!core/experimental/**/*.md"  # Exclude experimental files
```

## Navigation

- [Main Configuration Guide](config.md)
- [Configuration Reference](config-reference.md)
- [Input and Output Processing](input-output.md)
- [Task Types](task-types.md)
- [Template System](templates.md)
- [Command Line Interface](cli.md)
- [Troubleshooting](troubleshooting.md)
- [Glob Patterns](glob-patterns.md)
- [Logging](logging.md)

---

| Previous | Next |
|----------|------|
| [Logging](logging.md) | [Troubleshooting](troubleshooting.md) |

[Back to index](index.md)