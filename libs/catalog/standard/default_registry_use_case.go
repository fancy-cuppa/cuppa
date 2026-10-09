package standard

import (
	"github.com/meta-tui/cuppa/libs/catalog/definition"
	"github.com/meta-tui/cuppa/libs/catalog/registry"
)

// All returns every component definition Cuppa ships with.
func All() []definition.Definition {
	var defs []definition.Definition
	defs = append(defs, lipglossEntries()...)
	defs = append(defs, bubblesEntries()...)
	defs = append(defs, huhEntries()...)
	defs = append(defs, glamourEntries()...)
	defs = append(defs, ntchartsEntries()...)
	defs = append(defs, communityEntries()...)
	return defs
}

// Default builds the registry of shipped components. It panics if the shipped
// definitions are inconsistent, which is a programming error caught by tests.
func Default() *registry.Registry {
	r, err := registry.NewWithPacks(Packs(), All()...)
	if err != nil {
		panic(err)
	}
	return r
}
