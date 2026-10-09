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
}

func (b *builder) add(styled string, act func()) *builder {
	w := lipgloss.Width(styled)
	if act != nil {
		b.regions = append(b.regions, region{y: len(b.lines), x0: b.x, x1: b.x + w, act: act})
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
	for _, line := range wrap(text, m.w-indent-1) {
		b.text(strings.Repeat(" ", indent))
		b.add(style(line), act)
		b.end()
	}
}

func plain(s string) string { return s }

// Lines renders the pane as exactly h lines of w cells.
func (m *Model) Lines() []string {
	b := &builder{}
	b.text(theme.Title(" DETAILS")).end()
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
	m.total, m.regions = len(b.lines), b.regions
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
	doc := m.ed.Document()
	b.text(theme.Dim(" Nothing selected.")).end()
	b.text(theme.Faded(" Drag a component from")).end()
	b.text(theme.Faded(" the left bar onto the")).end()
	b.text(theme.Faded(" canvas.")).end()
	b.blank()
	b.text(theme.Bold(" Canvas")).end()
	b.text(fmt.Sprintf(" %d × %d cells", doc.Width, doc.Height)).end()
	b.text(fmt.Sprintf(" %d components", len(doc.Nodes))).end()
	b.blank()
	m.snapRow(b)
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
	m.numeric(b, "X", "x", n.Rect.X)
	m.numeric(b, "Y", "y", n.Rect.Y)
	m.numeric(b, "W", "w", n.Rect.W)
	m.numeric(b, "H", "h", n.Rect.H)
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

func (m *Model) numeric(b *builder, label, field string, v int) {
	b.text(" " + theme.Dim(label+"  "))
	b.add(theme.Button("[-]", true), func() { m.nudge(field, -1) }).text(" ")
	if m.editing == field {
		b.text(m.buf + "█")
	} else {
		b.add(theme.Bold(fmt.Sprintf("%4d", v)), func() { m.startEdit(field, strconv.Itoa(v)) })
	}
	b.text(" ").add(theme.Button("[+]", true), func() { m.nudge(field, 1) }).end()
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
}

// properties lists the component-specific properties from the catalog schema.
func (m *Model) properties(b *builder, n design.Node) {
	def, ok := m.cat.Get(n.Component)
	if !ok || len(def.Props) == 0 {
		return
	}
	b.blank()
	b.text(theme.Bold(" Properties")).end()
	for _, p := range def.Props {
		value, set := n.Props[p.Key]
		if !set {
			value = p.Default
		}
		m.block(b, p.Label, 1, theme.Dim, nil)
		if p.Kind == definition.PropText {
			m.textProp(b, "prop:"+p.Key, value)
			continue
		}
		b.text("  ")
		m.propValue(b, n.ID, p, value)
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
		b.add(theme.Button("◂", true), set(prev)).text(" "+theme.Bold(value)+" ").add(theme.Button("▸", true), set(next))
	case definition.PropInt:
		n, _ := strconv.Atoi(value)
		b.add(theme.Button("[-]", true), set(strconv.Itoa(n-1))).text(" ")
		if m.editing == field {
			b.text(m.buf + "█")
		} else {
			b.add(theme.Bold(value), func() { m.startEdit(field, value) })
		}
		b.text(" ").add(theme.Button("[+]", true), set(strconv.Itoa(n+1)))
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
