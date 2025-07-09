---
name: quick-note2
description: Create a timestamped quick note in your notes directory
---

# Quick Note

Creates a new Markdown note with a timestamp in your default notes folder.

## Commands

```bash
# define notes directory
notes_dir="${NOTES_DIR:-$HOME/notes}"
mkdir -p "$notes_dir"

# create filename with timestamp
timestamp=$(date +'%Y%m%d_%H%M%S')
file="$notes_dir/note_$timestamp.md"

# initialize note
echo "# Quick Note - $timestamp" > "$file"
echo "" >> "$file"
echo "> Add your note content here" >> "$file"

echo "Created quick note: $file"
```

## Environment Variables

- `NOTES_DIR` (optional): path to notes directory (default: `~/notes`)

## Troubleshooting

- Ensure you have write permissions to the notes directory.