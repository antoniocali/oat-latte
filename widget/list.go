package widget

import (
	oat "github.com/antoniocali/oat-latte"
	"github.com/antoniocali/oat-latte/latte"
	"github.com/gdamore/tcell/v2"
)

// ListItem is a single entry in a List.
type ListItem struct {
	Label string
	Value interface{}
}

// List is a vertically scrollable list of selectable items.
//
// Default keybindings (when focused):
//   - ↑ / ↓    Move — navigate up and down
//   - Enter     Select — invoke the onSelect callback
//   - Del       Delete — invoke the onDelete callback (if set)
//   - Home / ^A Top — jump to first item
//   - End  / ^E Bottom — jump to last item
type List struct {
	oat.BaseComponent
	oat.FocusBehavior

	items          []ListItem
	selected       int // currently highlighted index
	scrollOff      int // first visible item index
	onSelect       func(index int, item ListItem)
	onDelete       func(index int, item ListItem)
	onCursorChange func(index int, item ListItem)
	selectedStyle  latte.Style

	// highlight controls whether the selected row is filled with selectedStyle.
	// Defaults to true.
	highlight bool

	// cursor is the rune drawn in the gutter next to the selected row.
	// Defaults to ">". Set to "" to hide it entirely.
	cursor string
}

// NewList creates a List with the given items.
func NewList(items []ListItem) *List {
	l := &List{
		items:     items,
		highlight: true,
		cursor:    ">",
	}
	l.EnsureID()
	// Focus style: if List has a border, highlight it; otherwise the cursor
	// and selected row are sufficient inline cues.
	l.FocusStyle = latte.Style{
		BorderFG: latte.ColorBrightCyan,
	}
	l.selectedStyle = latte.Style{
		FG:   latte.ColorDefault,
		BG:   latte.ColorBlue,
		Bold: true,
	}
	return l
}

// WithStyle sets the display style for this List.
func (l *List) WithStyle(s latte.Style) *List { l.Style = s; return l }

// WithID sets a user-defined identifier on this component.
func (l *List) WithID(id string) *List { l.ID = id; return l }

// WithHighlight controls whether the selected row is filled with the
// selected-item style (background colour + bold). Defaults to true.
// Set to false when you want cursor-only indication without a background fill,
// for example in transparent or minimal UIs.
func (l *List) WithHighlight(enabled bool) *List { l.highlight = enabled; return l }

// WithCursor sets the gutter character drawn next to the selected row.
// The default is ">". Common alternatives: "▶", "→", "•", "❯", "*".
// Pass an empty string "" to hide the cursor entirely.
func (l *List) WithCursor(cursor string) *List { l.cursor = cursor; return l }

// GetValue implements oat.ValueGetter. Returns the Value field of the currently
// selected ListItem, or nil if the list is empty.
func (l *List) GetValue() interface{} {
	item, ok := l.SelectedItem()
	if !ok {
		return nil
	}
	return item.Value
}

// WithSelectedStyle overrides the style for the highlighted item.
func (l *List) WithSelectedStyle(s latte.Style) *List { l.selectedStyle = s; return l }

// WithOnSelect registers a callback when an item is confirmed (Enter pressed).
func (l *List) WithOnSelect(fn func(int, ListItem)) *List { l.onSelect = fn; return l }

// WithOnCursorChange registers a callback invoked whenever the highlighted
// index changes (Up, Down, Home, End). This fires on every cursor move,
// before the user confirms with Enter, making it suitable for live-preview
// panels that should update as the user browses the list.
func (l *List) WithOnCursorChange(fn func(int, ListItem)) *List {
	l.onCursorChange = fn
	return l
}

// WithOnDelete registers a callback when the Delete key is pressed on an item.
// The callback receives the current index and item. It is advertised in the
// StatusBar as "Del Delete".
func (l *List) WithOnDelete(fn func(int, ListItem)) *List { l.onDelete = fn; return l }

// ApplyTheme applies theme tokens to the List.
// The theme acts as the base; any style fields already set on the widget
// take precedence.
func (l *List) ApplyTheme(t latte.Theme) {
	l.Style = t.Text.Merge(l.Style)
	l.FocusStyle = latte.Style{BorderFG: t.FocusBorder}.Merge(l.FocusStyle)
	l.selectedStyle = t.ListSelected.Merge(l.selectedStyle)
}

// SetItems replaces the list contents.
func (l *List) SetItems(items []ListItem) {
	l.items = items
	if l.selected >= len(items) {
		l.selected = len(items) - 1
	}
	if l.selected < 0 {
		l.selected = 0
	}
}

// SelectedIndex returns the currently highlighted index.
func (l *List) SelectedIndex() int { return l.selected }

// SelectedItem returns the currently highlighted item, or zero value if empty.
func (l *List) SelectedItem() (ListItem, bool) {
	if l.selected < 0 || l.selected >= len(l.items) {
		return ListItem{}, false
	}
	return l.items[l.selected], true
}

