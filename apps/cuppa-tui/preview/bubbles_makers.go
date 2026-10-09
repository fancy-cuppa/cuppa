package preview

import (
	"strings"
	"time"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/paginator"
	"charm.land/bubbles/v2/progress"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/stopwatch"
	"charm.land/bubbles/v2/table"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/timer"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// maker starts the real Bubbles model for a component, sized w by h.
type maker func(p props, w, h int) *live

// makers are the components that run as real models in the preview, by catalog id.
var makers = map[string]maker{
	"bubbles.textinput": makeTextInput,
	"bubbles.textarea":  makeTextArea,
	"bubbles.list":      makeList,
	"bubbles.table":     makeTable,
	"bubbles.viewport":  makeViewport,
	"bubbles.paginator": makePaginator,
	"bubbles.spinner":   makeSpinner,
	"bubbles.progress":  makeProgress,
	"bubbles.stopwatch": makeStopwatch,
	"bubbles.timer":     makeTimer,
}

// Live reports whether a component runs as a real model in the preview.
func Live(componentID string) bool {
	_, ok := makers[componentID]
	return ok
}

func fg(c string) lipgloss.Style {
	if c == "" {
		return lipgloss.NewStyle()
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color(c))
}

func makeTextInput(p props, w, _ int) *live {
	m := textinput.New()
	m.Prompt = p.str("prompt")
	m.Placeholder = p.str("placeholder")
	m.SetValue(p.str("value"))
	m.SetWidth(max(w-lipgloss.Width(m.Prompt)-1, 1))
	st := m.Styles()
	st.Focused.Prompt = fg(p.str("color"))
	st.Blurred.Prompt = fg(p.str("color"))
	m.SetStyles(st)
	return &live{
		keys:   true,
		update: func(msg tea.Msg) tea.Cmd { var c tea.Cmd; m, c = m.Update(msg); return c },
		view:   func() string { return m.View() },
		focus: func(on bool) tea.Cmd {
			if on {
				return m.Focus()
			}
			m.Blur()
			return nil
		},
	}
}

func makeTextArea(p props, w, h int) *live {
	m := textarea.New()
	m.Placeholder = p.str("placeholder")
	m.ShowLineNumbers = p.flag("line_numbers")
	m.SetWidth(w)
	m.SetHeight(h)
	m.SetValue(strings.ReplaceAll(p.str("value"), "|", "\n"))
	return &live{
		keys:   true,
		update: func(msg tea.Msg) tea.Cmd { var c tea.Cmd; m, c = m.Update(msg); return c },
		view:   func() string { return m.View() },
		focus: func(on bool) tea.Cmd {
			if on {
				return m.Focus()
			}
			m.Blur()
			return nil
		},
	}
}

// word is a list item.
type word string

func (w word) FilterValue() string { return string(w) }

// Title and Description make a word a list.DefaultItem, which is what the default
// delegate draws.
func (w word) Title() string       { return string(w) }
func (w word) Description() string { return "" }

func split(s, sep string) []string {
	var out []string
	for _, part := range strings.Split(s, sep) {
		if part = strings.TrimSpace(part); part != "" {
			out = append(out, part)
		}
	}
	return out
}

func makeList(p props, w, h int) *live {
	var items []list.Item
	for _, s := range split(p.str("items"), ",") {
		items = append(items, word(s))
	}
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.SetSpacing(0)
	m := list.New(items, delegate, w, h)
	m.Title = p.str("title")
	m.SetShowFilter(p.flag("show_filter"))
	m.SetFilteringEnabled(p.flag("show_filter"))
	m.SetShowStatusBar(p.flag("show_status"))
	m.Select(p.integer("selected", 0))
	return &live{
		keys:   true,
		update: func(msg tea.Msg) tea.Cmd { var c tea.Cmd; m, c = m.Update(msg); return c },
		view:   func() string { return m.View() },
	}
}

func makeTable(p props, w, h int) *live {
	var cols []table.Column
	names := split(p.str("columns"), ",")
	each := max((w-2*len(names))/max(len(names), 1), 3)
	for _, n := range names {
		cols = append(cols, table.Column{Title: n, Width: each})
	}
	var rows []table.Row
	for _, line := range split(p.str("rows"), ";") {
		rows = append(rows, table.Row(split(line, ",")))
	}
	st := table.DefaultStyles()
	st.Selected = st.Selected.Foreground(lipgloss.Color("0")).Background(lipgloss.Color(orDefault(p.str("color"), "212")))
	m := table.New(table.WithColumns(cols), table.WithRows(rows), table.WithHeight(h), table.WithWidth(w), table.WithStyles(st))
	m.SetCursor(p.integer("selected", 0))
	return &live{
		keys:   true,
		update: func(msg tea.Msg) tea.Cmd { var c tea.Cmd; m, c = m.Update(msg); return c },
		view:   func() string { return m.View() },
		focus:  func(on bool) tea.Cmd { setTableFocus(&m, on); return nil },
	}
}

func setTableFocus(m *table.Model, on bool) {
	if on {
		m.Focus()
	} else {
		m.Blur()
	}
}

func orDefault(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}

func makeViewport(p props, w, h int) *live {
	m := viewport.New(viewport.WithWidth(w), viewport.WithHeight(h))
	m.SetContent(strings.ReplaceAll(p.str("content"), "|", "\n"))
	return &live{
		keys:   true,
		update: func(msg tea.Msg) tea.Cmd { var c tea.Cmd; m, c = m.Update(msg); return c },
		view:   func() string { return m.View() },
	}
}

func makePaginator(p props, _, _ int) *live {
	m := paginator.New(paginator.WithTotalPages(max(p.integer("total", 5), 1)))
	if p.str("style") == "arabic" {
		m.Type = paginator.Arabic
	} else {
		m.Type = paginator.Dots
	}
	m.Page = max(p.integer("page", 1)-1, 0)
	return &live{
		keys:   true,
		update: func(msg tea.Msg) tea.Cmd { var c tea.Cmd; m, c = m.Update(msg); return c },
		view:   func() string { return fg(p.str("color")).Render(m.View()) },
	}
}

var spinners = map[string]spinner.Spinner{
	"line": spinner.Line, "dot": spinner.Dot, "minidot": spinner.MiniDot, "jump": spinner.Jump,
	"pulse": spinner.Pulse, "points": spinner.Points, "globe": spinner.Globe, "moon": spinner.Moon,
	"monkey": spinner.Monkey,
}

func makeSpinner(p props, _, _ int) *live {
	kind, ok := spinners[p.str("style")]
	if !ok {
		kind = spinner.Dot
	}
	m := spinner.New(spinner.WithSpinner(kind), spinner.WithStyle(fg(p.str("color"))))
	return &live{
		start:  m.Tick,
		update: func(msg tea.Msg) tea.Cmd { var c tea.Cmd; m, c = m.Update(msg); return c },
		view:   func() string { return m.View() + " " + p.str("label") },
	}
}

func makeProgress(p props, w, _ int) *live {
	opts := []progress.Option{progress.WithWidth(max(w, 4)), progress.WithColors(lipgloss.Color(orDefault(p.str("color"), "212")))}
	if !p.flag("show_percentage") {
		opts = append(opts, progress.WithoutPercentage())
	}
	m := progress.New(opts...)
	percent := float64(p.integer("percent", 0)) / 100
	return &live{
		keys:  true,
		start: m.SetPercent(percent),
		update: func(msg tea.Msg) tea.Cmd {
			if k, ok := msg.(tea.KeyPressMsg); ok {
				switch k.String() {
				case "right", "up", "+":
					percent = min(percent+0.1, 1)
				case "left", "down", "-":
					percent = max(percent-0.1, 0)
				default:
					return nil
				}
				return m.SetPercent(percent)
			}
			var c tea.Cmd
			m, c = m.Update(msg)
			return c
		},
		view: func() string { return m.View() },
	}
}

func makeStopwatch(p props, _, _ int) *live {
	m := stopwatch.New(stopwatch.WithInterval(100 * time.Millisecond))
	l := &live{keys: true, view: func() string { return fg(p.str("color")).Render(m.View()) }}
	if p.flag("running") {
		l.start = m.Start()
	}
	l.update = func(msg tea.Msg) tea.Cmd {
		if k, ok := msg.(tea.KeyPressMsg); ok {
			if k.String() == "space" || k.String() == "enter" || k.String() == " " {
				return m.Toggle()
			}
			return nil
		}
		var c tea.Cmd
		m, c = m.Update(msg)
		return c
	}
	return l
}

func makeTimer(p props, _, _ int) *live {
	m := timer.New(parseClock(p.str("value")), timer.WithInterval(100*time.Millisecond))
	l := &live{keys: true, view: func() string { return fg(p.str("color")).Render(p.str("label") + " " + m.View()) }}
	if p.flag("running") {
		l.start = m.Start()
	}
	l.update = func(msg tea.Msg) tea.Cmd {
		if k, ok := msg.(tea.KeyPressMsg); ok {
			if k.String() == "space" || k.String() == "enter" || k.String() == " " {
				return m.Toggle()
			}
			return nil
		}
		var c tea.Cmd
		m, c = m.Update(msg)
		return c
	}
	return l
}

// parseClock reads "04:32" (minutes:seconds) or "1:02:03"; anything else is a minute.
func parseClock(s string) time.Duration {
	total := 0
	for _, part := range strings.Split(strings.TrimSpace(s), ":") {
		n := 0
		for _, r := range part {
			if r < '0' || r > '9' {
				return time.Minute
			}
			n = n*10 + int(r-'0')
		}
		total = total*60 + n
	}
	if total == 0 {
		return time.Minute
	}
	return time.Duration(total) * time.Second
}
