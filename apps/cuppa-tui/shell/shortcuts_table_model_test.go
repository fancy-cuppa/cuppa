package shell

import (
	"os"
	"strings"
	"testing"
)

// The user guide lists every key of the table, so the dialog and the guide
// cannot drift apart.
func TestTheUserGuideMentionsEveryKeyOfTheShortcutsTable(t *testing.T) {
	raw, err := os.ReadFile("../../../docs/user-guide.md")
	if err != nil {
		t.Skip("user guide not found: " + err.Error())
	}
	guide := strings.ReplaceAll(string(raw), "\r\n", "\n")
	for _, s := range shortcuts {
		for _, b := range s.bindings {
			if !strings.Contains(guide, b.keys) {
				t.Errorf("%s: the user guide does not mention %q (%s)", s.area, b.keys, b.what)
			}
		}
	}
}

func TestEveryAreaIsInTheDialogAndCtrlBecomesCmdOnAMac(t *testing.T) {
	text := shortcutsText(false)
	for _, s := range shortcuts {
		if !strings.Contains(text, s.area) {
			t.Errorf("dialog lacks the %q section", s.area)
		}
	}
	if mac := shortcutsText(true); strings.Contains(mac, "Ctrl") || !strings.Contains(mac, "Cmd+Shift+S") {
		t.Fatalf("on a Mac the dialog says: %s", mac)
	}
}
