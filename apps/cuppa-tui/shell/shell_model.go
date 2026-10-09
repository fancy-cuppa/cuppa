// Package shell is the root of the terminal app: it lays out the panes,
// translates Bubble Tea messages into pointer events and wires drag and drop
// from the palette to the stage.
package shell

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/fancy-cuppa/cuppa/apps/cuppa-tui/inspector"
	"github.com/fancy-cuppa/cuppa/apps/cuppa-tui/palette"
	"github.com/fancy-cuppa/cuppa/apps/cuppa-tui/pointer"
	"github.com/fancy-cuppa/cuppa/apps/cuppa-tui/stage"
	"github.com/fancy-cuppa/cuppa/apps/cuppa-tui/theme"
	"github.com/fancy-cuppa/cuppa/libs/canvas/editor"
	"github.com/fancy-cuppa/cuppa/libs/catalog/registry"
	"github.com/fancy-cuppa/cuppa/libs/document/design"
)

// pane identifies which pane owns a pointer gesture.
type pane int

const (
	nowhere pane = iota
	inPalette
	inStage
	inInspector
)

// Default canvas size of a new design, in cells.
const (
	defaultWidth  = 120
	defaultHeight = 40
)

// Model is the whole application.
type Model struct {
	cat *registry.Registry
	ed  *editor.Editor

	pal *palette.Model
	stg *stage.Model
	ins *inspector.Model

	w, h   int
	layout layout

	// owner is the pane a held-button gesture started in; it keeps receiving
	// events even when the pointer leaves it.
	owner pane
	// dragging is the component being dragged out of the palette, or "".
	dragging string
	mouseX   int
	mouseY   int
}

// New returns the app with an empty design.
func New(cat *registry.Registry) *Model {
	ed := editor.New(cat, design.NewDocument("Untitled", defaultWidth, defaultHeight))
	m := &Model{
		cat: cat,
		ed:  ed,
		pal: palette.New(cat),
		stg: stage.New(ed, cat),
		ins: inspector.New(ed, cat),
	}
	m.ins.BindSnap(m.stg.Snap, m.stg.SetSnap)
	return m
}

// Editor exposes the editor, mainly for tests.
func (m *Model) Editor() *editor.Editor { return m.ed }

// Init implements tea.Model.
func (m *Model) Init() tea.Cmd { return nil }

// Update implements tea.Model.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.resize(msg.Width, msg.Height)
	case tea.MouseClickMsg:
		m.mouse(toEvent(tea.Mouse(msg), pointer.Down))
	case tea.MouseReleaseMsg:
		m.mouse(toEvent(tea.Mouse(msg), pointer.Up))
	case tea.MouseMotionMsg:
		m.mouse(toEvent(tea.Mouse(msg), pointer.Move))
	case tea.MouseWheelMsg:
		m.mouse(toEvent(tea.Mouse(msg), pointer.Wheel))
	case tea.KeyPressMsg:
		return m, m.key(msg)
	}
	return m, nil
}

func (m *Model) resize(w, h int) {
	m.w, m.h = w, h
	m.layout = computeLayout(w, h)
	m.pal.SetSize(m.layout.palette.W, m.layout.palette.H)
	m.stg.SetSize(m.layout.stage.W, m.layout.stage.H)
	m.ins.SetSize(m.layout.inspector.W, m.layout.inspector.H)
}

func toEvent(mouse tea.Mouse, phase pointer.Phase) pointer.Event {
	e := pointer.Event{X: mouse.X, Y: mouse.Y, Phase: phase, Shift: mouse.Mod.Contains(tea.ModShift)}
	switch phase {
	case pointer.Down, pointer.Up:
		e.Left = mouse.Button == tea.MouseLeft
	case pointer.Move:
		e.Held = mouse.Button == tea.MouseLeft
	case pointer.Wheel:
		switch mouse.Button {
		case tea.MouseWheelUp:
			e.WheelY = -1
		case tea.MouseWheelDown:
			e.WheelY = 1
		case tea.MouseWheelLeft:
			e.WheelX = -1
		case tea.MouseWheelRight:
			e.WheelX = 1
		}
	}
	return e
}

