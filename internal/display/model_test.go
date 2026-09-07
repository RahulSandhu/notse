package display

import (
	"strings"
	"testing"

	"github.com/RahulSandhu/notse/internal/config"
)

func TestFlashStatusSetsMessage(t *testing.T) {
	m := Model{}
	cmd := m.flashStatus("Deleted: groceries")
	if m.statusMsg != "Deleted: groceries" {
		t.Fatalf("statusMsg = %q, want %q", m.statusMsg, "Deleted: groceries")
	}
	if m.statusMsgID != 1 {
		t.Fatalf("statusMsgID = %d, want 1", m.statusMsgID)
	}
	if cmd == nil {
		t.Fatal("flashStatus returned nil cmd")
	}
}

func TestStatusClearMatchingTickClearsMessage(t *testing.T) {
	m := Model{statusMsg: "Pinned: ideas", statusMsgID: 3}
	next, _ := m.Update(statusClearMsg(3))
	got := next.(Model)
	if got.statusMsg != "" {
		t.Fatalf("statusMsg = %q, want cleared", got.statusMsg)
	}
}

func TestStatusClearStaleTickKeepsMessage(t *testing.T) {
	m := Model{statusMsg: "Pinned: ideas", statusMsgID: 3}
	next, _ := m.Update(statusClearMsg(2))
	got := next.(Model)
	if got.statusMsg != "Pinned: ideas" {
		t.Fatalf("statusMsg = %q, want unchanged", got.statusMsg)
	}
}

func TestTruncateWithEllipsis(t *testing.T) {
	if got := truncateWithEllipsis("short", 50); got != "short" {
		t.Fatalf("short string was truncated: %q", got)
	}
	if got := truncateWithEllipsis("1234567890", 5); got != "1234…" {
		t.Fatalf("truncate = %q, want %q", got, "1234…")
	}
	if got := truncateWithEllipsis("", 5); got != "" {
		t.Fatalf("empty string = %q, want empty", got)
	}
}

func TestRenderListShowsFlashBesideTitle(t *testing.T) {
	SetTheme(config.DefaultTheme())
	m := Model{width: 100, statusMsg: "Created: hello"}
	out := m.renderListView()
	if !strings.Contains(out, "Created: hello") {
		t.Fatalf("renderListView missing flash text:\n%s", out)
	}
}
