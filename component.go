package oat

import (
	"fmt"
	"sync/atomic"

	"github.com/antoniocali/oat-latte/latte"
	"github.com/gdamore/tcell/v2"
)

// idCounter is a package-level counter used to generate unique component IDs.
// It is incremented atomically so construction is safe across goroutines.
var idCounter atomic.Int64

// genID returns a new unique component identifier of the form "oat-<n>".
func genID() string { return fmt.Sprintf("oat-%d", idCounter.Add(1)) }

// KeyEvent wraps a tcell.EventKey for use in component handlers.
type KeyEvent = tcell.EventKey

// MouseEvent wraps a tcell.EventMouse for use in component handlers.
type MouseEvent = tcell.EventMouse

// Component is the fundamental building block of an oat-latte UI.
// Every widget, layout, and container implements Component.
//
// The render pipeline has two passes:
//  1. Measure — the parent calls Measure to ask the child how large it wants to be.
//  2. Render  — the parent calls Render with the final allocated Region.
type Component interface {
	// Measure returns the desired Size given the available Constraint.
	// Children must never exceed constraint.MaxWidth / MaxHeight.
	// A Constraint value of -1 means unconstrained on that axis.
	Measure(c Constraint) Size

	// Render draws the component into buf within the given Region.
	// The Region is guaranteed to be >= the Size returned by Measure
	// (within the parent's discretion).
	Render(buf *Buffer, region Region)
}

// Focusable is an opt-in interface for components that can receive keyboard focus.
// The FocusManager automatically discovers all Focusable nodes in the tree.
type Focusable interface {
	Component

	// SetFocused is called by the FocusManager when focus is gained or lost.
	SetFocused(focused bool)

	// IsFocused returns the current focus state.
	IsFocused() bool

	// HandleKey is called when this component has focus and a key event fires.
	// Return true if the event was consumed (stops bubbling).
	HandleKey(ev *KeyEvent) bool
}

// Keybinder is an opt-in interface for components that want to advertise
// their available keyboard shortcuts to the StatusBar footer.
type Keybinder interface {
	KeyBindings() []KeyBinding
}

// KeyBinding describes a single keyboard shortcut and its effect.
type KeyBinding struct {
	// Key is the tcell key constant (e.g. tcell.KeyCtrlS).
	// Set to tcell.KeyRune and use Rune for printable characters.
	Key  tcell.Key
	Rune rune // used when Key == tcell.KeyRune
	Mod  tcell.ModMask

	// Label is the short key hint shown in the StatusBar (e.g. "^S").
	// Keep this brief — one to four characters.
	Label string

	// Description is the human-readable action name shown next to Label
	// in the StatusBar (e.g. "Save"). If empty, Label is shown alone.
	Description string

	// Handler is called when the key is pressed while this component is focused.
	// Bindings with a nil Handler are display-only hints for the StatusBar.
	Handler func()
}

// Scrollable is an opt-in interface for components with internal scroll state.
type Scrollable interface {
	Component

	// ScrollOffset returns the current scroll offset in lines (for vertical scroll).
	ScrollOffset() int

	// ScrollTo sets the scroll offset.
	ScrollTo(offset int)

	// ContentHeight returns the total height of the scrollable content.
	ContentHeight() int
}

// Layout is a Component that contains child Components.
// It is responsible for measuring and positioning its children.
type Layout interface {
	Component

	// Children returns the direct children of this layout.
	Children() []Component

	// AddChild appends a child to this layout.
	AddChild(child Component)
}

// --- FocusBehavior is a convenience embed for components that are Focusable ---

// FocusBehavior provides a default implementation of the Focusable boilerplate.
// Embed this in any component to gain focus tracking for free.
type FocusBehavior struct {
	focused bool
}

func (f *FocusBehavior) SetFocused(focused bool) { f.focused = focused }
func (f *FocusBehavior) IsFocused() bool         { return f.focused }

// --- BaseComponent is a convenience embed for the Style ---

