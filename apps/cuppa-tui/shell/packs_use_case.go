package shell

import (
	"os"
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
