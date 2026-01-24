package display

import (
	"fmt"
	"strings"

	"github.com/RahulSandhu/notse/internal/app"
	"github.com/RahulSandhu/notse/internal/storage"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the TUI state
type Model struct {
	notes         []app.Note
	storage       *storage.Storage
	cursor        int
	selected      map[int]struct{}
	mode          string // "list", "view", "create", "edit"
	err           error
	titleInput    textinput.Model
	contentArea   textarea.Model
	focusIndex    int    // 0 for title, 1 for content
	editingNoteID string // ID of note being edited
}

// NewModel creates a new TUI model
func NewModel(storage *storage.Storage) Model {
	// Load existing notes
	notes, _ := storage.Load()

	// Setup title input
	ti := textinput.New()
	ti.Placeholder = "Note title..."
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 50

	// Setup content area
	ta := textarea.New()
	ta.Placeholder = "Write your note here..."
	ta.SetWidth(60)
	ta.SetHeight(10)

	return Model{
		notes:         notes,
		storage:       storage,
		cursor:        0,
		selected:      make(map[int]struct{}),
		mode:          "list",
		titleInput:    ti,
		contentArea:   ta,
		focusIndex:    0,
		editingNoteID: "",
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates model state
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Handle create/edit mode separately
	if m.mode == "create" || m.mode == "edit" {
		return m.updateEditMode(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.mode == "list" && m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.mode == "list" && m.cursor < len(m.notes)-1 {
				m.cursor++
			}

		case "enter":
			if m.mode == "list" && len(m.notes) > 0 {
				// View the selected note
				m.mode = "view"
			}

		case "esc":
			// Return to list view
			m.mode = "list"

		case "e":
			if m.mode == "view" && len(m.notes) > 0 {
				// Edit the current note
				note := m.notes[m.cursor]
				m.mode = "edit"
				m.editingNoteID = note.ID
				m.focusIndex = 0
				m.titleInput.SetValue(note.Title)
				m.contentArea.SetValue(note.Content)
				m.titleInput.Focus()
			}

		case "n":
			if m.mode == "list" {
				// Create new note - enter create mode
				m.mode = "create"
				m.editingNoteID = ""
				m.focusIndex = 0
				m.titleInput.SetValue("")
				m.contentArea.SetValue("")
				m.titleInput.Focus()
			}

		case "d":
			if m.mode == "list" && len(m.notes) > 0 {
				// Delete selected note
				noteID := m.notes[m.cursor].ID
				m.storage.Delete(noteID)
				m.notes, _ = m.storage.Load()
				if m.cursor >= len(m.notes) && m.cursor > 0 {
					m.cursor--
				}
			}

		case "p":
			if m.mode == "list" && len(m.notes) > 0 {
				// Toggle pin
				note := m.notes[m.cursor]
				note.IsPinned = !note.IsPinned
				m.storage.Update(note.ID, note)
				m.notes, _ = m.storage.Load()
			}
		}
	}

	return m, cmd
}

// updateEditMode handles input when creating or editing a note
func (m Model) updateEditMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			// Cancel creation/editing
			m.mode = "list"
			m.editingNoteID = ""
			return m, nil

		case "tab":
			// Switch between title and content
			if m.focusIndex == 0 {
				m.focusIndex = 1
				m.titleInput.Blur()
				m.contentArea.Focus()
			} else {
				m.focusIndex = 0
				m.contentArea.Blur()
				m.titleInput.Focus()
			}
			return m, nil

		case "ctrl+s":
			// Save the note
			title := m.titleInput.Value()
			content := m.contentArea.Value()

			if title != "" || content != "" {
				if m.mode == "edit" && m.editingNoteID != "" {
					// Update existing note
					for _, note := range m.notes {
						if note.ID == m.editingNoteID {
							note.Title = title
							note.Content = content
							m.storage.Update(note.ID, note)
							break
						}
					}
				} else {
					// Create new note
					newNote := app.NewNote(title, content)
					m.storage.Add(newNote)
				}
				m.notes, _ = m.storage.Load()
				m.mode = "list"
				m.editingNoteID = ""
			}
			return m, nil
		}
	}

	// Update the focused input
	var cmd tea.Cmd
	if m.focusIndex == 0 {
		m.titleInput, cmd = m.titleInput.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.contentArea, cmd = m.contentArea.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the TUI
func (m Model) View() string {
	switch m.mode {
	case "view":
		if len(m.notes) > 0 {
			return m.renderNoteView()
		}
		return m.renderListView()
	case "create":
		return m.renderEditView()
	case "edit":
		return m.renderEditView()
	default:
		return m.renderListView()
	}
}

// renderListView shows the list of notes
func (m Model) renderListView() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		PaddingLeft(2)

	title := titleStyle.Render("📝 Notes")

	var sb strings.Builder
	sb.WriteString(title + "\n\n")

	if len(m.notes) == 0 {
		sb.WriteString("  No notes yet. Press 'n' to create one.\n")
	}

	for i, note := range m.notes {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		pinned := " "
		if note.IsPinned {
			pinned = "📌"
		}

		title := note.Title
		if title == "" {
			title = "(untitled)"
		}

		sb.WriteString(fmt.Sprintf("%s %s %s\n", cursor, pinned, title))
	}

	sb.WriteString("\n")
	sb.WriteString("  q: quit | n: new | d: delete | p: pin | enter: view\n")

	return sb.String()
}

// renderNoteView shows a single note in detail
func (m Model) renderNoteView() string {
	note := m.notes[m.cursor]

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4"))

	var sb strings.Builder
	sb.WriteString(titleStyle.Render(note.Title) + "\n\n")
	sb.WriteString(note.Content + "\n\n")
	sb.WriteString(fmt.Sprintf("Created: %s\n", note.CreatedAt.Format("2006-01-02 15:04")))
	sb.WriteString("\n")
	sb.WriteString("e: edit | esc: back | q: quit\n")

	return sb.String()
}

// renderEditView shows the note creation/edit form
func (m Model) renderEditView() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		PaddingLeft(2)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888"))

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Italic(true)

	var sb strings.Builder

	// Different title based on mode
	if m.mode == "edit" {
		sb.WriteString(titleStyle.Render("✏️  Edit Note") + "\n\n")
	} else {
		sb.WriteString(titleStyle.Render("✏️  Create New Note") + "\n\n")
	}

	// Title input
	sb.WriteString(labelStyle.Render("Title:") + "\n")
	sb.WriteString(m.titleInput.View() + "\n\n")

	// Content area
	sb.WriteString(labelStyle.Render("Content:") + "\n")
	sb.WriteString(m.contentArea.View() + "\n\n")

	// Help text
	sb.WriteString(helpStyle.Render("tab: switch fields | ctrl+s: save | esc: cancel") + "\n")

	return sb.String()
}
