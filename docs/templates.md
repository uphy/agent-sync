---
layout: default
title: Template System
---

[Back to index](index.md)

# Template System

agent-sync provides a powerful templating system that allows you to create dynamic content, include files, and handle agent-specific formatting.

## Template Path Resolution

Within templates, paths are resolved using special rules:

1. Absolute paths (like `/root` on Unix or `C:\` on Windows) are not allowed and will return an error

2. Paths starting with `@/` are relative to the configuration file directory (BasePath) with the `@/` prefix removed
   - Example: `@/shared/template.md` → `<config-dir>/shared/template.md`
   
3. All other paths (including those with `./` or `../` prefix) are relative to the including file's directory
   - Example from file `/templates/main.md`: `./partial.md` → `/templates/partial.md`
   - Example from file `/templates/main.md`: `../shared/partial.md` → `/shared/partial.md`
   - Example from file `/templates/main.md`: `common.md` → `/templates/common.md`

Examples assuming config file is at `/project/agent-sync.yml`:

| Template Path | Including File | Resolved Path |
|---------------|----------------|---------------|
| `@/shared/template.md` | Any file | `/project/shared/template.md` |
| `./partial.md` | `/project/templates/main.md` | `/project/templates/partial.md` |
| `../shared/common.md` | `/project/templates/main.md` | `/project/shared/common.md` |
| `utils/helper.md` | `/project/templates/main.md` | `/project/templates/utils/helper.md` |
| `/absolute/path.md` | Any file | Error: absolute paths are not allowed |

## Glob Pattern Support

All path-based template functions (`include`, `includeRaw`, `reference`, and `referenceRaw`) support glob patterns for dynamic file selection. This allows you to include multiple files that match a specified pattern without having to list them individually.

For example:
{% raw %}
```
{{ includeRaw "docs/*.md" }}           // All .md files in the docs directory
{{ include "templates/**/*.md" }}      // All .md files in templates directory and its subdirectories
{{ reference "configs/{dev,prod}.json" }} // Both dev.json and prod.json from the configs directory
```
{% endraw %}

For detailed information about glob pattern syntax and behavior, see the [Glob Patterns](glob-patterns.md) documentation.

## Template Functions

agent-sync supports the following template functions in source files:

| Function | Description | Example |
|----------|-------------|---------|
| `file "path/to/file"` | Formats a file reference according to the output agent | {% raw %}`{{ file "src/main.go" }}`{% endraw %} → `` `src/main.go` `` (Copilot) or `@/src/main.go` (Cline) |
| `include "path/to/file" ["path/to/another/file" ...]` | Includes content from one or more files with template processing. Supports glob patterns. | {% raw %}`{{ include "common/header.md" }}`{% endraw %} or {% raw %}`{{ include "common/*.md" }}`{% endraw %} |
| `includeRaw "path/to/file" ["path/to/another/file" ...]` | Includes content from one or more files without template processing. Supports glob patterns. | {% raw %}`{{ includeRaw "common/header.md" }}`{% endraw %} or {% raw %}`{{ includeRaw "templates/**/*.md" }}`{% endraw %} |
| `reference "path/to/file" ["path/to/another/file" ...]` | References content from one or more files with template processing. Supports glob patterns. | {% raw %}`{{ reference "data/config.json" }}`{% endraw %} or {% raw %}`{{ reference "config/*.json" }}`{% endraw %} |
| `referenceRaw "path/to/file" ["path/to/another/file" ...]` | References content from one or more files without template processing. Supports glob patterns. | {% raw %}`{{ referenceRaw "data/config.json" }}`{% endraw %} or {% raw %}`{{ referenceRaw "**/*.json" }}`{% endraw %} |
| `mcp "agent" "command" "arg1" "arg2"` | Formats an MCP command for the output agent | {% raw %}`{{ mcp "github" "get-issue" "owner" "repo" "123" }}`{% endraw %} |
| `agent` | Returns the current output agent identifier | {% raw %}`{{ if eq agent "claude" }}Claude-specific content{{ end }}`{% endraw %} |
| `ifAGENT "content"` | Conditionally includes content only for the specified agent | {% raw %}`{{ ifRoo "This will only appear in Roo output" }}`{% endraw %} |

## Template Function Examples

**File references:**
{% raw %}
```
For Claude: {{ file "src/main.go" }} → @src/main.go
For Roo: {{ file "src/main.go" }} → @/src/main.go
For Cline: {{ file "src/main.go" }} → @/src/main.go
For Copilot: {{ file "src/main.go" }} → `src/main.go`
```
{% endraw %}

**Including templates:**
{% raw %}
```
{{ include "header.md" }}
```
{% endraw %}
This will include and process the content of `header.md`, executing any template directives it contains.

**Including multiple templates with glob pattern:**
{% raw %}
```
{{ include "sections/*.md" }}
```
{% endraw %}
This will include and process all `.md` files in the `sections` directory in alphabetical order, executing any template directives they contain.

**Including multiple templates:**
{% raw %}
```
{{ include "header.md" "content.md" "footer.md" }}
```
{% endraw %}
This will include and process the content of all specified files in the given order, executing any template directives they contain.

**Including content without template processing:**
{% raw %}
```
{{ includeRaw "header.md" }}
```
{% endraw %}
This will include the content of `header.md` exactly as it is, without executing any template directives it contains.

**Including multiple files without template processing using glob:**
{% raw %}
```
{{ includeRaw "docs/**/*.md" }}
```
{% endraw %}
This will include the content of all `.md` files in the `docs` directory and its subdirectories, without executing any template directives they contain.

**Including files with exclusion pattern:**
{% raw %}
```
{{ include "src/*.go" "!src/*_test.go" }}
```
{% endraw %}
This will include all `.go` files in the `src` directory, excluding test files that match the `*_test.go` pattern.

**Referencing files with template processing:**
{% raw %}
```
{{ reference "data.json" }}
```
{% endraw %}
This will include the content of `data.json` with template processing.

**Referencing multiple files with template processing using glob:**
{% raw %}
```
{{ reference "configs/{production,staging,development}.json" }}
```
{% endraw %}
This will include the content of `production.json`, `staging.json`, and `development.json` from the `configs` directory, with template processing applied to each file.

**Referencing files without template processing:**
{% raw %}
```
{{ referenceRaw "data.json" }}
```
{% endraw %}
This will include the raw content of `data.json` without any template processing.

**Referencing multiple files without template processing:**
{% raw %}
```
{{ referenceRaw "config.json" "settings.json" "defaults.json" }}
```
{% endraw %}
This will include the raw content of all specified files in the given order, without any template processing.

**MCP commands:**
{% raw %}
```
{{ mcp "mcp_server_name" "tool_name" "arguments" }}
```
{% endraw %}
This formats an MCP command for the output agent.

**Agent-specific content:**
{% raw %}
```
{{ if eq agent "claude" }}
Claude-specific content here
{{ else if eq agent "roo" }}
Roo-specific content here
{{ end }}
```
{% endraw %}

**Simplified agent-specific content:**
{% raw %}
```
{{ ifRoo "This content only appears in Roo output" }}
{{ ifClaude "This content only appears in Claude output" }}
```
{% endraw %}

## Navigation

- [Main Configuration Guide](config.md)
- [Configuration Reference](config-reference.md)
- [Input and Output Processing](input-output.md)
- [Task Types](task-types.md)
- [Glob Patterns](glob-patterns.md)
- [Command Line Interface](cli.md)
- [Troubleshooting](troubleshooting.md)
- [Examples and Best Practices](examples.md)

---

| Previous | Next |
|----------|------|
| [Command Line Interface](cli.md) | [Input and Output Processing](input-output.md) |

[Back to index](index.md)