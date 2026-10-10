package modal

// Navigator is implemented by a dialog that takes navigation keys. The shell
// offers every key to Nav first, by its keystroke name ("tab", "shift+tab",
// "up", "down", "left", "right", "home", "end", "pgup", "pgdown", "space",
// "delete", or a single typed letter); the dialog says whether it used it.
// Keys it does not use then reach Key as typed text, as before.
type Navigator interface {
	Nav(name string) bool
}
