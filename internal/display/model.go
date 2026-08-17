package display

import (
	"fmt"
	"strings"

	"github.com/RahulSandhu/notse/internal/app"
	"github.com/RahulSandhu/notse/internal/config"
	"github.com/RahulSandhu/notse/internal/storage"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// theme holds the loaded color scheme
var theme *config.Theme

func SetTheme(t *config.Theme) {
	if t == nil {
		t = config.DefaultTheme()
	}
	theme = t
}

func accentColor() lipgloss.Color           { return lipgloss.Color(theme.Accent) }
func mutedColor() lipgloss.Color            { return lipgloss.Color(theme.Muted) }
func textOnAccent() lipgloss.Color          { return lipgloss.Color(theme.TextOnAccent) }
func selectedColor() lipgloss.Color         { return lipgloss.Color(theme.Selected) }
func titleInfoColor() lipgloss.Color        { return lipgloss.Color(theme.TitleInfo) }
func normalTitleColor() lipgloss.Color      { return lipgloss.Color(theme.NormalTitle) }
func pinIndicatorColor() lipgloss.Color     { return lipgloss.Color(theme.PinIndicatorColor) }
func helpKeyColor() lipgloss.Color          { return lipgloss.Color(theme.HelpKey) }
func pageActiveDotColor() lipgloss.Color    { return lipgloss.Color(theme.PageActiveDot) }

const (
	pageSize      = 6
	enterChar     = "↵"
	backspaceChar = "⌫"
	dotChar       = "•"
	pinChar       = "  "
	newChar       = "+"
	editChar      = "✎"
	backChar      = "esc"
)

func helpKeyStyle() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Foreground(helpKeyColor())
}

func helpDescStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(mutedColor())
}

func helpSepStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(mutedColor())
}

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

func (m *Model) totalPages() int {
	if len(m.notes) == 0 {
		return 1
	}
	return (len(m.notes) + pageSize - 1) / pageSize
}

func (m *Model) pageStart() int {
	return m.page * pageSize
}

func (m *Model) pageEnd() int {
	end := m.pageStart() + pageSize
	if end > len(m.notes) {
		end = len(m.notes)
	}
	return end
}

func (m *Model) clampCursor() {
	if m.cursor < 0 {
		m.cursor = 0
	}
	if m.cursor >= len(m.notes) {
		m.cursor = len(m.notes) - 1
		if m.cursor < 0 {
			m.cursor = 0
		}
	}
	m.page = m.cursor / pageSize
}

