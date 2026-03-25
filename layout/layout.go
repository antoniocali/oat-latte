// Package layout provides container components for the oat-latte TUI framework.
// Layouts position and size their children — they never render content themselves.
package layout

import (
	oat "github.com/antoniocali/oat-latte"
	"github.com/antoniocali/oat-latte/latte"
)

// childSlot pairs a Component with its computed flex weight for the box layouts.
type childSlot struct {
	component oat.Component
	flex      int // 0 = fixed/auto; >0 = flex weight for remaining space
}

// ---- VBox -----------------------------------------------------------------

// VBox stacks children vertically, top to bottom.
// Children with flex > 0 share the remaining vertical space proportionally.
type VBox struct {
	oat.BaseComponent
	slots []childSlot
	gap   int // rows of empty space between children
}

// NewVBox creates a VBox with optional children.
func NewVBox(children ...oat.Component) *VBox {
	v := &VBox{}
	for _, c := range children {
		v.slots = append(v.slots, childSlot{component: c, flex: 0})
	}
	return v
}

// WithStyle sets the style for this VBox.
func (v *VBox) WithStyle(s latte.Style) *VBox { v.Style = s; return v }

// WithGap sets the number of empty rows between children.
func (v *VBox) WithGap(n int) *VBox { v.gap = n; return v }

// ApplyTheme is a no-op for VBox — it carries no semantic role in the theme.
// The canvas tree-walk will recurse into children automatically.
func (v *VBox) ApplyTheme(_ latte.Theme) {}

// AddChild appends a child with flex weight 0 (fixed/auto size).
// If the child implements FlexSpacer (e.g. VFill), it is automatically
// promoted to a flex slot using the spacer's declared weight.
func (v *VBox) AddChild(c oat.Component) {
	if fs, ok := c.(FlexSpacer); ok {
		v.slots = append(v.slots, childSlot{c, fs.FlexWeight()})
		return
	}
	v.slots = append(v.slots, childSlot{c, 0})
}

// AddFlexChild appends a child that participates in flex space distribution.
func (v *VBox) AddFlexChild(c oat.Component, flex int) {
	if flex < 1 {
		flex = 1
	}
	v.slots = append(v.slots, childSlot{c, flex})
}

// Children satisfies oat.Layout.
func (v *VBox) Children() []oat.Component {
	out := make([]oat.Component, len(v.slots))
	for i, s := range v.slots {
		out[i] = s.component
	}
	return out
}

// Measure computes the VBox's desired size.
// Fixed/auto children contribute their measured height; flex children
// contribute zero (their height is only known at Render time).
func (v *VBox) Measure(c oat.Constraint) oat.Size {
	padded := c.Shrink(toOatInsets(v.Style.Padding))
	totalH := 0
	maxW := 0
	gaps := (len(v.slots) - 1) * v.gap
	if gaps < 0 {
		gaps = 0
	}
	totalH += gaps

	for _, slot := range v.slots {
		if slot.flex > 0 {
			// Flex children claim remaining space at render time; skip their height here.
			s := slot.component.Measure(oat.Constraint{MaxWidth: padded.MaxWidth, MaxHeight: -1})
			if s.Width > maxW {
				maxW = s.Width
			}
			continue
		}
		s := slot.component.Measure(padded)
		if s.Width > maxW {
			maxW = s.Width
		}
		totalH += s.Height
	}

	pad := v.Style.Padding
	return oat.Size{
		Width:  clamp(maxW+pad.Left+pad.Right, 0, c.MaxWidth),
		Height: clamp(totalH+pad.Top+pad.Bottom, 0, c.MaxHeight),
	}
}

