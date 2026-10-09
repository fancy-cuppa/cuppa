package stage

import (
	"github.com/fancy-cuppa/cuppa/apps/cuppa-tui/theme"
	"github.com/fancy-cuppa/cuppa/libs/document/design"
	"github.com/fancy-cuppa/cuppa/libs/render/grid"
	"github.com/fancy-cuppa/cuppa/libs/render/scene"
)

// render draws the visible part of the canvas, the selection and the drag
// feedback onto a grid the size of the pane.
func (m *Model) render() *grid.Grid {
	doc := m.ed.Document()
	rendered := scene.Render(doc, m.cat)
	view := grid.New(m.w, m.h)
	canvasBg := grid.Style{Bg: theme.Canvas}
	dot := grid.Style{Fg: theme.Faint, Bg: theme.Canvas}
	for vy := 0; vy < m.h; vy++ {
		for vx := 0; vx < m.w; vx++ {
			cx, cy := m.Canvas(vx, vy)
			if cx >= doc.Width || cy >= doc.Height {
				continue
			}
			c := rendered.At(cx, cy)
			switch {
			case c.Ch != 0:
				view.Set(vx, vy, c)
			case cx%4 == 0 && cy%2 == 0:
				view.Set(vx, vy, grid.Cell{Ch: '·', Style: dot})
			default:
				view.Set(vx, vy, grid.Cell{Ch: ' ', Style: canvasBg})
			}
		}
	}
	for _, id := range m.ed.Selected() {
		if n, ok := doc.Get(id); ok {
			m.drawSelection(view, n.Rect, len(m.ed.Selected()) == 1)
		}
	}
	m.drawGuides(view)
	if m.mode == marquee {
		m.outline(view, m.marqueeRect(), '·', grid.Style{Fg: theme.Accent})
	}
	if m.ghost != nil {
		m.outline(view, *m.ghost, '░', grid.Style{Fg: theme.Accent})
	}
	return view
}

// toView converts a canvas rectangle to pane coordinates.
func (m *Model) toView(r design.Rect) design.Rect { return r.Translate(-m.offX, -m.offY) }

// drawSelection recolours the outline of a selected node and, for a single
// selection, marks the resize handles at the corners.
func (m *Model) drawSelection(view *grid.Grid, canvas design.Rect, handles bool) {
	r := m.toView(canvas)
	accent := func(s grid.Style) grid.Style { s.Fg, s.Bold = theme.Accent, true; return s }
	for x := r.X; x < r.Right(); x++ {
		view.Restyle(x, r.Y, accent)
		view.Restyle(x, r.Bottom()-1, accent)
	}
	for y := r.Y; y < r.Bottom(); y++ {
		view.Restyle(r.X, y, accent)
		view.Restyle(r.Right()-1, y, accent)
	}
	if !handles {
		return
	}
	mark := func(x, y int) {
		view.Set(x, y, grid.Cell{Ch: '■', Style: grid.Style{Fg: theme.Accent, Bold: true}})
	}
	switch {
	case r.H == 1 && r.W >= 3:
		mark(r.X, r.Y)
		mark(r.Right()-1, r.Y)
	case r.H >= 2 && r.W >= 2:
		mark(r.X, r.Y)
		mark(r.Right()-1, r.Y)
		mark(r.X, r.Bottom()-1)
		mark(r.Right()-1, r.Bottom()-1)
	}
}

// outline draws a one-cell frame of ch around the canvas rectangle.
func (m *Model) outline(view *grid.Grid, canvas design.Rect, ch rune, s grid.Style) {
	r := m.toView(canvas)
	for x := r.X; x < r.Right(); x++ {
		view.Set(x, r.Y, grid.Cell{Ch: ch, Style: s})
		view.Set(x, r.Bottom()-1, grid.Cell{Ch: ch, Style: s})
	}
	for y := r.Y; y < r.Bottom(); y++ {
		view.Set(r.X, y, grid.Cell{Ch: ch, Style: s})
		view.Set(r.Right()-1, y, grid.Cell{Ch: ch, Style: s})
	}
}

// drawGuides draws the snap guides over empty canvas cells only, so they never
// hide a component.
func (m *Model) drawGuides(view *grid.Grid) {
	guide := grid.Style{Fg: theme.Accent, Dim: true}
	blank := func(x, y int) bool {
		c := view.At(x, y)
		return c.Ch == 0 || c.Ch == ' ' || c.Ch == '·'
	}
	for _, gd := range m.guides {
		if gd.Vertical {
			x := gd.Pos - m.offX
			for y := 0; y < view.H; y++ {
				if blank(x, y) {
					view.Set(x, y, grid.Cell{Ch: '┆', Style: guide})
				}
			}
			continue
		}
		y := gd.Pos - m.offY
		for x := 0; x < view.W; x++ {
			if blank(x, y) {
				view.Set(x, y, grid.Cell{Ch: '┄', Style: guide})
			}
		}
	}
}
