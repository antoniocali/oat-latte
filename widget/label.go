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
}

// NewLabel creates a Label displaying the given chip labels.
// Use WithStyle to override the tag chip appearance; by default ApplyTheme
// fills it in from the active theme.
func NewLabel(labels []string) *Label {
	l := &Label{
		labels:    labels,
		separator: '·',
	}
	l.EnsureID()
	return l
}

// WithStyle sets the display style for the chip labels.
func (l *Label) WithStyle(s latte.Style) *Label { l.Style = s; return l }

// WithID sets the widget's identifier for canvas lookup.
func (l *Label) WithID(id string) *Label { l.ID = id; return l }

// WithSeparator overrides the rune drawn between chips (default '·').
func (l *Label) WithSeparator(r rune) *Label { l.separator = r; return l }

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
// separator style.
func (l *Label) ApplyTheme(t latte.Theme) {
	l.Style = t.Tag.Merge(l.Style)
	if l.sepStyle == (latte.Style{}) {
		l.sepStyle = t.Muted
	}
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
	// Fill the entire row with the background style first so trailing cells
	// after the last chip (and the whole row when there are no labels) inherit
	// the correct background colour instead of leaking the canvas fill.
	bgStyle := latte.Style{BG: l.Style.BG}
	sub.FillBG(bgStyle)
	x := 0
	sepStr := " " + string(l.separator) + " "
	for i, lbl := range l.labels {
		chip := " " + lbl + " "
		sub.DrawText(x, 0, chip, l.Style)
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
