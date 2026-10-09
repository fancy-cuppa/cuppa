package menubar

// Action names what a menu item does; the shell maps them to behaviour.
type Action string

// Every menu action.
const (
	FileNew        Action = "file.new"
	FileOpen       Action = "file.open"
	FileSave       Action = "file.save"
	FileSaveAs     Action = "file.saveas"
	FileQuit       Action = "file.quit"
	EditUndo       Action = "edit.undo"
	EditRedo       Action = "edit.redo"
	EditDuplicate  Action = "edit.duplicate"
	EditDelete     Action = "edit.delete"
	EditGroup      Action = "edit.group"
	EditUngroup    Action = "edit.ungroup"
	EditComponent  Action = "edit.component"
	EditPacks      Action = "edit.packs"
	ViewPreview    Action = "view.preview"
	ExportPNG      Action = "export.png"
	ExportSVG      Action = "export.svg"
	ExportWebP     Action = "export.webp"
	ExportANSI     Action = "export.ansi"
	ExportText     Action = "export.text"
	HelpShortcuts  Action = "help.shortcuts"
	HelpAbout      Action = "help.about"
	nothing        Action = ""
	separatorLabel        = "-"
)

type item struct {
	label    string
	shortcut string
	action   Action
}

type menu struct {
	label string
	items []item
}

// menus is the whole menu bar, left to right.
var menus = []menu{
	{"File", []item{
		{"New", "Ctrl+N", FileNew},
		{"Open…", "Ctrl+O", FileOpen},
		{separatorLabel, "", nothing},
		{"Save", "Ctrl+S", FileSave},
		{"Save As…", "", FileSaveAs},
		{separatorLabel, "", nothing},
		{"Quit", "Ctrl+Q", FileQuit},
	}},
	{"Edit", []item{
		{"Undo", "Ctrl+Z", EditUndo},
		{"Redo", "Ctrl+Y", EditRedo},
		{separatorLabel, "", nothing},
		{"Duplicate", "", EditDuplicate},
		{"Delete", "Del", EditDelete},
		{separatorLabel, "", nothing},
		{"Group", "Ctrl+G", EditGroup},
		{"Ungroup", "Ctrl+U", EditUngroup},
		{"Save as component…", "", EditComponent},
		{separatorLabel, "", nothing},
		{"Component packs…", "", EditPacks},
	}},
	{"View", []item{
		{"Preview design", "Ctrl+P", ViewPreview},
	}},
	{"Export", []item{
		{"Image (PNG)…", "", ExportPNG},
		{"Image (SVG)…", "", ExportSVG},
		{"Image (WebP)…", "", ExportWebP},
		{separatorLabel, "", nothing},
		{"Colour text (ANSI)…", "", ExportANSI},
		{"Plain text…", "", ExportText},
	}},
	{"Help", []item{
		{"Shortcuts", "", HelpShortcuts},
		{"About Cuppa", "", HelpAbout},
	}},
}
