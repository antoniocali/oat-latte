package latte

// Theme is a named collection of semantic style tokens.
// Each token maps to a role in the UI rather than a specific component,
// so a single Theme value drives all components consistently.
//
// Usage — applying a built-in theme:
//
//	canvas := oat.NewCanvas(
//	    oat.WithTheme(latte.ThemeDark),
//	    oat.WithBody(myLayout),
//	)
//
// Usage — defining a custom theme:
//
//	myTheme := latte.Theme{
//	    Name:       "solarized",
//	    Canvas:     latte.Style{BG: latte.Hex("#002b36")},
//	    Text:       latte.Style{FG: latte.Hex("#839496")},
//	    // … fill in the rest
//	}
//	canvas := oat.NewCanvas(oat.WithTheme(myTheme), …)
type Theme struct {
	// Name is a human-readable identifier (e.g. "dark", "nord").
	Name string

	// Canvas is the full-screen background style.
	// Typically only BG is relevant here.
	Canvas Style

	// Text is used for plain body text (Text widget, List items).
	Text Style

	// Muted is used for secondary / de-emphasised text (labels, hints,
	// placeholder text, the status bar).
	Muted Style

	// Accent is used for highlighted / interactive elements that are not
	// in a destructive state (selected list item, focused border, titles,
	// primary action buttons).
	Accent Style

	// Success is used for positive states (done checkbox, low-priority indicator).
	Success Style

	// Warning is used for cautionary states (medium priority).
	Warning Style

	// Error is used for destructive / critical states (delete button, critical priority).
	Error Style

	// Panel is the base style for bordered container panels.
	// Typically sets Border and BorderFG (the inactive border colour).
	Panel Style

	// PanelTitle is the style for the title text stamped into a panel border.
	PanelTitle Style

	// Input is the style for single-line and multi-line text input fields.
	// Typically sets Border, BorderFG, and FG.
	Input Style

	// InputFocus is the style merged on top of Input when the field is focused.
	// Usually just changes the BorderFG to the accent colour.
	InputFocus Style

	// ListSelected is the style applied to the highlighted row in a List.
	ListSelected Style

	// Button is the base style for interactive buttons.
	Button Style

	// ButtonFocus is the style merged onto Button when focused.
	ButtonFocus Style

	// CheckBox is the base style for CheckBox widgets.
	CheckBox Style

	// CheckBoxFocus is the style merged onto CheckBox when focused.
	CheckBoxFocus Style

	// Header is the style for the Canvas header region.
	Header Style

	// Footer is the style for the Canvas footer / status bar.
	Footer Style

	// FocusBorder is the border colour used by layout.Border when any
	// descendant is focused. Applied to BorderFG automatically.
	FocusBorder Color

	// Dialog is the style for modal dialog overlays.
	// Typically includes a border, background, and foreground colour.
	Dialog Style

	// DialogTitle is the style for the title text inside a dialog border.
	DialogTitle Style

	// Tag is the style for inline badge/chip labels (e.g. the Label widget).
	Tag Style

	// NotificationInfo is the style for informational notification banners.
	NotificationInfo Style

	// NotificationSuccess is the style for success notification banners.
	NotificationSuccess Style

	// NotificationWarning is the style for warning notification banners.
	NotificationWarning Style

	// NotificationError is the style for error notification banners.
	NotificationError Style

	// Scrim is the style used to paint the full-screen dimming overlay that
	// appears behind modal dialogs.  Only the BG field matters — the scrim is
	// drawn as a solid fill of space characters.  Set BG to ColorDefault to
	// inherit the terminal default background (no visible overlay), or set it
	// to a semi-distinct color to visually separate the dialog from the rest of
	// the UI.
	Scrim Style
}

// ── Built-in themes ─────────────────────────────────────────────────────────

