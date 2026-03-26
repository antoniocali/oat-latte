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

// ── Theme builder methods ────────────────────────────────────────────────────
//
// Every method returns a new Theme by value, leaving the receiver unchanged.
// Style-typed tokens are merged using Style.Merge — non-zero fields in the
// supplied Style override the existing token; zero fields are left unchanged.
// This means you can make a targeted tweak without clobbering the rest of a
// token's properties:
//
//	// Nord but with no borders anywhere
//	borderless := latte.ThemeNord.
//	    WithPanel(latte.Style{Border: latte.BorderExplicitNone}).
//	    WithInput(latte.Style{Border: latte.BorderExplicitNone}).
//	    WithButton(latte.Style{Border: latte.BorderExplicitNone}).
//	    WithDialog(latte.Style{Border: latte.BorderExplicitNone})
//
//	// Dark theme with a hot-pink accent
//	pink := latte.ThemeDark.
//	    WithAccent(latte.Style{FG: latte.Hex("#ff69b4")}).
//	    WithFocusBorder(latte.Hex("#ff69b4")).
//	    WithName("dark-pink")

// WithName returns a copy of the theme with Name set to n.
func (t Theme) WithName(n string) Theme { t.Name = n; return t }

// WithCanvas merges s into the Canvas token and returns the updated theme.
func (t Theme) WithCanvas(s Style) Theme { t.Canvas = t.Canvas.Merge(s); return t }

// WithText merges s into the Text token and returns the updated theme.
func (t Theme) WithText(s Style) Theme { t.Text = t.Text.Merge(s); return t }

// WithMuted merges s into the Muted token and returns the updated theme.
func (t Theme) WithMuted(s Style) Theme { t.Muted = t.Muted.Merge(s); return t }

// WithAccent merges s into the Accent token and returns the updated theme.
func (t Theme) WithAccent(s Style) Theme { t.Accent = t.Accent.Merge(s); return t }

// WithSuccess merges s into the Success token and returns the updated theme.
func (t Theme) WithSuccess(s Style) Theme { t.Success = t.Success.Merge(s); return t }

// WithWarning merges s into the Warning token and returns the updated theme.
func (t Theme) WithWarning(s Style) Theme { t.Warning = t.Warning.Merge(s); return t }

// WithError merges s into the Error token and returns the updated theme.
func (t Theme) WithError(s Style) Theme { t.Error = t.Error.Merge(s); return t }

// WithPanel merges s into the Panel token and returns the updated theme.
func (t Theme) WithPanel(s Style) Theme { t.Panel = t.Panel.Merge(s); return t }

// WithPanelTitle merges s into the PanelTitle token and returns the updated theme.
func (t Theme) WithPanelTitle(s Style) Theme { t.PanelTitle = t.PanelTitle.Merge(s); return t }

// WithInput merges s into the Input token and returns the updated theme.
func (t Theme) WithInput(s Style) Theme { t.Input = t.Input.Merge(s); return t }

// WithInputFocus merges s into the InputFocus token and returns the updated theme.
func (t Theme) WithInputFocus(s Style) Theme { t.InputFocus = t.InputFocus.Merge(s); return t }

// WithListSelected merges s into the ListSelected token and returns the updated theme.
func (t Theme) WithListSelected(s Style) Theme {
	t.ListSelected = t.ListSelected.Merge(s)
	return t
}

// WithButton merges s into the Button token and returns the updated theme.
func (t Theme) WithButton(s Style) Theme { t.Button = t.Button.Merge(s); return t }

// WithButtonFocus merges s into the ButtonFocus token and returns the updated theme.
func (t Theme) WithButtonFocus(s Style) Theme { t.ButtonFocus = t.ButtonFocus.Merge(s); return t }

// WithCheckBox merges s into the CheckBox token and returns the updated theme.
func (t Theme) WithCheckBox(s Style) Theme { t.CheckBox = t.CheckBox.Merge(s); return t }

// WithCheckBoxFocus merges s into the CheckBoxFocus token and returns the updated theme.
func (t Theme) WithCheckBoxFocus(s Style) Theme {
	t.CheckBoxFocus = t.CheckBoxFocus.Merge(s)
	return t
}

// WithHeader merges s into the Header token and returns the updated theme.
func (t Theme) WithHeader(s Style) Theme { t.Header = t.Header.Merge(s); return t }

// WithFooter merges s into the Footer token and returns the updated theme.
func (t Theme) WithFooter(s Style) Theme { t.Footer = t.Footer.Merge(s); return t }

// WithFocusBorder sets the FocusBorder colour and returns the updated theme.
// Unlike the Style-typed tokens this field is a plain Color, so it is replaced
// rather than merged.
func (t Theme) WithFocusBorder(c Color) Theme { t.FocusBorder = c; return t }

// WithDialog merges s into the Dialog token and returns the updated theme.
func (t Theme) WithDialog(s Style) Theme { t.Dialog = t.Dialog.Merge(s); return t }

// WithDialogTitle merges s into the DialogTitle token and returns the updated theme.
func (t Theme) WithDialogTitle(s Style) Theme { t.DialogTitle = t.DialogTitle.Merge(s); return t }

