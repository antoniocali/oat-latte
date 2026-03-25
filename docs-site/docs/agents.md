---
sidebar_position: 9
title: Agent Reference
description: Complete API and pattern reference for building oat-latte TUI applications — designed for AI coding agents and developers who need accurate, complete context.
---

# Agent Reference

This page is the authoritative quick-reference for building TUI applications with oat-latte. It covers every package, type, and pattern in one place — useful for AI coding agents and developers who want the full picture without navigating multiple pages.

---

## Module path

```
github.com/antoniocali/oat-latte
```

Sub-packages:

| Import path | Contents |
|---|---|
| `github.com/antoniocali/oat-latte` | Core interfaces, `Canvas`, `Buffer`, `FocusManager`, geometry types |
| `github.com/antoniocali/oat-latte/latte` | `Style`, `Color`, `BorderStyle`, `Theme`, built-in themes |
| `github.com/antoniocali/oat-latte/layout` | `VBox`, `HBox`, `Grid`, `Stack`, `Border`, `Padding`, `VFill`, `HFill`, `FlexChild` |
| `github.com/antoniocali/oat-latte/widget` | `Text`, `Title`, `Button`, `CheckBox`, `EditText`, `List`, `Label`, `ProgressBar`, `StatusBar`, `NotificationManager`, `Dialog` |

---

## Core concepts

### Component

Every element in an oat-latte UI implements `Component`:

```go
type Component interface {
    Measure(c Constraint) Size
    Render(buf *Buffer, region Region)
}
```

The render pipeline is a strict two-pass system:

1. **Measure** — the parent asks the child for its desired size given a `Constraint` (available `MaxWidth`/`MaxHeight`, `-1` means unconstrained).
2. **Render** — the parent hands the child its allocated `Region` and the child draws into a `Buffer`.

Never skip Measure before Render. Never store the `Buffer` or `Region` between frames.

### Layout

A `Component` that holds children implements `Layout`:

```go
type Layout interface {
    Component
    Children() []Component
    AddChild(child Component)
}
```

The framework's tree walkers (theme propagation, focus collection, ID lookup) rely on `Layout.Children()`. Custom container types must implement it.

### Focusable

Interactive components implement `Focusable`:

```go
type Focusable interface {
    Component
    SetFocused(focused bool)
    IsFocused() bool
    HandleKey(ev *KeyEvent) bool // return true = event consumed
}
```

Embed `oat.FocusBehavior` to get `SetFocused`/`IsFocused` for free.

`HandleKey` must return `true` if it consumed the event. Returning `false` tells the canvas to try the next handler (focus cycling, global shortcuts).

### BaseComponent

Embed in every custom component:

```go
type MyWidget struct {
    oat.BaseComponent  // provides ID, Style, FocusStyle, Title, EnsureID(), EffectiveStyle()
    oat.FocusBehavior  // provides SetFocused(), IsFocused()
}
```

Call `e.EnsureID()` in the constructor to auto-assign a unique ID.

`EffectiveStyle(focused bool)` merges `FocusStyle` over `Style` when focused — use this in `Render`.

### Geometry

```go
Size{Width, Height int}                   // desired or allocated size in cells
Region{X, Y, Width, Height int}           // rectangle on screen
Constraint{MaxWidth, MaxHeight int}       // -1 = unconstrained
Insets{Top, Right, Bottom, Left int}      // padding / margin
```

### Anchor

`oat.Anchor` is a general-purpose iota for horizontal positioning used by `ProgressBar.WithPercentage` and `Border.WithTitle`.

```go
oat.AnchorLeft    // default — left edge
oat.AnchorCenter  // centred
oat.AnchorRight   // right edge
```

---

## Style system (latte)

### Style struct

```go
type Style struct {
    FG, BG              Color
    Bold, Italic, Underline, Blink, Reverse bool
    Padding, Margin     Insets
    Border              BorderStyle
    BorderFG, BorderBG  Color
    TextAlign           Alignment
}
```

Zero value means "inherit / use default". Non-zero fields override. Construct inline or with fluent builder methods (`WithFG`, `WithBG`, `WithBorder`, `WithPadding`, etc.).

### Style.Merge

```go
merged := base.Merge(override)
```

