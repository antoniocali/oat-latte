---
sidebar_position: 5
title: ComponentList
description: Scrollable, selectable list widget whose rows are arbitrary components.
---

# ComponentList

`widget.ComponentList` is the component-row counterpart of [`List`](./list.md). Instead of a plain `Label` string, each row renders an arbitrary `Component` — any widget or layout — so rows can contain rich, multi-column content like `HBox(Text, Flex(Text), Text)`.

Each item still carries a `Value interface{}` field, so the caller can store a row identifier (e.g. a record ID) and act on it from callbacks.

`ComponentList` implements `oat.Layout` so the theme tree-walker and the focus collector recurse into every row component automatically — no extra wiring needed.

## Basic usage

```go
items := []widget.ComponentListItem{
    {Component: layout.NewHBox(widget.NewText("Alice"),  widget.NewText("active")),   Value: 1},
    {Component: layout.NewHBox(widget.NewText("Bob"),    widget.NewText("inactive")), Value: 2},
    {Component: layout.NewHBox(widget.NewText("Charlie"),widget.NewText("active")),   Value: 3},
}

list := widget.NewComponentList(items).WithID("people-list")
```

## Building rich rows

The most common pattern is a factory function that produces each row component:

```go
makeRow := func(name, role, status string, id int) widget.ComponentListItem {
    row := layout.NewHBox(
        widget.NewText(name),
        layout.NewFlexChild(widget.NewText(role), 1),  // flex → fills available width
        widget.NewText(status),
    )
    return widget.ComponentListItem{Component: row, Value: id}
}

items := []widget.ComponentListItem{
    makeRow("Alice",   "Backend engineer",  "active",   1),
    makeRow("Bob",     "Frontend engineer", "inactive", 2),
    makeRow("Charlie", "DevOps",            "active",   3),
}

list := widget.NewComponentList(items).
    WithID("people-list").
    WithOnSelect(func(idx int, item widget.ComponentListItem) {
        id := item.Value.(int)
        // open detail view for this record
    }).
    WithOnCursorChange(func(idx int, item widget.ComponentListItem) {
        // update live preview
    })
```

## Builder options

| Method | Description |
|---|---|
| `WithStyle(s latte.Style)` | Override the display style (border, colours, padding) |
| `WithID(id string)` | Set a stable identifier for `Canvas.GetValue(id)` |
| `WithSelectedStyle(s latte.Style)` | Override the style for the highlighted row |
| `WithHighlight(enabled bool)` | Fill selected row with `selectedStyle` (default `true`) |
| `WithCursor(cursor string)` | Gutter character next to the selected row (default `>`) |
| `WithOnSelect(fn func(int, ComponentListItem))` | Callback fired on `Enter` |
| `WithOnDelete(fn func(int, ComponentListItem))` | Callback fired on `Delete` |
| `WithOnCursorChange(fn func(int, ComponentListItem))` | Callback fired on every cursor move |

## Callbacks

```go
// Fired when the user presses Enter.
list.WithOnSelect(func(idx int, item widget.ComponentListItem) {
    id := item.Value.(int)
    openRecord(id)
})

// Fired when the user presses Delete.
list.WithOnDelete(func(idx int, item widget.ComponentListItem) {
    confirmAndRemove(idx)
})

// Fired on every cursor move — use for live preview panels.
list.WithOnCursorChange(func(idx int, item widget.ComponentListItem) {
    detailView.SetText(fmt.Sprint(item.Value))
})
```

`WithOnCursorChange` is the key to two-panel "master/detail" layouts: as the user navigates the list, the detail panel updates in real time without requiring `Enter`.

## Cursor and highlight

```go
// Custom cursor glyph.
list.WithCursor("▶")   // or "→", "•", "❯", "*"

// Hide cursor entirely.
list.WithCursor("")

// Disable row background fill (cursor-only indication).
list.WithHighlight(false)
```

## Updating items

Call `list.SetItems(newItems)` to replace the full item slice at any time.

`SetItems` preserves the cursor position if the new slice is long enough; otherwise it clamps to the last item.

## Reading the selection

```go
item, ok := list.SelectedItem()   // (widget.ComponentListItem, bool)
idx       := list.SelectedIndex() // int
```

`ComponentList` also implements `oat.ValueGetter`: `Canvas.GetValue(id)` returns the `Value` field of the currently selected item.

## Variable row heights

Unlike `List`, rows in a `ComponentList` can span more than one terminal line. Each row's component is measured with `Measure` to determine its height; the total list height is the sum of all row heights (plus border/padding). Scroll is tracked by item index so the viewport always shows complete rows — partial rows are never rendered.

