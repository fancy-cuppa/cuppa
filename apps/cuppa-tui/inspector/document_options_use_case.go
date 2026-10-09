package inspector

import (
	"fmt"
	"strconv"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/meta-tui/cuppa/apps/cuppa-tui/theme"
	"github.com/meta-tui/cuppa/libs/canvas/editor"
	"github.com/meta-tui/cuppa/libs/document/design"
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
	b.blank()

	b.text(theme.Bold(" Terminal")).end()
	m.terminalRow(b, doc.Light)
	m.profileRow(b, doc.Profile)
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

// profiles are the colour profiles in the order the arrows step through them.
var profiles = []struct{ value, name string }{
	{"", "true colour"}, {design.Profile256, "256 colours"}, {design.Profile16, "16 colours"}, {design.ProfileNone, "no colour"},
}

// terminalRow chooses whether the canvas previews a dark or a light terminal.
func (m *Model) terminalRow(b *builder, light bool) {
	pick := func(label string, active bool, set bool) {
		text := "[" + label + "]"
		if active {
			b.add(theme.Selected(text), func() { m.ed.SetLight(set) })
		} else {
			b.add(theme.Button(text, true), func() { m.ed.SetLight(set) })
		}
	}
	b.text(" " + theme.Dim("Preview on  "))
	pick("Dark", !light, false)
	b.text(" ")
	pick("Light", light, true)
	b.end()
}

// profileRow steps through the colour profiles with arrows.
func (m *Model) profileRow(b *builder, current string) {
	at := 0
	for i, p := range profiles {
		if p.value == current {
			at = i
		}
	}
	step := func(d int) func() {
		return func() { m.report(m.ed.SetProfile(profiles[(at+d+len(profiles))%len(profiles)].value)) }
	}
	b.text(" " + theme.Dim("Colours  "))
	b.add(theme.Button("◂", true), step(-1)).text(" " + theme.Bold(profiles[at].name) + " ").add(theme.Button("▸", true), step(1)).end()
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
