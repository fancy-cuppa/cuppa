package format

import "encoding/json"

// migration upgrades a body from one version to the next.
type migration func(body rawEnvelope) (rawEnvelope, error)

// migrations[v] turns a version-v body into a version-v+1 body. Version 1 is
// the first format, so the table starts empty; add an entry whenever
// CurrentVersion grows.
var migrations = map[uint16]migration{}

// migrate upgrades body from version `from` to CurrentVersion.
func migrate(body rawEnvelope, from uint16) (rawEnvelope, error) {
	return migrateWith(migrations, body, from, CurrentVersion)
}

func migrateWith(table map[uint16]migration, body rawEnvelope, from, to uint16) (rawEnvelope, error) {
	for v := from; v < to; v++ {
		step, ok := table[v]
		if !ok {
			return nil, ErrCorrupt
		}
		next, err := step(body)
		if err != nil {
			return nil, err
		}
		body = next
	}
	if !json.Valid(body) {
		return nil, ErrCorrupt
	}
	return body, nil
}
