package inspector

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/meta-tui/cuppa/apps/cuppa-tui/theme"
	"github.com/meta-tui/cuppa/libs/canvas/editor"
	"github.com/meta-tui/cuppa/libs/catalog/definition"
	"github.com/meta-tui/cuppa/libs/document/design"
)

// builder assembles content lines and records which spans are clickable.
type builder struct {
	lines   []string
	regions []region
	cur     string
	x       int
	// focus is the keyboard stop to draw as focused (-1 for none) and stops
	// counts the stops added so far.
	focus, stops int
}

// spoken is how a span reads aloud: the text before it on its line (or, when
// there is none, the line above, which names a property) and its own text,
// without the stepper buttons.
func (b *builder) spoken(styled string) string {
	before := ansi.Strip(b.cur)
	if strings.TrimSpace(before) == "" && len(b.lines) > 0 {
		before = ansi.Strip(b.lines[len(b.lines)-1])
	}
	text := before + " " + ansi.Strip(styled)
	for _, junk := range []string{"[-]", "[+]", "◂", "▸", "█"} {
		text = strings.ReplaceAll(text, junk, "")
	}
	return strings.Join(strings.Fields(text), " ")
}

// add is a clickable span that is also a stop for the keyboard.
func (b *builder) add(styled string, act func()) *builder {
	return b.span(styled, act, nil, false)
}

// addStep is a stop that Left and Right (or - and +) also change by step(±1).
func (b *builder) addStep(styled string, act func(), step func(delta int)) *builder {
	return b.span(styled, act, step, false)
}

// addMouse is a clickable span the keyboard does not stop on, because
// another stop does what it does (the [-] and [+] of a number).
func (b *builder) addMouse(styled string, act func()) *builder {
	return b.span(styled, act, nil, true)
}

func (b *builder) span(styled string, act func(), step func(int), mouseOnly bool) *builder {
	w := lipgloss.Width(styled)
	label := b.spoken(styled)
	if act != nil {
		if !mouseOnly {
			if b.stops == b.focus {
				styled = theme.Selected(ansi.Strip(styled))
			}
			b.stops++
		}
		b.regions = append(b.regions, region{y: len(b.lines), x0: b.x, x1: b.x + w, act: act, step: step, mouseOnly: mouseOnly, label: label})
	}
	b.cur += styled
	b.x += w
	return b
}

func (b *builder) text(s string) *builder { return b.add(s, nil) }

func (b *builder) end() {
	b.lines = append(b.lines, b.cur)
	b.cur, b.x = "", 0
}

func (b *builder) blank() { b.end() }

// wrap splits text into lines of at most width cells, breaking at spaces,
// commas and slashes, so a long list or path stays inside the pane.
func wrap(text string, width int) []string {
	return strings.Split(ansi.Wrap(text, max(width, 4), ",/"), "\n")
}

// block writes text wrapped to the pane, each line indented. When act is set,
// every wrapped line is clickable.
func (m *Model) block(b *builder, text string, indent int, style func(string) string, act func()) {
	for i, line := range wrap(text, m.w-indent-1) {
		b.text(strings.Repeat(" ", indent))
		// One property is one stop for the keyboard, however many lines it wraps to.
		if i == 0 {
			b.add(style(line), act)
		} else {
			b.addMouse(style(line), act)
		}
		b.end()
	}
}

func plain(s string) string { return s }

// Lines renders the pane as exactly h lines of w cells.
func (m *Model) Lines() []string {
	b := &builder{focus: m.focusStop()}
	title := theme.Title(" DETAILS")
	if m.focused {
		title = theme.Selected(" DETAILS ")
	}
	b.text(title).end()
	m.history(b)
	b.blank()
	if n, ok := m.ed.Primary(); !ok {
		m.empty(b)
	} else if len(m.ed.Selected()) > 1 {
		m.group(b)
	} else {
		m.single(b, n)
	}
	m.layers(b)
	m.total, m.regions, m.stops = len(b.lines), b.regions, b.stops
	m.scroll = min(max(m.scroll, 0), max(m.total-m.h, 0))
	end := min(m.scroll+m.h, len(b.lines))
	return theme.Block(b.lines[m.scroll:end], m.w, m.h)
}

func (m *Model) history(b *builder) {
	b.text(" ").
		add(theme.Button("[Undo]", m.ed.CanUndo()), func() { m.ed.Undo() }).text(" ").
		add(theme.Button("[Redo]", m.ed.CanRedo()), func() { m.ed.Redo() }).end()
}

