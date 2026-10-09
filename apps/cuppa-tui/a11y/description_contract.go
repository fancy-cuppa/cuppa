// Package a11y is how the panes of the Cuppa terminal app describe themselves
// to assistive technology. It mirrors the shape of a screen snapshot (roles are
// ARIA roles) without importing any front end: the desktop app maps it onto
// TReactUI, and a pane that has nothing to say simply returns no nodes.
package a11y

// Node is one element of a description. Role is an ARIA role: heading, list,
// listitem, listbox, option, button, textbox, status or text.
type Node struct {
	Role     string
	Label    string
	Value    string
	Selected bool
	Focused  bool
	Children []Node
}

// Snapshot describes the whole screen.
type Snapshot struct {
	Title string
	// Status is the line of feedback at the bottom, for announcing changes.
	Status string
	Nodes  []Node
}

// Describer is implemented by a pane or dialog that can describe itself.
type Describer interface {
	Describe() []Node
}

// Heading, Text and the other constructors keep the panes' code short.
func Heading(label string) Node { return Node{Role: "heading", Label: label} }

// Text is a line of plain text.
func Text(label string) Node { return Node{Role: "text", Label: label} }

// Button is something that can be pressed.
func Button(label string) Node { return Node{Role: "button", Label: label} }

// Field is a labelled value that can be edited.
func Field(label, value string) Node { return Node{Role: "textbox", Label: label, Value: value} }

// List groups options or items under a label.
func List(label string, children ...Node) Node {
	return Node{Role: "list", Label: label, Children: children}
}

// Item is one entry of a list.
func Item(label string, selected bool) Node {
	return Node{Role: "listitem", Label: label, Selected: selected}
}
