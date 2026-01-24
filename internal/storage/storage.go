package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/RahulSandhu/notse/internal/app"
)

// Storage handles note persistence
type Storage struct {
	FilePath string
}

// NewStorage creates a new storage instance
func NewStorage(filePath string) (*Storage, error) {
	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	// Create file if it doesn't exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		emptyList := app.NotesList{Notes: []app.Note{}}
		data, err := json.MarshalIndent(emptyList, "", "  ")
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile(filePath, data, 0644); err != nil {
			return nil, err
		}
	}

	return &Storage{FilePath: filePath}, nil
}

// Load reads all notes from file
func (s *Storage) Load() ([]app.Note, error) {
	data, err := os.ReadFile(s.FilePath)
	if err != nil {
		return nil, err
	}

	var notesList app.NotesList
	if err := json.Unmarshal(data, &notesList); err != nil {
		return nil, err
	}

	return notesList.Notes, nil
}

// Save writes all notes to file
func (s *Storage) Save(notes []app.Note) error {
	notesList := app.NotesList{Notes: notes}
	data, err := json.MarshalIndent(notesList, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.FilePath, data, 0644)
}

// Add appends a new note
func (s *Storage) Add(note app.Note) error {
	notes, err := s.Load()
	if err != nil {
		return err
	}

	notes = append(notes, note)
	return s.Save(notes)
}

// Delete removes a note by ID
func (s *Storage) Delete(id string) error {
	notes, err := s.Load()
	if err != nil {
		return err
	}

	filtered := make([]app.Note, 0)
	for _, note := range notes {
		if note.ID != id {
			filtered = append(filtered, note)
		}
	}

	return s.Save(filtered)
}

// Update modifies an existing note
func (s *Storage) Update(id string, updatedNote app.Note) error {
	notes, err := s.Load()
	if err != nil {
		return err
	}

	for i, note := range notes {
		if note.ID == id {
			notes[i] = updatedNote
			return s.Save(notes)
		}
	}

	return nil // Note not found
}
