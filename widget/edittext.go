package widget

import (
	"strings"
	"unicode/utf8"

	oat "github.com/antoniocali/oat-latte"
	"github.com/antoniocali/oat-latte/latte"
	"github.com/gdamore/tcell/v2"
)

// EditText is a single-line or multi-line text input component.
// It is Focusable, Keybinder, and optionally Scrollable.
//
// Default keybindings:
//   - ^S  Save    — calls OnSave with the current text (if set)
//   - ^G  Cancel  — calls OnCancel (if set)
//   - ^K  Kill    — delete from cursor to end of line
//   - ^U  Clear   — delete from start of line to cursor
//   - ^A  Start   — move cursor to start of line
//   - ^E  End     — move cursor to end of line
//
// A hint label (set via WithHint) is rendered as a persistent muted line
// directly above the editable content.  Unlike a placeholder — which only
// appears when the field is empty — the hint is always visible and serves as
// an inline field label, removing the need for a separate Text widget.
type EditText struct {
	oat.BaseComponent
	oat.FocusBehavior

	lines       []string // internal line buffer
	cursorRow   int
	cursorCol   int // rune-counted column in the current line
	scrollRow   int // first visible row (for multi-line scrolling)
	multiLine   bool
	placeholder string
	hint        string // persistent muted label rendered above the content
	hintStyle   latte.Style
	maxLength   int // 0 = unlimited (single-line only)

	onChange func(text string)
	onSave   func(text string)
	onCancel func()

	// callerStyle and callerFocusStyle preserve the styles set by the caller
	// (via WithStyle) before any theme application. ApplyTheme always merges
	// the current theme token with these originals so that switching themes
	// fully replaces the previous theme's colours rather than accumulating
	// stale values from the prior theme.
	callerStyle      latte.Style
	callerFocusStyle latte.Style
}

// NewEditText creates a single-line text input.
func NewEditText() *EditText {
	e := &EditText{lines: []string{""}}
	e.EnsureID()
	// FocusStyle is intentionally NOT pre-seeded here.
	// ApplyTheme sets it from the active theme's InputFocus token.
	// Pre-seeding with hardcoded values would cause them to survive theme
	// switches via Merge, blocking the new theme's colours from taking effect.
	return e
}

// NewMultiLineEditText creates a multi-line text editor.
func NewMultiLineEditText() *EditText {
	e := NewEditText()
	e.multiLine = true
	return e
}

// WithStyle sets the display style for this EditText.
func (e *EditText) WithStyle(s latte.Style) *EditText { e.Style = s; e.callerStyle = s; return e }

// WithID sets a user-defined identifier on this component.
// Use this to retrieve the widget's value later via Canvas.GetValue(id).
func (e *EditText) WithID(id string) *EditText { e.ID = id; return e }

// GetValue implements oat.ValueGetter. Returns the current text as a string.
func (e *EditText) GetValue() interface{} { return e.GetText() }

// WithHAlign sets the horizontal alignment for this widget within a VBox slot.
// No argument (or HAlignFill) resets to the default fill behaviour.
func (e *EditText) WithHAlign(a ...oat.HAlign) *EditText {
	e.BaseComponent.HAlign = oat.HAlignFill
	if len(a) > 0 {
		e.BaseComponent.HAlign = a[0]
	}
	return e
}

// WithVAlign sets the vertical alignment for this widget within an HBox slot.
// No argument (or VAlignFill) resets to the default fill behaviour.
func (e *EditText) WithVAlign(a ...oat.VAlign) *EditText {
	e.BaseComponent.VAlign = oat.VAlignFill
	if len(a) > 0 {
		e.BaseComponent.VAlign = a[0]
	}
	return e
}

// WithPlaceholder sets placeholder text shown when the input is empty.
func (e *EditText) WithPlaceholder(p string) *EditText { e.placeholder = p; return e }

// WithHint sets a persistent muted label rendered as a single line directly
// above the editable content area.  The hint is always visible regardless of
// whether the field has content, making it suitable as an inline field label.
// It removes the need for a companion Text widget when using borderless inputs.
//
//	input := widget.NewEditText(latte.Style{Border: latte.BorderExplicitNone}).
//	    WithHint("Title")
func (e *EditText) WithHint(hint string) *EditText { e.hint = hint; return e }

