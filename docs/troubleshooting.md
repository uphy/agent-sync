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
| `AGENT_DEF_LOG_FILE` | Log file path (overridden by `--log-file` flag) |
| `AGENT_DEF_LOG_LEVEL` | Log level (overridden by `--log-level` flag) |

## Debugging Tips

1. **Use `--dry-run`**: Preview what would be generated without writing files
2. **Enable debug logging**: Use `--debug` or `--log-level=debug` for detailed logs
3. **Check output path formatting**: Remember that paths ending with `/` are treated as directories, while paths without are treated as files
4. **Dry-run**: Run `agent-sync apply --dry-run` to check your configuration before applying
5. **Check file permissions**: Ensure you have read permissions for input files and write permissions for output directories
6. **Examine template processing**: If you suspect issues with templates, try creating simpler templates first to isolate the problem
7. **Check for conflicting paths**: Make sure you're not writing to the same output file from different tasks

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