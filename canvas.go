package oat

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/antoniocali/oat-latte/latte"
	"github.com/gdamore/tcell/v2"
)

// Canvas is the root container of an oat-latte application.
// It owns the tcell screen, the focus manager, and the render loop.
//
// Usage:
//
//	app := oat.NewCanvas(
//	    oat.WithHeader(myHeader),
//	    oat.WithBody(myLayout),
//	)
//	if err := app.Run(); err != nil {
//	    log.Fatal(err)
//	}
type Canvas struct {
	header     Component
	footer     Component // if nil, an auto StatusBar is injected
	body       Component
	style      latte.Style
	focusStyle latte.Style  // injected into all focusable components that have no per-component FocusStyle
	theme      *latte.Theme // optional global theme applied to the whole component tree
	headerH    int
	footerH    int
	autoBar    statusBarSetter // reference to the injected StatusBar
	focus      *FocusManager
	screen     tcell.Screen
	quit       chan struct{}
	primary    Focusable // if set, receives initial focus instead of DFS-first

	// overlay stack — last element is rendered (and focused) on top.
	// These are modal overlays dismissed by Esc.
	overlays []Component

	// persistentOverlays are rendered above everything (including modal overlays)
	// but are never dismissed by Esc. Use ShowPersistentOverlay for components
	// like NotificationManager that should always remain visible.
	persistentOverlays []Component

	// notifyCh receives time.Time ticks produced by notification timers.
	// Receiving on this channel triggers a re-render so expiring notifications
	// are removed from the screen without requiring a key event.
	notifyCh chan time.Time

	// globalBindings are key bindings that fire regardless of which component
	// currently holds focus.  They are checked after the focused component has
	// had a chance to handle the key, so a focused widget can still shadow a
	// global binding when that is the desired behaviour (e.g. Esc in EditText).
	globalBindings []KeyBinding
}

// statusBarSetter is the interface StatusBar satisfies so Canvas can update it
// without importing the widget package (which would create a cycle).
type statusBarSetter interface {
	Component
	SetBindings([]KeyBinding)
}

// NotificationOverlay is implemented by any component that can be mounted as a
// Canvas notification manager.  Satisfying this interface is the only
// requirement for use with WithNotificationManager.
//
// widget.NotificationManager implements this interface.
type NotificationOverlay interface {
	Component
	// SetNotifyChannel receives the Canvas's internal notify channel so that
	// timer goroutines can trigger re-renders when a timed notification expires.
	// Called once by WithNotificationManager during Canvas construction.
	SetNotifyChannel(ch chan<- time.Time)
}

// CanvasOption configures a Canvas.
type CanvasOption func(*Canvas)

// WithHeader sets a component rendered at the top of the screen.
func WithHeader(c Component) CanvasOption {
	return func(cv *Canvas) { cv.header = c }
}

// WithFooter sets a component rendered at the bottom of the screen.
// If not set, an auto-populated StatusBar is used.
func WithFooter(c Component) CanvasOption {
	return func(cv *Canvas) { cv.footer = c }
}

// WithBody sets the main content component.
func WithBody(c Component) CanvasOption {
	return func(cv *Canvas) { cv.body = c }
}

// WithStyle sets the background style for the Canvas.
func WithStyle(s latte.Style) CanvasOption {
	return func(cv *Canvas) { cv.style = s }
}

// WithAutoStatusBar injects a StatusBar into the footer that is automatically
// populated with keybindings from the currently focused component.
// This is the default behaviour when no footer is provided via WithFooter.
func WithAutoStatusBar(bar statusBarSetter) CanvasOption {
	return func(cv *Canvas) {
		cv.autoBar = bar
		cv.footer = bar
	}
}

// WithPrimary sets the component that receives focus when the application starts.
// If not set, focus lands on the first Focusable node in DFS order.
// Use this to direct the user's attention to the most important interactive element.
func WithPrimary(f Focusable) CanvasOption {
	return func(cv *Canvas) { cv.primary = f }
}

// WithFocusStyle sets a global focus highlight style that is injected into every
// focusable component that has not set its own FocusStyle.
// Per-component FocusStyle always takes precedence over the global style.
func WithFocusStyle(s latte.Style) CanvasOption {
	return func(cv *Canvas) { cv.focusStyle = s }
}

