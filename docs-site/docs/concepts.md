---
sidebar_position: 3
title: Core Concepts
description: The three building blocks of every oat-latte application — Canvas, Layouts, and Widgets.
---

# Core Concepts

Every oat-latte application is built from three kinds of thing: a **Canvas**, **Layouts**, and **Widgets**. Understanding what each one does — and how they relate — is everything you need to build UIs with oat-latte.

---

## Canvas

The Canvas is the application root. It owns the terminal screen, runs the event loop, manages keyboard focus, and handles overlays (dialogs). You create exactly one Canvas per application.

```go
app := oat.NewCanvas(
    oat.WithTheme(latte.ThemeDark),
    oat.WithBody(body),           // the root layout
)
if err := app.Run(); err != nil {
    log.Fatal(err)
}
```

`Run()` takes over the terminal and blocks until the user quits. The Canvas divides the screen into three vertical regions: an optional **header**, the **body** (fills all remaining space), and an optional **footer** (used by the status bar).

### Useful Canvas options

| Option | What it does |
|---|---|
| `WithTheme(t)` | Apply a colour theme to the entire tree |
| `WithBody(c)` | Set the root layout component |
| `WithHeader(c)` | Set an optional header component |
| `WithAutoStatusBar(bar)` | Mount a status bar in the footer |
| `WithPrimary(f)` | Which widget receives focus first |
| `WithNotificationManager(n)` | Wire up toast notifications |
| `WithGlobalKeyBinding(b)` | Register an app-wide keyboard shortcut |

### Useful Canvas methods

| Method | What it does |
|---|---|
| `Quit()` | Exit the event loop gracefully |
| `SetTheme(t)` | Switch themes at runtime |
| `GetTheme()` | Return a pointer to the active theme (`nil` if none set) |
| `ShowDialog(d)` | Push a modal dialog; focus moves into it |
| `HideDialog()` | Dismiss the topmost dialog |
| `FocusByRef(f)` | Jump focus directly to a specific widget |

---

## Layouts

Layouts are containers. They hold other components — other layouts, or widgets — and decide how to position and size them. You nest layouts to build any UI structure.

oat-latte ships six layout primitives:

### VBox — vertical stack

Arranges children top-to-bottom. Use `AddFlexChild` to let a child stretch and fill remaining space.

```go
vbox := layout.NewVBox(
    widget.NewText("Label"),
    widget.NewEditText().WithHint("Name"),
)
vbox.AddFlexChild(widget.NewText(longContent), 1) // stretches to fill
```

Use `WithHAlign` to align children horizontally within their slot (instead of filling the full width):

```go
// Every child centred horizontally
vbox := layout.NewVBox(title, body, footer).WithHAlign(oat.HAlignCenter)
```

### HBox — horizontal row

Arranges children left-to-right. Same flex system as VBox.

```go
hbox := layout.NewHBox()
hbox.AddFlexChild(listPanel, 1)
hbox.AddFlexChild(detailPanel, 3)  // 3× wider than the list
```

Use `WithVAlign` to align children vertically within their slot:

```go
// Status badges pinned to the bottom of each column
hbox := layout.NewHBox(iconText, nameText, statusChip).WithVAlign(oat.VAlignBottom)
```

### Border — framed panel

Wraps any component in a titled border. Border panels automatically highlight their frame when a descendant is focused.

```go
panel := layout.NewBorder(innerComponent).
    WithTitle("My Panel").
    WithRoundedCorner(true)
```

### Padding

Adds whitespace around a component without drawing a border.

```go
padded := layout.NewPaddingUniform(child, 1)  // 1 cell on all sides
```

### Grid

Places children into a fixed row × column grid with optional spanning.

```go
g := layout.NewGrid(2, 3) // 2 rows, 3 cols
g.AddChildAt(widget, 0, 0, 1, 2) // row 0, col 0, span 1 row × 2 cols
```

### Spacers — VFill and HFill

Empty spacers that consume leftover space. Use them to push widgets to the edges of a box.

