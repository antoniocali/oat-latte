package oat

// Size represents a width and height in terminal cells.
type Size struct {
	Width  int
	Height int
}

// Region represents a rectangular area on the terminal screen.
// X and Y are the top-left corner (0-indexed).
type Region struct {
	X      int
	Y      int
	Width  int
	Height int
}

// Constraint is passed from parent to child during the Measure pass.
// It describes the maximum space available. A value of -1 means unconstrained.
type Constraint struct {
	MaxWidth  int
	MaxHeight int
}

// Insets represents padding or margin on all four sides.
type Insets struct {
	Top    int
	Right  int
	Bottom int
	Left   int
}

// Horizontal returns the total horizontal inset (Left + Right).
func (i Insets) Horizontal() int { return i.Left + i.Right }

// Vertical returns the total vertical inset (Top + Bottom).
func (i Insets) Vertical() int { return i.Top + i.Bottom }

// Uniform returns an Insets with the same value on all four sides.
func Uniform(n int) Insets {
	return Insets{Top: n, Right: n, Bottom: n, Left: n}
}

// Symmetric returns an Insets with separate horizontal and vertical values.
func Symmetric(vertical, horizontal int) Insets {
	return Insets{Top: vertical, Right: horizontal, Bottom: vertical, Left: horizontal}
}

// Shrink reduces a Constraint by the given insets, ensuring it never goes below 0.
func (c Constraint) Shrink(insets Insets) Constraint {
	w := c.MaxWidth - insets.Horizontal()
	h := c.MaxHeight - insets.Vertical()
	if w < 0 {
		w = 0
	}
	if h < 0 {
		h = 0
	}
	return Constraint{MaxWidth: w, MaxHeight: h}
}

// Clamp ensures a Size fits within the Constraint.
func (c Constraint) Clamp(s Size) Size {
	if c.MaxWidth >= 0 && s.Width > c.MaxWidth {
		s.Width = c.MaxWidth
	}
	if c.MaxHeight >= 0 && s.Height > c.MaxHeight {
		s.Height = c.MaxHeight
	}
	return s
}

// Inner returns the Region inset by the given Insets.
func (r Region) Inner(insets Insets) Region {
	x := r.X + insets.Left
	y := r.Y + insets.Top
	w := r.Width - insets.Horizontal()
	h := r.Height - insets.Vertical()
	if w < 0 {
		w = 0
	}
	if h < 0 {
		h = 0
	}
	return Region{X: x, Y: y, Width: w, Height: h}
}

// ToConstraint converts a Region into a Constraint matching its size.
func (r Region) ToConstraint() Constraint {
	return Constraint{MaxWidth: r.Width, MaxHeight: r.Height}
}

// Dimension describes how a component sizes itself on a single axis.
type Dimension interface{ dimension() }

type (
	// FixedDim requests exactly n cells.
	FixedDim struct{ N int }
	// FillDim requests all available space.
	FillDim struct{}
	// AutoDim requests as much space as the content needs.
	AutoDim struct{}
)

func (FixedDim) dimension() {}
func (FillDim) dimension()  {}
func (AutoDim) dimension()  {}

// Fixed is a convenience constructor for FixedDim.
func Fixed(n int) Dimension { return FixedDim{N: n} }

// Fill is the singleton FillDim value.
var Fill Dimension = FillDim{}

// Auto is the singleton AutoDim value.
var Auto Dimension = AutoDim{}
