package shell

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/meta-tui/cuppa/apps/cuppa-tui/modal"
	"github.com/meta-tui/cuppa/apps/cuppa-tui/packsdialog"
	"github.com/meta-tui/cuppa/libs/catalog/cupp"
	"github.com/meta-tui/cuppa/libs/document/design"
)

func teaFile(t *testing.T, dir string) string {
	t.Helper()
	data, err := cupp.Encode(cupp.Pack{ID: "tea-shop", Name: "Tea shop", Components: []design.Composite{{
		ID: "card", Name: "Card", W: 20, H: 5,
		Nodes: []design.Node{{ID: "a", Component: "lipgloss.box", Name: "Frame", Rect: design.Rect{W: 20, H: 5}}},
	}}})
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "download.cupp")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func packsShell(t *testing.T) (*Model, string) {
	t.Helper()
	m := newShell(t)
	m.userPacks = filepath.Join(t.TempDir(), "packs")
	return m, m.userPacks
}

func TestInstallingAPackCopiesItAndAddsItsComponents(t *testing.T) {
	m, dir := packsShell(t)
	if err := m.installPack(teaFile(t, t.TempDir())); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "tea-shop.cupp")); err != nil {
		t.Fatalf("the pack is copied under its own id: %v", err)
	}
	if _, ok := m.cat.Get("tea-shop.card"); !ok {
		t.Fatal("the component is in the catalog now")
	}
	m.pal.SetSize(m.layout.palette.W, 80)
	if !strings.Contains(paletteText(m), "Tea shop") {
		t.Fatalf("the palette lists the pack:\n%s", paletteText(m))
	}
	if err := m.installPack(teaFile(t, t.TempDir())); err == nil {
		t.Fatal("the same pack twice is refused")
	}
}

func TestAddingSomethingThatIsNotAPackIsRefused(t *testing.T) {
	m, dir := packsShell(t)
	bad := filepath.Join(t.TempDir(), "bad.cupp")
	if err := os.WriteFile(bad, []byte("nope"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := m.installPack(bad); err == nil {
		t.Fatal("a damaged file is refused")
	}
	if _, err := os.Stat(dir); err == nil {
		t.Fatal("nothing is created for a refused pack")
	}
}

func TestRemovingAPackDeletesItsFileAfterAsking(t *testing.T) {
	m, dir := packsShell(t)
	if err := m.installPack(teaFile(t, t.TempDir())); err != nil {
		t.Fatal(err)
	}
	m.afterPacksDialog(modal.Outcome{Button: packsdialog.RemoveButton, Value: "tea-shop"})
	dlg := m.flow.Modal()
	if dlg == nil {
		t.Fatal("removing asks first")
	}
	if _, err := os.Stat(filepath.Join(dir, "tea-shop.cupp")); err != nil {
		t.Fatal("nothing is deleted before the answer")
	}
	dlg.Key("", false, true, false) // Enter: the first button, Remove
	m.flow.Resolve()
	if _, err := os.Stat(filepath.Join(dir, "tea-shop.cupp")); err == nil {
		t.Fatal("the file is deleted")
	}
	if _, ok := m.cat.Get("tea-shop.card"); ok {
		t.Fatal("the pack is gone from the catalog")
	}
}