// ThemeDefault uses the terminal's native ANSI-16 palette.
// It works on any terminal regardless of true-color support.
var ThemeDefault = Theme{
	Name:   "default",
	Canvas: Style{BG: ColorDefault},

	Text:  Style{FG: ColorDefault},
	Muted: Style{FG: ColorBrightBlack},

	Accent:  Style{FG: ColorBrightCyan, Bold: true},
	Success: Style{FG: ColorBrightGreen},
	Warning: Style{FG: ColorBrightYellow},
	Error:   Style{FG: ColorBrightRed},

	Panel:      Style{Border: BorderRounded, BorderFG: ColorBrightBlack},
	PanelTitle: Style{FG: ColorBrightCyan, Bold: true},

	Input:      Style{Border: BorderSingle, BorderFG: ColorBrightBlack, FG: ColorDefault},
	InputFocus: Style{BorderFG: ColorBrightCyan},

	ListSelected:  Style{FG: ColorBrightCyan, Bold: true},
	Button:        Style{FG: ColorDefault, Border: BorderSingle, BorderFG: ColorBrightBlack},
	ButtonFocus:   Style{Reverse: true, BorderFG: ColorBrightCyan},
	CheckBox:      Style{FG: ColorDefault},
	CheckBoxFocus: Focused,

	Header: Style{FG: ColorBrightWhite, BG: ColorBlue, Bold: true},
	Footer: Style{FG: ColorBrightBlack, BG: ColorDefault},

	FocusBorder: ColorBrightCyan,

	Dialog:      Style{Border: BorderRounded, BorderFG: ColorBrightCyan, BG: ColorDefault},
	DialogTitle: Style{FG: ColorBrightCyan, Bold: true},
	Tag:         Style{FG: ColorBlack, BG: ColorBrightCyan},

	NotificationInfo:    Style{FG: ColorBlack, BG: ColorBrightCyan, Bold: true},
	NotificationSuccess: Style{FG: ColorBlack, BG: ColorBrightGreen, Bold: true},
	NotificationWarning: Style{FG: ColorBlack, BG: ColorBrightYellow, Bold: true},
	NotificationError:   Style{FG: ColorBlack, BG: ColorBrightRed, Bold: true},

	Scrim: Style{BG: ColorDefault},
}

// ThemeDark is a true-color dark theme with a deep navy-black background and
// a blue-cyan accent palette. Requires a terminal with true-color support.
var ThemeDark = Theme{
	Name:   "dark",
	Canvas: Style{BG: Hex("#0f1117")},

	Text:  Style{FG: ColorBrightWhite, BG: Hex("#0f1117")},
	Muted: Style{FG: Hex("#4a5068"), BG: Hex("#0f1117")},

	Accent:  Style{FG: Hex("#7c9cff"), BG: Hex("#0f1117"), Bold: true},
	Success: Style{FG: Hex("#4ade80"), BG: Hex("#0f1117")},
	Warning: Style{FG: Hex("#facc15"), BG: Hex("#0f1117")},
	Error:   Style{FG: Hex("#f87171"), BG: Hex("#0f1117")},

	Panel:      Style{Border: BorderRounded, BorderFG: Hex("#2e3247")},
	PanelTitle: Style{FG: Hex("#7c9cff"), Bold: true},

	Input:      Style{Border: BorderSingle, BorderFG: Hex("#4a5068"), FG: ColorBrightWhite},
	InputFocus: Style{BorderFG: Hex("#7c9cff")},

	ListSelected:  Style{FG: Hex("#7c9cff"), Bold: true},
	Button:        Style{FG: ColorDefault, Border: BorderSingle, BorderFG: Hex("#4a5068")},
	ButtonFocus:   Style{Reverse: true, BorderFG: Hex("#7c9cff")},
	CheckBox:      Style{FG: Hex("#4ade80")},
	CheckBoxFocus: Focused,

	Header: Style{FG: Hex("#7c9cff"), BG: Hex("#0f1117"), Bold: true},
	Footer: Style{FG: Hex("#4a5068"), BG: Hex("#0f1117")},

	FocusBorder: Hex("#7c9cff"),

	Dialog:      Style{Border: BorderRounded, BorderFG: Hex("#7c9cff"), BG: Hex("#1a1f2e")},
	DialogTitle: Style{FG: Hex("#7c9cff"), Bold: true},
	Tag:         Style{FG: Hex("#0f1117"), BG: Hex("#7c9cff"), Bold: true},

	NotificationInfo:    Style{FG: Hex("#0f1117"), BG: Hex("#7c9cff"), Bold: true},
	NotificationSuccess: Style{FG: Hex("#0f1117"), BG: Hex("#4ade80"), Bold: true},
	NotificationWarning: Style{FG: Hex("#0f1117"), BG: Hex("#facc15"), Bold: true},
	NotificationError:   Style{FG: Hex("#0f1117"), BG: Hex("#f87171"), Bold: true},

	Scrim: Style{BG: Hex("#080b10")},
}

