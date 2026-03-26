---
sidebar_position: 11
title: Divider
description: Horizontal and vertical rule widgets for separating layout regions in oat-latte.
---

# Divider

`widget.Divider` renders a single-cell-thick rule — either a horizontal line (`─`) placed between rows in a `VBox`, or a vertical line (`│`) placed between columns in an `HBox`. It is a display-only widget with no interactive key bindings.

## Constructors

```go
// Horizontal rule — use inside a VBox
hd := widget.NewHDivider()

// Vertical rule — use inside an HBox
vd := widget.NewVDivider()

// Axis-explicit constructor
d := widget.NewDivider(widget.AxisHorizontal)
d := widget.NewDivider(widget.AxisVertical)
```

## Basic usage

```go
vbox := layout.NewVBox(
    widget.NewText("Section A"),
    widget.NewHDivider(),
    widget.NewText("Section B"),
)

hbox := layout.NewHBox(
    layout.NewFlexChild(leftPanel, 1),
    widget.NewVDivider(),
    layout.NewFlexChild(rightPanel, 2),
)
```

## Custom line character

```go
// Double horizontal rule
hd := widget.NewHDivider().WithRune('═')

// Heavy vertical rule
vd := widget.NewVDivider().WithRune('┃')
```

## Partial-width / partial-height rules

By default the rule spans the full width (horizontal) or full height (vertical) of the space allocated by the parent. Use `DividerSize` to limit how much of that space the visible line occupies.

| Constructor | Meaning |
|---|---|
| `widget.DividerFill` | Full allocated length (default) |
| `widget.DividerFixed(n)` | Exactly `n` terminal cells |
| `widget.DividerPercent(p)` | `p`% of the allocated length (1–100) |

### Horizontal dividers — `WithMaxSize`

`WithMaxSize(size, anchor)` controls the **width** of the visible rule and its horizontal position within the allocated space:

```go
// 60% wide, centred
hd := widget.NewHDivider().
    WithMaxSize(widget.DividerPercent(60), oat.AnchorCenter)

// 20 cells wide, pushed to the right
hd := widget.NewHDivider().
    WithMaxSize(widget.DividerFixed(20), oat.AnchorRight)
```

The `anchor` argument is optional and defaults to `oat.AnchorLeft`.

| Anchor | Behaviour |
|---|---|
| `oat.AnchorLeft` (default) | Rule starts at the left edge |
| `oat.AnchorCenter` | Rule centred in the allocated width |
| `oat.AnchorRight` | Rule pushed to the right edge |

### Vertical dividers — `WithMaxSizeV`

`WithMaxSizeV(size, anchor)` controls the **height** of the visible rule and its vertical position:

```go
// 8 cells tall, centred vertically
vd := widget.NewVDivider().
    WithMaxSizeV(widget.DividerFixed(8), oat.VAnchorMiddle)

// 50% height, aligned to bottom
vd := widget.NewVDivider().
    WithMaxSizeV(widget.DividerPercent(50), oat.VAnchorBottom)
```

The `anchor` argument is optional and defaults to `oat.VAnchorTop`.

| VAnchor | Behaviour |
|---|---|
| `oat.VAnchorTop` (default) | Rule starts at the top edge |
| `oat.VAnchorMiddle` | Rule centred in the allocated height |
| `oat.VAnchorBottom` | Rule pushed to the bottom edge |

:::note Axis safety
`WithMaxSize` accepts `oat.Anchor` (H-axis). `WithMaxSizeV` accepts `oat.VAnchor` (V-axis). Passing the wrong type is a compile error — the two anchor types are deliberately kept separate.
:::

## Theming

`ApplyTheme` maps the `Muted` theme token onto the divider, so it automatically adopts the theme's secondary colour. Override the colour with `WithStyle`:

```go
hd := widget.NewHDivider().
    WithStyle(latte.Style{FG: latte.ColorRed})
```

## Builder reference

| Method | Description |
|---|---|
| `WithRune(r rune)` | Override the line character (default `'─'` / `'│'`) |
| `WithMaxSize(size, anchor ...oat.Anchor)` | H-axis: limit width and set horizontal anchor |
| `WithMaxSizeV(size, anchor ...oat.VAnchor)` | V-axis: limit height and set vertical anchor |
| `WithStyle(s latte.Style)` | Override display style (colour, attributes) |
| `WithID(id string)` | Set a stable identifier |
