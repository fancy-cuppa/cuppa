package gosource

import (
	"github.com/meta-tui/cuppa/libs/catalog/definition"
	"github.com/meta-tui/cuppa/libs/document/design"
)

// Catalog is what the generator needs to know about components.
type Catalog interface {
	Get(id string) (definition.Definition, bool)
}

// maxDepth stops a component that contains itself from expanding forever.
const maxDepth = 8

// leaf is one component the generated program places: a real component with
// its absolute position and final property values.
type leaf struct {
	Kind  string
	Name  string
	Rect  design.Rect
	Props map[string]string
}

// leaves flattens a document back to front: groups are replaced by their
// parts, and components made from other components (pack components) by the
// parts they are made of, scaled to the size they were placed at.
func leaves(doc design.Document, cat Catalog) []leaf {
	return expand(doc.Nodes, design.Rect{W: doc.Width, H: doc.Height}, doc.Width, doc.Height, 0, nil, cat)
}

// expand places nodes whose rectangles are relative to a box of baseW by baseH
// that sits at origin; override carries values a pack component hands down.
func expand(nodes []design.Node, origin design.Rect, baseW, baseH, depth int, override func(design.Node) map[string]string, cat Catalog) []leaf {
	var out []leaf
	for _, n := range nodes {
		if n.Hidden {
			continue
		}
		rect := n.Rect
		if origin.W != baseW || origin.H != baseH {
			rect = rect.Scale(baseW, baseH, origin.W, origin.H)
		}
		rect = rect.Translate(origin.X, origin.Y)
		props := withDefaults(n, cat)
		if override != nil {
			for k, v := range override(n) {
				props[k] = v
			}
		}
		switch {
		case n.IsGroup() && depth < maxDepth:
			out = append(out, expand(n.Children, rect, n.BaseW, n.BaseH, depth+1, nil, cat)...)
		case depth < maxDepth && isComposite(n, cat):
			def, _ := cat.Get(n.Component)
			inner := *def.Inner
			hand := func(part design.Node) map[string]string {
				vals := map[string]string{}
				for _, e := range inner.Props {
					if e.Target != part.ID {
						continue
					}
					if v, ok := props[e.Key]; ok {
						vals[e.TargetProp] = v
					} else {
						vals[e.TargetProp] = e.Default
					}
				}
				return vals
			}
			out = append(out, expand(inner.Nodes, rect, inner.W, inner.H, depth+1, hand, cat)...)
		default:
			out = append(out, leaf{Kind: n.Component, Name: n.Name, Rect: rect, Props: props})
		}
	}
	return out
}

func isComposite(n design.Node, cat Catalog) bool {
	def, ok := cat.Get(n.Component)
	return ok && def.Inner != nil
}

// withDefaults is the node's properties over the catalog defaults.
func withDefaults(n design.Node, cat Catalog) map[string]string {
	props := map[string]string{}
	if def, ok := cat.Get(n.Component); ok {
		for k, v := range def.Defaults() {
			props[k] = v
		}
	}
	for k, v := range n.Props {
		props[k] = v
	}
	return props
}