```go
btnRow := layout.NewHBox()
btnRow.AddChild(layout.NewHFill())  // pushes buttons to the right
btnRow.AddChild(cancelBtn)
btnRow.AddChild(layout.NewHFill().WithMaxSize(2)) // fixed 2-cell gap
btnRow.AddChild(okBtn)
```

### AlignChild — per-child alignment override

`AlignChild` wraps a single component and overrides its alignment within the parent box — useful when one child needs different alignment from the rest, or when you want alignment set inline rather than on the widget itself.

```go
// Save pinned right, Cancel pinned left — no spacer widgets needed
vbox := layout.NewVBox(
    layout.NewAlignChild(saveBtn,   oat.HAlignRight, oat.VAlignFill),
    layout.NewAlignChild(cancelBtn, oat.HAlignLeft,  oat.VAlignFill),
)
```

See the [Layout reference](./layout.md#alignchild) for the full API and more examples.

### Composing layouts

Layouts nest freely. A typical two-panel app:

```go
list   := widget.NewList(items)
detail := widget.NewText("")

list.WithOnCursorChange(func(_ int, item widget.ListItem) {
    detail.SetText(fmt.Sprint(item.Value))
})

body := layout.NewHBox()
body.AddFlexChild(layout.NewBorder(list).WithTitle("Items"),   1)
body.AddFlexChild(layout.NewBorder(detail).WithTitle("Detail"), 3)
```

---

## Widgets

Widgets are the leaves of the tree — the elements users actually see and interact with. oat-latte ships a full set out of the box. See the **Widgets** section in the sidebar for the full reference for each one.

### Display widgets

| Widget | What it renders |
|---|---|
| `widget.NewText(s)` | Static text with word-wrap and scroll |
| `widget.NewTitle(s)` | Bold heading with an optional rule beneath |
| `widget.NewLabel(tags)` | Row of inline badge chips |
| `widget.NewProgressBar()` | Horizontal fill bar with optional percent label |
| `widget.NewDivider(axis)` | Horizontal or vertical rule |

### Interactive widgets

| Widget | What it does |
|---|---|
| `widget.NewButton(label, fn)` | Pressable button; fires `fn` on Enter or Space |
| `widget.NewEditText()` | Single-line text input |
| `widget.NewMultiLineEditText()` | Multi-line text input |
| `widget.NewCheckBox(label)` | Boolean toggle |
| `widget.NewList(items)` | Scrollable, selectable list |

### App-level widgets

| Widget | What it does |
|---|---|
| `widget.NewDialog(title)` | Modal overlay with scrim backdrop |
| `widget.NewStatusBar()` | Footer bar that shows key-binding hints |
| `widget.NewNotificationManager()` | Toast notifications pushed from anywhere |

### Wiring widgets together

Widgets communicate through callbacks, not through a central model. Set a callback on one widget to update another:

```go
input := widget.NewEditText().WithHint("Search")
list  := widget.NewList(allItems)

input.WithOnChange(func(text string) {
    list.SetItems(filter(allItems, text))
})
```

### Reading widget values

Every interactive widget implements `ValueGetter`. The Canvas can look up any widget's current value by its ID:

```go
nameInput := widget.NewEditText().WithID("name")
// ... later ...
val, _ := app.GetValue("name") // returns the current text as interface{}
```

---

## Putting it together

A complete application assembles these three layers:

```go
func main() {
    // 1. Widgets
    input  := widget.NewEditText().WithHint("Name").WithID("name")
    result := widget.NewText("—")
    btn    := widget.NewButton("Save", func() {
        result.SetText("Saved: " + input.GetText())
    })

    // 2. Layouts
    body := layout.NewBorder(
        layout.NewPaddingUniform(
            layout.NewVBox(input, btn, result),
            1,
        ),
    ).WithTitle("New item")

    // 3. Canvas
    app := oat.NewCanvas(
        oat.WithTheme(latte.ThemeDark),
        oat.WithBody(body),
        oat.WithPrimary(input),
    )
    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

Once you are comfortable with these three layers, the [Tutorial](/docs/tutorial) walks you through building a full task-list application step by step.