// Render draws all children top-to-bottom within region.
func (v *VBox) Render(buf *oat.Buffer, region oat.Region) {
	inner := region.Inner(toOatInsets(v.Style.Padding))
	sub := buf.Sub(inner)

	// First pass: measure all fixed children to find remaining space for flex.
	totalFixed := 0
	totalFlex := 0
	fixedHeights := make([]int, len(v.slots))
	for i, slot := range v.slots {
		if slot.flex == 0 {
			s := slot.component.Measure(inner.ToConstraint())
			fixedHeights[i] = s.Height
			totalFixed += s.Height
		} else {
			totalFlex += slot.flex
		}
	}

	gaps := (len(v.slots) - 1) * v.gap
	if gaps < 0 {
		gaps = 0
	}
	remaining := inner.Height - totalFixed - gaps
	if remaining < 0 {
		remaining = 0
	}

	// Second pass: assign heights and render.
	y := 0
	for i, slot := range v.slots {
		h := fixedHeights[i]
		if slot.flex > 0 && totalFlex > 0 {
			h = (remaining * slot.flex) / totalFlex
			// Respect a VFill's optional max-size cap.
			if vf, ok := slot.component.(*VFill); ok && vf.maxSize > 0 && h > vf.maxSize {
				h = vf.maxSize
			}
		}
		childRegion := oat.Region{X: 0, Y: y, Width: inner.Width, Height: h}
		slot.component.Render(sub, childRegion)
		y += h
		if i < len(v.slots)-1 {
			y += v.gap
		}
	}
}

// ---- HBox -----------------------------------------------------------------

// HBox places children side by side horizontally, left to right.
// Children with flex > 0 share the remaining horizontal space proportionally.
type HBox struct {
	oat.BaseComponent
	slots []childSlot
	gap   int // columns of empty space between children
}

// NewHBox creates an HBox with optional children.
func NewHBox(children ...oat.Component) *HBox {
	h := &HBox{}
	for _, c := range children {
		h.slots = append(h.slots, childSlot{component: c, flex: 0})
	}
	return h
}

// WithStyle sets the style for this HBox.
func (h *HBox) WithStyle(s latte.Style) *HBox { h.Style = s; return h }

// WithGap sets the number of empty columns between children.
func (h *HBox) WithGap(n int) *HBox { h.gap = n; return h }

// ApplyTheme is a no-op for HBox — it carries no semantic role in the theme.
// The canvas tree-walk will recurse into children automatically.
func (h *HBox) ApplyTheme(_ latte.Theme) {}

// AddChild appends a child with flex weight 0.
// If the child implements FlexSpacer (e.g. HFill), it is automatically
// promoted to a flex slot using the spacer's declared weight.
func (h *HBox) AddChild(c oat.Component) {
	if fs, ok := c.(FlexSpacer); ok {
		h.slots = append(h.slots, childSlot{c, fs.FlexWeight()})
		return
	}
	h.slots = append(h.slots, childSlot{c, 0})
}

// AddFlexChild appends a child that participates in flex space distribution.
func (h *HBox) AddFlexChild(c oat.Component, flex int) {
	if flex < 1 {
		flex = 1
	}
	h.slots = append(h.slots, childSlot{c, flex})
}

// Children satisfies oat.Layout.
func (h *HBox) Children() []oat.Component {
	out := make([]oat.Component, len(h.slots))
	for i, s := range h.slots {
		out[i] = s.component
	}
	return out
}

// Measure computes the HBox's desired size.
// Fixed/auto children contribute their measured width; flex children contribute zero.
func (h *HBox) Measure(c oat.Constraint) oat.Size {
	padded := c.Shrink(toOatInsets(h.Style.Padding))
	totalW := 0
	maxH := 0
	gaps := (len(h.slots) - 1) * h.gap
	if gaps < 0 {
		gaps = 0
	}
	totalW += gaps

	for _, slot := range h.slots {
		if slot.flex > 0 {
			s := slot.component.Measure(oat.Constraint{MaxWidth: -1, MaxHeight: padded.MaxHeight})
			if s.Height > maxH {
				maxH = s.Height
			}
			continue
		}
		s := slot.component.Measure(padded)
		if s.Height > maxH {
			maxH = s.Height
		}
		totalW += s.Width
	}

	pad := h.Style.Padding
	return oat.Size{
		Width:  clamp(totalW+pad.Left+pad.Right, 0, c.MaxWidth),
		Height: clamp(maxH+pad.Top+pad.Bottom, 0, c.MaxHeight),
	}
}