Non-zero fields from `override` replace those in `base`. Use this pattern in `ApplyTheme` to let theme act as base while preserving caller-set fields:

```go
func (w *MyWidget) ApplyTheme(t latte.Theme) {
    w.Style = t.Input.Merge(w.Style)         // theme is base; caller wins
    w.FocusStyle = t.InputFocus.Merge(w.FocusStyle)
}
```

Never do `w.Style = t.Input` — that overwrites explicit overrides such as `BorderExplicitNone`.

### Border sentinels

| Constant | Value | Meaning |
|---|---|---|
| `latte.BorderNone` | 0 | Unset; inherits from theme |
| `latte.BorderExplicitNone` | -1 | Actively suppress border (no box drawn) |
| `latte.BorderSingle` | 1 | `┌─┐│└─┘` |
| `latte.BorderRounded` | 2 | `╭─╮│╰─╯` |
| `latte.BorderDouble` | 3 | `╔═╗║╚═╝` |
| `latte.BorderThick` | 4 | `┏━┓┃┗━┛` |
| `latte.BorderDashed` | 5 | `┌╌┐╎└╌┘` |

Check for borders in `Render`:

```go
if style.Border != latte.BorderNone && style.Border != latte.BorderExplicitNone {
    sub.DrawBorderTitle(style.Border, e.Title, latte.Style{}, style, oat.AnchorLeft)
}
```

### Colors

```go
latte.RGB(255, 200, 100)   // true-color
latte.Hex("#FF6600")        // true-color from hex string
latte.ColorRed              // ANSI-16
latte.ColorBrightCyan       // ANSI-16
```

### Themes

Five built-in themes, all applied via `oat.WithTheme(t)` at construction time or switched at runtime via `app.SetTheme(t)`:

```go
latte.ThemeDefault   // ANSI-16, works everywhere
latte.ThemeDark      // true-color, deep navy/blue-cyan
latte.ThemeLight     // true-color, warm off-white
latte.ThemeDracula   // true-color, Dracula palette
latte.ThemeNord      // true-color, Nord arctic palette
```

Theme tokens (fields on `latte.Theme`): `Canvas`, `Text`, `Muted`, `Accent`, `Success`, `Warning`, `Error`, `Panel`, `PanelTitle`, `Input`, `InputFocus`, `ListSelected`, `Button`, `ButtonFocus`, `CheckBox`, `CheckBoxFocus`, `Header`, `Footer`, `FocusBorder`, `Dialog`, `DialogTitle`, `Scrim`, `Tag`, `NotificationInfo`, `NotificationSuccess`, `NotificationWarning`, `NotificationError`.

### Theme builder methods

Every `Theme` value exposes a `With<Token>` method for each field. All methods return a new `Theme` by value — the originals are never mutated. Style-typed methods use `Style.Merge` internally so only non-zero fields of the supplied `Style` are applied.

```go
// Nord but with no borders anywhere
borderless := latte.ThemeNord.
    WithPanel(latte.Style{Border: latte.BorderExplicitNone}).
    WithInput(latte.Style{Border: latte.BorderExplicitNone}).
    WithButton(latte.Style{Border: latte.BorderExplicitNone}).
    WithDialog(latte.Style{Border: latte.BorderExplicitNone}).
    WithName("nord-borderless")

// Dark theme with a custom accent colour
pink := latte.ThemeDark.
    WithAccent(latte.Style{FG: latte.Hex("#ff69b4")}).
    WithFocusBorder(latte.Hex("#ff69b4")).
    WithName("dark-pink")
```

`WithFocusBorder` accepts a `Color` directly (not a `Style`) because `FocusBorder` is a plain `Color` field. All other token methods accept a `Style`.

---

## Canvas

### Construction

```go
app := oat.NewCanvas(
    oat.WithTheme(latte.ThemeDark),
    oat.WithHeader(headerComponent),
    oat.WithBody(bodyComponent),
    oat.WithAutoStatusBar(statusBar),   // auto-populates footer with key hints
    oat.WithPrimary(firstFocusable),    // override DFS-first focus
    oat.WithNotificationManager(notifs), // wires channel + mounts as persistent overlay
    oat.WithGlobalKeyBinding(           // app-wide shortcut (see below)
        oat.KeyBinding{
            Key: tcell.KeyCtrlT, Label: "^T", Description: "Toggle theme",
            Handler: func() { /* ... */ },
        },
    ),
)
if err := app.Run(); err != nil {
    log.Fatal(err)
}
```

