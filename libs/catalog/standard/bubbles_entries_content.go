package standard

import "github.com/fancy-cuppa/cuppa/libs/catalog/definition"

const bubblesImport = "charm.land/bubbles/v2"

func bubblesEntries() []definition.Definition {
	return []definition.Definition{
		{
			ID: "bubbles.textinput", Name: "Text input", Family: definition.FamilyBubbles,
			Description: "A single-line text field.",
			DefaultSize: size(28, 1), MinSize: size(8, 1), Import: bubblesImport + "/textinput", Status: definition.StatusSupported,
			Props: []definition.PropSpec{
				textProp("prompt", "Prompt", "> "),
				textProp("placeholder", "Placeholder", "Type here"),
				textProp("value", "Value", ""),
				colorProp("color", "Color", "212"),
			},
		},
		{
			ID: "bubbles.spinner", Name: "Spinner", Family: definition.FamilyBubbles,
			Description: "An animated activity indicator.",
			DefaultSize: size(16, 1), MinSize: size(2, 1), Import: bubblesImport + "/spinner", Status: definition.StatusSupported,
			Props: []definition.PropSpec{
				choiceProp("style", "Spinner", "dot", "line", "dot", "minidot", "jump", "pulse", "points", "globe", "moon", "monkey"),
				textProp("label", "Label", "Brewing..."),
				colorProp("color", "Color", "205"),
			},
		},
		{
			ID: "bubbles.progress", Name: "Progress", Family: definition.FamilyBubbles,
			Description: "A progress bar.",
			DefaultSize: size(30, 1), MinSize: size(6, 1), Import: bubblesImport + "/progress", Status: definition.StatusSupported,
			Props: []definition.PropSpec{
				intProp("percent", "Percent", 60, 0, 100),
				boolProp("show_percentage", "Show percentage", true),
				colorProp("color", "Fill color", "212"),
			},
		},
	}
}
