package widget

import (
	"strings"

	oat "github.com/antoniocali/oat-latte"
	"github.com/antoniocali/oat-latte/latte"
)

// ── Axis ─────────────────────────────────────────────────────────────────────

// Axis controls the orientation of a Divider.
type Axis int

const (
	// AxisHorizontal draws a full-width horizontal rule (─────).
	// Use inside a VBox to visually separate vertically-stacked items.
	AxisHorizontal Axis = iota

	// AxisVertical draws a full-height vertical rule (│).
	// Use inside an HBox to visually separate horizontally-stacked items.
	AxisVertical
)

// ── DividerSize ──────────────────────────────────────────────────────────────

// DividerSizeKind distinguishes how a Divider's length is expressed.
type DividerSizeKind int

const (
	// DividerSizeFill means the divider spans the full allocated length (default).
	DividerSizeFill DividerSizeKind = iota

	// DividerSizeFixed means the divider is exactly N terminal cells long.
	DividerSizeFixed

	// DividerSizePercent means the divider is P percent of the allocated length (1–100).
	DividerSizePercent
)

// DividerSize describes how long the visible portion of a Divider should be.
//
// Use the constructors DividerFixed, DividerPercent, or the sentinel
// DividerFill instead of building this struct directly.
//
//	// Fixed 20-cell rule centred in the available space
//	widget.NewHDivider().WithMaxSize(widget.DividerFixed(20), oat.VAnchorMiddle)
//
//	// 60% of available width, anchored to the left
//	widget.NewHDivider().WithMaxSize(widget.DividerPercent(60), oat.VAnchorTop)
type DividerSize struct {
	Kind  DividerSizeKind
	Value int // cells when Fixed; 1–100 when Percent; ignored for Fill
}

// DividerFill is the default: the divider spans the full allocated length.
var DividerFill = DividerSize{Kind: DividerSizeFill}

// DividerFixed returns a DividerSize that requests exactly n terminal cells.
// n is clamped to ≥ 1.
func DividerFixed(n int) DividerSize {
	if n < 1 {
		n = 1
	}
	return DividerSize{Kind: DividerSizeFixed, Value: n}
}

// DividerPercent returns a DividerSize that requests p percent of the
// available terminal length on the divider's primary axis (1–100).
func DividerPercent(p int) DividerSize {
	if p < 1 {
		p = 1
	}
	if p > 100 {
		p = 100
	}
	return DividerSize{Kind: DividerSizePercent, Value: p}
}

// resolve converts a DividerSize to an absolute cell count given the total
// available cells on the primary axis.
func (s DividerSize) resolve(available int) int {
	switch s.Kind {
	case DividerSizePercent:
		v := available * s.Value / 100
		if v < 1 {
			v = 1
		}
		return v
	case DividerSizeFixed:
		if s.Value > available {
			return available
		}
		return s.Value
	default: // DividerSizeFill
		return available
	}
}

// ── Divider ──────────────────────────────────────────────────────────────────

