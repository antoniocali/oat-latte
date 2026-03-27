---
sidebar_position: 7
title: Focus System
description: How keyboard focus works in oat-latte — automatic collection, Tab cycling, key dispatch, global bindings, and programmatic jumps.
---

# Focus System

## How focus is collected

When `Canvas.Run()` starts, and again whenever the overlay stack changes or you call `InvalidateLayout()`, the framework performs a **depth-first search** over the component tree and registers every `Focusable` node in order. DFS left-to-right corresponds to visual top-left to bottom-right.

The first node in the list receives initial focus, unless you override this with `oat.WithPrimary(f)`.

## Cycling

| Key | Action |
|---|---|
| `Tab` | Move to the next focusable (`FocusManager.Next()`) |
| `Shift+Tab` | Move to the previous focusable (`FocusManager.Prev()`) |
| `←` / `↑` | Cycle backward if the focused component does not consume the key |
| `→` / `↓` | Cycle forward if the focused component does not consume the key |

This means interactive widgets (like `List`, which handles `↑`/`↓` internally) keep those keys for themselves. Non-interactive containers let arrows fall through, making arrow navigation "free" for non-text layouts.

## Key dispatch

Every key event travels through a fixed priority chain:

```
Tab / Shift+Tab
  └─ FocusManager.Next() / Prev()

Any other key
  └─ FocusManager.Dispatch(ev)
       ├─ Walk KeyBindings() of focused component
       │    └─ If binding has a Handler and matches → call Handler, consumed
       └─ else → focused.HandleKey(ev)
                  ├─ true  → consumed, done
                  └─ false → canvas.dispatchGlobal(ev)
                               ├─ matching global binding → call Handler, consumed
                               └─ no match → canvas tries ←/→ focus cycling
```

The focused widget **always gets first refusal**. Global bindings only fire when the widget does not consume the key. This means a widget can shadow a global shortcut when that's the right behaviour — for example, `Esc` in an `EditText` cancels editing rather than quitting the app.

## Global key bindings

Register app-level shortcuts that fire regardless of which widget is focused using `oat.WithGlobalKeyBinding`:

```go
themes := []latte.Theme{latte.ThemeDark, latte.ThemeLight, latte.ThemeDracula, latte.ThemeNord}
current := 0

app := oat.NewCanvas(
    oat.WithTheme(themes[current]),
    oat.WithBody(body),
    oat.WithGlobalKeyBinding(
        oat.KeyBinding{
            Key:         tcell.KeyCtrlT,
            Mod:         tcell.ModCtrl,
            Label:       "^T",
            Description: "Toggle theme",
            Handler: func() {
                current = (current + 1) % len(themes)
                app.SetTheme(themes[current])
            },
        },
        oat.KeyBinding{
            Key:         tcell.KeyCtrlH,
            Label:       "^H",
            Description: "Help",
            Handler:     func() { app.ShowDialog(helpDialog) },
        },
    ),
)
```

Key behaviours:

- **Variadic and accumulating** — pass multiple bindings in one call, or call `WithGlobalKeyBinding` multiple times; bindings are always appended, never replaced.
- **Status bar integration** — global bindings are shown in the footer status bar alongside the focused widget's own hints.
- **Shadowed by focused widgets** — if the focused widget returns `true` from `HandleKey` for the same key, the global binding never fires.

:::tip When to use global bindings vs the proxy pattern
Use `WithGlobalKeyBinding` for shortcuts that truly belong to the whole application — theme switching, a global help overlay, quit confirmation. Use the [proxy pattern](#the-proxy-pattern) when a shortcut only makes sense in the context of one specific widget (e.g. `n` to create a new item when a list is focused).
:::

## The proxy pattern

When you need a widget to handle extra keys without modifying the widget itself, wrap it in a thin proxy:

```go
type myListProxy struct {
    *widget.List
    app *App
}

func (p *myListProxy) HandleKey(ev *oat.KeyEvent) bool {
    if ev.Key() == tcell.KeyRune && ev.Rune() == 'n' {
        p.app.showNewItemDialog()
        return true // consumed
    }
    return p.List.HandleKey(ev) // delegate everything else
}

func (p *myListProxy) KeyBindings() []oat.KeyBinding {
    extra := []oat.KeyBinding{
        {Key: tcell.KeyRune, Rune: 'n', Label: "n", Description: "New item"},
    }
    return append(extra, p.List.KeyBindings()...)
}
```

Use the proxy in place of the original widget everywhere:

```go
proxy := &myListProxy{List: rawList, app: a}
vbox.AddFlexChild(proxy, 1)                // layout sees the proxy
app := oat.NewCanvas(..., oat.WithPrimary(proxy))
```

## Programmatic focus

Jump focus directly to any widget by pointer identity by calling `app.FocusByRef(target)`:

The target must be in the current focus tree (body or active dialog). Use this to direct the user to a specific field after an action — for example, pressing a shortcut key that opens an editor and immediately focuses the title input.

## Dialog focus confinement

While any dialog is visible, all key events (including Tab and arrows) are routed exclusively to the dialog's focus tree. The body receives nothing. `HideDialog()` restores the body focus tree automatically.

:::note Global bindings inside dialogs
Global bindings still fire inside dialogs — if the key is not consumed by the focused widget within the dialog, the global binding chain runs as normal.
:::

## FocusGuard — context-aware Tab cycling

Implement `oat.FocusGuard` to dynamically exclude a component and its entire subtree from Tab cycling:

```go
type FocusGuard interface {
    IsFocusable() bool
}
```

When `IsFocusable()` returns `false`, the tree walker skips the node **and all its descendants**. This is the right tool when whole panels should be unreachable depending on application state.

**Example — two-mode editor where only one panel is reachable at a time:**

```go
// Guard wrapping the list panel — reachable only when NOT in editor mode.
type listGuard struct {
    oat.Component
    app *App
}
func (g *listGuard) IsFocusable() bool { return !g.app.editorMode }

// Guard wrapping the editor panel — reachable only when in editor mode.
type editorGuard struct {
    *widget.EditText
    app *App
}
func (g *editorGuard) IsFocusable() bool { return g.app.editorMode }
```

After toggling the mode, call `InvalidateLayout()` so the focus tree is rebuilt, then set the desired initial focus:

```go
func (a *App) setEditorMode(on bool) {
    if a.editorMode == on {
        return
    }
    a.editorMode = on
    a.canvas.InvalidateLayout()
    if on {
        a.canvas.FocusByRef(a.titleInput) // jump into editor
    } else {
        a.canvas.FocusByRef(a.list)       // return to list
    }
}
```

## Advertising shortcuts in the status bar

Implement `KeyBindings()` on any focusable component to advertise its shortcuts. The `StatusBar` widget reads these and renders them as `[key] Description` hints.

```go
func (w *MyWidget) KeyBindings() []oat.KeyBinding {
    return []oat.KeyBinding{
        // Handler nil = display-only hint; non-nil = executed by Dispatch.
        {Key: tcell.KeyCtrlS, Label: "^S", Description: "Save"},
        {Key: tcell.KeyCtrlG, Label: "^G", Description: "Cancel"},
    }
}
```

Global bindings registered with `WithGlobalKeyBinding` are automatically appended to the status bar after the focused widget's hints — you do not need to advertise them manually.

## Goroutine safety

Key event handlers run on the **main goroutine**. Do not update UI state from a background goroutine. Instead, send to `app.NotifyChannel()` to trigger a re-render after the background work is done:

```go
go func() {
    result := doExpensiveWork()
    myText.SetText(result)                 // safe: set state
    app.NotifyChannel() <- time.Now()      // trigger re-render
}()
```
