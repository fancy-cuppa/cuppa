// Package standard assembles the registry of components Cuppa ships with.
package standard

import (
	"strconv"

	"github.com/fancy-cuppa/cuppa/libs/catalog/definition"
)

func textProp(key, label, def string) definition.PropSpec {
	return definition.PropSpec{Key: key, Label: label, Kind: definition.PropText, Default: def}
}

func colorProp(key, label, def string) definition.PropSpec {
	return definition.PropSpec{Key: key, Label: label, Kind: definition.PropColor, Default: def}
}

func boolProp(key, label string, def bool) definition.PropSpec {
	v := "false"
	if def {
		v = "true"
	}
	return definition.PropSpec{Key: key, Label: label, Kind: definition.PropBool, Default: v}
}

func intProp(key, label string, def, lo, hi int) definition.PropSpec {
	return definition.PropSpec{Key: key, Label: label, Kind: definition.PropInt, Default: strconv.Itoa(def), Min: lo, Max: hi}
}

func choiceProp(key, label, def string, choices ...string) definition.PropSpec {
	return definition.PropSpec{Key: key, Label: label, Kind: definition.PropChoice, Default: def, Choices: choices}
}

func size(w, h int) definition.Size { return definition.Size{W: w, H: h} }