// mouse routes one pointer event to the pane that should handle it.
func (m *Model) mouse(e pointer.Event) {
	m.mouseX, m.mouseY = e.X, e.Y
	target := m.layout.paneAt(e.X, e.Y)
	if m.owner != nowhere && (e.Phase == pointer.Move && e.Held || e.Phase == pointer.Up) {
		target = m.owner
	}
	if e.Phase != pointer.Move || !e.Held {
		m.pal.Leave()
	}
	switch target {
	case inPalette:
		local := e.Translate(m.layout.palette.X, m.layout.palette.Y)
		if id := m.pal.Handle(local); id != "" {
			m.dragging = id
			m.owner = inPalette
		}
	case inStage:
		if m.dragging == "" {
			m.stg.Handle(e.Translate(m.layout.stage.X, m.layout.stage.Y))
		}
	case inInspector:
		m.ins.Handle(e.Translate(m.layout.inspector.X, m.layout.inspector.Y))
	}
	if m.dragging != "" {
		m.dragFromPalette(e)
	}
	if e.Phase == pointer.Down && m.dragging == "" {
		m.owner = target
	}
	if e.Phase == pointer.Up {
		m.owner = nowhere
	}
}

// dragFromPalette shows a ghost over the stage and drops on release.
func (m *Model) dragFromPalette(e pointer.Event) {
	cx, cy, over := m.stageCell(e.X, e.Y)
	def, _ := m.cat.Get(m.dragging)
	switch {
	case e.Phase == pointer.Up:
		if over {
			if _, err := m.ed.Add(m.dragging, cx, cy); err != nil {
				return
			}
		}
		m.endDrag()
	case over:
		r := design.Rect{X: cx, Y: cy, W: def.DefaultSize.W, H: def.DefaultSize.H}
		m.stg.SetGhost(&r)
	default:
		m.stg.SetGhost(nil)
	}
}

func (m *Model) endDrag() {
	m.dragging = ""
	m.stg.SetGhost(nil)
}

// stageCell converts screen coordinates to a canvas cell and reports whether
// the pointer is over the stage.
func (m *Model) stageCell(x, y int) (int, int, bool) {
	r := m.layout.stage
	lx, ly := x-r.X, y-r.Y
	if !m.stg.Contains(lx, ly) {
		return 0, 0, false
	}
	cx, cy := m.stg.Canvas(lx, ly)
	return cx, cy, true
}

func (m *Model) key(msg tea.KeyPressMsg) tea.Cmd {
	k := msg.Key()
	text := msg.String()
	switch text {
	case "ctrl+c", "ctrl+q":
		return tea.Quit
	}
	enter, esc, back := k.Code == tea.KeyEnter, k.Code == tea.KeyEscape, k.Code == tea.KeyBackspace
	switch {
	case m.ins.Editing():
		m.ins.Key(k.Text, back, enter, esc)
	case m.pal.Searching():
		m.pal.Key(k.Text, back, enter, esc)
	case esc:
		m.endDrag()
		m.ed.Clear()
	case k.Code == tea.KeyDelete:
		m.ed.Delete()
	case text == "ctrl+z":
		m.ed.Undo()
	case text == "ctrl+y":
		m.ed.Redo()
	}
	return nil
}

// View implements tea.Model.
func (m *Model) View() tea.View {
	v := tea.NewView(m.render())
	v.AltScreen = true
	v.MouseMode = tea.MouseModeAllMotion
	v.WindowTitle = "Cuppa"
	return v
}

func (m *Model) render() string {
	if m.w == 0 || m.h == 0 {
		return ""
	}
	l := m.layout
	pal, stg, ins := m.pal.Lines(), m.stg.Lines(), m.ins.Lines()
	sep := theme.Faded("│")
	out := make([]string, 0, m.h)
	out = append(out, m.titleBar())
	for i := 0; i < l.stage.H; i++ {
		out = append(out, pal[i]+sep+stg[i]+sep+ins[i])
	}
	out = append(out, m.statusBar())
	return strings.Join(out, "\n")
}

func (m *Model) titleBar() string {
	name := m.ed.Document().Name
	if m.ed.Dirty() {
		name += " •"
	}
	return theme.Fit(" "+theme.Title("☕ Cuppa")+theme.Faded("  ·  ")+name, m.w)
}

func (m *Model) statusBar() string {
	hint := "Drag a component from the left bar onto the canvas"
	switch {
	case m.dragging != "":
		hint = "Release over the canvas to place it · Esc cancels"
	case m.stg.Busy():
		hint = "Release to finish"
	default:
		if n, ok := m.ed.Primary(); ok {
			hint = fmt.Sprintf("%s  (%d,%d)  %d×%d", n.Name, n.Rect.X, n.Rect.Y, n.Rect.W, n.Rect.H)
		}
	}
	return theme.Fit(" "+theme.Dim(hint), m.w)
}
