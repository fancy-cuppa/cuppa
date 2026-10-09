package shell

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

func TestLayoutHonoursChosenWidthsWithinLimits(t *testing.T) {
	l := computeLayout(140, 40, 40, 50)
	if l.palette.W != 40 || l.inspector.W != 50 || l.stage.W != 140-40-50-2 {
		t.Fatalf("layout = %+v", l)
	}
	// Too narrow a choice is raised to the minimum.
	l = computeLayout(140, 40, 3, 3)
	if l.palette.W != minPaletteW || l.inspector.W != minInspectorW {
		t.Fatalf("layout = %+v", l)
	}
	// The stage keeps its room: the side bars give way on a narrow terminal.
	l = computeLayout(70, 20, 40, 40)
	if l.stage.W < minStageW {
		t.Fatalf("stage squeezed to %d", l.stage.W)
	}
}

func TestDividersAreTheColumnsBetweenPanes(t *testing.T) {
	l := computeLayout(140, 40, 0, 0)
	if l.dividerAt(l.palette.W, 5) != leftDivider {
		t.Fatal("column after the palette is the left divider")
	}
	if l.dividerAt(l.inspector.X-1, 5) != rightDivider {
		t.Fatal("column before the details bar is the right divider")
	}
	if l.dividerAt(l.palette.W, 0) != noDivider || l.dividerAt(l.palette.W, 39) != noDivider {
		t.Fatal("the menu bar and status bar rows are not dividers")
	}
	if l.dividerAt(5, 5) != noDivider {
		t.Fatal("a pane is not a divider")
	}
}

func TestWidthsForClampsToTheLimits(t *testing.T) {
	pal, ins := widthsFor(140, 26, 32, leftDivider, 2)
	if pal != minPaletteW || ins != 32 {
		t.Fatalf("too far left: %d %d", pal, ins)
	}
	pal, _ = widthsFor(140, 26, 32, leftDivider, 139)
	if pal != 140-32-2-minStageW {
		t.Fatalf("too far right keeps the stage: %d", pal)
	}
	_, ins = widthsFor(140, 26, 32, rightDivider, 139)
	if ins != minInspectorW {
		t.Fatalf("details bar minimum: %d", ins)
	}
	_, ins = widthsFor(140, 26, 32, rightDivider, 0)
	if ins != 140-26-2-minStageW {
		t.Fatalf("details bar maximum: %d", ins)
	}
}

func TestDraggingTheLeftDividerResizesThePalette(t *testing.T) {
	m := newShell(t)
	x := m.layout.palette.W
	send(m, click(x, 10))
	send(m, motion(x+8, 10))
	if m.layout.palette.W != x+8 {
		t.Fatalf("palette width %d, want %d", m.layout.palette.W, x+8)
	}
	send(m, motion(x+14, 12))
	send(m, release(x+14, 12))
	if m.layout.palette.W != x+14 || m.grab != noDivider {
		t.Fatalf("after release: width %d grab %v", m.layout.palette.W, m.grab)
	}
	if len(m.Editor().Document().Nodes) != 0 {
		t.Fatal("dragging a divider must not touch the design")
	}
	// Every row still fills the terminal exactly.
	for i, l := range strings.Split(m.render(), "\n") {
		if w := ansi.StringWidth(l); w != 140 {
			t.Fatalf("row %d is %d cells wide", i, w)
		}
	}
}

func TestDraggingTheRightDividerResizesTheDetailsBar(t *testing.T) {
	m := newShell(t)
	x := m.layout.inspector.X - 1
	send(m, click(x, 10))
	send(m, motion(x-10, 10))
	send(m, release(x-10, 10))
	if m.layout.inspector.W != 32+10 {
		t.Fatalf("details width %d", m.layout.inspector.W)
	}
	if m.layout.stage.W != 140-m.layout.palette.W-m.layout.inspector.W-2 {
		t.Fatalf("the stage takes the rest: %+v", m.layout)
	}
}

func TestDividerLightsUpOnHoverAndWhileDragging(t *testing.T) {
	m := newShell(t)
	x := m.layout.palette.W
	if strings.Contains(m.render(), "┃") {
		t.Fatal("dividers are quiet by default")
	}
	send(m, motion2(x, 10)) // hover without a button
	if !strings.Contains(m.render(), "┃") {
		t.Fatal("hovering a divider should light it")
	}
	if !strings.Contains(strings.Join(screenText(m), "\n"), "width of the panel") {
		t.Fatal("the status line should say what dragging does")
	}
	send(m, motion2(60, 10))
	if strings.Contains(m.render(), "┃") {
		t.Fatal("moving away calms it again")
	}
}

func TestPressingAPaneNextToADividerStillWorks(t *testing.T) {
	m := newShell(t)
	x := m.layout.palette.W
	send(m, click(x-2, 4)) // inside the palette: starts a component drag as before
	if m.dragging == "" {
		t.Fatal("a click one cell off the divider belongs to the pane")
	}
}

func TestWidthsAreRememberedAndRestored(t *testing.T) {
	file := filepath.Join(t.TempDir(), "cuppa", "layout.json")
	m := newShell(t)
	m.restoreLayoutFrom(file) // nothing saved yet: defaults stay
	if m.layout.palette.W != preferredPaletteW {
		t.Fatalf("palette %d", m.layout.palette.W)
	}
	x := m.layout.palette.W
	send(m, click(x, 10))
	send(m, motion(x+9, 10))
	send(m, release(x+9, 10))

	again := newShell(t)
	again.restoreLayoutFrom(file)
	if again.layout.palette.W != x+9 {
		t.Fatalf("restored palette width %d, want %d", again.layout.palette.W, x+9)
	}
}

func TestLayoutIsNotWrittenUnlessAFrontEndAsksForIt(t *testing.T) {
	m := newShell(t) // no RestoreLayout call: layoutFile stays empty
	x := m.layout.palette.W
	send(m, click(x, 10))
	send(m, motion(x+5, 10))
	send(m, release(x+5, 10))
	if m.layoutFile != "" {
		t.Fatal("tests and embedders must not touch the user's settings")
	}
}
