# Template Function Tests

## File Function Tests
- `file(relative)`: @/path/to/file
- `file(absolute)`: /path/to/file
- `file(parent)`: @/../path/to/file

## Include Function Tests

-- include_test.md --
- agent: roo
-- end of include_test.md --

## Reference Function Tests

[参考: reference_test.md]

## MCP Function Tests

- `mcp(basic)`: MCP tool (Server: `test-server`, Tool: `get-data`)
- `mcp(with args)`: MCP tool (Server: `test-server`, Tool: `get-data`, Arguments: `arg1, arg2`)

## Roo Only

This is roo only section.

## Variables

- `Current agent`: roo

## References

### reference_test.md

-- reference_test.md --
- agent: roo
-- end of reference_test.md --

