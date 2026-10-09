package terminal

import tea "charm.land/bubbletea/v2"

// OnQuit wraps a model and calls fn when the model asks the program to quit,
// so the window can close when the user chooses Quit in the app.
func OnQuit(inner tea.Model, fn func()) tea.Model {
	return &quitAware{inner: inner, onQuit: fn}
}

type quitAware struct {
	inner  tea.Model
	onQuit func()
}

func (m *quitAware) Init() tea.Cmd { return m.watch(m.inner.Init()) }

func (m *quitAware) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	next, cmd := m.inner.Update(msg)
	m.inner = next
	return m, m.watch(cmd)
}

func (m *quitAware) View() tea.View { return m.inner.View() }

// watch runs cmd as usual and reports when it produced a quit.
func (m *quitAware) watch(cmd tea.Cmd) tea.Cmd {
	if cmd == nil {
		return nil
	}
	return func() tea.Msg {
		msg := cmd()
		if _, quit := msg.(tea.QuitMsg); quit {
			m.onQuit()
		}
		return msg
	}
}
