---
sidebar_position: 9
title: StatusBar
description: Auto-populated key hint footer.
---

# StatusBar

`widget.StatusBar` renders a single-row footer showing `[key] Description` hints for the currently focused component's advertised `KeyBindings`. It is updated automatically by the canvas on every focus change.

## Setup

```go
statusBar := widget.NewStatusBar()

app := oat.NewCanvas(
    oat.WithTheme(latte.ThemeDark),
    oat.WithBody(body),
    oat.WithAutoStatusBar(statusBar),  // mounts as footer; auto-updates
)
```

That is all you need to do. The canvas calls `statusBar.SetBindings(bindings)` whenever focus changes, so the hints always reflect the active widget.

## Builder options

| Method | Description |
|---|---|
| `WithStyle(s latte.Style)` | Override the footer display style |
| `WithID(id string)` | Set a stable identifier |

The theme supplies the default footer style via the `Footer` token. Use `WithStyle` only if you want to deviate from the theme.

## What appears

The status bar shows all `KeyBindings` returned by the currently focused component's `KeyBindings()` method, plus global canvas hints (`Tab · Next`, `Esc · Quit`).

Each binding is rendered as `[^S] Save` — the bracketed key in bold accent colour, the description in muted colour. Bindings with an empty `Description` show only the bracketed key.

## Advertising shortcuts

Any focusable component can expose shortcuts by implementing `KeyBindings()`:

```go
func (w *MyWidget) KeyBindings() []oat.KeyBinding {
    return []oat.KeyBinding{
        {Key: tcell.KeyCtrlS, Label: "^S", Description: "Save"},
        {Key: tcell.KeyCtrlG, Label: "^G", Description: "Cancel"},
        // Handler: nil means display-only; the actual logic is in HandleKey.
    }
}
```

Bindings with a non-nil `Handler` are also dispatched by the canvas key router without needing custom `HandleKey` code.
