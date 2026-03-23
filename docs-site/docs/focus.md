---
sidebar_position: 7
title: Focus System
description: How keyboard focus works in oat-latte — automatic collection, Tab cycling, key dispatch, and programmatic jumps.
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

```
Tab / Shift+Tab
  └─ FocusManager.Next() / Prev()

Any other key
  └─ FocusManager.Dispatch(ev)
       ├─ Walk KeyBindings() of focused component
       │    └─ If binding has a Handler and matches → call Handler, consumed
       └─ else → focused.HandleKey(ev)
                  ├─ true  → consumed, done
                  └─ false → canvas tries ←/→ focus cycling
```

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

Jump focus directly to any widget by pointer identity:

```go
app.FocusByRef(myEditText)
```

The target must be in the current focus tree (body or active dialog). Use this to direct the user to a specific field after an action — for example, pressing a shortcut key that opens an editor and immediately focuses the title input.

## Dialog focus confinement

While any dialog is visible, all key events (including Tab and arrows) are routed exclusively to the dialog's focus tree. The body receives nothing. `HideDialog()` restores the body focus tree automatically.

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

## Goroutine safety

Key event handlers run on the **main goroutine**. Do not update UI state from a background goroutine. Instead, send to `app.NotifyChannel()` to trigger a re-render after the background work is done:

```go
go func() {
    result := doExpensiveWork()
    myText.SetText(result)                 // safe: set state
    app.NotifyChannel() <- time.Now()      // trigger re-render
}()
```
