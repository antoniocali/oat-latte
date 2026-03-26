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
	value         float64 // 0.0 – 1.0
	fillChar      rune
	emptyChar     rune
	showPercent   bool
	percentAnchor oat.Anchor // where the "XX%" label is placed
}

// NewProgressBar creates a ProgressBar.
// By default the percentage label is shown at the left edge (AnchorLeft).
func NewProgressBar() *ProgressBar {
	p := &ProgressBar{
		fillChar:      '█',
		emptyChar:     '░',
		showPercent:   true,
		percentAnchor: oat.AnchorLeft,
	}
	p.EnsureID()
	return p
}

// WithStyle sets the display style for this ProgressBar.
func (p *ProgressBar) WithStyle(s latte.Style) *ProgressBar { p.Style = s; return p }

// SetStyle replaces the bar's display style.
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

// WithShowPercent controls whether a percentage label is shown.
// Deprecated: use WithPercentage instead (it also sets the anchor).
func (p *ProgressBar) WithShowPercent(show bool) *ProgressBar { p.showPercent = show; return p }

// WithPercentage controls whether a percentage label is rendered and where it
// appears relative to the bar. anchor is an oat.Anchor (H-axis) and is
// optional; it defaults to oat.AnchorLeft.
//
//	pb.WithPercentage(true)                    // " 42% ███░░░"  label at the left (default)
//	pb.WithPercentage(true, oat.AnchorRight)   // "███░░░░ 42%"  label at the right
//	pb.WithPercentage(true, oat.AnchorCenter)  // "███ 42% ░░░"  label centred inside the bar
//	pb.WithPercentage(false)                   // no label
//
// Note: Anchor is the horizontal-axis type. ProgressBar is always a
// horizontal widget so only H-axis positioning is relevant here.
func (p *ProgressBar) WithPercentage(show bool, anchor ...oat.Anchor) *ProgressBar {
	p.showPercent = show
	if len(anchor) > 0 {
		p.percentAnchor = anchor[0]
	}
	return p
}

// ApplyTheme applies theme tokens to the ProgressBar.
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

	total := region.Width
	label := ""
	if p.showPercent {
		label = fmt.Sprintf(" %3d%%", int(p.value*100)) // " XX%" — always 5 chars
	}
	labelLen := len([]rune(label))

	// barWidth is the number of cells dedicated to the fill/empty characters.
	barWidth := total - labelLen
	if barWidth < 1 {
		barWidth = 1
	}

	filled := int(float64(barWidth) * p.value)

	var bar strings.Builder
	for i := 0; i < barWidth; i++ {
		if i < filled {
			bar.WriteRune(p.fillChar)
		} else {
			bar.WriteRune(p.emptyChar)
		}
	}
	barStr := bar.String()

	// Compose the full line according to the percent anchor.
	var line string
	switch p.percentAnchor {
	case oat.AnchorRight:
		// bar ... label
		line = barStr + label
	case oat.AnchorCenter:
		// Split bar in half, stamp label in the middle.
		// barWidth is guaranteed ≥ 1; we put the label after the first half.
		half := barWidth / 2
		line = barStr[:half] + label + barStr[half:]
	default: // AnchorLeft
		// label ... bar
		line = label + barStr
	}

	sub.DrawText(0, 0, line, p.Style)
}
