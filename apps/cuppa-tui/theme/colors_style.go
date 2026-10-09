// Package theme holds the colours and text styles every pane shares.
package theme

import "charm.land/lipgloss/v2"

// Colours as Lip Gloss color specs (ANSI 256).
const (
	Accent    = "212"
	Muted     = "245"
	Faint     = "240"
	Canvas    = "235"
	Highlight = "238"
	Good      = "114"
	Warn      = "215"
)

// Title styles a pane heading.
func Title(s string) string {
	return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(Accent)).Render(s)
}

// Dim styles secondary text.
func Dim(s string) string { return lipgloss.NewStyle().Foreground(lipgloss.Color(Muted)).Render(s) }

// Faded styles hints and separators.
func Faded(s string) string { return lipgloss.NewStyle().Foreground(lipgloss.Color(Faint)).Render(s) }

// Bold styles emphasised text.
func Bold(s string) string { return lipgloss.NewStyle().Bold(true).Render(s) }

// Button styles a clickable label; active buttons are highlighted.
func Button(s string, enabled bool) string {
	st := lipgloss.NewStyle().Foreground(lipgloss.Color(Accent))
	if !enabled {
		st = lipgloss.NewStyle().Foreground(lipgloss.Color(Faint))
	}
	return st.Render(s)
}

// Selected styles the currently selected row.
func Selected(s string) string {
	return lipgloss.NewStyle().Bold(true).Background(lipgloss.Color(Highlight)).Render(s)
}

// Hovered styles the row under the pointer.
func Hovered(s string) string {
	return lipgloss.NewStyle().Background(lipgloss.Color(Highlight)).Render(s)
}

// Fit pads or truncates a single line to exactly w cells.
func Fit(s string, w int) string {
	if w <= 0 {
		return ""
	}
	return lipgloss.NewStyle().Width(w).MaxWidth(w).Render(s)
}

// Block returns exactly h lines of exactly w cells, padding with blank lines.
func Block(lines []string, w, h int) []string {
	out := make([]string, h)
	for i := range out {
		if i < len(lines) {
			out[i] = Fit(lines[i], w)
		} else {
			out[i] = Fit("", w)
		}
	}
	return out
}
