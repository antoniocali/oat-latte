package widget

import (
	"fmt"

	oat "github.com/antoniocali/oat-latte"
	"github.com/antoniocali/oat-latte/latte"
	"github.com/gdamore/tcell/v2"
)

// Button is a clickable action trigger.
//
// Default keybindings (when focused):
//   - Enter / Space  Press — invoke the onPress callback
type Button struct {
	oat.BaseComponent
	oat.FocusBehavior
	label         string
	onPress       func()
	roundedCorner bool
}

// NewButton creates a Button with the given label and press handler.
func NewButton(label string, onPress func()) *Button {
	b := &Button{label: label, onPress: onPress}
	b.EnsureID()
	// Focus: Reverse swaps FG/BG — clearly signals the button is active
	// without requiring a surrounding border.
	b.FocusStyle = latte.Focused
	return b
}

// WithStyle sets the display style for this Button.
func (b *Button) WithStyle(s latte.Style) *Button {
	validateButtonBorder(s.Border)
	b.Style = s
	return b
}

// WithID sets a user-defined identifier on this component.
func (b *Button) WithID(id string) *Button { b.ID = id; return b }

// WithRoundedCorner controls whether the button border uses arc corners
// (╭─╮ / ╰─╯) instead of the default square corners (┌─┐ / └─┘).
//
// This only affects buttons that have a border set (either via WithStyle or
// from the active theme's ButtonFocus token). Calling WithRoundedCorner(true)
// on a button whose effective border style is BorderDouble, BorderThick, or
// BorderDashed will panic at render time — arc corner codepoints exist only
// for light-weight strokes.
//
// To use rounded corners on a button that has no border by default, first set
// the border style:
//
//	widget.NewButton("OK", fn).
//	    WithStyle(latte.Style{Border: latte.BorderSingle}).
//	    WithRoundedCorner(true)
func (b *Button) WithRoundedCorner(rounded bool) *Button {
	b.roundedCorner = rounded
	return b
}

// GetValue implements oat.ValueGetter. Returns the button's label as a string.
func (b *Button) GetValue() interface{} { return b.label }

// ApplyTheme applies theme tokens to the Button.
// The theme acts as the base; any style fields already set on the widget
// take precedence.
func (b *Button) ApplyTheme(t latte.Theme) {
	b.Style = t.Button.Merge(b.Style)
	b.FocusStyle = t.ButtonFocus.Merge(b.FocusStyle)
}

func (b *Button) Measure(c oat.Constraint) oat.Size {
	// Border presence is determined by the base (unfocused) style only.
	// FocusStyle carries colour/attribute overrides and must not affect layout shape.
	hasBorder := b.Style.Border != latte.BorderNone && b.Style.Border != latte.BorderExplicitNone
	labelW := len([]rune(b.label))
	var w, h int
	if hasBorder {
		// ╭─ label ─╮  →  2 border cols + 1 pad each side + label
		w = labelW + 4
		h = 3
	} else {
		w = labelW + 4 // "[ label ]"
		h = 1
	}
	if c.MaxWidth >= 0 && w > c.MaxWidth {
		w = c.MaxWidth
	}
	return oat.Size{Width: w, Height: h}
}

func (b *Button) Render(buf *oat.Buffer, region oat.Region) {
	// EffectiveStyle provides the colour/attribute overrides for the current
	// focus state, but border presence (shape) is always from b.Style so that
	// the layout is stable regardless of focus.
	style := b.EffectiveStyle(b.IsFocused())
	sub := buf.Sub(region)
	sub.FillBG(style)

	hasBorder := b.Style.Border != latte.BorderNone && b.Style.Border != latte.BorderExplicitNone
	if hasBorder {
		borderStyle := b.Style.Border
		if b.roundedCorner {
			switch borderStyle {
			case latte.BorderDouble, latte.BorderThick, latte.BorderDashed:
				panic(fmt.Sprintf(
					"oat-latte: Button.WithRoundedCorner(true) is not compatible with border style %d "+
						"(arc corner codepoints exist only for BorderSingle / BorderRounded); "+
						"use WithStyle(latte.Style{Border: latte.BorderRounded}) to switch style entirely",
					borderStyle,
				))
			default:
				borderStyle = latte.BorderRounded
			}
		}
		sub.DrawBorderTitle(borderStyle, "", latte.Style{}, style, oat.AnchorLeft)
		// Draw the label centred on the middle row (y=1), inside the border.
		// Use sub.Sub so the label is clipped to the button's own region and
		// never bleeds into adjacent rows if the button is given less height
		// than its measured 3 rows (e.g. in an overflow scenario).
		innerWidth := region.Width - 2
		if innerWidth > 0 {
			sub2 := sub.Sub(oat.Region{X: 1, Y: 1, Width: innerWidth, Height: 1})
			sub2.DrawTextAligned(0, 0, innerWidth, b.label, latte.AlignCenter, style)
		}
	} else {
		text := "[ " + b.label + " ]"
		sub.DrawTextAligned(0, 0, region.Width, text, style.TextAlign, style)
	}
}

func (b *Button) HandleKey(ev *oat.KeyEvent) bool {
	if ev.Key() == tcell.KeyEnter || (ev.Key() == tcell.KeyRune && ev.Rune() == ' ') {
		if b.onPress != nil {
			b.onPress()
		}
		return true
	}
	return false
}

func (b *Button) KeyBindings() []oat.KeyBinding {
	return []oat.KeyBinding{
		{Key: tcell.KeyEnter, Label: "Enter", Description: "Press"},
	}
}

// validateButtonBorder panics if the given border style is incompatible with
// Button. Only BorderNone, BorderExplicitNone, BorderSingle, and BorderRounded
// are valid; Double, Thick, and Dashed are rejected because they don't pair
// well with the button's single-row label or arc corners.
func validateButtonBorder(b latte.BorderStyle) {
	switch b {
	case latte.BorderNone, latte.BorderExplicitNone, latte.BorderSingle, latte.BorderRounded:
		// ok
	case latte.BorderDouble, latte.BorderThick, latte.BorderDashed:
		panic(fmt.Sprintf(
			"oat-latte: Button does not support border style %d; "+
				"allowed values are BorderNone, BorderExplicitNone, BorderSingle, BorderRounded",
			b,
		))
	}
}
