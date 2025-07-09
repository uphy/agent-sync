---
name: test-command
description: A test command for testing
parameters:
  - name: param1
    type: string
    description: First parameter
    required: true
  - name: param2
    type: integer
    description: Second parameter
    required: false
    default: 42
examples:
  - "test-command --param1 value"
  - "test-command --param1 value --param2 100"
---

# Test Command

This is a test command file with YAML frontmatter.

## Usage

Use this command to test the command processing functionality.

```
test-command --param1 <value> [--param2 <number>]
```

## Examples

```
test-command --param1 value
test-command --param1 value --param2 100