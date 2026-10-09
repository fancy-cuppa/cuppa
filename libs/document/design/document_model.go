package design

import "fmt"

// Document is a design: a fixed-size canvas and the nodes placed on it.
// Nodes are ordered back to front, so the last node is drawn on top.
type Document struct {
	Name   string `json:"name"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Nodes  []Node `json:"nodes"`
	// Background is the canvas colour ("" for the terminal's own, "0" to "255"
	// or "#rrggbb"). It fills every cell no component colours itself.
	Background string `json:"background,omitempty"`
	// HideGrid turns off the dotted grid shown on an empty canvas in the editor.
	// The grid is never part of an export.
	HideGrid bool    `json:"hideGrid,omitempty"`
	Effects  Effects `json:"effects,omitempty"`
	// Seq is the counter behind generated node ids; it only ever grows.
	Seq int `json:"seq"`
}

// Effects are looks applied over the whole canvas, in the editor and in exports.
type Effects struct {
	// Shadow drops a one-cell shadow below and to the right of every component.
	Shadow bool `json:"shadow,omitempty"`
	// Scanlines dims every other row, like a CRT.
	Scanlines bool `json:"scanlines,omitempty"`
	// Vignette dims the edges of the canvas.
	Vignette bool `json:"vignette,omitempty"`
}

// Canvas size limits, in cells.
const (
	MaxWidth  = 400
	MaxHeight = 200
)

// NewDocument returns an empty document of the given canvas size.
func NewDocument(name string, width, height int) Document {
	return Document{Name: name, Width: width, Height: height}
}

// Bounds is the canvas as a rectangle.
func (d Document) Bounds() Rect { return Rect{W: d.Width, H: d.Height} }

// Clone returns a deep copy, so edits to it never reach the original.
func (d Document) Clone() Document {
	nodes := make([]Node, len(d.Nodes))
	for i, n := range d.Nodes {
		nodes[i] = n.Clone()
	}
	d.Nodes = nodes
	return d
}

// NewID returns an id no node has had in this document.
func (d *Document) NewID() NodeID {
	d.Seq++
	return NodeID(fmt.Sprintf("n%d", d.Seq))
}

// Add appends the node on top. A node without an id gets a fresh one.
// The stored node is returned.
func (d *Document) Add(n Node) Node {
	if n.ID == "" {
		n.ID = d.NewID()
	}
	d.Nodes = append(d.Nodes, n)
	return n
}

// Index is the position of the node in z-order (0 is the back), or -1.
func (d Document) Index(id NodeID) int {
	for i, n := range d.Nodes {
		if n.ID == id {
			return i
		}
	}
	return -1
}

// Get returns a copy of the node with the given id.
func (d Document) Get(id NodeID) (Node, bool) {
	if i := d.Index(id); i >= 0 {
		return d.Nodes[i], true
	}
	return Node{}, false
}

// Update applies fn to the node with the given id and reports whether it exists.
func (d *Document) Update(id NodeID, fn func(*Node)) bool {
	i := d.Index(id)
	if i < 0 {
		return false
	}
	fn(&d.Nodes[i])
	return true
}

// Remove deletes the node and reports whether it existed.
func (d *Document) Remove(id NodeID) bool {
	i := d.Index(id)
	if i < 0 {
		return false
	}
	d.Nodes = append(d.Nodes[:i], d.Nodes[i+1:]...)
	return true
}

// Raise moves the node one step toward the front.
func (d *Document) Raise(id NodeID) bool { return d.moveTo(id, d.Index(id)+1) }

// Lower moves the node one step toward the back.
func (d *Document) Lower(id NodeID) bool { return d.moveTo(id, d.Index(id)-1) }

// ToFront moves the node above every other node.
func (d *Document) ToFront(id NodeID) bool { return d.moveTo(id, len(d.Nodes)-1) }

// ToBack moves the node below every other node.
func (d *Document) ToBack(id NodeID) bool { return d.moveTo(id, 0) }

// MoveToIndex places the node at position to in z-order (0 is the back) and
// reports whether the order changed.
func (d *Document) MoveToIndex(id NodeID, to int) bool { return d.moveTo(id, to) }

// moveTo places the node at position to in z-order and reports whether the
// order changed.
func (d *Document) moveTo(id NodeID, to int) bool {
	from := d.Index(id)
	if from < 0 || to < 0 || to >= len(d.Nodes) || to == from {
		return false
	}
	n := d.Nodes[from]
	d.Nodes = append(d.Nodes[:from], d.Nodes[from+1:]...)
	d.Nodes = append(d.Nodes[:to], append([]Node{n}, d.Nodes[to:]...)...)
	return true
}
