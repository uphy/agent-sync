# Template System

agent-def provides a powerful templating system that allows you to create dynamic content, include files, and handle agent-specific formatting.

## Template Path Resolution

Within templates, paths are resolved using special rules:

1. Absolute paths (like `/root` on Unix or `C:\` on Windows) are not allowed and will return an error

2. Paths starting with `@/` are relative to the configuration file directory (BasePath) with the `@/` prefix removed
   - Example: `@/shared/template.md` → `<config-dir>/shared/template.md`
   
3. All other paths (including those with `./` or `../` prefix) are relative to the including file's directory
   - Example from file `/templates/main.md`: `./partial.md` → `/templates/partial.md`
   - Example from file `/templates/main.md`: `../shared/partial.md` → `/shared/partial.md`
   - Example from file `/templates/main.md`: `common.md` → `/templates/common.md`

Examples assuming config file is at `/project/agent-def.yml`:

| Template Path | Including File | Resolved Path |
|---------------|----------------|---------------|
| `@/shared/template.md` | Any file | `/project/shared/template.md` |
| `./partial.md` | `/project/templates/main.md` | `/project/templates/partial.md` |
| `../shared/common.md` | `/project/templates/main.md` | `/project/shared/common.md` |
| `utils/helper.md` | `/project/templates/main.md` | `/project/templates/utils/helper.md` |
| `/absolute/path.md` | Any file | Error: absolute paths are not allowed |

## Template Functions

agent-def supports the following template functions in source files:

| Function | Description | Example |
|----------|-------------|---------|
| `file "path/to/file"` | Formats a file reference according to the output agent | `{{ file "src/main.go" }}` → `` `src/main.go` `` (Copilot) or `@/src/main.go` (Cline) |
| `include "path/to/file" ["path/to/another/file" ...]` | Includes content from one or more files with template processing | `{{ include "common/header.md" }}` or `{{ include "header.md" "footer.md" }}` |
| `includeRaw "path/to/file" ["path/to/another/file" ...]` | Includes content from one or more files without template processing | `{{ includeRaw "common/header.md" }}` or `{{ includeRaw "header.md" "footer.md" }}` |
| `reference "path/to/file" ["path/to/another/file" ...]` | References content from one or more files with template processing | `{{ reference "data/config.json" }}` or `{{ reference "config.json" "settings.json" }}` |
| `referenceRaw "path/to/file" ["path/to/another/file" ...]` | References content from one or more files without template processing | `{{ referenceRaw "data/config.json" }}` or `{{ referenceRaw "config.json" "settings.json" }}` |
| `mcp "agent" "command" "arg1" "arg2"` | Formats an MCP command for the output agent | `{{ mcp "github" "get-issue" "owner" "repo" "123" }}` |
| `agent` | Returns the current output agent identifier | `{{ if eq agent "claude" }}Claude-specific content{{ end }}` |
| `ifAGENT "content"` | Conditionally includes content only for the specified agent | `{{ ifRoo "This will only appear in Roo output" }}` |

## Template Function Examples

**File references:**
```
For Claude: {{ file "src/main.go" }} → @src/main.go
For Roo: {{ file "src/main.go" }} → @/src/main.go
For Cline: {{ file "src/main.go" }} → @/src/main.go
For Copilot: {{ file "src/main.go" }} → `src/main.go`
```

**Including templates:**
```
{{ include "header.md" }}
```
This will include and process the content of `header.md`, executing any template directives it contains.

**Including multiple templates:**
```
{{ include "header.md" "content.md" "footer.md" }}
```
This will include and process the content of all specified files in the given order, executing any template directives they contain.

**Including content without template processing:**
```
{{ includeRaw "header.md" }}
```
This will include the content of `header.md` exactly as it is, without executing any template directives it contains.

**Including multiple files without template processing:**
```
{{ includeRaw "header.md" "content.md" "footer.md" }}
```
This will include the content of all specified files in the given order, without executing any template directives they contain.

**Referencing files with template processing:**
```
{{ reference "data.json" }}
```
This will include the content of `data.json` with template processing.

**Referencing multiple files with template processing:**
```
{{ reference "config.json" "settings.json" "defaults.json" }}
```
This will include the content of all specified files in the given order, with template processing applied to each file.

**Referencing files without template processing:**
```
{{ referenceRaw "data.json" }}
```
This will include the raw content of `data.json` without any template processing.

**Referencing multiple files without template processing:**
```
{{ referenceRaw "config.json" "settings.json" "defaults.json" }}
```
This will include the raw content of all specified files in the given order, without any template processing.

**MCP commands:**
```
{{ mcp "mcp_server_name" "tool_name" "arguments" }}
```
This formats an MCP command for the output agent.

**Agent-specific content:**
```
{{ if eq agent "claude" }}
Claude-specific content here
{{ else if eq agent "roo" }}
Roo-specific content here
{{ end }}
```

**Simplified agent-specific content:**
```
{{ ifRoo "This content only appears in Roo output" }}
{{ ifClaude "This content only appears in Claude output" }}
```

## Navigation

- [Main Configuration Guide](config.md)
- [Configuration Reference](config-reference.md)
- [Input and Output Processing](input-output.md)
- [Task Types](task-types.md)
- [Command Line Interface](cli.md)
- [Troubleshooting](troubleshooting.md)
- [Examples and Best Practices](examples.md)