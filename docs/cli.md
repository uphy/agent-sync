# Command Line Interface

agent-def provides a command-line interface (CLI) with several commands and options to control its behavior.

## Global Flags

These flags can be used with any command:

| Flag | Description |
|------|-------------|
| `--output, -o` | Output format (json, yaml, text) |
| `--verbose, -v` | Enable verbose output |
| `--log-file` | Log file path |
| `--log-level` | Log level (debug, info, warn, error) |
| `--debug` | Set log level to debug (shorthand for --log-level=debug) |

## Commands

### `build`

Generates files based on agent-def.yml configuration.

Usage: `agent-def build [project...]`

If no project names are provided, all projects will be processed.

Flags:
- `--config, -c`: Path to agent-def.yml file or directory containing it (default: ".")
- `--dry-run`: Show what would be generated without writing files
- `--force, -f`: Force overwrite without prompting for confirmation

### `validate`

Validates the agent-def.yml configuration without generating files.

Usage: `agent-def validate`

Flags:
- `--config, -c`: Path to agent-def.yml file or directory containing it (default: ".")

### `list`

Lists available agents or project configurations.

Usage: `agent-def list [agents|projects]`

### `init`

Initializes a new agent-def.yml configuration and sample files.

Usage: `agent-def init`

Flags:
- `--force, -f`: Force overwrite of existing files

## Examples

**Building all projects with verbose output:**
```bash
agent-def build --verbose
```

**Building specific projects:**
```bash
agent-def build project1 project2 --verbose
```

**Validating configuration:**
```bash
agent-def validate --config /path/to/agent-def.yml
```

**Initializing a new configuration:**
```bash
agent-def init
```

**Listing supported agents:**
```bash
agent-def list agents
```

**Building with debug logs:**
```bash
agent-def build --debug --log-file build.log
```

**Dry run to preview changes:**
```bash
agent-def build --dry-run
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