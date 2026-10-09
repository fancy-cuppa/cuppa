package disk

import (
	"os"

	"github.com/fancy-cuppa/cuppa/libs/cuppafile/format"
	"github.com/fancy-cuppa/cuppa/libs/document/design"
)

// Load reads the design stored at path.
func Load(path string) (design.Document, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return design.Document{}, err
	}
	return format.Decode(data)
}