### Key methods

```go
app.Run() error                        // start event loop; blocks until quit
app.Quit()                             // signal graceful exit
app.SetTheme(t latte.Theme)            // replace active theme and re-apply to full tree
app.ShowDialog(d Component)            // push modal overlay, steal focus; dismissed by Esc
app.ShowPersistentOverlay(d Component) // render on top always; never dismissed by Esc
app.HideDialog()                       // pop topmost overlay, restore body focus
app.HasOverlay() bool                  // true while any dialog is visible
app.FocusByRef(target Focusable)       // jump focus to a specific widget
app.GetWidgetByID(id string) Component
app.GetValue(id string) (interface{}, bool)
app.InvalidateLayout()                 // force full focus re-collection (after tree mutation)
```

### Layout regions

The canvas divides the screen vertically: header → body → footer. Header and footer heights are measured each frame; body fills the remainder.

---

## Layout containers

### VBox / HBox

```go
vbox := layout.NewVBox()
vbox.AddChild(widget.NewText("Label"))
vbox.AddFlexChild(editText, 1)          // flex weight 1 = share remaining space
vbox.AddChild(layout.NewVFill())        // spacer; equivalent to AddFlexChild weight 1
vbox.AddChild(layout.NewVFill().WithMaxSize(1))  // fixed 1-row gap

hbox := layout.NewHBox(child1, child2)  // variadic shorthand
hbox.AddFlexChild(progressBar, 1)
```

### FlexChild

`layout.NewFlexChild(child, weight...)` wraps any `Component` as a flex slot so it can be passed inline to variadic constructors:

```go
// Equivalent to: vbox := layout.NewVBox(); vbox.AddChild(title); vbox.AddFlexChild(body, 1); vbox.AddChild(btnRow)
vbox := layout.NewVBox(
    title,
    layout.NewFlexChild(body),   // weight defaults to 1
    btnRow,
)
```

- Weight defaults to `1`; minimum effective weight is `1`.
- Implements `oat.Layout` via `Children()` — theme propagation and focus collection recurse into the wrapped component automatically.

### Border

```go
panel := layout.NewBorder(innerComponent).
    WithTitle("My Panel").
    WithTitleStyle(latte.Style{Bold: true})

// Centred title:
panel := layout.NewBorder(innerComponent).
    WithTitle("My Panel", oat.AnchorCenter)

// Rounded corners (╭─╮ / ╰─╯) — only valid for BorderSingle:
panel := layout.NewBorder(innerComponent).
    WithTitle("My Panel").
    WithRoundedCorner(true)

// Custom style (e.g. explicit padding):
panel := layout.NewBorder(innerComponent).
    WithStyle(latte.Style{Padding: latte.Insets{Bottom: 1}}).
    WithTitle("My Panel")
```

`Border` automatically sets its border color to `t.FocusBorder` when any descendant is focused (after theme application).

#### WithTitle anchor

`WithTitle` accepts an optional `oat.Anchor` as a second argument:

```go
func (b *Border) WithTitle(title string, anchor ...oat.Anchor) *Border
```

| Anchor | Result |
|---|---|
| `oat.AnchorLeft` (default) | `╭─ My Panel ──────────╮` |
| `oat.AnchorCenter` | `╭──── My Panel ────────╮` |
| `oat.AnchorRight` | `╭────────── My Panel ──╮` |

Omitting the anchor defaults to `oat.AnchorLeft`, keeping all existing call sites unchanged.

#### WithRoundedCorner

```go
func (b *Border) WithRoundedCorner(rounded bool) *Border
```

- `true` — switches the border style to `BorderRounded` (`╭─╮│╰─╯`).
- `false` — restores `BorderSingle` if the current style is `BorderRounded`; no-op otherwise.
- **Panics** if called with `true` when the current border style is `BorderDouble`, `BorderThick`, or `BorderDashed`. Unicode provides arc corner codepoints (`╭╮╰╯`) only for light-weight strokes (`─` `│`); they do not connect visually to double (`═` `║`), heavy (`━` `┃`), or dashed (`╌` `╎`) lines. Use `WithStyle(latte.Style{Border: latte.BorderRounded})` to switch style entirely instead.

