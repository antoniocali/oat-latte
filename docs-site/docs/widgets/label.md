---
sidebar_position: 6
title: Label
description: Inline tag/chip badge row.
---

# Label

`widget.Label` renders a horizontal row of inline badge chips, separated by a `·` character. It is a non-focusable display component useful for showing tags or categories alongside list items or editor fields.

## Basic usage

```go
lbl := widget.NewLabel([]string{"go", "tui", "terminal"})
```

Renders as: `go · tui · terminal`

## Builder options

| Method | Description |
|---|---|
| `WithStyle(s latte.Style)` | Override the chip display style |
| `WithID(id string)` | Set a stable identifier for `Canvas.GetValue(id)` |
| `WithSeparator(r rune)` | Override the separator rune (default `·`) |
| `WithHighlight(bool)` | Control whether chips render with a background colour fill (default `true`) |

## Custom separator

```go
lbl := widget.NewLabel([]string{"backend", "api"}).
    WithSeparator('|')
// Renders as:  backend | api
```

## Updating labels

```go
lbl.SetLabels([]string{"work", "urgent"})
lbl.SetLabels(nil) // clears all chips
```

## Reading labels

```go
tags := lbl.GetLabels() // []string
```

`Label` also implements `oat.ValueGetter`: `Canvas.GetValue(id)` returns all chips joined by the separator as a `string`.

## Typical use

```go
// In a note editor, below the tags input field:
tagsLabel := widget.NewLabel(nil)

tagsInput.WithOnChange(func(text string) {
    parts := strings.Split(text, ",")
    var tags []string
    for _, t := range parts {
        if t = strings.TrimSpace(t); t != "" {
            tags = append(tags, t)
        }
    }
    tagsLabel.SetLabels(tags)
})
```

Each chip is styled with the theme's `Tag` token. The separator uses the theme's `Muted` token.

## WithHighlight

`WithHighlight(false)` strips the chip background colour while keeping the foreground colour and text attributes:

```go
// Default — filled chip background (uses Tag token BG).
lbl := widget.NewLabel([]string{"go", "tui"})

// No background fill — FG colour and bold/italic are preserved.
lbl := widget.NewLabel([]string{"go", "tui"}).WithHighlight(false)
```

Useful for minimal or transparent UIs where a coloured badge background would clash with the surrounding panel.
