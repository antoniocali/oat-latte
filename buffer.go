package oat

import (
	"github.com/antoniocali/oat-latte/latte"
	"github.com/gdamore/tcell/v2"
)

// Buffer is an abstraction over tcell.Screen that provides
// bounds-checked, clipped cell writes. It also handles the
// conversion from latte.Style to tcell.Style.
//
// All rendering in oat-latte goes through a Buffer — no component
// ever writes directly to tcell.Screen.
type Buffer struct {
	screen tcell.Screen
	clip   Region // current clipping region
}

// newBuffer creates a Buffer wrapping the given tcell.Screen.
// The initial clip region covers the full screen.
func newBuffer(screen tcell.Screen) *Buffer {
	w, h := screen.Size()
	return &Buffer{
		screen: screen,
		clip:   Region{X: 0, Y: 0, Width: w, Height: h},
	}
}

// Sub returns a new Buffer whose clipping region is restricted to region.
// Coordinates passed to the sub-buffer are relative to region's origin.
func (b *Buffer) Sub(region Region) *Buffer {
	// Convert region to absolute screen coordinates, clipped to parent.
	absX := b.clip.X + region.X
	absY := b.clip.Y + region.Y

	// Clamp width/height so we never exceed the parent clip.
	w := region.Width
	if absX+w > b.clip.X+b.clip.Width {
		w = b.clip.X + b.clip.Width - absX
	}
	h := region.Height
	if absY+h > b.clip.Y+b.clip.Height {
		h = b.clip.Y + b.clip.Height - absY
	}
	if w < 0 {
		w = 0
	}
	if h < 0 {
		h = 0
	}

	return &Buffer{
		screen: b.screen,
		clip:   Region{X: absX, Y: absY, Width: w, Height: h},
	}
}

// SetCell writes a single rune at (x, y) relative to the buffer's clip origin.
// Out-of-bounds writes are silently dropped.
func (b *Buffer) SetCell(x, y int, ch rune, style latte.Style) {
	ax := b.clip.X + x
	ay := b.clip.Y + y
	if ax < b.clip.X || ax >= b.clip.X+b.clip.Width {
		return
	}
	if ay < b.clip.Y || ay >= b.clip.Y+b.clip.Height {
		return
	}
	b.screen.SetContent(ax, ay, ch, nil, style.ToTcell())
}

// SetCellTcell writes a cell using a raw tcell.Style (used internally for borders).
func (b *Buffer) SetCellTcell(x, y int, ch rune, style tcell.Style) {
	ax := b.clip.X + x
	ay := b.clip.Y + y
	if ax < b.clip.X || ax >= b.clip.X+b.clip.Width {
		return
	}
	if ay < b.clip.Y || ay >= b.clip.Y+b.clip.Height {
		return
	}
	b.screen.SetContent(ax, ay, ch, nil, style)
}

// Fill fills the entire buffer region with the given rune and style.
func (b *Buffer) Fill(ch rune, style latte.Style) {
	ts := style.ToTcell()
	for y := 0; y < b.clip.Height; y++ {
		for x := 0; x < b.clip.Width; x++ {
			b.screen.SetContent(b.clip.X+x, b.clip.Y+y, ch, nil, ts)
		}
	}
}

// FillBG fills the buffer region with spaces using the given background color.
func (b *Buffer) FillBG(style latte.Style) {
	b.Fill(' ', style)
}

// Width returns the width of this buffer's clip region.
func (b *Buffer) Width() int { return b.clip.Width }

// Height returns the height of this buffer's clip region.
func (b *Buffer) Height() int { return b.clip.Height }

// Region returns the current clipping region in absolute screen coordinates.
func (b *Buffer) Region() Region { return b.clip }

// DrawText writes a string starting at (x, y), clipped to the buffer bounds.
// Returns the x position after the last character written.
func (b *Buffer) DrawText(x, y int, text string, style latte.Style) int {
	ts := style.ToTcell()
	cx := x
	for _, ch := range text {
		if cx >= b.clip.Width {
			break
		}
		b.screen.SetContent(b.clip.X+cx, b.clip.Y+y, ch, nil, ts)
		cx++
	}
	return cx
}

