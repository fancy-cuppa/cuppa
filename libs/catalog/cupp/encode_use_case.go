package cupp

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"encoding/json"
	"fmt"
)

// Encode returns the bytes of a .cupp file holding p. A pack that could not be
// read back is refused here, so a bad file is never written.
func Encode(p Pack) ([]byte, error) {
	if err := p.check(); err != nil {
		return nil, err
	}
	body, err := json.Marshal(envelope{Pack: p})
	if err != nil {
		return nil, err
	}
	var out bytes.Buffer
	out.WriteString(Magic)
	if err := binary.Write(&out, binary.BigEndian, CurrentVersion); err != nil {
		return nil, err
	}
	zw := gzip.NewWriter(&out)
	if _, err := zw.Write(body); err != nil {
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, fmt.Errorf("cupp: %w", err)
	}
	return out.Bytes(), nil
}