// WithMaxLength sets the maximum number of characters (single-line only).
func (e *EditText) WithMaxLength(n int) *EditText { e.maxLength = n; return e }

// WithOnChange registers a callback invoked whenever the content changes.
func (e *EditText) WithOnChange(fn func(string)) *EditText { e.onChange = fn; return e }

// WithOnSave registers a callback invoked when the user presses ^S.
// The callback receives the current text. It is also advertised in the StatusBar
// as "^S Save".
func (e *EditText) WithOnSave(fn func(string)) *EditText { e.onSave = fn; return e }

// WithOnCancel registers a callback invoked when the user presses ^G.
// Useful for discarding edits and returning to a previous state.
func (e *EditText) WithOnCancel(fn func()) *EditText { e.onCancel = fn; return e }

// ApplyTheme applies theme tokens to the EditText.
// The theme acts as the base; any style fields already set on the widget
// (e.g. Border: BorderExplicitNone passed via WithStyle) take precedence.
// ApplyTheme always re-derives Style from the theme token merged with the
// original callerStyle so that switching themes fully replaces the previous
// theme's colours rather than accumulating stale values.
func (e *EditText) ApplyTheme(t latte.Theme) {
	e.Style = t.Input.Merge(e.callerStyle)
	e.FocusStyle = t.InputFocus.Merge(e.callerFocusStyle)
	e.hintStyle = t.Muted
}

// SetText replaces all content.
func (e *EditText) SetText(text string) {
	if e.multiLine {
		e.lines = strings.Split(text, "\n")
	} else {
		e.lines = []string{strings.ReplaceAll(text, "\n", " ")}
	}
	e.cursorRow = len(e.lines) - 1
	e.cursorCol = utf8.RuneCountInString(e.lines[e.cursorRow])
}

// GetText returns the current content as a string.
func (e *EditText) GetText() string {
	return strings.Join(e.lines, "\n")
}

// --- oat.Scrollable --------------------------------------------------------
func (e *EditText) ScrollOffset() int  { return e.scrollRow }
func (e *EditText) ContentHeight() int { return len(e.lines) }
func (e *EditText) ScrollTo(off int) {
	if off < 0 {
		off = 0
	}
	if off >= len(e.lines) {
		off = len(e.lines) - 1
	}
	e.scrollRow = off
}

// --- oat.Focusable ---------------------------------------------------------

func (e *EditText) HandleKey(ev *oat.KeyEvent) bool {
	switch ev.Key() {
	case tcell.KeyRune:
		e.insertRune(ev.Rune())
		return true
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		e.deleteBackward()
		return true
	case tcell.KeyDelete:
		e.deleteForward()
		return true
	case tcell.KeyLeft:
		e.moveCursorLeft()
		return true
	case tcell.KeyRight:
		e.moveCursorRight()
		return true
	case tcell.KeyHome, tcell.KeyCtrlA:
		e.cursorCol = 0
		return true
	case tcell.KeyEnd, tcell.KeyCtrlE:
		e.cursorCol = utf8.RuneCountInString(e.lines[e.cursorRow])
		return true
	case tcell.KeyUp:
		if e.multiLine && e.cursorRow > 0 {
			e.cursorRow--
			lineLen := utf8.RuneCountInString(e.lines[e.cursorRow])
			if e.cursorCol > lineLen {
				e.cursorCol = lineLen
			}
		}
		return e.multiLine
	case tcell.KeyDown:
		if e.multiLine && e.cursorRow < len(e.lines)-1 {
			e.cursorRow++
			lineLen := utf8.RuneCountInString(e.lines[e.cursorRow])
			if e.cursorCol > lineLen {
				e.cursorCol = lineLen
			}
		}
		return e.multiLine
	case tcell.KeyEnter:
		if e.multiLine {
			e.splitLine()
			return true
		}
		// Single-line: Enter triggers Save if registered, otherwise let Canvas handle it.
		if e.onSave != nil {
			e.onSave(e.GetText())
			return true
		}
		return false
	case tcell.KeyCtrlS:
		if e.onSave != nil {
			e.onSave(e.GetText())
		}
		return true
	case tcell.KeyCtrlG:
		if e.onCancel != nil {
			e.onCancel()
		}
		return true
	case tcell.KeyCtrlK:
		// Kill to end of line.
		runes := []rune(e.lines[e.cursorRow])
		e.lines[e.cursorRow] = string(runes[:e.cursorCol])
		e.notifyChange()
		return true
	case tcell.KeyCtrlU:
		// Kill from start of line to cursor.
		runes := []rune(e.lines[e.cursorRow])
		e.lines[e.cursorRow] = string(runes[e.cursorCol:])
		e.cursorCol = 0
		e.notifyChange()
		return true
	}
	return false
}

