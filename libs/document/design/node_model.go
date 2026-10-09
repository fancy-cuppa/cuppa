package design

// NodeID identifies one component instance inside a document.
type NodeID string

// Node is one placed component: which catalog entry it is, where it sits and
// the property values the designer changed.
type Node struct {
	ID        NodeID            `json:"id"`
	Component string            `json:"component"`
	Name      string            `json:"name"`
	Rect      Rect              `json:"rect"`
	Props     map[string]string `json:"props,omitempty"`
	// Hidden nodes are not drawn, exported or hit by the pointer; they stay in
	// the layer list.
	Hidden bool `json:"hidden,omitempty"`
	// Locked nodes are drawn but cannot be moved, resized, deleted or have their
	// properties changed.
	Locked bool `json:"locked,omitempty"`
	// Children are set on a group: the components it holds, with rectangles
	// relative to the group's top-left, laid out for a box of BaseW x BaseH. A
	// placed group scales them with its own size.
	Children []Node `json:"children,omitempty"`
	BaseW    int    `json:"baseW,omitempty"`
	BaseH    int    `json:"baseH,omitempty"`
}

// GroupComponent is the component id of a group.
const GroupComponent = "cuppa.group"

// IsGroup reports whether the node is a group of other nodes.
func (n Node) IsGroup() bool { return n.Component == GroupComponent && len(n.Children) > 0 }

// Clone returns a deep copy of the node.
func (n Node) Clone() Node {
	if n.Props != nil {
		props := make(map[string]string, len(n.Props))
		for k, v := range n.Props {
			props[k] = v
		}
		n.Props = props
	}
	if n.Children != nil {
		children := make([]Node, len(n.Children))
		for i, c := range n.Children {
			children[i] = c.Clone()
		}
		n.Children = children
	}
	return n
}
