// Command people is a "People Directory" TUI built with oat-latte.
// It demonstrates ComponentList with multi-column rows (name, role, status),
// a live-preview detail panel, and an "Add person" dialog.
//
// Layout:
//
//	╭─ People ─────────────────────────────╮ ╭─ Detail ─────────────────────╮
//	│ > Alice   Backend engineer   active   │ │ Name:   Alice                │
//	│   Bob     Frontend engineer  inactive │ │ Role:   Backend engineer     │
//	│   Charlie DevOps             active   │ │ Status: active               │
//	│                                       │ │                              │
//	╰───────────────────────────────────────╯ ╰──────────────────────────────╯
//	  ↑↓ Move   Enter Select   n New   Del Delete   ^T Theme   Tab Next
//
// Key bindings:
//   - Up / Down   Navigate the list
//   - n           Open the "New Person" dialog
//   - Del         Delete the selected person (with confirmation)
//   - ^T          Cycle through built-in themes
//   - Tab         Cycle focus
//   - Esc         Dismiss dialog / quit
package main

import (
	"fmt"
	"log"
	"time"

	oat "github.com/antoniocali/oat-latte"
	"github.com/antoniocali/oat-latte/latte"
	"github.com/antoniocali/oat-latte/layout"
	"github.com/antoniocali/oat-latte/widget"
	"github.com/gdamore/tcell/v2"
)

// ── domain ────────────────────────────────────────────────────────────────────

var nextID = 1

// Person is the application domain model.
type Person struct {
	ID     int
	Name   string
	Role   string
	Status string // "active" | "inactive"
}

func newPerson(name, role, status string) Person {
	p := Person{ID: nextID, Name: name, Role: role, Status: status}
	nextID++
	return p
}

func seed() []Person {
	return []Person{
		newPerson("Alice", "Backend engineer", "active"),
		newPerson("Bob", "Frontend engineer", "inactive"),
		newPerson("Charlie", "DevOps", "active"),
		newPerson("Diana", "Product manager", "active"),
		newPerson("Eve", "Data scientist", "inactive"),
	}
}

// ── app ───────────────────────────────────────────────────────────────────────

// App holds all application state and widget references.
type App struct {
	canvas *oat.Canvas
	list   *listProxy
	detail *widget.Text
	notifs *widget.NotificationManager
	people []Person
}

// makeRow builds a ComponentListItem for a Person.
// The row is an HBox: name (16-char padded) | role (flex) | status (coloured).
func makeRow(p Person) widget.ComponentListItem {
	statusStyle := latte.Style{FG: latte.ColorGreen}
	if p.Status != "active" {
		statusStyle = latte.Style{FG: latte.ColorBrightBlack}
	}

	row := layout.NewHBox(
		widget.NewText(fmt.Sprintf("%-16s", p.Name)),
		layout.NewFlexChild(widget.NewText(p.Role), 1),
		widget.NewText(p.Status).WithStyle(statusStyle),
	)
	return widget.ComponentListItem{Component: row, Value: p.ID}
}

func (a *App) rebuildItems() []widget.ComponentListItem {
	items := make([]widget.ComponentListItem, len(a.people))
	for i, p := range a.people {
		items[i] = makeRow(p)
	}
	return items
}

func (a *App) personByID(id int) (Person, int, bool) {
	for i, p := range a.people {
		if p.ID == id {
			return p, i, true
		}
	}
	return Person{}, -1, false
}

func (a *App) refreshDetail(p Person) {
	a.detail.SetText(fmt.Sprintf(
		"Name:   %s\nRole:   %s\nStatus: %s",
		p.Name, p.Role, p.Status,
	))
}

// ── dialogs ───────────────────────────────────────────────────────────────────

func (a *App) showNewDialog() {
	nameInput := widget.NewEditText().WithHint("Name").WithPlaceholder("Full name...")
	roleInput := widget.NewEditText().WithHint("Role").WithPlaceholder("Job title...")

	doCreate := func() {
		name := nameInput.GetText()
		role := roleInput.GetText()
		if name == "" {
			a.notifs.Push("Name cannot be empty", widget.NotificationKindError, 2*time.Second)
			return
		}
		p := newPerson(name, role, "active")
		a.people = append(a.people, p)
		a.list.SetItems(a.rebuildItems())
		a.canvas.HideDialog()
		a.notifs.Push("Added "+name, widget.NotificationKindSuccess, 2*time.Second)
	}

	nameInput.WithOnSave(func(_ string) { doCreate() })

	cancelBtn := widget.NewButton("Cancel", func() { a.canvas.HideDialog() })
	addBtn := widget.NewButton("Add", doCreate)

	btnRow := layout.NewHBox()
	btnRow.AddChild(layout.NewHFill())
	btnRow.AddChild(cancelBtn)
	btnRow.AddChild(layout.NewHFill().WithMaxSize(2))
	btnRow.AddChild(addBtn)

	body := layout.NewPaddingUniform(layout.NewVBox(
		nameInput,
		roleInput,
		layout.NewVFill().WithMaxSize(1),
		btnRow,
	), 1)

	a.canvas.ShowDialog(
		widget.NewDialog("New Person").
			WithChild(body).
			WithMaxSize(50, 13),
	)
}