// Render draws all children left-to-right within region.
func (h *HBox) Render(buf *oat.Buffer, region oat.Region) {
	inner := region.Inner(toOatInsets(h.Style.Padding))
	sub := buf.Sub(inner)

	// First pass: measure fixed children.
	totalFixed := 0
	totalFlex := 0
	fixedWidths := make([]int, len(h.slots))
	for i, slot := range h.slots {
		if slot.flex == 0 {
			s := slot.component.Measure(inner.ToConstraint())
			fixedWidths[i] = s.Width
			totalFixed += s.Width
		} else {
			totalFlex += slot.flex
		}
	}

	gaps := (len(h.slots) - 1) * h.gap
	if gaps < 0 {
		gaps = 0
	}
	remaining := inner.Width - totalFixed - gaps
	if remaining < 0 {
		remaining = 0
	}

	// Second pass: assign widths and render.
	x := 0
	for i, slot := range h.slots {
		w := fixedWidths[i]
		if slot.flex > 0 && totalFlex > 0 {
			w = (remaining * slot.flex) / totalFlex
			// Respect an HFill's optional max-size cap.
			if hf, ok := slot.component.(*HFill); ok && hf.maxSize > 0 && w > hf.maxSize {
				w = hf.maxSize
			}
		}
		childRegion := oat.Region{X: x, Y: 0, Width: w, Height: inner.Height}
		slot.component.Render(sub, childRegion)
		x += w
		if i < len(h.slots)-1 {
			x += h.gap
		}
	}
}

// ---- Grid -----------------------------------------------------------------

// GridChild wraps a component with its grid position and optional span.
type GridChild struct {
	Component oat.Component
	Row, Col  int
	RowSpan   int // default 1
	ColSpan   int // default 1
}

// Grid arranges children in a fixed rows×cols grid.
// Each cell has equal width and height (terminal grids rarely need unequal cells).
type Grid struct {
	oat.BaseComponent
	rows     int
	cols     int
	children []GridChild
	rowGap   int
	colGap   int
}

// NewGrid creates a Grid with the given number of rows and columns.
func NewGrid(rows, cols int) *Grid {
	return &Grid{rows: rows, cols: cols}
}

// WithStyle sets the style.
func (g *Grid) WithStyle(s latte.Style) *Grid { g.Style = s; return g }

// WithGap sets the gap between rows and columns.
func (g *Grid) WithGap(rowGap, colGap int) *Grid { g.rowGap = rowGap; g.colGap = colGap; return g }

// ApplyTheme is a no-op for Grid — it carries no semantic role in the theme.
func (g *Grid) ApplyTheme(_ latte.Theme) {}

// Add places a component at (row, col) with span (1,1).
func (g *Grid) Add(row, col int, c oat.Component) *Grid {
	g.children = append(g.children, GridChild{Component: c, Row: row, Col: col, RowSpan: 1, ColSpan: 1})
	return g
}

// AddSpan places a component at (row, col) with the given row and column span.
func (g *Grid) AddSpan(row, col, rowSpan, colSpan int, c oat.Component) *Grid {
	g.children = append(g.children, GridChild{Component: c, Row: row, Col: col, RowSpan: rowSpan, ColSpan: colSpan})
	return g
}

// AddChild satisfies oat.Layout (appends at next available cell).
func (g *Grid) AddChild(c oat.Component) {
	pos := len(g.children)
	row := pos / g.cols
	col := pos % g.cols
	g.Add(row, col, c)
}

// Children satisfies oat.Layout.
func (g *Grid) Children() []oat.Component {
	out := make([]oat.Component, len(g.children))
	for i, gc := range g.children {
		out[i] = gc.Component
	}
	return out
}

