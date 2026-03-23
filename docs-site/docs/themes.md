---
sidebar_position: 6
title: Themes & Styles
description: How the Style struct, Style.Merge cascade, and built-in themes work in oat-latte.
---

# Themes & Styles

## The Style struct

`latte.Style` is the complete visual description of any component:

```go
type Style struct {
    FG, BG          latte.Color
    Bold, Italic, Underline, Blink, Reverse bool
    Padding, Margin latte.Insets
    Border          latte.BorderStyle
    BorderFG, BorderBG latte.Color
    TextAlign       latte.Alignment
}
```

The **zero value** of `Style` means "inherit / use default". A field only overrides the default when it is non-zero. This makes styles composable without touching fields you do not care about.

## Colors

```go
latte.RGB(30, 120, 255)    // true-color (24-bit)
latte.Hex("#1E78FF")        // true-color from a hex string
latte.ColorCyan             // ANSI-16 named colors
latte.ColorBrightWhite      // ANSI-16 bright variants
```

True-color requires a modern terminal (e.g. iTerm2, kitty, Windows Terminal, most Linux terminals). `latte.ThemeDefault` uses only ANSI-16 and works everywhere.

## Borders

| Constant | Runes | Notes |
|---|---|---|
| `latte.BorderNone` | — | Zero value; inherits from theme |
| `latte.BorderExplicitNone` | — | Actively suppresses the border |
| `latte.BorderSingle` | `┌─┐│└─┘` | |
| `latte.BorderRounded` | `╭─╮│╰─╯` | Default for dialogs |
| `latte.BorderDouble` | `╔═╗║╚═╝` | |
| `latte.BorderThick` | `┏━┓┃┗━┛` | |
| `latte.BorderDashed` | `┌╌┐╎└╌┘` | |

Use `latte.BorderExplicitNone` to suppress a border that a theme would otherwise add:

```go
// No box drawn, even if the theme sets a border on Input.
input := widget.NewEditText(latte.Style{Border: latte.BorderExplicitNone})
```

## Style.Merge

```go
result := base.Merge(override)
```

Non-zero fields from `override` replace the corresponding fields in `base`. Zero fields in `override` leave `base` untouched. This is how theme tokens cascade without clobbering explicit per-widget settings.

:::warning
In `ApplyTheme`, always use `Merge` — never assign the theme token directly:

```go
// Correct — theme is the base; per-widget overrides survive.
func (w *MyWidget) ApplyTheme(t latte.Theme) {
    w.Style = t.Input.Merge(w.Style)
}

// Wrong — overwrites BorderExplicitNone and any other caller-set field.
func (w *MyWidget) ApplyTheme(t latte.Theme) {
    w.Style = t.Input
}
```
:::

## Built-in themes

Apply a theme once at construction time with `oat.WithTheme(t)`. The canvas walks the entire component tree and calls `ApplyTheme` on every node that implements `ThemeReceiver`.

| Theme | Palette | Terminal requirement |
|---|---|---|
| `latte.ThemeDefault` | ANSI-16 | Any terminal |
| `latte.ThemeDark` | True-color, deep navy / blue-cyan | True-color terminal |
| `latte.ThemeLight` | True-color, warm off-white / indigo | True-color terminal |
| `latte.ThemeDracula` | True-color, Dracula palette | True-color terminal |
| `latte.ThemeNord` | True-color, Nord arctic palette | True-color terminal |

```go
app := oat.NewCanvas(
    oat.WithTheme(latte.ThemeDark),
    oat.WithBody(body),
)
```

## Semantic theme tokens

Themes expose named style tokens — roles rather than components — so the same theme drives the entire UI consistently:

| Token | Used by |
|---|---|
| `Canvas` | Full-screen background |
| `Text`, `Muted` | Body copy, secondary text |
| `Accent`, `Success`, `Warning`, `Error` | State colors |
| `Panel`, `PanelTitle` | `layout.Border` containers |
| `Input`, `InputFocus` | `widget.EditText` |
| `ListSelected` | `widget.List` selected row |
| `Button`, `ButtonFocus` | `widget.Button` |
| `CheckBox`, `CheckBoxFocus` | `widget.CheckBox` |
| `Header`, `Footer` | Canvas header / footer |
| `FocusBorder` | `layout.Border` when a descendant is focused |
| `Dialog`, `DialogTitle`, `Scrim` | `widget.Dialog` and its backdrop |
| `Tag` | `widget.Label` chips |
| `NotificationInfo/Success/Warning/Error` | `widget.NotificationManager` banners |

## Applying a theme to a custom widget

```go
func (w *MyWidget) ApplyTheme(t latte.Theme) {
    // Pick the token that best describes your widget's role.
    w.Style      = t.Input.Merge(w.Style)
    w.FocusStyle = t.InputFocus.Merge(w.FocusStyle)
}
```

The framework calls `ApplyTheme` automatically when the canvas is constructed with `WithTheme`. You do not need to call it yourself.
