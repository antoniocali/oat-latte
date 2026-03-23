package oat

import "github.com/gdamore/tcell/v2"

// FocusManager maintains the ordered list of all Focusable components
// in the component tree (DFS order) and tracks which one is active.
//
// It is created by Canvas during Run() via a tree walk, and rebuilt on
// every component-tree mutation.
type FocusManager struct {
	nodes   []Focusable
	current int // index into nodes; -1 means no focusable components
}

// NewFocusManager creates an empty FocusManager.
func NewFocusManager() *FocusManager {
	return &FocusManager{current: -1}
}

// Collect walks the component tree rooted at root and populates the
// manager with all Focusable nodes in DFS order.
// Any previously tracked nodes are cleared first.
func (fm *FocusManager) Collect(root Component) {
	// Defocus current before rebuilding
	if fm.current >= 0 && fm.current < len(fm.nodes) {
		fm.nodes[fm.current].SetFocused(false)
	}

	fm.nodes = fm.nodes[:0]
	fm.current = -1
	walkFocusable(root, &fm.nodes)

	if len(fm.nodes) > 0 {
		fm.current = 0
		fm.nodes[0].SetFocused(true)
	}
}

// walkFocusable performs a DFS over the component tree collecting Focusable nodes.
// If a node implements FocusGuard and IsFocusable() returns false, both the node
// and its entire subtree are skipped.
func walkFocusable(c Component, out *[]Focusable) {
	if g, ok := c.(FocusGuard); ok && !g.IsFocusable() {
		return
	}
	if f, ok := c.(Focusable); ok {
		*out = append(*out, f)
	}
	if l, ok := c.(Layout); ok {
		for _, child := range l.Children() {
			walkFocusable(child, out)
		}
	}
}

// Current returns the currently focused component, or nil.
func (fm *FocusManager) Current() Focusable {
	if fm.current < 0 || fm.current >= len(fm.nodes) {
		return nil
	}
	return fm.nodes[fm.current]
}

// Nodes returns all focusable nodes collected from the component tree.
// The slice is in DFS order and should not be mutated by callers.
func (fm *FocusManager) Nodes() []Focusable {
	return fm.nodes
}

// Next moves focus to the next component (Tab behaviour).
func (fm *FocusManager) Next() {
	if len(fm.nodes) == 0 {
		return
	}
	fm.nodes[fm.current].SetFocused(false)
	fm.current = (fm.current + 1) % len(fm.nodes)
	fm.nodes[fm.current].SetFocused(true)
}

// Prev moves focus to the previous component (Shift+Tab behaviour).
func (fm *FocusManager) Prev() {
	if len(fm.nodes) == 0 {
		return
	}
	fm.nodes[fm.current].SetFocused(false)
	fm.current = (fm.current - 1 + len(fm.nodes)) % len(fm.nodes)
	fm.nodes[fm.current].SetFocused(true)
}

// FocusIndex moves focus to the component at a specific index.
func (fm *FocusManager) FocusIndex(i int) {
	if i < 0 || i >= len(fm.nodes) {
		return
	}
	if fm.current >= 0 && fm.current < len(fm.nodes) {
		fm.nodes[fm.current].SetFocused(false)
	}
	fm.current = i
	fm.nodes[i].SetFocused(true)
}

// FocusByRef moves focus to the first node that is pointer-identical to target.
// This is used by Canvas.WithPrimary to honour a designated initial focus component.
func (fm *FocusManager) FocusByRef(target Focusable) {
	for i, n := range fm.nodes {
		if n == target {
			fm.FocusIndex(i)
			return
		}
	}
}

// Dispatch sends a key event to the currently focused component.
// Returns true if the event was consumed.
//
// Dispatch flow:
//  1. Walk the focused component's KeyBindings(). A binding with a non-nil
//     Handler is executed and the event is consumed. Bindings with a nil
//     Handler are display-only labels for the StatusBar and are skipped here.
//  2. If no binding handler fired, delegate to HandleKey for full key handling.
func (fm *FocusManager) Dispatch(ev *KeyEvent) bool {
	focused := fm.Current()
	if focused == nil {
		return false
	}

	// Only invoke bindings that have an explicit handler attached.
	// Bindings without a Handler are StatusBar hints only — the real logic
	// lives in HandleKey and must not be short-circuited here.
	if kb, ok := focused.(Keybinder); ok {
		for _, binding := range kb.KeyBindings() {
			if binding.Handler != nil && matchesBinding(ev, binding) {
				binding.Handler()
				return true
			}
		}
	}

	// Delegate all key handling to the component's HandleKey method.
	return focused.HandleKey(ev)
}

// CurrentKeyBindings returns the KeyBindings of the currently focused component,
// or nil if the focused component doesn't implement Keybinder.
func (fm *FocusManager) CurrentKeyBindings() []KeyBinding {
	focused := fm.Current()
	if focused == nil {
		return nil
	}
	if kb, ok := focused.(Keybinder); ok {
		return kb.KeyBindings()
	}
	return nil
}

// matchesBinding checks whether a KeyEvent corresponds to a KeyBinding.
func matchesBinding(ev *KeyEvent, b KeyBinding) bool {
	if ev.Modifiers() != b.Mod {
		return false
	}
	if b.Key == tcell.KeyRune {
		return ev.Key() == tcell.KeyRune && ev.Rune() == b.Rune
	}
	return ev.Key() == b.Key
}

// --- Convenience constructors for common KeyBindings ----------------------

// Bind creates a KeyBinding for a special key (non-rune).
// label is the short key hint (e.g. "^S"); desc is the action name (e.g. "Save").
func Bind(key tcell.Key, label, desc string, handler func()) KeyBinding {
	return KeyBinding{Key: key, Label: label, Description: desc, Handler: handler}
}

// BindRune creates a KeyBinding for a printable character.
func BindRune(r rune, mod tcell.ModMask, label, desc string, handler func()) KeyBinding {
	return KeyBinding{Key: tcell.KeyRune, Rune: r, Mod: mod, Label: label, Description: desc, Handler: handler}
}
