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

var (
	accentColor  = lipgloss.Color("#b8d9ae")
	mutedColor   = lipgloss.Color("#6e6e70")
	textOnAccent = lipgloss.Color("#000000")
	selectedColor = lipgloss.Color("#445c3d")
)

// statusIcon returns the unicode icon for a note status
func statusIcon(status string) string {
	switch status {
	case "done":
		return "✓"
	case "missed":
		return "✗"
	default:
		return "○"
	}
}

// nextStatus cycles through pending -> done -> missed -> pending
func nextStatus(status string) string {
	switch status {
	case "pending":
		return "done"
	case "done":
		return "missed"
	default:
		return "pending"
	}
}

// Model represents the TUI state
type Model struct {
	notes         []app.Note
	storage       *storage.Storage
	cursor        int
	mode          string // "list", "view", "create", "edit"
	width         int
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
		mode:          "list",
		width:         80,
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

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.contentArea.SetWidth(msg.Width - 10)
		m.contentArea.SetHeight(msg.Height - 12)
		return m, nil
	}

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

		case "s":
			if m.mode == "view" && len(m.notes) > 0 {
				// Cycle note status
				note := m.notes[m.cursor]
				note.Status = nextStatus(note.Status)
				m.storage.Update(note.ID, note)
				m.notes, _ = m.storage.Load()
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
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(textOnAccent).
		Background(accentColor).
		Padding(0, 1)

	selectedBar := "│"
	selectedBarStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(selectedColor)

	selectedTitleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(selectedColor)

	selectedMetaStyle := lipgloss.NewStyle().
		Foreground(selectedColor)

	metaStyle := lipgloss.NewStyle().
		Foreground(mutedColor)

	helpStyle := lipgloss.NewStyle().
		Foreground(mutedColor)

	var sb strings.Builder
	sb.WriteString(headerStyle.Render("Notes History") + "\n\n")

	if len(m.notes) == 0 {
		sb.WriteString("  No notes yet. Press 'n' to create one.\n")
	}

	for i, note := range m.notes {
		icon := lipgloss.NewStyle().Foreground(statusColor(note.Status)).Render(statusIcon(note.Status))
		title := note.Title
		if title == "" {
			title = "(untitled)"
		}

		if note.IsPinned {
			title = title + " *"
		}

		line := fmt.Sprintf("%s %s", icon, title)
		timestamp := note.CreatedAt.Format("2006-01-02 15:04")

		if m.cursor == i {
			sb.WriteString(fmt.Sprintf("%s %s\n", selectedBarStyle.Render(selectedBar), selectedTitleStyle.Render(line)))
			sb.WriteString(fmt.Sprintf("%s %s\n", selectedBarStyle.Render(selectedBar), selectedMetaStyle.Render(timestamp)))
		} else {
			sb.WriteString(fmt.Sprintf("  %s\n", line))
			sb.WriteString(fmt.Sprintf("  %s\n", metaStyle.Render(timestamp)))
		}
	}

	sb.WriteString("\n")
	sb.WriteString(helpStyle.Render("q quit • n new • d delete • p pin • enter view") + "\n")

	return sb.String()
}

// statusColor returns the color for a note status
func statusColor(status string) lipgloss.Color {
	switch status {
	case "done":
		return lipgloss.Color("#b8d9ae")
	case "missed":
		return lipgloss.Color("#e0e0e0")
	default:
		return lipgloss.Color("#6e6e70")
	}
}

// renderNoteView shows a single note in detail
func (m Model) renderNoteView() string {
	note := m.notes[m.cursor]

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(textOnAccent).
		Background(accentColor).
		Padding(0, 1)

	metaStyle := lipgloss.NewStyle().
		Foreground(mutedColor)

	statusStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(statusColor(note.Status))

	helpStyle := lipgloss.NewStyle().
		Foreground(mutedColor)

	var sb strings.Builder
	sb.WriteString(headerStyle.Render(note.Title) + "\n\n")
	sb.WriteString(note.Content + "\n\n")
	sb.WriteString("Status: " + statusStyle.Render(fmt.Sprintf("%s %s", statusIcon(note.Status), note.Status)) + "\n")
	sb.WriteString(metaStyle.Render(fmt.Sprintf("Created: %s", note.CreatedAt.Format("2006-01-02 15:04"))) + "\n")
	sb.WriteString("\n")
	sb.WriteString(helpStyle.Render("s status • e edit • esc back • q quit") + "\n")

	return sb.String()
}

// renderEditView shows the note creation/edit form
func (m Model) renderEditView() string {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(textOnAccent).
		Background(accentColor).
		Padding(0, 1)

	helpStyle := lipgloss.NewStyle().
		Foreground(mutedColor)

	var sb strings.Builder

	if m.mode == "edit" {
		sb.WriteString(headerStyle.Render("Edit Note") + "\n\n")
	} else {
		sb.WriteString(headerStyle.Render("Create New Note") + "\n\n")
	}

	sb.WriteString(m.titleInput.View() + "\n\n")
	sb.WriteString(m.contentArea.View() + "\n\n")
	sb.WriteString(helpStyle.Render("tab switch fields • ctrl+s save • esc cancel") + "\n")

	return sb.String()
}
