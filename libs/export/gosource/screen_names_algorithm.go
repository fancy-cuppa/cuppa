package gosource

import (
	"strings"
	"unicode"
)

// words splits a name into its letters-and-digits words.
func words(s string) []string {
	return strings.FieldsFunc(s, func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) })
}

// pascal turns a name into an exported Go identifier: "Show status" becomes
// "ShowStatus". A name with no letters gives "".
func pascal(s string) string {
	var b strings.Builder
	for _, w := range words(s) {
		r := []rune(w)
		b.WriteRune(unicode.ToUpper(r[0]))
		b.WriteString(string(r[1:]))
	}
	out := b.String()
	if out != "" && unicode.IsDigit([]rune(out)[0]) {
		out = "X" + out
	}
	return out
}

// camel is pascal with a lower-case first word: "mvdColours".
func camel(s string) string {
	ws := words(s)
	if len(ws) == 0 {
		return ""
	}
	first := strings.ToLower(ws[0])
	if unicode.IsDigit([]rune(first)[0]) {
		first = "x" + first
	}
	rest := pascal(strings.Join(ws[1:], " "))
	return first + rest
}

// snake is the file name stem of a name: "mvd_colours".
func snake(s string) string {
	ws := words(s)
	for i := range ws {
		ws[i] = strings.ToLower(ws[i])
	}
	return strings.Join(ws, "_")
}

// ScreenName is the Go name the export gives a design: the function that draws
// it.
func ScreenName(design string) string { return pascal(design) }
