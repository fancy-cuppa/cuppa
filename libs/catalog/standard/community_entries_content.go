package standard

import "github.com/meta-tui/cuppa/libs/catalog/definition"

// Community entries are drawn as faithful approximations. Which Bubble Tea
// major version each library targets is recorded in docs/catalog/community-components.md.
func communityEntries() []definition.Definition {
	accent := colorProp("color", "Color", "212")
	return []definition.Definition{
		{
			ID: "community.bubbletable", Name: "Bubble table", Family: definition.FamilyCommunity,
			Description: "Evertras/bubble-table: paginated, sortable table with a footer.",
			DefaultSize: size(44, 9), MinSize: size(16, 5), Import: "github.com/evertras/bubble-table/table", Status: definition.StatusPlaceholder,
			Props: []definition.PropSpec{
				textProp("columns", "Columns (comma separated)", "Name,Kind,Qty"),
				textProp("rows", "Rows (; and , separated)", "Earl Grey,Black,2;Sencha,Green,1;Rooibos,Herbal,4"),
				textProp("footer", "Footer", "Page 1/3"),
				accent,
			},
		},
		{
			ID: "community.flexbox", Name: "Flex box", Family: definition.FamilyCommunity,
			Description: "76creates/stickers FlexBox: responsive grid of cells (Bubble Tea v1).",
			DefaultSize: size(40, 8), MinSize: size(10, 3), Import: "github.com/76creates/stickers/flexbox", Status: definition.StatusPlaceholder,
			Props: []definition.PropSpec{intProp("rows", "Rows", 2, 1, 6), intProp("columns", "Columns", 3, 1, 6), accent},
		},
		{
			ID: "community.boxer", Name: "Boxer layout", Family: definition.FamilyCommunity,
			Description: "treilik/bubbleboxer: a layout tree of side-by-side models (Bubble Tea v1-era).",
			DefaultSize: size(44, 10), MinSize: size(12, 3), Import: "github.com/treilik/bubbleboxer", Status: definition.StatusPlaceholder,
			Props: []definition.PropSpec{choiceProp("orientation", "Orientation", "horizontal", "horizontal", "vertical"), intProp("children", "Children", 2, 2, 5), accent},
		},
		{
			ID: "community.datepicker", Name: "Date picker", Family: definition.FamilyCommunity,
			Description: "EthanEFung/bubble-datepicker: a month calendar (Bubble Tea v1-era).",
			DefaultSize: size(26, 9), MinSize: size(22, 8), Import: "github.com/EthanEFung/bubble-datepicker", Status: definition.StatusPlaceholder,
			Props: []definition.PropSpec{textProp("month", "Month label", "October 2026"), intProp("day", "Selected day", 9, 1, 31), accent},
		},
		{
			ID: "community.overlay", Name: "Overlay / modal", Family: definition.FamilyCommunity,
			Description: "rmhubbert/bubbletea-overlay: a modal window (Bubble Tea v1; v2 users composite with Lip Gloss).",
			DefaultSize: size(36, 7), MinSize: size(12, 3), Import: "github.com/rmhubbert/bubbletea-overlay", Status: definition.StatusPlaceholder,
			Props: []definition.PropSpec{textProp("title", "Title", "Are you sure?"), textProp("body", "Body", "This cannot be undone."), accent},
		},
		{
			ID: "community.statusbar", Name: "Status bar", Family: definition.FamilyCommunity,
			Description: "knipferrc/teacup statusbar: four colored segments (Bubble Tea v1-era).",
			DefaultSize: size(60, 1), MinSize: size(12, 1), Import: "github.com/knipferrc/teacup/statusbar", Status: definition.StatusPlaceholder,
			Props: []definition.PropSpec{
				textProp("left", "Left", "main.go"), textProp("middle", "Middle", "UTF-8"),
				textProp("right", "Right", "Ln 12, Col 4"), textProp("end", "End", "100%"), accent,
			},
		},
		{
			ID: "community.filetree", Name: "File tree", Family: definition.FamilyCommunity,
			Description: "knipferrc/teacup filetree: a navigable directory tree (Bubble Tea v1-era).",
			DefaultSize: size(28, 10), MinSize: size(10, 3), Import: "github.com/knipferrc/teacup/filetree", Status: definition.StatusPlaceholder,
			Props: []definition.PropSpec{
				textProp("entries", "Entries (comma separated, / = dir)", "apps/,libs/,docs/,go.work,README.md"),
				intProp("selected", "Selected", 1, 0, 99), accent,
			},
		},
	}
}
