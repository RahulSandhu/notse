#!/bin/bash

NOTES_FILE="$HOME/.config/notse/notes.json"

mkdir -p "$(dirname "$NOTES_FILE")"

cat > "$NOTES_FILE" << 'EOF'
{
  "notes": [
    {
      "id": "20240101120000",
      "title": "Welcome to Notse",
      "content": "This is your first note! Press 'enter' to view it in detail.\n\nYou can:\n- Navigate with j/k or arrow keys\n- Press 'd' to delete\n- Press 'p' to pin\n- Press 'q' to quit",
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z",
      "is_pinned": true,
      "tags": ["welcome", "tutorial"]
    },
    {
      "id": "20240101130000",
      "title": "Shopping List",
      "content": "- Milk\n- Eggs\n- Bread\n- Coffee\n- Butter",
      "created_at": "2024-01-01T13:00:00Z",
      "updated_at": "2024-01-01T13:00:00Z",
      "is_pinned": false,
      "tags": ["shopping", "todo"]
    },
    {
      "id": "20240101140000",
      "title": "Project Ideas",
      "content": "1. Build a CLI notes app ✅\n2. Create a habit tracker\n3. Make a terminal file browser\n4. Build a simple task manager",
      "created_at": "2024-01-01T14:00:00Z",
      "updated_at": "2024-01-01T14:00:00Z",
      "is_pinned": false,
      "tags": ["ideas", "projects"]
    },
    {
      "id": "20240101150000",
      "title": "Meeting Notes",
      "content": "Discussed the new features for Q1:\n- User authentication\n- Dark mode support\n- Export functionality\n\nAction items:\n- Review designs by Friday\n- Setup development environment",
      "created_at": "2024-01-01T15:00:00Z",
      "updated_at": "2024-01-01T15:00:00Z",
      "is_pinned": true,
      "tags": ["work", "meetings"]
    }
  ]
}
EOF

echo "Sample notes created at $NOTES_FILE"
echo "Run 'make run' to view them!"
