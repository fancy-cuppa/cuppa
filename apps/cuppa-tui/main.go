package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/fancy-cuppa/cuppa/apps/cuppa-tui/shell"
	"github.com/fancy-cuppa/cuppa/libs/catalog/standard"
)

// version is stamped at build time via -ldflags "-X main.version=...".
var version = "dev"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Println("cuppa", version)
		return
	}
	if _, err := tea.NewProgram(shell.New(standard.Default())).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "cuppa:", err)
		os.Exit(1)
	}
}
