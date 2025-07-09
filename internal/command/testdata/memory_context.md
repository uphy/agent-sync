---
title: Test Memory Context
description: A test memory context file for testing
tags:
  - test
  - memory
---

# Test Memory Context

This is a test memory context file with YAML frontmatter.

## Section with Include

{{include "testdata/included_content.md"}}

## Section with Reference

{{reference "testdata/referenced_content.md"}}

## Section with MCP Command

{{mcp "test-agent" "test-command" "arg1" "arg2"}}