// BaseComponent holds common fields shared by all concrete components.
type BaseComponent struct {
	ID         string // unique identifier; auto-generated if not set via WithID
	Style      latte.Style
	FocusStyle latte.Style // applied when focused (merged over Style)
	Title      string      // optional title rendered above the component
	HAlign     HAlign      // horizontal alignment within a VBox slot (default HAlignFill)
	VAlign     VAlign      // vertical alignment within an HBox slot (default VAlignFill)
}

// initID ensures the component has an ID, generating one if not already set.
// Called by widget constructors to guarantee every component has a non-empty ID.
func (b *BaseComponent) initID() {
	if b.ID == "" {
		b.ID = genID()
	}
}

// EnsureID is the exported version of initID for use by sub-packages.
// Widget constructors outside the oat package call this to guarantee every
// component has a non-empty ID without exposing the genID implementation.
func (b *BaseComponent) EnsureID() {
	b.initID()
}

// EffectiveStyle merges FocusStyle on top of Style when focused.
func (b *BaseComponent) EffectiveStyle(focused bool) latte.Style {
	if focused {
		return b.Style.Merge(b.FocusStyle)
	}
	return b.Style
}

// SetFocusStyle sets the focus style, but only if no per-component focus style
// has been explicitly configured (i.e. FocusStyle is still the zero value).
// This allows Canvas to inject a global focus style without overriding
// per-component customisation.
func (b *BaseComponent) SetFocusStyle(s latte.Style) {
	if b.FocusStyle == (latte.Style{}) {
		b.FocusStyle = s
	}
}

// FocusStyleInjector is implemented by any component that embeds BaseComponent.
// Canvas uses this to inject a global focus highlight style at startup.
type FocusStyleInjector interface {
	SetFocusStyle(latte.Style)
}

// ValueGetter is an opt-in interface for components that hold a user-visible
// value (text content, checked state, selected index, etc.).
// Canvas.GetValue(id) walks the tree and returns the value from the first
// component whose ID matches.
//
// Return types by widget:
//   - EditText    → string  (current text)
//   - CheckBox    → bool    (checked state)
//   - List        → interface{} (Value field of selected ListItem, or nil if empty)
//   - Button      → string  (button label)
//   - ProgressBar → float64 (current value 0.0–1.0)
//   - Text/Title  → string  (displayed text)
type ValueGetter interface {
	GetValue() interface{}
}

// IDer is implemented by any component that embeds BaseComponent.
// Canvas uses this to find a component by its ID.
type IDer interface {
	GetID() string
}

// GetID returns the component's unique identifier.
func (b *BaseComponent) GetID() string { return b.ID }

// GetHAlign returns the component's horizontal alignment preference.
// Satisfies AlignProvider. The zero value is HAlignFill (full-width, unchanged behaviour).
func (b *BaseComponent) GetHAlign() HAlign { return b.HAlign }

// GetVAlign returns the component's vertical alignment preference.
// Satisfies AlignProvider. The zero value is VAlignFill (full-height, unchanged behaviour).
func (b *BaseComponent) GetVAlign() VAlign { return b.VAlign }

// WithHAlign sets the horizontal alignment for this component within a VBox slot.
// Variadic so WithHAlign() with no arguments is valid (resets to HAlignFill).
func (b *BaseComponent) WithHAlign(a ...HAlign) *BaseComponent {
	if len(a) > 0 {
		b.HAlign = a[0]
	} else {
		b.HAlign = HAlignFill
	}
	return b
}

// WithVAlign sets the vertical alignment for this component within an HBox slot.
// Variadic so WithVAlign() with no arguments is valid (resets to VAlignFill).
func (b *BaseComponent) WithVAlign(a ...VAlign) *BaseComponent {
	if len(a) > 0 {
		b.VAlign = a[0]
	} else {
		b.VAlign = VAlignFill
	}
	return b
}

// Anchor controls the horizontal position of an element within its container.
//
// It is the H-axis anchor type. All APIs that place content along the
// horizontal axis accept Anchor:
//   - DrawBorderTitle / Border.WithTitle — title position inside the top rule
//   - ProgressBar.WithPercentage — where the "XX%" label appears
//   - Divider.WithAnchor (AxisVertical) — which section of a vertical divider to render
//
// For the vertical axis see VAnchor.
// For positioning a widget within its allocated cross-axis region see HAlign / VAlign.
type Anchor int

