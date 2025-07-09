---
title: Example Memory Context
description: An example memory context file demonstrating includes, references, and MCP commands
tags:
  - example
  - memory
---

# Example Memory Context

This is an example memory context file demonstrating includes, references, and MCP commands.

## Section with Include

{{include "included_content.md"}}

## Section with Reference

{{reference "referenced_content.md"}}

## Section with MCP Command

{{mcp "example-agent" "do-something" "paramA" "paramB"}}