// DrawTextAligned writes text within a fixed-width cell [x, x+width),
// aligned according to latte.Alignment.
func (b *Buffer) DrawTextAligned(x, y, width int, text string, align latte.Alignment, style latte.Style) {
	runes := []rune(text)
	textLen := len(runes)
	if textLen > width {
		runes = runes[:width]
		textLen = width
	}

	startX := x
	switch align {
	case latte.AlignCenter:
		startX = x + (width-textLen)/2
	case latte.AlignEnd:
		startX = x + width - textLen
	}

	ts := style.ToTcell()
	for i, ch := range runes {
		cx := startX + i
		if cx < x || cx >= x+width || cx >= b.clip.Width {
			continue
		}
		b.screen.SetContent(b.clip.X+cx, b.clip.Y+y, ch, nil, ts)
	}
}

// DrawBorder draws a border around the full buffer region using the given style.
func (b *Buffer) DrawBorder(borderStyle latte.BorderStyle, style latte.Style) {
	b.DrawBorderTitle(borderStyle, "", latte.Style{}, style, AnchorLeft)
}

// DrawBorderTitle draws a border and optionally stamps " title " into the top rule.
// anchor (oat.Anchor, H-axis) controls the horizontal position of the title:
// AnchorLeft (after the opening corner), AnchorCenter, or AnchorRight (before
// the closing corner). titleStyle is used for the title text; if its FG is
// ColorDefault the border FG is used.
//
// Note: Anchor is the horizontal-axis type. The title always appears in the
// top border row; there is no vertical placement variant for this function.
// For vertical positioning see oat.VAnchor (used by Divider).
func (b *Buffer) DrawBorderTitle(borderStyle latte.BorderStyle, title string, titleStyle latte.Style, style latte.Style, anchor Anchor) {
	if borderStyle == latte.BorderNone || borderStyle == latte.BorderExplicitNone || b.clip.Width < 2 || b.clip.Height < 2 {
		return
	}
	runes := borderStyle.Runes()
	bs := style.BorderTcell()

	w := b.clip.Width
	h := b.clip.Height

	// Top and bottom rows
	for x := 1; x < w-1; x++ {
		b.SetCellTcell(x, 0, runes.Top, bs)
		b.SetCellTcell(x, h-1, runes.Bottom, bs)
	}
	// Left and right columns
	for y := 1; y < h-1; y++ {
		b.SetCellTcell(0, y, runes.Left, bs)
		b.SetCellTcell(w-1, y, runes.Right, bs)
	}
	// Corners
	b.SetCellTcell(0, 0, runes.TopLeft, bs)
	b.SetCellTcell(w-1, 0, runes.TopRight, bs)
	b.SetCellTcell(0, h-1, runes.BottomLeft, bs)
	b.SetCellTcell(w-1, h-1, runes.BottomRight, bs)

	// Stamp title into the top border line.
	if title == "" || w < 6 {
		return
	}

	// Build the padded label: " Title "
	label := " " + title + " "
	labelRunes := []rune(label)

	// Available interior width: leave 2 cells for corners + 1 guard on each side.
	maxLen := w - 4
	if len(labelRunes) > maxLen {
		labelRunes = labelRunes[:maxLen]
	}

	ts := titleStyle.ToTcell()
	// If the caller didn't set a title FG, inherit the border FG.
	if titleStyle.FG == latte.ColorDefault && titleStyle.BG == latte.ColorDefault &&
		!titleStyle.Bold && !titleStyle.Italic {
		ts = bs
	}

	// Compute startX based on anchor.
	// Interior runs from x=1 to x=w-2 (inclusive). Guard of 1 on each side gives
	// usable range [2, w-2-len(label)].
	var startX int
	labelLen := len(labelRunes)
	switch anchor {
	case AnchorRight:
		startX = w - 2 - labelLen // just before the right corner
	case AnchorCenter:
		startX = 1 + (w-2-labelLen)/2
		if startX < 2 {
			startX = 2
		}
	default: // AnchorLeft
		startX = 2 // after ╭─
	}

	for i, r := range labelRunes {
		b.SetCellTcell(startX+i, 0, r, ts)
	}
}

// ShowCursor positions the terminal cursor at (x, y) within this buffer.
// Used by EditText to show the insertion point.
func (b *Buffer) ShowCursor(x, y int) {
	b.screen.ShowCursor(b.clip.X+x, b.clip.Y+y)
}

// HideCursor hides the terminal cursor.
func (b *Buffer) HideCursor() {
	b.screen.HideCursor()
}
