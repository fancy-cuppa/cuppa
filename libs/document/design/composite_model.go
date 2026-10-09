package design

import (
	"fmt"
	"regexp"
)

// Composite is a component the user made out of other components: a small
// design with a default size, plus the properties it offers as its own.
type Composite struct {
	// ID is unique inside its pack: lowercase letters, digits and dashes.
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	// W and H are the size it has when first placed. Inner rectangles are
	// relative to a box of this size and scale with the placed node.
	W int `json:"w"`
	H int `json:"h"`
	// Nodes are the inner components, back to front.
	Nodes []Node `json:"nodes"`
	// Props are the properties shown in the details bar; each one drives a
	// property of an inner node.
	Props []Exposed `json:"props,omitempty"`
}

// Exposed is one property a composite offers. Kind is one of text, int, bool,
// color or choice, like a catalog property.
type Exposed struct {
	Key     string   `json:"key"`
	Label   string   `json:"label"`
	Kind    string   `json:"kind"`
	Default string   `json:"default,omitempty"`
	Choices []string `json:"choices,omitempty"`
	// Target is the inner node and TargetProp the property of it that this one sets.
	Target     NodeID `json:"target"`
	TargetProp string `json:"targetProp"`
}

var (
	idPattern   = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
	propKinds   = map[string]bool{"text": true, "int": true, "bool": true, "color": true, "choice": true}
	maxInnerLen = 500
)

// ValidID reports whether s can be a pack or composite id.
func ValidID(s string) bool { return len(s) <= 64 && idPattern.MatchString(s) }

// Check says why the composite cannot be used, or returns nil.
func (c Composite) Check() error {
	switch {
	case !ValidID(c.ID):
		return fmt.Errorf("component id %q must be lowercase letters, digits and dashes", c.ID)
	case c.Name == "":
		return fmt.Errorf("component %q has no name", c.ID)
	case c.W < 1 || c.H < 1:
		return fmt.Errorf("component %q has no size", c.ID)
	case len(c.Nodes) > maxInnerLen:
		return fmt.Errorf("component %q has too many parts", c.ID)
	}
	ids := map[NodeID]bool{}
	for _, n := range c.Nodes {
		if n.ID == "" || ids[n.ID] || n.Rect.W < 1 || n.Rect.H < 1 {
			return fmt.Errorf("component %q has a damaged part %q", c.ID, n.ID)
		}
		ids[n.ID] = true
	}
	keys := map[string]bool{}
	for _, p := range c.Props {
		switch {
		case p.Key == "" || keys[p.Key]:
			return fmt.Errorf("component %q has a repeated or empty property key", c.ID)
		case !propKinds[p.Kind]:
			return fmt.Errorf("component %q property %q has unknown kind %q", c.ID, p.Key, p.Kind)
		case p.Kind == "choice" && len(p.Choices) == 0:
			return fmt.Errorf("component %q property %q has no choices", c.ID, p.Key)
		case !ids[p.Target] || p.TargetProp == "":
			return fmt.Errorf("component %q property %q drives nothing", c.ID, p.Key)
		}
		keys[p.Key] = true
	}
	return nil
}

// Clone returns a deep copy.
func (c Composite) Clone() Composite {
	nodes := make([]Node, len(c.Nodes))
	for i, n := range c.Nodes {
		nodes[i] = n.Clone()
	}
	c.Nodes = nodes
	c.Props = append([]Exposed(nil), c.Props...)
	return c
}
