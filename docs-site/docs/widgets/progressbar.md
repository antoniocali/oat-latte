---
sidebar_position: 7
title: ProgressBar
description: Horizontal progress indicator.
---

# ProgressBar

`widget.ProgressBar` renders a horizontal progress indicator for values between `0.0` and `1.0`. It is a display-only widget with no interactive key bindings.

## Basic usage

```go
pb := widget.NewProgressBar()
pb.SetValue(0.65)   // 65%
```

## Builder options

| Method | Description |
|---|---|
| `WithStyle(s latte.Style)` | Override the display style (colour, background) |
| `WithID(id string)` | Set a stable identifier for `Canvas.GetValue(id)` |
| `WithShowPercent(show bool)` | Show or hide the `100%` label at the right edge (default `true`) |
| `WithFillChar(r rune)` | Rune used for the filled portion (default `█`) |
| `WithEmptyChar(r rune)` | Rune used for the empty portion (default `░`) |

## Hide the percentage label

```go
pb := widget.NewProgressBar().WithShowPercent(false)
```

## Custom fill characters

```go
pb := widget.NewProgressBar().
    WithFillChar('=').
    WithEmptyChar('-')
// Renders as:  ======-------- 60%
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

`SetValue` is safe to call from any goroutine. After updating, send to `app.NotifyChannel()` to trigger a re-render:

```go
go func() {
    for i := 0; i <= 100; i++ {
        pb.SetValue(float64(i) / 100.0)
        app.NotifyChannel() <- time.Now()
        time.Sleep(50 * time.Millisecond)
    }
}()
```
