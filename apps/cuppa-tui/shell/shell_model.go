// Package shell is the root of the terminal app: it lays out the panes,
// translates Bubble Tea messages into pointer events and wires drag and drop
// from the palette to the stage.
package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/meta-tui/cuppa/apps/cuppa-tui/colorpicker"
	"github.com/meta-tui/cuppa/apps/cuppa-tui/fileflow"
	"github.com/meta-tui/cuppa/apps/cuppa-tui/inspector"
	"github.com/meta-tui/cuppa/apps/cuppa-tui/logo"
	"github.com/meta-tui/cuppa/apps/cuppa-tui/menubar"
	"github.com/meta-tui/cuppa/apps/cuppa-tui/modal"
	"github.com/meta-tui/cuppa/apps/cuppa-tui/palette"
	"github.com/meta-tui/cuppa/apps/cuppa-tui/pointer"
	"github.com/meta-tui/cuppa/apps/cuppa-tui/preview"
	"github.com/meta-tui/cuppa/apps/cuppa-tui/stage"
	"github.com/meta-tui/cuppa/apps/cuppa-tui/theme"
	"github.com/meta-tui/cuppa/apps/cuppa-tui/themepicker"
	"github.com/meta-tui/cuppa/libs/canvas/editor"
	"github.com/meta-tui/cuppa/libs/catalog/packstate"
	"github.com/meta-tui/cuppa/libs/catalog/registry"
	"github.com/meta-tui/cuppa/libs/document/design"
	"github.com/meta-tui/cuppa/libs/export/image"
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
	// cat is what every pane reads components from; base is the built-in set
	// it is rebuilt from when packs are added, removed or created.
	cat  *registry.Live
	base *registry.Registry
	// packs says which component packs the palette lists.
	packs *packstate.State
	// userPacks is the folder of installed packs ("" until LoadUserPacks) and
	// packProblems what could not be loaded from it.
	userPacks    string
	embedded     []design.Embedded
	// run is the design being previewed, or nil while editing; pending holds
	// commands it asked for until Update returns them.
	run     *preview.Session
	pending []tea.Cmd
	packProblems []error
	ed    *editor.Editor
	bar   *menubar.Model
	flow  *fileflow.Flow

	pal *palette.Model
	stg *stage.Model
	ins *inspector.Model

	w, h   int
	layout layout

	// owner is the pane a held-button gesture started in; it keeps receiving
	// events even when the pointer leaves it.
	owner pane
	// wantPalette and wantInspector are the side bar widths the user chose (0 for the default).
	wantPalette, wantInspector int
	// command is true on a Mac, where the keys are named Cmd.
	command bool
	// grab is the divider being dragged and hover the one under the pointer.
	grab, hover divider
	// layoutFile is where the widths are remembered, or "" for nowhere.
	layoutFile string
	// dragging is the component being dragged out of the palette, or "".
	dragging string
	mouseX   int
	mouseY   int
}

// New returns the app with an empty design.
func New(cat *registry.Registry) *Model {
	live := registry.NewLive(cat)
	ed := editor.New(live, design.NewDocument("Untitled", defaultWidth, defaultHeight))
	m := &Model{
		cat:   live,
		base:  cat,
		ed:    ed,
		bar:   menubar.New(),
		flow:  fileflow.New(ed, live),
		pal:   palette.New(live),
		packs: packstate.Open(""),
		stg:   stage.New(ed, live),
		ins:   inspector.New(ed, live),
	}
	m.flow.SetOnLoad(m.adoptEmbedded)
	m.ins.BindSnap(m.stg.Snap, m.stg.SetSnap)
	m.ins.BindColorPicker(m.pickColor)
	m.ins.BindThemePicker(m.pickTheme)
	return m
}

// pickColor opens the colour dialog and applies the choice when it ends.
func (m *Model) pickColor(title, current string, apply func(color string)) {
	m.flow.Show(colorpicker.New(title, current), func(o modal.Outcome) {
		if !o.Canceled {
			apply(o.Value)
		}
	})
}

// pickTheme opens the dialog that fills the theme from a colour scheme and
// applies the choice, as one undo step, when it ends.
func (m *Model) pickTheme(apply func(background string, t design.Theme)) {
	m.flow.Show(themepicker.New("Theme from a colour scheme", "Dracula"), func(o modal.Outcome) {
		if o.Canceled {
			return
		}
		if background, t, ok := themepicker.Parse(o.Value); ok {
			apply(background, t)
		}
	})
}

// OpenFile loads a design before the first frame, e.g. from the command line.
func (m *Model) OpenFile(path string) error { return m.flow.OpenPath(path) }

