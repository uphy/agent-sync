---
layout: default
title: Installation
---

[Back to index](index.md)

# Installation

This guide covers the different methods available for installing agent-sync.

## Homebrew Installation (macOS)

You can install `agent-sync` using Homebrew:

```bash
brew install uphy/tap/agent-sync
```

## Go Installation

Install via Go:

```bash
go install github.com/uphy/agent-sync/cmd/agent-sync@latest
```

## Binary Downloads

Download pre-compiled binaries for your platform from [GitHub Releases](https://github.com/uphy/agent-sync/releases).

The release process automatically builds binaries for multiple platforms:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64, arm64)

Simply download the appropriate archive for your platform, extract it, and place the binary in a directory included in your system's PATH.

## Verifying the Installation

After installation, verify that agent-sync is installed correctly by running:

```bash
agent-sync --version
```

You should see the version number of the installed agent-sync tool.

## Next Steps

Once installed, you can create a new agent-sync project using the `init` command:

```bash
mkdir my-ai-instructions
cd my-ai-instructions
agent-sync init
```

For more information about using agent-sync, see the [Command Line Interface](cli.md) documentation.

## Navigation

- [Main Configuration Guide](config.md)
- [Configuration Reference](config-reference.md)
- [Command Line Interface](cli.md)
- [Task Types](task-types.md)
- [Template System](templates.md)
- [Input and Output Processing](input-output.md)
- [Troubleshooting](troubleshooting.md)
- [Examples and Best Practices](examples.md)

---

| Previous | Next |
|----------|------|
| [Index](index.md) | [Configuration Guide](config.md) |

[Back to index](index.md)