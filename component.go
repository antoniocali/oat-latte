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
