---
layout: default
title: Troubleshooting
---

[Back to index](index.md)

# Error Handling and Troubleshooting

This document covers common errors, environment variables, debugging tips, and OS-specific considerations for agent-sync.

## Common Errors

| Error | Description | Solution |
|-------|-------------|----------|
| `configuration file not found` | The agent-sync.yml file could not be found | Ensure you're in the correct directory or specify the path with `--config` |
| `invalid configuration format` | The configuration file has syntax errors | Check the YAML syntax for errors |
| `mixed configuration format` | Both simplified and standard formats were detected | Use either the simplified format or the standard format, not both |
| `agent not found` | The specified agent is not supported | Check for typos or use `agent-sync list agents` to see supported agents |
| `file access denied` | Permission issues when reading/writing files | Check file permissions |

## Environment Variables

agent-sync recognizes the following environment variables:

| Variable | Description |
|----------|-------------|
| `AGENT_SYNC_LOG_FILE` | Log file path (overridden by `--log-file` flag) |
| `AGENT_SYNC_LOG_LEVEL` | Log level (overridden by `--log-level` flag) |

## Debugging Tips

1. **Use `--dry-run`**: Preview what would be generated without writing files
2. **Enable debug logging**: Use `--debug` or `--log-level=debug` for detailed logs
3. **Check output path formatting**: Remember that paths ending with `/` are treated as directories, while paths without are treated as files
4. **Check file permissions**: Ensure you have read permissions for input files and write permissions for output directories
5. **Examine template processing**: If you suspect issues with templates, try creating simpler templates first to isolate the problem
6. **Check for conflicting paths**: Make sure you're not writing to the same output file from different tasks

## Using Enhanced Dry-Run Mode

The `--dry-run` flag provides a powerful way to preview and debug your configuration without making actual changes. The enhanced output format offers:

1. **Agent-based organization**: Output is grouped by agent, making it easy to see how files will be processed for each target system

2. **Status indicators**: Each file is prefixed with one of these status tags:
   - `[CREATE]`: File does not exist and would be created
   - `[MODIFY]`: File exists and its content would be changed
   - `[UNCHANGED]`: File exists and its content would remain the same

3. **File details**: For each file, you'll see:
   - Full file path
   - File size in human-readable format (B, KB, MB)

4. **Agent summaries**: For each agent, a summary is provided showing:
   - Number of files that would be created
   - Number of files that would be modified
   - Number of files that would remain unchanged

Example of using dry-run to debug configuration issues:

```bash
# Run in dry-run mode to check configuration
agent-sync apply --dry-run

# If you need more details, combine with debug logging
agent-sync apply --dry-run --debug
```

This helps identify issues such as:
- Files being created in unexpected locations
- Expected modifications not being detected
- Unnecessary file duplication
- Size discrepancies indicating template processing issues

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

## Troubleshooting Log Issues

If you're having issues with logs:

1. **Log file not being created**:
   - Check if the directory exists and has write permissions
   - Use an absolute path to ensure correct location

2. **Missing detailed information**:
   - Ensure log level is set to `debug`
   - Enable verbose output with `--verbose`

3. **Too much log information**:
   - Increase log level to `warn` or `error`
   - Disable verbose output

For more information about logging configuration, see the [Logging Guide](logging.md).

## Navigation

- [Main Configuration Guide](config.md)
- [Configuration Reference](config-reference.md)
- [Input and Output Processing](input-output.md)
- [Task Types](task-types.md)
- [Template System](templates.md)
- [Command Line Interface](cli.md)
- [Examples and Best Practices](examples.md)
- [Glob Patterns](glob-patterns.md)
- [Logging](logging.md)

---

| Previous | Next |
|----------|------|
| [Examples and Best Practices](examples.md) | |

[Back to index](index.md)