// ThemeLight is a true-color light theme with a warm off-white background.
// Requires a terminal with true-color support.
var ThemeLight = Theme{
	Name:   "light",
	Canvas: Style{BG: Hex("#f8f5f0")},

	Text:  Style{FG: Hex("#1a1a2e"), BG: Hex("#f8f5f0")},
	Muted: Style{FG: Hex("#8e8e9a"), BG: Hex("#f8f5f0")},

	Accent:  Style{FG: Hex("#4361ee"), BG: Hex("#f8f5f0"), Bold: true},
	Success: Style{FG: Hex("#2d6a4f"), BG: Hex("#f8f5f0")},
	Warning: Style{FG: Hex("#e07b00"), BG: Hex("#f8f5f0")},
	Error:   Style{FG: Hex("#d62828"), BG: Hex("#f8f5f0")},

	Panel:      Style{Border: BorderRounded, BorderFG: Hex("#c8c4be")},
	PanelTitle: Style{FG: Hex("#4361ee"), Bold: true},

	Input:      Style{Border: BorderSingle, BorderFG: Hex("#c8c4be"), FG: Hex("#1a1a2e")},
	InputFocus: Style{BorderFG: Hex("#4361ee")},

	ListSelected:  Style{FG: Hex("#4361ee"), Bold: true},
	Button:        Style{FG: Hex("#1a1a2e"), Border: BorderSingle, BorderFG: Hex("#c8c4be")},
	ButtonFocus:   Style{Reverse: true, BorderFG: Hex("#4361ee")},
	CheckBox:      Style{FG: Hex("#2d6a4f")},
	CheckBoxFocus: Style{Reverse: true, Border: BorderSingle, BorderFG: Hex("#4361ee")},

	Header: Style{FG: Hex("#4361ee"), BG: Hex("#f8f5f0"), Bold: true},
	Footer: Style{FG: Hex("#8e8e9a"), BG: Hex("#f8f5f0")},

	FocusBorder: Hex("#4361ee"),

	Dialog:      Style{Border: BorderRounded, BorderFG: Hex("#4361ee"), BG: Hex("#eeebe5")},
	DialogTitle: Style{FG: Hex("#4361ee"), Bold: true},
	Tag:         Style{FG: Hex("#f8f5f0"), BG: Hex("#4361ee"), Bold: true},

	NotificationInfo:    Style{FG: Hex("#f8f5f0"), BG: Hex("#4361ee"), Bold: true},
	NotificationSuccess: Style{FG: Hex("#f8f5f0"), BG: Hex("#2d6a4f"), Bold: true},
	NotificationWarning: Style{FG: Hex("#1a1a2e"), BG: Hex("#e07b00"), Bold: true},
	NotificationError:   Style{FG: Hex("#f8f5f0"), BG: Hex("#d62828"), Bold: true},

	Scrim: Style{BG: Hex("#e0dcd7")},
}