// WithTheme applies a latte.Theme to every component in the tree that implements
// oat.ThemeReceiver. The theme is applied once during Run(), before the first
// render, so components constructed after Run() is called are not affected.
//
// WithTheme also sets the canvas background style from Theme.Canvas and the
// global focus border colour from Theme.FocusBorder.
// Per-component styles set explicitly after WithTheme will not be overridden.
func WithTheme(t latte.Theme) CanvasOption {
	return func(cv *Canvas) { cv.theme = &t }
}

// WithGlobalKeyBinding registers one or more key bindings that are active
// regardless of which component currently holds focus.  The focused widget
// always gets first refusal: global bindings are only checked when the focused
// widget did not consume the key event.  Calling this option multiple times
// (or passing multiple bindings) accumulates rather than replaces bindings.
//
// Example — toggle theme with ^T from anywhere:
//
//	oat.WithGlobalKeyBinding(oat.KeyBinding{
//	    Key:         tcell.KeyCtrlT,
//	    Label:       "^T",
//	    Description: "Toggle theme",
//	    Handler:     func() { app.SetTheme(latte.ThemeDark) },
//	})
func WithGlobalKeyBinding(bindings ...KeyBinding) CanvasOption {
	return func(cv *Canvas) {
		cv.globalBindings = append(cv.globalBindings, bindings...)
	}
}

// WithNotificationManager mounts nm as a persistent overlay and wires it to
// the canvas event loop so expiring notifications trigger automatic re-renders.
// This replaces the previous two-step pattern of calling SetNotifyChannel and
// ShowPersistentOverlay manually.
//
// The caller constructs the manager and retains the reference for pushing
// notifications; the canvas owns the mounting and channel wiring:
//
//	notifs := widget.NewNotificationManager()
//
//	app := oat.NewCanvas(
//	    oat.WithTheme(latte.ThemeDark),
//	    oat.WithBody(body),
//	    oat.WithNotificationManager(notifs),
//	)
//
//	// later, from any callback:
//	notifs.Push("Saved", widget.NotificationKindSuccess, 2*time.Second)
func WithNotificationManager(nm NotificationOverlay) CanvasOption {
	return func(cv *Canvas) {
		nm.SetNotifyChannel(cv.notifyCh)
		cv.persistentOverlays = append(cv.persistentOverlays, nm)
	}
}

// NewCanvas constructs a Canvas from the given options.
func NewCanvas(opts ...CanvasOption) *Canvas {
	cv := &Canvas{
		focus:    NewFocusManager(),
		quit:     make(chan struct{}),
		notifyCh: make(chan time.Time, 8),
	}
	for _, opt := range opts {
		opt(cv)
	}
	return cv
}

// SetHeader replaces the header component after construction.
func (cv *Canvas) SetHeader(c Component) { cv.header = c }

// SetFooter replaces the footer component after construction.
func (cv *Canvas) SetFooter(c Component) { cv.footer = c }

// SetBody replaces the body component after construction.
func (cv *Canvas) SetBody(c Component) { cv.body = c }

// SetAutoBar registers a StatusBar to be auto-populated with focus keybindings.
func (cv *Canvas) SetAutoBar(bar statusBarSetter) {
	cv.autoBar = bar
	cv.footer = bar
}

// Quit signals the event loop to stop gracefully.
func (cv *Canvas) Quit() {
	select {
	case <-cv.quit:
	default:
		close(cv.quit)
	}
}

// ShowDialog pushes a dialog (or any Component) onto the overlay stack and
// redirects keyboard focus into it.  The underlying body and header remain
// visible but are visually dimmed and receive no key events.
//
// Call HideDialog to dismiss the topmost overlay.
func (cv *Canvas) ShowDialog(d Component) {
	cv.overlays = append(cv.overlays, d)
	// Apply theme to the new overlay so it inherits the active colour scheme.
	if cv.theme != nil {
		applyThemeTree(d, *cv.theme)
	}
	// Collect focusable nodes inside the dialog.
	cv.focus.Collect(d)
	if cv.focusStyle != (latte.Style{}) {
		cv.injectFocusStyle()
	}
	cv.updateStatusBar()
}

// HideDialog removes the topmost overlay from the stack and restores focus to
// the body.  If the stack is already empty this is a no-op.
func (cv *Canvas) HideDialog() {
	if len(cv.overlays) == 0 {
		return
	}
	cv.overlays = cv.overlays[:len(cv.overlays)-1]
	// Restore focus to the body tree.
	if cv.body != nil {
		cv.focus.Collect(cv.body)
	}
	if cv.focusStyle != (latte.Style{}) {
		cv.injectFocusStyle()
	}
	if cv.primary != nil && len(cv.overlays) == 0 {
		cv.focus.FocusByRef(cv.primary)
	}
	cv.updateStatusBar()
}

