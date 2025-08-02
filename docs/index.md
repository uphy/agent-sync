---
layout: default
title: agent-sync Documentation
---

# Agent Sync Documentation

## Introduction

Agent Sync (`agent-sync`) is a tool that streamlines working with multiple AI assistants by letting you write your context and commands once and automatically converting them to formats compatible with different AI agents.

Rather than maintaining separate files for each AI assistant, agent-sync allows you to define a single source of truth for your AI instructions and automatically generate the appropriate files for Claude, Roo, Cline, Copilot, and other supported agents.

## Key Features

- **Multi-Agent Support**: Convert context (memory), commands (Roo slash commands, Claude commands, etc.), and modes (Claude Code subagents and Roo custom modes) to formats compatible with Claude, Roo, Cline, Copilot, and more
- **Write Once, Deploy Everywhere**: Maintain a single source of truth for AI instructions instead of separate files for each assistant
- **Multi-Agent Support**: Convert context (memory) and commands to formats compatible with Claude, Roo, Cline, Copilot, and more
- **Template System**: Use powerful templating to dynamically generate content with agent-specific formatting
- **Project and User Separation**: Maintain both project-specific and global user-level configurations
- **Flexible Output**: Control how content is organized (concatenated into single files or split into multiple files)

## How It Works

1. You write context (memory) and commands in standard markdown files
2. You configure how these files should be processed in `agent-sync.yml`
3. You run `agent-sync apply` to generate agent-specific files in the correct formats and locations
4. Your AI assistants now have the same knowledge and capabilities across platforms

## Documentation Contents

### Core Concepts

- [**Installation**](installation.md) - Installation methods for different platforms
- [**Configuration Guide**](config.md) - Main configuration concepts and structure
- [**Configuration Reference**](config-reference.md) - Detailed configuration parameters and settings
- [**Task Types**](task-types.md) - Memory and command task types with agent-specific details

### Features and Usage

- [**Command Line Interface**](cli.md) - Command line options and usage
- [**Template System**](templates.md) - Template functions and path resolution
- [**Input and Output Processing**](input-output.md) - Input handling, output paths, and file processing

### Advanced Topics

- [**Glob Patterns**](glob-patterns.md) - Pattern syntax for selecting input files
- [**Logging**](logging.md) - Configuring and using the logging system
- [**Examples and Best Practices**](examples.md) - Configuration examples and recommendations
- [**Troubleshooting**](troubleshooting.md) - Error handling and debugging tips

## Getting Started

To get started with agent-sync, install it using one of the methods described in the [Installation Guide](installation.md) and run the `init` command to create a sample configuration:

```bash
agent-sync init
```

Then customize the generated `agent-sync.yml` file and the content in the `memories/` and `commands/` directories according to your needs, and run the `apply` command:

```bash
agent-sync apply
```

For more details on installation, configuration, and usage, explore the documentation links above.