// Divider renders a single-cell-wide rule that visually separates adjacent
// widgets inside a VBox or HBox.
//
// Orientation is set at construction time via Axis:
//
//   - AxisHorizontal — draws a ─────── rule, 1 row tall, filling the full
//     allocated width. Place between items in a VBox.
//   - AxisVertical   — draws a │ column, 1 cell wide, filling the full
//     allocated height. Place between items in an HBox.
//
// # Partial-length dividers
//
// By default the rule spans the entire allocated length. Use WithMaxSize to
// render only a portion of it, and the corresponding anchor to position that
// portion within the full space:
//
//	// Horizontal: 60% width, centred vertically (VAnchor positions within
//	// the 1-row space — future-proofs for padding scenarios)
//	hd := widget.NewHDivider().
//	    WithMaxSize(widget.DividerPercent(60), oat.VAnchorMiddle)
//
//	// Vertical: fixed 8-cell rule, anchored to the top of the column
//	vd := widget.NewVDivider().
//	    WithMaxSize(widget.DividerFixed(8), oat.VAnchorTop)
//
// Anchor semantics per axis:
//
//	AxisHorizontal → the size controls WIDTH; anchor is oat.Anchor (H-axis):
//	    AnchorLeft   — rule starts at the left edge of the region
//	    AnchorCenter — rule is horizontally centred
//	    AnchorRight  — rule ends at the right edge
//
//	AxisVertical → the size controls HEIGHT; anchor is oat.VAnchor (V-axis):
//	    VAnchorTop    — rule starts at the top of the region
//	    VAnchorMiddle — rule is vertically centred
//	    VAnchorBottom — rule ends at the bottom
//
// Using the wrong anchor type for the axis is prevented by the API: horizontal
// dividers have WithMaxSize(DividerSize, oat.Anchor) and vertical dividers have
// WithMaxSize(DividerSize, oat.VAnchor).
//
// # Styling
//
// The line rune defaults to '─' (horizontal) or '│' (vertical).
// Both can be overridden with WithRune. The style (FG color etc.) defaults to
// the theme's Muted token and can be overridden with WithStyle.
type Divider struct {
	oat.BaseComponent // ID, Style, EnsureID(), EffectiveStyle()

	axis     Axis
	lineRune rune

	// Primary-axis size specification. For AxisHorizontal this controls width;
	// for AxisVertical it controls height.
	size DividerSize

	// hAnchor positions the visible portion on the H-axis.
	// Only meaningful for AxisHorizontal when size is not Fill.
	hAnchor oat.Anchor

	// vAnchor positions the visible portion on the V-axis.
	// Only meaningful for AxisVertical when size is not Fill.
	vAnchor oat.VAnchor
}

// NewDivider returns a Divider with the given axis and sensible defaults.
// The default rune is '─' for horizontal and '│' for vertical.
// The default size is DividerFill (spans the full allocated space).
func NewDivider(axis Axis) *Divider {
	d := &Divider{
		axis: axis,
		size: DividerFill,
	}
	d.EnsureID()
	switch axis {
	case AxisVertical:
		d.lineRune = '│'
	default:
		d.lineRune = '─'
	}
	return d
}

// NewHDivider is a convenience constructor for AxisHorizontal.
func NewHDivider() *Divider { return NewDivider(AxisHorizontal) }

// NewVDivider is a convenience constructor for AxisVertical.
func NewVDivider() *Divider { return NewDivider(AxisVertical) }

// ── Fluent builder methods ────────────────────────────────────────────────────

// WithID sets the component identifier.
func (d *Divider) WithID(id string) *Divider { d.ID = id; return d }

// WithStyle overrides the visual style (FG color, BG color, etc.).
// Fields that remain zero are filled in by ApplyTheme.
func (d *Divider) WithStyle(s latte.Style) *Divider { d.Style = s; return d }

// WithRune replaces the default line rune.
// Examples: '═' (double horizontal), '┄' (dashed), '╍' (heavy dashed).
func (d *Divider) WithRune(r rune) *Divider { d.lineRune = r; return d }

// WithMaxSize controls how much of the allocated space is occupied by the
// visible rule, and where it is anchored within that space.
//
// This method is valid for AxisHorizontal dividers. The size controls the
// width of the rule; anchor (oat.Anchor) positions it horizontally:
//
//	AnchorLeft   — rule starts at the left edge (default)
//	AnchorCenter — rule is horizontally centred
//	AnchorRight  — rule ends at the right edge
//
// For AxisVertical dividers use WithMaxSizeV instead.
func (d *Divider) WithMaxSize(size DividerSize, anchor ...oat.Anchor) *Divider {
	d.size = size
	if len(anchor) > 0 {
		d.hAnchor = anchor[0]
	}
	return d
}

// WithMaxSizeV controls how much of the allocated space is occupied by the
// visible rule, and where it is anchored within that space.
//
// This method is valid for AxisVertical dividers. The size controls the
// height of the rule; anchor (oat.VAnchor) positions it vertically:
//
//	VAnchorTop    — rule starts at the top edge (default)
//	VAnchorMiddle — rule is vertically centred
//	VAnchorBottom — rule ends at the bottom edge
//
// For AxisHorizontal dividers use WithMaxSize instead.
func (d *Divider) WithMaxSizeV(size DividerSize, anchor ...oat.VAnchor) *Divider {
	d.size = size
	if len(anchor) > 0 {
		d.vAnchor = anchor[0]
	}
	return d
}