// ShowPersistentOverlay registers a component that is always rendered on top of
// everything (including modal dialogs) but is never dismissed by Esc.
// Use this for components like NotificationManager that must remain visible
// for the lifetime of the application.
func (cv *Canvas) ShowPersistentOverlay(d Component) {
	cv.persistentOverlays = append(cv.persistentOverlays, d)
	if cv.theme != nil {
		applyThemeTree(d, *cv.theme)
	}
}

// HasOverlay reports whether any overlay is currently displayed.
func (cv *Canvas) HasOverlay() bool { return len(cv.overlays) > 0 }

// FocusByRef moves keyboard focus to the first Focusable node in the current
// focus tree that is pointer-identical to target.
// This lets application code direct focus to a specific widget in response to
// a key event (e.g. pressing 'e' to jump directly into the editor's title field).
// It is a no-op if target is not found in the current focus tree.
func (cv *Canvas) FocusByRef(target Focusable) {
	cv.focus.FocusByRef(target)
	cv.updateStatusBar()
}

// Run initialises the terminal screen, enters the event loop, and blocks until
// the application quits (Ctrl+C, Esc, or cv.Quit()).
func (cv *Canvas) Run() error {
	screen, err := tcell.NewScreen()
	if err != nil {
		return err
	}
	if err := screen.Init(); err != nil {
		return err
	}
	defer screen.Fini()

	screen.EnableMouse()
	screen.Clear()
	cv.screen = screen

	// Collect focusable nodes from the component tree.
	if cv.body != nil {
		cv.focus.Collect(cv.body)
	}
	// Apply a global theme to every component in the tree that opts in.
	// This must happen before focus-style injection so the theme's focus
	// colours are in place before the injection guard runs.
	if cv.theme != nil {
		cv.applyTheme()
	}
	// Inject the global focus style into any focusable component that hasn't
	// set its own FocusStyle. Per-component styles are preserved.
	if cv.focusStyle != (latte.Style{}) {
		cv.injectFocusStyle()
	}
	// If a primary component was designated, move initial focus to it.
	if cv.primary != nil {
		cv.focus.FocusByRef(cv.primary)
	}
	cv.updateStatusBar()

	// Handle OS signals for graceful shutdown.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// tcell delivers events on its own goroutine; we poll with PollEvent.
	eventCh := make(chan tcell.Event, 64)
	go func() {
		for {
			ev := screen.PollEvent()
			if ev == nil {
				return
			}
			select {
			case eventCh <- ev:
			case <-cv.quit:
				return
			}
		}
	}()

	// Initial render.
	cv.render()
	screen.Show()

	for {
		select {
		case <-cv.quit:
			return nil
		case <-sigCh:
			return nil
		case <-cv.notifyCh:
			// A notification timer fired — re-render so expired notifications
			// are removed from the display without waiting for a key event.
			cv.render()
			screen.Show()
		case ev := <-eventCh:
			if cv.handleEvent(ev) {
				cv.render()
				screen.Show()
			}
		}
	}
}

