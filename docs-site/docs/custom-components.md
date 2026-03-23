---
sidebar_position: 8
title: Custom Components
description: How to build your own widgets and layout containers in oat-latte.
---

# Custom Components

## Anatomy of a widget

Every custom component embeds `BaseComponent` and (if interactive) `FocusBehavior`:

```go
type CounterWidget struct {
    oat.BaseComponent
    oat.FocusBehavior
    count int
}

func NewCounterWidget() *CounterWidget {
    w := &CounterWidget{}
    w.EnsureID() // auto-assigns a unique ID
    return w
}
```

## Measure

Return the size your component needs. Respect `Constraint.MaxWidth` and `MaxHeight` — never return a size larger than the constraint allows.

```go
func (w *CounterWidget) Measure(c oat.Constraint) oat.Size {
    h := 1 // one row of text
    width := c.MaxWidth
    if width < 0 {
        width = 20 // sensible default when unconstrained
    }
    return oat.Size{Width: width, Height: h}
}
```

## Render

Draw into a clipped sub-buffer. Always call `buf.Sub(region)` first — never write outside the region you were given.

```go
func (w *CounterWidget) Render(buf *oat.Buffer, region oat.Region) {
    style := w.EffectiveStyle(w.IsFocused()) // merges FocusStyle when focused
    sub := buf.Sub(region)
    sub.FillBG(style)
    sub.DrawText(0, 0, fmt.Sprintf("Count: %d", w.count), style)
}
```

## HandleKey

Return `true` if you consumed the event; `false` to let the canvas try next.

```go
func (w *CounterWidget) HandleKey(ev *oat.KeyEvent) bool {
    switch ev.Key() {
    case tcell.KeyRune:
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

## KeyBindings

Advertise shortcuts to the `StatusBar`. Bindings with a non-nil `Handler` are also executed automatically by `FocusManager.Dispatch`.

```go
func (w *CounterWidget) KeyBindings() []oat.KeyBinding {
    return []oat.KeyBinding{
        {Key: tcell.KeyRune, Rune: '+', Label: "+", Description: "Increment"},
        {Key: tcell.KeyRune, Rune: '-', Label: "-", Description: "Decrement"},
    }
}
```

## ApplyTheme

Pick the theme token that best describes your widget's role and use `Merge` to preserve caller-set overrides:

```go
func (w *CounterWidget) ApplyTheme(t latte.Theme) {
    w.Style      = t.Panel.Merge(w.Style)
    w.FocusStyle = t.FocusBorder.Merge(w.FocusStyle) // or t.InputFocus, etc.
}
```

:::warning
Never assign `w.Style = t.SomeToken` directly. This clobbers `BorderExplicitNone` and any other field the caller set before the theme was applied. Always use `Merge`.
:::

## Custom layout containers

If your component holds children, also implement `Layout`:

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
    leftRegion  := oat.Region{X: region.X,        Y: region.Y, Width: half,             Height: region.Height}
    rightRegion := oat.Region{X: region.X + half, Y: region.Y, Width: region.Width-half, Height: region.Height}
    c.left.Measure(oat.Constraint{MaxWidth: half, MaxHeight: region.Height})
    c.right.Measure(oat.Constraint{MaxWidth: region.Width - half, MaxHeight: region.Height})
    c.left.Render(buf, leftRegion)
    c.right.Render(buf, rightRegion)
}
```

:::tip
`Children()` must return all direct children. The framework tree-walkers (focus collection, theme propagation) depend on it. A container that hides children from `Children()` will break focus and theming.
:::

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
