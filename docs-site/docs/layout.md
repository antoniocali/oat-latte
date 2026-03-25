---
sidebar_position: 4
title: Layout
description: VBox, HBox, Border, Padding, Grid — building screens with oat-latte layout containers.
---

# Layout

Layout containers hold children and handle all size negotiation between them. Every container implements `Component` (so it can be nested anywhere) and `Layout` (so the framework can walk its children).

## VBox and HBox

`VBox` stacks children **vertically**; `HBox` places them **horizontally**.

```go
// Fixed children take their natural size.
// Flex children share whatever space remains.
vbox := layout.NewVBox()
vbox.AddChild(widget.NewText("Label"))
vbox.AddFlexChild(myEditText, 1)       // weight 1
vbox.AddFlexChild(anotherField, 2)     // weight 2 = twice as tall

// VFill is a shorthand flex spacer with weight 1.
vbox.AddChild(layout.NewVFill())

// Fixed gap of exactly 1 row.
vbox.AddChild(layout.NewVFill().WithMaxSize(1))

// HBox variadic constructor for when all children are fixed-size.
hbox := layout.NewHBox(child1, child2, child3)

// Add a flex progress bar that fills remaining horizontal space.
hbox.AddFlexChild(progressBar, 1)
```

:::tip
Use `AddFlexChild` with weight `1` on the main content area and fixed-size `AddChild` for everything else (headers, footers, button rows). This is the most common layout pattern.
:::

## Border

Wraps a single child with a configurable box border and an optional title stamped into the top rule.

```go
panel := layout.NewBorder(innerComponent).
    WithTitle("My Panel").
    WithTitleStyle(latte.Style{Bold: true})

// Rounded corners — ╭─ My Panel ──╮ / ╰──────────╯
panel := layout.NewBorder(innerComponent).
    WithTitle("My Panel").
    WithRoundedCorner(true)

// Custom style (e.g. explicit padding):
panel := layout.NewBorder(innerComponent).
    WithStyle(latte.Style{Padding: latte.Insets{Bottom: 1}}).
    WithTitle("My Panel")
```

When any descendant of the border is focused, `Border` automatically promotes its border color to the theme's `FocusBorder` token. No extra code needed.

### WithRoundedCorner

`WithRoundedCorner(true)` switches the border style to `BorderRounded`, giving it arc corners (`╭╮╰╯`) instead of the default square ones (`┌┐└┘`). `WithRoundedCorner(false)` restores `BorderSingle`.

**Compatibility:** Arc corners exist in Unicode only for light-weight strokes (`─` `│`). Calling `WithRoundedCorner(true)` on a border whose style is `BorderDouble`, `BorderThick`, or `BorderDashed` will **panic** at construction time, because those strokes have no matching arc corner codepoints and the result would be visually broken.

| Border style | `WithRoundedCorner(true)` |
|---|---|
| `BorderSingle` (default) | Allowed → becomes `BorderRounded` |
| `BorderDouble` | Panics — no double-stroke arc corners in Unicode |
| `BorderThick` | Panics — no heavy-stroke arc corners in Unicode |
| `BorderDashed` | Panics — arc corners don't connect to dashed strokes |

To use rounded corners with a non-default style, switch the entire border style via `WithStyle`:

```go
// This is already BorderRounded — no need for WithRoundedCorner.
panel := layout.NewBorder(innerComponent).
    WithStyle(latte.Style{Border: latte.BorderRounded}).
    WithTitle("My Panel")
```

## Padding

Adds blank space around a single child.

```go
// Uniform padding — 1 cell on all sides.
padded := layout.NewPaddingUniform(child, 1)

// Asymmetric padding.
padded := layout.NewPadding(child, latte.Insets{Top: 1, Right: 2, Bottom: 1, Left: 2})
```

## Grid

Positions children in a fixed-size rows × columns grid with equal cell sizes.

```go
g := layout.NewGrid(2, 3)              // 2 rows, 3 cols
g.AddChildAt(widget, 0, 0, 1, 1)      // row 0, col 0, rowSpan 1, colSpan 1
g.AddChildAt(wide,   1, 0, 1, 3)      // spans all 3 columns
g.WithGap(0, 1)                        // 0 row gap, 1 col gap
```

## HFill and VFill

Flex spacers for use inside `HBox` / `VBox`.

```go
// Push content to either end of an HBox.
row := layout.NewHBox()
row.AddChild(leftLabel)
row.AddChild(layout.NewHFill())   // fills the gap
row.AddChild(rightLabel)

// Fixed-size spacer (e.g. a 1-row vertical gap).
vbox.AddChild(layout.NewVFill().WithMaxSize(1))
```
