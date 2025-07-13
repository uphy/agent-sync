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
---

# Deploy Command

This command is used to deploy the application.