// OpenFileOrNotify loads a design before the first frame, like OpenFile, but
// tells the user in a dialog when it cannot, for front ends with no console.
func (m *Model) OpenFileOrNotify(path string) {
	if err := m.flow.OpenPath(path); err != nil {
		m.flow.Notice("Cannot open file", err.Error())
	}
}

// Welcome shows the first-run notice, if one is due.
func (m *Model) Welcome() {
	dir, err := os.UserConfigDir()
	if err != nil {
		return
	}
	m.flow.Welcome(filepath.Join(dir, "cuppa", "freeze-notice"))
}

// Editor exposes the editor, mainly for tests.
func (m *Model) Editor() *editor.Editor { return m.ed }

// UseIconFont draws the logo in the menu bar from the Cuppa icon font. The
// desktop app and the web page call it, because they ship the font; a
// terminal would show empty boxes.
func (m *Model) UseIconFont() { m.bar.SetIconFont(true) }

// UseCommandKey makes the menus and the shortcuts list say Cmd where they say
// Ctrl. The desktop app and the browser call it on a Mac; their page turns a
// Cmd press into the Ctrl one the editor listens for.
func (m *Model) UseCommandKey() {
	m.command = true
	m.bar.SetCommandKey(true)
}

// Init implements tea.Model.
func (m *Model) Init() tea.Cmd { return nil }

// exportDoneMsg reports a finished background export.
type exportDoneMsg struct{ err error }

// Update implements tea.Model.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.run != nil {
		if handled, cmd := m.run.Update(msg); handled {
			return m, tea.Batch(cmd, m.settle())
		}
	}
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.resize(msg.Width, msg.Height)
	case tea.MouseClickMsg:
		m.route(toEvent(tea.Mouse(msg), pointer.Down))
	case tea.MouseReleaseMsg:
		m.route(toEvent(tea.Mouse(msg), pointer.Up))
	case tea.MouseMotionMsg:
		m.route(toEvent(tea.Mouse(msg), pointer.Move))
	case tea.MouseWheelMsg:
		m.route(toEvent(tea.Mouse(msg), pointer.Wheel))
	case tea.KeyPressMsg:
		m.key(msg)
	case exportDoneMsg:
		m.flow.Finish(msg.err)
	}
	return m, tea.Batch(m.settle(), m.takePending())
}

// settle runs what the last input set in motion: a finished dialog, a
// background export, quitting. It also keeps the menu items' enabled state
// in step with the editor.
func (m *Model) settle() tea.Cmd {
	m.flow.Resolve()
	m.bar.SetEnabled(menubar.EditUndo, m.ed.CanUndo())
	m.bar.SetEnabled(menubar.EditRedo, m.ed.CanRedo())
	hasSel := len(m.ed.Selected()) > 0
	m.bar.SetEnabled(menubar.EditDuplicate, hasSel)
	m.bar.SetEnabled(menubar.EditCopy, hasSel)
	m.bar.SetEnabled(menubar.EditPaste, m.ed.CanPaste())
	m.bar.SetEnabled(menubar.EditDelete, hasSel)
	for _, a := range []menubar.Action{menubar.EditToFront, menubar.EditForward, menubar.EditBackward, menubar.EditToBack} {
		m.bar.SetEnabled(a, hasSel)
	}
	m.bar.SetEnabled(menubar.EditGroup, m.ed.CanGroup())
	m.bar.SetEnabled(menubar.EditUngroup, m.ed.CanUngroup())
	m.bar.SetEnabled(menubar.EditComponent, m.canSaveComponent())
	for _, a := range []menubar.Action{menubar.ExportPNG, menubar.ExportSVG, menubar.ExportWebP} {
		m.bar.SetUnavailable(a, !m.flow.FreezeAvailable())
	}
	if m.flow.Quitting() {
		return tea.Quit
	}
	if job := m.flow.TakeJob(); job != nil {
		return func() tea.Msg { return exportDoneMsg{err: job.Run()} }
	}
	return nil
}

// route sends a pointer event to the open dialog, else the menu bar, else
// the panes.
func (m *Model) route(e pointer.Event) {
	if dlg := m.flow.Modal(); dlg != nil {
		m.mouseX, m.mouseY = e.X, e.Y
		dlg.Handle(e)
		return
	}
	if m.owner == nowhere && m.dragging == "" {
		if act, used := m.bar.Handle(e); used {
			m.perform(act)
			return
		}
	}
	if m.handleDivider(e) {
		return
	}
	m.mouse(e)
}