// Measure returns the size needed to fit all rows and columns.
func (g *Grid) Measure(c oat.Constraint) oat.Size {
	pad := toOatInsets(g.Style.Padding)
	inner := c.Shrink(pad)

	cellW := 0
	cellH := 0
	if g.cols > 0 {
		cellW = (inner.MaxWidth - g.colGap*(g.cols-1)) / g.cols
	}
	if g.rows > 0 {
		cellH = (inner.MaxHeight - g.rowGap*(g.rows-1)) / g.rows
	}

	totalW := cellW*g.cols + g.colGap*(g.cols-1) + g.Style.Padding.Left + g.Style.Padding.Right
	totalH := cellH*g.rows + g.rowGap*(g.rows-1) + g.Style.Padding.Top + g.Style.Padding.Bottom
	return oat.Size{Width: clamp(totalW, 0, c.MaxWidth), Height: clamp(totalH, 0, c.MaxHeight)}
}

// Render draws each child in its assigned cell.
func (g *Grid) Render(buf *oat.Buffer, region oat.Region) {
	inner := region.Inner(toOatInsets(g.Style.Padding))
	sub := buf.Sub(inner)

	if g.cols == 0 || g.rows == 0 {
		return
	}

	cellW := (inner.Width - g.colGap*(g.cols-1)) / g.cols
	cellH := (inner.Height - g.rowGap*(g.rows-1)) / g.rows

	for _, gc := range g.children {
		if gc.Row >= g.rows || gc.Col >= g.cols {
			continue
		}
		rs := gc.RowSpan
		if rs < 1 {
			rs = 1
		}
		cs := gc.ColSpan
		if cs < 1 {
			cs = 1
		}

		x := gc.Col * (cellW + g.colGap)
		y := gc.Row * (cellH + g.rowGap)
		w := cellW*cs + g.colGap*(cs-1)
		h := cellH*rs + g.rowGap*(rs-1)

		childRegion := oat.Region{X: x, Y: y, Width: w, Height: h}
		gc.Component.Render(sub, childRegion)
	}
}

// ---- Stack ----------------------------------------------------------------

// Stack layers children on top of each other (Z-axis).
// The last child renders on top. Used for modals and overlays.
type Stack struct {
	oat.BaseComponent
	children []oat.Component
}

// NewStack creates a Stack with optional children.
func NewStack(children ...oat.Component) *Stack {
	return &Stack{children: children}
}

// WithStyle sets the style.
func (s *Stack) WithStyle(st latte.Style) *Stack { s.Style = st; return s }

// ApplyTheme is a no-op for Stack — it carries no semantic role in the theme.
func (s *Stack) ApplyTheme(_ latte.Theme) {}

// AddChild appends a child layer.
func (s *Stack) AddChild(c oat.Component) { s.children = append(s.children, c) }

// Children satisfies oat.Layout.
func (s *Stack) Children() []oat.Component { return s.children }

// Measure returns the size of the largest child.
func (s *Stack) Measure(c oat.Constraint) oat.Size {
	maxW, maxH := 0, 0
	for _, child := range s.children {
		sz := child.Measure(c)
		if sz.Width > maxW {
			maxW = sz.Width
		}
		if sz.Height > maxH {
			maxH = sz.Height
		}
	}
	return oat.Size{Width: maxW, Height: maxH}
}

// Render draws each child in the full region (later children paint over earlier ones).
func (s *Stack) Render(buf *oat.Buffer, region oat.Region) {
	for _, child := range s.children {
		child.Render(buf, region)
	}
}

// ---- Border ---------------------------------------------------------------

// Border wraps a single child and draws a border around it.
// It reserves one cell on each side for the border lines.
// An optional Title is stamped into the top border line: ╭─ Title ──╮
type Border struct {
	oat.BaseComponent
	child       oat.Component
	titleStyle  latte.Style // style for the title text in the border line
	titleAnchor oat.Anchor  // horizontal position of the title (default AnchorLeft)
}

