package theme

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

// Overlay draws repl over base starting at column x, keeping the styled text
// on either side. base must already be padded to the full width.
func Overlay(base, repl string, x int) string {
	w := ansi.StringWidth(repl)
	left := ansi.Truncate(base, x, "")
	if pad := x - ansi.StringWidth(left); pad > 0 {
		left += strings.Repeat(" ", pad)
	}
	return left + "\x1b[m" + repl + "\x1b[m" + ansi.TruncateLeft(base, x+w, "")
}

// Panel frames content lines in a rounded box with a title, exactly w cells
// wide. The result has len(content)+2 lines.
func Panel(title string, content []string, w int) []string {
	if w < 6 {
		w = 6
	}
	edge := lipgloss.NewStyle().Foreground(lipgloss.Color(Accent))
	head := "─ " + title + " "
	fill := max(w-2-ansi.StringWidth(head), 0)
	out := make([]string, 0, len(content)+2)
	out = append(out, edge.Render("╭"+head+strings.Repeat("─", fill)+"╮"))
	for _, l := range content {
		out = append(out, edge.Render("│")+Fit(l, w-2)+edge.Render("│"))
	}
	out = append(out, edge.Render("╰"+strings.Repeat("─", w-2)+"╯"))
	return out
}