func (m *Model) goToPage(page int) {
	if page < 0 {
		page = 0
	}
	maxPage := m.totalPages() - 1
	if page > maxPage {
		page = maxPage
	}
	m.page = page
	m.cursor = m.pageStart()
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
	page          int    // current page index for paginated list
	showHelp      bool   // expanded help visible in list mode
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
		page:          0,
		showHelp:      false,
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
			if m.mode == "list" && m.cursor > m.pageStart() {
				m.cursor--
			}

		case "down", "j":
			if m.mode == "list" && m.cursor < m.pageEnd()-1 {
				m.cursor++
			}

		case "left", "h":
			if m.mode == "list" && m.page > 0 {
				m.goToPage(m.page - 1)
			}

		case "right", "l":
			if m.mode == "list" && m.page < m.totalPages()-1 {
				m.goToPage(m.page + 1)
			}

		case "?":
			if m.mode == "list" {
				m.showHelp = !m.showHelp
			}

		case "enter":
			if m.mode == "list" && len(m.notes) > 0 {
				// View the selected note
				m.mode = "view"
			}

		case "esc":
			if m.mode == "list" {
				return m, tea.Quit
			}
			// Return to list view
			m.mode = "list"
			m.showHelp = false

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

		case "backspace":
			if m.mode == "list" && len(m.notes) > 0 {
				// Delete selected note
				noteID := m.notes[m.cursor].ID
				m.storage.Delete(noteID)
				m.notes, _ = m.storage.Load()
				m.clampCursor()
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
	render := lipgloss.NewStyle().Padding(1, 1).Render
	switch m.mode {
	case "view":
		if len(m.notes) > 0 {
			return render(m.renderNoteView())
		}
		return render(m.renderListView())
	case "create":
		return render(m.renderEditView())
	case "edit":
		return render(m.renderEditView())
	default:
		return render(m.renderListView())
	}
}

// renderListView shows the list of notes
func (m Model) renderListView() string {
	contentPadding := lipgloss.NewStyle().PaddingLeft(2)
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(textOnAccent()).
		Background(accentColor()).
		Padding(0, 1)

	countStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(titleInfoColor())

	selectedStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(selectedColor()).
		BorderLeft(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(selectedColor()).
		PaddingLeft(1)

	normalTitleStyle := lipgloss.NewStyle().
		Foreground(normalTitleColor())

	metaStyle := lipgloss.NewStyle().
		Foreground(mutedColor())

	pinStyle := lipgloss.NewStyle().
		Foreground(pinIndicatorColor())

	var sb strings.Builder
	sb.WriteString(contentPadding.Render(headerStyle.Render("Notes History")) + "\n\n")
	sb.WriteString(contentPadding.Render(countStyle.Render(m.noteCountLabel())) + "\n\n")

	if len(m.notes) == 0 {
		sb.WriteString(contentPadding.Render("No notes yet. Press 'n' to create one.") + "\n")
	}

	start := m.pageStart()
	end := m.pageEnd()
	for i := start; i < end; i++ {
		note := m.notes[i]
		icon := lipgloss.NewStyle().Foreground(statusColor(note.Status)).Render(statusIcon(note.Status))
		title := note.Title
		if title == "" {
			title = "(untitled)"
		}

		line := fmt.Sprintf("%s %s", icon, normalTitleStyle.Render(title))
		timestamp := note.CreatedAt.Format("2006-01-02 15:04")

		pin := ""
		if note.IsPinned {
			pin = pinStyle.Render(pinChar)
		}

		rawLine := fmt.Sprintf("%s %s", statusIcon(note.Status), title)
		if m.cursor == i {
			sb.WriteString(selectedStyle.Render(rawLine) + "\n")
			rawTimestamp := timestamp
			if note.IsPinned {
				rawTimestamp += pinChar
			}
			sb.WriteString(selectedStyle.Render(rawTimestamp) + "\n")
		} else {
			sb.WriteString(contentPadding.Render(line) + "\n")
			styledTimestamp := timestamp
			if pin != "" {
				styledTimestamp += pin
			}
			sb.WriteString(contentPadding.Render(metaStyle.Render(styledTimestamp)) + "\n")
		}
	}

	sb.WriteString("\n")
	sb.WriteString(m.renderPaginationDots())

	if m.showHelp {
		sb.WriteString(m.renderListFullHelp())
	} else {
		sb.WriteString(m.renderListShortHelp())
	}

	return sb.String()
}

func (m Model) noteCountLabel() string {
	count := len(m.notes)
	if count == 1 {
		return "1 item"
	}
	return fmt.Sprintf("%d items", count)
}

func (m Model) renderPaginationDots() string {
	activeStyle := lipgloss.NewStyle().Foreground(pageActiveDotColor())
	inactiveStyle := lipgloss.NewStyle().Foreground(mutedColor())

	dotsStyle := lipgloss.NewStyle().MarginLeft(2).MarginBottom(1)

	var sb strings.Builder
	for i := 0; i < m.totalPages(); i++ {
		if i == m.page {
			sb.WriteString(activeStyle.Render(dotChar))
		} else {
			sb.WriteString(inactiveStyle.Render(dotChar))
		}
	}
	return dotsStyle.Render(sb.String()) + "\n"
}

func (m Model) renderListShortHelp() string {
	moreLabel := "more"
	if m.showHelp {
		moreLabel = "less"
	}

	helpStyle := lipgloss.NewStyle().PaddingLeft(2)

	return helpStyle.Render(lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.helpItem(newChar, "new note"),
		helpSepStyle().Render(" • "),
		m.helpItem(backspaceChar, "delete"),
		helpSepStyle().Render(" • "),
		m.helpItem(pinChar, "pin/unpin"),
		helpSepStyle().Render(" • "),
		m.helpItem("?", moreLabel),
	) + "\n")
}

func (m Model) renderListFullHelp() string {
	moreLabel := "less"
	fullHelpStyle := lipgloss.NewStyle().PaddingLeft(2)

	return fullHelpStyle.Render(lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.helpColumn(
			m.helpRow("k/j", "up/down"),
			m.helpRow("h/l", "page"),
			m.helpRow("home", "start"),
			m.helpRow("end", "end"),
			m.helpRow("esc/q", "quit"),
		),
		"  ",
		m.helpColumn(
			m.helpRow(enterChar, "view"),
			m.helpRow("n", "new note"),
			m.helpRow(backspaceChar, "delete"),
			m.helpRow("p", "pin/unpin"),
		),
		"  ",
		m.helpItem("?", moreLabel),
	) + "\n")
}

