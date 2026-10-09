package inspector

import (
	"fmt"
	"strconv"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/meta-tui/cuppa/apps/cuppa-tui/theme"
	"github.com/meta-tui/cuppa/libs/canvas/editor"
)

// documentOptions is what the pane shows when nothing is selected: the canvas
// size, its background colour and the effects applied over the whole design.
func (m *Model) documentOptions(b *builder) {
	doc := m.ed.Document()
	b.text(theme.Bold(" Canvas")).end()
	m.canvasNumeric(b, "W", "canvas-w", doc.Width)
	m.canvasNumeric(b, "H", "canvas-h", doc.Height)
	b.text(" " + theme.Faded(fmt.Sprintf("%d components", len(doc.Nodes)))).end()
	b.blank()

	b.text(theme.Bold(" Background")).end()
	m.backgroundRow(b, doc.Background)
	b.blank()

	b.text(theme.Bold(" Effects")).end()
	m.optionRow(b, "Grid dots", editor.EffectGrid)
	m.optionRow(b, "Shadows", editor.EffectShadow)
	m.optionRow(b, "Scanlines", editor.EffectScanlines)
	m.optionRow(b, "Vignette", editor.EffectVignette)
	if m.message != "" {
		b.blank()
		b.text(" " + lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Warn)).Render(m.message)).end()
	}
	b.blank()
	m.snapRow(b)
}

// canvasNumeric is a size field with [-] [+] and click-to-type, like a node's.
func (m *Model) canvasNumeric(b *builder, label, field string, v int) {
	b.text(" " + theme.Dim(label+"  "))
	b.add(theme.Button("[-]", true), func() { m.nudgeCanvas(field, -1) }).text(" ")
	if m.editing == field {
		b.text(m.buf + "█")
	} else {
		b.add(theme.Bold(fmt.Sprintf("%4d", v)), func() { m.startEdit(field, strconv.Itoa(v)) })
	}
	b.text(" ").add(theme.Button("[+]", true), func() { m.nudgeCanvas(field, 1) }).end()
}

func (m *Model) nudgeCanvas(field string, delta int) {
	doc := m.ed.Document()
	m.setCanvas(field, map[string]int{"canvas-w": doc.Width, "canvas-h": doc.Height}[field]+delta)
}

func (m *Model) setCanvas(field string, v int) {
	doc := m.ed.Document()
	w, h := doc.Width, doc.Height
	if field == "canvas-w" {
		w = v
	} else {
		h = v
	}
	m.report(m.ed.SetCanvasSize(w, h))
}

func (m *Model) commitCanvas(field, value string) {
	v, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		m.message = "enter a whole number"
		return
	}
	m.setCanvas(field, v)
}

// backgroundRow shows the colour as a swatch; clicking opens the picker (or a
// text box when no picker is bound).
func (m *Model) backgroundRow(b *builder, value string) {
	if m.editing == "background" {
		b.text(" " + m.buf + "█").end()
		return
	}
	swatch := theme.Faded("--")
	if value != "" {
		swatch = lipgloss.NewStyle().Background(lipgloss.Color(value)).Render("  ")
	}
	shown := value
	if shown == "" {
		shown = "terminal default"
	}
	open := func() { m.startEdit("background", value) }
	if m.pickColor != nil {
		open = func() {
			m.pickColor("Background", value, func(v string) { m.report(m.ed.SetBackground(v)) })
		}
	}
	b.text(" " + swatch + " ").add(shown+theme.Faded(" ✎"), open).end()
}

// optionRow is a checkbox for one canvas option.
func (m *Model) optionRow(b *builder, label, name string) {
	mark := "[ ]"
	if m.ed.Effect(name) {
		mark = "[x]"
	}
	b.text(" ").add(theme.Button(mark+" "+label, true), func() {
		m.report(m.ed.SetEffect(name, !m.ed.Effect(name)))
	}).end()
}
