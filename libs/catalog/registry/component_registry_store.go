// Package registry stores component definitions and answers lookups.
package registry

import (
	"fmt"
	"sort"
	"strings"

	"github.com/fancy-cuppa/cuppa/libs/catalog/definition"
)

// Registry is an immutable set of component definitions.
type Registry struct {
	byID map[string]definition.Definition
	list []definition.Definition
}

// New builds a registry. It fails on an empty or duplicate id.
func New(defs ...definition.Definition) (*Registry, error) {
	r := &Registry{byID: make(map[string]definition.Definition, len(defs))}
	for _, d := range defs {
		if d.ID == "" {
			return nil, fmt.Errorf("registry: definition %q has no id", d.Name)
		}
		if _, dup := r.byID[d.ID]; dup {
			return nil, fmt.Errorf("registry: duplicate id %q", d.ID)
		}
		r.byID[d.ID] = d
		r.list = append(r.list, d)
	}
	order := make(map[definition.Family]int, len(definition.FamilyOrder))
	for i, f := range definition.FamilyOrder {
		order[f] = i
	}
	sort.SliceStable(r.list, func(i, j int) bool {
		a, b := r.list[i], r.list[j]
		if order[a.Family] != order[b.Family] {
			return order[a.Family] < order[b.Family]
		}
		return strings.ToLower(a.Name) < strings.ToLower(b.Name)
	})
	return r, nil
}

// Get returns the definition with the given id.
func (r *Registry) Get(id string) (definition.Definition, bool) {
	d, ok := r.byID[id]
	return d, ok
}

// List returns every definition ordered by family, then name.
func (r *Registry) List() []definition.Definition {
	return append([]definition.Definition(nil), r.list...)
}

// Families returns the families that have at least one definition, in display order.
func (r *Registry) Families() []definition.Family {
	var out []definition.Family
	for _, f := range definition.FamilyOrder {
		if len(r.ByFamily(f)) > 0 {
			out = append(out, f)
		}
	}
	return out
}

// ByFamily returns the definitions of one family ordered by name.
func (r *Registry) ByFamily(f definition.Family) []definition.Definition {
	var out []definition.Definition
	for _, d := range r.list {
		if d.Family == f {
			out = append(out, d)
		}
	}
	return out
}

// Search returns the definitions whose name, id or description contain query,
// ignoring case. An empty query returns everything.
func (r *Registry) Search(query string) []definition.Definition {
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return r.List()
	}
	var out []definition.Definition
	for _, d := range r.list {
		hay := strings.ToLower(d.Name + " " + d.ID + " " + d.Description)
		if strings.Contains(hay, q) {
			out = append(out, d)
		}
	}
	return out
}