// perform carries out a menu choice.
func (m *Model) perform(a menubar.Action) {
	if !keepsPreview(a) {
		m.stopPreview()
	}
	switch a {
	case menubar.ViewPreview:
		m.togglePreview()
	case menubar.FileNew:
		m.flow.NewDesign()
	case menubar.FileOpen:
		m.flow.Open()
	case menubar.FileSave:
		m.flow.Save()
	case menubar.FileSaveAs:
		m.flow.SaveAs()
	case menubar.FileQuit:
		m.flow.Quit()
	case menubar.EditUndo:
		m.ed.Undo()
	case menubar.EditRedo:
		m.ed.Redo()
	case menubar.EditCopy:
		m.ed.Copy()
	case menubar.EditPaste:
		m.ed.Paste()
	case menubar.EditDuplicate:
		m.ed.Duplicate()
	case menubar.EditDelete:
		m.ed.Delete()
	case menubar.EditToFront:
		m.ed.Reorder(editor.BringToFront)
	case menubar.EditForward:
		m.ed.Reorder(editor.BringForward)
	case menubar.EditBackward:
		m.ed.Reorder(editor.SendBackward)
	case menubar.EditToBack:
		m.ed.Reorder(editor.SendToBack)
	case menubar.EditGroup:
		m.ed.Group()
	case menubar.EditUngroup:
		m.ed.Ungroup()
	case menubar.EditComponent:
		m.saveAsComponent()
	case menubar.EditPacks:
		m.openPacks()
	case menubar.ExportPNG:
		m.flow.ExportImage(image.PNG)
	case menubar.ExportSVG:
		m.flow.ExportImage(image.SVG)
	case menubar.ExportWebP:
		m.flow.ExportImage(image.WebP)
	case menubar.ExportANSI:
		m.flow.ExportText(true)
	case menubar.ExportText:
		m.flow.ExportText(false)
	case menubar.ExportGo:
		m.flow.ExportGoSource()
	case menubar.HelpShortcuts:
		m.flow.Notice("Shortcuts", shortcutsText(m.command))
	case menubar.HelpAbout:
		m.flow.NoticeWithArt("About Cuppa", aboutText, logo.Dialog())
	}
}

// shortcutsTemplate is the Shortcuts dialog; Ctrl becomes Cmd on a Mac.
const shortcutsTemplate = "Ctrl+N  New            Ctrl+O  Open\n" +
	"Ctrl+S  Save           Ctrl+Shift+S  Save As\n" +
	"Ctrl+Q  Quit           Ctrl+F  Search components\n" +
	"Ctrl+Z  Undo           Ctrl+Y  Redo (also Ctrl+Shift+Z)\n" +
	"Ctrl+C  Copy           Ctrl+V  Paste\n" +
	"Ctrl+D  Duplicate      Del     Delete\n" +
	"Ctrl+G  Group          Ctrl+U  Ungroup\n" +
	"Ctrl+]  Forward one    Ctrl+[  Back one\n" +
	"Ctrl+Shift+]  To front Ctrl+Shift+[  To back\n" +
	"Arrows  Move the selection 1 cell (Shift: 10)\n" +
	"Up/Down in a number: +1 / -1 (Shift: 10)\n" +
	"Esc     Deselect / cancel\n\n" +
	"Everything else is the mouse: drag components from the left onto\n" +
	"the canvas, drag the corners to resize."

func shortcutsText(command bool) string {
	if command {
		return strings.ReplaceAll(shortcutsTemplate, "Ctrl", "Cmd")
	}
	return shortcutsTemplate
}

const aboutText = "A designer for Bubble Tea interfaces, made with Bubble Tea.\n" +
	"Export images need Freeze: github.com/charmbracelet/freeze"

func (m *Model) resize(w, h int) {
	m.w, m.h = w, h
	m.bar.SetWidth(w)
	m.flow.SetScreen(w, h)
	m.relayout()
}

