package app

import (
	"testing"
)

// TestNewNote tests the NewNote function to ensure it initializes a Note correctly
func TestNewNote(t *testing.T) {
	// Define test inputs
	title := "Test Note"
	content := "This is a test note"

	// Create a new note
	note := NewNote(title, content)

	// Validate the note fields
	if note.Title != title {
		t.Errorf("Expected title %s, got %s", title, note.Title)
	}

	if note.Content != content {
		t.Errorf("Expected content %s, got %s", content, note.Content)
	}

	if note.ID == "" {
		t.Error("Expected non-empty ID")
	}

	if note.IsPinned {
		t.Error("Expected note to not be pinned by default")
	}

	if note.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}

	if note.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}
}

// TestGenerateID tests the generateID function to ensure it produces unique IDs
func TestGenerateID(t *testing.T) {
	// Generate two IDs
	id1 := generateID()
	id2 := generateID()

	// IDs should be unique
	if id1 == id2 {
		t.Errorf("Expected unique IDs, got id1=%s id2=%s", id1, id2)
	}

	// ID should not be empty
	if len(id1) == 0 {
		t.Error("Expected non-empty ID")
	}

	// ID should be reasonable length
	if len(id1) < 14 {
		t.Errorf("Expected ID length >= 14, got %d", len(id1))
	}

	// Generate multiple IDs to ensure they're all unique
	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id := generateID()
		if ids[id] {
			t.Errorf("Duplicate ID generated: %s", id)
		}
		ids[id] = true
	}
}
