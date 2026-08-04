package app

import (
	"fmt"
	"time"
)

// Note represents a single note entry
type Note struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsPinned  bool      `json:"is_pinned"`
	Status    string    `json:"status"`
	Tags      []string  `json:"tags,omitempty"`
}

// NotesList holds all notes
type NotesList struct {
	Notes []Note `json:"notes"`
}

// NewNote creates a new note with generated ID
func NewNote(title, content string) Note {
	now := time.Now()
	return Note{
		ID:        generateID(),
		Title:     title,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
		IsPinned:  false,
		Status:    "pending",
		Tags:      []string{},
	}
}

// generateID creates a simple unique ID based on timestamp with nanoseconds
func generateID() string {
	now := time.Now()
	return now.Format("20060102150405") + fmt.Sprintf("%09d", now.Nanosecond())
}
