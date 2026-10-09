package stage

import (
	"strings"
	"testing"

	"github.com/fancy-cuppa/cuppa/apps/cuppa-tui/pointer"
	"github.com/fancy-cuppa/cuppa/libs/canvas/editor"
	"github.com/fancy-cuppa/cuppa/libs/catalog/standard"
	"github.com/fancy-cuppa/cuppa/libs/document/design"
)

func setup(t *testing.T) (*Model, *editor.Editor, design.NodeID) {
	t.Helper()
	cat := standard.Default()
	ed := editor.New(cat, design.NewDocument("t", 100, 40))
	id, err := ed.Add("lipgloss.box", 10, 5) // 24x6 at (10,5)
	if err != nil {
		t.Fatal(err)
	}
	ed.Clear()
	m := New(ed, cat)
	m.SetSize(60, 20)
	return m, ed, id
}

func down(m *Model, x, y int) { m.Handle(pointer.Event{X: x, Y: y, Phase: pointer.Down, Left: true}) }
func drag(m *Model, x, y int) { m.Handle(pointer.Event{X: x, Y: y, Phase: pointer.Move, Held: true}) }
func up(m *Model, x, y int)   { m.Handle(pointer.Event{X: x, Y: y, Phase: pointer.Up, Left: true}) }
func rect(ed *editor.Editor, id design.NodeID) design.Rect {
	n, _ := ed.Document().Get(id)
	return n.Rect
}

func TestClickSelectsAndEmptyClickClears(t *testing.T) {
	m, ed, id := setup(t)
	down(m, 15, 8)
	up(m, 15, 8)
	if sel := ed.Selected(); len(sel) != 1 || sel[0] != id {
		t.Fatalf("selection = %v", sel)
	}
	down(m, 50, 18)
	up(m, 50, 18)
	if len(ed.Selected()) != 0 {
		t.Fatal("clicking empty canvas should deselect")
	}
}

func TestDragMovesByThePointerDelta(t *testing.T) {
	m, ed, id := setup(t)
	down(m, 15, 8) // inner cell, grabbed 5,3 from the top-left
	drag(m, 20, 10)
	drag(m, 25, 12)
	up(m, 25, 12)
	if got := rect(ed, id); got.X != 20 || got.Y != 9 {
		t.Fatalf("rect = %+v", got)
	}
	ed.Undo()
	if got := rect(ed, id); got.X != 10 || got.Y != 5 {
		t.Fatalf("drag should undo as one step: %+v", got)
	}
}

func TestPlainClickLeavesNoUndoStep(t *testing.T) {
	m, ed, _ := setup(t)
	ed.MarkSaved()
	down(m, 15, 8)
	up(m, 15, 8)
	if ed.Dirty() {
		t.Fatal("clicking changed the document")
	}
}

func TestDragCornerResizes(t *testing.T) {
	m, ed, id := setup(t)
	down(m, 15, 8)
	up(m, 15, 8)    // select
	down(m, 33, 10) // bottom-right corner cell of (10,5,24,6)
	drag(m, 38, 13)
	up(m, 38, 13)
	if got := rect(ed, id); got.W != 29 || got.H != 9 || got.X != 10 || got.Y != 5 {
		t.Fatalf("rect = %+v", got)
	}
}

func TestMarqueeSelectsIntersectingNodes(t *testing.T) {
	m, ed, a := setup(t)
	b, _ := ed.Add("lipgloss.label", 40, 15)
	ed.Clear()
	down(m, 5, 2)
	drag(m, 45, 16)
	up(m, 45, 16)
	sel := ed.Selected()
	if len(sel) != 2 {
		t.Fatalf("selection = %v (want %v and %v)", sel, a, b)
	}
}

func TestShiftClickTogglesSelection(t *testing.T) {
	m, ed, a := setup(t)
	b, _ := ed.Add("lipgloss.label", 40, 15)
	ed.Select(a)
	m.Handle(pointer.Event{X: 41, Y: 15, Phase: pointer.Down, Left: true, Shift: true})
	m.Handle(pointer.Event{X: 41, Y: 15, Phase: pointer.Up, Left: true})
	if sel := ed.Selected(); len(sel) != 2 {
		t.Fatalf("selection = %v", sel)
	}
	m.Handle(pointer.Event{X: 41, Y: 15, Phase: pointer.Down, Left: true, Shift: true})
	m.Handle(pointer.Event{X: 41, Y: 15, Phase: pointer.Up, Left: true})
	if ed.IsSelected(b) {
		t.Fatal("second shift-click should deselect")
	}
}

func TestWheelScrollsTheViewport(t *testing.T) {
	m, _, _ := setup(t)
	m.Handle(pointer.Event{Phase: pointer.Wheel, WheelY: 1})
	if m.offY != 2 {
		t.Fatalf("offY = %d", m.offY)
	}
	for i := 0; i < 50; i++ {
		m.Handle(pointer.Event{Phase: pointer.Wheel, WheelY: 1})
	}
	if m.offY != 20 { // doc height 40 - pane height 20
		t.Fatalf("offY not clamped: %d", m.offY)
	}
	if x, y := m.Canvas(0, 0); x != 0 || y != 20 {
		t.Fatalf("Canvas(0,0) = %d,%d", x, y)
	}
}

func TestRenderIsExactlyPaneSized(t *testing.T) {
	m, ed, id := setup(t)
	ed.Select(id)
	g := design.Rect{X: 3, Y: 3, W: 5, H: 2}
	m.SetGhost(&g)
	lines := m.Lines()
	if len(lines) != 20 {
		t.Fatalf("lines = %d", len(lines))
	}
	if !strings.Contains(strings.Join(lines, ""), "■") {
		t.Fatal("selection handles not drawn")
	}
}

func TestDraggingSnapsToANeighbourAndShowsAGuide(t *testing.T) {
	m, ed, id := setup(t) // box at (10,5) 24x6
	other, _ := ed.Add("lipgloss.label", 60, 20)
	ed.Clear()
	// Move the box so its left edge lands one cell off the label's left edge (60).
	down(m, 15, 8)
	drag(m, 15+49, 8) // box.X -> 59
	if got := rect(ed, id); got.X != 60 {
		t.Fatalf("X = %d, want 60 (snapped to the label at %v)", got.X, other)
	}
	if len(m.guides) == 0 {
		t.Fatal("no guide recorded while snapped")
	}
	up(m, 15+49, 8)
	if len(m.guides) != 0 {
		t.Fatal("guides should clear on release")
	}
}

func TestSnapCanBeSwitchedOff(t *testing.T) {
	m, ed, id := setup(t)
	if _, err := ed.Add("lipgloss.label", 60, 20); err != nil {
		t.Fatal(err)
	}
	ed.Clear()
	m.SetSnap(false)
	down(m, 15, 8)
	drag(m, 15+49, 8)
	up(m, 15+49, 8)
	if got := rect(ed, id); got.X != 59 {
		t.Fatalf("X = %d, want the unsnapped 59", got.X)
	}
}
