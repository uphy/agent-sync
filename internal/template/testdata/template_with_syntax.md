# Template with Syntax

This file contains template syntax that should not be processed by raw functions:

{{.Value}} - This is a variable
{{include "testdata/content.md"}} - This is an include directive
{{agent}} - This is a function call

The raw functions should preserve these as-is.