// KeyBindings advertises available shortcuts to the StatusBar.
func (e *EditText) KeyBindings() []oat.KeyBinding {
	bindings := []oat.KeyBinding{
		{Key: tcell.KeyCtrlA, Label: "^A", Description: "Start"},
		{Key: tcell.KeyCtrlE, Label: "^E", Description: "End"},
		{Key: tcell.KeyCtrlK, Label: "^K", Description: "Kill line"},
		{Key: tcell.KeyCtrlU, Label: "^U", Description: "Clear"},
	}
	if e.onSave != nil {
		bindings = append([]oat.KeyBinding{
			{Key: tcell.KeyCtrlS, Label: "^S", Description: "Save"},
		}, bindings...)
	}
	if e.onCancel != nil {
		bindings = append(bindings,
			oat.KeyBinding{Key: tcell.KeyCtrlG, Label: "^G", Description: "Cancel"},
		)
	}
	return bindings
}

// --- oat.Component ---------------------------------------------------------

func (e *EditText) Measure(c oat.Constraint) oat.Size {
	h := 1
	if e.multiLine {
		h = len(e.lines)
	}
	if e.hint != "" {
		h++ // one extra row for the hint label
	}
	// Account for border (1 cell each side) and padding.
	borderInset := 0
	if e.Style.Border != latte.BorderNone && e.Style.Border != latte.BorderExplicitNone {
		borderInset = 1
	}
	pad := toOatInsets(e.Style.Padding)
	h += pad.Vertical() + borderInset*2
	if c.MaxHeight >= 0 && h > c.MaxHeight {
		h = c.MaxHeight
	}
	w := c.MaxWidth
	if w < 0 {
		w = 20 // sensible default when unconstrained
	}
	return oat.Size{Width: w, Height: h}
}

func (e *EditText) Render(buf *oat.Buffer, region oat.Region) {
	style := e.EffectiveStyle(e.IsFocused())
	sub := buf.Sub(region)
	sub.FillBG(style)

	// Border inset: when a border is drawn it occupies the outermost cell on
	// all four sides. Content must be offset by 1 to stay inside the border.
	borderInset := 0
	if style.Border != latte.BorderNone && style.Border != latte.BorderExplicitNone {
		sub.DrawBorderTitle(style.Border, e.Title, latte.Style{}, style, oat.AnchorLeft)
		borderInset = 1
	}

	// Total content offset = border inset + style padding.
	pad := toOatInsets(style.Padding)
	offX := pad.Left + borderInset
	offY := pad.Top + borderInset
	contentW := region.Width - pad.Horizontal() - borderInset*2
	contentH := region.Height - pad.Vertical() - borderInset*2
	if contentW < 0 {
		contentW = 0
	}
	if contentH < 0 {
		contentH = 0
	}

	// Render hint label above the content area.
	if e.hint != "" {
		hs := e.hintStyle
		if hs == (latte.Style{}) {
			hs = style
			hs.FG = latte.ColorBrightBlack
		}
		sub.DrawText(offX, offY, e.hint, hs)
		offY++
		contentH--
	}

	visibleRows := contentH
	startRow := e.scrollRow

	// Adjust scroll to keep cursor visible.
	if e.cursorRow < startRow {
		startRow = e.cursorRow
		e.scrollRow = startRow
	}
	if e.cursorRow >= startRow+visibleRows {
		startRow = e.cursorRow - visibleRows + 1
		e.scrollRow = startRow
	}

	for y := 0; y < visibleRows; y++ {
		rowIdx := startRow + y
		if rowIdx >= len(e.lines) {
			break
		}
		line := e.lines[rowIdx]

		if line == "" && rowIdx == 0 && e.placeholder != "" {
			phStyle := style
			phStyle.FG = latte.ColorBrightBlack
			sub.DrawText(offX, offY+y, e.placeholder, phStyle)
			continue
		}

		// Horizontal scroll for single-line inputs.
		runes := []rune(line)
		visibleStart := 0
		if !e.multiLine && e.cursorCol >= contentW {
			visibleStart = e.cursorCol - contentW + 1
		}
		if visibleStart > len(runes) {
			visibleStart = len(runes)
		}
		visibleRunes := runes[visibleStart:]
		if len(visibleRunes) > contentW {
			visibleRunes = visibleRunes[:contentW]
		}
		sub.DrawText(offX, offY+y, string(visibleRunes), style)
	}

	// Show cursor when focused.
	if e.IsFocused() {
		curX := e.cursorCol + offX
		if !e.multiLine && e.cursorCol >= contentW {
			curX = contentW - 1 + offX
		}
		curY := (e.cursorRow - e.scrollRow) + offY
		if curX < region.Width && curY < region.Height {
			sub.ShowCursor(curX, curY)
		}
	}
}

