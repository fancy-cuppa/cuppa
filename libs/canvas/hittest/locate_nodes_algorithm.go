// Package hittest answers "what is under this cell" for the canvas.
package hittest

import "github.com/meta-tui/cuppa/libs/document/design"

// Node returns the topmost visible node covering the cell (x, y). Hidden
// nodes cannot be hit.
func Node(doc design.Document, x, y int) (design.NodeID, bool) {
	for i := len(doc.Nodes) - 1; i >= 0; i-- {
		if !doc.Nodes[i].Hidden && doc.Nodes[i].Rect.Contains(x, y) {
			return doc.Nodes[i].ID, true
		}
	}
	return "", false
}

// NodesIn returns the visible nodes whose rectangle touches r, back to front.
func NodesIn(doc design.Document, r design.Rect) []design.NodeID {
	var out []design.NodeID
	for _, n := range doc.Nodes {
		if !n.Hidden && n.Rect.Intersects(r) {
			out = append(out, n.ID)
		}
	}
	return out
}