func (m Model) helpColumn(rows ...string) string {
	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (m Model) helpRow(key, desc string) string {
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		helpKeyStyle().Render(key),
		" ",
		helpDescStyle().Render(desc),
	)
}

func (m Model) helpItem(key, desc string) string {
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		helpKeyStyle().Render(key),
		" ",
		helpDescStyle().Render(desc),
	)
}

// statusColor returns the color for a note status
func statusColor(status string) lipgloss.Color {
	switch status {
	case "done":
		return lipgloss.Color(theme.Accent)
	case "missed":
		return lipgloss.Color("#e0e0e0")
	default:
		return lipgloss.Color(theme.Muted)
	}
}

// renderNoteView shows a single note in detail
func (m Model) renderNoteView() string {
	note := m.notes[m.cursor]

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(textOnAccent()).
		Background(accentColor()).
		Padding(0, 1)

	contentPadding := lipgloss.NewStyle().PaddingLeft(2)

	metaStyle := lipgloss.NewStyle().
		Foreground(mutedColor())

	statusStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(statusColor(note.Status))

	var sb strings.Builder
	sb.WriteString(contentPadding.Render(headerStyle.Render(note.Title)) + "\n\n")
	for _, line := range strings.Split(note.Content, "\n") {
		sb.WriteString(contentPadding.Render(line) + "\n")
	}
	sb.WriteString("\n")
	sb.WriteString(contentPadding.Render("Status: " + statusStyle.Render(fmt.Sprintf("%s %s", statusIcon(note.Status), note.Status))) + "\n")
	sb.WriteString(contentPadding.Render(metaStyle.Render(fmt.Sprintf("Created: %s", note.CreatedAt.Format("2006-01-02 15:04")))) + "\n")
	sb.WriteString("\n")
	sb.WriteString(m.renderViewHelp())

	return sb.String()
}

func (m Model) renderViewHelp() string {
	statusIconStr := statusIcon(m.notes[m.cursor].Status)

	return lipgloss.NewStyle().PaddingLeft(2).Render(lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.helpItem(statusIconStr, "status"),
		helpSepStyle().Render(" • "),
		m.helpItem(editChar, "edit"),
		helpSepStyle().Render(" • "),
		m.helpItem(backChar, "back"),
	) + "\n")
}

// renderEditView shows the note creation/edit form
func (m Model) renderEditView() string {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(textOnAccent()).
		Background(accentColor()).
		Padding(0, 1)

	contentPadding := lipgloss.NewStyle().PaddingLeft(2)
	helpStyle := lipgloss.NewStyle().
		Foreground(mutedColor()).
		PaddingLeft(2)

	var sb strings.Builder

	if m.mode == "edit" {
		sb.WriteString(contentPadding.Render(headerStyle.Render("Edit Note")) + "\n\n")
	} else {
		sb.WriteString(contentPadding.Render(headerStyle.Render("Create New Note")) + "\n\n")
	}

	sb.WriteString(contentPadding.Render(m.titleInput.View()) + "\n\n")
	sb.WriteString(contentPadding.Render(m.contentArea.View()) + "\n\n")
	sb.WriteString(helpStyle.Render("tab switch fields • ctrl+s save • esc cancel") + "\n")

	return sb.String()
}
