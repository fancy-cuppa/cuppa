package scene

import (
	"strconv"
	"strings"

	"github.com/meta-tui/cuppa/libs/catalog/definition"
	"github.com/meta-tui/cuppa/libs/render/grid"
)

// Props are a node's property values merged over the component defaults.
type Props map[string]string

// Str returns the value of key.
func (p Props) Str(key string) string { return p[key] }

// Int returns the value of key as an integer, or def if it does not parse.
func (p Props) Int(key string, def int) int {
	n, err := strconv.Atoi(strings.TrimSpace(p[key]))
	if err != nil {
		return def
	}
	return n
}

// Bool returns the value of key as a boolean.
func (p Props) Bool(key string) bool { return p[key] == "true" }

// List splits a comma separated value into trimmed, non-empty items.
func (p Props) List(key string) []string { return splitList(p[key], ",") }

func splitList(s, sep string) []string {
	var out []string
	for _, part := range strings.Split(s, sep) {
		if part = strings.TrimSpace(part); part != "" {
			out = append(out, part)
		}
	}
	return out
}

// Dim is the style of quiet text: the design's muted colour, or the default grey
// when the theme sets none.
func (p Props) Dim() grid.Style {
	if c := p[definition.MutedKey]; c != "" {
		return grid.Style{Fg: c}
	}
	return dim
}
