---
layout: default
title: Logging Guide
---

[Back to index](index.md)

# agent-sync Logging Guide

## Overview

The agent-sync tool includes a flexible logging system that helps you:

1. **See what's happening** with clear, color-coded messages in the terminal when logging is enabled
2. **Keep records** of operations for troubleshooting and auditing
3. **Control the amount of information** displayed based on your needs

This guide explains how to configure and use logging features in agent-sync.

## Quick Start

By default, agent-sync operates silently with no console output. To enable output, use one of these options:

1. Use the `--verbose` flag to show output in your terminal:

```bash
agent-sync apply --verbose
```

The `--verbose` flag does two important things:
- Enables the logging system (sets `Enabled = true`)
- Sets console output to be displayed in your terminal

2. Use the `--log-file` flag to save logs to a file:

```bash
agent-sync apply --log-file apply.log
```

When enabled, output includes visual indicators:
- `✓` Success messages (green)
- `➜` Progress indicators (blue)
- `ℹ` General information
- `✗` Error messages (red)

**Note:** There is a known issue where DEBUG level logs may not appear even when logging is enabled with the `--verbose` flag. This issue is under investigation.

## Configuration Methods

You can configure logging behavior in two ways:

### 1. Command-Line Flags

These flags provide immediate control for a single command:

```bash
# Increase logging detail
agent-sync apply --log-level debug

# Save logs to a file
agent-sync apply --log-file ./logs/agent-sync.log

# Display more detailed messages
agent-sync apply --verbose
```

### 2. Environment Variables

For persistent settings across multiple commands:

```bash
# Set these variables in your shell or .bashrc/.zshrc file
export AGENT_DEF_LOG_LEVEL=debug
export AGENT_DEF_LOG_FILE=./logs/agent-sync.log
```

## Configuration Priority

When multiple configuration sources exist, agent-sync uses this priority order:

1. Command-line flags (highest priority)
2. Environment variables
3. Default values (lowest priority)

By default, logging is disabled. You must explicitly enable it using one of the configuration methods above.

## Log Levels

agent-sync uses four log levels, controlling how much information is recorded:

| Level | Description | Use When |
|-------|-------------|---------|
| `debug` | Most detailed information | Troubleshooting or investigating issues |
| `info` | General operational messages | Normal operation, default setting |
| `warn` | Potential issues | Want to see only warnings and errors |
| `error` | Errors and failures | Want to see only errors |

## Log File Settings

When logging to a file, you can control:

| Setting | Description | Default |
|---------|-------------|---------|
| `file` | Path to the log file | None (no file logging) |
| `max_size` | Maximum file size in MB before rotation | 10 |
| `max_age` | Days to retain old log files | 30 |
| `max_files` | Number of old log files to keep | 3 |

## Common Use Cases

### Troubleshooting a Problem

When something isn't working as expected:

```bash
agent-sync apply --log-level debug --log-file debug.log --verbose
```

This will:
1. Show detailed logs in the terminal
2. Save comprehensive logs to debug.log for later analysis

### Running in a Script

When using agent-sync in automated scripts:

```bash
agent-sync apply --log-level error --log-file apply.log
```

This will:
1. Only show errors in the terminal
2. Save all logs to apply.log

### Development Environment

For day-to-day development work:

```bash
# Set environment variables
export AGENT_DEF_LOG_LEVEL=info
# Or use command-line flags
agent-sync apply --verbose
```

### Production Environment

For server or CI/CD environments:

```bash
# Set environment variables
export AGENT_DEF_LOG_LEVEL=info
export AGENT_DEF_LOG_FILE=/var/log/agent-sync/agent-sync.log
# Or use command-line flags
agent-sync apply --log-level info --log-file /var/log/agent-sync/agent-sync.log
```

### Silent Operation

For minimal output when running in scripts or automated environments, simply don't set any logging options. By default, agent-sync will not produce any log output unless explicitly enabled.

## Troubleshooting

### Common Issues

1. **Log file not being created**
   - Check if the directory exists and has write permissions
   - Use an absolute path to ensure correct location

2. **Missing detailed information**
   - Ensure log level is set to `debug`
   - Enable verbose output with `--verbose`

3. **Too much log information**
   - Increase log level to `warn` or `error`
   - Disable verbose output

4. **No color in terminal output**
   - Some terminals don't support colors - set `color: false` in config

### Getting Help

If you encounter issues with logging, try running with debug level to get more information:

```bash
agent-sync --log-level debug --verbose apply
```

## Testing Logs

If you want to verify that your logging configuration is working correctly, you can use the `log-test` command:

```bash
# Test logging with default settings (silent)
agent-sync log-test

# Test logging with verbose output enabled
agent-sync --verbose log-test

# Test logging with a specific log level
agent-sync --log-level debug --verbose log-test

# Save test logs to a file
agent-sync --log-file test.log log-test
```

The `log-test` command:
- Outputs messages at all log levels (debug, info, warn, error)
- Includes structured logging examples with additional fields
- Helps verify that logging is configured correctly
- Is particularly useful for testing the `--verbose` flag behavior

This command is valuable for developers who need to verify the logging system is working as expected or to troubleshoot logging issues.

## Reference

### Complete Configuration Options

| Option | CLI Flag | Env Variable | Description | Default |
|--------|----------|--------------|-------------|---------|
| Logging Enabled | `--verbose` | - | Enable/disable logging | `false` |
| Log Level | `--log-level` | `AGENT_DEF_LOG_LEVEL` | Logging detail (debug, info, warn, error) | `info` |
| Log File | `--log-file` | `AGENT_DEF_LOG_FILE` | Path to save logs | - |
| Verbose | `--verbose` | - | Show detailed messages | `false` |

## Navigation

- [Main Configuration Guide](config.md)
- [Configuration Reference](config-reference.md)
- [Task Types](task-types.md)
- [Template System](templates.md)
- [Input and Output Processing](input-output.md)
- [Glob Patterns](glob-patterns.md)
- [Command Line Interface](cli.md)
- [Examples and Best Practices](examples.md)
- [Troubleshooting](troubleshooting.md)

---

| Previous | Next |
|----------|------|
| [Glob Patterns](glob-patterns.md) | [Examples and Best Practices](examples.md) |

[Back to index](index.md)