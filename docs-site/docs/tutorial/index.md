---
sidebar_position: 1
title: "Tutorial: Task List"
description: Build a fully working terminal task-list app from scratch, one step at a time.
---

# Tutorial: Task List

This tutorial builds a real terminal app from scratch. By the end you will have a task list with:

- a scrollable list of tasks
- an "Add task" dialog with a text input
- delete with confirmation
- toast notifications
- a status bar showing key hints

Each step is self-contained and runnable. You can stop at any step and have a working program.

---

## What you will build

```
┌─ Tasks ──────────────────────────────────────────────────────────┐
│  ▶  Buy groceries                                                 │
│     Write tutorial                                                │
│     Ship v0.1.0                                                   │
│                                                                   │
└──────────────────────────────────────────────────────────────────┘
  n · New    Del · Delete    Esc · Quit
```

---

## Prerequisites

- Go 1.21 or later
- A true-color terminal (iTerm2, Ghostty, Windows Terminal, etc.)
- oat-latte installed:

```sh
go get github.com/antoniocali/oat-latte
```

Create a new module for the tutorial:

```sh
mkdir tasklist && cd tasklist
go mod init tasklist
go get github.com/antoniocali/oat-latte
```

---

## Step 1 — Hello, terminal

Create `main.go` and get a window on screen.

```go
package main

import (
	"log"

	oat "github.com/antoniocali/oat-latte"
	"github.com/antoniocali/oat-latte/latte"
	"github.com/antoniocali/oat-latte/widget"
)

func main() {
	body := widget.NewText("Hello, terminal!")

	app := oat.NewCanvas(
		oat.WithTheme(latte.ThemeDark),
		oat.WithBody(body),
	)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
```

Run it:

```sh
go run .
```

You should see `Hello, terminal!` centred in a dark terminal. Press **Esc** to quit.

:::tip What just happened?
`NewCanvas` owns the tcell screen and event loop. `WithBody` sets the component that fills the middle of the screen. `WithTheme` propagates a colour scheme to every component in the tree.
:::

---

## Step 2 — A list of tasks

Replace the `Text` with a `List`. A `List` is focusable and handles keyboard navigation automatically.

```go
package main

import (
	"log"

	oat "github.com/antoniocali/oat-latte"
	"github.com/antoniocali/oat-latte/latte"
	"github.com/antoniocali/oat-latte/layout"
	"github.com/antoniocali/oat-latte/widget"
)

func main() {
	items := []widget.ListItem{
		{Label: "Buy groceries"},
		{Label: "Write tutorial"},
		{Label: "Ship v0.1.0"},
	}

	list := widget.NewList(items)

	body := layout.NewBorder(list).WithTitle("Tasks")

	app := oat.NewCanvas(
		oat.WithTheme(latte.ThemeDark),
		oat.WithBody(body),
		oat.WithPrimary(list),
	)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
```

Run it. Use **↑ ↓** to move the cursor, **Esc** to quit.

:::tip What changed?
`layout.NewBorder` wraps any component in a titled box. `WithPrimary(list)` tells the canvas to give focus to the list on startup rather than picking the first focusable it finds in the tree.
:::

:::tip Customising the list cursor and highlight
By default the list draws `>` next to the selected row and fills it with the theme's selection colour. Both are optional:

```go
list := widget.NewList(items).
    WithCursor("▶").         // any string — "→", "•", "❯", or "" to hide it
    WithHighlight(false)     // keep the cursor glyph but drop the background fill
```

`WithCursor("")` hides the gutter symbol entirely; `WithHighlight(false)` keeps it but removes the coloured background — useful for minimal or transparent UIs.
:::

:::tip Positioning the border title
`WithTitle` accepts an optional anchor as a second argument:

```go
layout.NewBorder(list).WithTitle("Tasks")                      // left (default)
layout.NewBorder(list).WithTitle("Tasks", oat.AnchorCenter)    // centred
layout.NewBorder(list).WithTitle("Tasks", oat.AnchorRight)     // right
```

