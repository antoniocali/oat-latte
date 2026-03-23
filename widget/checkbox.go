package widget

import (
	oat "github.com/antoniocali/oat-latte"
	"github.com/antoniocali/oat-latte/latte"
	"github.com/gdamore/tcell/v2"
)

// CheckBox is a boolean toggle widget.
//
// Default keybindings (when focused):
//   - Space / Enter  Toggle — flip the checked state
type CheckBox struct {
	oat.BaseComponent
	oat.FocusBehavior
	label    string
	checked  bool
	onToggle func(checked bool)
}

// NewCheckBox creates a CheckBox with the given label.
func NewCheckBox(label string) *CheckBox {
	c := &CheckBox{label: label}
	c.EnsureID()
	// Focus: Reverse swaps FG/BG making the whole row inverted — clearly
	// visible on any terminal background without needing a border.
	c.FocusStyle = latte.Focused
	return c
}

// WithStyle sets the display style for this CheckBox.
func (c *CheckBox) WithStyle(s latte.Style) *CheckBox { c.Style = s; return c }

// WithID sets a user-defined identifier on this component.
func (c *CheckBox) WithID(id string) *CheckBox { c.ID = id; return c }

// GetValue implements oat.ValueGetter. Returns the checked state as a bool.
func (c *CheckBox) GetValue() interface{} { return c.checked }

// WithOnToggle registers a callback invoked when the checkbox state changes.
func (c *CheckBox) WithOnToggle(fn func(bool)) *CheckBox { c.onToggle = fn; return c }

// SetChecked sets the checked state programmatically.
func (c *CheckBox) SetChecked(v bool) { c.checked = v }

// IsChecked returns the current checked state.
func (c *CheckBox) IsChecked() bool { return c.checked }

// ApplyTheme applies theme tokens to the CheckBox.
// The theme acts as the base; any style fields already set on the widget
// take precedence.
func (c *CheckBox) ApplyTheme(t latte.Theme) {
	c.Style = t.CheckBox.Merge(c.Style)
	c.FocusStyle = t.CheckBoxFocus.Merge(c.FocusStyle)
}

func (c *CheckBox) Measure(con oat.Constraint) oat.Size {
	w := 4 + len([]rune(c.label)) // "[ ] " + label
	if con.MaxWidth >= 0 && w > con.MaxWidth {
		w = con.MaxWidth
	}
	return oat.Size{Width: w, Height: 1}
}

func (c *CheckBox) Render(buf *oat.Buffer, region oat.Region) {
	style := c.EffectiveStyle(c.IsFocused())
	sub := buf.Sub(region)
	sub.FillBG(style)
	check := "[ ]"
	if c.checked {
		check = "[x]"
	}
	text := check + " " + c.label
	sub.DrawText(0, 0, text, style)
}

func (c *CheckBox) HandleKey(ev *oat.KeyEvent) bool {
	if ev.Key() == tcell.KeyEnter || (ev.Key() == tcell.KeyRune && ev.Rune() == ' ') {
		c.checked = !c.checked
		if c.onToggle != nil {
			c.onToggle(c.checked)
		}
		return true
	}
	return false
}

func (c *CheckBox) KeyBindings() []oat.KeyBinding {
	return []oat.KeyBinding{
		{Key: tcell.KeyRune, Rune: ' ', Label: "Space", Description: "Toggle"},
	}
}
