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
	label   string
	onPress func()
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
func (b *Button) WithStyle(s latte.Style) *Button { b.Style = s; return b }

// WithID sets a user-defined identifier on this component.
func (b *Button) WithID(id string) *Button { b.ID = id; return b }

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
	w := len([]rune(b.label)) + 4 // "[ label ]"
	if c.MaxWidth >= 0 && w > c.MaxWidth {
		w = c.MaxWidth
	}
	return oat.Size{Width: w, Height: 1}
}

func (b *Button) Render(buf *oat.Buffer, region oat.Region) {
	style := b.EffectiveStyle(b.IsFocused())
	sub := buf.Sub(region)
	sub.FillBG(style)
	text := "[ " + b.label + " ]"
	sub.DrawTextAligned(0, 0, region.Width, text, style.TextAlign, style)
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
