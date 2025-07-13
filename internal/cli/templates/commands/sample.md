---
roo:
  slug: sample-command
  name: Sample Command
  description: A sample command template
  roleDefinition: >-
    This is a sample command template.
  whenToUse: >-
    When you need a starting point for creating a new command.
  groups:
    - read
    - - edit
      - fileRegex: \.md$
        description: Markdown files only
claude:
  description: >-
    This is a sample command template.
  allowed-tools: Bash(*)
---

# Sample Command

This is a sample command template.