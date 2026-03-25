---
sidebar_position: 7
title: ProgressBar
description: Horizontal progress indicator with configurable percentage label positioning.
---

# ProgressBar

`widget.ProgressBar` renders a horizontal progress indicator for values between `0.0` and `1.0`. It is a display-only widget with no interactive key bindings.

## Basic usage

```go
pb := widget.NewProgressBar()
pb.SetValue(0.65)   // 65%
```

## Percentage label

By default the percentage label is shown at the left edge. Use `WithPercentage` to control both visibility and position:

```go
// Left (default) — " 65% ██████░░░░"
pb := widget.NewProgressBar().WithPercentage(true)

// Centre — "███ 65% ░░░░"
pb := widget.NewProgressBar().WithPercentage(true, oat.AnchorCenter)

// Right — "██████░░░░ 65%"
pb := widget.NewProgressBar().WithPercentage(true, oat.AnchorRight)

// Hidden
pb := widget.NewProgressBar().WithPercentage(false)
```

The `anchor` parameter is optional — omit it to get `oat.AnchorLeft` by default.

| Anchor | Layout |
|---|---|
| `oat.AnchorLeft` (default) | Label precedes the bar: `" 65% ██░░░░"` |
| `oat.AnchorCenter` | Label stamped into the midpoint: `"██ 65% ░░"` |
| `oat.AnchorRight` | Label appended after the bar: `"████░░ 65%"` |

:::note
`WithShowPercent(bool)` still works as a deprecated shorthand but does not let you set the anchor. Prefer `WithPercentage`.
:::

## Builder options

| Method | Description |
|---|---|
| `WithPercentage(show bool, anchor ...oat.Anchor)` | Control label visibility and position (default `true`, `AnchorLeft`) |
| `WithStyle(s latte.Style)` | Override the display style (colour, background) |
| `WithID(id string)` | Set a stable identifier for `Canvas.GetValue(id)` |
| `WithFillChar(r rune)` | Rune used for the filled portion (default `█`) |
| `WithEmptyChar(r rune)` | Rune used for the empty portion (default `░`) |
| `WithShowPercent(show bool)` | Deprecated — use `WithPercentage` instead |

## Custom fill characters

```go
pb := widget.NewProgressBar().
    WithFillChar('=').
    WithEmptyChar('-').
    WithPercentage(true, oat.AnchorRight)
// Renders as:  ======---- 60%
```

## Per-item colour

After theme application you can override the fill colour per-item using `WithStyle`:

```go
pb := widget.NewProgressBar().
    WithStyle(latte.Style{FG: latte.ColorRed})
```

`SetStyle` is still available as a deprecated alias for `WithStyle` when you need to mutate the style after construction (e.g. in a list renderer).

## In a flex layout

`ProgressBar` grows to fill available horizontal space when given a flex weight:

```go
row := layout.NewHBox()
row.AddChild(widget.NewText("Progress: "))
row.AddFlexChild(pb, 1)
row.AddChild(widget.NewText(" 4/10"))
```

## Reading the value

```go
v := pb.Progress() // float64, 0.0–1.0
```

`ProgressBar` also implements `oat.ValueGetter`: `Canvas.GetValue(id)` returns the current value as a `float64`.

## Updating from a goroutine

`SetValue` is safe to call from any goroutine. Use `widget.NotificationManager` (via `oat.WithNotificationManager`) to trigger a re-render:

```go
notifs := widget.NewNotificationManager()
app := oat.NewCanvas(
    oat.WithBody(body),
    oat.WithNotificationManager(notifs),
)

go func() {
    for i := 0; i <= 100; i++ {
        pb.SetValue(float64(i) / 100.0)
        notifs.Push("", widget.NotificationKindInfo, time.Millisecond) // triggers re-render
        time.Sleep(50 * time.Millisecond)
    }
}()
```
