---
sidebar_position: 3
title: EditText
description: Single-line and multi-line text input fields.
---

# EditText

`widget.EditText` is a focusable text input. It comes in two modes:

- **Single-line** — `widget.NewEditText()` — horizontal scroll when content exceeds width
- **Multi-line** — `widget.NewMultiLineEditText()` — vertical scroll; `Enter` inserts a line break

## Basic usage

```go
// Single-line
input := widget.NewEditText().
    WithPlaceholder("Enter a value…")

// Multi-line
body := widget.NewMultiLineEditText().
    WithPlaceholder("Write here…")
```

## Builder options

| Method | Description |
|---|---|
| `WithStyle(s latte.Style)` | Override the display style (border, colours, padding) |
| `WithID(id string)` | Set a stable identifier for `Canvas.GetValue(id)` |
| `WithPlaceholder(text string)` | Muted text shown when the field is empty |
| `WithHint(text string)` | Persistent muted label rendered above the content area |
| `WithMaxLength(n int)` | Maximum character count (single-line only; 0 = unlimited) |
| `WithOnChange(fn func(string))` | Callback fired on every keystroke |
| `WithOnSave(fn func(string))` | Callback fired when the user presses `^S` |
| `WithOnCancel(fn func())` | Callback fired when the user presses `^G` |

## Hint label

`WithHint` renders a persistent muted label directly above the content area. Unlike a placeholder, it is always visible and serves as an inline field label. Use it instead of a separate `Text` widget.

```go
input := widget.NewEditText().
    WithHint("Email address")
```

## Borderless fields

Use `WithStyle` with `latte.BorderExplicitNone` to suppress the border entirely. Combined with `WithHint`, this produces a clean, label-only field well-suited for editor panels:

```go
titleInput := widget.NewEditText().
    WithStyle(latte.Style{Border: latte.BorderExplicitNone}).
    WithHint("Title").
    WithPlaceholder("Untitled…")

bodyInput := widget.NewMultiLineEditText().
    WithStyle(latte.Style{Border: latte.BorderExplicitNone}).
    WithHint("Body").
    WithPlaceholder("Write here…")
```

## Callbacks

```go
input.WithOnChange(func(text string) {
    // fired on every keystroke
})

input.WithOnSave(func(text string) {
    // fired when the user presses ^S (or Enter on a single-line field)
})

input.WithOnCancel(func() {
    // fired when the user presses ^G
})
```

When `WithOnSave` is registered, `^S Save` is advertised in the `StatusBar` automatically.

## Reading and setting text

```go
input.SetText("initial value")
current := input.GetText()
```

`EditText` also implements `oat.ValueGetter`: `Canvas.GetValue(id)` returns the current text as a `string`.

## Max length

```go
input := widget.NewEditText().WithMaxLength(64)
```

## Built-in key bindings

| Key | Action |
|---|---|
| `^S` | Save — fires `onSave` callback |
| `^G` | Cancel — fires `onCancel` callback |
| `^K` | Kill to end of line |
| `^U` | Kill from start of line to cursor |
| `^A` / `Home` | Move cursor to start of line |
| `^E` / `End` | Move cursor to end of line |
| `←` / `→` | Move cursor left / right |
| `↑` / `↓` | Move cursor up / down (multi-line only) |
| `Enter` | Insert line break (multi-line) or fire `onSave` (single-line, if set) |
| `Backspace` | Delete character before cursor |
| `Delete` | Delete character after cursor |