Omitting the anchor always defaults to `oat.AnchorLeft`, so existing code needs no changes.
:::

---

## Step 3 — Status bar with key hints

Add a `StatusBar` in the footer. It auto-populates from the focused component's `KeyBindings()`.

```go
package main

import (
	"log"

	oat "github.com/antoniocali/oat-latte"
	"github.com/antoniocali/oat-latte/latte"
	"github.com/antoniocali/oat-latte/layout"
	"github.com/antoniocali/oat-latte/widget"
)

func main() {
	items := []widget.ListItem{
		{Label: "Buy groceries"},
		{Label: "Write tutorial"},
		{Label: "Ship v0.1.0"},
	}

	list := widget.NewList(items)
	body := layout.NewBorder(list).WithTitle("Tasks")

	// highlight-start
	statusBar := widget.NewStatusBar()
	// highlight-end

	app := oat.NewCanvas(
		oat.WithTheme(latte.ThemeDark),
		oat.WithBody(body),
		// highlight-start
		oat.WithAutoStatusBar(statusBar),
		// highlight-end
		oat.WithPrimary(list),
	)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
```

The footer now shows the list's built-in key hints (↑↓ to move, Enter to select, Del to delete).

---

## Step 4 — Add tasks with a dialog

When the user presses **n**, show a dialog with a text input.

This is the first time you need application-level state, so introduce an `App` struct.

```go
package main

import (
	"log"

	oat "github.com/antoniocali/oat-latte"
	"github.com/antoniocali/oat-latte/latte"
	"github.com/antoniocali/oat-latte/layout"
	"github.com/antoniocali/oat-latte/widget"
	"github.com/gdamore/tcell/v2"
)

// App holds all shared state.
type App struct {
	canvas *oat.Canvas
	list   *widget.List
	items  []widget.ListItem
}

// showNewDialog opens a modal to add a task.
func (a *App) showNewDialog() {
	input := widget.NewEditText().
		WithHint("Task name").
		WithPlaceholder("What needs doing?")

	cancelBtn := widget.NewButton("Cancel", func() {
		a.canvas.HideDialog()
	})

	doAdd := func() {
		name := input.GetText()
		if name == "" {
			return
		}
		a.items = append(a.items, widget.ListItem{Label: name})
		a.list.SetItems(a.items)
		a.canvas.HideDialog()
	}

	// Wire ^S on the input to trigger Add as well.
	input.WithOnSave(func(_ string) { doAdd() })

	createBtn := widget.NewButton("Add", doAdd)

	btnRow := layout.NewHBox()
	btnRow.AddChild(layout.NewHFill())
	btnRow.AddChild(cancelBtn)
	btnRow.AddChild(layout.NewHFill().WithMaxSize(2))
	btnRow.AddChild(createBtn)

	body := layout.NewPaddingUniform(
		layout.NewVBox(
			widget.NewText("Enter a name for the new task."),
			layout.NewVFill().WithMaxSize(1),
			input,
			layout.NewVFill().WithMaxSize(1),
			btnRow,
		), 1)

	dlg := widget.NewDialog("New Task").
		WithChild(body).
		WithMaxSize(50, 13)

	a.canvas.ShowDialog(dlg)
}

// listProxy wraps the list to intercept the 'n' key.
type listProxy struct {
	*widget.List
	app *App
}

func (p *listProxy) HandleKey(ev *oat.KeyEvent) bool {
	if ev.Key() == tcell.KeyRune && ev.Rune() == 'n' {
		p.app.showNewDialog()
		return true
	}
	return p.List.HandleKey(ev)
}

func (p *listProxy) KeyBindings() []oat.KeyBinding {
	extra := []oat.KeyBinding{
		{Key: tcell.KeyRune, Rune: 'n', Label: "n", Description: "New task"},
	}
	return append(extra, p.List.KeyBindings()...)
}

func main() {
	a := &App{}
	a.items = []widget.ListItem{
		{Label: "Buy groceries"},
		{Label: "Write tutorial"},
		{Label: "Ship v0.1.0"},
	}

	a.list = widget.NewList(a.items)
	proxy := &listProxy{List: a.list, app: a}

	body := layout.NewBorder(proxy).WithTitle("Tasks")
	statusBar := widget.NewStatusBar()

	a.canvas = oat.NewCanvas(
		oat.WithTheme(latte.ThemeDark),
		oat.WithBody(body),
		oat.WithAutoStatusBar(statusBar),
		oat.WithPrimary(proxy),
	)

	if err := a.canvas.Run(); err != nil {
		log.Fatal(err)
	}
}
```

