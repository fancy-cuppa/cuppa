package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/meta-tui/cuppa/apps/cuppa-tui/shell"
	"github.com/meta-tui/cuppa/libs/catalog/standard"
)

// version is stamped at build time via -ldflags "-X main.version=...".
var version = "dev"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Println("cuppa", version)
		return
	}
	app := shell.New(standard.Default())
	app.RestoreLayout()
	app.LoadUserPacks()
	if len(os.Args) > 1 {
		if err := app.OpenFile(os.Args[1]); err != nil {
			fmt.Fprintln(os.Stderr, "cuppa:", err)
			os.Exit(1)
		}
	}
	app.Welcome()
	app.ReportPackProblems()
	if _, err := tea.NewProgram(app).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "cuppa:", err)
		os.Exit(1)
	}
}
