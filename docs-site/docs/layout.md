---
sidebar_position: 4
title: Layout
description: VBox, HBox, Border, Padding, Grid — building screens with oat-latte layout containers.
---

# Layout

Layout containers hold children and handle all size negotiation between them. Every container implements `Component` (so it can be nested anywhere) and `Layout` (so the framework can walk its children).

## VBox and HBox

`VBox` stacks children **vertically**; `HBox` places them **horizontally**.

### AddChild vs AddFlexChild

This is the core sizing decision you make for every child.

**`AddChild(c)`** — the child takes its **natural size** and no more. The box asks the child how big it wants to be (via `Measure`) and allocates exactly that. Use this for things whose size is inherently fixed: labels, buttons, a single-line input, a status row.

**`AddFlexChild(c, weight)`** — the child participates in **flex distribution**. After all fixed children have taken their space, whatever is left over is divided among flex children proportionally by weight. A child with weight `2` gets twice the space of one with weight `1`. Use this for the main content area, multi-line editors, lists — anything that should grow to fill available space.

```
┌─────────────────────────────┐
│  AddChild  → natural size   │  ← e.g. a title row: always 1 row tall
│  AddChild  → natural size   │  ← e.g. a button row: always 1 row tall
│                             │
│  AddFlexChild weight 1      │  ← fills all remaining space
│                             │
│  AddChild  → natural size   │  ← e.g. a status bar: always 1 row tall
└─────────────────────────────┘
```

When **multiple** flex children are present, space is split by the ratio of their weights:

```
remaining space = 30 rows
  flex child A  weight 1  → 10 rows  (1/3)
  flex child B  weight 2  → 20 rows  (2/3)
```

:::tip
A `VFill` or `HFill` added via `AddChild` auto-promotes itself to a flex slot with weight `1` — it is a shorthand for `AddFlexChild(layout.NewVFill(), 1)`.
:::

```go
vbox := layout.NewVBox()
vbox.AddChild(widget.NewText("Label"))        // fixed: takes its natural height
vbox.AddFlexChild(myEditText, 1)              // flex weight 1
vbox.AddFlexChild(anotherField, 2)            // flex weight 2 → twice as tall as myEditText

// VFill shorthand — equivalent to AddFlexChild(layout.NewVFill(), 1)
vbox.AddChild(layout.NewVFill())

// Fixed gap of exactly 1 row — WithMaxSize caps the spacer's flex growth
vbox.AddChild(layout.NewVFill().WithMaxSize(1))

// HBox variadic constructor: all children are fixed-size
hbox := layout.NewHBox(child1, child2, child3)

// Flex progress bar that fills remaining horizontal space
hbox.AddFlexChild(progressBar, 1)
```

:::tip
The most common pattern: one `AddFlexChild(mainContent, 1)` for the central area, `AddChild` for everything else (header, footer, button rows). Everything snaps to its natural size; the main area fills the rest.
:::

## Border

Wraps a single child with a configurable box border and an optional title stamped into the top rule.

```go
// Default — title at the left edge.
panel := layout.NewBorder(innerComponent).
    WithTitle("My Panel").
    WithTitleStyle(latte.Style{Bold: true})

// Centred title.
panel := layout.NewBorder(innerComponent).
    WithTitle("My Panel", oat.AnchorCenter)

// Title at the right edge.
panel := layout.NewBorder(innerComponent).
    WithTitle("My Panel", oat.AnchorRight)

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

### WithTitle anchor

`WithTitle` accepts an optional `oat.Anchor` as a second argument that controls where the title text is stamped in the top border rule:

```go
func (b *Border) WithTitle(title string, anchor ...oat.Anchor) *Border
```

| Anchor | Result |
|---|---|
| `oat.AnchorLeft` (default) | `╭─ My Panel ──────────╮` |
| `oat.AnchorCenter` | `╭──── My Panel ────────╮` |
| `oat.AnchorRight` | `╭────────── My Panel ──╮` |

Omitting the anchor is the same as passing `oat.AnchorLeft`. This keeps all existing call sites (`WithTitle("My Panel")`) unchanged.

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
