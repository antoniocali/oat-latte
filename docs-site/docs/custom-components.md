---
sidebar_position: 8
title: Custom Components
description: How to extend oat-latte by building your own widgets and layout containers.
---

# Custom Components

The built-in widgets cover most needs, but oat-latte is designed to be extended. This page explains the underlying model and shows you how to build your own widgets and layout containers.

---

## How the framework works under the hood

Before writing custom components it helps to understand the three low-level concepts the framework is built on.

### The Component interface

Every element in the tree — layouts, widgets, spacers — implements `Component`:

```go
type Component interface {
    Measure(c Constraint) Size
    Render(buf *Buffer, region Region)
}
```

The render pipeline is a strict **two-pass** process on every frame:

1. **Measure** — the parent calls `child.Measure(constraint)` to ask how much space the child wants. `Constraint` carries `MaxWidth` and `MaxHeight`; `-1` means unconstrained.
2. **Render** — the parent calls `child.Render(buf, region)` to hand the child its allocated rectangle. The child draws into `buf`, clipped to `region`.

Always call `Measure` before `Render` in the same pass. Never cache or store the `Buffer` or `Region` between frames.

### Geometry types

```go
Size{Width, Height int}               // desired or allocated size in cells
Region{X, Y, Width, Height int}       // a rectangle on screen
Constraint{MaxWidth, MaxHeight int}   // available space; -1 = unconstrained
Insets{Top, Right, Bottom, Left int}  // padding / margin
```

### Focusable

Interactive components also implement `Focusable`, which lets the focus manager tab between them and route key events to them:

```go
type Focusable interface {
    Component
    SetFocused(focused bool)
    IsFocused() bool
    HandleKey(ev *KeyEvent) bool // return true = event consumed
}
```

`HandleKey` must return `true` if the component consumed the key event. Returning `false` passes it up the chain (global shortcuts, focus cycling).

### BaseComponent and FocusBehavior

Embed these two types in every custom component to get the boilerplate for free:

```go
type MyWidget struct {
    oat.BaseComponent  // ID, Style, FocusStyle, Title, EnsureID(), EffectiveStyle()
    oat.FocusBehavior  // SetFocused(), IsFocused()
}
```

- Call `w.EnsureID()` in the constructor to auto-assign a unique ID.
- Call `w.EffectiveStyle(w.IsFocused())` in `Render` to get the correctly merged style (focus overlay applied when focused).

### Buffer

`Buffer` wraps the tcell screen with bounds-checked, clipped writes. Components never write to tcell directly. Always call `buf.Sub(region)` to get a clipped sub-buffer before writing:

```go
func (w *MyWidget) Render(buf *oat.Buffer, region oat.Region) {
    sub := buf.Sub(region)
    sub.FillBG(w.Style)
    sub.DrawText(0, 0, "Hello", w.Style)
}
```

---

## Building a custom widget

### 1. Define the struct

```go
type CounterWidget struct {
    oat.BaseComponent
    oat.FocusBehavior
    count int
}

func NewCounterWidget() *CounterWidget {
    w := &CounterWidget{}
    w.EnsureID()
    return w
}
```

### 2. Implement Measure

Return the size your component needs. Never return a size larger than the constraint allows.

```go
func (w *CounterWidget) Measure(c oat.Constraint) oat.Size {
    width := c.MaxWidth
    if width < 0 {
        width = 20 // sensible default when unconstrained
    }
    return oat.Size{Width: width, Height: 1}
}
```

### 3. Implement Render

Draw into a clipped sub-buffer. Use `EffectiveStyle` to get the right colours when focused.

```go
func (w *CounterWidget) Render(buf *oat.Buffer, region oat.Region) {
    style := w.EffectiveStyle(w.IsFocused())
    sub := buf.Sub(region)
    sub.FillBG(style)
    sub.DrawText(0, 0, fmt.Sprintf("Count: %d", w.count), style)
}
```

### 4. Implement HandleKey

Return `true` if you consumed the event; `false` to pass it up the chain.

```go
func (w *CounterWidget) HandleKey(ev *oat.KeyEvent) bool {
    if ev.Key() == tcell.KeyRune {
        switch ev.Rune() {
        case '+':
            w.count++
            return true
        case '-':
            w.count--
            return true
        }
    }
    return false
}
```

### 5. Advertise shortcuts

