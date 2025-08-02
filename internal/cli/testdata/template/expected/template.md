# Template Function Tests

## File Function Tests
- `file(relative)`: {{ file "path/to/file" }}
- `file(absolute)`: {{ file "/path/to/file" }}
- `file(parent)`: {{ file "../path/to/file" }}

## Include Function Tests

{{ include "include_test.md" }}

## Reference Function Tests

{{ reference "reference_test.md" }}

## MCP Function Tests

- `mcp(basic)`: {{ mcp "test-server" "get-data" }}
- `mcp(with args)`: {{ mcp "test-server" "get-data" "arg1" "arg2" }}

{{- if isRoo }}

## Roo Only

This is roo only section.
{{- end }}
{{- if isClaude }}

## Claude Only

This is claude only section.
{{- end }}

## Variables

- `Current agent`: {{ agent }}