// handleEvent processes a tcell.Event. Returns true if a re-render is needed.
func (cv *Canvas) handleEvent(ev tcell.Event) bool {
	switch e := ev.(type) {
	case *tcell.EventResize:
		cv.screen.Sync()
		return true

	case *tcell.EventKey:
		// Global quit bindings.
		if e.Key() == tcell.KeyCtrlC {
			cv.Quit()
			return false
		}

		// Escape: dismiss topmost overlay if one is open; otherwise quit.
		if e.Key() == tcell.KeyEscape {
			if cv.HasOverlay() {
				cv.HideDialog()
				return true
			}
			cv.Quit()
			return false
		}

		// While an overlay is visible, key navigation is confined to it.
		// Tab/Shift-Tab and arrows cycle within the overlay's focus nodes.
		if cv.HasOverlay() {
			if e.Key() == tcell.KeyTab {
				cv.focus.Next()
				cv.updateStatusBar()
				return true
			}
			if e.Key() == tcell.KeyBacktab {
				cv.focus.Prev()
				cv.updateStatusBar()
				return true
			}
			consumed := cv.focus.Dispatch(e)
			if !consumed {
				// Try global bindings before falling back to arrow-key cycling.
				if cv.dispatchGlobal(e) {
					return true
				}
				switch e.Key() {
				case tcell.KeyUp, tcell.KeyLeft:
					cv.focus.Prev()
					cv.updateStatusBar()
				case tcell.KeyDown, tcell.KeyRight:
					cv.focus.Next()
					cv.updateStatusBar()
				}
			}
			return true
		}

		// Normal (no overlay) focus cycling.
		if e.Key() == tcell.KeyTab {
			cv.focus.Next()
			cv.updateStatusBar()
			return true
		}
		if e.Key() == tcell.KeyBacktab {
			cv.focus.Prev()
			cv.updateStatusBar()
			return true
		}

		// Dispatch to focused component first.
		// If the component consumes the key (returns true), we're done.
		// If it does NOT consume an arrow key, fall through to inter-widget
		// navigation so the user can move focus with arrows in addition to Tab.
		consumed := cv.focus.Dispatch(e)
		if !consumed {
			// Try global bindings before falling back to arrow-key cycling.
			if cv.dispatchGlobal(e) {
				return true
			}
			switch e.Key() {
			case tcell.KeyUp, tcell.KeyLeft:
				cv.focus.Prev()
				cv.updateStatusBar()
			case tcell.KeyDown, tcell.KeyRight:
				cv.focus.Next()
				cv.updateStatusBar()
			}
		}
		return true
	}
	return false
}

// render performs a full Measure → Render pass and writes to the screen buffer.
func (cv *Canvas) render() {
	w, h := cv.screen.Size()
	if w == 0 || h == 0 {
		return
	}

	buf := newBuffer(cv.screen)
	buf.Fill(' ', cv.style)

	// Hide cursor before the render tree runs. Any focused EditText will
	// re-show it during its own Render call; this ensures the cursor is
	// hidden when no editable component is focused.
	buf.HideCursor()

	// Calculate region allocations.
	headerH := 0
	footerH := 0

	if cv.header != nil {
		s := cv.header.Measure(Constraint{MaxWidth: w, MaxHeight: h})
		headerH = s.Height
	}
	if cv.footer != nil {
		s := cv.footer.Measure(Constraint{MaxWidth: w, MaxHeight: h})
		footerH = s.Height
	}

	bodyH := h - headerH - footerH
	if bodyH < 0 {
		bodyH = 0
	}

	// Render header.
	if cv.header != nil && headerH > 0 {
		headerRegion := Region{X: 0, Y: 0, Width: w, Height: headerH}
		cv.header.Render(buf, headerRegion)
	}

	// Render body.
	if cv.body != nil && bodyH > 0 {
		bodyRegion := Region{X: 0, Y: headerH, Width: w, Height: bodyH}
		cv.body.Render(buf, bodyRegion)
	}

	// Render footer.
	if cv.footer != nil && footerH > 0 {
		footerRegion := Region{X: 0, Y: h - footerH, Width: w, Height: footerH}
		cv.footer.Render(buf, footerRegion)
	}

	// Render overlays (dialogs) on top of everything.
	// Each overlay receives the full screen region so it can centre itself.
	fullRegion := Region{X: 0, Y: 0, Width: w, Height: h}
	for _, overlay := range cv.overlays {
		overlay.Render(buf, fullRegion)
	}
	// Render persistent overlays (e.g. NotificationManager) above modal overlays.
	for _, overlay := range cv.persistentOverlays {
		overlay.Render(buf, fullRegion)
	}
}

// updateStatusBar pushes current focus keybindings to the auto StatusBar.
func (cv *Canvas) updateStatusBar() {
	if cv.autoBar == nil {
		return
	}
	bindings := cv.focus.CurrentKeyBindings()
	// Append global bindings so they appear in the status bar alongside the
	// focused widget's own hints.  Global bindings come after widget-specific
	// ones so the widget's shortcuts are listed first.
	bindings = append(bindings, cv.globalBindings...)
	// Always append navigation shortcuts.
	bindings = append(bindings,
		KeyBinding{Key: tcell.KeyTab, Label: "Tab", Description: "Next"},
		KeyBinding{Key: tcell.KeyEscape, Label: "Esc", Description: "Quit"},
	)
	cv.autoBar.SetBindings(bindings)
}

// InvalidateLayout rebuilds the focus tree after the component tree changes.
// Call this after dynamically adding/removing components.
func (cv *Canvas) InvalidateLayout() {
	if cv.body != nil {
		cv.focus.Collect(cv.body)
	}
	if cv.focusStyle != (latte.Style{}) {
		cv.injectFocusStyle()
	}
	cv.updateStatusBar()
}

