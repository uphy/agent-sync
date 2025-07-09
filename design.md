## 1\. 概要

### 1.1. 解決する課題

AIコーディングエージェントは、プロジェクトのコンテキスト（アーキテクチャ、コーディング規約、主要なファイル構造など）を正確に理解することで、その真価を発揮する。しかし、各エージェントは独自のコンテキスト形式（例: `CLAUDE.md`, `.roo/rules/`, `.aider.conf.yml`）を採用しており、開発者はプロジェクトごとに類似した内容の設定ファイルを複数作成・管理する必要がある。この状況は以下の問題点を引き起こす。

- **設定の重複と不整合**: 同じコンテキスト情報を、エージェントごとに異なるフォーマットで記述する必要があり、メンテナンスコストが増大し、更新漏れによる情報の不整合が発生しやすい。
- **ワークフローの複雑化**: `git worktree` などを用いて複数のブランチで並行作業を行う際、バージョン管理下にないコンテキストファイルやコマンド定義を各ワークツリーに手動でコピー＆ペーストする必要があり、非常に手間がかかる。
- **チーム開発での共有の困難**: プロジェクト固有の便利なカスタムコマンドや重要なコンテキスト情報を、口頭やドキュメントベースで共有するのは非効率であり、チームメンバー間での活用度に差が生まれる。

### 1.2. 提案する解決策

この課題を解決するため、\*\*単一の宣言的な設定ファイル (`agent-def.yml`) に基づいて、プロジェクト固有の設定とOSユーザー全体の設定を一元管理し、サポートする各AIエージェントの形式に一括変換するCLIツール「Agent Definition (`agent-def`)」\*\*を提案する。

このツールは以下の思想に基づいている。