### Padding

```go
padded := layout.NewPaddingUniform(child, 1)          // 1 cell all sides
padded := layout.NewPadding(child, latte.Insets{Top: 1, Left: 2})
```

### Dialog

```go
dlg := widget.NewDialog("Title").
    WithChild(bodyComponent).
    WithSize(widget.DialogPercent(50), widget.DialogPercent(60))  // 50% × 60% of terminal

// OR fixed size:
dlg := widget.NewDialog("Confirm").
    WithChild(bodyComponent).
    WithMaxSize(52, 9)   // exactly 52 × 9 cells

app.ShowDialog(dlg)
// inside a button callback:
app.HideDialog()
```

Dialog always centres itself and paints a full-screen scrim behind it.

### Grid

```go
g := layout.NewGrid(2, 3)          // 2 rows, 3 cols
g.AddChildAt(widget, 0, 0, 1, 1)   // row, col, rowSpan, colSpan
g.WithGap(0, 1)                     // rowGap, colGap
```

---

## Widgets

### Text

```go
t := widget.NewText("Hello, world!")
t.SetText("Updated text")
t.GetText() string
```

Supports word-wrap (bounded by render width) and vertical scroll (`Scrollable`).

### Button

```go
btn := widget.NewButton("Save", func() {
    // pressed
}).WithID("save-btn")

// With rounded border corners (╭─╮ / ╰─╯):
btn := widget.NewButton("OK", fn).
    WithStyle(latte.Style{Border: latte.BorderSingle}).
    WithRoundedCorner(true)
```

Activated by `Enter` or `Space`.

Border presence is determined solely by `b.Style` (the unfocused base style). `FocusStyle` / `ButtonFocus` carry only colour and attribute overrides (e.g. `Reverse: true`, `BorderFG`). This means the button's layout shape — and therefore `Measure` output — is stable regardless of focus state.

When `b.Style` has a border set, `Measure` returns `Height: 3` (top border + label row + bottom border) and the border is drawn at render time. Without a border, `Measure` returns `Height: 1` and the label is rendered as `"[ label ]"`.

All built-in themes set `Button.Border: BorderSingle` so buttons always render with a visible border. `ButtonFocus` carries only `Reverse: true` and an accent `BorderFG` to highlight the active button.

#### WithRoundedCorner

```go
func (b *Button) WithRoundedCorner(rounded bool) *Button
```

- `true` — draws border arc corners (`╭╮╰╯`) instead of square ones (`┌┐└┘`).
- `false` — no-op; arc corners are only applied when `true`.
- **Panics at render time** if `WithRoundedCorner(true)` is set and the effective border style is `BorderDouble`, `BorderThick`, or `BorderDashed`. Arc corner codepoints exist only for light-weight strokes (`─` `│`).
- **`WithStyle` panics immediately** at construction time if called with `BorderDouble`, `BorderThick`, or `BorderDashed` — only `BorderNone`, `BorderExplicitNone`, `BorderSingle`, and `BorderRounded` are valid for Button.

### CheckBox

```go
cb := widget.NewCheckBox("Enable feature").
    WithOnToggle(func(checked bool) { /* ... */ })
cb.IsChecked() bool
cb.SetChecked(true)
```

### EditText

```go
// Single-line
input := widget.NewEditText().
    WithID("username").
    WithPlaceholder("Enter username…").
    WithHint("Username").           // persistent label above the field
    WithMaxLength(64).
    WithOnChange(func(s string) { /* live */ }).
    WithOnSave(func(s string)   { /* ^S pressed */ }).
    WithOnCancel(func()         { /* ^G pressed */ })

// Multi-line
body := widget.NewMultiLineEditText().
    WithHint("Description").
    WithPlaceholder("Write here…")

// Borderless (hint replaces the border's visual framing)
field := widget.NewEditText().
    WithStyle(latte.Style{Border: latte.BorderExplicitNone}).
    WithHint("Email")

input.SetText("hello")
input.GetText() string
```

