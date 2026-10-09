package scene

import (
	"github.com/meta-tui/cuppa/libs/document/design"
	"github.com/meta-tui/cuppa/libs/render/grid"
)

// maxCompositeDepth stops a component that contains itself, directly or
// through others, from recursing forever.
const maxCompositeDepth = 8

// paintComposite paints the parts of a user-made component onto g. Parts are
// scaled from the component's default size to the size it was placed at, and
// the component's own property values are handed down to the parts they drive.
func paintComposite(g *grid.Grid, c design.Composite, props Props, cat Catalog, depth int) {
	if depth >= maxCompositeDepth {
		paintGeneric(g, c.Name)
		return
	}
	for _, part := range c.Nodes {
		if part.Hidden {
			continue
		}
		part = part.Clone()
		if part.Props == nil {
			part.Props = map[string]string{}
		}
		for _, e := range c.Props {
			if e.Target != part.ID {
				continue
			}
			if v, ok := props[e.Key]; ok {
				part.Props[e.TargetProp] = v
			} else {
				part.Props[e.TargetProp] = e.Default
			}
		}
		part.Rect = scaleRect(part.Rect, c.W, c.H, g.W, g.H)
		g.Blit(renderNode(part, cat, depth+1), part.Rect.X, part.Rect.Y)
	}
}

// scaleRect maps r from a box of size fromW x fromH onto one of toW x toH,
// keeping every edge on a whole cell and the result at least 1 x 1.
func scaleRect(r design.Rect, fromW, fromH, toW, toH int) design.Rect {
	x0, x1 := r.X*toW/fromW, (r.X+r.W)*toW/fromW
	y0, y1 := r.Y*toH/fromH, (r.Y+r.H)*toH/fromH
	return design.Rect{X: x0, Y: y0, W: max(x1-x0, 1), H: max(y1-y0, 1)}
}
