---
sidebar_position: 1
title: Text
description: Static text display with word-wrap and scroll.
---

# Text

`widget.Text` displays static text. It word-wraps to fit its render width and supports vertical scrolling for content taller than its allocated region.

## Basic usage

Create a text widget with `widget.NewText("Hello, world!")` and add it to any layout container.

## Builder options

| Method | Description |
|---|---|
| `.WithID(id string)` | Assign a stable ID for `Canvas.GetValue` lookup |
| `.WithStyle(s latte.Style)` | Override the theme-supplied text style |
| `.WithWordWrap(bool)` | Enable or disable word-wrap (default: off) |
| `.WithScrollable(bool)` | Enable vertical scrolling (default: off) |

## Updating content

```go
t.SetText("New content here")
current := t.GetText()
```

## Styling

Style can be set at construction time or after:

```go
// At construction:
t := widget.NewText("Centered heading").
    WithStyle(latte.Style{TextAlign: latte.AlignCenter, Bold: true})

// Or via the zero-value default and theme override:
t := widget.NewText("Body text")   // theme fills in latte.Text style
```

## Word-wrap and scroll

```go
scrollable := widget.NewText(longContent).
    WithWordWrap(true).
    WithScrollable(true)
```

## In a layout

`Text` is a non-focusable display component. Wrap it in a flex child to let it fill available space, e.g. `vbox.AddFlexChild(widget.NewText(longContent), 1)`.

## Title widget

For a prominent heading with an optional horizontal rule beneath it, use `widget.NewTitle("My Application")` instead of a plain `Text`. It uses the theme's `Accent` token by default. To override:

```go
heading := widget.NewTitle("My Application").
    WithStyle(latte.Style{FG: latte.ColorBrightCyan, Bold: true}).
    WithSeparator(true)
```

### Title builder options

| Method | Description |
|---|---|
| `.WithStyle(s latte.Style)` | Override the default accent style |
| `.WithSeparator(bool)` | Show/hide the horizontal rule below (default: on) |
| `.WithID(id string)` | Assign a stable ID |
