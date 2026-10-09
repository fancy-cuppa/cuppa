package definition

// Pack is a named group of components that can be switched on or off as a
// whole, like a shape library. A pack's id is the Family its components carry.
type Pack struct {
	ID          Family
	Name        string
	Version     string
	Description string
	// Builtin packs ship with the app: they can be switched off, never removed.
	Builtin bool
}
