---
sidebar_position: 4
title: Layout
description: VBox, HBox, Border, Padding, Grid, AlignChild — building screens with oat-latte layout containers.
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

### Cross-axis alignment

By default every child in a VBox fills the full allocated width and every child in an HBox fills the full allocated height. Cross-axis alignment lets you **shrink a child to its natural size** and pin it to one edge (or centre it) within its slot — without adding spacer widgets.

#### Box-wide default — WithHAlign / WithVAlign

`VBox.WithHAlign` sets a default horizontal alignment for all children that do not declare their own. `HBox.WithVAlign` sets a default vertical alignment for all children.

```go
// Every child right-aligned inside this VBox
vbox := layout.NewVBox(titleText, bodyText, footerText).
    WithHAlign(oat.HAlignRight)

// Every child pinned to the bottom inside this HBox
hbox := layout.NewHBox(iconText, nameText, statusText).
    WithVAlign(oat.VAlignBottom)
```

The zero value (`HAlignFill` / `VAlignFill`) is the default and keeps the existing full-stretch behaviour — no breaking change.

**HAlign values** (used by VBox):

| Value | Effect |
|---|---|
| `oat.HAlignFill` | Fill the full allocated width (default) |
| `oat.HAlignLeft` | Shrink to natural width, pin to the left |
| `oat.HAlignCenter` | Shrink to natural width, centre horizontally |
| `oat.HAlignRight` | Shrink to natural width, pin to the right |

**VAlign values** (used by HBox):

| Value | Effect |
|---|---|
| `oat.VAlignFill` | Fill the full allocated height (default) |
| `oat.VAlignTop` | Shrink to natural height, pin to the top |
| `oat.VAlignMiddle` | Shrink to natural height, centre vertically |
| `oat.VAlignBottom` | Shrink to natural height, pin to the bottom |

#### Per-widget self-alignment

Every built-in widget exposes fluent `WithHAlign` / `WithVAlign` methods that return the concrete widget type. This is the idiomatic way to set alignment inline in a builder chain:

```go
saveBtn   := widget.NewButton("Save",   fn).WithHAlign(oat.HAlignRight)
cancelBtn := widget.NewButton("Cancel", fn).WithHAlign(oat.HAlignLeft)

vbox := layout.NewVBox(saveBtn, cancelBtn)
// saveBtn  → right-aligned
// cancelBtn → left-aligned
// box default remains HAlignFill — no effect because widgets declare their own
```

For custom widgets that embed `BaseComponent` but do not yet have their own `WithHAlign`/`WithVAlign` builders, set the field directly:

```go
myWidget.BaseComponent.HAlign = oat.HAlignRight
```

#### Per-child wrapper — AlignChild

`layout.NewAlignChild` is the cleanest way to set alignment inline, without touching the widget's own `BaseComponent`:

```go
vbox := layout.NewVBox(
    layout.NewAlignChild(saveBtn,   oat.HAlignRight,  oat.VAlignFill),
    layout.NewAlignChild(cancelBtn, oat.HAlignLeft,   oat.VAlignFill),
)
```

`AlignChild` takes precedence over both the child's own `BaseComponent` alignment and the box-wide default — it is the highest-priority override.

#### Resolution order

For each child in a VBox or HBox, the effective alignment is resolved as follows:

1. **`AlignChild` wrapper** — if the child is an `AlignChild`, its declared alignment wins.
2. **Child's own `AlignProvider`** — if the child embeds `BaseComponent` and has a non-fill value set, that value is used.
3. **Box-wide default** — the value passed to `WithHAlign` / `WithVAlign` on the box itself.
4. **`HAlignFill` / `VAlignFill`** — the fallback; full-stretch behaviour, identical to before alignment was added.

#### Example — centred button row

```go
// A row of buttons that are each pinned to their natural width
// and the whole group is centred by wrapping in a VBox set to HAlignCenter.
btnRow := layout.NewHBox(
    cancelBtn,
    layout.NewHFill().WithMaxSize(2),
    okBtn,
)

vbox := layout.NewVBox(
    layout.NewFlexChild(contentArea, 1),
    layout.NewAlignChild(btnRow, oat.HAlignCenter, oat.VAlignFill),
)
```

#### Example — mixed alignment in one box

```go
// Title centred, body fills width, action link right-aligned
vbox := layout.NewVBox(
    layout.NewAlignChild(widget.NewText("Welcome"), oat.HAlignCenter, oat.VAlignFill),
    layout.NewFlexChild(bodyContent, 1),
    layout.NewAlignChild(widget.NewText("Sign in →"), oat.HAlignRight, oat.VAlignFill),
)
```

#### Example — bottom-aligned status chip in an HBox

```go
// Icon + name fill full height; status chip is pinned to the bottom of the row
row := layout.NewHBox(
    layout.NewAlignChild(iconText, oat.HAlignFill, oat.VAlignFill),
    layout.NewFlexChild(nameText, 1),
    layout.NewAlignChild(statusChip, oat.HAlignFill, oat.VAlignBottom),
)
```

#### Combining alignment with flex

`AlignChild` is not a flex spacer on its own — wrap it with `AddFlexChild` or `NewFlexChild` when you also need it to participate in flex distribution:

```go
// Flex child that is also right-aligned within its slot
vbox.AddFlexChild(
    layout.NewAlignChild(myWidget, oat.HAlignRight, oat.VAlignFill),
    1,
)

// Or inline:
vbox := layout.NewVBox(
    layout.NewFlexChild(
        layout.NewAlignChild(myWidget, oat.HAlignRight, oat.VAlignFill),
    ),
)
```


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

