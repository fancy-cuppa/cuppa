// Package screenscmd is the "cuppa screens" command line: turn a folder of
// designs into a Go package of data-agnostic screens, the way a program that
// already exists takes its views from Cuppa (ADR 0007).
package screenscmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/meta-tui/cuppa/libs/catalog/cupp"
	"github.com/meta-tui/cuppa/libs/catalog/standard"
	"github.com/meta-tui/cuppa/libs/cuppafile/disk"
	"github.com/meta-tui/cuppa/libs/cuppafile/format"
	"github.com/meta-tui/cuppa/libs/document/design"
	"github.com/meta-tui/cuppa/libs/export/gosource"
)

const usage = `usage:
  cuppa screens <design.cuppa|folder>... -o <folder> [-p <package>]

Writes one contract and one view file per design, and the shared runtime, into
<folder> (a package of its own). Files it wrote earlier are replaced and the
files of designs that are gone are removed; any other file stops the export.
`

// Run executes "cuppa screens" with the arguments after the word screens and
// returns the process exit code.
func Run(args []string, stdout, stderr io.Writer) int {
	var inputs []string
	out, pkg := "", ""
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-o", "--out":
			i++
			if i < len(args) {
				out = args[i]
			}
		case "-p", "--package":
			i++
			if i < len(args) {
				pkg = args[i]
			}
		default:
			inputs = append(inputs, args[i])
		}
	}
	if out == "" || len(inputs) == 0 {
		_, _ = fmt.Fprint(stderr, usage)
		return 2
	}
	if err := run(inputs, out, pkg, stdout); err != nil {
		_, _ = fmt.Fprintln(stderr, "cuppa screens:", err)
		return 1
	}
	return 0
}

func run(inputs []string, out, pkg string, stdout io.Writer) error {
	paths, err := designFiles(inputs)
	if err != nil {
		return err
	}
	var docs []design.Document
	var embedded []design.Embedded
	for _, path := range paths {
		doc, err := disk.Load(path)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		docs = append(docs, doc)
		embedded = append(embedded, doc.Embedded...)
	}
	if pkg == "" {
		pkg = filepath.Base(filepath.Clean(out))
	}
	cat := cupp.Adopt(standard.Default(), embedded)
	project := gosource.GenerateScreens(docs, cat, pkg)
	if err := gosource.ReplaceScreens(out, project); err != nil {
		return err
	}
	names := make([]string, 0, len(project.Files))
	for name := range project.Files {
		names = append(names, name)
	}
	sort.Strings(names)
	_, _ = fmt.Fprintf(stdout, "%d designs, %d files written to %s (package %s)\n", len(docs), len(names), out, project.Module)
	for _, n := range project.Notes {
		_, _ = fmt.Fprintln(stdout, "note:", n)
	}
	for _, r := range project.Requires {
		_, _ = fmt.Fprintln(stdout, "needs in go.mod:", r)
	}
	return nil
}

// designFiles expands folders into the designs they hold, in name order.
func designFiles(inputs []string) ([]string, error) {
	var paths []string
	for _, in := range inputs {
		info, err := os.Stat(in)
		if err != nil {
			return nil, err
		}
		if !info.IsDir() {
			paths = append(paths, in)
			continue
		}
		entries, err := os.ReadDir(in)
		if err != nil {
			return nil, err
		}
		for _, e := range entries {
			if !e.IsDir() && strings.EqualFold(filepath.Ext(e.Name()), format.Extension) {
				paths = append(paths, filepath.Join(in, e.Name()))
			}
		}
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("no %s designs found", format.Extension)
	}
	return paths, nil
}
