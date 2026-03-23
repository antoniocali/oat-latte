---
sidebar_position: 1
title: Introduction
description: What oat-latte is, why it exists, and how it compares to Bubble Tea.
---

# oat-latte

oat-latte is a component-based TUI (terminal UI) framework for Go. It gives you a declarative component model, a two-pass layout engine, and a full widget library — all rendered through [tcell](https://github.com/gdamore/tcell).

## Why it exists

After spending time with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and genuinely enjoying it, I found myself wondering whether there was space for a different mental model — one that felt more familiar to developers coming from GUI frameworks.

My background includes a good stretch of Flutter development, and Flutter's approach to UI always clicked with me: a component tree, a constraint-based layout pass, widgets that know how to measure and paint themselves, and a clear separation between the "what does it want to be?" and "where does it actually live?" phases. When I sat down to build terminal tools in Go, that model felt like the natural fit. oat-latte is the result of trying to bring that ideology to the terminal.

## How it compares to Bubble Tea

[Bubble Tea](https://github.com/charmbracelet/bubbletea) is a wonderful framework and the right choice for many programs. It is worth understanding the difference in philosophy before choosing one over the other.

### Bubble Tea's model

Bubble Tea is built on the [Elm architecture](https://guide.elm-lang.org/architecture/): your entire application is a single `Model`, events arrive as `Msg` values, a `Update` function returns the next model, and a `View` function renders it as a string. The framework is intentional about keeping you close to the metal — you compose views by concatenating strings and ANSI sequences, which gives you complete control over every character on screen.

This is a powerful model. It is explicit, pure, and very testable. If you want to deeply understand what your TUI is doing at every step, or if you need fine-grained control over rendering, Bubble Tea is likely the better tool.

### oat-latte's model

oat-latte takes a different approach. Instead of a central model and a view function, the UI is a **tree of components**. Each component implements two methods:

- `Measure` — given available space as a constraint, return the desired size
- `Render` — given an allocated region, draw into a buffer

Layout containers own their children, pass constraints down, and allocate regions. Widgets are self-contained: they own their state, handle their own key events, and paint themselves. You wire them together with callbacks rather than routing messages through a central update function.

This makes it straightforward to build structured UIs — forms, dashboards, list-detail views — without writing layout arithmetic by hand. The trade-off is that you are further from the raw terminal than you would be with Bubble Tea.

### When to use which

| | oat-latte | Bubble Tea |
|---|---|---|
| Mental model | Component tree (Flutter-like) | Elm architecture |
| Layout | Automatic (Measure/Render) | Manual string composition |
| State management | Per-widget, callback-driven | Central model + Update |
| Best for | Structured UIs, forms, dashboards | Any TUI; especially custom or artistic layouts |
| Control over rendering | High-level | Low-level |
| Go deep into TUI internals | Less so | Yes — this is its strength |

If you want to go deep into how terminal UIs work at the core — raw sequences, full control over the render loop, and a functional reactive style — Bubble Tea is still the better approach and I would wholeheartedly recommend it. oat-latte is for when you want to skip the layout plumbing and think in components.

## What you get

- **Component model** — every element implements a two-method interface: `Measure` and `Render`. Parents ask children for their size, then hand them a region to draw into.
- **Layout primitives** — `VBox`, `HBox`, `Grid`, `Border`, `Padding`, `Dialog`, and flex spacers (`VFill`, `HFill`) cover the vast majority of real layouts without custom sizing code.
- **Widget library** — `Text`, `Button`, `CheckBox`, `EditText` (single- and multi-line), `List`, `Label`, `ProgressBar`, `StatusBar`, and `NotificationManager` out of the box.
- **Focus system** — automatic DFS-ordered Tab/Shift-Tab cycling; `HandleKey` returns a boolean so components can pass events up the chain; programmatic `FocusByRef` for instant jumps.
- **Theme system** — five built-in themes (`Default`, `Dark`, `Light`, `Dracula`, `Nord`); a `Style.Merge` cascade lets per-widget overrides survive theme application.

## Next step

[Install oat-latte →](./installation)
