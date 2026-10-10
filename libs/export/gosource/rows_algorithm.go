package gosource

import (
	"strconv"
	"strings"
)

// rowsKind is the component whose bound rows are a typed list in a screen's
// contract (ADR 0008).
const rowsKind = "lipgloss.rows"

// colourPickerKind is the component whose bound value is a ColourPicker the
// program owns.
const colourPickerKind = "lipgloss.colourpicker"

// rowColumn is one column of a Rows component.
type rowColumn struct {
	name   string
	colour bool
}

// rowColumns reads the component's columns property: "Name[:width][:colour]"
// items separated by commas.
func rowColumns(spec string) []rowColumn {
	var out []rowColumn
	for _, item := range rowSplit(spec, ",", true) {
		parts := strings.Split(strings.TrimSpace(item), ":")
		col := rowColumn{name: strings.TrimSpace(parts[0])}
		for _, extra := range parts[1:] {
			extra = strings.ToLower(strings.TrimSpace(extra))
			if _, err := strconv.Atoi(extra); err != nil && (extra == "colour" || extra == "color") {
				col.colour = true
			}
		}
		out = append(out, col)
	}
	return out
}

// rowSplit cuts s at sep keeping empty pieces and whitespace, with the same
// escapes as the lists (a backslash before a comma, a semicolon or a backslash).
// With unescape false the backslashes stay, for the first cut of the rows.
func rowSplit(s, sep string, unescape bool) []string {
	var out []string
	var cur strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '\\' && i+1 < len(s) && (s[i+1] == ',' || s[i+1] == ';' || s[i+1] == '\\') {
			if !unescape {
				cur.WriteByte(c)
			}
			cur.WriteByte(s[i+1])
			i++
			continue
		}
		if strings.HasPrefix(s[i:], sep) {
			out = append(out, cur.String())
			cur.Reset()
			i += len(sep) - 1
			continue
		}
		cur.WriteByte(c)
	}
	return append(out, cur.String())
}

// rowCells reads the rows property: cells separated by commas, rows by
// semicolons.
func rowCells(s string) [][]string {
	var rows [][]string
	for _, raw := range rowSplit(s, ";", false) {
		if raw != "" {
			rows = append(rows, rowSplit(raw, ",", true))
		}
	}
	return rows
}

// rowStyleConst is the Go constant of a row style name.
func rowStyleConst(style string) string {
	switch style {
	case "selected":
		return "RowSelected"
	case "dim":
		return "RowDim"
	case "accent":
		return "RowAccent"
	}
	return "RowNormal"
}