// ThemeDracula follows the popular Dracula color scheme.
// Requires a terminal with true-color support.
// Reference: https://draculatheme.com/contribute
var ThemeDracula = Theme{
	Name:   "dracula",
	Canvas: Style{BG: Hex("#282a36")},

	Text:  Style{FG: Hex("#f8f8f2"), BG: Hex("#282a36")},
	Muted: Style{FG: Hex("#6272a4"), BG: Hex("#282a36")},

	Accent:  Style{FG: Hex("#bd93f9"), BG: Hex("#282a36"), Bold: true}, // purple
	Success: Style{FG: Hex("#50fa7b"), BG: Hex("#282a36")},             // green
	Warning: Style{FG: Hex("#ffb86c"), BG: Hex("#282a36")},             // orange
	Error:   Style{FG: Hex("#ff5555"), BG: Hex("#282a36")},             // red

	Panel:      Style{Border: BorderRounded, BorderFG: Hex("#44475a")},
	PanelTitle: Style{FG: Hex("#bd93f9"), Bold: true},

	Input:      Style{Border: BorderSingle, BorderFG: Hex("#6272a4"), FG: Hex("#f8f8f2")},
	InputFocus: Style{BorderFG: Hex("#8be9fd")}, // cyan

	ListSelected:  Style{FG: Hex("#8be9fd"), Bold: true},
	Button:        Style{FG: Hex("#f8f8f2"), Border: BorderSingle, BorderFG: Hex("#6272a4")},
	ButtonFocus:   Style{Reverse: true, BorderFG: Hex("#bd93f9")},
	CheckBox:      Style{FG: Hex("#50fa7b")},
	CheckBoxFocus: Style{Reverse: true, Border: BorderSingle, BorderFG: Hex("#bd93f9")},

	Header: Style{FG: Hex("#bd93f9"), BG: Hex("#282a36"), Bold: true},
	Footer: Style{FG: Hex("#6272a4"), BG: Hex("#282a36")},

	FocusBorder: Hex("#8be9fd"),

	Dialog:      Style{Border: BorderRounded, BorderFG: Hex("#bd93f9"), BG: Hex("#1e1f29")},
	DialogTitle: Style{FG: Hex("#bd93f9"), Bold: true},
	Tag:         Style{FG: Hex("#282a36"), BG: Hex("#bd93f9"), Bold: true},

	NotificationInfo:    Style{FG: Hex("#282a36"), BG: Hex("#8be9fd"), Bold: true},
	NotificationSuccess: Style{FG: Hex("#282a36"), BG: Hex("#50fa7b"), Bold: true},
	NotificationWarning: Style{FG: Hex("#282a36"), BG: Hex("#ffb86c"), Bold: true},
	NotificationError:   Style{FG: Hex("#282a36"), BG: Hex("#ff5555"), Bold: true},

	Scrim: Style{BG: Hex("#1a1b24")},
}

// ThemeNord follows the Nord color scheme (arctic, bluish palette).
// Requires a terminal with true-color support.
// Reference: https://www.nordtheme.com
var ThemeNord = Theme{
	Name:   "nord",
	Canvas: Style{BG: Hex("#2e3440")},

	Text:  Style{FG: Hex("#eceff4"), BG: Hex("#2e3440")},
	Muted: Style{FG: Hex("#4c566a"), BG: Hex("#2e3440")},

	Accent:  Style{FG: Hex("#88c0d0"), BG: Hex("#2e3440"), Bold: true}, // nord8
	Success: Style{FG: Hex("#a3be8c"), BG: Hex("#2e3440")},             // nord14
	Warning: Style{FG: Hex("#ebcb8b"), BG: Hex("#2e3440")},             // nord13
	Error:   Style{FG: Hex("#bf616a"), BG: Hex("#2e3440")},             // nord11

	Panel:      Style{Border: BorderRounded, BorderFG: Hex("#3b4252")},
	PanelTitle: Style{FG: Hex("#88c0d0"), Bold: true},

	Input:      Style{Border: BorderSingle, BorderFG: Hex("#4c566a"), FG: Hex("#eceff4")},
	InputFocus: Style{BorderFG: Hex("#88c0d0")},

	ListSelected:  Style{FG: Hex("#88c0d0"), Bold: true},
	Button:        Style{FG: Hex("#eceff4"), Border: BorderSingle, BorderFG: Hex("#4c566a")},
	ButtonFocus:   Style{Reverse: true, BorderFG: Hex("#88c0d0")},
	CheckBox:      Style{FG: Hex("#a3be8c")},
	CheckBoxFocus: Style{Reverse: true, Border: BorderSingle, BorderFG: Hex("#88c0d0")},

	Header: Style{FG: Hex("#88c0d0"), BG: Hex("#2e3440"), Bold: true},
	Footer: Style{FG: Hex("#4c566a"), BG: Hex("#2e3440")},

	FocusBorder: Hex("#88c0d0"),

	Dialog:      Style{Border: BorderRounded, BorderFG: Hex("#88c0d0"), BG: Hex("#272c36")},
	DialogTitle: Style{FG: Hex("#88c0d0"), Bold: true},
	Tag:         Style{FG: Hex("#2e3440"), BG: Hex("#88c0d0"), Bold: true},

	NotificationInfo:    Style{FG: Hex("#2e3440"), BG: Hex("#88c0d0"), Bold: true},
	NotificationSuccess: Style{FG: Hex("#2e3440"), BG: Hex("#a3be8c"), Bold: true},
	NotificationWarning: Style{FG: Hex("#2e3440"), BG: Hex("#ebcb8b"), Bold: true},
	NotificationError:   Style{FG: Hex("#eceff4"), BG: Hex("#bf616a"), Bold: true},

	Scrim: Style{BG: Hex("#1e2128")},
}
