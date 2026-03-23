package widget

import (
	"fmt"
	"strings"

	oat "github.com/antoniocali/oat-latte"
	"github.com/antoniocali/oat-latte/latte"
)

// ProgressBar renders a horizontal progress indicator.
//
// There are no interactive keybindings; ProgressBar is a display-only widget.
type ProgressBar struct {
	oat.BaseComponent
	value       float64 // 0.0 – 1.0
	fillChar    rune
	emptyChar   rune
	showPercent bool
}

// NewProgressBar creates a ProgressBar.
func NewProgressBar() *ProgressBar {
	p := &ProgressBar{
		fillChar:    '█',
		emptyChar:   '░',
		showPercent: true,
	}
	p.EnsureID()
	return p
}

// WithStyle sets the display style for this ProgressBar.
// Call this after construction (or after ApplyTheme) to override the
// theme-supplied accent colour with a per-instance style.
func (p *ProgressBar) WithStyle(s latte.Style) *ProgressBar { p.Style = s; return p }

// SetStyle replaces the bar's display style (e.g. to change fill colour per-item).
// Deprecated: use WithStyle instead.
func (p *ProgressBar) SetStyle(s latte.Style) { p.Style = s }

// WithID sets a user-defined identifier on this component.
func (p *ProgressBar) WithID(id string) *ProgressBar { p.ID = id; return p }

// GetValue implements oat.ValueGetter. Returns the current progress as a float64 (0.0–1.0).
func (p *ProgressBar) GetValue() interface{} { return p.value }

// WithFillChar sets the rune used for the filled portion.
func (p *ProgressBar) WithFillChar(r rune) *ProgressBar { p.fillChar = r; return p }

// WithEmptyChar sets the rune used for the empty portion.
func (p *ProgressBar) WithEmptyChar(r rune) *ProgressBar { p.emptyChar = r; return p }

// WithShowPercent controls whether a percentage label is shown at the right edge.
func (p *ProgressBar) WithShowPercent(show bool) *ProgressBar { p.showPercent = show; return p }

// ApplyTheme applies theme tokens to the ProgressBar.
// The bar uses the Accent colour by default; callers that need a specific
// colour (e.g. per-priority) should call SetStyle after the theme is applied.
func (p *ProgressBar) ApplyTheme(t latte.Theme) {
	p.Style = t.Accent
}

// SetValue sets the progress value (0.0–1.0). Values are clamped.
func (p *ProgressBar) SetValue(v float64) {
	if v < 0 {
		v = 0
	}
	if v > 1 {
		v = 1
	}
	p.value = v
}

// Progress returns the current progress value as float64 (0.0–1.0).
func (p *ProgressBar) Progress() float64 { return p.value }

func (p *ProgressBar) Measure(c oat.Constraint) oat.Size {
	w := c.MaxWidth
	if w < 0 {
		w = 20
	}
	return oat.Size{Width: w, Height: 1}
}

func (p *ProgressBar) Render(buf *oat.Buffer, region oat.Region) {
	sub := buf.Sub(region)
	sub.FillBG(p.Style)

	barWidth := region.Width
	if p.showPercent {
		barWidth -= 5 // reserve "100% "
	}
	if barWidth < 1 {
		barWidth = 1
	}

	filled := int(float64(barWidth) * p.value)
	var sb strings.Builder
	for i := 0; i < barWidth; i++ {
		if i < filled {
			sb.WriteRune(p.fillChar)
		} else {
			sb.WriteRune(p.emptyChar)
		}
	}
	if p.showPercent {
		sb.WriteString(fmt.Sprintf(" %3d%%", int(p.value*100)))
	}
	sub.DrawText(0, 0, sb.String(), p.Style)
}
