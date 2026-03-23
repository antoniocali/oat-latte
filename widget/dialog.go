package widget

import (
	oat "github.com/antoniocali/oat-latte"
	"github.com/antoniocali/oat-latte/latte"
)

// ── DialogSize ────────────────────────────────────────────────────────────────

// DialogSizeKind distinguishes how a Dialog dimension is expressed.
type DialogSizeKind int

const (
	// DialogSizeFixed means the value is an absolute number of terminal cells.
	DialogSizeFixed DialogSizeKind = iota

	// DialogSizePercent means the value is a percentage of the available
	// terminal dimension (0–100).  The dialog is resized on every render pass
	// so it tracks terminal resizes automatically.
	DialogSizePercent
)

// DialogSize describes the desired width or height of a Dialog.
// Use the constructors DialogFixed and DialogPercent instead of building this
// struct directly.
//
//	dlg.WithSize(widget.DialogFixed(60), widget.DialogPercent(70))
//	dlg.WithSize(widget.DialogPercent(80), widget.DialogPercent(60))
type DialogSize struct {
	Kind  DialogSizeKind
	Value int // cells when Fixed; 0–100 when Percent
}

// DialogFixed returns a DialogSize that requests exactly n terminal cells.
func DialogFixed(n int) DialogSize { return DialogSize{Kind: DialogSizeFixed, Value: n} }

// DialogPercent returns a DialogSize that requests p percent of the available
// terminal dimension (clamped to 1–100).
func DialogPercent(p int) DialogSize {
	if p < 1 {
		p = 1
	}
	if p > 100 {
		p = 100
	}
	return DialogSize{Kind: DialogSizePercent, Value: p}
}

// resolve converts a DialogSize to an absolute cell count given the total
// available cells on that axis.
func (s DialogSize) resolve(available int) int {
	if s.Kind == DialogSizePercent {
		v := available * s.Value / 100
		if v < 1 {
			v = 1
		}
		return v
	}
	return s.Value
}

// ── Dialog ────────────────────────────────────────────────────────────────────

// Dialog is a modal overlay component that renders a bordered, centred panel
// on top of whatever is beneath it in the Z-stack.
//
// Dialog steals keyboard focus: while visible it intercepts all key events and
// forwards them only to its own focusable children.  The surrounding UI is
// visually dimmed via a scrim (solid fill with the theme's Scrim.BG colour) and
// is not reachable by Tab or arrow navigation.
//
// Sizing — use WithSize to control the dialog dimensions:
//
//	// 60-cell wide, 20-cell tall (fixed, the default)
//	dlg.WithMaxSize(60, 20)
//
//	// 80% of terminal width, 70% of terminal height
//	dlg.WithSize(widget.DialogPercent(80), widget.DialogPercent(70))
//
//	// Mix: fixed width, percentage height
//	dlg.WithSize(widget.DialogFixed(60), widget.DialogPercent(60))
//
// When both WithMaxSize and WithSize are called the last call wins.
//
// The dialog is always centred in the available screen area.
type Dialog struct {
	oat.BaseComponent
	child      oat.Component
	titleStyle latte.Style
	scrimStyle latte.Style // populated by ApplyTheme; defaults to BG: ColorDefault

	// sizeW / sizeH control the maximum terminal cells the dialog may occupy.
	// They are resolved at Render time so percent-based sizes track terminal
	// resizes automatically.
	sizeW DialogSize
	sizeH DialogSize
}

// NewDialog constructs a Dialog with the given title.
// Call WithStyle to set border/colour options, or leave them to ApplyTheme.
func NewDialog(title string) *Dialog {
	d := &Dialog{
		sizeW: DialogFixed(60),
		sizeH: DialogFixed(20),
	}
	d.Title = title
	if d.Style.Border == latte.BorderNone {
		d.Style.Border = latte.BorderRounded
	}
	return d
}

// WithStyle sets the visual style for the dialog border and background.
// Fields that remain zero will be filled in by ApplyTheme.
func (d *Dialog) WithStyle(s latte.Style) *Dialog {
	d.Style = s
	if d.Style.Border == latte.BorderNone {
		d.Style.Border = latte.BorderRounded
	}
	return d
}

// WithChild sets the body component rendered inside the dialog border.
func (d *Dialog) WithChild(c oat.Component) *Dialog {
	d.child = c
	return d
}

// WithID sets the dialog's identifier for canvas lookup.
func (d *Dialog) WithID(id string) *Dialog {
	d.ID = id
	return d
}

// WithTitle sets (or overrides) the dialog's title text.
func (d *Dialog) WithTitle(title string) *Dialog {
	d.Title = title
	return d
}

