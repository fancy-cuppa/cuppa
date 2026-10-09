package cupp

import (
	"encoding/json"
	"fmt"
)

// ToJSON returns the pack as indented JSON, the form you edit by hand. It is
// the body of a .cupp file before compression.
func ToJSON(p Pack) ([]byte, error) {
	return json.MarshalIndent(envelope{Pack: p}, "", "  ")
}

// FromJSON reads the JSON form back and checks it like Encode does, so the
// first error is the one to fix. Unlike Decode it refuses a pack with a bad
// component instead of dropping it.
func FromJSON(data []byte) (Pack, error) {
	var env envelope
	if err := json.Unmarshal(data, &env); err != nil {
		return Pack{}, fmt.Errorf("not valid pack JSON: %w", err)
	}
	if err := env.Pack.check(); err != nil {
		return Pack{}, err
	}
	return env.Pack, nil
}
