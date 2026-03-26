---
sidebar_position: 6
title: Themes & Styles
description: How the Style struct, Style.Merge cascade, built-in themes, and runtime theme switching work in oat-latte.
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

## Named color palette

`latte/colors.go` provides ~120 named `Color` constants so you can reference semantic names instead of raw hex strings. All constants are `latte.Color` values produced via `latte.RGB` and work with any API that accepts a `latte.Color`.

### Utility scales

Tailwind-style shade scales (50–950): `Slate`, `Zinc`, `Stone`, `Sky`, `Blue`, `Indigo`, `Cornflower`, `Cyan`, `Teal`, `Emerald`, `Green`, `Lime`, `Yellow`, `Amber`, `Orange`, `Red`, `Rose`, `Pink`, `Violet`, `Purple`, `Fuchsia`.

```go
latte.Pink500       // RGB(236, 72, 153)
latte.Slate800      // RGB(30, 41, 59)
latte.Emerald400    // RGB(52, 211, 153)
latte.Cornflower400 // RGB(100, 149, 237) — classic TUI focus blue
```

### Design-system palettes

Named constants extracted directly from the four true-color built-in themes, grouped by theme family:

| Family | Example constants |
|---|---|
| `Dark*` | `DarkBg`, `DarkAccent`, `DarkMuted`, `DarkSuccess`, `DarkWarning`, `DarkError`, `DarkBgElevated`, `DarkBgScrim`, `DarkBorder` |
| `Light*` | `LightBg`, `LightAccent`, `LightText`, `LightMuted`, `LightSuccess`, `LightWarning`, `LightError`, `LightBgElevated`, `LightBgScrim`, `LightBorder` |
| `Dracula*` | `DraculaBg`, `DraculaPurple`, `DraculaCyan`, `DraculaGreen`, `DraculaFg`, `DraculaComment`, `DraculaSelection`, `DraculaOrange`, `DraculaRed`, `DraculaYellow`, `DraculaPink` |
| `Nord0`–`Nord15` | `NordBg`, `NordBgElevated`, `NordBgScrim` |

```go
// Derive a custom theme using named constants instead of raw hex literals.
myTheme := latte.ThemeDark.
    WithAccent(latte.Style{FG: latte.Pink500}).
    WithFocusBorder(latte.Pink500).
    WithName("dark-pink")

// Mix palette families — e.g. Nord background with Dracula accents.
hybrid := latte.ThemeNord.
    WithAccent(latte.Style{FG: latte.DraculaPurple}).
    WithFocusBorder(latte.DraculaCyan).
    WithName("nord-dracula")
```

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
input := widget.NewEditText().
    WithStyle(latte.Style{Border: latte.BorderExplicitNone})
```

## Style.Merge

Non-zero fields from `override` replace the corresponding fields in `base`. Zero fields in `override` leave `base` untouched. Calling `base.Merge(override)` is how theme tokens cascade without clobbering explicit per-widget settings.

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

Five themes ship out of the box:

| Theme | Palette | Terminal requirement |
|---|---|---|
| `latte.ThemeDefault` | ANSI-16 | Any terminal |
| `latte.ThemeDark` | True-color, deep navy / blue-cyan | True-color terminal |
| `latte.ThemeLight` | True-color, warm off-white / indigo | True-color terminal |
| `latte.ThemeDracula` | True-color, Dracula palette | True-color terminal |
| `latte.ThemeNord` | True-color, Nord arctic palette | True-color terminal |

### Applying a theme at construction

Pass a theme via `oat.WithTheme(t)`. The canvas walks the entire component tree and calls `ApplyTheme` on every node that implements `ThemeReceiver`.

```go
app := oat.NewCanvas(
    oat.WithTheme(latte.ThemeDark),
    oat.WithBody(body),
)
```

### Switching themes at runtime

Call `app.SetTheme(t)` to replace the active theme at any time — for example from a key binding. The new theme is immediately re-applied to the entire tree, including all mounted overlays and persistent overlays. The canvas background is also updated.

```go
themes := []latte.Theme{
    latte.ThemeDark,
    latte.ThemeLight,
    latte.ThemeDracula,
    latte.ThemeNord,
}
current := 0