// NewBorder wraps child with a border using the default (single) border style.
// Use the builder methods (WithStyle, WithTitle, WithTitleStyle) to customise.
func NewBorder(child oat.Component) *Border {
	b := &Border{child: child}
	b.Style.Border = latte.BorderSingle
	return b
}

// WithStyle sets a custom style on this Border.
// If style.Border is BorderNone the border type defaults to BorderSingle.
func (b *Border) WithStyle(s latte.Style) *Border {
	b.Style = s
	if b.Style.Border == latte.BorderNone {
		b.Style.Border = latte.BorderSingle
	}
	return b
}

// WithTitle sets the label stamped into the top border rule.
// anchor is optional and defaults to oat.AnchorLeft when omitted.
// Pass oat.AnchorCenter or oat.AnchorRight to reposition the title.
func (b *Border) WithTitle(title string, anchor ...oat.Anchor) *Border {
	b.Title = title
	if len(anchor) > 0 {
		b.titleAnchor = anchor[0]
	}
	return b
}

// WithTitleStyle sets a custom style for the title text (e.g. bold cyan).
// If not set, the border colour is inherited.
func (b *Border) WithTitleStyle(s latte.Style) *Border {
	b.titleStyle = s
	return b
}

// WithRoundedCorner controls whether the border uses rounded corners (╭╮╰╯)
// instead of the default square ones (┌┐└┘).
//
// Calling WithRoundedCorner(true) switches the border style to BorderRounded.
// Calling WithRoundedCorner(false) on a rounded border restores BorderSingle.
//
// Panics if rounded is true and the current border style is BorderDouble,
// BorderThick, or BorderDashed: Unicode provides arc corner codepoints only
// for light-weight strokes (─ │), so they cannot connect to double (═ ║),
// heavy (━ ┃), or dashed (╌ ╎) lines without producing a broken visual.
// Use WithStyle(latte.Style{Border: latte.BorderRounded}) to switch style
// entirely, or keep the incompatible border style without rounded corners.
func (b *Border) WithRoundedCorner(rounded bool) *Border {
	if rounded {
		switch b.Style.Border {
		case latte.BorderDouble:
			panic("oat-latte: WithRoundedCorner(true) is not supported for BorderDouble — " +
				"Unicode has no double-stroke arc corners (╔╗╚╝ cannot be rounded)")
		case latte.BorderThick:
			panic("oat-latte: WithRoundedCorner(true) is not supported for BorderThick — " +
				"Unicode has no heavy-stroke arc corners (┏┓┗┛ cannot be rounded)")
		case latte.BorderDashed:
			panic("oat-latte: WithRoundedCorner(true) is not supported for BorderDashed — " +
				"arc corners (╭╮╰╯) do not connect to dashed strokes (╌ ╎)")
		}
		b.Style.Border = latte.BorderRounded
	} else if b.Style.Border == latte.BorderRounded {
		b.Style.Border = latte.BorderSingle
	}
	return b
}

// ApplyTheme applies Panel and PanelTitle tokens from the theme to this Border.
// The FocusBorder colour is stored in FocusStyle.BorderFG so the focus-aware
// Render logic picks it up automatically.
func (b *Border) ApplyTheme(t latte.Theme) {
	// Preserve the border shape if the caller set it explicitly, but apply
	// colours from the theme.
	if b.Style.Border != latte.BorderNone {
		b.Style.BorderFG = t.Panel.BorderFG
		b.Style.BG = t.Panel.BG
		b.Style.BorderBG = t.Panel.BG
	} else {
		b.Style = t.Panel
		b.Style.BorderBG = t.Panel.BG
	}
	b.titleStyle = t.PanelTitle
	b.FocusStyle = latte.Style{BorderFG: t.FocusBorder}
}

// AddChild sets the single child (replaces any existing child).
func (b *Border) AddChild(c oat.Component) { b.child = c }

// Children satisfies oat.Layout.
func (b *Border) Children() []oat.Component {
	if b.child == nil {
		return nil
	}
	return []oat.Component{b.child}
}

