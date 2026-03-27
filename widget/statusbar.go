package widget

import (
	oat "github.com/antoniocali/oat-latte"
	"github.com/antoniocali/oat-latte/latte"
)

// StatusBar renders a single row of key-hint hints, typically placed in the
// Canvas footer. Its content is updated automatically by the Canvas whenever
// focus changes.
//
// Each binding is rendered as:  [^S] Save
// where the bracketed key label is shown in bold+accent style and the
// description is shown in muted style. Bindings are separated by a dim │.
type StatusBar struct {
	oat.BaseComponent
	bindings    []oat.KeyBinding
	accentColor latte.Color // derived from theme Accent token in ApplyTheme
}

// NewStatusBar creates a StatusBar.
func NewStatusBar() *StatusBar {
	s := &StatusBar{}
	s.EnsureID()
	return s
}

// WithStyle sets the display style for this StatusBar.
func (s *StatusBar) WithStyle(st latte.Style) *StatusBar { s.Style = st; return s }

// WithID sets a user-defined identifier on this component.
func (s *StatusBar) WithID(id string) *StatusBar { s.ID = id; return s }

// SetBindings replaces the displayed key bindings. Called by Canvas on focus change.
func (s *StatusBar) SetBindings(bindings []oat.KeyBinding) {
	s.bindings = bindings
}

// ApplyTheme applies theme tokens to the StatusBar.
func (s *StatusBar) ApplyTheme(t latte.Theme) {
	s.Style = t.Footer
	s.accentColor = t.Accent.FG
}

func (s *StatusBar) Measure(c oat.Constraint) oat.Size {
	w := c.MaxWidth
	if w < 0 {
		w = 0
	}
	return oat.Size{Width: w, Height: 1}
}

// Render draws the status bar. Each binding is shown as "[Label] Description"
// with the bracketed key in bold+accent and the description in the base (muted)
// style. Bindings are separated by a dim │ divider.
//
// If a binding has no Description, only "[Label]" is rendered.
func (s *StatusBar) Render(buf *oat.Buffer, region oat.Region) {
	sub := buf.Sub(region)
	sub.FillBG(s.Style)

	// Derive accent and separator styles from the base footer style.
	// We keep the background from the footer style so colours stay consistent.
	accentFG := s.accentColor
	if accentFG == latte.ColorDefault {
		accentFG = latte.ColorBrightCyan // safe fallback for ThemeDefault (ANSI-16)
	}
	keyStyle := latte.Style{
		FG:   accentFG,
		BG:   s.Style.BG,
		Bold: true,
	}
	// If the base style has an explicit FG we use it for description text;
	// otherwise fall back to a dim colour.
	descStyle := s.Style // inherits FG (muted) and BG from footer
	sepStyle := latte.Style{
		FG: latte.ColorBrightBlack,
		BG: s.Style.BG,
	}

	x := 1 // one cell left margin
	for i, b := range s.bindings {
		if x >= region.Width {
			break
		}
		// Separator between bindings.
		if i > 0 {
			x = sub.DrawText(x, 0, " │ ", sepStyle)
		}
		// "[Label]" in bold accent.
		x = sub.DrawText(x, 0, "["+b.Label+"]", keyStyle)
		// " Description" in muted style (only if non-empty).
		if b.Description != "" {
			x = sub.DrawText(x, 0, " "+b.Description, descStyle)
		}
	}
}