func (a *App) showDeleteDialog(p Person, personIdx int) {
	cancelBtn := widget.NewButton("Cancel", func() { a.canvas.HideDialog() })
	deleteBtn := widget.NewButton("Delete", func() {
		a.people = append(a.people[:personIdx], a.people[personIdx+1:]...)
		a.list.SetItems(a.rebuildItems())
		a.detail.SetText("")
		a.canvas.HideDialog()
		a.notifs.Push("Deleted "+p.Name, widget.NotificationKindWarning, 2*time.Second)
	})

	btnRow := layout.NewHBox()
	btnRow.AddChild(layout.NewHFill())
	btnRow.AddChild(cancelBtn)
	btnRow.AddChild(layout.NewHFill().WithMaxSize(2))
	btnRow.AddChild(deleteBtn)

	body := layout.NewPaddingUniform(layout.NewVBox(
		widget.NewText(fmt.Sprintf("Delete %q?", p.Name)),
		widget.NewText("This action cannot be undone."),
		layout.NewVFill().WithMaxSize(1),
		btnRow,
	), 1)

	a.canvas.ShowDialog(
		widget.NewDialog("Confirm Delete").
			WithChild(body).
			WithMaxSize(50, 9),
	)
}

// ── proxy ─────────────────────────────────────────────────────────────────────

// listProxy wraps ComponentList to intercept the 'n' shortcut for adding a
// new person without modifying the widget itself.
type listProxy struct {
	*widget.ComponentList
	app *App
}

func (p *listProxy) HandleKey(ev *oat.KeyEvent) bool {
	if ev.Key() == tcell.KeyRune && ev.Rune() == 'n' {
		p.app.showNewDialog()
		return true
	}
	return p.ComponentList.HandleKey(ev)
}

func (p *listProxy) KeyBindings() []oat.KeyBinding {
	extra := []oat.KeyBinding{
		{Key: tcell.KeyRune, Rune: 'n', Label: "n", Description: "New person"},
	}
	return append(extra, p.ComponentList.KeyBindings()...)
}

// ── build ─────────────────────────────────────────────────────────────────────

func (a *App) build() {
	a.people = seed()
	a.notifs = widget.NewNotificationManager()

	// ── ComponentList ─────────────────────────────────────────────────────────

	cl := widget.NewComponentList(a.rebuildItems()).WithID("people")

	a.detail = widget.NewText("")

	// Live-preview: update the detail panel as the cursor moves.
	cl.WithOnCursorChange(func(_ int, item widget.ComponentListItem) {
		if id, ok := item.Value.(int); ok {
			if person, _, ok := a.personByID(id); ok {
				a.refreshDetail(person)
			}
		}
	})

	// Delete with confirmation dialog.
	cl.WithOnDelete(func(_ int, item widget.ComponentListItem) {
		if id, ok := item.Value.(int); ok {
			if person, personIdx, ok := a.personByID(id); ok {
				a.showDeleteDialog(person, personIdx)
			}
		}
	})

	a.list = &listProxy{ComponentList: cl, app: a}

	// Seed the detail panel with the first person.
	if len(a.people) > 0 {
		a.refreshDetail(a.people[0])
	}

	// ── Layout: list (2/3) + detail (1/3) ────────────────────────────────────

	listPanel := layout.NewBorder(a.list).WithTitle("People")
	detailPanel := layout.NewBorder(layout.NewPaddingUniform(a.detail, 1)).WithTitle("Detail")

	body := layout.NewHBox()
	body.AddFlexChild(listPanel, 2)
	body.AddFlexChild(detailPanel, 1)

	// ── Canvas ────────────────────────────────────────────────────────────────

	statusBar := widget.NewStatusBar()

	themes := []latte.Theme{
		latte.ThemeDark,
		latte.ThemeLight,
		latte.ThemeDracula,
		latte.ThemeNord,
	}
	themeIdx := 0

	a.canvas = oat.NewCanvas(
		oat.WithTheme(themes[themeIdx]),
		oat.WithBody(body),
		oat.WithAutoStatusBar(statusBar),
		oat.WithPrimary(a.list),
		oat.WithNotificationManager(a.notifs),
		oat.WithGlobalKeyBinding(oat.KeyBinding{
			Key:         tcell.KeyCtrlT,
			Mod:         tcell.ModCtrl,
			Label:       "^T",
			Description: "Toggle theme",
			Handler: func() {
				themeIdx = (themeIdx + 1) % len(themes)
				a.canvas.SetTheme(themes[themeIdx])
			},
		}),
	)

	a.notifs.Push("Welcome! Press [n] to add a person.", widget.NotificationKindInfo, 3*time.Second)
}

func main() {
	a := &App{}
	a.build()
	if err := a.canvas.Run(); err != nil {
		log.Fatal(err)
	}
}