// Measure adds 2 to each dimension for the border.
func (b *Border) Measure(c oat.Constraint) oat.Size {
	if b.child == nil {
		return oat.Size{}
	}
	inner := c.Shrink(oat.Insets{Top: 1, Right: 1, Bottom: 1, Left: 1})
	s := b.child.Measure(inner)
	return oat.Size{Width: s.Width + 2, Height: s.Height + 2}
}

// Render draws the border (with optional title) and then the child inside it.
// When any descendant component has focus the border colour is promoted to
// the focused border colour (cyan) so the user always knows which panel is active.
func (b *Border) Render(buf *oat.Buffer, region oat.Region) {
	style := b.Style
	if b.child != nil && containsFocus(b.child) {
		// Override the border colour to the focus highlight colour.
		// If a custom FocusStyle was set on the Border itself, use its BorderFG;
		// otherwise fall back to the framework default (bright cyan).
		fg := b.FocusStyle.BorderFG
		if fg == latte.ColorDefault {
			fg = latte.ColorBrightCyan
		}
		style.BorderFG = fg
	}
	sub := buf.Sub(region)
	// Fill the full region with the panel background before drawing the border
	// runes and child. This ensures the interior padding area and any cells not
	// covered by the child have the correct background colour.
	if style.BG != latte.ColorDefault {
		sub.FillBG(style)
	}
	sub.DrawBorderTitle(style.Border, b.Title, b.titleStyle, style, b.titleAnchor)
	if b.child != nil {
		innerRegion := oat.Region{X: 1, Y: 1, Width: region.Width - 2, Height: region.Height - 2}
		b.child.Render(sub, innerRegion)
	}
}

// ---- Padding --------------------------------------------------------------

// Padding wraps a single child and adds configurable space around it.
type Padding struct {
	oat.BaseComponent
	child  oat.Component
	insets oat.Insets
}

// NewPadding wraps child with the given insets.
func NewPadding(child oat.Component, insets oat.Insets) *Padding {
	return &Padding{child: child, insets: insets}
}

// NewPaddingUniform wraps child with uniform padding on all sides.
func NewPaddingUniform(child oat.Component, n int) *Padding {
	return NewPadding(child, oat.Uniform(n))
}

// ApplyTheme is a no-op for Padding — it carries no semantic role in the theme.
func (p *Padding) ApplyTheme(_ latte.Theme) {}

// AddChild sets the single child.
func (p *Padding) AddChild(c oat.Component) { p.child = c }

// Children satisfies oat.Layout.
func (p *Padding) Children() []oat.Component {
	if p.child == nil {
		return nil
	}
	return []oat.Component{p.child}
}

// Measure adds padding to the child's measured size.
func (p *Padding) Measure(c oat.Constraint) oat.Size {
	if p.child == nil {
		return oat.Size{Width: p.insets.Horizontal(), Height: p.insets.Vertical()}
	}
	inner := c.Shrink(p.insets)
	s := p.child.Measure(inner)
	return oat.Size{
		Width:  s.Width + p.insets.Horizontal(),
		Height: s.Height + p.insets.Vertical(),
	}
}

// Render draws the child offset by the padding insets.
func (p *Padding) Render(buf *oat.Buffer, region oat.Region) {
	if p.child == nil {
		return
	}
	innerRegion := region.Inner(p.insets)
	sub := buf.Sub(region)
	p.child.Render(sub, innerRegion)
}

// ---- helpers --------------------------------------------------------------