Built-in key bindings: `^S` Save, `^G` Cancel, `^K` kill-to-EOL, `^U` kill-from-SOL, `^A`/Home start-of-line, `^E`/End end-of-line.

### List

```go
items := []widget.ListItem{
    {Label: "First item",  Value: 1},
    {Label: "Second item", Value: 2},
}
list := widget.NewList(items).
    WithID("my-list").
    WithOnSelect(func(idx int, item widget.ListItem) { /* Enter pressed */ }).
    WithOnDelete(func(idx int, item widget.ListItem) { /* Del pressed */ }).
    WithOnCursorChange(func(idx int, item widget.ListItem) { /* live preview */ })

list.SetItems(newItems)
list.SelectedItem() widget.ListItem
list.SelectedIndex() int
```

### Label (tag chips)

```go
lbl := widget.NewLabel([]string{"go", "tui"})
lbl.SetLabels(tags)

// Without background fill (keeps FG colour, strips BG):
lbl := widget.NewLabel(tags).WithHighlight(false)
```

Renders inline chips separated by `·`.

`WithHighlight(false)` strips the chip background colour while keeping the foreground colour and text attributes. Default is `true` (chips render with filled background). The `highlight` setting does not affect the FG colour or bold/italic attributes.

### ProgressBar

```go
pb := widget.NewProgressBar().
    WithPercentage(true)          // show percent at left (default)
pb.SetValue(0.75)    // 0.0 – 1.0
```

`WithPercentage(show bool, anchor ...oat.Anchor)` controls whether a `" XX%"` label is rendered and where. The anchor is optional and defaults to `oat.AnchorLeft`:

| Anchor | Result |
|---|---|
| `oat.AnchorLeft` | `" 75% ████░░░░"` — label at the left |
| `oat.AnchorCenter` | `"████ 75% ░░░░"` — label stamped into the middle of the bar |
| `oat.AnchorRight` | `"████░░░░ 75%"` — label at the right |

`WithShowPercent(bool)` still works but does not change the anchor. Prefer `WithPercentage`.

### NotificationManager

```go
notifs := widget.NewNotificationManager()

app := oat.NewCanvas(
    oat.WithTheme(latte.ThemeDark),
    oat.WithBody(body),
    oat.WithNotificationManager(notifs),  // wires channel + mounts as persistent overlay
)

notifs.Push("Saved", widget.NotificationKindSuccess, 2*time.Second)
notifs.Push("Error!", widget.NotificationKindError, 0)  // 0 = no auto-dismiss
```

### StatusBar

```go
bar := widget.NewStatusBar()
// Pass to canvas:
oat.WithAutoStatusBar(bar)
```

Auto-populates with the focused component's `KeyBindings()` plus any registered global bindings.

---

## Focus system

### Automatic collection

On `Canvas.Run()` the framework performs a DFS over the component tree and registers every `Focusable` node. Order is DFS (depth-first, children left-to-right), which corresponds to visual top-left to bottom-right order.

### Cycling

`Tab` → next focusable. `Shift+Tab` → previous. Arrow keys cycle if the focused component does not consume them (returns `false` from `HandleKey`).

### Keyboard dispatch

```
Tab / Shift+Tab          → FocusManager.Next() / Prev()
Any key                  → FocusManager.Dispatch(ev)
  ├─ walk KeyBindings() with Handler != nil → invoke handler (consumed)
  └─ else → focused.HandleKey(ev)
       ├─ true  → consumed
       └─ false → canvas.dispatchGlobal(ev)
            ├─ matching global binding found → invoke handler (consumed)
            └─ no match → canvas tries Left/Right focus cycling
```

Global bindings are checked **after** the focused widget so that widgets can shadow them when needed (e.g. an `EditText` consuming `Esc` to cancel editing rather than triggering the app-level quit).

### Global key bindings

Register app-wide shortcuts with `WithGlobalKeyBinding` at construction time:

```go
app := oat.NewCanvas(
    oat.WithTheme(latte.ThemeDark),
    oat.WithBody(body),
    oat.WithGlobalKeyBinding(
        oat.KeyBinding{
            Key:         tcell.KeyCtrlT,
            Label:       "^T",
            Description: "Toggle theme",
            Handler: func() {
                current = (current + 1) % len(themes)
                app.SetTheme(themes[current])
            },
        },
        oat.KeyBinding{
            Key:         tcell.KeyCtrlH,
            Label:       "^H",
            Description: "Help",
            Handler:     func() { app.ShowDialog(helpDialog) },
        },
    ),
)
```

