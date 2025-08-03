---
layout: default
title: Command Line Interface
---

[Back to index](index.md)

# Command Line Interface

agent-sync provides a command-line interface (CLI) with several commands and options to control its behavior.

## Global Flags

These flags can be used with any command:

| Flag | Description |
|------|-------------|
| `--output, -o` | Output format (json, yaml, text) |
| `--verbose, -v` | Enable verbose output |
| `--log-file` | Log file path |
| `--log-level` | Log level (debug, info, warn, error) |
| `--debug` | Set log level to debug (shorthand for --log-level=debug and enables console output same as --verbose) |

## Environment Variables

Most command-line flags can be configured via environment variables using the `AGENT_SYNC_` prefix.

For details on available environment variables, run:
```bash
agent-sync --help
```

## Commands

### `apply`

Generates files based on agent-sync.yml configuration.

Usage: `agent-sync apply`

All projects will be processed.

Flags:
- `--config, -c`: Path to agent-sync.yml file or directory containing it (default: ".")
- `--dry-run`: Show what would be generated without writing files. The output provides detailed information organized by agent, including file status ([CREATE], [MODIFY], or [UNCHANGED]), file paths, sizes, and summaries showing counts of created, modified, and unchanged files
- `--force, -f`: Force overwrite without prompting for confirmation

### `init`

Initializes a new agent-sync.yml configuration and sample files.

Usage: `agent-sync init`

Flags:
- `--force, -f`: Force overwrite of existing files

## Examples

**Applying with verbose output:**
```bash
agent-sync apply --verbose
```

**Initializing a new configuration:**
```bash
agent-sync init
```

**Applying with debug logs:**
```bash
# Shorthand: sets level=debug and enables console output (implies verbose)
agent-sync apply --debug

# Equivalent explicit form (debug level + console output)
agent-sync apply --log-level debug --verbose
```

Note: `--log-level debug` without `--verbose` does not enable console output; `--debug` implies verbose and will enable console output.

**Dry run to preview changes:**
```bash
agent-sync apply --dry-run
```

For more information about logging configuration, see the [Logging Guide](logging.md).

## Navigation

- [Main Configuration Guide](config.md)
- [Configuration Reference](config-reference.md)
- [Input and Output Processing](input-output.md)
- [Task Types](task-types.md)
- [Template System](templates.md)
- [Troubleshooting](troubleshooting.md)
- [Examples and Best Practices](examples.md)

---

| Previous | Next |
|----------|------|
| [Task Types](task-types.md) | [Template System](templates.md) |

[Back to index](index.md)