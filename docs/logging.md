# Agent-def Logging Guide

## Overview

The agent-def tool includes a flexible logging system that helps you:

1. **See what's happening** with clear, color-coded messages in the terminal when logging is enabled
2. **Keep records** of operations for troubleshooting and auditing
3. **Control the amount of information** displayed based on your needs

This guide explains how to configure and use logging features in agent-def.

## Quick Start

By default, agent-def operates silently with no console output. To enable output, use one of these options:

1. Use the `--verbose` flag to show output in your terminal:

```bash
agent-def build --verbose
```

The `--verbose` flag does two important things:
- Enables the logging system (sets `Enabled = true`)
- Sets console output to be displayed in your terminal

2. Use the `--log-file` flag to save logs to a file:

```bash
agent-def build --log-file build.log
```

When enabled, output includes visual indicators:
- `✓` Success messages (green)
- `➜` Progress indicators (blue)
- `ℹ` General information
- `✗` Error messages (red)

**Note:** There is a known issue where DEBUG level logs may not appear even when logging is enabled with the `--verbose` flag. This issue is under investigation.

## Configuration Methods

You can configure logging behavior in three ways:

### 1. Command-Line Flags

These flags provide immediate control for a single command:

```bash
# Increase logging detail
agent-def build --log-level debug

# Save logs to a file
agent-def build --log-file ./logs/agent-def.log

# Display more detailed messages
agent-def build --verbose
```

### 2. Environment Variables

For persistent settings across multiple commands:

```bash
# Set these variables in your shell or .bashrc/.zshrc file
export AGENT_DEF_LOG_LEVEL=debug
export AGENT_DEF_LOG_FILE=./logs/agent-def.log
```

### 3. Configuration File

For project-specific settings, add a `logging` section to your `agent-def.yml`:

```yaml
logging:
  enabled: true
  level: info
  file: ./logs/agent-def.log
  max_size: 10
  max_age: 30
  max_files: 3
  format: text
  console_output: false
  color: true
  verbose: false
```

## Configuration Priority

When multiple configuration sources exist, agent-def uses this priority order:

1. Command-line flags (highest priority)
2. Environment variables
3. Configuration file
4. Default values (lowest priority)

By default, logging is disabled. You must explicitly enable it using one of the configuration methods above.

## Log Levels

Agent-def uses four log levels, controlling how much information is recorded:

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
agent-def build --log-level debug --log-file debug.log --verbose
```

This will:
1. Show detailed logs in the terminal
2. Save comprehensive logs to debug.log for later analysis

### Running in a Script

When using agent-def in automated scripts:

```bash
agent-def build --log-level error --log-file build.log
```

This will:
1. Only show errors in the terminal
2. Save all logs to build.log

### Development Environment

For day-to-day development work:

```yaml
# In agent-def.yml
logging:
  enabled: true
  level: info
  format: text
  console_output: true
  color: true
  verbose: false
```

### Production Environment

For server or CI/CD environments:

```yaml
# In agent-def.yml
logging:
  enabled: true
  level: info
  file: /var/log/agent-def/agent-def.log
  max_size: 100
  max_age: 30
  max_files: 10
  format: json
  console_output: false
  verbose: false
```

### Silent Operation

For minimal output when running in scripts or automated environments:

```yaml
# In agent-def.yml
logging:
  enabled: false  # This is the default
```

With this configuration (or by default), agent-def will not produce any log output.

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
agent-def --log-level debug --verbose build
```

## Testing Logs

If you want to verify that your logging configuration is working correctly, you can use the `log-test` command:

```bash
# Test logging with default settings (silent)
agent-def log-test

# Test logging with verbose output enabled
agent-def --verbose log-test

# Test logging with a specific log level
agent-def --log-level debug --verbose log-test

# Save test logs to a file
agent-def --log-file test.log log-test
```

The `log-test` command:
- Outputs messages at all log levels (debug, info, warn, error)
- Includes structured logging examples with additional fields
- Helps verify that logging is configured correctly
- Is particularly useful for testing the `--verbose` flag behavior

This command is valuable for developers who need to verify the logging system is working as expected or to troubleshoot logging issues.

## Reference

### Complete Configuration Options

| Option | CLI Flag | Env Variable | Config File | Description | Default |
|--------|----------|--------------|-------------|-------------|---------|
| Logging Enabled | - | - | `enabled` | Enable/disable logging | `false` |
| Log Level | `--log-level` | `AGENT_DEF_LOG_LEVEL` | `level` | Logging detail (debug, info, warn, error) | `info` |
| Log File | `--log-file` | `AGENT_DEF_LOG_FILE` | `file` | Path to save logs | - |
| Verbose | `--verbose` | - | `verbose` | Show detailed messages | `false` |
| - | - | - | `format` | Output format (text, json) | `text` |
| - | - | - | `color` | Enable colored output | `true` |
| - | - | - | `console_output` | Enable console output | `false` |
| - | - | - | `max_size` | Max log file size (MB) | `10` |
| - | - | - | `max_age` | Days to retain logs | `30` |
| - | - | - | `max_files` | Number of log files to keep | `3` |