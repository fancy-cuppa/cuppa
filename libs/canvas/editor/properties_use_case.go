package editor

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/meta-tui/cuppa/libs/catalog/definition"
	"github.com/meta-tui/cuppa/libs/document/design"
)

// Rename sets the node's display name.
func (e *Editor) Rename(id design.NodeID, name string) bool {
	name = strings.TrimSpace(name)
	n, ok := e.doc.Get(id)
	if !ok || name == "" || name == n.Name {
		return false
	}
	return e.apply(func() bool {
		return e.doc.Update(id, func(n *design.Node) { n.Name = name })
	})
}

// SetProp sets one property after validating the value against the
// component's property specification.
func (e *Editor) SetProp(id design.NodeID, key, value string) error {
	n, ok := e.doc.Get(id)
	if !ok {
		return fmt.Errorf("editor: no node %q", id)
	}
	def, ok := e.cat.Get(n.Component)
	if !ok {
		return fmt.Errorf("editor: unknown component %q", n.Component)
	}
	spec, ok := def.Prop(key)
	if !ok {
		return fmt.Errorf("editor: %s has no property %q", def.ID, key)
	}
	value, err := normalise(spec, value)
	if err != nil {
		return err
	}
	if cur, set := n.Props[key]; (set && cur == value) || (!set && spec.Default == value) {
		return nil
	}
	e.apply(func() bool {
		return e.doc.Update(id, func(n *design.Node) {
			if n.Props == nil {
				n.Props = map[string]string{}
			}
			n.Props[key] = value
		})
	})
	return nil
}

// normalise validates value for spec and returns its canonical form.
func normalise(spec definition.PropSpec, value string) (string, error) {
	switch spec.Kind {
	case definition.PropInt:
		n, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil {
			return "", fmt.Errorf("editor: %s must be a whole number", spec.Label)
		}
		if spec.Max > spec.Min {
			n = min(max(n, spec.Min), spec.Max)
		}
		return strconv.Itoa(n), nil
	case definition.PropBool:
		switch strings.ToLower(strings.TrimSpace(value)) {
		case "true", "yes", "on", "1":
			return "true", nil
		case "false", "no", "off", "0":
			return "false", nil
		}
		return "", fmt.Errorf("editor: %s must be true or false", spec.Label)
	case definition.PropChoice:
		if !slices.Contains(spec.Choices, value) {
			return "", fmt.Errorf("editor: %s must be one of %s", spec.Label, strings.Join(spec.Choices, ", "))
		}
	}
	return value, nil
}