app := oat.NewCanvas(
    oat.WithTheme(themes[current]),
    oat.WithBody(body),
    oat.WithGlobalKeyBinding(oat.KeyBinding{
        Key:         tcell.KeyCtrlT,
        Label:       "^T",
        Description: "Toggle theme",
        Handler: func() {
            current = (current + 1) % len(themes)
            app.SetTheme(themes[current])
        },
    }),
)
```

`SetTheme` is safe to call from any key-event callback — it runs on the main goroutine and the event loop re-renders on the next tick automatically.

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

Pick the token that best describes your widget's role and apply it via `Merge`:

```go
func (w *MyWidget) ApplyTheme(t latte.Theme) {
    // Pick the token that best describes your widget's role.
    w.Style      = t.Input.Merge(w.Style)
    w.FocusStyle = t.InputFocus.Merge(w.FocusStyle)
}
```

The framework calls `ApplyTheme` automatically when:
- the canvas starts via `WithTheme`, and
- `app.SetTheme(t)` is called at runtime.

You do not need to call it yourself.

## Deriving a custom theme

Every built-in theme exposes a `With<Token>` builder method for each of its fields. All methods return a new `Theme` by value — the originals are never mutated.

**Style-typed tokens** accept a `latte.Style` and use `Style.Merge` semantics: only the non-zero fields you supply are overridden; the rest of the token is left exactly as the base theme defined it.

```go
// Start from a built-in theme and make targeted tweaks.
myTheme := latte.ThemeNord.
    WithAccent(latte.Style{FG: latte.Hex("#ff69b4")}).   // swap accent colour only
    WithFocusBorder(latte.Hex("#ff69b4")).                // keep Nord's FG, just change focus ring
    WithName("nord-pink")
```

### Removing borders globally

Pass `latte.Style{Border: latte.BorderExplicitNone}` to the tokens that carry borders. `BorderExplicitNone` propagates through `Merge` and actively suppresses any border that the base theme would draw:

```go
borderless := latte.ThemeNord.
    WithPanel(latte.Style{Border: latte.BorderExplicitNone}).
    WithInput(latte.Style{Border: latte.BorderExplicitNone}).
    WithButton(latte.Style{Border: latte.BorderExplicitNone}).
    WithDialog(latte.Style{Border: latte.BorderExplicitNone}).
    WithName("nord-borderless")

app := oat.NewCanvas(oat.WithTheme(borderless), oat.WithBody(body))
```

### Changing the background

```go
deepBlack := latte.ThemeDark.
    WithCanvas(latte.Style{BG: latte.Hex("#000000")}).
    WithName("dark-deep")
```

### Full list of builder methods

| Method | Token type | Notes |
|---|---|---|
| `WithName(string)` | — | Sets the theme name |
| `WithCanvas(Style)` | `Style` | Full-screen background |
| `WithText(Style)` | `Style` | Body text |
| `WithMuted(Style)` | `Style` | Secondary / de-emphasised text |
| `WithAccent(Style)` | `Style` | Primary interactive colour |
| `WithSuccess(Style)` | `Style` | Positive state |
| `WithWarning(Style)` | `Style` | Cautionary state |
| `WithError(Style)` | `Style` | Destructive / critical state |
| `WithPanel(Style)` | `Style` | `layout.Border` containers |
| `WithPanelTitle(Style)` | `Style` | Title text in panel borders |
| `WithInput(Style)` | `Style` | `widget.EditText` base |
| `WithInputFocus(Style)` | `Style` | `widget.EditText` focused overlay |
| `WithListSelected(Style)` | `Style` | `widget.List` selected row |
| `WithButton(Style)` | `Style` | `widget.Button` base |
| `WithButtonFocus(Style)` | `Style` | `widget.Button` focused overlay |
| `WithCheckBox(Style)` | `Style` | `widget.CheckBox` base |
| `WithCheckBoxFocus(Style)` | `Style` | `widget.CheckBox` focused overlay |
| `WithHeader(Style)` | `Style` | Canvas header region |
| `WithFooter(Style)` | `Style` | Canvas footer / status bar |
| `WithFocusBorder(Color)` | `Color` | Border colour when a descendant is focused |
| `WithDialog(Style)` | `Style` | `widget.Dialog` panel |
| `WithDialogTitle(Style)` | `Style` | Title text inside a dialog |
| `WithScrim(Style)` | `Style` | Full-screen backdrop behind dialogs |
| `WithTag(Style)` | `Style` | `widget.Label` chip badges |
| `WithNotificationInfo(Style)` | `Style` | Info notification banner |
| `WithNotificationSuccess(Style)` | `Style` | Success notification banner |
| `WithNotificationWarning(Style)` | `Style` | Warning notification banner |
| `WithNotificationError(Style)` | `Style` | Error notification banner |

:::tip WithFocusBorder takes a Color, not a Style
`FocusBorder` is a plain `Color` field (used as a `BorderFG` override inside `layout.Border`), so its builder accepts a `Color` directly rather than a `Style`. All other tokens accept a `Style`.
:::