- **Single Source of Truth**: プロジェクトのコンテキストやカスタムコマンドは、すべてMarkdown形式のソースファイルとして管理する。`agent-def.yml` は、それらのソースを「どのエージェント」の「どの形式」に変換するかを定義する単一の真実点となる。
- **宣言的な設定 (Declarative Configuration)**: 「何を」「どこから」「どこへ」生成するかを、`agent-def.yml` で宣言的に管理する。開発者は「どのように」変換するかを意識する必要はない。
- **スコープの分離**: 設定を `projects`（プロジェクト単位）と `user`（OSユーザー単位）に明確に分離し、異なるスコメインテナンスのスコープを明確に管理する。
- **抽象化とテンプレート化 (Abstraction & Templating)**: `{{ include "path" }}` のようなヘルパー関数を持つMarkdownをソースとすることで、定義の再利用性と表現力を高め、DRY (Don't Repeat Yourself) の原則を徹底する。
- **複雑なワークフローへの標準対応**: `git worktree` で作成された複数の作業ディレクトリへの設定ファイルの同時展開や、柔軟な出力先指定に標準で対応する。

## 2\. 設定ファイル仕様 (`agent-def.yml`)

本ツールの中心となるYAML形式の設定ファイル。プロジェクトルート、またはホームディレクトリ (`~/.config/agent-def/`) に配置する。

### 2.1. トップレベルの構造

```
# agent-def.yml

# 設定ファイルのスキーマバージョン。将来的な破壊的変更に対応するため。
configVersion: "1.0"

# プロジェクト固有の設定を定義します。
projects:
  # ... projectオブジェクトのリスト ...

# OSユーザー全体で共通のグローバル設定を定義します。
user:
  # ... userオブジェクト ...
```

### 2.2. `projects` ブロック

複数のプロジェクト設定を管理する。キーには任意のプロジェクト名を指定する。

```
projects:
  # "my-awesome-app" はこの設定ブロックを識別するための任意の名前。
  my-awesome-app:
    # [任意] ソースパスの基点となるディレクトリ。
    # 省略した場合は \`agent-def.yml\` があるディレクトリが基準となる。
    # 絶対パス、または設定ファイルからの相対パスで指定。
    root: ./
    
    # [必須] 生成されたファイルの配置先となるプロジェクトのルートディレクトリ。
    # git worktreeなどを想定し、複数指定可能。
    # 安全性のため、パスは必ず絶対パス、または \`~\` から始まるパスで指定すること。
    # 環境変数も利用可能 (例: '$PROJECT_HOME/my-awesome-app')。
    destinations:
      - /Users/your-name/src/my-awesome-app
      - /Users/your-name/src/my-awesome-app-feature-x
    
    # [必須] このプロジェクトで実行される生成タスクのリスト。
    tasks:
      # ... taskオブジェクトのリスト ...
```

### 2.3. `user` ブロック

OSユーザー全体に適用されるグローバルな設定を管理する。プロジェクトに依存しない共通のコマンドやコンテキストを定義するのに便利です。出力先はユーザーのホームディレクトリ (`~`) が基準となる。

```
user:
  # [必須] グローバルに実行される生成タスクのリスト。
  tasks:
    # ... taskオブジェクトのリスト ...
```

### 2.4. `task` オブジェクト

`projects` と `user` の両方で使われる、単一の生成タスクを定義するオブジェクト。

| キー | 型 | 必須 | 説明 |
| --- | --- | --- | --- |
| `name` | string | 否 | タスクを識別するための任意の名前。ログ出力などで利用され、可読性が向上する。 |
| `type` | string | 済 | タスクの種類。`command` (カスタムコマンド) または `memory` (コンテキスト情報) を指定。 |
| `sources` | list | 済 | ソースとなるファイルやディレクトリのパスのリスト。パスは `projects.root` (または `agent-def.yml` の場所) からの相対パス。ディレクトリを指定した場合、その配下の `.md` ファイルがすべて対象となる。ワイルドカード (`**/*`) も利用可能。 |
| `concat` | boolean | 否 | `true` の場合、`sources` で指定された複数のMarkdownファイルを1つのファイルに結合して出力する。`false` または省略した場合、ソースファイルごとに個別に出力する。主に `memory` タイプで利用。デフォルトは `false`。 |
| `targets` | list | 済 | 1つ以上の出力先エージェントを定義する `target` オブジェクトのリスト。 |

**`task` オブジェクトの例:**

```
# 複数のコンテキストファイルを1つに結合してClaudeとRooに出力するタスク
- name: "Project Architecture Memory"
  type: memory
  sources:
    - context/
  concat: true # context/配下の.mdを1つにまとめる
  targets:
    - agent: claude
    - agent: roo
      # デフォルトの出力先 \`.roo/rules/\` を上書き
      target: ".roo/memory/architecture.md"

# コマンド定義ディレクトリから、個別のコマンドファイルを生成するタスク
- name: "Project Custom Commands"
  type: command
  sources:
    - commands/ # このディレクトリ配下の.mdがそれぞれコマンドになる
  # concatはfalse(デフォルト)なので、ファイルごとに生成される
  targets:
    - agent: roo
    - agent: claude
```

### 2.5. `target` オブジェクト

`tasks` の中で、具体的な出力先エージェントとその設定を定義する。

| キー | 型 | 必須 | 説明 |
| --- | --- | --- | --- |
| `agent` | string | 済 | ターゲットとなるエージェント名 (例: `roo`, `claude`, `cursor`)。`agent-def list` で対応一覧を確認できる。 |
| `target` | string | 否 | 出力先のパス。省略した場合、[6.2. 出力先の変換](https://www.google.com/search?q=%2362-%E5%87%BA%E5%8A%9B%E5%85%88%E3%81%AE%E5%A4%89%E6%8F%9B "null")で定義された規約ベースのデフォルトパスに出力される。<br>・`projects` タスクの場合: `destinations` からの相対パス。<br>・`user` タスクの場合: `~` (ホームディレクトリ) からの相対パス。 |

## 3\. 共通フォーマット仕様（ソースファイル）

`agent-def.yml` の `sources` で指定されるMarkdownファイルの仕様。

### 3.1. コンテキスト定義ファイル (例: `context/architecture.md`)

プロジェクトの概要、規約、アーキテクチャ、ファイル構造など、AIエージェントに記憶させたい情報を自由に記述するMarkdownファイル。特別な形式はない。

```
# My Awesome App アーキテクチャ概要

このアプリケーションは、フロントエンドとバックエンドから構成されるマイクロサービスアーキテクチャを採用しています。

## フロントエンド

-   **フレームワーク**: React (Next.js)
-   **言語**: TypeScript
-   **主要ライブラリ**:
    -   \`zustand\`: 状態管理
    -   \`react-query\`: データフェッチ
    -   \`tailwindcss\`: スタイリング

## コーディング規約

-   コンポーネント名は \`PascalCase\` で記述してください。
-   \`src/components/ui\` ディレクトリには、汎用的なUIコンポーネントを配置します。{{ file "src/components/ui/Button.tsx" }} を参考にしてください。
```

### 3.2. コマンド定義ファイル (例: `commands/create_component.md`)

1コマンドにつき1ファイルで定義する。YAML frontmatterでメタデータを、Markdown本体でプロンプトを記述する。

- **ファイル名**: 任意。管理しやすい名前をつける。
- **Frontmatter**: コマンドのメタデータを定義する。
	- `command` (string, 必須): 実際に呼び出す際のコマンド名。
	- `description` (string, 任意): コマンドの説明。
- **Markdown本体**: AIエージェントに渡すプロンプト本体。テンプレート構文を利用できる。
```
---
command: create_component
description: "新しいReactコンポーネントを生成する"
---

\`src/components/{{ args[0] }}/{{ args[1] }}.tsx\` というパスに、新しいReactコンポーネントを作成してください。

コンポーネント名は \`{{ args[1] }}\` です。
ファイルには、基本的なコンポーネントの雛形と、対応するStorybookのファイルも含めてください。

**制約:**
-   関数コンポーネントとして実装してください。
-   \`tailwindcss\` を使用して基本的なスタイルを適用してください。
-   以下のファイル構造を参考にしてください。
    {{ include "templates/component_structure.md" }}
```

## 4\. テンプレートエンジンとヘルパー関数

ソースファイルを処理する際に、以下のヘルパー関数を `{{ }}` で囲んで使用できる。

| ヘルパー関数 | 説明 |
| --- | --- |
| `{{ file "path" }}` | ファイルパスをエージェント固有の参照形式に変換する（例: `@/path`）。パスは `root` からの相対パスとして解釈される。 |
| `{{ include "path" }}` | 指定された `path` のファイル内容をその場にインライン展開する。展開された内容も再帰的にテンプレート処理される。無限再帰を避けるため、最大深度が設定される。 |
| `{{ reference "path" }}` | 指定された `path` の内容を、エージェントが「参照情報」として認識できる形式で追記する。例えば、プロンプトの最後に `See Also:` セクションを追加するなど。 |
| `{{ mcp "server" "tool" "args..." }}` | 他のエージェントやツールを呼び出すためのコマンドを生成する（Model Context Protocol用）。エージェント間の連携を目的とする。 |
| `{{ env "VAR_NAME" }}` | 指定された環境変数の値を展開する。変数が存在しない場合は空文字列になる。 |
| `{{ args[n] }}` | (コマンド定義ファイル内のみ) コマンド実行時に渡された引数を参照する。`{{ args[0] }}` は最初の引数。 |

## 5\. CLIツール仕様 (`agent-def`)

### 5.1. コマンド一覧

| コマンド | 説明 |
| --- | --- |
| `agent-def build [project...]` | `agent-def.yml` に基づき、ファイルを生成する。<br>・`[project...]`: 対象のプロジェクト名を指定。省略した場合は `projects` 内の全プロジェクトが対象となる。<br>・`--user`: `user` 設定のみをビルドする。<br>・`--watch`: ソースファイルの変更を監視し、自動で再ビルドする。<br>・`--dry-run`: ファイルを実際に書き込まず、生成される内容と出力先パスをコンソールに表示する。 |
| `agent-def validate` | `agent-def.yml` の構文と設定値の妥当性を検証する。 |
| `agent-def list` | 対応しているエージェント名、`type` ごとのデフォルト出力先など、利用可能な設定の一覧を表示する。 |
| `agent-def init` | カレントディレクトリに `agent-def.yml` の雛形と、`context/`, `commands/` といった基本的なディレクトリ構造を生成する。 |

### 5.2. プロジェクトでの運用

`package.json` の `scripts` に登録することで、チームでの利用が容易になる。

```
{
  "name": "my-awesome-app",
  "scripts": {
    "agent:build": "agent-def build my-awesome-app",
    "agent:watch": "agent-def build my-awesome-app --watch"
  }
}
```

また、`git` の pre-commit フックと連携させることで、コミット前に設定が最新であることを保証できる。

```
# .pre-commit-config.yaml
repos:
-   repo: local
    hooks:
    -   id: agent-def-build
        name: agent-def build
        entry: npm run agent:build
        language: system
        files: ^(agent-def.yml|context/.*|commands/.*)$
```

## 6\. 変換ルール

### 6.1. ヘルパー関数の変換例

| ヘルパー | Roo Code | Claude Code | Cursor Code |
| --- | --- | --- | --- |
| `{{ file "path" }}` | `@/path` | `@path` | `@path` |
| `{{ reference "path" }}` | (プロンプト末尾に追記)<br>`---<br>Reference: path<br>...file content...` | (プロンプト末尾に追記)<br>`[Reference: path]\n...file content...` | (プロンプト末尾に追記)<br>`// Reference: path\n...file content...` |
| `{{ mcp "..." }}` | `MCP tool (...)` | `MCP tool (...)` | `MCP tool (...)` |

### 6.2. 出力先の変換

`agent-def.yml` の `targets` に `target` キーが指定されている場合、そのパスが最優先される。`target` キーが省略された場合のデフォルトパスは以下の通り。

| スコープ | `type` | `agent` | `concat` | デフォルトの出力先パス |
| --- | --- | --- | --- | --- |
| `projects` | `command` | `roo` | `false` | `destinations` 内の `.roomodes/` |
| `projects` | `command` | `claude` | `false` | `destinations` 内の `.claude/commands/` |
| `projects` | `memory` | `roo` | `true` | `destinations` 内の `.roo/rules/context.md` |
| `projects` | `memory` | `claude` | `true` | `destinations` 内の `CLAUDE.md` |
| `user` | `command` | `roo` | `false` | `~/.config/roo/modes/` |
| `user` | `command` | `claude` | `false` | `~/.config/claude/commands/` |
| `user` | `memory` | `roo` | `true` | `~/.config/roo/rules/global.md` |
| `user` | `memory` | `claude` | `true` | `~/.claude_global_context.md` |

*注: `concat: false` の場合、出力先はディレクトリとなり、その中にソースファイルごとのファイルが生成される。*

## 7\. 今後の展望

- **ヘルパー関数の拡充**: ファイルシステムやバージョン管理システムと連携する動的なヘルパー関数を追加する。
	- `{{ ls "dir" }}`: ディレクトリ構造をツリー形式で出力する。
	- `{{ git_diff "commit" }}`: 指定したコミットの差分をコンテキストに含める。
- **対応エージェントの拡充**: 主要なAIコーディングエージェントへの対応を順次進める。
	- GitHub Copilot Workspace
	- Aider
	- その他、コミュニティからの要望が高いエージェント
- **設定のインポート機能**: 共通の設定を別のYAMLファイルに切り出し、`imports` キーで読み込めるようにすることで、複数の `agent-def.yml` 間での設定の再利用を促進する。
- **プラグインアーキテクチャ**: コミュニティが独自のエージェント対応やヘルパー関数を簡単に追加できるプラグイン機構を導入する。
- **VS Code拡張機能**: `agent-def.yml` やソースファイルの入力補完、バリデーション、プレビュー機能を提供し、開発体験を向上させる。