package shell

import (
	"errors"
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

func TestPackProblemsAreListedInANotice(t *testing.T) {
	m := newShell(t)
	m.ReportPackProblems()
	if m.flow.Modal() != nil {
		t.Fatal("no problems, no notice")
	}
	m.packProblems = []error{errors.New("tea.cupp: file is damaged")}
	m.ReportPackProblems()
	if m.flow.Modal() == nil {
		t.Fatal("a notice should open")
	}
	var screen []string
	for _, l := range m.flow.Modal().Lines() {
		screen = append(screen, ansi.Strip(l))
	}
	if !strings.Contains(strings.Join(screen, "\n"), "tea.cupp: file is damaged") {
		t.Fatalf("the notice names the file and the problem:\n%s", strings.Join(screen, "\n"))
	}
}