func (m *Model) empty(b *builder) {
	b.text(theme.Dim(" Nothing selected.")).end()
	b.text(theme.Faded(" Drag a component from")).end()
	b.text(theme.Faded(" the left bar onto the")).end()
	b.text(theme.Faded(" canvas.")).end()
	b.blank()
	m.documentOptions(b)
}

// snapRow is the "snap to guides" checkbox.
func (m *Model) snapRow(b *builder) {
	if m.snapGet == nil {
		return
	}
	mark := "[ ]"
	if m.snapGet() {
		mark = "[x]"
	}
	b.text(" ").add(theme.Button(mark+" Snap to guides", true), func() { m.snapSet(!m.snapGet()) }).end()
}

func (m *Model) group(b *builder) {
	b.text(theme.Bold(fmt.Sprintf(" %d selected", len(m.ed.Selected())))).end()
	b.blank()
	m.order(b)
	m.actions(b)
	b.blank()
	m.snapRow(b)
}

func (m *Model) single(b *builder, n design.Node) {
	b.text(" ").add(theme.Dim("Name  "), nil)
	m.editable(b, "name", n.Name, n.Name)
	b.end()
	b.text(" ").add(theme.Dim("Type  "), nil).text(n.Component).end()
	if n.Locked {
		b.text(" " + theme.Faded("Locked: unlock it in Layers to edit")).end()
	}
	b.blank()
	l := n.Layout
	m.numeric(b, n.ID, "X", "x", n.Rect.X, l.X)
	m.numeric(b, n.ID, "Y", "y", n.Rect.Y, l.Y)
	m.numeric(b, n.ID, "W", "w", n.Rect.W, l.W)
	m.numeric(b, n.ID, "H", "h", n.Rect.H, l.H)
	b.blank()
	doc := m.ed.Document()
	b.text(fmt.Sprintf(" Layer %d of %d", doc.Index(n.ID)+1, len(doc.Nodes))).end()
	m.order(b)
	m.actions(b)
	b.blank()
	m.snapRow(b)
	m.properties(b, n)
	if m.message != "" {
		b.blank()
		b.text(" " + lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Warn)).Render(m.message)).end()
	}
}

// editable renders a value that turns into a text box when clicked.
func (m *Model) editable(b *builder, field, value, initial string) {
	if m.editing == field {
		b.text(m.buf + "█")
		return
	}
	if value == "" {
		value = theme.Faded("(empty)")
	}
	b.add(value, func() { m.startEdit(field, initial) })
}

// numeric is one geometry row: [-] value [+] and a unit button. A value that
// follows the canvas shows its expression ("50%", "100% - 10"); the button
// switches the axis between fixed cells and a percentage of the canvas.
func (m *Model) numeric(b *builder, id design.NodeID, label, field string, v int, expression string) {
	shown, initial := fmt.Sprintf("%4d", v), strconv.Itoa(v)
	if expression != "" {
		shown, initial = fmt.Sprintf("%4s", expression), expression
	}
	b.text(" " + theme.Dim(label+"  "))
	b.addMouse(theme.Button("[-]", true), func() { m.nudge(field, -1) }).text(" ")
	if m.editing == field {
		b.text(m.buf + "█")
	} else {
		b.addStep(theme.Bold(shown), func() { m.startEdit(field, initial) },
			func(d int) { m.nudge(field, d) })
	}
	b.text(" ").addMouse(theme.Button("[+]", true), func() { m.nudge(field, 1) })
	unit := "[%]"
	if expression != "" {
		unit = "[#]"
	}
	b.text(" ").add(theme.Button(unit, true), func() { m.report(m.ed.ToggleLayoutUnit(id, field)) })
	b.end()
}

func (m *Model) order(b *builder) {
	b.text(" ").
		add(theme.Button("[Front]", true), func() { m.ed.Reorder(editor.BringToFront) }).text(" ").
		add(theme.Button("[Fwd]", true), func() { m.ed.Reorder(editor.BringForward) }).end()
	b.text(" ").
		add(theme.Button("[Bwd]", true), func() { m.ed.Reorder(editor.SendBackward) }).text(" ").
		add(theme.Button("[Back]", true), func() { m.ed.Reorder(editor.SendToBack) }).end()
}