```go
// A row that is 3 cells tall because it contains multi-line content.
makeRow := func(title, body string, id int) widget.ComponentListItem {
    col := layout.NewVBox(
        widget.NewText(title),
        widget.NewText(body),
        widget.NewHDivider(),
    )
    return widget.ComponentListItem{Component: col, Value: id}
}
```

## Two-panel example

```go
list   := widget.NewComponentList(makeItems()).WithID("records")
detail := widget.NewText("")

list.WithOnCursorChange(func(_ int, item widget.ComponentListItem) {
    if r, ok := item.Value.(Record); ok {
        detail.SetText(r.Description)
    }
})

body := layout.NewHBox()
body.AddFlexChild(layout.NewBorder(list).WithTitle("Records"),  1)
body.AddFlexChild(layout.NewBorder(detail).WithTitle("Detail"), 3)
```

## Theme

`ComponentList` uses the same theme tokens as `List`:

| Token | Applied to |
|---|---|
| `t.Text` | Base list style |
| `t.ListSelected` | Selected row highlight |
| `t.FocusBorder` | Border colour when focused |

Override any of these with `WithStyle` / `WithSelectedStyle` after construction.

## Key bindings

| Key | Action |
|---|---|
| `↑` / `↓` | Move cursor up / down |
| `Home` / `^A` | Jump to first item |
| `End` / `^E` | Jump to last item |
| `Enter` | Fire `onSelect` callback |
| `Delete` | Fire `onDelete` callback (only advertised when set) |
| `Tab` / `Shift+Tab` | Move focus to next / previous widget |

## Example app — People Directory

`cmd/example/people` is a runnable reference app that demonstrates `ComponentList` end-to-end:

- Multi-column rows: **name** (fixed width) · **role** (flex) · **status** (coloured)
- Live-preview detail panel — updates as the cursor moves
- "New Person" dialog with two `EditText` inputs
- Delete with a confirmation dialog
- Theme cycling with `^T`

```
╭─ People ─────────────────────────────╮ ╭─ Detail ─────────────────────╮
│ > Alice   Backend engineer   active   │ │ Name:   Alice                │
│   Bob     Frontend engineer  inactive │ │ Role:   Backend engineer     │
│   Charlie DevOps             active   │ │ Status: active               │
│   Diana   Product manager    active   │ │                              │
│   Eve     Data scientist     inactive │ │                              │
╰───────────────────────────────────────╯ ╰──────────────────────────────╯
  ↑↓ Move   n New   Del Delete   ^T Theme   Tab Next   Esc Quit
```

Run it from the repo root:

```sh
go run ./cmd/example/people
```

### How the rows are built

```go
func makeRow(p Person) widget.ComponentListItem {
    statusStyle := latte.Style{FG: latte.ColorGreen}
    if p.Status != "active" {
        statusStyle = latte.Style{FG: latte.ColorBrightBlack}
    }

    row := layout.NewHBox(
        widget.NewText(fmt.Sprintf("%-16s", p.Name)),
        layout.NewFlexChild(widget.NewText(p.Role), 1),
        widget.NewText(p.Status).WithStyle(statusStyle),
    )
    return widget.ComponentListItem{Component: row, Value: p.ID}
}
```

The `FlexChild` on the role column ensures the role text stretches to fill the available width, keeping the status tag flush-right regardless of terminal width.

### Live-preview wiring

```go
cl.WithOnCursorChange(func(_ int, item widget.ComponentListItem) {
    if id, ok := item.Value.(int); ok {
        if person, _, ok := personByID(id); ok {
            detail.SetText(fmt.Sprintf(
                "Name:   %s\nRole:   %s\nStatus: %s",
                person.Name, person.Role, person.Status,
            ))
        }
    }
})
```

### Extending with a proxy

A thin `listProxy` intercepts the `n` key to open the "New Person" dialog, without modifying `ComponentList` itself:

```go
type listProxy struct {
    *widget.ComponentList
    app *App
}

func (p *listProxy) HandleKey(ev *oat.KeyEvent) bool {
    if ev.Key() == tcell.KeyRune && ev.Rune() == 'n' {
        p.app.showNewDialog()
        return true
    }
    return p.ComponentList.HandleKey(ev)
}

func (p *listProxy) KeyBindings() []oat.KeyBinding {
    extra := []oat.KeyBinding{
        {Key: tcell.KeyRune, Rune: 'n', Label: "n", Description: "New person"},
    }
    return append(extra, p.ComponentList.KeyBindings()...)
}
```

---

## Choosing between List and ComponentList

| | List | ComponentList |
|---|---|---|
| Row content | Plain label string | Any `Component` |
| Row height | Always 1 cell | Variable (measured per row) |
| Theme propagation | n/a | Automatic via `Layout.Children()` |
| Use when | Simple text items | Multi-column, multi-line, or styled rows |