// --- oat.Scrollable --------------------------------------------------------
func (l *List) ScrollOffset() int  { return l.scrollOff }
func (l *List) ContentHeight() int { return len(l.items) }
func (l *List) ScrollTo(off int) {
	if off < 0 {
		off = 0
	}
	if off >= len(l.items) {
		off = len(l.items) - 1
	}
	l.scrollOff = off
}

func (l *List) Measure(c oat.Constraint) oat.Size {
	h := len(l.items)
	if c.MaxHeight >= 0 && h > c.MaxHeight {
		h = c.MaxHeight
	}
	w := c.MaxWidth
	if w < 0 {
		w = 20
	}
	return oat.Size{Width: w, Height: h}
}

func (l *List) Render(buf *oat.Buffer, region oat.Region) {
	style := l.EffectiveStyle(l.IsFocused())
	sub := buf.Sub(region)
	sub.FillBG(style)

	borderInset := 0
	if style.Border != latte.BorderNone && style.Border != latte.BorderExplicitNone {
		sub.DrawBorderTitle(style.Border, l.Title, latte.Style{}, style)
		borderInset = 1
	}

	pad := toOatInsets(style.Padding)
	inner := oat.Region{
		X:      pad.Left + borderInset,
		Y:      pad.Top + borderInset,
		Width:  region.Width - pad.Horizontal() - borderInset*2,
		Height: region.Height - pad.Vertical() - borderInset*2,
	}
	if inner.Width < 0 {
		inner.Width = 0
	}
	if inner.Height < 0 {
		inner.Height = 0
	}

	visibleRows := inner.Height

	// Adjust scroll to keep selection visible.
	if l.selected < l.scrollOff {
		l.scrollOff = l.selected
	}
	if l.selected >= l.scrollOff+visibleRows {
		l.scrollOff = l.selected - visibleRows + 1
	}

	// cursorWidth is 1 when a cursor glyph is configured, 0 when hidden.
	cursorWidth := 1
	if l.cursor == "" {
		cursorWidth = 0
	}

	// Cursor style: accent when focused, dimmed when not.
	cursorStyle := latte.Style{FG: latte.ColorBrightCyan, Bold: true}
	if !l.IsFocused() {
		cursorStyle = latte.Style{FG: latte.ColorBrightBlack}
	}

	for y := 0; y < visibleRows; y++ {
		idx := l.scrollOff + y
		if idx >= len(l.items) {
			break
		}
		itemStyle := style
		cursorRune := " "
		if idx == l.selected {
			if l.IsFocused() && l.highlight {
				itemStyle = l.selectedStyle
			}
			if l.cursor != "" {
				cursorRune = l.cursor
			}
		}

		// Draw cursor glyph (or space) then the label.
		if cursorWidth > 0 {
			sub.DrawText(inner.X, inner.Y+y, cursorRune, cursorStyle)
		}
		text := " " + l.items[idx].Label
		sub.DrawText(inner.X+cursorWidth, inner.Y+y, text, itemStyle)
		// Pad remainder of row to fill background colour.
		used := cursorWidth + len([]rune(text))
		for x := used; x < inner.Width; x++ {
			sub.SetCell(inner.X+x, inner.Y+y, ' ', itemStyle)
		}
	}
}

func (l *List) HandleKey(ev *oat.KeyEvent) bool {
	switch ev.Key() {
	case tcell.KeyUp:
		if l.selected > 0 {
			l.moveCursor(l.selected - 1)
		}
		return true
	case tcell.KeyDown:
		if l.selected < len(l.items)-1 {
			l.moveCursor(l.selected + 1)
		}
		return true
	case tcell.KeyHome, tcell.KeyCtrlA:
		l.moveCursor(0)
		return true
	case tcell.KeyEnd, tcell.KeyCtrlE:
		l.moveCursor(len(l.items) - 1)
		return true
	case tcell.KeyEnter:
		if l.onSelect != nil && l.selected < len(l.items) {
			l.onSelect(l.selected, l.items[l.selected])
		}
		return true
	case tcell.KeyDelete:
		if l.onDelete != nil && l.selected < len(l.items) {
			l.onDelete(l.selected, l.items[l.selected])
		}
		return true
	}
	return false
}

// moveCursor updates the selected index and fires onCursorChange if registered.
func (l *List) moveCursor(idx int) {
	if idx < 0 {
		idx = 0
	}
	if idx >= len(l.items) {
		idx = len(l.items) - 1
	}
	l.selected = idx
	if l.onCursorChange != nil && idx >= 0 && idx < len(l.items) {
		l.onCursorChange(idx, l.items[idx])
	}
}

func (l *List) KeyBindings() []oat.KeyBinding {
	bindings := []oat.KeyBinding{
		{Key: tcell.KeyUp, Label: "↑", Description: "Up"},
		{Key: tcell.KeyDown, Label: "↓", Description: "Down"},
		{Key: tcell.KeyEnter, Label: "Enter", Description: "Select"},
	}
	if l.onDelete != nil {
		bindings = append(bindings,
			oat.KeyBinding{Key: tcell.KeyDelete, Label: "Del", Description: "Delete"},
		)
	}
	return bindings
}
