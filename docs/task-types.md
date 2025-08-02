---
layout: default
title: Task Types
---

[Back to index](index.md)

# Task Types

agent-sync supports the following two task types:

## 1. Memory (`type: memory`)

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

## 2. Command (`type: command`)

Provides custom command definitions for AI agents. This creates shortcuts for performing specific tasks.

**Default output locations:**

| Agent | Scope | Default Path | Description |
|-------|-------|-------------|-------------|
| Claude | User | `~/.claude/commands/{filename}.md` | User's global Claude command file |
| Claude | Project | `.claude/commands/{filename}.md` | Project-specific Claude command file |
| Roo | User | `~/.roo/commands/{filename}.md` | User's global Roo slash command file (per file) |
| Roo | Project | `.roo/commands/{filename}.md` | Project-specific Roo slash command file (per file) |
| Cline | User | `~/Documents/Cline/Workflows/{filename}.md` | User's global Cline workflow file |
| Cline | Project | `.clinerules/workflows/{filename}.md` | Project-specific Cline workflow file |
| Copilot | User | `~/.vscode/prompts/{filename}.prompt.md` | User's global Copilot prompt file |
| Copilot | Project | `.github/prompts/{filename}.prompt.md` | Project-specific Copilot prompt file |

Note: For Claude, Cline, Copilot, and similar agents, `{filename}` is derived from the input file's basename (the filename without its directory path). For example, an input file named `my-project/commands/deploy.md` would result in an output file named `deploy.md` in the appropriate output directory.

## Agent-specific Command Frontmatter

Each agent may support specific frontmatter attributes for commands:

### Claude Command Frontmatter

- `claude.description`: Brief description of the command shown in the help menu
- `claude.allowed-tools`: List of tools the command is permitted to use

### Roo Frontmatter (slash commands)

Supported frontmatter:
- Top-level:
  - `description`: Brief description of what the slash command does
- Roo-specific (under `roo:`):
  - `argument-hint` (optional): Short hint string for expected arguments (e.g., "[path]" or "[owner repo]")

Legacy compatibility:
- If no top-level `description` is provided, agent-sync maps legacy Roo fields:
  - `description` is synthesized from `roo.roleDefinition` if present, otherwise from `roo.name`
- Legacy `roo.whenToUse` and `roo.groups` are ignored in slash command output

### Cline Frontmatter

- Cline workflows use standard markdown without special frontmatter requirements

### Copilot Command Frontmatter

- `mode`: The operational mode for Copilot (e.g., "chat", "inline")
- `model`: The specific model to use
- `tools`: List of available tools
- `description`: Brief description of the prompt's purpose

## Navigation

- [Main Configuration Guide](config.md)
- [Configuration Reference](config-reference.md)
- [Input and Output Processing](input-output.md)
- [Template System](templates.md)
- [Command Line Interface](cli.md)
- [Troubleshooting](troubleshooting.md)
- [Examples and Best Practices](examples.md)

---

| Previous | Next |
|----------|------|
| [Configuration Reference](config-reference.md) | [Command Line Interface](cli.md) |

[Back to index](index.md)