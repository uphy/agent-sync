---
roo:
  slug: deploy
  name: Deploy
  description: deploy app
  roleDefinition: >-
    This command is used to deploy the application.
  whenToUse: >-
    When you need to deploy the application
  groups:
    - read
    - - edit
      - fileRegex: \.md$
        description: Markdown files only
claude:
  description: >-
    This command is used to deploy the application.
  allowed-tools: Bash(git add:*), Bash(git status:*), Bash(git commit:*)
---

# Deploy Command

This command is used to deploy the application.
