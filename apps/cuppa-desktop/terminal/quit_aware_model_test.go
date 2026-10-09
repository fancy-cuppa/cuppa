package terminal

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

type quitter struct{}

func (quitter) Init() tea.Cmd { return nil }
func (q quitter) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyPressMsg); ok && key.String() == "q" {
		return q, tea.Quit
	}
	return q, nil
}
func (quitter) View() tea.View { return tea.NewView("") }

func TestOnQuitFiresWhenTheModelQuits(t *testing.T) {
	quits := 0
	model := OnQuit(quitter{}, func() { quits++ })

	_, cmd := model.Update(tea.KeyPressMsg{Code: 'x', Text: "x"})
	if cmd != nil {
		t.Fatal("no command expected for a harmless key")
	}
	_, cmd = model.Update(tea.KeyPressMsg{Code: 'q', Text: "q"})
	if cmd == nil {
		t.Fatal("expected the quit command to be passed through")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatal("the quit command must still produce a QuitMsg")
	}
	if quits != 1 {
		t.Fatalf("onQuit called %d times, want 1", quits)
	}
}
