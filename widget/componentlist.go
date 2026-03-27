package widget

import (
	oat "github.com/antoniocali/oat-latte"
	"github.com/antoniocali/oat-latte/latte"
	"github.com/gdamore/tcell/v2"
)

// ComponentListItem is a single entry in a ComponentList.
// Component is the widget rendered for this row; Value is an opaque identifier
// the caller can use to correlate a row with application data (e.g. a record ID).
type ComponentListItem struct {
	Component oat.Component
	Value     interface{}
}

// ComponentList is a vertically scrollable, keyboard-navigable list whose rows
// are arbitrary Components rather than plain label strings. This allows each
// row to contain rich layouts such as HBox(Text, Flex(Text), Text).
//
// Row heights are variable: each row's Component is measured during Measure /
// Render to determine how many terminal rows it occupies. Scroll is tracked
// by item index, not by pixel offset, so the list always shows full rows.
//
// Default keybindings (when focused):
//   - ↑ / ↓    Move — navigate up and down
//   - Enter     Select — invoke the onSelect callback
//   - Del       Delete — invoke the onDelete callback (if set)
//   - Home / ^A Top — jump to first item
//   - End  / ^E Bottom — jump to last item
//
// ComponentList implements oat.Layout so the theme tree-walker and the
// focus collector recurse into each row's component automatically.
type ComponentList struct {
	oat.BaseComponent
	oat.FocusBehavior

	items          []ComponentListItem
	selected       int // currently highlighted index
	scrollOff      int // first visible item index
	onSelect       func(index int, item ComponentListItem)
	onDelete       func(index int, item ComponentListItem)
	onCursorChange func(index int, item ComponentListItem)
	selectedStyle  latte.Style

	// highlight controls whether the selected row's background is filled with
	// selectedStyle. Defaults to true.
	highlight bool

	// cursor is the rune drawn in the gutter next to the selected row.
	// Defaults to ">". Set to "" to hide it entirely.
	cursor string

	// rowHeight caches the measured height of each row from the last Measure
	// call. Re-populated on every Measure so it is always current.
	rowHeight []int
}