// ── Theme ─────────────────────────────────────────────────────────────────────

// ApplyTheme maps the Muted token onto the divider so it renders as
// de-emphasised structural chrome. Caller-set style fields win via Merge.
func (d *Divider) ApplyTheme(t latte.Theme) {
	d.Style = t.Muted.Merge(d.Style)
}

// ── Component interface ───────────────────────────────────────────────────────

// Measure returns the desired size of the divider.
//
// The divider always claims a fixed 1 cell on its secondary axis and the full
// available space on its primary axis (so it fills the container and the
// parent can place siblings either side of it):
//
//	AxisHorizontal → Width=MaxWidth (fill), Height=1
//	AxisVertical   → Width=1, Height=MaxHeight (fill)
//
// When the relevant constraint is unconstrained (-1) a fallback of 1 is used.
func (d *Divider) Measure(c oat.Constraint) oat.Size {
	switch d.axis {
	case AxisVertical:
		h := c.MaxHeight
		if h < 0 {
			h = 1
		}
		return oat.Size{Width: 1, Height: h}
	default: // AxisHorizontal
		w := c.MaxWidth
		if w < 0 {
			w = 1
		}
		return oat.Size{Width: w, Height: 1}
	}
}

// Render draws the divider line into the allocated region.
//
// When size is DividerFill the entire region is covered.
// When size is Fixed or Percent only a portion of the region is drawn,
// positioned according to hAnchor (horizontal) or vAnchor (vertical).
func (d *Divider) Render(buf *oat.Buffer, region oat.Region) {
	sub := buf.Sub(region)

	switch d.axis {
	case AxisVertical:
		d.renderVertical(sub, region)
	default:
		d.renderHorizontal(sub, region)
	}
}

// renderHorizontal draws a horizontal rule inside sub.
// The rule spans d.size.resolve(region.Width) cells, positioned by d.hAnchor.
func (d *Divider) renderHorizontal(sub *oat.Buffer, region oat.Region) {
	length := d.size.resolve(region.Width)
	startX := anchorStart(length, region.Width, d.hAnchor)

	rule := strings.Repeat(string(d.lineRune), length)
	sub.DrawText(startX, 0, rule, d.Style)
}

// renderVertical draws a vertical rule inside sub.
// The rule spans d.size.resolve(region.Height) rows, positioned by d.vAnchor.
func (d *Divider) renderVertical(sub *oat.Buffer, region oat.Region) {
	length := d.size.resolve(region.Height)
	startY := vanchorStart(length, region.Height, d.vAnchor)

	r := string(d.lineRune)
	for y := startY; y < startY+length; y++ {
		sub.DrawText(0, y, r, d.Style)
	}
}

// ── helpers ──────────────────────────────────────────────────────────────────

// anchorStart returns the starting X offset for an element of length `len`
// inside a container of width `total`, according to the given Anchor.
func anchorStart(length, total int, a oat.Anchor) int {
	switch a {
	case oat.AnchorRight:
		start := total - length
		if start < 0 {
			return 0
		}
		return start
	case oat.AnchorCenter:
		start := (total - length) / 2
		if start < 0 {
			return 0
		}
		return start
	default: // AnchorLeft
		return 0
	}
}

// vanchorStart returns the starting Y offset for an element of length `len`
// inside a container of height `total`, according to the given VAnchor.
func vanchorStart(length, total int, a oat.VAnchor) int {
	switch a {
	case oat.VAnchorBottom:
		start := total - length
		if start < 0 {
			return 0
		}
		return start
	case oat.VAnchorMiddle:
		start := (total - length) / 2
		if start < 0 {
			return 0
		}
		return start
	default: // VAnchorTop
		return 0
	}
}
