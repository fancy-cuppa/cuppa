package shell

import (
	"os"
	"strings"
	"path/filepath"

	"github.com/meta-tui/cuppa/apps/cuppa-tui/modal"
	"github.com/meta-tui/cuppa/apps/cuppa-tui/packsdialog"
	"github.com/meta-tui/cuppa/libs/catalog/definition"
	"github.com/meta-tui/cuppa/libs/catalog/packstate"
)

// restorePacks brings back which packs the user switched off and starts
// remembering changes. RestoreLayout calls it, so tests never touch settings.
func (m *Model) restorePacks() {
	dir, err := os.UserConfigDir()
	if err != nil {
		return
	}
	m.packs = packstate.Open(filepath.Join(dir, "cuppa", "packs.json"))
	m.refreshPalette()
}

// refreshPalette lists only the components of enabled packs. The full catalog
// stays behind the editor, so a design that uses a disabled pack still renders.
func (m *Model) refreshPalette() {
	m.pal.SetCatalog(m.cat.Only(m.packs.Enabled))
}

// ReportPackProblems tells the user which installed packs could not be used.
// Front ends call it once at start, after loading them.
func (m *Model) ReportPackProblems(problems []error) {
	if len(problems) == 0 {
		return
	}
	var lines []string
	for _, p := range problems {
		lines = append(lines, "• "+p.Error())
	}
	m.flow.Notice("Some packs could not be loaded", strings.Join(lines, "\n"))
}

// openPacks shows the Packs dialog; each click applies at once.
func (m *Model) openPacks() {
	var entries []packsdialog.Entry
	for _, p := range m.cat.Packs() {
		entries = append(entries, packsdialog.Entry{Pack: p, Count: len(m.cat.ByFamily(p.ID))})
	}
	toggle := func(id definition.Family) {
		m.packs.SetEnabled(id, !m.packs.Enabled(id))
		m.refreshPalette()
	}
	m.flow.Show(packsdialog.New(entries, m.packs.Enabled, toggle), func(modal.Outcome) {})
}
