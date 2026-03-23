// Package latte provides the style and theming system for oat-latte TUI components.
// Styles are composable and can be merged, similar to CSS cascading.
package latte

import "github.com/gdamore/tcell/v2"

// Color wraps tcell's color type for a cleaner public API.
type Color = tcell.Color

// Predefined colors — terminal's default colors (ANSI-16).
var (
	ColorDefault Color = tcell.ColorDefault

	ColorBlack   Color = tcell.ColorBlack
	ColorRed     Color = tcell.ColorRed
	ColorGreen   Color = tcell.ColorGreen
	ColorYellow  Color = tcell.ColorYellow
	ColorBlue    Color = tcell.ColorBlue
	ColorMagenta Color = tcell.ColorDarkMagenta
	ColorCyan    Color = tcell.ColorTeal
	ColorWhite   Color = tcell.ColorWhite

	// Bright variants
	ColorBrightBlack   Color = tcell.ColorDarkGray
	ColorBrightRed     Color = tcell.ColorOrangeRed
	ColorBrightGreen   Color = tcell.ColorLime
	ColorBrightYellow  Color = tcell.ColorYellow
	ColorBrightBlue    Color = tcell.ColorCornflowerBlue
	ColorBrightMagenta Color = tcell.ColorFuchsia
	ColorBrightCyan    Color = tcell.ColorAqua
	ColorBrightWhite   Color = tcell.ColorSilver
)

// RGB constructs a 24-bit RGB color.
// Values are 0–255. Requires a terminal with true-color support.
func RGB(r, g, b uint8) Color {
	return tcell.NewRGBColor(int32(r), int32(g), int32(b))
}

// Hex parses a hex color string like "#FF8800" or "FF8800".
// Falls back to ColorDefault on invalid input.
func Hex(h string) Color {
	if len(h) > 0 && h[0] == '#' {
		h = h[1:]
	}
	if len(h) != 6 {
		return ColorDefault
	}
	parse := func(s string) uint8 {
		var v uint8
		for _, c := range s {
			v <<= 4
			switch {
			case c >= '0' && c <= '9':
				v |= uint8(c - '0')
			case c >= 'a' && c <= 'f':
				v |= uint8(c-'a') + 10
			case c >= 'A' && c <= 'F':
				v |= uint8(c-'A') + 10
			}
		}
		return v
	}
	return RGB(parse(h[0:2]), parse(h[2:4]), parse(h[4:6]))
}

// BorderStyle describes the type of border drawn around a component.
type BorderStyle int

const (
	// BorderNone is the zero value — "not set / inherit from parent or theme".
	// When used in a Style that is being Merge'd, it means "leave the existing
	// border unchanged".  To explicitly suppress a border, use BorderExplicitNone.
	BorderNone    BorderStyle = iota
	BorderSingle              // ─ │ ┌ ┐ └ ┘
	BorderDouble              // ═ ║ ╔ ╗ ╚ ╝
	BorderRounded             // ─ │ ╭ ╮ ╰ ╯
	BorderThick               // ━ ┃ ┏ ┓ ┗ ┛
	BorderDashed              // ╌ ╎ ┌ ┐ └ ┘

	// BorderExplicitNone explicitly disables the border on a component.
	// Unlike BorderNone (the zero value, meaning "unset"), this value actively
	// clears any border that would otherwise be inherited from a theme or base
	// style.  Use it in a Style literal when you want no visible border:
	//
	//	widget.NewEditText(latte.Style{Border: latte.BorderExplicitNone})
	//
	// Rendering code treats BorderExplicitNone identically to BorderNone (no
	// runes are drawn), but Style.Merge propagates it to clear the receiver's
	// border.
	BorderExplicitNone BorderStyle = -1
)

// BorderRunes holds the runes for a particular border style.
type BorderRunes struct {
	Top, Bottom rune // horizontal line
	Left, Right rune // vertical line
	TopLeft     rune
	TopRight    rune
	BottomLeft  rune
	BottomRight rune
}

// Runes returns the rune set for this BorderStyle.
// BorderNone and BorderExplicitNone both return the single-line rune set as a
// safe fallback — callers must guard against BorderNone / BorderExplicitNone
// before calling Runes (see Buffer.DrawBorderTitle).
func (b BorderStyle) Runes() BorderRunes {
	switch b {
	case BorderDouble:
		return BorderRunes{'═', '═', '║', '║', '╔', '╗', '╚', '╝'}
	case BorderRounded:
		return BorderRunes{'─', '─', '│', '│', '╭', '╮', '╰', '╯'}
	case BorderThick:
		return BorderRunes{'━', '━', '┃', '┃', '┏', '┓', '┗', '┛'}
	case BorderDashed:
		return BorderRunes{'╌', '╌', '╎', '╎', '┌', '┐', '└', '┘'}
	default: // BorderSingle
		return BorderRunes{'─', '─', '│', '│', '┌', '┐', '└', '┘'}
	}
}

// Insets mirrors the oat.Insets type so latte stays self-contained.
// oat converts between the two where needed.
type Insets struct {
	Top, Right, Bottom, Left int
}

// Uniform returns an Insets with the same value on all sides.
func Uniform(n int) Insets {
	return Insets{n, n, n, n}
}

// Symmetric returns an Insets with separate vertical and horizontal values.
func Symmetric(v, h int) Insets {
	return Insets{v, h, v, h}
}

// Alignment controls text/content alignment inside a component.
type Alignment int