// relayout sizes the panes for the current terminal and chosen widths.
func (m *Model) relayout() {
	m.layout = computeLayout(m.w, m.h, m.wantPalette, m.wantInspector)
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
	if m.run != nil {
		m.previewMouse(e)
		return
	}
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

func (m *Model) key(msg tea.KeyPressMsg) {
	k := msg.Key()
	// Terminals that report the Cmd key (Kitty protocol) send it as Super; it
	// does what Ctrl does, so Cmd+S saves on a Mac.
	if k.Mod&tea.ModSuper != 0 {
		k.Mod = k.Mod&^tea.ModSuper | tea.ModCtrl
	}
	text := k.Keystroke()
	enter, esc, back := k.Code == tea.KeyEnter, k.Code == tea.KeyEscape, k.Code == tea.KeyBackspace
	if dlg := m.flow.Modal(); dlg != nil {
		dlg.Key(k.Text, back, enter, esc)
		return
	}
	switch text {
	case "ctrl+q":
		m.flow.Quit()
		return
	case "ctrl+p":
		m.togglePreview()
		return
	case "ctrl+n":
		m.stopPreview()
		m.flow.NewDesign()
		return
	case "ctrl+o":
		m.stopPreview()
		m.flow.Open()
		return
	case "ctrl+s":
		m.flow.Save()
		return
	case "ctrl+shift+s":
		m.flow.SaveAs()
		return
	}
	switch {
	case m.bar.Open():
		if esc {
			m.bar.Close()
		}
	case m.run != nil:
		m.previewKey(msg)
	case m.ins.Dragging():
		if esc {
			m.ins.CancelDrag()
		}
	case m.ins.Editing():
		switch k.Code {
		case tea.KeyUp:
			m.ins.Step(stepSize(k))
		case tea.KeyDown:
			m.ins.Step(-stepSize(k))
		default:
			m.ins.Key(k.Text, back, enter, esc)
		}
	case m.pal.Searching():
		m.pal.Key(k.Text, back, enter, esc)
	case esc:
		m.endDrag()
		m.ed.Clear()
	case k.Code == tea.KeyDelete:
		m.ed.Delete()
	case text == "ctrl+d":
		m.ed.Duplicate()
	case text == "ctrl+f":
		m.pal.FocusSearch()
	case text == "ctrl+shift+]":
		m.ed.Reorder(editor.BringToFront)
	case text == "ctrl+]":
		m.ed.Reorder(editor.BringForward)
	case text == "ctrl+[":
		m.ed.Reorder(editor.SendBackward)
	case text == "ctrl+shift+[":
		m.ed.Reorder(editor.SendToBack)
	case isArrow(k.Code) && len(m.ed.Selected()) > 0:
		m.nudge(k)
	case text == "ctrl+g":
		m.ed.Group()
	case text == "ctrl+u":
		m.ed.Ungroup()
	case text == "ctrl+c":
		m.ed.Copy()
	case text == "ctrl+v":
		m.ed.Paste()
	case text == "ctrl+z":
		m.ed.Undo()
	case text == "ctrl+y", text == "ctrl+shift+z":
		m.ed.Redo()
	}
}

// stepSize is 1, or 10 with Shift held.
func stepSize(k tea.Key) int {
	if k.Mod&tea.ModShift != 0 {
		return 10
	}
	return 1
}

func isArrow(code rune) bool {
	return code == tea.KeyUp || code == tea.KeyDown || code == tea.KeyLeft || code == tea.KeyRight
}

// nudge moves the selection with an arrow key: 1 cell, or 10 with Shift.
func (m *Model) nudge(k tea.Key) {
	step := stepSize(k)
	dx, dy := 0, 0
	switch k.Code {
	case tea.KeyLeft:
		dx = -step
	case tea.KeyRight:
		dx = step
	case tea.KeyUp:
		dy = -step
	case tea.KeyDown:
		dy = step
	}
	m.ed.Nudge(dx, dy)
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
	left, right := m.separator(leftDivider), m.separator(rightDivider)
	out := make([]string, 0, m.h)
	out = append(out, m.bar.Line(m.titleText()))
	for i := 0; i < l.stage.H; i++ {
		out = append(out, pal[i]+left+stg[i]+right+ins[i])
	}
	out = append(out, m.statusBar())
	if x, drop := m.bar.Dropdown(); drop != nil {
		m.overlay(out, drop, x, 1)
	}
	if dlg := m.flow.Modal(); dlg != nil {
		r := dlg.Rect()
		m.overlay(out, dlg.Lines(), r.X, r.Y)
	}
	return strings.Join(out, "\n")
}

// overlay draws lines over the screen rows starting at (x, y), clipped to the screen.
func (m *Model) overlay(screen, lines []string, x, y int) {
	for i, l := range lines {
		if row := y + i; row >= 0 && row < len(screen) {
			screen[row] = theme.Overlay(screen[row], l, x)
		}
	}
}

// titleText is the design name at the right of the menu bar.
func (m *Model) titleText() string {
	name := m.flow.Title()
	if m.ed.Dirty() {
		name += " •"
	}
	return theme.Dim(name)
}

func (m *Model) statusBar() string {
	hint := "Drag a component from the left bar onto the canvas"
	switch {
	case m.grab != noDivider || m.hover != noDivider:
		hint = "Drag to change the width of the panel"
	case m.dragging != "":
		hint = "Release over the canvas to place it · Esc cancels"
	case m.stg.Busy():
		hint = "Release to finish"
	default:
		if n, ok := m.ed.Primary(); ok {
			hint = fmt.Sprintf("%s  (%d,%d)  %d×%d", n.Name, n.Rect.X, n.Rect.Y, n.Rect.W, n.Rect.H)
		} else if st := m.flow.Status(); st != "" {
			hint = st
		}
	}
	return theme.Fit(" "+theme.Dim(hint), m.w)
}
