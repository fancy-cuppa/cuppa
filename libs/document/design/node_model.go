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
}

// Clone returns a deep copy of the node.
func (n Node) Clone() Node {
	if n.Props != nil {
		props := make(map[string]string, len(n.Props))
		for k, v := range n.Props {
			props[k] = v
		}
		n.Props = props
	}
	return n
}
