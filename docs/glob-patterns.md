# Glob Pattern Support

agent-sync uses glob patterns in multiple contexts:

1. In the `inputs` field of configuration files to select input files
2. In template functions (`include`, `includeRaw`, `reference`, and `referenceRaw`) to dynamically include files

## Supported Glob Patterns

| Pattern | Description | Example |
|---------|-------------|---------|
| `*` | Matches any sequence of characters within a single path component | `docs/*.md` matches all `.md` files in the `docs` directory |
| `**` | Matches any sequence of characters across multiple path components (recursive) | `docs/**/*.md` matches all `.md` files in `docs` and all its subdirectories |
| `?` | Matches exactly one character | `file?.md` matches `file1.md` but not `file10.md` |
| `[abc]` | Matches any character within the brackets | `file[123].md` matches `file1.md`, `file2.md`, or `file3.md` |
| `{a,b,c}` | Matches any of the comma-separated patterns | `file.{md,txt}` matches `file.md` or `file.txt` |
| `!pattern` | Excludes files that match the pattern | `!*_test.md` excludes all files ending with `_test.md` |

## Pattern Order and Precedence

1. All include patterns (without `!` prefix) are processed first
2. Then exclude patterns (with `!` prefix) are applied to filter the results
3. The final list of files is sorted alphabetically by path

For example:
```
files/**/*.md      # Include all Markdown files in the files directory and subdirectories
!files/**/*_test.md # Exclude all test files
```

## File Processing Order

The final list of files will be sorted alphabetically by path. If you need files to be processed in a specific order, list them individually without glob patterns.

## Examples

### Basic Pattern Matching

```
*.md              # All Markdown files in the current directory
docs/*.md         # All Markdown files in the docs directory
```

### Recursive Directory Traversal

```
**/*.md           # All Markdown files in any directory
docs/**/*.md      # All Markdown files in the docs directory and its subdirectories
```

### Pattern Sets

```
{config,settings}.json  # Matches both config.json and settings.json
docs/*.{md,txt}         # Matches .md and .txt files in the docs directory
```

### Exclusion Patterns

```
*.md              # Include all Markdown files
!*_test.md        # Exclude test Markdown files
```

### Complex Patterns

```
src/**/*.{js,ts}  # All JavaScript and TypeScript files in src directory and subdirectories
!src/**/test/**   # Exclude all files in test directories
```

## Notes

- When no files match a glob pattern, the result is an empty list (no error)
- File paths are always returned in alphabetical order
- Pattern matching is case-sensitive on Unix-like systems and case-insensitive on Windows