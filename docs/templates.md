# Template System

agent-def provides a powerful templating system that allows you to create dynamic content, include files, and handle agent-specific formatting.

## Template Path Resolution

Within templates, paths are resolved using special rules:

1. Paths starting with `/` are relative to the configuration file directory with the leading slash removed
   - Example: `/shared/template.md` → `<config-dir>/shared/template.md`
   
2. Paths starting with `./` or `../` are relative to the including file's directory
   - Example from file `/templates/main.md`: `./partial.md` → `/templates/partial.md`
   - Example from file `/templates/main.md`: `../shared/partial.md` → `/shared/partial.md`
   
3. Other paths without `./` or `../` prefix are relative to the configuration file directory
   - Example: `shared/template.md` → `<config-dir>/shared/template.md`
   
4. OS-absolute paths (like C:\ on Windows) are preserved as-is
   - Example: `C:\templates\file.md` remains unchanged

Examples assuming config file is at `/project/agent-def.yml`:

| Template Path | Including File | Resolved Path |
|---------------|----------------|---------------|
| `/shared/template.md` | Any file | `/project/shared/template.md` |
| `./partial.md` | `/project/templates/main.md` | `/project/templates/partial.md` |
| `../shared/common.md` | `/project/templates/main.md` | `/project/shared/common.md` |
| `utils/helper.md` | Any file | `/project/utils/helper.md` |
| `C:\absolute\path.md` | Any file (Windows) | `C:\absolute\path.md` |

## Template Functions

agent-def supports the following template functions in source files:

| Function | Description | Example |
|----------|-------------|---------|
| `file "path/to/file"` | Formats a file reference according to the output agent | `{{ file "src/main.go" }}` → `` `src/main.go` `` (Copilot) or `@/src/main.go` (Cline) |
| `include "path/to/file"` | Includes content from another file with template processing | `{{ include "common/header.md" }}` |
| `includeRaw "path/to/file"` | Includes content from another file without template processing | `{{ includeRaw "common/header.md" }}` |
| `reference "path/to/file"` | References another file's content with template processing | `{{ reference "data/config.json" }}` |
| `referenceRaw "path/to/file"` | References another file's content without template processing | `{{ referenceRaw "data/config.json" }}` |
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

**Including content without template processing:**
```
{{ includeRaw "header.md" }}
```
This will include the content of `header.md` exactly as it is, without executing any template directives it contains.

**Referencing files with template processing:**
```
{{ reference "data.json" }}
```
This will include the content of `data.json` with template processing.

**Referencing files without template processing:**
```
{{ referenceRaw "data.json" }}
```
This will include the raw content of `data.json` without any template processing.

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