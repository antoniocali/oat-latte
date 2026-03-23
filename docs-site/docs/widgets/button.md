---
sidebar_position: 2
title: Button
description: Clickable action trigger widget.
---

# Button

`widget.Button` is a focusable action trigger rendered as `[ label ]`. It fires a callback when the user presses `Enter` or `Space`.

## Basic usage

```go
btn := widget.NewButton("Save", func() {
    // handle the press
})
```

## Builder options

| Method | Description |
|---|---|
| `WithStyle(s latte.Style)` | Override the display style |
| `WithID(id string)` | Set a stable identifier for `Canvas.GetValue(id)` |

## With an ID

Assign an ID to retrieve the button's label later via `Canvas.GetValue(id)`:

```go
btn := widget.NewButton("Delete", onDelete).WithID("delete-btn")
```

## Custom style

```go
btn := widget.NewButton("Danger", onDelete).
    WithStyle(latte.Style{FG: latte.ColorRed, Bold: true})
```

The focus style (reversed colours by default) is always applied on top of `Style` when the button is focused. Use `WithStyle` to control the unfocused appearance only; the theme supplies focus colours.

## In a button row

The typical pattern is an `HBox` with `HFill` to push buttons to the right:

```go
btnRow := layout.NewHBox()
btnRow.AddChild(layout.NewHFill())              // pushes buttons rightward
btnRow.AddChild(cancelBtn)
btnRow.AddChild(layout.NewHFill().WithMaxSize(2)) // 2-cell gap between buttons
btnRow.AddChild(okBtn)
```

## Reading the value

`Button` implements `oat.ValueGetter`. `Canvas.GetValue(id)` returns the button's label as a `string`.

## Keyboard behavior

| Key | Action |
|---|---|
| `Enter` | Fires the press callback |
| `Space` | Fires the press callback |
| `Tab` / `Shift+Tab` | Move focus to next / previous widget |
