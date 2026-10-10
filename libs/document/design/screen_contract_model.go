package design

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

// KeyBinding is a key that raises a screen event. When the design is exported
// as a screen, the key is read by the screen's handle function and the event
// comes back to the program that owns the screen.
type KeyBinding struct {
	// Key is a key as Bubble Tea writes it: "s", "esc", "ctrl+s", "enter".
	Key string `json:"key"`
	// Event is the name of the event the key raises.
	Event string `json:"event"`
	// Label is what the key bar shows for the key; it may be empty.
	Label string `json:"label,omitempty"`
}

// ValidInputName reports whether s can name a screen input or event: it starts
// with a letter and holds letters, digits, spaces, dashes and underscores. The
// export turns it into a Go identifier.
func ValidInputName(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	for i, r := range s {
		switch {
		case unicode.IsLetter(r):
		case i > 0 && (unicode.IsDigit(r) || r == ' ' || r == '-' || r == '_'):
		default:
			return false
		}
	}
	return true
}

// CleanInputName trims a name and collapses its inner spaces.
func CleanInputName(s string) string { return strings.Join(strings.Fields(s), " ") }

// ParseKeys reads a list of key bindings written as "key=Event" or
// "key=Event:Label", separated by commas: "s=Save:save, esc=Back".
func ParseKeys(spec string) ([]KeyBinding, error) {
	var out []KeyBinding
	seen := map[string]bool{}
	for _, part := range strings.Split(spec, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		key, rest, ok := strings.Cut(part, "=")
		key, rest = strings.TrimSpace(key), strings.TrimSpace(rest)
		event, label, _ := strings.Cut(rest, ":")
		event, label = CleanInputName(event), strings.TrimSpace(label)
		if !ok || key == "" || !ValidInputName(event) {
			return nil, fmt.Errorf("design: %q is not key=Event or key=Event:label", part)
		}
		if seen[key] {
			return nil, fmt.Errorf("design: key %q is used twice", key)
		}
		seen[key] = true
		out = append(out, KeyBinding{Key: key, Event: event, Label: label})
	}
	return out, nil
}

// FormatKeys is ParseKeys the other way round.
func FormatKeys(keys []KeyBinding) string {
	parts := make([]string, len(keys))
	for i, k := range keys {
		parts[i] = k.Key + "=" + k.Event
		if k.Label != "" {
			parts[i] += ":" + k.Label
		}
	}
	return strings.Join(parts, ", ")
}

// BoundKeys is the property keys of the node that are bound to an input, in
// order.
func (n Node) BoundKeys() []string {
	keys := make([]string, 0, len(n.Bind))
	for k := range n.Bind {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
