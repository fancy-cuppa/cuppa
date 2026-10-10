package gosource

import "strings"

// splitList cuts s at sep and trims the pieces, leaving out empty ones. A
// backslash before a comma, a semicolon or another backslash is that
// character itself, so an item or a table cell can hold one.
func splitList(s, sep string) []string { return splitEscaped(s, sep, true) }

// splitRaw is splitList that keeps the backslashes, for the first cut of a
// table: the rows are cut at semicolons, then each row at commas.
func splitRaw(s, sep string) []string { return splitEscaped(s, sep, false) }

func splitEscaped(s, sep string, unescape bool) []string {
	var out []string
	var cur strings.Builder
	flush := func() {
		if piece := strings.TrimSpace(cur.String()); piece != "" {
			out = append(out, piece)
		}
		cur.Reset()
	}
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
			flush()
			i += len(sep) - 1
			continue
		}
		cur.WriteByte(c)
	}
	flush()
	return out
}
