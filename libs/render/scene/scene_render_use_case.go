// Package scene paints a whole document onto a cell grid.
package scene

import (
	"github.com/fancy-cuppa/cuppa/libs/catalog/definition"
	"github.com/fancy-cuppa/cuppa/libs/document/design"
	"github.com/fancy-cuppa/cuppa/libs/render/grid"
)

// Catalog is what the renderer needs to know about components.
type Catalog interface {
	Get(id string) (definition.Definition, bool)
}

type painter func(g *grid.Grid, p Props)

// painters maps a component id to the function that paints it.
var painters = map[string]painter{
	"lipgloss.box":      paintBox,
	"lipgloss.label":    paintLabel,
	"lipgloss.list":     paintList,
	"lipgloss.tree":     paintTree,
	"lipgloss.table":    paintTable,
	"lipgloss.joinh":    paintJoinH,
	"lipgloss.joinv":    paintJoinV,
	"lipgloss.place":    paintPlace,
	"bubbles.textinput": paintTextInput,
	"bubbles.spinner":   paintSpinner,
	"bubbles.progress":  paintProgress,
}

// Render paints every node of doc, back to front, onto a new grid the size of
// the canvas. Cells no node covers stay unpainted.
func Render(doc design.Document, cat Catalog) *grid.Grid {
	out := grid.New(doc.Width, doc.Height)
	for _, n := range doc.Nodes {
		out.Blit(RenderNode(n, cat), n.Rect.X, n.Rect.Y)
	}
	return out
}

// RenderNode paints one node onto its own grid, sized like the node.
func RenderNode(n design.Node, cat Catalog) *grid.Grid {
	g := grid.New(n.Rect.W, n.Rect.H)
	def, known := cat.Get(n.Component)
	props := Props{}
	if known {
		for k, v := range def.Defaults() {
			props[k] = v
		}
	}
	for k, v := range n.Props {
		props[k] = v
	}
	if p, ok := painters[n.Component]; ok {
		// Every painter owns the whole node rectangle: clear it first so
		// nodes below never show through.
		g.Fill(design.Rect{W: g.W, H: g.H}, ' ', grid.Style{})
		p(g, props)
		return g
	}
	g.Fill(design.Rect{W: g.W, H: g.H}, ' ', grid.Style{})
	name := n.Name
	if known {
		name = def.Name
	}
	paintGeneric(g, name)
	return g
}

// Painted reports whether the component has a dedicated painter.
func Painted(componentID string) bool {
	_, ok := painters[componentID]
	return ok
}
