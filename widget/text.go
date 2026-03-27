// Package widget provides leaf-node UI components for the oat-latte TUI framework.
package widget

import (
	"strings"
	"unicode/utf8"

	oat "github.com/antoniocali/oat-latte"
	"github.com/antoniocali/oat-latte/latte"
)

// ---- Text -----------------------------------------------------------------

// Text renders a static string. Supports optional word-wrap and scrolling.
type Text struct {
	oat.BaseComponent
	text       string
	wordWrap   bool
	scrollable bool
	scrollOff  int      // scroll offset in lines
	lines      []string // cached wrapped lines
	lastWidth  int
}

// NewText creates a Text widget with the given content.
func NewText(text string) *Text {
	t := &Text{text: text}
	t.EnsureID()
	return t
}

// WithStyle sets the display style for this Text widget.
func (t *Text) WithStyle(s latte.Style) *Text { t.Style = s; return t }

// WithID sets a user-defined identifier on this component.
func (t *Text) WithID(id string) *Text { t.ID = id; return t }

// GetValue implements oat.ValueGetter. Returns the current text as a string.
func (t *Text) GetValue() interface{} { return t.text }

// WithWordWrap enables word-wrapping.
func (t *Text) WithWordWrap(enabled bool) *Text { t.wordWrap = enabled; return t }

// WithScrollable makes the text vertically scrollable.
func (t *Text) WithScrollable(enabled bool) *Text { t.scrollable = enabled; return t }

// SetText updates the displayed text.
func (t *Text) SetText(text string) { t.text = text; t.lines = nil }

// GetText returns the current text.
func (t *Text) GetText() string { return t.text }

// ApplyTheme applies the Text semantic token from the theme.
// The theme acts as the base; any style fields already set on the widget
// (via WithStyle) take precedence via Merge.
func (t *Text) ApplyTheme(th latte.Theme) {
	t.Style = th.Text.Merge(t.Style)
}

// --- oat.Scrollable --------------------------------------------------------
func (t *Text) ScrollOffset() int  { return t.scrollOff }
func (t *Text) ContentHeight() int { return len(t.wrappedLines(t.lastWidth)) }
func (t *Text) ScrollTo(off int) {
	max := t.ContentHeight()
	if off < 0 {
		off = 0
	}
	if off >= max {
		off = max - 1
	}
	if off < 0 {
		off = 0
	}
	t.scrollOff = off
}

// Measure returns the desired size.
func (t *Text) Measure(c oat.Constraint) oat.Size {
	pad := toOatInsets(t.Style.Padding)
	innerW := c.MaxWidth - pad.Horizontal()
	if innerW < 0 {
		innerW = 0
	}
	lines := t.wrappedLines(innerW)
	h := len(lines) + pad.Vertical()
	if c.MaxHeight >= 0 && h > c.MaxHeight {
		h = c.MaxHeight
	}
	w := 0
	for _, l := range lines {
		if utf8.RuneCountInString(l) > w {
			w = utf8.RuneCountInString(l)
		}
	}
	w += pad.Horizontal()
	if c.MaxWidth >= 0 && w > c.MaxWidth {
		w = c.MaxWidth
	}
	return oat.Size{Width: w, Height: h}
}

// Render draws the text lines into the buffer.
func (t *Text) Render(buf *oat.Buffer, region oat.Region) {
	pad := toOatInsets(t.Style.Padding)
	// inner dimensions relative to the sub-buffer's origin (not the parent buffer).
	innerW := region.Width - pad.Horizontal()
	innerH := region.Height - pad.Vertical()
	if innerW < 0 {
		innerW = 0
	}
	if innerH < 0 {
		innerH = 0
	}
	sub := buf.Sub(region)

	// Draw border if specified.
	if t.Style.Border != latte.BorderNone && t.Style.Border != latte.BorderExplicitNone {
		sub.DrawBorderTitle(t.Style.Border, t.Title, latte.Style{}, t.Style, oat.AnchorLeft)
	}

	sub.FillBG(t.Style)

	lines := t.wrappedLines(innerW)
	t.lastWidth = innerW

	start := 0
	if t.scrollable {
		start = t.scrollOff
	}

	for y := 0; y < innerH; y++ {
		lineIdx := start + y
		if lineIdx >= len(lines) {
			break
		}
		sub.DrawTextAligned(pad.Left, pad.Top+y, innerW, lines[lineIdx], t.Style.TextAlign, t.Style)
	}
}

// wrappedLines returns the text split into lines, optionally word-wrapped.
func (t *Text) wrappedLines(width int) []string {
	if t.lines != nil && t.lastWidth == width {
		return t.lines
	}
	raw := strings.Split(t.text, "\n")
	if !t.wordWrap || width <= 0 {
		t.lines = raw
		return t.lines
	}
	var out []string
	for _, line := range raw {
		out = append(out, wrapLine(line, width)...)
	}
	t.lines = out
	return t.lines
}

// wrapLine splits a single line into segments of at most width runes,
// breaking on word boundaries where possible.
func wrapLine(line string, width int) []string {
	if width <= 0 {
		return []string{line}
	}
	words := strings.Fields(line)
	if len(words) == 0 {
		return []string{""}
	}
	var lines []string
	current := ""
	for _, word := range words {
		if current == "" {
			current = word
			continue
		}
		if utf8.RuneCountInString(current)+1+utf8.RuneCountInString(word) <= width {
			current += " " + word
		} else {
			lines = append(lines, current)
			current = word
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

// ---- Title ----------------------------------------------------------------

// Title renders a styled heading, optionally with a separator line below it.
type Title struct {
	oat.BaseComponent
	text      string
	separator bool
}

// NewTitle creates a Title widget.
func NewTitle(text string) *Title {
	t := &Title{text: text, separator: true}
	t.Style = latte.Title
	t.EnsureID()
	return t
}

// WithID sets a user-defined identifier on this component.
func (t *Title) WithID(id string) *Title { t.ID = id; return t }

// GetValue implements oat.ValueGetter. Returns the title text as a string.
func (t *Title) GetValue() interface{} { return t.text }

// WithStyle overrides the default title style.
func (t *Title) WithStyle(s latte.Style) *Title { t.Style = s; return t }

// ApplyTheme applies the Accent semantic token to the Title.
func (t *Title) ApplyTheme(th latte.Theme) {
	t.Style = th.Accent
}

// WithSeparator controls whether a line is drawn below the title text.
func (t *Title) WithSeparator(enabled bool) *Title { t.separator = enabled; return t }

// Measure returns the size: width = text length, height = 1 (+ 1 for separator).
func (t *Title) Measure(c oat.Constraint) oat.Size {
	h := 1
	if t.separator {
		h = 2
	}
	w := utf8.RuneCountInString(t.text)
	if c.MaxWidth >= 0 && w > c.MaxWidth {
		w = c.MaxWidth
	}
	if c.MaxHeight >= 0 && h > c.MaxHeight {
		h = c.MaxHeight
	}
	return oat.Size{Width: w, Height: h}
}

// Render draws the title text and an optional separator.
func (t *Title) Render(buf *oat.Buffer, region oat.Region) {
	sub := buf.Sub(region)
	sub.DrawText(0, 0, t.text, t.Style)
	if t.separator && region.Height > 1 {
		sep := strings.Repeat("─", region.Width)
		sepStyle := latte.Style{FG: t.Style.FG}
		sub.DrawText(0, 1, sep, sepStyle)
	}
}

// ---- helpers shared by this package --------------------------------------
