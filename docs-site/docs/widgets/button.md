---
sidebar_position: 2
title: Button
description: Clickable action trigger widget.
---

# Button

`widget.Button` is a focusable action trigger. It fires a callback when the user presses `Enter` or `Space`.

All built-in themes set `Border: BorderSingle` and `RoundedCorner: true` on the `Button` token, so buttons automatically render with rounded arc corners:

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
| `WithStyle(s latte.Style)` | Override the display style |
| `WithID(id string)` | Set a stable identifier for `Canvas.GetValue(id)` |
| `WithRoundedCorner(bool)` | Explicitly enable or disable arc corners (`╭╮╰╯`) |

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

The theme supplies focus colours on top of whatever style you set. Use `WithStyle` to control the unfocused appearance only.

## WithRoundedCorner

```go
// Explicit opt-in — overrides the theme's RoundedCorner setting for this button.
btn := widget.NewButton("OK", fn).
    WithRoundedCorner(true)

// Explicit opt-out — disables arc corners even when the active theme has RoundedCorner: true.
btn := widget.NewButton("Cancel", fn).
    WithRoundedCorner(false)
```

- `true` — draws arc corners (`╭╮╰╯`) instead of square ones (`┌┐└┘`).
- `false` — disables arc corners regardless of the theme's `RoundedCorner` setting.
- **Once called, this explicit choice overrides the theme** for this button. Buttons that never call `WithRoundedCorner` inherit the theme's `RoundedCorner` value automatically via `ApplyTheme`.
- If the button's resolved border style is **incompatible** with arc corners (`BorderDouble`, `BorderThick`, `BorderDashed`), the rounded-corner request is a **silent no-op** — the button keeps its original square corners without panicking.

```go
// Silent no-op example — dashed border stays square, no panic.
btn := widget.NewButton("Dashed", fn).
    WithStyle(latte.Style{Border: latte.BorderDashed}).
    WithRoundedCorner(true)  // silently ignored; dashed stays square
```

:::tip
Because all built-in themes set `RoundedCorner: true`, you do not need to call `WithRoundedCorner(true)` on every button — it happens automatically when a theme is applied. Only call it explicitly when you want to override the theme's default.
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