// injectFocusStyle propagates the canvas-level focus style into every collected
// focusable component that has not set its own per-component FocusStyle.
func (cv *Canvas) injectFocusStyle() {
	for _, node := range cv.focus.Nodes() {
		if inj, ok := node.(FocusStyleInjector); ok {
			inj.SetFocusStyle(cv.focusStyle)
		}
	}
}

// applyTheme walks the full component tree (header, body, footer) and calls
// ApplyTheme on every node that implements ThemeReceiver.
// It also sets the canvas background style from the theme so the caller does
// not need to pass WithStyle separately.
func (cv *Canvas) applyTheme() {
	t := *cv.theme

	// Propagate theme canvas background if the caller has not set an explicit style.
	if cv.style == (latte.Style{}) {
		cv.style = t.Canvas
	}

	applyThemeTree(cv.header, t)
	applyThemeTree(cv.body, t)
	applyThemeTree(cv.footer, t)
}

// applyThemeTree recursively applies t to c and all its descendants.
func applyThemeTree(c Component, t latte.Theme) {
	if c == nil {
		return
	}
	if tr, ok := c.(ThemeReceiver); ok {
		tr.ApplyTheme(t)
	}
	if l, ok := c.(Layout); ok {
		for _, child := range l.Children() {
			applyThemeTree(child, t)
		}
	}
}

// GetWidgetByID performs a depth-first search of the component tree (header,
// body, footer) and returns the first Component whose ID matches id.
// Returns nil if no match is found.
func (cv *Canvas) GetWidgetByID(id string) Component {
	for _, root := range []Component{cv.header, cv.body, cv.footer} {
		if c := findByID(root, id); c != nil {
			return c
		}
	}
	return nil
}

// GetValue retrieves the value held by the component with the given ID.
// The component must implement oat.ValueGetter (all built-in widgets do).
// Returns (value, true) on success, or (nil, false) if no widget is found or
// if the widget does not implement ValueGetter.
func (cv *Canvas) GetValue(id string) (interface{}, bool) {
	c := cv.GetWidgetByID(id)
	if c == nil {
		return nil, false
	}
	if vg, ok := c.(ValueGetter); ok {
		return vg.GetValue(), true
	}
	return nil, false
}

// findByID performs a DFS search of the subtree rooted at c.
func findByID(c Component, id string) Component {
	if c == nil {
		return nil
	}
	if ider, ok := c.(IDer); ok && ider.GetID() == id {
		return c
	}
	if l, ok := c.(Layout); ok {
		for _, child := range l.Children() {
			if found := findByID(child, id); found != nil {
				return found
			}
		}
	}
	return nil
}

// dispatchGlobal checks ev against the canvas-level global bindings and invokes
// the first matching handler.  Returns true if a binding consumed the event.
// This is called after the focused widget has already had a chance to handle the
// key, so focused widgets can shadow global bindings when needed.
func (cv *Canvas) dispatchGlobal(ev *tcell.EventKey) bool {
	for _, b := range cv.globalBindings {
		if b.Handler != nil && matchesBinding(ev, b) {
			b.Handler()
			return true
		}
	}
	return false
}

// SetTheme replaces the active theme and immediately re-applies it to the
// entire component tree (header, body, footer, and any currently mounted
// overlays or persistent overlays).  The canvas background style is also
// updated from the new theme.
//
// SetTheme is safe to call from inside key-event callbacks (it runs on the
// main goroutine).  To trigger a re-render after the call, the event loop
// will automatically re-render on the next tick; you do not need to call
// NotifyChannel manually.
func (cv *Canvas) SetTheme(t latte.Theme) {
	cv.theme = &t

	// Reset canvas background so applyTheme re-derives it from the new theme.
	cv.style = latte.Style{}

	cv.applyTheme()

	// Re-apply to overlay stack and persistent overlays.
	for _, o := range cv.overlays {
		applyThemeTree(o, t)
	}
	for _, o := range cv.persistentOverlays {
		applyThemeTree(o, t)
	}
}

// GetTheme returns a pointer to the active theme. The pointer is the same
// value stored internally by SetTheme — do not mutate the pointed-to Theme;
// use SetTheme with a new value instead.
//
// Returns nil if no theme has been set (i.e. NewCanvas was called without
// WithTheme and SetTheme has never been called).
func (cv *Canvas) GetTheme() *latte.Theme {
	return cv.theme
}
