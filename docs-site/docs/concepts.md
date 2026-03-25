---
sidebar_position: 3
title: Core Concepts
description: The mental model behind oat-latte — components, the render pipeline, and how the pieces fit together.
---

# Core Concepts

## Component

Everything in an oat-latte UI is a `Component`:

```go
type Component interface {
    Measure(c Constraint) Size
    Render(buf *Buffer, region Region)
}
```

The render pipeline is strictly two-pass on every frame:

1. **Measure** — a parent calls `child.Measure(constraint)` to ask how much space the child wants. The `Constraint` carries `MaxWidth` and `MaxHeight`; `-1` means unconstrained.
2. **Render** — the parent calls `child.Render(buf, region)` to hand the child its allocated rectangle. The child draws into `buf`, clipped to `region`.

:::warning
Always call `Measure` before `Render` in the same pass. Never cache or store the `Buffer` or `Region` between frames.
:::

## Geometry types

```go
Size{Width, Height int}          // a component's desired or allocated size
Region{X, Y, Width, Height int}  // a rectangle on screen
Constraint{MaxWidth, MaxHeight int} // available space; -1 = unconstrained
Insets{Top, Right, Bottom, Left int} // padding / margin
```

`oat.Anchor` is a horizontal-position enum (`AnchorLeft`, `AnchorCenter`, `AnchorRight`) used by `Border.WithTitle` and `ProgressBar.WithPercentage` to control where text is placed inside a bar or border rule.

## Layout

A `Component` that holds children also implements `Layout`:

```go
type Layout interface {
    Component
    Children() []Component
    AddChild(child Component)
}
```

The framework's tree walkers (theme propagation, focus collection, widget lookup by ID) all rely on `Children()`. Any custom container type must implement this interface.

## Focusable

Interactive components implement `Focusable`:

```go
type Focusable interface {
    Component
    SetFocused(focused bool)
    IsFocused() bool
    HandleKey(ev *KeyEvent) bool
}
```

`HandleKey` returns `true` if the component consumed the event, or `false` to let the canvas try the next handler (focus cycling, global shortcuts like Esc).

## BaseComponent and FocusBehavior

Every custom component should embed these two types:

```go
type MyWidget struct {
    oat.BaseComponent  // ID, Style, FocusStyle, Title, EnsureID(), EffectiveStyle()
    oat.FocusBehavior  // SetFocused(), IsFocused()
}
```

- Call `e.EnsureID()` in the constructor to auto-assign a unique ID.
- Call `e.EffectiveStyle(e.IsFocused())` in `Render` to get the merged style (focus style overlaid on base style).

## Canvas

`Canvas` is the application root. It owns the tcell screen, the focus manager, the overlay stack, and the event loop.

```go
app := oat.NewCanvas(
    oat.WithTheme(latte.ThemeDark),
    oat.WithHeader(headerComponent),
    oat.WithBody(bodyComponent),
    oat.WithAutoStatusBar(statusBar),
    oat.WithPrimary(firstFocusable),
    oat.WithGlobalKeyBinding(oat.KeyBinding{   // app-wide shortcut
        Key: tcell.KeyCtrlT, Label: "^T",
        Description: "Toggle theme",
        Handler: func() { /* ... */ },
    }),
)
if err := app.Run(); err != nil {
    log.Fatal(err)
}
```

The canvas divides the terminal vertically: **header → body → footer**. Header and footer heights are measured each frame; the body fills what remains.

### Key Canvas methods

| Method | Description |
|---|---|
| `Run()` | Start the event loop; blocks until quit |
| `Quit()` | Signal a graceful exit |
| `SetTheme(t)` | Replace the active theme and re-apply to the full tree |
| `ShowDialog(d)` | Push a modal overlay; focus moves into it |
| `HideDialog()` | Pop the topmost overlay; focus returns to body |
| `FocusByRef(f)` | Jump focus directly to a specific widget |
| `GetWidgetByID(id)` | Look up a widget by its string ID |
| `GetValue(id)` | Get the current value of a widget by ID |
| `InvalidateLayout()` | Force focus re-collection after tree mutation |

## Buffer

`Buffer` wraps `tcell.Screen` with bounds-checked, clipped writes. Components never write to tcell directly.

Use `buf.Sub(region)` to get a sub-buffer clipped to a child's region before passing it down:

```go
func (w *MyWidget) Render(buf *oat.Buffer, region oat.Region) {
    sub := buf.Sub(region)
    sub.FillBG(w.Style)
    sub.DrawText(0, 0, "Hello", w.Style)
}
```

## Putting it together

A typical application follows this shape:

```go
type App struct {
    canvas *oat.Canvas
    notifs *widget.NotificationManager
    list   *widget.List
    detail *widget.Text
}

func (a *App) build() {
    a.list   = widget.NewList(items)
    a.detail = widget.NewText("")

    a.list.WithOnCursorChange(func(_ int, item widget.ListItem) {
        a.detail.SetText(fmt.Sprint(item.Value))
    })

    body := layout.NewHBox()
    body.AddFlexChild(layout.NewBorder(a.list).WithTitle("Items"),   1)
    body.AddFlexChild(layout.NewBorder(a.detail).WithTitle("Detail"), 3)

    statusBar := widget.NewStatusBar()
    a.notifs  = widget.NewNotificationManager()

    themes := []latte.Theme{latte.ThemeDark, latte.ThemeLight, latte.ThemeDracula, latte.ThemeNord}
    themeIdx := 0

    a.canvas = oat.NewCanvas(
        oat.WithTheme(themes[themeIdx]),
        oat.WithBody(body),
        oat.WithAutoStatusBar(statusBar),
        oat.WithPrimary(a.list),
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