- Variadic: pass multiple bindings in one call, or call `WithGlobalKeyBinding` multiple times — bindings accumulate.
- Global bindings appear in the status bar alongside the focused widget's own hints.
- A focused widget can shadow a global binding by returning `true` from `HandleKey` for the same key.

### Custom Focusable (proxy pattern)

Wrap an existing widget to intercept specific keys without modifying the widget itself:

```go
type myProxy struct {
    *widget.List
    app *App
}

func (p *myProxy) HandleKey(ev *oat.KeyEvent) bool {
    if ev.Key() == tcell.KeyRune && ev.Rune() == 'n' {
        p.app.showNewDialog()
        return true   // consumed
    }
    return p.List.HandleKey(ev)   // delegate the rest
}

func (p *myProxy) KeyBindings() []oat.KeyBinding {
    extra := []oat.KeyBinding{
        {Key: tcell.KeyRune, Rune: 'n', Label: "n", Description: "New item"},
    }
    return append(extra, p.List.KeyBindings()...)
}
```

### Programmatic focus

```go
app.FocusByRef(myTitleInput)   // jump to a specific widget
```

`FocusByRef` is pointer identity. The target must be in the current focus tree (body or active dialog).

### FocusGuard — context-aware Tab cycling

Implement `oat.FocusGuard` to dynamically exclude a component (and its whole subtree) from Tab cycling:

```go
type FocusGuard interface {
    IsFocusable() bool
}
```

When `IsFocusable()` returns `false`, `walkFocusable` skips the node **and all its descendants**. This is the correct way to build context-sensitive panels where entire subtrees should be unreachable depending on application state.

**Pattern — mode-gated inputs:**

```go
// Thin wrapper that only participates in Tab cycling when editorMode is active.
type editorInputGuard struct {
    *widget.EditText
    app *App
}
func (g *editorInputGuard) IsFocusable() bool { return g.app.editorMode }

// A custom component can implement FocusGuard directly.
func (s *myShim) IsFocusable() bool { return !s.app.editorMode }
```

Call `canvas.InvalidateLayout()` whenever the mode changes so the focus tree is rebuilt:

```go
func (a *App) setEditorMode(on bool) {
    if a.editorMode == on { return }
    a.editorMode = on
    a.canvas.InvalidateLayout()
}
```

After `InvalidateLayout()`, call `canvas.FocusByRef(target)` to set the desired initial focus for the new mode.

### KeyBinding

```go
type KeyBinding struct {
    Key         tcell.Key
    Rune        rune       // only used when Key == tcell.KeyRune
    Mod         tcell.ModMask
    Label       string     // short hint, e.g. "^S"
    Description string     // e.g. "Save"
    Handler     func()     // nil = display-only hint; non-nil = executed by Dispatch
}
```

---

## Theming a custom widget

```go
func (w *MyWidget) ApplyTheme(t latte.Theme) {
    // Theme as base, caller-set fields take precedence via Merge.
    w.Style = t.Input.Merge(w.Style)
    w.FocusStyle = t.InputFocus.Merge(w.FocusStyle)
}
```

Register `ApplyTheme` on the type (not a pointer receiver) so the framework's tree walker can call it on the embedded value.

---

## Common patterns

### Basic two-panel app

```go
list := widget.NewList(items).WithID("list")
detail := widget.NewText("")

list.WithOnCursorChange(func(_ int, item widget.ListItem) {
    detail.SetText(fmt.Sprint(item.Value))
})

body := layout.NewHBox()
body.AddFlexChild(layout.NewBorder(list).WithTitle("Items"), 1)
body.AddFlexChild(layout.NewBorder(detail).WithTitle("Detail"), 3)
```

### Modal dialog