// --- internal editing ops --------------------------------------------------

func (e *EditText) insertRune(r rune) {
	line := []rune(e.lines[e.cursorRow])
	if e.maxLength > 0 && len(line) >= e.maxLength {
		return
	}
	newLine := make([]rune, 0, len(line)+1)
	newLine = append(newLine, line[:e.cursorCol]...)
	newLine = append(newLine, r)
	newLine = append(newLine, line[e.cursorCol:]...)
	e.lines[e.cursorRow] = string(newLine)
	e.cursorCol++
	e.notifyChange()
}

func (e *EditText) deleteBackward() {
	if e.cursorCol > 0 {
		line := []rune(e.lines[e.cursorRow])
		newLine := append(line[:e.cursorCol-1], line[e.cursorCol:]...)
		e.lines[e.cursorRow] = string(newLine)
		e.cursorCol--
		e.notifyChange()
	} else if e.multiLine && e.cursorRow > 0 {
		// Merge with previous line.
		prev := e.lines[e.cursorRow-1]
		cur := e.lines[e.cursorRow]
		e.cursorCol = utf8.RuneCountInString(prev)
		e.lines[e.cursorRow-1] = prev + cur
		e.lines = append(e.lines[:e.cursorRow], e.lines[e.cursorRow+1:]...)
		e.cursorRow--
		e.notifyChange()
	}
}

func (e *EditText) deleteForward() {
	line := []rune(e.lines[e.cursorRow])
	if e.cursorCol < len(line) {
		newLine := append(line[:e.cursorCol], line[e.cursorCol+1:]...)
		e.lines[e.cursorRow] = string(newLine)
		e.notifyChange()
	} else if e.multiLine && e.cursorRow < len(e.lines)-1 {
		e.lines[e.cursorRow] += e.lines[e.cursorRow+1]
		e.lines = append(e.lines[:e.cursorRow+1], e.lines[e.cursorRow+2:]...)
		e.notifyChange()
	}
}

func (e *EditText) moveCursorLeft() {
	if e.cursorCol > 0 {
		e.cursorCol--
	} else if e.multiLine && e.cursorRow > 0 {
		e.cursorRow--
		e.cursorCol = utf8.RuneCountInString(e.lines[e.cursorRow])
	}
}

func (e *EditText) moveCursorRight() {
	lineLen := utf8.RuneCountInString(e.lines[e.cursorRow])
	if e.cursorCol < lineLen {
		e.cursorCol++
	} else if e.multiLine && e.cursorRow < len(e.lines)-1 {
		e.cursorRow++
		e.cursorCol = 0
	}
}

func (e *EditText) splitLine() {
	line := []rune(e.lines[e.cursorRow])
	before := string(line[:e.cursorCol])
	after := string(line[e.cursorCol:])
	e.lines[e.cursorRow] = before
	newLines := make([]string, 0, len(e.lines)+1)
	newLines = append(newLines, e.lines[:e.cursorRow+1]...)
	newLines = append(newLines, after)
	newLines = append(newLines, e.lines[e.cursorRow+1:]...)
	e.lines = newLines
	e.cursorRow++
	e.cursorCol = 0
	e.notifyChange()
}

func (e *EditText) notifyChange() {
	if e.onChange != nil {
		e.onChange(e.GetText())
	}
}