const (
	// AnchorLeft aligns the element to the left edge. This is the default.
	AnchorLeft Anchor = iota
	// AnchorCenter centres the element horizontally.
	AnchorCenter
	// AnchorRight aligns the element to the right edge.
	AnchorRight
)

// VAnchor controls the vertical position of an element within its container.
//
// It is the V-axis counterpart of Anchor. APIs that place content along the
// vertical axis accept VAnchor:
//   - Divider.WithAnchor (AxisHorizontal) — which section of a horizontal divider to render
//
// VAnchor is intentionally V-axis only so that the cross-axis alignment types
// (HAlign / VAlign) can be built without mixing placement semantics.
type VAnchor int

const (
	// VAnchorTop aligns the element to the top edge. This is the default.
	VAnchorTop VAnchor = iota
	// VAnchorMiddle centres the element vertically.
	VAnchorMiddle
	// VAnchorBottom aligns the element to the bottom edge.
	VAnchorBottom
)

// HAlign controls how a widget is positioned along the horizontal (cross) axis
// within its allocated slot in a VBox (or any container that distributes space
// vertically and needs to decide how wide each child is).
//
// The zero value HAlignFill preserves the existing behaviour: the child is
// given the full allocated width, exactly as before this type was added.
//
// Non-fill values cause the container to measure the child's desired width and
// then position it within the slot according to the chosen alignment, leaving
// the remaining width empty.
type HAlign int

const (
	// HAlignFill gives the child the full allocated width. This is the default
	// (zero value) and matches the pre-alignment behaviour of VBox.
	HAlignFill HAlign = iota
	// HAlignLeft shrinks the child to its desired width and pins it to the left.
	HAlignLeft
	// HAlignCenter shrinks the child to its desired width and centres it.
	HAlignCenter
	// HAlignRight shrinks the child to its desired width and pins it to the right.
	HAlignRight
)

// VAlign controls how a widget is positioned along the vertical (cross) axis
// within its allocated slot in an HBox (or any container that distributes space
// horizontally and needs to decide how tall each child is).
//
// The zero value VAlignFill preserves the existing behaviour: the child is
// given the full allocated height, exactly as before this type was added.
//
// Non-fill values cause the container to measure the child's desired height and
// then position it within the slot according to the chosen alignment, leaving
// the remaining height empty.
type VAlign int

const (
	// VAlignFill gives the child the full allocated height. This is the default
	// (zero value) and matches the pre-alignment behaviour of HBox.
	VAlignFill VAlign = iota
	// VAlignTop shrinks the child to its desired height and pins it to the top.
	VAlignTop
	// VAlignMiddle shrinks the child to its desired height and centres it.
	VAlignMiddle
	// VAlignBottom shrinks the child to its desired height and pins it to the bottom.
	VAlignBottom
)

// AlignProvider is an opt-in interface for components that carry their own
// cross-axis alignment preferences. Any component embedding BaseComponent
// satisfies this interface automatically via the GetHAlign / GetVAlign getters.
//
// VBox inspects each child's HAlign (via AlignProvider or its own default).
// HBox inspects each child's VAlign (via AlignProvider or its own default).
type AlignProvider interface {
	GetHAlign() HAlign
	GetVAlign() VAlign
}

// FocusGuard is an opt-in interface that lets a component dynamically opt out
// of receiving keyboard focus. When a component implements FocusGuard and
// IsFocusable returns false, the DFS focus walk skips that node entirely —
// neither it nor any of its descendants will be added to the focus list.
//
// This is useful for context-sensitive panels where entire subtrees should
// be unreachable via Tab cycling depending on application state.
type FocusGuard interface {
	IsFocusable() bool
}

// ThemeReceiver is an opt-in interface for components that can be styled by a
// latte.Theme. Canvas.WithTheme walks the entire component tree and calls
// ApplyTheme on every node that implements this interface.
//
// Widgets apply the relevant semantic tokens to their Style / FocusStyle fields.
// Layout containers apply the theme to themselves and then propagate it to their
// children, so a single WithTheme call at the canvas level is sufficient to
// theme the whole application.
type ThemeReceiver interface {
	ApplyTheme(t latte.Theme)
}
