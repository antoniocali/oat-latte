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
	addBtn := widget.NewButton("Add", doAdd)

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
			WithMaxSize(50, 11),
	)
}

// listProxy intercepts 'n' to open the new-task dialog.
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

	proxy := &listProxy{List: a.list, app: a}
	body := layout.NewBorder(proxy).WithTitle("Tasks")
	statusBar := widget.NewStatusBar()

	a.notifs = widget.NewNotificationManager()

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