func (m *Model) actions(b *builder) {
	b.text(" ").
		add(theme.Button("[Duplicate]", true), func() { m.ed.Duplicate() }).text(" ").
		add(theme.Button("[Delete]", true), func() { m.ed.Delete() }).end()
	b.text(" ").
		add(theme.Button("[Group]", m.ed.CanGroup()), func() { m.ed.Group() }).text(" ").
		add(theme.Button("[Ungroup]", m.ed.CanUngroup()), func() { m.ed.Ungroup() }).end()
}

// properties lists the component-specific properties from the catalog schema.
func (m *Model) properties(b *builder, n design.Node) {
	def, ok := m.cat.Get(n.Component)
	if !ok || len(def.Props) == 0 {
		return
	}
	b.blank()
	b.text(theme.Bold(" Properties")).end()
	doc := m.ed.Document()
	effective := def.Effective(n.Props, doc.Theme, doc.Background)
	for _, p := range def.Props {
		_, own := n.Props[p.Key]
		value := effective[p.Key]
		m.block(b, p.Label, 1, theme.Dim, nil)
		if p.Kind == definition.PropText {
			m.textProp(b, "prop:"+p.Key, value)
			continue
		}
		b.text("  ")
		m.propValue(b, n.ID, p, value)
		if p.Role != definition.RoleNone {
			m.themeMark(b, n.ID, p, own)
		}
		b.end()
	}
}

// textProp shows a text property wrapped to the pane; clicking any line edits it.
func (m *Model) textProp(b *builder, field, value string) {
	if m.editing == field {
		m.block(b, m.buf+"█", 2, plain, nil)
		return
	}
	edit := func() { m.startEdit(field, value) }
	if value == "" {
		m.block(b, "(empty)", 2, theme.Faded, edit)
		return
	}
	m.block(b, value, 2, plain, edit)
}

func (m *Model) propValue(b *builder, id design.NodeID, p definition.PropSpec, value string) {
	field := "prop:" + p.Key
	set := func(v string) func() {
		return func() { m.report(m.ed.SetProp(id, p.Key, v)) }
	}
	switch p.Kind {
	case definition.PropBool:
		mark := "[ ]"
		if value == "true" {
			mark = "[x]"
		}
		next := "true"
		if value == "true" {
			next = "false"
		}
		b.add(theme.Button(mark, true), set(next))
	case definition.PropChoice:
		i := slices.Index(p.Choices, value)
		prev := p.Choices[(i+len(p.Choices)-1)%len(p.Choices)]
		next := p.Choices[(i+1)%len(p.Choices)]
		b.addMouse(theme.Button("◂", true), set(prev)).text(" ").
			addStep(theme.Bold(value), set(next), func(d int) {
				set(p.Choices[((i+d)%len(p.Choices)+len(p.Choices))%len(p.Choices)])()
			}).text(" ").addMouse(theme.Button("▸", true), set(next))
	case definition.PropInt:
		n, _ := strconv.Atoi(value)
		b.addMouse(theme.Button("[-]", true), set(strconv.Itoa(n-1))).text(" ")
		if m.editing == field {
			b.text(m.buf + "█")
		} else {
			b.addStep(theme.Bold(value), func() { m.startEdit(field, value) },
				func(d int) { set(strconv.Itoa(n + d))() })
		}
		b.text(" ").addMouse(theme.Button("[+]", true), set(strconv.Itoa(n+1)))
	case definition.PropColor:
		if m.editing == field {
			b.text(m.buf + "█")
		} else {
			swatch := lipgloss.NewStyle().Background(lipgloss.Color(value)).Render("  ")
			if value == "" {
				swatch = theme.Faded("--")
			}
			open := func() { m.startEdit(field, value) }
			if m.pickColor != nil {
				open = func() { m.pickColor(p.Label, value, func(v string) { set(v)() }) }
			}
			b.text(swatch+" ").add(value+theme.Faded(" ✎"), open)
		}
	default:
		m.editable(b, field, value, value)
	}
}

// themeMark says whether a colour follows the theme or is the component's own,
// and for an own one offers to go back to the theme.
func (m *Model) themeMark(b *builder, id design.NodeID, p definition.PropSpec, own bool) {
	if !own {
		b.text(" " + theme.Faded("theme"))
		return
	}
	b.text(" ").add(theme.Button("[theme]", true), func() { m.ed.ClearProp(id, p.Key) })
}
