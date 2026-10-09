package format

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"encoding/json"
	"io"

	"github.com/fancy-cuppa/cuppa/libs/document/design"
)

// Decode reads a .cuppa file, migrating older versions to the current one.
func Decode(data []byte) (design.Document, error) {
	if !bytes.HasPrefix(data, []byte(Magic)) {
		return design.Document{}, ErrNotCuppa
	}
	rest := data[len(Magic):]
	if len(rest) < 2 {
		return design.Document{}, ErrCorrupt
	}
	version := binary.BigEndian.Uint16(rest)
	if version > CurrentVersion {
		return design.Document{}, ErrTooNew
	}
	if version == 0 {
		return design.Document{}, ErrCorrupt
	}
	zr, err := gzip.NewReader(bytes.NewReader(rest[2:]))
	if err != nil {
		return design.Document{}, ErrCorrupt
	}
	body, err := io.ReadAll(io.LimitReader(zr, MaxDecodedSize+1))
	if err != nil || len(body) > MaxDecodedSize {
		return design.Document{}, ErrCorrupt
	}
	body, err = migrate(body, version)
	if err != nil {
		return design.Document{}, err
	}
	var env envelope
	if err := json.Unmarshal(body, &env); err != nil {
		return design.Document{}, ErrCorrupt
	}
	return normalise(env.Document), nil
}

// normalise repairs what a hand-edited or damaged file can get wrong, so the
// editor never sees an impossible document.
func normalise(doc design.Document) design.Document {
	doc.Width, doc.Height = max(doc.Width, 1), max(doc.Height, 1)
	seen := map[design.NodeID]bool{}
	nodes := make([]design.Node, 0, len(doc.Nodes))
	for _, n := range doc.Nodes {
		if n.ID == "" || seen[n.ID] || n.Rect.W < 1 || n.Rect.H < 1 {
			continue
		}
		seen[n.ID] = true
		nodes = append(nodes, n)
	}
	doc.Nodes = nodes
	return doc
}
