package widget

import (
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
	label            string
	onPress          func()
	roundedCorner    bool // effective value used at render time
	roundedCornerSet bool // true when the caller called WithRoundedCorner explicitly

	// callerStyle and callerFocusStyle preserve the styles set explicitly by
	// the caller (via WithStyle / WithFocusStyle) before any theme application.
	// ApplyTheme always merges the current theme token with these originals so
	// that switching themes fully replaces the previous theme's colours.
	callerStyle      latte.Style
	callerFocusStyle latte.Style
}

// NewButton creates a Button with the given label and press handler.
func NewButton(label string, onPress func()) *Button {
	b := &Button{label: label, onPress: onPress}
	b.EnsureID()
	// FocusStyle is intentionally NOT pre-seeded here.
	// ApplyTheme sets it from the active theme's ButtonFocus token.
	// Pre-seeding with a hardcoded style (e.g. latte.Focused) would cause its
	// non-zero fields to survive theme switches via Merge, blocking the new
	// theme's colours from taking effect.
	return b
}

// WithStyle sets the display style for this Button.
func (b *Button) WithStyle(s latte.Style) *Button {
	b.Style = s
	b.callerStyle = s
	return b
}

// WithID sets a user-defined identifier on this component.
func (b *Button) WithID(id string) *Button { b.ID = id; return b }

// WithRoundedCorner controls whether the button border uses arc corners
// (╭─╮ / ╰─╯) instead of the default square corners (┌─┐ / └─┘).
//
// Once called, this explicit choice overrides the theme's RoundedCorner
// setting for this button. If the button's resolved border style is
// incompatible with arc corners (BorderDouble, BorderThick, BorderDashed)
// the rounded-corner request is silently ignored — no panic is raised.
func (b *Button) WithRoundedCorner(rounded bool) *Button {
	b.roundedCorner = rounded
	b.roundedCornerSet = true
	return b
}

// GetValue implements oat.ValueGetter. Returns the button's label as a string.
func (b *Button) GetValue() interface{} { return b.label }

// WithHAlign sets the horizontal alignment for this widget within a VBox slot.
// No argument (or HAlignFill) resets to the default fill behaviour.
func (b *Button) WithHAlign(a ...oat.HAlign) *Button {
	b.BaseComponent.HAlign = oat.HAlignFill
	if len(a) > 0 {
		b.BaseComponent.HAlign = a[0]
	}
	return b
}

// WithVAlign sets the vertical alignment for this widget within an HBox slot.
// No argument (or VAlignFill) resets to the default fill behaviour.
func (b *Button) WithVAlign(a ...oat.VAlign) *Button {
	b.BaseComponent.VAlign = oat.VAlignFill
	if len(a) > 0 {
		b.BaseComponent.VAlign = a[0]
	}
	return b
}

// ApplyTheme applies theme tokens to the Button.
// The theme acts as the base; any style fields explicitly set by the caller
// (via WithStyle / WithFocusStyle) take precedence via Merge.
// If the caller has not explicitly set rounded-corner preference via
// WithRoundedCorner, the theme's RoundedCorner field drives the behaviour.
func (b *Button) ApplyTheme(t latte.Theme) {
	b.Style = t.Button.Merge(b.callerStyle)
	b.FocusStyle = t.ButtonFocus.Merge(b.callerFocusStyle)
	if !b.roundedCornerSet {
		b.roundedCorner = t.RoundedCorner
	}
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
			// Arc corners are only compatible with light-weight single strokes.
			// For incompatible styles silently keep the original border shape.
			switch borderStyle {
			case latte.BorderDouble, latte.BorderThick, latte.BorderDashed:
				// incompatible — leave borderStyle unchanged
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