```go
func (b *Border) WithRoundedCorner(rounded bool) *Border
```

Stores the rounded-corner intent in an internal field; does **not** mutate `Style.Border`. The effective corner shape is resolved at render time: `BorderSingle` ↔ `BorderRounded` is toggled based on this field.

**Incompatible styles** (`BorderDouble`, `BorderThick`, `BorderDashed`) are **silently left unchanged** in both directions — no panic is raised.

| Border style | `WithRoundedCorner(true)` |
|---|---|
| `BorderSingle` (default) | Renders with arc corners `╭╮╰╯` |
| `BorderRounded` | Already rounded — no change |
| `BorderDouble` | Silently ignored — no arc codepoints for double strokes |
| `BorderThick` | Silently ignored — no arc codepoints for heavy strokes |
| `BorderDashed` | Silently ignored — arc corners don't connect to dashed strokes |

Once called, this explicit choice overrides the theme's `RoundedCorner` setting for this `Border` — `ApplyTheme` will not overwrite it on theme switches. Calling `WithRoundedCorner(false)` explicitly opts out of rounded corners even when `theme.RoundedCorner` is `true`.

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

## FlexChild

`FlexChild` wraps any `Component` as a flex slot. It is a convenience type that lets you pass flex children directly to the variadic `NewVBox` / `NewHBox` constructors without needing a separate `AddFlexChild` call after construction.

```go
// Without FlexChild — two-step pattern:
vbox := layout.NewVBox(titleText)
vbox.AddFlexChild(bodyEditor, 1)
vbox.AddChild(btnRow)

// With FlexChild — one-liner:
vbox := layout.NewVBox(
    titleText,
    layout.NewFlexChild(bodyEditor),   // weight defaults to 1
    btnRow,
)
```

`NewFlexChild(child, weight)` — weight is variadic and defaults to `1`. The minimum effective weight is `1`.

```go
// Explicit weight:
layout.NewFlexChild(leftPanel, 1)
layout.NewFlexChild(rightPanel, 3)   // right panel gets 3× the space
```

`FlexChild` implements `oat.Layout` via `Children()`, so theme propagation and focus collection recurse into the wrapped component automatically. You can use it anywhere `AddFlexChild` would be used, including inside `HBox`.

:::tip When to use FlexChild vs AddFlexChild
They are equivalent in effect. Prefer `NewFlexChild` when building the child list inline (e.g. passing to a variadic constructor). Use `AddFlexChild` when you need to add a flex child to an already-constructed box.
:::

## AlignChild

`AlignChild` wraps a single component with explicit cross-axis alignment — the per-child counterpart to `VBox.WithHAlign` / `HBox.WithVAlign`. Use it when one child in a box needs different alignment from the rest, or when the child component does not expose its own `WithHAlign`/`WithVAlign` builder.

```go
import "github.com/antoniocali/oat-latte/layout"

// Save pinned right, Cancel pinned left — no spacers needed
vbox := layout.NewVBox(
    layout.NewAlignChild(saveBtn,   oat.HAlignRight, oat.VAlignFill),
    layout.NewAlignChild(cancelBtn, oat.HAlignLeft,  oat.VAlignFill),
)
```

### Constructor

```go
func NewAlignChild(child oat.Component, h oat.HAlign, v oat.VAlign) *AlignChild
```

- `h` — `oat.HAlignFill`, `HAlignLeft`, `HAlignCenter`, or `HAlignRight`
- `v` — `oat.VAlignFill`, `VAlignTop`, `VAlignMiddle`, or `VAlignBottom`

Pass `oat.HAlignFill` or `oat.VAlignFill` for the axis you do not want to constrain.

### Combining with flex

`AlignChild` is **not** a `FlexSpacer` — it does not claim flex space on its own. Wrap it with `AddFlexChild` or `NewFlexChild` when you also need it to participate in flex distribution:

```go
// A flex child that is also right-aligned in its slot
vbox.AddFlexChild(
    layout.NewAlignChild(headerWidget, oat.HAlignRight, oat.VAlignFill),
    1,
)

// Inline using NewFlexChild:
vbox := layout.NewVBox(
    layout.NewFlexChild(
        layout.NewAlignChild(headerWidget, oat.HAlignRight, oat.VAlignFill),
    ),
    btnRow,
)
```

### Theme and focus propagation

`AlignChild` implements `oat.Layout` via `Children()`, so theme propagation and focus collection recurse into the wrapped component automatically — exactly like `FlexChild`.


## Divider

`widget.Divider` renders a single-cell-thick rule between layout children. Place a horizontal divider in a `VBox` to separate rows; place a vertical divider in an `HBox` to separate columns.

```go
// Horizontal rule between two sections
vbox := layout.NewVBox(
    widget.NewText("Section A"),
    widget.NewHDivider(),
    widget.NewText("Section B"),
)

// Vertical rule between two panels
hbox := layout.NewHBox(
    layout.NewFlexChild(leftPanel, 1),
    widget.NewVDivider(),
    layout.NewFlexChild(rightPanel, 2),
)
```

The divider fills the full allocated span by default. Use `WithMaxSize` / `WithMaxSizeV` with a `DividerSize` to limit it and anchor it within the available space:

```go
// 60% wide, centred horizontally
widget.NewHDivider().WithMaxSize(widget.DividerPercent(60), oat.AnchorCenter)

// 8 cells tall, centred vertically
widget.NewVDivider().WithMaxSizeV(widget.DividerFixed(8), oat.VAnchorMiddle)
```

See the [Divider widget page](./widgets/divider.md) for the full API reference.
