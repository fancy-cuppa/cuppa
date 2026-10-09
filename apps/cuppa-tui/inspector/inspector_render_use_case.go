package inspector

import (
	"fmt"
	"slices"
	"strconv"

	"charm.land/lipgloss/v2"
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
		b.text(" " + theme.Dim(p.Label)).end()
		b.text("  ")
		m.propValue(b, n.ID, p, value)
		b.end()
	}
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
			b.text(swatch+" ").add(value+theme.Faded(" ✎"), func() { m.startEdit(field, value) })
		}
	default:
		m.editable(b, field, value, value)
	}
}

// layers lists nodes front to back; clicking one selects it.
func (m *Model) layers(b *builder) {
	doc := m.ed.Document()
	b.blank()
	b.text(theme.Bold(" Layers")).end()
	if len(doc.Nodes) == 0 {
		b.text(theme.Faded(" (none)")).end()
		return
	}
	for i := len(doc.Nodes) - 1; i >= 0; i-- {
		n := doc.Nodes[i]
		id := n.ID
		label := " " + n.Name
		if m.ed.IsSelected(id) {
			b.add(theme.Selected(theme.Fit(label, m.w)), func() { m.ed.Select(id) })
		} else {
			b.add(label, func() { m.ed.Select(id) })
		}
		b.end()
	}
}