// NewComponentList creates a ComponentList with the given items.
func NewComponentList(items []ComponentListItem) *ComponentList {
	l := &ComponentList{
		items:     items,
		highlight: true,
		cursor:    ">",
	}
	l.EnsureID()
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

// WithStyle sets the display style for this ComponentList.
func (l *ComponentList) WithStyle(s latte.Style) *ComponentList { l.Style = s; return l }

// WithID sets a user-defined identifier on this component.
func (l *ComponentList) WithID(id string) *ComponentList { l.ID = id; return l }

// WithHighlight controls whether the selected row is filled with the
// selected-item style. Defaults to true.
func (l *ComponentList) WithHighlight(enabled bool) *ComponentList {
	l.highlight = enabled
	return l
}

// WithCursor sets the gutter character drawn next to the selected row.
// The default is ">". Pass "" to hide the cursor entirely.
func (l *ComponentList) WithCursor(cursor string) *ComponentList { l.cursor = cursor; return l }

// WithSelectedStyle overrides the highlight style for the selected row.
func (l *ComponentList) WithSelectedStyle(s latte.Style) *ComponentList {
	l.selectedStyle = s
	return l
}

// WithOnSelect registers a callback invoked when the user presses Enter on a row.
func (l *ComponentList) WithOnSelect(fn func(int, ComponentListItem)) *ComponentList {
	l.onSelect = fn
	return l
}

// WithOnCursorChange registers a callback invoked whenever the highlighted index
// changes (Up, Down, Home, End). Suitable for live-preview panels.
func (l *ComponentList) WithOnCursorChange(fn func(int, ComponentListItem)) *ComponentList {
	l.onCursorChange = fn
	return l
}

// WithOnDelete registers a callback invoked when the user presses Delete on a row.
func (l *ComponentList) WithOnDelete(fn func(int, ComponentListItem)) *ComponentList {
	l.onDelete = fn
	return l
}

// SetItems replaces the list contents.
func (l *ComponentList) SetItems(items []ComponentListItem) {
	l.items = items
	if l.selected >= len(items) {
		l.selected = len(items) - 1
	}
	if l.selected < 0 {
		l.selected = 0
	}
	l.rowHeight = nil // invalidate cache
}

// SelectedIndex returns the currently highlighted index.
func (l *ComponentList) SelectedIndex() int { return l.selected }

// SelectedItem returns the currently highlighted item, or zero value + false if empty.
func (l *ComponentList) SelectedItem() (ComponentListItem, bool) {
	if l.selected < 0 || l.selected >= len(l.items) {
		return ComponentListItem{}, false
	}
	return l.items[l.selected], true
}

// GetValue implements oat.ValueGetter. Returns the Value field of the currently
// selected ComponentListItem, or nil if the list is empty.
func (l *ComponentList) GetValue() interface{} {
	item, ok := l.SelectedItem()
	if !ok {
		return nil
	}
	return item.Value
}

// --- oat.Layout ------------------------------------------------------------

// Children returns a flat slice of all row components so the framework's
// tree walkers (theme propagation, focus collection) recurse into them.
func (l *ComponentList) Children() []oat.Component {
	children := make([]oat.Component, len(l.items))
	for i, item := range l.items {
		children[i] = item.Component
	}
	return children
}

// AddChild appends a new item whose Component is the supplied child.
// Value is left nil; use SetItems when you need a non-nil Value.
func (l *ComponentList) AddChild(child oat.Component) {
	l.items = append(l.items, ComponentListItem{Component: child})
	l.rowHeight = nil
}

// --- oat.Scrollable --------------------------------------------------------

func (l *ComponentList) ScrollOffset() int  { return l.scrollOff }
func (l *ComponentList) ContentHeight() int { return len(l.items) }
func (l *ComponentList) ScrollTo(off int) {
	if off < 0 {
		off = 0
	}
	if off >= len(l.items) {
		off = len(l.items) - 1
	}
	l.scrollOff = off
}

// --- oat.Component ---------------------------------------------------------

// Measure asks every row component for its preferred height (unconstrained on
// the Y axis, bounded by MaxWidth on the X axis) and sums them to return the
// total desired size. A cache of per-row heights is stored so Render can reuse
// them without re-measuring.
func (l *ComponentList) Measure(c oat.Constraint) oat.Size {
	style := l.EffectiveStyle(l.IsFocused())

	borderInset := 0
	if style.Border != latte.BorderNone && style.Border != latte.BorderExplicitNone {
		borderInset = 1
	}
	pad := toOatInsets(style.Padding)

	cursorWidth := 1
	if l.cursor == "" {
		cursorWidth = 0
	}

	// Width available for each row's content component.
	innerW := c.MaxWidth - pad.Horizontal() - borderInset*2 - cursorWidth
	if innerW < 0 {
		innerW = 0
	}

	rowC := oat.Constraint{MaxWidth: innerW, MaxHeight: -1}

	l.rowHeight = make([]int, len(l.items))
	totalH := 0
	for i, item := range l.items {
		if item.Component == nil {
			l.rowHeight[i] = 1
		} else {
			s := item.Component.Measure(rowC)
			if s.Height < 1 {
				s.Height = 1
			}
			l.rowHeight[i] = s.Height
		}
		totalH += l.rowHeight[i]
	}

	// Add border and padding overhead.
	totalH += pad.Vertical() + borderInset*2

	if c.MaxHeight >= 0 && totalH > c.MaxHeight {
		totalH = c.MaxHeight
	}
	w := c.MaxWidth
	if w < 0 {
		w = 20
	}
	return oat.Size{Width: w, Height: totalH}
}

// Render draws the visible rows into buf within the given region.
// Each row's Component is rendered into a sub-region sized to its measured
// height. If the row is selected and highlight is enabled, the row's
// background is filled with selectedStyle before delegating to the component.
func (l *ComponentList) Render(buf *oat.Buffer, region oat.Region) {
	style := l.EffectiveStyle(l.IsFocused())
	sub := buf.Sub(region)
	sub.FillBG(style)

	borderInset := 0
	if style.Border != latte.BorderNone && style.Border != latte.BorderExplicitNone {
		sub.DrawBorderTitle(style.Border, l.Title, latte.Style{}, style, oat.AnchorLeft)
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

	cursorWidth := 1
	if l.cursor == "" {
		cursorWidth = 0
	}

	// Cursor styles: accent when focused, dimmed when not.
	cursorStyle := latte.Style{FG: latte.ColorBrightCyan, Bold: true}
	if !l.IsFocused() {
		cursorStyle = latte.Style{FG: latte.ColorBrightBlack}
	}

	// Ensure the row-height cache is populated. If Measure was not called with
	// the current inner width (e.g. after a terminal resize), re-measure now.
	rowContentW := inner.Width - cursorWidth
	if rowContentW < 0 {
		rowContentW = 0
	}
	rowC := oat.Constraint{MaxWidth: rowContentW, MaxHeight: -1}
	if len(l.rowHeight) != len(l.items) {
		l.rowHeight = make([]int, len(l.items))
		for i, item := range l.items {
			if item.Component == nil {
				l.rowHeight[i] = 1
			} else {
				s := item.Component.Measure(rowC)
				if s.Height < 1 {
					s.Height = 1
				}
				l.rowHeight[i] = s.Height
			}
		}
	}

	visibleHeight := inner.Height

	// Scroll invariant: the selected item must be fully visible.
	// Step 1 — scroll up if selection is above the viewport.
	if l.selected < l.scrollOff {
		l.scrollOff = l.selected
	}
	// Step 2 — scroll down until the selection fits within visibleHeight.
	for {
		h := 0
		for i := l.scrollOff; i <= l.selected && i < len(l.items); i++ {
			h += l.rowHeight[i]
		}
		if h <= visibleHeight || l.scrollOff == l.selected {
			break
		}
		l.scrollOff++
	}

	// Render rows starting at scrollOff until we exhaust the visible area.
	curY := inner.Y
	for idx := l.scrollOff; idx < len(l.items); idx++ {
		rowH := l.rowHeight[idx]
		if curY+rowH > inner.Y+visibleHeight {
			// Skip rows that would overflow; always render complete rows only.
			break
		}

		isSelected := idx == l.selected
		rowStyle := style
		if isSelected && l.IsFocused() && l.highlight {
			rowStyle = l.selectedStyle
		}

		// Fill the row background across the full inner width.
		for dy := 0; dy < rowH; dy++ {
			for dx := 0; dx < inner.Width; dx++ {
				sub.SetCell(inner.X+dx, curY+dy, ' ', rowStyle)
			}
		}

		// Draw the cursor glyph (or blank) in the gutter column.
		if cursorWidth > 0 {
			cursorRune := " "
			if isSelected && l.cursor != "" {
				cursorRune = l.cursor
			}
			sub.DrawText(inner.X, curY, cursorRune, cursorStyle)
		}

		// Render the row's component into the region to the right of the gutter.
		if item := l.items[idx]; item.Component != nil {
			rowRegion := oat.Region{
				X:      inner.X + cursorWidth,
				Y:      curY,
				Width:  inner.Width - cursorWidth,
				Height: rowH,
			}
			// Re-measure with the exact allocated width so flex children adapt.
			item.Component.Measure(oat.Constraint{
				MaxWidth:  rowRegion.Width,
				MaxHeight: rowH,
			})
			item.Component.Render(sub, rowRegion)
		}

		curY += rowH
	}
}

// ApplyTheme applies theme tokens to the ComponentList.
// The theme acts as the base; any style fields already set on the widget
// take precedence.
func (l *ComponentList) ApplyTheme(t latte.Theme) {
	l.Style = t.Text.Merge(l.Style)
	l.FocusStyle = latte.Style{BorderFG: t.FocusBorder}.Merge(l.FocusStyle)
	l.selectedStyle = t.ListSelected.Merge(l.selectedStyle)
}

func (l *ComponentList) HandleKey(ev *oat.KeyEvent) bool {
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
func (l *ComponentList) moveCursor(idx int) {
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

func (l *ComponentList) KeyBindings() []oat.KeyBinding {
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
