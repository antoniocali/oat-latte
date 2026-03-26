---
sidebar_position: 1
title: Introduction
description: What oat-latte is, why it exists, and how it compares to Bubble Tea.
---

# oat-latte

oat-latte is a component-based TUI (terminal UI) framework for Go. It gives you a declarative component model, a two-pass layout engine, and a full widget library — all rendered through [tcell](https://github.com/gdamore/tcell).

## Three things to know

Every oat-latte application is built from exactly three kinds of thing:

| | What it does |
|---|---|
| **Canvas** | Owns the terminal screen and runs the event loop. You create one Canvas per application and call `Run()`. |
| **Layouts** | Containers that position their children — `VBox`, `HBox`, `Border`, `Grid`, and more. You nest them to build any UI structure. |
| **Widgets** | The interactive and display elements — `Text`, `Button`, `EditText`, `List`, `CheckBox`, `ProgressBar`, and more. Widgets live inside layouts. |

That's the whole model. A minimal app looks like this:

```go
body := layout.NewBorder(
    layout.NewVBox(
        widget.NewText("Hello!"),
        widget.NewButton("Quit", func() { app.Quit() }),
    ),
).WithTitle("My App")

app := oat.NewCanvas(
    oat.WithTheme(latte.ThemeDark),
    oat.WithBody(body),
)
app.Run()
```

You pick widgets, arrange them in layouts, hand the root layout to a Canvas, and call `Run()`. Everything else — focus, key handling, theming, redraws — is handled for you.

## Why it exists

After spending time with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and genuinely enjoying it, I found myself wondering whether there was space for a different mental model — one that felt more familiar to developers coming from GUI frameworks.

My background includes a good stretch of Flutter development, and Flutter's approach to UI always clicked with me: a component tree, a constraint-based layout pass, widgets that know how to measure and paint themselves, and a clear separation between the "what does it want to be?" and "where does it actually live?" phases. When I sat down to build terminal tools in Go, that model felt like the natural fit. oat-latte is the result of trying to bring that ideology to the terminal.

## How it compares to Bubble Tea

[Bubble Tea](https://github.com/charmbracelet/bubbletea) is a wonderful framework and the right choice for many programs. It is worth understanding the difference in philosophy before choosing one over the other.

### Bubble Tea's model

Bubble Tea is built on the [Elm architecture](https://guide.elm-lang.org/architecture/): your entire application is a single `Model`, events arrive as `Msg` values, an `Update` function returns the next model, and a `View` function renders it as a string. The framework is intentional about keeping you close to the metal — you compose views by concatenating strings and ANSI sequences, which gives you complete control over every character on screen.

This is a powerful model. It is explicit, pure, and very testable. If you want to deeply understand what your TUI is doing at every step, or if you need fine-grained control over rendering, Bubble Tea is likely the better tool.

### oat-latte's model

oat-latte takes a different approach. The UI is a **tree of components**. Layout containers own their children, pass constraints down, and allocate regions. Widgets are self-contained: they own their state, handle their own key events, and paint themselves. You wire them together with callbacks rather than routing messages through a central update function.

This makes it straightforward to build structured UIs — forms, dashboards, list-detail views — without writing layout arithmetic by hand. The trade-off is that you are further from the raw terminal than you would be with Bubble Tea.

### When to use which

| | oat-latte | Bubble Tea |
|---|---|---|
| Mental model | Component tree (Flutter-like) | Elm architecture |
| Layout | Automatic | Manual string composition |
| State management | Per-widget, callback-driven | Central model + Update |
| Best for | Structured UIs, forms, dashboards | Any TUI; especially custom or artistic layouts |
| Control over rendering | High-level | Low-level |

If you want to go deep into how terminal UIs work at the core — raw sequences, full control over the render loop, and a functional reactive style — Bubble Tea is still the better approach. oat-latte is for when you want to skip the layout plumbing and think in components.

## Next step

[Install oat-latte →](./installation)