// WithMaxSize is a convenience wrapper for fixed pixel dimensions equivalent to
// WithSize(DialogFixed(w), DialogFixed(h)).
func (d *Dialog) WithMaxSize(w, h int) *Dialog {
	d.sizeW = DialogFixed(w)
	d.sizeH = DialogFixed(h)
	return d
}

// WithSize sets the dialog's desired width and height.
// Each dimension can be either DialogFixed(n) (exact terminal cells) or
// DialogPercent(p) (0–100% of the available terminal dimension).
//
//	dlg.WithSize(widget.DialogPercent(80), widget.DialogPercent(70))
//	dlg.WithSize(widget.DialogFixed(60), widget.DialogPercent(60))
func (d *Dialog) WithSize(w, h DialogSize) *Dialog {
	d.sizeW = w
	d.sizeH = h
	return d
}

// ApplyTheme applies Dialog, DialogTitle, and Scrim tokens from the active theme.
// Theme acts as base; fields already set on the dialog take precedence via Merge.
func (d *Dialog) ApplyTheme(t latte.Theme) {
	d.Style = t.Dialog.Merge(d.Style)
	d.titleStyle = t.DialogTitle
	d.scrimStyle = t.Scrim
}

// AddChild sets the inner child (satisfies oat.Layout).
func (d *Dialog) AddChild(c oat.Component) { d.child = c }

// Children satisfies oat.Layout.
func (d *Dialog) Children() []oat.Component {
	if d.child == nil {
		return nil
	}
	return []oat.Component{d.child}
}

// maxDimensions resolves sizeW / sizeH given the available region, returning
// the absolute (maxW, maxH) the dialog box may occupy.
func (d *Dialog) maxDimensions(available oat.Region) (int, int) {
	maxW := d.sizeW.resolve(available.Width)
	maxH := d.sizeH.resolve(available.Height)
	if maxW < 4 {
		maxW = 4
	}
	if maxH < 3 {
		maxH = 3
	}
	return maxW, maxH
}

// Measure returns the dialog's desired size: content + 2 border cells on each
// axis, clamped to the resolved max dimensions.
func (d *Dialog) Measure(c oat.Constraint) oat.Size {
	// Build a dummy Region so maxDimensions can resolve percentages.
	available := oat.Region{Width: c.MaxWidth, Height: c.MaxHeight}
	if available.Width < 0 {
		available.Width = 0
	}
	if available.Height < 0 {
		available.Height = 0
	}
	maxW, maxH := d.maxDimensions(available)

	inner := oat.Constraint{MaxWidth: maxW - 2, MaxHeight: maxH - 2}
	contentSize := oat.Size{}
	if d.child != nil {
		contentSize = d.child.Measure(inner)
	}
	w := clamp(contentSize.Width+2, 4, maxW)
	h := clamp(contentSize.Height+2, 3, maxH)
	return oat.Size{Width: w, Height: h}
}

// Render draws the dialog centred in region, painting a scrim over the rest
// of the area to visually separate the dialog from the background.
func (d *Dialog) Render(buf *oat.Buffer, region oat.Region) {
	// 1. Paint scrim — use the theme's Scrim.BG if available, else ColorDefault.
	scrim := d.scrimStyle
	buf.Sub(region).Fill(' ', scrim)

	// 2. Resolve size and centre the dialog.
	maxW, maxH := d.maxDimensions(region)
	inner := oat.Constraint{MaxWidth: maxW - 2, MaxHeight: maxH - 2}
	contentSize := oat.Size{}
	if d.child != nil {
		contentSize = d.child.Measure(inner)
	}
	w := clamp(contentSize.Width+2, 4, maxW)
	h := clamp(contentSize.Height+2, 3, maxH)

	x := region.X + (region.Width-w)/2
	y := region.Y + (region.Height-h)/2
	if x < region.X {
		x = region.X
	}
	if y < region.Y {
		y = region.Y
	}

	dialogRegion := oat.Region{X: x, Y: y, Width: w, Height: h}
	sub := buf.Sub(dialogRegion)

	// 3. Fill the dialog background.
	sub.FillBG(d.Style)

	// 4. Draw border + title.
	sub.DrawBorderTitle(d.Style.Border, d.Title, d.titleStyle, d.Style)

	// 5. Render child inside the border.
	if d.child != nil {
		innerRegion := oat.Region{
			X: 1, Y: 1,
			Width:  w - 2,
			Height: h - 2,
		}
		d.child.Render(sub, innerRegion)
	}
}
