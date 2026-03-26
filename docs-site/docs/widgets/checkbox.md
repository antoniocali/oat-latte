---
sidebar_position: 5
title: CheckBox
description: Boolean toggle widget.
---

# CheckBox

`widget.CheckBox` is a focusable boolean toggle rendered as `[ ] label` (unchecked) or `[x] label` (checked).

## Basic usage

Create a checkbox with `widget.NewCheckBox("Enable notifications")` and add it to a layout container.

## Builder options

| Method | Description |
|---|---|
| `WithStyle(s latte.Style)` | Override the display style |
| `WithID(id string)` | Set a stable identifier for `Canvas.GetValue(id)` |
| `WithOnToggle(fn func(bool))` | Callback fired when the checked state changes |

## Toggle callback

```go
cb := widget.NewCheckBox("Dark mode").
    WithOnToggle(func(checked bool) {
        if checked {
            // apply dark mode
        }
    })
```

## Reading and setting state

```go
cb.SetChecked(true)
isOn := cb.IsChecked()
```

`CheckBox` also implements `oat.ValueGetter`: `Canvas.GetValue(id)` returns the current checked state as a `bool`.

## Key bindings

| Key | Action |
|---|---|
| `Space` or `Enter` | Toggle checked state |
| `Tab` / `Shift+Tab` | Move focus |