func clamp(v, min, max int) int {
	if max < 0 {
		return v // unconstrained
	}
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// toOatInsets converts a latte.Insets to an oat.Insets.
// The two types are structurally identical; this avoids a circular import.
func toOatInsets(i latte.Insets) oat.Insets {
	return oat.Insets{Top: i.Top, Right: i.Right, Bottom: i.Bottom, Left: i.Left}
}

// containsFocus reports whether c or any of its descendants is focused.
// Used by Border so it can highlight its border when any child is active.
func containsFocus(c oat.Component) bool {
	if f, ok := c.(oat.Focusable); ok && f.IsFocused() {
		return true
	}
	if l, ok := c.(oat.Layout); ok {
		for _, child := range l.Children() {
			if containsFocus(child) {
				return true
			}
		}
	}
	return false
}

// ---- FlexSpacer -----------------------------------------------------------

// FlexSpacer is an optional interface for components that should be treated as
// flex children when added to a VBox or HBox via AddChild.
// VFill and HFill implement this interface so callers can use the simpler
// AddChild API without calling AddFlexChild manually.
type FlexSpacer interface {
	oat.Component
	// FlexWeight returns the flex weight for this spacer (always ≥ 1).
	FlexWeight() int
}

// ---- VFill ----------------------------------------------------------------

// VFill is a vertical spacer that expands to fill remaining space in a VBox.
// When added to a VBox via AddChild or AddFlexChild, it consumes all remaining
// vertical space and pushes its siblings toward the edges.
//
// Inside a scrollable container VFill degrades gracefully: use WithMaxSize to
// cap the height to a safe bound, or it will simply take zero height when
// remaining space is negative.
//
//	vbox.AddChild(topWidget)
//	vbox.AddChild(layout.NewVFill())   // expands to fill gap
//	vbox.AddChild(bottomWidget)
type VFill struct {
	weight  int
	maxSize int // 0 = uncapped
}

// NewVFill creates a VFill spacer with flex weight 1.
func NewVFill() *VFill { return &VFill{weight: 1} }

// WithWeight sets the flex weight (default 1). Higher weights claim
// proportionally more space when multiple flex children compete.
func (f *VFill) WithWeight(w int) *VFill {
	if w < 1 {
		w = 1
	}
	f.weight = w
	return f
}

// WithMaxSize caps the maximum height this spacer will consume.
// This is recommended when VFill is used inside a Scrollable container,
// where uncapped spacers can produce unexpectedly large content heights.
func (f *VFill) WithMaxSize(n int) *VFill { f.maxSize = n; return f }

// FlexWeight satisfies FlexSpacer.
func (f *VFill) FlexWeight() int { return f.weight }

// Measure returns zero — VFill claims space only at Render time via the flex
// pass in VBox.
func (f *VFill) Measure(_ oat.Constraint) oat.Size { return oat.Size{} }

// Render does nothing; VFill is a pure spacer.
func (f *VFill) Render(_ *oat.Buffer, _ oat.Region) {}

// ---- HFill ----------------------------------------------------------------

// HFill is a horizontal spacer that expands to fill remaining space in an HBox.
// When added to an HBox via AddChild or AddFlexChild, it consumes all remaining
// horizontal space and pushes its siblings toward the edges.
//
//	hbox.AddChild(leftWidget)
//	hbox.AddChild(layout.NewHFill())   // expands to fill gap
//	hbox.AddChild(rightWidget)
type HFill struct {
	weight  int
	maxSize int // 0 = uncapped
}

// NewHFill creates an HFill spacer with flex weight 1.
func NewHFill() *HFill { return &HFill{weight: 1} }

// WithWeight sets the flex weight (default 1).
func (f *HFill) WithWeight(w int) *HFill {
	if w < 1 {
		w = 1
	}
	f.weight = w
	return f
}

// WithMaxSize caps the maximum width this spacer will consume.
func (f *HFill) WithMaxSize(n int) *HFill { f.maxSize = n; return f }

// FlexWeight satisfies FlexSpacer.
func (f *HFill) FlexWeight() int { return f.weight }

// Measure returns zero — HFill claims space only at Render time via the flex
// pass in HBox.
func (f *HFill) Measure(_ oat.Constraint) oat.Size { return oat.Size{} }

// Render does nothing; HFill is a pure spacer.
func (f *HFill) Render(_ *oat.Buffer, _ oat.Region) {}
