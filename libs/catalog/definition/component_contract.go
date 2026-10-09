// Package definition holds the contracts that describe a designable component.
package definition

import "github.com/meta-tui/cuppa/libs/document/design"

// Family groups components by where they come from.
type Family string

// Component families, in the order the palette lists them.
const (
	FamilyLipgloss  Family = "lipgloss"
	FamilyBubbles   Family = "bubbles"
	FamilyHuh       Family = "huh"
	FamilyGlamour   Family = "glamour"
	FamilyNtcharts  Family = "ntcharts"
	FamilyCommunity Family = "community"
)

// FamilyOrder is the display order of the families.
var FamilyOrder = []Family{FamilyLipgloss, FamilyBubbles, FamilyHuh, FamilyGlamour, FamilyNtcharts, FamilyCommunity}

// Title is the human-readable family name.
func (f Family) Title() string {
	switch f {
	case FamilyLipgloss:
		return "Lip Gloss"
	case FamilyBubbles:
		return "Bubbles"
	case FamilyHuh:
		return "Huh forms"
	case FamilyGlamour:
		return "Glamour"
	case FamilyNtcharts:
		return "ntcharts"
	case FamilyCommunity:
		return "Community"
	}
	return string(f)
}

// Status says how faithfully the designer can show a component.
type Status string

// Component statuses.
const (
	// StatusSupported components have a faithful preview.
	StatusSupported Status = "supported"
	// StatusPlaceholder components are drawn as a labelled approximation.
	StatusPlaceholder Status = "placeholder"
	// StatusPlanned components are known but not yet designable.
	StatusPlanned Status = "planned"
)

// PropKind decides how a property is edited and parsed.
type PropKind string

// Property kinds.
const (
	PropText   PropKind = "text"
	PropInt    PropKind = "int"
	PropBool   PropKind = "bool"
	PropColor  PropKind = "color"
	PropChoice PropKind = "choice"
)

// PropSpec describes one editable property of a component. Values are stored
// as strings and interpreted according to Kind.
type PropSpec struct {
	Key     string
	Label   string
	Kind    PropKind
	Default string
	// Choices lists the allowed values of a PropChoice property.
	Choices []string
	// Min and Max bound a PropInt property when Max > Min.
	Min, Max int
}

// Size is a width and height in cells.
type Size struct {
	W, H int
}

// Definition describes one component the designer can place.
type Definition struct {
	// ID is stable and unique, for example "bubbles.spinner".
	ID          string
	Name        string
	Family      Family
	Description string
	DefaultSize Size
	MinSize     Size
	Props       []PropSpec
	// Import is the Go import path the real component lives at.
	Import string
	Status Status
	// Inner is set on a component made of other components (from a .cupp
	// pack): it is drawn by painting these parts, not by a dedicated painter.
	Inner *design.Composite
}

// Defaults returns the default value of every property.
func (d Definition) Defaults() map[string]string {
	out := make(map[string]string, len(d.Props))
	for _, p := range d.Props {
		out[p.Key] = p.Default
	}
	return out
}

// Prop returns the specification of the property with the given key.
func (d Definition) Prop(key string) (PropSpec, bool) {
	for _, p := range d.Props {
		if p.Key == key {
			return p, true
		}
	}
	return PropSpec{}, false
}
