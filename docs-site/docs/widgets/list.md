---
sidebar_position: 4
title: List
description: Scrollable, selectable list widget.
---

# List

`widget.List` is a vertically scrollable list of items. Each item has a `Label` (displayed string) and a `Value` (arbitrary `interface{}`). A `>` cursor marks the selected row and the selected row is highlighted by default.

## Basic usage

```go
items := []widget.ListItem{
    {Label: "First item",  Value: 1},
    {Label: "Second item", Value: 2},
    {Label: "Third item",  Value: "three"},
}

list := widget.NewList(items).WithID("my-list")
```

## Builder options

| Method | Description |
|---|---|
| `WithStyle(s latte.Style)` | Override the display style (border, colours, padding) |
| `WithID(id string)` | Set a stable identifier for `Canvas.GetValue(id)` |
| `WithSelectedStyle(s latte.Style)` | Override the style for the highlighted row |
| `WithHighlight(enabled bool)` | Fill selected row with `selectedStyle` (default `true`) |
| `WithCursor(cursor string)` | Gutter character next to the selected row (default `>`) |
| `WithOnSelect(fn func(int, ListItem))` | Callback fired on `Enter` |
| `WithOnDelete(fn func(int, ListItem))` | Callback fired on `Delete` |
| `WithOnCursorChange(fn func(int, ListItem))` | Callback fired on every cursor move |

## Callbacks

```go
// Fired when the user presses Enter.
list.WithOnSelect(func(idx int, item widget.ListItem) {
    fmt.Printf("Selected: %v\n", item.Value)
})

// Fired when the user presses Delete.
list.WithOnDelete(func(idx int, item widget.ListItem) {
    // confirm and remove
})

// Fired on every cursor move — use for live preview panels.
list.WithOnCursorChange(func(idx int, item widget.ListItem) {
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

```go
list.SetItems(newItems)
```

`SetItems` preserves the cursor position if the new slice is long enough; otherwise it clamps to the last item.

## Reading the selection

```go
item, ok := list.SelectedItem()  // (widget.ListItem, bool)
idx       := list.SelectedIndex() // int
```

`List` also implements `oat.ValueGetter`: `Canvas.GetValue(id)` returns the `Value` field of the currently selected item.

## Live-preview example

```go
list   := widget.NewList(notes).WithID("notes")
detail := widget.NewText("")

list.WithOnCursorChange(func(_ int, item widget.ListItem) {
    if note, ok := item.Value.(Note); ok {
        detail.SetText(note.Body)
    }
})

body := layout.NewHBox()
body.AddFlexChild(layout.NewBorder(list).WithTitle("Notes"),  1)
body.AddFlexChild(layout.NewBorder(detail).WithTitle("Detail"), 3)
```

## Key bindings

| Key | Action |
|---|---|
| `↑` / `↓` | Move cursor up / down |
| `Home` / `^A` | Jump to first item |
| `End` / `^E` | Jump to last item |
| `Enter` | Fire `onSelect` callback |
| `Delete` | Fire `onDelete` callback (only advertised when set) |
| `Tab` / `Shift+Tab` | Move focus to next / previous widget |
