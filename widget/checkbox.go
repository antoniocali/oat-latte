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

	// callerStyle and callerFocusStyle preserve the styles set by the caller
	// (via WithStyle) before any theme application. ApplyTheme always merges
	// the current theme token with these originals so that switching themes
	// fully replaces the previous theme's colours rather than accumulating
	// stale values from the prior theme.
	callerStyle      latte.Style
	callerFocusStyle latte.Style
}

// NewCheckBox creates a CheckBox with the given label.
func NewCheckBox(label string) *CheckBox {
	c := &CheckBox{label: label}
	c.EnsureID()
	// FocusStyle is intentionally NOT pre-seeded here.
	// ApplyTheme sets it from the active theme's CheckBoxFocus token.
	// Pre-seeding with latte.Focused would cause its non-zero fields to survive
	// theme switches via Merge, blocking the new theme's colours from taking effect.
	return c
}

// WithStyle sets the display style for this CheckBox.
func (c *CheckBox) WithStyle(s latte.Style) *CheckBox { c.Style = s; c.callerStyle = s; return c }

// WithID sets a user-defined identifier on this component.
func (c *CheckBox) WithID(id string) *CheckBox { c.ID = id; return c }

// GetValue implements oat.ValueGetter. Returns the checked state as a bool.
func (c *CheckBox) GetValue() interface{} { return c.checked }

// WithHAlign sets the horizontal alignment for this widget within a VBox slot.
// No argument (or HAlignFill) resets to the default fill behaviour.
func (c *CheckBox) WithHAlign(a ...oat.HAlign) *CheckBox {
	c.BaseComponent.HAlign = oat.HAlignFill
	if len(a) > 0 {
		c.BaseComponent.HAlign = a[0]
	}
	return c
}

// WithVAlign sets the vertical alignment for this widget within an HBox slot.
// No argument (or VAlignFill) resets to the default fill behaviour.
func (c *CheckBox) WithVAlign(a ...oat.VAlign) *CheckBox {
	c.BaseComponent.VAlign = oat.VAlignFill
	if len(a) > 0 {
		c.BaseComponent.VAlign = a[0]
	}
	return c
}

// WithOnToggle registers a callback invoked when the checkbox state changes.
func (c *CheckBox) WithOnToggle(fn func(bool)) *CheckBox { c.onToggle = fn; return c }

// SetChecked sets the checked state programmatically.
func (c *CheckBox) SetChecked(v bool) { c.checked = v }

// IsChecked returns the current checked state.
func (c *CheckBox) IsChecked() bool { return c.checked }

// ApplyTheme applies theme tokens to the CheckBox.
// The theme acts as the base; any style fields explicitly set by the caller
// (via WithStyle) take precedence via Merge.
// ApplyTheme always re-derives Style from the theme token merged with the
// original callerStyle so that switching themes fully replaces the previous
// theme's colours rather than accumulating stale values.
func (c *CheckBox) ApplyTheme(t latte.Theme) {
	c.Style = t.CheckBox.Merge(c.callerStyle)
	c.FocusStyle = t.CheckBoxFocus.Merge(c.callerFocusStyle)
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