const (
	AlignStart  Alignment = iota // left / top
	AlignCenter                  // centered
	AlignEnd                     // right / bottom
)

// Style is the complete visual description of a component.
// Fields are pointers so zero-value Style means "inherit / use default".
// Use the builder methods for ergonomic construction.
type Style struct {
	FG        Color
	BG        Color
	Bold      bool
	Italic    bool
	Underline bool
	Blink     bool
	Reverse   bool // swap FG/BG (used for focus highlight)

	Padding Insets
	Margin  Insets

	Border   BorderStyle
	BorderFG Color
	BorderBG Color

	TextAlign Alignment
}

// ToTcell converts Style into a tcell.Style for cell-level rendering.
func (s Style) ToTcell() tcell.Style {
	ts := tcell.StyleDefault.
		Foreground(s.FG).
		Background(s.BG).
		Bold(s.Bold).
		Italic(s.Italic).
		Underline(s.Underline).
		Blink(s.Blink).
		Reverse(s.Reverse)
	return ts
}

// BorderTcell returns a tcell.Style for drawing border runes.
func (s Style) BorderTcell() tcell.Style {
	fg := s.BorderFG
	if fg == ColorDefault {
		fg = s.FG
	}
	return tcell.StyleDefault.Foreground(fg).Background(s.BorderBG)
}

// Merge returns a new Style where non-zero fields from other override s.
func (s Style) Merge(other Style) Style {
	if other.FG != ColorDefault {
		s.FG = other.FG
	}
	if other.BG != ColorDefault {
		s.BG = other.BG
	}
	if other.Bold {
		s.Bold = true
	}
	if other.Italic {
		s.Italic = true
	}
	if other.Underline {
		s.Underline = true
	}
	if other.Blink {
		s.Blink = true
	}
	if other.Reverse {
		s.Reverse = true
	}
	if other.Border != BorderNone {
		// BorderExplicitNone clears the border; any positive value sets it.
		if other.Border == BorderExplicitNone {
			s.Border = BorderNone
			s.BorderFG = ColorDefault
			s.BorderBG = ColorDefault
		} else {
			s.Border = other.Border
			s.BorderFG = other.BorderFG
			s.BorderBG = other.BorderBG
		}
	}
	if other.Padding != (Insets{}) {
		s.Padding = other.Padding
	}
	if other.Margin != (Insets{}) {
		s.Margin = other.Margin
	}
	if other.TextAlign != AlignStart {
		s.TextAlign = other.TextAlign
	}
	return s
}

// --- Fluent builder methods ------------------------------------------------

// WithFG sets the foreground color.
func (s Style) WithFG(c Color) Style { s.FG = c; return s }

// WithBG sets the background color.
func (s Style) WithBG(c Color) Style { s.BG = c; return s }

// WithBold enables bold text.
func (s Style) WithBold() Style { s.Bold = true; return s }

// WithItalic enables italic text.
func (s Style) WithItalic() Style { s.Italic = true; return s }

// WithUnderline enables underlined text.
func (s Style) WithUnderline() Style { s.Underline = true; return s }

// WithBlink enables blinking text.
func (s Style) WithBlink() Style { s.Blink = true; return s }

// WithReverse swaps FG and BG (commonly used for focus highlight).
func (s Style) WithReverse() Style { s.Reverse = true; return s }

// WithPadding sets uniform padding on all sides.
func (s Style) WithPadding(n int) Style { s.Padding = Uniform(n); return s }

// WithPaddingInsets sets padding with individual sides.
func (s Style) WithPaddingInsets(i Insets) Style { s.Padding = i; return s }

// WithMargin sets uniform margin on all sides.
func (s Style) WithMargin(n int) Style { s.Margin = Uniform(n); return s }

// WithMarginInsets sets margin with individual sides.
func (s Style) WithMarginInsets(i Insets) Style { s.Margin = i; return s }

// WithBorder sets the border style using the component's FG color.
func (s Style) WithBorder(b BorderStyle) Style { s.Border = b; return s }

// WithBorderColor sets the border foreground color explicitly.
func (s Style) WithBorderColor(c Color) Style { s.BorderFG = c; return s }

// WithTextAlign sets the horizontal text alignment.
func (s Style) WithTextAlign(a Alignment) Style { s.TextAlign = a; return s }

// --- Preset themes ---------------------------------------------------------

// Default is a plain style with terminal defaults.
var Default = Style{
	FG: ColorDefault,
	BG: ColorDefault,
}

// Focused is the style merged on top of a component's base style when it has focus.
// It sets a bright cyan border (for bordered components) and applies Reverse
// to swap FG/BG, making focus visible on any component including borderless
// ones like Button and CheckBox.
var Focused = Style{
	FG:       ColorDefault,
	BG:       ColorDefault,
	Reverse:  true,
	Border:   BorderSingle,
	BorderFG: ColorBrightCyan,
}

// Header is a suggested style for Canvas headers.
var Header = Style{
	FG:   ColorBrightWhite,
	BG:   ColorBlue,
	Bold: true,
}

// Footer is a suggested style for Canvas footers / status bars.
var Footer = Style{
	FG: ColorBrightBlack,
	BG: ColorDefault,
}

// Title is a suggested style for Title widgets.
var Title = Style{
	FG:   ColorBrightCyan,
	Bold: true,
}

// Error is a style for error messages.
var Error = Style{
	FG:   ColorBrightRed,
	Bold: true,
}

// Success is a style for success messages.
var Success = Style{
	FG: ColorBrightGreen,
}
