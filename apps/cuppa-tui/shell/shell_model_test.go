package shell

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/fancy-cuppa/cuppa/libs/catalog/standard"
)

func newShell(t *testing.T) *Model {
	t.Helper()
	m := New(standard.Default())
	m.Update(tea.WindowSizeMsg{Width: 140, Height: 40})
	return m
}

func send(m *Model, msg tea.Msg) { m.Update(msg) }

func click(x, y int) tea.Msg   { return tea.MouseClickMsg{X: x, Y: y, Button: tea.MouseLeft} }
func motion(x, y int) tea.Msg  { return tea.MouseMotionMsg{X: x, Y: y, Button: tea.MouseLeft} }
func release(x, y int) tea.Msg { return tea.MouseReleaseMsg{X: x, Y: y, Button: tea.MouseLeft} }

func TestLayoutFillsTheTerminal(t *testing.T) {
	m := newShell(t)
	out := strings.Split(m.render(), "\n")
	if len(out) != 40 {
		t.Fatalf("rows = %d", len(out))
	}
}

func TestDragFromPaletteAndDropCreatesNode(t *testing.T) {
	m := newShell(t)
	// Palette: line 0 is the shell title bar, so pane row 2 is screen row 3
	// (first family header); the first item sits on screen row 4.
	send(m, click(5, 4))
	if m.dragging == "" {
		t.Fatal("pressing a palette item should start a drag")
	}
	stageX := m.layout.stage.X + 10
	send(m, motion(stageX, 8))
	if m.stg == nil || m.dragging == "" {
		t.Fatal("drag lost")
	}
	send(m, release(stageX, 8))
	doc := m.Editor().Document()
	if len(doc.Nodes) != 1 {
		t.Fatalf("nodes = %d", len(doc.Nodes))
	}
	if r := doc.Nodes[0].Rect; r.X != 10 || r.Y != 7 {
		t.Fatalf("dropped at %+v, want cell (10,7)", r)
	}
	if m.dragging != "" {
		t.Fatal("drag state not cleared")
	}
}

func TestReleasingOutsideTheCanvasCancelsTheDrag(t *testing.T) {
	m := newShell(t)
	send(m, click(5, 4))
	send(m, motion(m.layout.stage.X+5, 8))
	send(m, release(2, 10)) // back over the palette
	if n := len(m.Editor().Document().Nodes); n != 0 {
		t.Fatalf("nodes = %d", n)
	}
	if m.dragging != "" {
		t.Fatal("drag state not cleared")
	}
}

func TestDroppedNodeCanBeSelectedMovedAndInspected(t *testing.T) {
	m := newShell(t)
	send(m, click(5, 4))
	x0 := m.layout.stage.X
	send(m, motion(x0+10, 8))
	send(m, release(x0+10, 8))
	// Select by clicking inside, then drag 6 cells right.
	send(m, click(x0+12, 9))
	send(m, motion(x0+18, 9))
	send(m, release(x0+18, 9))
	n, ok := m.Editor().Primary()
	if !ok || n.Rect.X != 16 {
		t.Fatalf("node = %+v ok=%v", n, ok)
	}
	if !strings.Contains(m.render(), n.Name) {
		t.Fatal("inspector/layers should mention the selected node")
	}
}

func TestKeyboardBasics(t *testing.T) {
	m := newShell(t)
	send(m, click(5, 4))
	send(m, release(m.layout.stage.X+3, 5))
	if len(m.Editor().Document().Nodes) != 1 {
		t.Fatal("setup failed")
	}
	send(m, tea.KeyPressMsg{Code: tea.KeyDelete})
	if len(m.Editor().Document().Nodes) != 0 {
		t.Fatal("delete key should remove the selection")
	}
	_, cmd := m.Update(tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl})
	if cmd == nil {
		t.Fatal("ctrl+c should quit")
	}
}

func TestNarrowTerminalStillHasAStage(t *testing.T) {
	l := computeLayout(60, 20)
	if l.stage.W < 1 || l.palette.W < 16 || l.inspector.W < 20 {
		t.Fatalf("layout = %+v", l)
	}
}