```go
func showConfirm(app *oat.Canvas, msg string, onConfirm func()) {
    cancelBtn := widget.NewButton("Cancel", func() { app.HideDialog() })
    okBtn     := widget.NewButton("OK",     func() { onConfirm(); app.HideDialog() })

    btnRow := layout.NewHBox()
    btnRow.AddChild(layout.NewHFill())
    btnRow.AddChild(cancelBtn)
    btnRow.AddChild(layout.NewHFill().WithMaxSize(2))
    btnRow.AddChild(okBtn)

    body := layout.NewPaddingUniform(layout.NewVBox(
        widget.NewText(msg),
        layout.NewVFill().WithMaxSize(1),
        btnRow,
    ), 1)

    app.ShowDialog(
        widget.NewDialog("Confirm").
            WithChild(body).
            WithMaxSize(50, 9),
    )
}
```

### Borderless editor with hints

```go
titleInput := widget.NewEditText().
    WithStyle(latte.Style{Border: latte.BorderExplicitNone}).
    WithHint("Title").
    WithPlaceholder("Untitled…")

bodyInput := widget.NewMultiLineEditText().
    WithStyle(latte.Style{Border: latte.BorderExplicitNone}).
    WithHint("Body").
    WithPlaceholder("Write here…")

editorVBox := layout.NewVBox(titleInput)
editorVBox.AddFlexChild(bodyInput, 1)

panel := layout.NewBorder(editorVBox).WithTitle("Editor")
```

### Notification toasts

```go
notifs := widget.NewNotificationManager()

app := oat.NewCanvas(
    oat.WithTheme(latte.ThemeDark),
    oat.WithBody(body),
    oat.WithNotificationManager(notifs),  // wires channel + mounts as persistent overlay
)

// later, from any callback:
notifs.Push("Saved successfully", widget.NotificationKindSuccess, 2*time.Second)
```

### Full application skeleton

```go
type App struct {
    canvas *oat.Canvas
    notifs *widget.NotificationManager
    // ... widgets
}

func (a *App) build() {
    statusBar := widget.NewStatusBar()
    a.notifs = widget.NewNotificationManager()

    // ... build component tree ...

    themes := []latte.Theme{latte.ThemeDark, latte.ThemeLight, latte.ThemeDracula, latte.ThemeNord}
    themeIdx := 0

    a.canvas = oat.NewCanvas(
        oat.WithTheme(themes[themeIdx]),
        oat.WithHeader(header),
        oat.WithBody(body),
        oat.WithAutoStatusBar(statusBar),
        oat.WithPrimary(primaryFocusable),
        oat.WithNotificationManager(a.notifs),  // wires channel + mounts as persistent overlay
        oat.WithGlobalKeyBinding(oat.KeyBinding{
            Key:         tcell.KeyCtrlT,
            Label:       "^T",
            Description: "Toggle theme",
            Handler: func() {
                themeIdx = (themeIdx + 1) % len(themes)
                a.canvas.SetTheme(themes[themeIdx])
            },
        }),
    )
}

func main() {
    a := &App{}
    a.build()
    if err := a.canvas.Run(); err != nil {
        log.Fatal(err)
    }
}
```

---

## Constraints and invariants

- Never call `Render` without having called `Measure` first in the same pass.
- Never write to a `Buffer` outside the `Region` passed to `Render` — use `buf.Sub(region)` to get a clipped sub-buffer and write into that.
- `BorderExplicitNone` (`-1`) actively suppresses a border. Check both `BorderNone` and `BorderExplicitNone` in render guards.
- `Style.Merge` preserves `BorderExplicitNone` through the cascade — do not use direct struct assignment in `ApplyTheme`.
- `Canvas.InvalidateLayout()` must be called after any dynamic addition or removal of components from the tree to re-collect focusable nodes.
- Key event handlers run on the main goroutine. Use background goroutines safely by pushing to `notifs` directly; `WithNotificationManager` wires the re-render channel automatically.
- `oat.WithNotificationManager(notifs)` mounts `NotificationManager` as a persistent (non-modal) overlay. It is never dismissed by Esc and always renders on top of modal dialogs.
- Global bindings registered with `WithGlobalKeyBinding` fire **after** the focused widget. A widget can shadow a global binding by returning `true` from `HandleKey`.
- `app.SetTheme(t)` re-applies the theme to the entire tree including all overlays and persistent overlays. It resets the canvas background style so the new theme's `Canvas` token takes effect.
