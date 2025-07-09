## 1\. 概要

### 1.1. 解決する課題

AIコーディングエージェントは、プロジェクトのコンテキストを正確に理解することで真価を発揮するが、各エージェントは独自のコンテキスト形式（`CLAUDE.md`, `.roo/rules`等）を採用しており、開発者は設定ファイルの重複作成やメンテナンスコストの増大という課題に直面しています。

また、`git worktree` などを用いて複数のワークスペースで同時並行に作業する際、バージョン管理されていないコンテキストファイルやコマンド定義を各環境に手動で配置・同期するのは手間がかかり、更新漏れの原因にもなります。

### 1.2. 提案する解決策

この課題を解決するため、**コンテキスト定義とコマンド定義を分離**し、それぞれを\*\*特定のAIエージェント形式に変換する、UNIX哲学に基づいたCLIツール「Agent Definition (`agent-def`)」\*\*を提案する。

このツールは以下の思想に基づいている。

- **関心の分離 (Separation of Concerns)**: プロジェクトの「永続的な記憶」となるコンテキストは任意のMarkdownファイルに、「特定のタスクを実行する」コマンドは**コマンドごとに独立したファイル**に記述する。
- **抽象化とテンプレート化 (Abstraction & Templating)**: `{{ file "path" }}` や `{{ include "path" }}` のような抽象的なテンプレート構文を導入する。これにより、エージェントごとのフォーマットの違いをツールが吸収し、利用者は本質的な定義に集中できる。
- **モジュール性 (Modularity)**: 各コマンドは「入力ファイル/ディレクトリを受け取り、変換し、出力する」という単一の責務を持つ。これにより、シェルのパイプラインやスクリプトとの連携が容易になる。
- **柔軟なファイル構成**: `init` コマンドを廃止し、特定のファイル構成を強制しない。利用者は自由に定義ファイルを配置できる。

## 2\. テンプレートエンジンとヘルパー関数

`agent-def` は、入力ファイルを処理する際にテンプレートエンジンを使用する。

| ヘルパー関数 | 説明 |
| --- | --- |
| `{{ file "path/to/file" }}` | ファイルやディレクトリへの参照を、ターゲットのエージェントに最適化された形式で出力する。 |
| `{{ include "path/to/file" }}` | 指定されたファイルの内容を読み込み、その場に展開（インライン展開）する。**展開される内容も再帰的にテンプレート処理される。** |
| `{{ reference "path/to/file" }}` | その場には参照（例: `[参考: path/to/file]`）のみを記述し、ドキュメントの末尾にファイル内容を追記する。長大なファイル向け。**追記される内容も再帰的にテンプレート処理される。** |
| `{{ mcp "agent" "command" "args..." }}` | 他のエージェントやツールを呼び出すためのコマンドを生成する。エージェント間連携（Multi-agent Collaboration Platform）用。 |

## 3\. 共通フォーマット仕様

### 3.1. コンテキスト定義ファイル (例: `context.md`)

プロジェクトの概要、規約、アーキテクチャなどを記述する。ファイル名は任意。

```
<!-- context.md -->

# プロジェクト名: My Awesome Web App

## 1. プロジェクト概要
...

## 2. コーディング規約
...

## 3. 主要アーキテクチャ
...

### データベーススキーマ詳細
{{ include "prisma/schema.prisma" }}
```

### 3.2. コマンド定義ファイル (例: `commands/commit.md`)

**1コマンドにつき1ファイル**で定義する。ファイル名はコマンド名と一致させることが推奨される。`agent-def` ツールは、frontmatterに記述されたキーを、ターゲットエージェントの仕様に合わせてマッピングする。

```
<!-- commands/commit.md -->
---
# Frontmatterでコマンドのメタデータを定義
# これらのキーは agent-def によって各エージェントの形式に変換される
slug: "commit"
name: "📝 Commit"
roleDefinition: "現在のgitのステージング差分を分析し、変更理由を尋ね、セマンティックなコミットメッセージを生成し、git commitを実行します。"
whenToUse: "ステージングされた変更に基づいてgit commitを自動生成して実行したい場合に使用します。"
groups:
  - "read"
  - "command"
  - "mcp"
---
# Markdown本体にプロンプトを記述
ステージングされているファイルの差分を分析し、Conventional Commits規約に準拠したコミットメッセージを生成してください。
```

## 4\. CLIツール仕様 (`agent-def`)

### 4.1. コマンド一覧

| コマンド | 説明 |
| --- | --- |
| `agent-def memory` | コンテキスト定義ファイルをテンプレート処理し、各エージェント形式に変換する。 |
| `agent-def command` | **コマンド定義ファイル群**をテンプレート処理し、各エージェント形式に集約・変換する。 |
| `agent-def list` | 対応しているエージェントの一覧を表示する。 |

### 4.2. 使用例

#### コンテキスト情報の生成

```
agent-def memory --type claude --input context.md --output CLAUDE.md
```

#### コマンド情報の生成

`commands` ディレクトリ配下の全 `.md` ファイルを読み込み、Roo Code用の単一の設定ファイルに集約して出力する。

```
agent-def command --type roo --input-dir ./commands --output .roo/custom_modes.json
```

### 4.3. プロジェクトでの運用

`package.json` の `scripts` に登録する運用方法は引き続き推奨される。

```
{
  "name": "my-awesome-web-app",
  "scripts": {
    "agent:def": "npm run agent:def:claude && npm run agent:def:roo",
    "agent:def:claude": "agent-def memory -t claude -i context.md -o CLAUDE.md",
    "agent:def:roo": "agent-def memory -t roo -i context.md -o .roo/rules && agent-def command -t roo -i ./commands -o .roo/custom_modes.json"
  },
  "devDependencies": {
    "agent-def": "latest" 
  }
}
```

## 5\. 変換ルール (Conversion Rules)

`agent-def` は、ターゲットとなるエージェントに応じてヘルパー関数を以下のルールで変換します。

### 5.1. `{{ file "path/to/file" }}`

ファイルやディレクトリへの参照を、各エージェントが解釈しやすい形式に変換します。

- **Roo Code**: `@/path/to/file`
- **Claude Code**: `@path/to/file`

### 5.2. `{{ include "path/to/file" }}`

このヘルパーの出力形式にエージェント間の差異はありません。指定されたファイルの内容がそのまま展開されます。

### 5.3. `{{ mcp "agent" "command" "args..." }}`

このヘルパーの出力形式にエージェント間の差異はありません。以下の固定フォーマットで出力されます。

`MCP tool (MCP Server: <agent>, Tool: <command>, Arguments: <args...>)`

## 6\. 今後の展望

- **ヘルパー関数の拡充**:
	- `{{ ls "dir" }}`: ディレクトリ内のファイル一覧を出力する。
	- `{{ git_diff "commit" }}`: 特定のコミットの差分をコンテキストに含める。
- **対応エージェントの拡充**: GitHub Copilot Workspace, Cursor などに順次対応。
- **VS Code拡張機能**: テンプレート構文のシンタックスハイライトや入力補完を提供。