Implement `KeyBindings()` to have hints appear in the status bar. Bindings with a non-nil `Handler` are also executed automatically by the focus manager.

```go
func (w *CounterWidget) KeyBindings() []oat.KeyBinding {
    return []oat.KeyBinding{
        {Key: tcell.KeyRune, Rune: '+', Label: "+", Description: "Increment"},
        {Key: tcell.KeyRune, Rune: '-', Label: "-", Description: "Decrement"},
    }
}
```

### 6. Apply themes

Pick the theme token that best describes your widget's role and use `Merge` to preserve any caller-set overrides:

```go
func (w *CounterWidget) ApplyTheme(t latte.Theme) {
    w.Style      = t.Panel.Merge(w.Style)
    w.FocusStyle = t.InputFocus.Merge(w.FocusStyle)
}
```

:::warning
Never assign `w.Style = t.SomeToken` directly. This overwrites `BorderExplicitNone` and any other field the caller set before the theme was applied. Always use `Merge`.
:::

---

## Building a custom layout container

If your component holds children, implement `Layout` in addition to `Component`. The framework's tree-walkers (focus collection, theme propagation) depend on `Children()` — a container that hides children from it will break focus and theming.

```go
type TwoColumn struct {
    oat.BaseComponent
    left, right oat.Component
}

func (c *TwoColumn) Children() []oat.Component {
    return []oat.Component{c.left, c.right}
}

func (c *TwoColumn) AddChild(child oat.Component) {
    if c.left == nil {
        c.left = child
    } else {
        c.right = child
    }
}

func (c *TwoColumn) Measure(con oat.Constraint) oat.Size {
    half := oat.Constraint{MaxWidth: con.MaxWidth / 2, MaxHeight: con.MaxHeight}
    ls := c.left.Measure(half)
    rs := c.right.Measure(half)
    h := ls.Height
    if rs.Height > h {
        h = rs.Height
    }
    return oat.Size{Width: con.MaxWidth, Height: h}
}

func (c *TwoColumn) Render(buf *oat.Buffer, region oat.Region) {
    half := region.Width / 2
    leftRegion  := oat.Region{X: region.X,        Y: region.Y, Width: half,              Height: region.Height}
    rightRegion := oat.Region{X: region.X + half, Y: region.Y, Width: region.Width - half, Height: region.Height}
    c.left.Measure(oat.Constraint{MaxWidth: half,                MaxHeight: region.Height})
    c.right.Measure(oat.Constraint{MaxWidth: region.Width - half, MaxHeight: region.Height})
    c.left.Render(buf, leftRegion)
    c.right.Render(buf, rightRegion)
}
```

---

## Full example: editable counter

```go
package main

import (
    "fmt"
    "log"

    "github.com/gdamore/tcell/v2"
    oat "github.com/antoniocali/oat-latte"
    "github.com/antoniocali/oat-latte/latte"
    "github.com/antoniocali/oat-latte/layout"
)

type Counter struct {
    oat.BaseComponent
    oat.FocusBehavior
    count int
}

func NewCounter() *Counter { c := &Counter{}; c.EnsureID(); return c }

func (c *Counter) Measure(con oat.Constraint) oat.Size {
    w := con.MaxWidth
    if w < 0 {
        w = 20
    }
    return oat.Size{Width: w, Height: 1}
}

func (c *Counter) Render(buf *oat.Buffer, region oat.Region) {
    style := c.EffectiveStyle(c.IsFocused())
    buf.Sub(region).FillBG(style)
    buf.Sub(region).DrawText(1, 0, fmt.Sprintf("Count: %d  [+/-]", c.count), style)
}

func (c *Counter) HandleKey(ev *oat.KeyEvent) bool {
    if ev.Key() == tcell.KeyRune {
        switch ev.Rune() {
        case '+':
            c.count++
            return true
        case '-':
            c.count--
            return true
        }
    }
    return false
}

func (c *Counter) ApplyTheme(t latte.Theme) {
    c.Style      = t.Panel.Merge(c.Style)
    c.FocusStyle = t.InputFocus.Merge(c.FocusStyle)
}

func main() {
    counter := NewCounter()
    body    := layout.NewBorder(counter).WithTitle("Counter")

    app := oat.NewCanvas(
        oat.WithTheme(latte.ThemeDark),
        oat.WithBody(body),
        oat.WithPrimary(counter),
    )
    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```
