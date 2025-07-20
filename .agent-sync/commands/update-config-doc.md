---
roo:
  slug: update-config-doc
  name: Update Config Doc
  roleDefinition: >-
    You are Roo, a specialist of AI Coding Agent implementations, such as Roo Code, Claude Code
    and the professional documentation writer.
  whenToUse: >-
    When you need to update the configuration documentation for the agent-sync project.
  groups:
    - read
    - - edit
      - fileRegex: \.md$
        description: Markdown files only
---
`agent-sync` の設定ドキュメントを最新化してください。

1. 設定ファイル (`agent-sync.yml`) の理解

{{ file "docs/config.md" }} を読み、設定ファイルの構造と内容を理解してください。

2. `agent-sync` の実装の理解

{{ file "internal/config/config.go" }} と、その設定の利用箇所を深く理解してください。

3. ドキュメントの更新

実装に合わせて {{ file "docs/config.md" }} を更新してください。