// WithScrim merges s into the Scrim token and returns the updated theme.
func (t Theme) WithScrim(s Style) Theme { t.Scrim = t.Scrim.Merge(s); return t }

// WithTag merges s into the Tag token and returns the updated theme.
func (t Theme) WithTag(s Style) Theme { t.Tag = t.Tag.Merge(s); return t }

// WithNotificationInfo merges s into the NotificationInfo token and returns the updated theme.
func (t Theme) WithNotificationInfo(s Style) Theme {
	t.NotificationInfo = t.NotificationInfo.Merge(s)
	return t
}

// WithNotificationSuccess merges s into the NotificationSuccess token and returns the updated theme.
func (t Theme) WithNotificationSuccess(s Style) Theme {
	t.NotificationSuccess = t.NotificationSuccess.Merge(s)
	return t
}

// WithNotificationWarning merges s into the NotificationWarning token and returns the updated theme.
func (t Theme) WithNotificationWarning(s Style) Theme {
	t.NotificationWarning = t.NotificationWarning.Merge(s)
	return t
}

// WithNotificationError merges s into the NotificationError token and returns the updated theme.
func (t Theme) WithNotificationError(s Style) Theme {
	t.NotificationError = t.NotificationError.Merge(s)
	return t
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
	Canvas: Style{BG: DarkBg},

	Text:  Style{FG: ColorBrightWhite, BG: DarkBg},
	Muted: Style{FG: DarkMuted, BG: DarkBg},

	Accent:  Style{FG: DarkAccent, BG: DarkBg, Bold: true},
	Success: Style{FG: DarkSuccess, BG: DarkBg},
	Warning: Style{FG: DarkWarning, BG: DarkBg},
	Error:   Style{FG: DarkError, BG: DarkBg},

	Panel:      Style{Border: BorderRounded, BorderFG: DarkBorder},
	PanelTitle: Style{FG: DarkAccent, Bold: true},

	Input:      Style{Border: BorderSingle, BorderFG: DarkMuted, FG: ColorBrightWhite},
	InputFocus: Style{BorderFG: DarkAccent},

	ListSelected:  Style{FG: DarkAccent, Bold: true},
	Button:        Style{FG: ColorDefault, Border: BorderSingle, BorderFG: DarkMuted},
	ButtonFocus:   Style{Reverse: true, BorderFG: DarkAccent},
	CheckBox:      Style{FG: DarkSuccess},
	CheckBoxFocus: Focused,

	Header: Style{FG: DarkAccent, BG: DarkBg, Bold: true},
	Footer: Style{FG: DarkMuted, BG: DarkBg},

	FocusBorder: DarkAccent,

	Dialog:      Style{Border: BorderRounded, BorderFG: DarkAccent, BG: DarkBgElevated},
	DialogTitle: Style{FG: DarkAccent, Bold: true},
	Tag:         Style{FG: DarkBg, BG: DarkAccent, Bold: true},

	NotificationInfo:    Style{FG: DarkBg, BG: DarkAccent, Bold: true},
	NotificationSuccess: Style{FG: DarkBg, BG: DarkSuccess, Bold: true},
	NotificationWarning: Style{FG: DarkBg, BG: DarkWarning, Bold: true},
	NotificationError:   Style{FG: DarkBg, BG: DarkError, Bold: true},

	Scrim: Style{BG: DarkBgScrim},
}

// ThemeLight is a true-color light theme with a warm off-white background.
// Requires a terminal with true-color support.
var ThemeLight = Theme{
	Name:   "light",
	Canvas: Style{BG: LightBg},

	Text:  Style{FG: LightText, BG: LightBg},
	Muted: Style{FG: LightMuted, BG: LightBg},

	Accent:  Style{FG: LightAccent, BG: LightBg, Bold: true},
	Success: Style{FG: LightSuccess, BG: LightBg},
	Warning: Style{FG: LightWarning, BG: LightBg},
	Error:   Style{FG: LightError, BG: LightBg},

	Panel:      Style{Border: BorderRounded, BorderFG: LightBorder},
	PanelTitle: Style{FG: LightAccent, Bold: true},

	Input:      Style{Border: BorderSingle, BorderFG: LightBorder, FG: LightText},
	InputFocus: Style{BorderFG: LightAccent},

	ListSelected:  Style{FG: LightAccent, Bold: true},
	Button:        Style{FG: LightText, Border: BorderSingle, BorderFG: LightBorder},
	ButtonFocus:   Style{Reverse: true, BorderFG: LightAccent},
	CheckBox:      Style{FG: LightSuccess},
	CheckBoxFocus: Style{Reverse: true, Border: BorderSingle, BorderFG: LightAccent},

	Header: Style{FG: LightAccent, BG: LightBg, Bold: true},
	Footer: Style{FG: LightMuted, BG: LightBg},

	FocusBorder: LightAccent,

	Dialog:      Style{Border: BorderRounded, BorderFG: LightAccent, BG: LightBgElevated},
	DialogTitle: Style{FG: LightAccent, Bold: true},
	Tag:         Style{FG: LightBg, BG: LightAccent, Bold: true},

	NotificationInfo:    Style{FG: LightBg, BG: LightAccent, Bold: true},
	NotificationSuccess: Style{FG: LightBg, BG: LightSuccess, Bold: true},
	NotificationWarning: Style{FG: LightText, BG: LightWarning, Bold: true},
	NotificationError:   Style{FG: LightBg, BG: LightError, Bold: true},

	Scrim: Style{BG: LightBgScrim},
}

