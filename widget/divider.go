package widget

import (
	"strings"

	oat "github.com/antoniocali/oat-latte"
	"github.com/antoniocali/oat-latte/latte"
)

// Axis controls the orientation of a Divider.
type Axis int

const (
	// AxisHorizontal draws a horizontal rule (─────).
	// Use inside a VBox to separate vertically-stacked items.
	AxisHorizontal Axis = iota

	// AxisVertical draws a vertical rule (│).
	// Use inside an HBox to separate horizontally-stacked items.
	AxisVertical
)

// Divider renders a single-cell-wide line that visually separates adjacent
// widgets. The orientation is controlled by Axis:
//
//   - AxisHorizontal — a full-width ─────── rule, 1 row tall.
//     Place it between items in a VBox.
//   - AxisVertical — a full-height │ column, 1 cell wide.
//     Place it between items in an HBox.
//
// The rune and color can be customised via the fluent builder methods.
//
// Basic usage:
//
//	vbox.AddChild(widget.NewDivider(widget.AxisHorizontal))
//	hbox.AddChild(widget.NewDivider(widget.AxisVertical))
type Divider struct {
	oat.BaseComponent // ID, Style, EnsureID(), EffectiveStyle()
	axis              Axis
	// rune drawn for each cell of the line.
	// Defaults: '─' for AxisHorizontal, '│' for AxisVertical.
	lineRune rune
}

// NewDivider returns a new Divider with the given axis and sensible defaults.
func NewDivider(axis Axis) *Divider {
	d := &Divider{axis: axis}
	d.EnsureID()
	switch axis {
	case AxisVertical:
		d.lineRune = '│'
	default:
		d.lineRune = '─'
	}
	return d
}

// NewHDivider is a convenience constructor for a horizontal divider.
func NewHDivider() *Divider { return NewDivider(AxisHorizontal) }

// NewVDivider is a convenience constructor for a vertical divider.
func NewVDivider() *Divider { return NewDivider(AxisVertical) }

// --- Fluent builder methods -------------------------------------------------

// WithID sets the component ID and returns the Divider for chaining.
func (d *Divider) WithID(id string) *Divider { d.ID = id; return d }

// WithStyle overrides the base style (FG color, BG color, etc.).
func (d *Divider) WithStyle(s latte.Style) *Divider { d.Style = s; return d }

// WithRune replaces the default line rune with r.
// For example, use '═' for a double horizontal rule, '┄' for dashes, etc.
func (d *Divider) WithRune(r rune) *Divider { d.lineRune = r; return d }

// --- Theme ------------------------------------------------------------------

// ApplyTheme maps the Muted token onto the divider so it renders as a
// de-emphasised separator. Caller-set style fields are preserved via Merge.
func (d *Divider) ApplyTheme(t latte.Theme) {
	d.Style = t.Muted.Merge(d.Style)
}

// --- Component interface ----------------------------------------------------

// Measure returns the desired size of the divider:
//   - AxisHorizontal: fills available width, always 1 row tall.
//   - AxisVertical: always 1 cell wide, fills available height.
//
// When the relevant constraint axis is unconstrained (-1) a fallback of 1 is
// used so the widget still renders sensibly in an unconstrained context.
func (d *Divider) Measure(c oat.Constraint) oat.Size {
	switch d.axis {
	case AxisVertical:
		h := c.MaxHeight
		if h < 0 {
			h = 1 // unconstrained fallback
		}
		return oat.Size{Width: 1, Height: h}
	default: // AxisHorizontal
		w := c.MaxWidth
		if w < 0 {
			w = 1 // unconstrained fallback
		}
		return oat.Size{Width: w, Height: 1}
	}
}

// Render draws the divider line into the allocated region.
//   - AxisHorizontal: draws a single row of lineRune repeated across the full width.
//   - AxisVertical: draws lineRune in column 0 for every row in the region.
func (d *Divider) Render(buf *oat.Buffer, region oat.Region) {
	sub := buf.Sub(region)

	switch d.axis {
	case AxisVertical:
		r := string(d.lineRune)
		for y := 0; y < region.Height; y++ {
			sub.DrawText(0, y, r, d.Style)
		}
	default: // AxisHorizontal
		line := strings.Repeat(string(d.lineRune), region.Width)
		sub.DrawText(0, 0, line, d.Style)
	}
}
