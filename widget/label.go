package widget

import (
	"strings"

	oat "github.com/antoniocali/oat-latte"
	"github.com/antoniocali/oat-latte/latte"
)

// Label renders a horizontal row of inline badge chips separated by a
// configurable separator rune (default '·').
//
// Each chip is rendered with the tag style (padded by one space on each side)
// and the separator is rendered with the muted style between chips.
//
// Example output with two chips and the default separator:
//
//	backend · api
//
// Typical usage:
//
//	lbl := widget.NewLabel([]string{"backend", "api"}, latte.Style{})
//	lbl.SetLabels([]string{"infra", "db", "new"})
type Label struct {
	oat.BaseComponent
	labels    []string
	separator rune
	sepStyle  latte.Style // style for the separator
	highlight bool        // whether chips render with their background colour

	// callerStyle preserves the style set by the caller (via WithStyle) before
	// any theme application. ApplyTheme always merges the current theme token
	// with this original so that switching themes fully replaces the previous
	// theme's colours rather than accumulating stale values.
	callerStyle latte.Style
}

// NewLabel creates a Label displaying the given chip labels.
// Chip highlight (background colour fill) is enabled by default.
// Use WithStyle to override the tag chip appearance; by default ApplyTheme
// fills it in from the active theme.
func NewLabel(labels []string) *Label {
	l := &Label{
		labels:    labels,
		separator: '·',
		highlight: true,
	}
	l.EnsureID()
	return l
}

// WithStyle sets the display style for the chip labels.
func (l *Label) WithStyle(s latte.Style) *Label { l.Style = s; l.callerStyle = s; return l }

// WithID sets the widget's identifier for canvas lookup.
func (l *Label) WithID(id string) *Label { l.ID = id; return l }

// WithSeparator overrides the rune drawn between chips (default '·').
func (l *Label) WithSeparator(r rune) *Label { l.separator = r; return l }

// WithHighlight controls whether chips are rendered with their background
// colour fill. Defaults to true.
// Set to false for a plain text appearance where only the foreground colour
// of the chip style is used — useful when embedding labels inside rows that
// already have a coloured background.
func (l *Label) WithHighlight(enabled bool) *Label { l.highlight = enabled; return l }

// WithHAlign sets the horizontal alignment for this widget within a VBox slot.
// No argument (or HAlignFill) resets to the default fill behaviour.
func (l *Label) WithHAlign(a ...oat.HAlign) *Label {
	l.BaseComponent.HAlign = oat.HAlignFill
	if len(a) > 0 {
		l.BaseComponent.HAlign = a[0]
	}
	return l
}

// WithVAlign sets the vertical alignment for this widget within an HBox slot.
// No argument (or VAlignFill) resets to the default fill behaviour.
func (l *Label) WithVAlign(a ...oat.VAlign) *Label {
	l.BaseComponent.VAlign = oat.VAlignFill
	if len(a) > 0 {
		l.BaseComponent.VAlign = a[0]
	}
	return l
}

// SetLabels replaces the displayed chips.
func (l *Label) SetLabels(labels []string) { l.labels = labels }

// GetLabels returns the current chip labels.
func (l *Label) GetLabels() []string { return l.labels }

// GetValue satisfies oat.ValueGetter — returns the labels joined by the separator.
func (l *Label) GetValue() interface{} {
	parts := make([]string, len(l.labels))
	for i, lbl := range l.labels {
		parts[i] = lbl
	}
	return strings.Join(parts, string(l.separator))
}

// ApplyTheme applies the Tag token as the chip style and Muted as the
// separator style. The theme acts as the base; any style fields explicitly set
// by the caller (via WithStyle) take precedence via Merge.
// ApplyTheme always re-derives Style from the theme token merged with the
// original callerStyle so that switching themes fully replaces the previous
// theme's colours rather than accumulating stale values.
func (l *Label) ApplyTheme(t latte.Theme) {
	l.Style = t.Tag.Merge(l.callerStyle)
	l.sepStyle = t.Muted
}

// Measure returns the single-row width of all chips and separators.
func (l *Label) Measure(c oat.Constraint) oat.Size {
	w := l.totalWidth()
	if c.MaxWidth > 0 && w > c.MaxWidth {
		w = c.MaxWidth
	}
	return oat.Size{Width: w, Height: 1}
}

// Render draws each chip with padding and the separators between them.
func (l *Label) Render(buf *oat.Buffer, region oat.Region) {
	sub := buf.Sub(region)
	// Fill with transparent/default background so the parent's background
	// shows through the gaps between chips and after the last chip.
	// Using l.Style.BG here would bleed the chip colour across the entire row.
	sub.FillBG(latte.Style{})

	chipStyle := l.Style
	if !l.highlight {
		// Strip background — keep only the foreground colour and text attributes.
		chipStyle.BG = latte.ColorDefault
	}

	x := 0
	sepStr := " " + string(l.separator) + " "
	for i, lbl := range l.labels {
		chip := " " + lbl + " "
		sub.DrawText(x, 0, chip, chipStyle)
		x += len([]rune(chip))
		if i < len(l.labels)-1 {
			sub.DrawText(x, 0, sepStr, l.sepStyle)
			x += len([]rune(sepStr))
		}
	}
}

// totalWidth computes the full rendered width of all chips + separators.
func (l *Label) totalWidth() int {
	if len(l.labels) == 0 {
		return 0
	}
	w := 0
	sepW := len([]rune(" " + string(l.separator) + " "))
	for _, lbl := range l.labels {
		w += len([]rune(lbl)) + 2 // chip = " label "
	}
	w += sepW * (len(l.labels) - 1)
	return w
}