// ThemeDracula follows the popular Dracula color scheme.
// Requires a terminal with true-color support.
// Reference: https://draculatheme.com/contribute
var ThemeDracula = Theme{
	Name:   "dracula",
	Canvas: Style{BG: DraculaBg},

	Text:  Style{FG: DraculaFg, BG: DraculaBg},
	Muted: Style{FG: DraculaComment, BG: DraculaBg},

	Accent:  Style{FG: DraculaPurple, BG: DraculaBg, Bold: true},
	Success: Style{FG: DraculaGreen, BG: DraculaBg},
	Warning: Style{FG: DraculaOrange, BG: DraculaBg},
	Error:   Style{FG: DraculaRed, BG: DraculaBg},

	Panel:      Style{Border: BorderRounded, BorderFG: DraculaSelection},
	PanelTitle: Style{FG: DraculaPurple, Bold: true},

	Input:      Style{Border: BorderSingle, BorderFG: DraculaComment, FG: DraculaFg},
	InputFocus: Style{BorderFG: DraculaCyan},

	ListSelected:  Style{FG: DraculaCyan, Bold: true},
	Button:        Style{FG: DraculaFg, Border: BorderSingle, BorderFG: DraculaComment},
	ButtonFocus:   Style{Reverse: true, BorderFG: DraculaPurple},
	CheckBox:      Style{FG: DraculaGreen},
	CheckBoxFocus: Style{Reverse: true, Border: BorderSingle, BorderFG: DraculaPurple},

	Header: Style{FG: DraculaPurple, BG: DraculaBg, Bold: true},
	Footer: Style{FG: DraculaComment, BG: DraculaBg},

	FocusBorder: DraculaCyan,

	Dialog:      Style{Border: BorderRounded, BorderFG: DraculaPurple, BG: DraculaBgElevated},
	DialogTitle: Style{FG: DraculaPurple, Bold: true},
	Tag:         Style{FG: DraculaBg, BG: DraculaPurple, Bold: true},

	NotificationInfo:    Style{FG: DraculaBg, BG: DraculaCyan, Bold: true},
	NotificationSuccess: Style{FG: DraculaBg, BG: DraculaGreen, Bold: true},
	NotificationWarning: Style{FG: DraculaBg, BG: DraculaOrange, Bold: true},
	NotificationError:   Style{FG: DraculaBg, BG: DraculaRed, Bold: true},

	Scrim: Style{BG: DraculaBgScrim},
}

// ThemeNord follows the Nord color scheme (arctic, bluish palette).
// Requires a terminal with true-color support.
// Reference: https://www.nordtheme.com
var ThemeNord = Theme{
	Name:   "nord",
	Canvas: Style{BG: Nord0},

	Text:  Style{FG: Nord6, BG: Nord0},
	Muted: Style{FG: Nord3, BG: Nord0},

	Accent:  Style{FG: Nord8, BG: Nord0, Bold: true},
	Success: Style{FG: Nord14, BG: Nord0},
	Warning: Style{FG: Nord13, BG: Nord0},
	Error:   Style{FG: Nord11, BG: Nord0},

	Panel:      Style{Border: BorderRounded, BorderFG: Nord1},
	PanelTitle: Style{FG: Nord8, Bold: true},

	Input:      Style{Border: BorderSingle, BorderFG: Nord3, FG: Nord6},
	InputFocus: Style{BorderFG: Nord8},

	ListSelected:  Style{FG: Nord8, Bold: true},
	Button:        Style{FG: Nord6, Border: BorderSingle, BorderFG: Nord3},
	ButtonFocus:   Style{Reverse: true, BorderFG: Nord8},
	CheckBox:      Style{FG: Nord14},
	CheckBoxFocus: Style{Reverse: true, Border: BorderSingle, BorderFG: Nord8},

	Header: Style{FG: Nord8, BG: Nord0, Bold: true},
	Footer: Style{FG: Nord3, BG: Nord0},

	FocusBorder: Nord8,

	Dialog:      Style{Border: BorderRounded, BorderFG: Nord8, BG: NordBgElevated},
	DialogTitle: Style{FG: Nord8, Bold: true},
	Tag:         Style{FG: Nord0, BG: Nord8, Bold: true},

	NotificationInfo:    Style{FG: Nord0, BG: Nord8, Bold: true},
	NotificationSuccess: Style{FG: Nord0, BG: Nord14, Bold: true},
	NotificationWarning: Style{FG: Nord0, BG: Nord13, Bold: true},
	NotificationError:   Style{FG: Nord6, BG: Nord11, Bold: true},

	Scrim: Style{BG: NordBgScrim},
}