Press **n** to open the dialog, type a name, press **Enter** on the Add button (or Tab to reach it), and the list updates.

:::tip The proxy pattern
oat-latte uses a **proxy pattern** to add key handling to existing widgets without modifying them. Wrap the widget, override `HandleKey` for the keys you want to intercept, and delegate everything else back to the wrapped widget. `KeyBindings()` exposes the extra hints to the status bar.

See [Focus — The proxy pattern](../focus#the-proxy-pattern) for the full explanation.
:::

:::tip Sizing dialogs with bordered buttons
All built-in themes set `Border: BorderSingle` on buttons, so each button always measures `Height: 3` (top border + label + bottom border). When computing `WithMaxSize`, budget:

- **2** for the dialog border
- **2** for `NewPaddingUniform(..., 1)`
- **1** per `Text` row
- **4** per `EditText` with a hint (`WithHint`)
- **1** per `VFill.WithMaxSize(1)` spacer
- **3** for the button row

The `showNewDialog` above: 2 + 2 + 1 + 1 + 4 + 1 + 3 = **14** minimum. `WithMaxSize(50, 13)` is slightly under that so the two `VFill` spacers share the one available flex row; the layout still fits comfortably.
:::

---

## Step 5 — Delete with confirmation

Wire the list's built-in `OnDelete` callback to show a confirm dialog before removing the item.

Add a `showConfirmDialog` method to `App` and hook `WithOnDelete` on the list:

```go
// showConfirmDialog shows a yes/no dialog. onConfirm is called if the user
// chooses Yes.
func (a *App) showConfirmDialog(msg string, onConfirm func()) {
	noBtn := widget.NewButton("No", func() {
		a.canvas.HideDialog()
	})
	yesBtn := widget.NewButton("Yes", func() {
		onConfirm()
		a.canvas.HideDialog()
	})

	btnRow := layout.NewHBox()
	btnRow.AddChild(layout.NewHFill())
	btnRow.AddChild(noBtn)
	btnRow.AddChild(layout.NewHFill().WithMaxSize(2))
	btnRow.AddChild(yesBtn)

	body := layout.NewPaddingUniform(
		layout.NewVBox(
			widget.NewText(msg),
			layout.NewVFill().WithMaxSize(1),
			btnRow,
		), 1)

	dlg := widget.NewDialog("Confirm").
		WithChild(body).
		WithMaxSize(48, 9)

	a.canvas.ShowDialog(dlg)
}
```

Then, when building the list, attach the callback:

```go
a.list = widget.NewList(a.items).
    WithOnDelete(func(idx int, item widget.ListItem) {
        a.showConfirmDialog(
            "Delete \""+item.Label+"\"?",
            func() {
                a.items = append(a.items[:idx], a.items[idx+1:]...)
                a.list.SetItems(a.items)
            },
        )
    })
```

Press **Del** on a task — a confirmation box appears. **No** dismisses it; **Yes** removes the item.

---

## Step 6 — Toast notifications

Add a `NotificationManager` so the app can show transient success/error toasts.

```go
// Add to App struct:
notifs *widget.NotificationManager
```

Wire it via `oat.WithNotificationManager` when constructing the canvas:

```go
a.notifs = widget.NewNotificationManager()

a.canvas = oat.NewCanvas(
    oat.WithTheme(latte.ThemeDark),
    oat.WithBody(body),
    oat.WithAutoStatusBar(statusBar),
    oat.WithPrimary(proxy),
    oat.WithNotificationManager(a.notifs),  // wires channel + mounts as persistent overlay
)
```

Now push a notification after any state change. In `showNewDialog`'s Add callback:

```go
a.items = append(a.items, widget.ListItem{Label: name})
a.list.SetItems(a.items)
a.canvas.HideDialog()
// highlight-next-line
a.notifs.Push("Task added", widget.NotificationKindSuccess, 2*time.Second)
```

And in the confirm delete callback:

```go
a.items = append(a.items[:idx], a.items[idx+1:]...)
a.list.SetItems(a.items)
// highlight-next-line
a.notifs.Push("Task deleted", widget.NotificationKindSuccess, 2*time.Second)
```

:::tip WithNotificationManager
`oat.WithNotificationManager(notifs)` handles both wiring the timer channel and mounting the manager as a persistent overlay. There is no need to call `SetNotifyChannel` or `ShowPersistentOverlay` manually.
:::

---

## Step 7 — Delete with `d` as well as `Del`

The list already fires the `onDelete` callback on the `Del` key. To also support `d`, extend the proxy's `HandleKey` to intercept that rune and forward a synthetic delete event to the underlying list:

```go
func (p *listProxy) HandleKey(ev *oat.KeyEvent) bool {
	if ev.Key() == tcell.KeyRune {
		switch ev.Rune() {
		case 'n':
			p.app.showNewDialog()
			return true
		// highlight-start
		case 'd':
			return p.List.HandleKey(tcell.NewEventKey(tcell.KeyDelete, 0, tcell.ModNone))
		// highlight-end
		}
	}
	return p.List.HandleKey(ev)
}
```

Add the hint to `KeyBindings` so it appears in the status bar:

```go
func (p *listProxy) KeyBindings() []oat.KeyBinding {
	return append(
		[]oat.KeyBinding{
			{Key: tcell.KeyRune, Rune: 'n', Label: "n", Description: "New task"},
			// highlight-next-line
			{Key: tcell.KeyRune, Rune: 'd', Label: "d", Description: "Delete task"},
		},
		p.List.KeyBindings()...,
	)
}
```

:::tip Synthetic events
`tcell.NewEventKey(tcell.KeyDelete, 0, tcell.ModNone)` constructs a key event identical to what the terminal sends when the user presses `Del`. The proxy simply re-routes the `d` keypress through the same `List.HandleKey` path — no duplication of delete logic required.
:::

---

## Complete program

Here is the finished `main.go` with all steps assembled:

```go
package main

import (
	"log"
	"time"

	oat "github.com/antoniocali/oat-latte"
	"github.com/antoniocali/oat-latte/latte"
	"github.com/antoniocali/oat-latte/layout"
	"github.com/antoniocali/oat-latte/widget"
	"github.com/gdamore/tcell/v2"
)

type App struct {
	canvas *oat.Canvas
	list   *widget.List
	notifs *widget.NotificationManager
	items  []widget.ListItem
}

func (a *App) showConfirmDialog(msg string, onConfirm func()) {
	noBtn := widget.NewButton("No", func() {
		a.canvas.HideDialog()
	})
	yesBtn := widget.NewButton("Yes", func() {
		onConfirm()
		a.canvas.HideDialog()
	})

	btnRow := layout.NewHBox()
	btnRow.AddChild(layout.NewHFill())
	btnRow.AddChild(noBtn)
	btnRow.AddChild(layout.NewHFill().WithMaxSize(2))
	btnRow.AddChild(yesBtn)

	body := layout.NewPaddingUniform(
		layout.NewVBox(
			widget.NewText(msg),
			layout.NewVFill().WithMaxSize(1),
			btnRow,
		), 1)

	a.canvas.ShowDialog(
		widget.NewDialog("Confirm").
			WithChild(body).
			WithMaxSize(48, 9),
	)
}

func (a *App) showNewDialog() {
	input := widget.NewEditText().
		WithHint("Task name").
		WithPlaceholder("What needs doing?")

	doAdd := func() {
		name := input.GetText()
		if name == "" {
			return
		}
		a.items = append(a.items, widget.ListItem{Label: name})
		a.list.SetItems(a.items)
		a.canvas.HideDialog()
		a.notifs.Push("Task added", widget.NotificationKindSuccess, 2*time.Second)
	}

	input.WithOnSave(func(_ string) { doAdd() })

	cancelBtn := widget.NewButton("Cancel", func() { a.canvas.HideDialog() })
	addBtn    := widget.NewButton("Add",    doAdd)

	btnRow := layout.NewHBox()
	btnRow.AddChild(layout.NewHFill())
	btnRow.AddChild(cancelBtn)
	btnRow.AddChild(layout.NewHFill().WithMaxSize(2))
	btnRow.AddChild(addBtn)

	body := layout.NewPaddingUniform(
		layout.NewVBox(
			widget.NewText("Enter a name for the new task."),
			layout.NewVFill().WithMaxSize(1),
			input,
			layout.NewVFill().WithMaxSize(1),
			btnRow,
		), 1)

	a.canvas.ShowDialog(
		widget.NewDialog("New Task").
			WithChild(body).
			WithMaxSize(50, 13),
	)
}

// listProxy intercepts 'n' and 'd' before delegating to the list.
type listProxy struct {
	*widget.List
	app *App
}

func (p *listProxy) HandleKey(ev *oat.KeyEvent) bool {
	if ev.Key() == tcell.KeyRune {
		switch ev.Rune() {
		case 'n':
			p.app.showNewDialog()
			return true
		case 'd':
			return p.List.HandleKey(tcell.NewEventKey(tcell.KeyDelete, 0, tcell.ModNone))
		}
	}
	return p.List.HandleKey(ev)
}

func (p *listProxy) KeyBindings() []oat.KeyBinding {
	return append(
		[]oat.KeyBinding{
			{Key: tcell.KeyRune, Rune: 'n', Label: "n", Description: "New task"},
			{Key: tcell.KeyRune, Rune: 'd', Label: "d", Description: "Delete task"},
		},
		p.List.KeyBindings()...,
	)
}

func main() {
	a := &App{
		items: []widget.ListItem{
			{Label: "Buy groceries"},
			{Label: "Write tutorial"},
			{Label: "Ship v0.1.0"},
		},
	}

	a.list = widget.NewList(a.items).
		WithOnDelete(func(idx int, item widget.ListItem) {
			a.showConfirmDialog(
				"Delete \""+item.Label+"\"?",
				func() {
					a.items = append(a.items[:idx], a.items[idx+1:]...)
					a.list.SetItems(a.items)
					a.notifs.Push("Task deleted", widget.NotificationKindSuccess, 2*time.Second)
				},
			)
		})

	proxy     := &listProxy{List: a.list, app: a}
	body      := layout.NewBorder(proxy).WithTitle("Tasks")
	statusBar := widget.NewStatusBar()
	a.notifs   = widget.NewNotificationManager()

	a.canvas = oat.NewCanvas(
		oat.WithTheme(latte.ThemeDark),
		oat.WithBody(body),
		oat.WithAutoStatusBar(statusBar),
		oat.WithPrimary(proxy),
		oat.WithNotificationManager(a.notifs),
	)

	if err := a.canvas.Run(); err != nil {
		log.Fatal(err)
	}
}
```

---

## Run this example

The finished task-list app is included in the oat-latte repository. You can run it directly without cloning:

```sh
go run github.com/antoniocali/oat-latte/cmd/example/tasklist
```

Or, if you have the repo checked out:

```sh
make run-tasklist
```

---

## What's next

You now know how oat-latte's core pieces fit together. From here:

- Read [Core Concepts](../concepts) to understand the Measure/Render pipeline in depth.
- Read [Layout](../layout) for `Grid`, `Stack`, `Padding`, and flex sizing.
- Read [Focus](../focus) for the proxy pattern, custom `KeyBindings`, and `FocusByRef`.
- Read [Custom Components](../custom-components) to build your own widgets.
