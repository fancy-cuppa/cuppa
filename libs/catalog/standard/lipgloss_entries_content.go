package standard

import "github.com/fancy-cuppa/cuppa/libs/catalog/definition"

const lipglossImport = "charm.land/lipgloss/v2"

var borders = []string{"rounded", "normal", "thick", "double", "ascii", "hidden"}

func lipglossEntries() []definition.Definition {
	return []definition.Definition{
		{
			ID: "lipgloss.box", Name: "Box", Family: definition.FamilyLipgloss,
			Description: "A bordered block with an optional title.",
			DefaultSize: size(24, 6), MinSize: size(3, 3), Import: lipglossImport, Status: definition.StatusSupported,
			Props: []definition.PropSpec{
				textProp("title", "Title", ""),
				choiceProp("border", "Border", "rounded", borders...),
				colorProp("color", "Border color", "212"),
			},
		},
		{
			ID: "lipgloss.label", Name: "Label", Family: definition.FamilyLipgloss,
			Description: "A line of styled text.",
			DefaultSize: size(14, 1), MinSize: size(1, 1), Import: lipglossImport, Status: definition.StatusSupported,
			Props: []definition.PropSpec{
				textProp("text", "Text", "Hello, Cuppa"),
				colorProp("color", "Foreground", "255"),
				colorProp("background", "Background", ""),
				boolProp("bold", "Bold", false),
				choiceProp("align", "Align", "left", "left", "center", "right"),
			},
		},
		{
			ID: "lipgloss.list", Name: "List", Family: definition.FamilyLipgloss,
			Description: "An enumerated list (lipgloss/list).",
			DefaultSize: size(20, 5), MinSize: size(6, 1), Import: lipglossImport + "/list", Status: definition.StatusSupported,
			Props: []definition.PropSpec{
				textProp("items", "Items (comma separated)", "Tea,Milk,Sugar"),
				choiceProp("enumerator", "Enumerator", "bullet", "bullet", "arabic", "alphabet", "dash"),
				colorProp("color", "Color", "255"),
			},
		},
		{
			ID: "lipgloss.table", Name: "Table", Family: definition.FamilyLipgloss,
			Description: "A static table (lipgloss/table).",
			DefaultSize: size(32, 7), MinSize: size(8, 3), Import: lipglossImport + "/table", Status: definition.StatusSupported,
			Props: []definition.PropSpec{
				textProp("headers", "Headers (comma separated)", "Name,Kind,Qty"),
				textProp("rows", "Rows (; and , separated)", "Earl Grey,Black,2;Sencha,Green,1"),
				choiceProp("border", "Border", "normal", borders...),
				colorProp("color", "Border color", "240"),
			},
		},
		{
			ID: "lipgloss.tree", Name: "Tree", Family: definition.FamilyLipgloss,
			Description: "A tree of items (lipgloss/tree).",
			DefaultSize: size(22, 6), MinSize: size(8, 2), Import: lipglossImport + "/tree", Status: definition.StatusSupported,
			Props: []definition.PropSpec{
				textProp("root", "Root", "Pantry"),
				textProp("items", "Items (comma separated)", "Tea,Milk,Biscuits"),
				colorProp("color", "Color", "255"),
			},
		},
		{
			ID: "lipgloss.joinh", Name: "Join horizontal", Family: definition.FamilyLipgloss,
			Description: "Lays out blocks side by side (JoinHorizontal).",
			DefaultSize: size(30, 5), MinSize: size(6, 1), Import: lipglossImport, Status: definition.StatusPlaceholder,
			Props: []definition.PropSpec{
				intProp("columns", "Columns", 3, 2, 8),
				choiceProp("align", "Align", "top", "top", "center", "bottom"),
			},
		},
		{
			ID: "lipgloss.joinv", Name: "Join vertical", Family: definition.FamilyLipgloss,
			Description: "Stacks blocks (JoinVertical).",
			DefaultSize: size(20, 7), MinSize: size(4, 2), Import: lipglossImport, Status: definition.StatusPlaceholder,
			Props: []definition.PropSpec{
				intProp("rows", "Rows", 3, 2, 8),
				choiceProp("align", "Align", "left", "left", "center", "right"),
			},
		},
		{
			ID: "lipgloss.place", Name: "Place", Family: definition.FamilyLipgloss,
			Description: "Positions content inside a larger area (Place).",
			DefaultSize: size(24, 7), MinSize: size(5, 3), Import: lipglossImport, Status: definition.StatusPlaceholder,
			Props: []definition.PropSpec{
				choiceProp("horizontal", "Horizontal", "center", "left", "center", "right"),
				choiceProp("vertical", "Vertical", "center", "top", "center", "bottom"),
				textProp("text", "Content", "centered"),
			},
		},
	}
}
