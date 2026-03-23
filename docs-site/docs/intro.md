---
sidebar_position: 1
title: Introduction
description: What oat-latte is and why you might use it.
---

# oat-latte

oat-latte is a component-based TUI (terminal UI) framework for Go. It gives you a declarative component model, a two-pass layout engine, and a full widget library — all rendered through [tcell](https://github.com/gdamore/tcell).

## What you get

- **Component model** — every element implements a two-method interface: `Measure` and `Render`. Parents ask children for their size, then hand them a region to draw into.
- **Layout primitives** — `VBox`, `HBox`, `Grid`, `Border`, `Padding`, `Dialog`, and flex spacers (`VFill`, `HFill`) cover the vast majority of real layouts without custom sizing code.
- **Widget library** — `Text`, `Button`, `CheckBox`, `EditText` (single- and multi-line), `List`, `Label`, `ProgressBar`, `StatusBar`, and `NotificationManager` out of the box.
- **Focus system** — automatic DFS-ordered Tab/Shift-Tab cycling; `HandleKey` returns a boolean so components can pass events up the chain; programmatic `FocusByRef` for instant jumps.
- **Theme system** — five built-in themes (`Default`, `Dark`, `Light`, `Dracula`, `Nord`); a `Style.Merge` cascade lets per-widget overrides survive theme application.

## When to use it

oat-latte is a good fit when you need to ship a terminal tool with a real UI — a database browser, a task manager, a log viewer, a dashboard — and you want a structured component model rather than raw terminal drawing calls.

It is not a web framework and has no browser dependency. Everything runs in a terminal.

## Next step

[Install oat-latte →](./installation)
