---
sidebar_position: 2
title: Button
description: Clickable action trigger widget.
---

# Button

`widget.Button` is a focusable action trigger. It fires a callback when the user presses `Enter` or `Space`.

All built-in themes set `Border: BorderSingle` on the `Button` token, so buttons always render with a visible border box:

```
╭────────╮
│  Save  │
╰────────╯
```

## Basic usage

```go
btn := widget.NewButton("Save", func() {
    // handle the press
})
```

## Builder options

| Method | Description |
|---|---|
| `WithStyle(s latte.Style)` | Override the display style (see validation rules below) |
| `WithID(id string)` | Set a stable identifier for `Canvas.GetValue(id)` |
| `WithRoundedCorner(bool)` | Draw arc corners (`╭╮╰╯`) instead of square ones |

## Border and height

`Button.Measure` derives border presence from `b.Style` (the unfocused base style) only — focus state never affects the layout shape.

| `b.Style.Border` | `Measure` height | Rendered as |
|---|---|---|
| `BorderSingle` (theme default) | 3 rows | bordered box |
| `BorderRounded` | 3 rows | bordered box with arc corners |
| `BorderNone` / `BorderExplicitNone` | 1 row | `[ label ]` inline |

`FocusStyle` / `ButtonFocus` carry only colour and attribute overrides (`Reverse: true`, accent `BorderFG`). The layout shape is stable regardless of focus state — a button row in a dialog always measures the same height whether or not a button inside it has focus.

## WithStyle

```go
btn := widget.NewButton("Danger", onDelete).
    WithStyle(latte.Style{FG: latte.ColorRed, Bold: true})
```

`WithStyle` validates the border field immediately and **panics** if an incompatible style is passed:

| Border value | Allowed |
|---|---|
| `BorderNone` (0) | Yes |
| `BorderExplicitNone` (-1) | Yes |
| `BorderSingle` (1) | Yes |
| `BorderRounded` (2) | Yes |
| `BorderDouble` (3) | **Panics** |
| `BorderThick` (4) | **Panics** |
| `BorderDashed` (5) | **Panics** |

The theme supplies focus colours on top of whatever style you set. Use `WithStyle` to control the unfocused appearance only.

## WithRoundedCorner

```go
btn := widget.NewButton("OK", fn).
    WithRoundedCorner(true)
```

Draws arc corners (`╭╮╰╯`) instead of square ones (`┌┐└┘`).

- `true` — arc corners active. The border style must be `BorderSingle` or `BorderRounded` at render time, otherwise **panics at render time**.
- `false` — no-op (square corners are the default).

:::tip
`WithRoundedCorner(true)` is equivalent to `WithStyle(latte.Style{Border: latte.BorderRounded})`. The explicit style approach is more portable.
:::

## With an ID

Assign an ID to retrieve the button's label later via `Canvas.GetValue(id)`:

```go
btn := widget.NewButton("Delete", onDelete).WithID("delete-btn")
```

## In a button row

The typical pattern is an `HBox` with `HFill` to push buttons to the right:

```go
btnRow := layout.NewHBox()
btnRow.AddChild(layout.NewHFill())              // pushes buttons rightward
btnRow.AddChild(cancelBtn)
btnRow.AddChild(layout.NewHFill().WithMaxSize(2)) // 2-cell gap between buttons
btnRow.AddChild(okBtn)
```

Because buttons always measure `Height: 3` (with the default bordered theme), a dialog body containing a button row needs at least 3 rows allocated for that row. Account for this when sizing dialogs with `WithMaxSize`.

## Reading the value

`Button` implements `oat.ValueGetter`. `Canvas.GetValue(id)` returns the button's label as a `string`.

## Keyboard behavior

| Key | Action |
|---|---|
| `Enter` | Fires the press callback |
| `Space` | Fires the press callback |
| `Tab` / `Shift+Tab` | Move focus to next / previous widget |
