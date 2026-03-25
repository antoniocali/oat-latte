// Command kanban is a Kanban board TUI built with oat-latte.
//
// Layout:
//
//	╭─────────────────────────────────────────────────────────────╮
//	│  ◆ oat-latte  ·  Kanban Board               Tab · Esc·Quit  │
//	│  Done  ████████████░░░░░░░░  3 / 7                          │
//	╰─────────────────────────────────────────────────────────────╯
//	╭─ Backlog ────╮ ╭─ In Progress ╮ ╭─ Done ──────────────────╮
//	│ > Task A     │ │ > Task C     │ │ > Task E                 │
//	│   Task B     │ │   Task D     │ │   Task F                 │
//	│              │ │              │ │                          │
//	╰──────────────╯ ╰──────────────╯ ╰──────────────────────────╯
//	  [n] New  [Enter] View  [⇧→] Move right  [⇧←] Move left  [←/→] Switch column  [Tab] Next
//
// Shortcuts:
//   - n         New task dialog
//   - Enter     View/edit task dialog
//   - Shift+→   Move focused task to next column
//   - Shift+←   Move focused task to previous column
//   - ←/→       Switch focus to the previous/next column
//   - Del       Delete task (with confirmation dialog)
//   - Tab       Cycle focus between columns
//   - Esc       Dismiss dialog / quit
package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	oat "github.com/antoniocali/oat-latte"
	"github.com/antoniocali/oat-latte/latte"
	"github.com/antoniocali/oat-latte/layout"
	"github.com/antoniocali/oat-latte/widget"
	"github.com/gdamore/tcell/v2"
)

// ── domain ────────────────────────────────────────────────────────────────────

type ColumnID int

const (
	ColBacklog ColumnID = iota
	ColInProgress
	ColDone
	numCols
)

func (c ColumnID) Title() string {
	switch c {
	case ColInProgress:
		return "In Progress"
	case ColDone:
		return "Done"
	default:
		return "Backlog"
	}
}

var nextCardID = 1

type Card struct {
	ID   int
	Name string
	Desc string
	Col  ColumnID
}

func newCard(name, desc string, col ColumnID) Card {
	c := Card{ID: nextCardID, Name: name, Desc: desc, Col: col}
	nextCardID++
	return c
}

// ── app ───────────────────────────────────────────────────────────────────────

type App struct {
	canvas    *oat.Canvas
	statusBar *widget.StatusBar
	notifs    *widget.NotificationManager

	cards []Card

	// one List per column; proxies wrap the lists for extra key bindings
	lists       [numCols]*widget.List
	proxies     [numCols]*columnProxy
	listPanels  [numCols]*layout.Border
	focusedCol  ColumnID
	progressBar *widget.ProgressBar
	progressTxt *widget.Text
}

func seedCards() []Card {
	return []Card{
		newCard("Define user stories", "Gather requirements from stakeholders.", ColBacklog),
		newCard("Set up CI pipeline", "Configure GitHub Actions for build and test.", ColBacklog),
		newCard("Design DB schema", "ERD for the v2 data model.", ColBacklog),
		newCard("Build auth service", "JWT-based login and refresh tokens.", ColInProgress),
		newCard("API gateway routing", "Configure Kong rules for v2 endpoints.", ColInProgress),
		newCard("Write unit tests", "Cover all service-layer functions.", ColDone),
		newCard("Migrate staging DB", "Run schema migrations on staging.", ColDone),
	}
}

// ── helpers ───────────────────────────────────────────────────────────────────

func (a *App) cardsInCol(col ColumnID) []Card {
	var out []Card
	for _, c := range a.cards {
		if c.Col == col {
			out = append(out, c)
		}
	}
	return out
}

func (a *App) listItems(col ColumnID) []widget.ListItem {
	cards := a.cardsInCol(col)
	items := make([]widget.ListItem, len(cards))
	for i, c := range cards {
		items[i] = widget.ListItem{Label: c.Name, Value: c.ID}
	}
	return items
}

func (a *App) refreshAll() {
	for col := ColumnID(0); col < numCols; col++ {
		a.lists[col].SetItems(a.listItems(col))
	}
	a.refreshProgress()
}

func (a *App) refreshProgress() {
	total := len(a.cards)
	done := len(a.cardsInCol(ColDone))
	pct := 0.0
	if total > 0 {
		pct = float64(done) / float64(total)
	}
	a.progressBar.SetValue(pct)
	a.progressTxt.SetText(fmt.Sprintf(" %d / %d done", done, total))
}

func (a *App) cardByID(id int) (Card, int, bool) {
	for i, c := range a.cards {
		if c.ID == id {
			return c, i, true
		}
	}
	return Card{}, -1, false
}

func (a *App) selectedCardInCol(col ColumnID) (Card, bool) {
	item, ok := a.lists[col].SelectedItem()
	if !ok {
		return Card{}, false
	}
	id, ok := item.Value.(int)
	if !ok {
		return Card{}, false
	}
	c, _, ok := a.cardByID(id)
	return c, ok
}

// ── dialogs ───────────────────────────────────────────────────────────────────

func (a *App) showNewCardDialog() {
	nameInput := widget.NewEditText().
		WithPlaceholder("Card name…").
		WithHint("Name")
	descInput := widget.NewMultiLineEditText().
		WithHint("Description")

	var dlg *widget.Dialog

	doCreate := func() {
		name := strings.TrimSpace(nameInput.GetText())
		if name == "" {
			a.notifs.Push("Name cannot be empty", widget.NotificationKindError, 2*time.Second)
			return
		}
		desc := strings.TrimSpace(descInput.GetText())
		a.cards = append(a.cards, newCard(name, desc, a.focusedCol))
		a.refreshAll()
		a.canvas.HideDialog()
		a.notifs.Push("Card created: "+name, widget.NotificationKindSuccess, 2*time.Second)
	}

	cancelBtn := widget.NewButton("Cancel", func() {
		a.canvas.HideDialog()
	})
	createBtn := widget.NewButton("Create", doCreate)

	nameInput.WithOnSave(func(_ string) { doCreate() })

	btnRow := layout.NewHBox()
	btnRow.AddChild(layout.NewHFill())
	btnRow.AddChild(cancelBtn)
	btnRow.AddChild(layout.NewHFill().WithMaxSize(2))
	btnRow.AddChild(createBtn)

	body := layout.NewPaddingUniform(layout.NewVBox(
		nameInput,
		descInput,
		layout.NewVFill().WithMaxSize(1),
		btnRow,
	), 1)

	dlg = widget.NewDialog("New Card").
		WithChild(body).
		WithSize(widget.DialogFixed(52), widget.DialogPercent(60))
	_ = dlg

	a.canvas.ShowDialog(dlg)
}

func (a *App) showViewDialog(card Card) {
	nameInput := widget.NewEditText().WithHint("Name")
	nameInput.SetText(card.Name)
	descInput := widget.NewMultiLineEditText().WithHint("Description")
	descInput.SetText(card.Desc)

	var dlg *widget.Dialog

	cancelBtn := widget.NewButton("Cancel", func() {
		a.canvas.HideDialog()
	})
	saveBtn := widget.NewButton("Save", func() {
		name := strings.TrimSpace(nameInput.GetText())
		if name == "" {
			a.notifs.Push("Name cannot be empty", widget.NotificationKindError, 2*time.Second)
			return
		}
		_, idx, ok := a.cardByID(card.ID)
		if ok {
			a.cards[idx].Name = name
			a.cards[idx].Desc = strings.TrimSpace(descInput.GetText())
		}
		a.refreshAll()
		a.canvas.HideDialog()
		a.notifs.Push("Card saved", widget.NotificationKindSuccess, 2*time.Second)
	})

	colLabel := widget.NewText("Column: " + card.Col.Title())

	btnRow := layout.NewHBox()
	btnRow.AddChild(colLabel)
	btnRow.AddChild(layout.NewHFill())
	btnRow.AddChild(cancelBtn)
	btnRow.AddChild(layout.NewHFill().WithMaxSize(2))
	btnRow.AddChild(saveBtn)

	body := layout.NewPaddingUniform(layout.NewVBox(
		nameInput,
		descInput,
		layout.NewVFill().WithMaxSize(1),
		btnRow,
	), 1)

	dlg = widget.NewDialog("Edit Card").
		WithChild(body).
		WithMaxSize(60, 16)
	_ = dlg

	a.canvas.ShowDialog(dlg)
}

func (a *App) showDeleteDialog(card Card) {
	var dlg *widget.Dialog

	cancelBtn := widget.NewButton("Cancel", func() {
		a.canvas.HideDialog()
	})
	deleteBtn := widget.NewButton("Delete", func() {
		_, idx, ok := a.cardByID(card.ID)
		if ok {
			a.cards = append(a.cards[:idx], a.cards[idx+1:]...)
		}
		a.refreshAll()
		a.canvas.HideDialog()
		a.notifs.Push("Card deleted", widget.NotificationKindWarning, 2*time.Second)
	})

	msg := widget.NewText(fmt.Sprintf("Delete \"%s\"?", card.Name))
	hint := widget.NewText("This action cannot be undone.")

	btnRow := layout.NewHBox()
	btnRow.AddChild(layout.NewHFill())
	btnRow.AddChild(cancelBtn)
	btnRow.AddChild(layout.NewHFill().WithMaxSize(2))
	btnRow.AddChild(deleteBtn)

	body := layout.NewPaddingUniform(layout.NewVBox(
		msg, hint,
		layout.NewVFill().WithMaxSize(1),
		btnRow,
	), 1)

	dlg = widget.NewDialog("Confirm Delete").
		WithChild(body).
		WithMaxSize(50, 9)
	_ = dlg

	a.canvas.ShowDialog(dlg)
}

// ── build ─────────────────────────────────────────────────────────────────────

func (a *App) build() {
	a.statusBar = widget.NewStatusBar()
	a.notifs = widget.NewNotificationManager()
	a.cards = seedCards()

	// ── columns ──────────────────────────────────────────────────────────────

	var colPanels []oat.Component

	for col := ColumnID(0); col < numCols; col++ {
		col := col // capture

		list := widget.NewList(a.listItems(col)).
			WithID(fmt.Sprintf("col-%d", col))

		list.WithOnSelect(func(_ int, item widget.ListItem) {
			if id, ok := item.Value.(int); ok {
				if c, _, ok := a.cardByID(id); ok {
					a.showViewDialog(c)
				}
			}
		})
		list.WithOnDelete(func(_ int, item widget.ListItem) {
			if id, ok := item.Value.(int); ok {
				if c, _, ok := a.cardByID(id); ok {
					a.showDeleteDialog(c)
				}
			}
		})

		a.lists[col] = list

		panel := layout.NewBorder(list).
			WithTitle(col.Title())
		a.listPanels[col] = panel
		colPanels = append(colPanels, panel)
	}

	// ── body: 3 equal columns ────────────────────────────────────────────────

	body := layout.NewHBox()
	for _, p := range colPanels {
		body.AddFlexChild(p, 1)
	}

	// ── header: title + progress bar ─────────────────────────────────────────

	a.progressBar = widget.NewProgressBar().
		WithFillChar('█').
		WithEmptyChar('░').
		WithShowPercent(false)
	a.progressTxt = widget.NewText("")
	a.refreshProgress()

	titleTxt := widget.NewText("  ◆ oat-latte  ·  Kanban Board")
	hintTxt := widget.NewText("n·New  Enter·View  ⇧→·Move right  ⇧←·Move left  ←/→·Switch  Tab·Next  Esc·Quit  ")

	topRow := layout.NewHBox(titleTxt)
	topRow.AddChild(layout.NewHFill())
	topRow.AddChild(hintTxt)

	progressLabel := widget.NewText("  Progress  ")
	progressRow := layout.NewHBox(progressLabel)
	progressRow.AddFlexChild(a.progressBar, 1)
	progressRow.AddChild(a.progressTxt)

	headerBox := layout.NewVBox(
		topRow,
		progressRow,
		widget.NewText(strings.Repeat("─", 200)),
	)

	// ── wrap list keys (must happen before canvas so proxies are in the tree) ──

	for col := ColumnID(0); col < numCols; col++ {
		col := col
		a.wrapListKeys(a.lists[col], col)
	}

	// ── canvas ───────────────────────────────────────────────────────────────

	a.canvas = oat.NewCanvas(
		oat.WithTheme(latte.ThemeDark),
		oat.WithHeader(headerBox),
		oat.WithBody(body),
		oat.WithAutoStatusBar(a.statusBar),
		oat.WithPrimary(a.proxies[ColBacklog]),
		oat.WithNotificationManager(a.notifs),
	)

	a.notifs.Push("Welcome to Kanban! Press [n] to create a card.", widget.NotificationKindInfo, 4*time.Second)
}

// wrapListKeys wraps list in a columnProxy that intercepts extra keys
// ('n' for new card, →/← for move between columns) and stores the proxy
// on a.proxies[col] so it can be referenced by WithPrimary.
func (a *App) wrapListKeys(list *widget.List, col ColumnID) {
	moveRight := func() {
		if col >= ColDone {
			a.notifs.Push("Already in the last column", widget.NotificationKindWarning, 2*time.Second)
			return
		}
		card, ok := a.selectedCardInCol(col)
		if !ok {
			return
		}
		_, idx, ok := a.cardByID(card.ID)
		if !ok {
			return
		}
		a.cards[idx].Col = col + 1
		a.refreshAll()
		a.notifs.Push(fmt.Sprintf("Moved to %s", (col+1).Title()), widget.NotificationKindInfo, 2*time.Second)
	}

	moveLeft := func() {
		if col <= ColBacklog {
			a.notifs.Push("Already in the first column", widget.NotificationKindWarning, 2*time.Second)
			return
		}
		card, ok := a.selectedCardInCol(col)
		if !ok {
			return
		}
		_, idx, ok := a.cardByID(card.ID)
		if !ok {
			return
		}
		a.cards[idx].Col = col - 1
		a.refreshAll()
		a.notifs.Push(fmt.Sprintf("Moved to %s", (col-1).Title()), widget.NotificationKindInfo, 2*time.Second)
	}

	newCard := func() {
		a.focusedCol = col
		a.showNewCardDialog()
	}

	// Re-wire onSelect/onDelete so they go through the app dialogs.
	list.WithOnSelect(func(_ int, item widget.ListItem) {
		if id, ok := item.Value.(int); ok {
			if c, _, ok := a.cardByID(id); ok {
				a.showViewDialog(c)
			}
		}
	})

	// Use a column proxy that wraps the list and intercepts extra keys.
	proxy := &columnProxy{
		List:      list,
		col:       col,
		app:       a,
		moveRight: moveRight,
		moveLeft:  moveLeft,
		newCard:   newCard,
	}
	a.proxies[col] = proxy

	// Replace the list in the column's border child.
	a.listPanels[col].AddChild(proxy)
}

// columnProxy wraps a *widget.List and intercepts extra keys.
type columnProxy struct {
	*widget.List
	col       ColumnID
	app       *App
	moveRight func()
	moveLeft  func()
	newCard   func()
}

// HandleKey intercepts extra bindings then delegates to the underlying list.
func (p *columnProxy) HandleKey(ev *oat.KeyEvent) bool {
	switch {
	case ev.Key() == tcell.KeyRune && ev.Rune() == 'n':
		p.newCard()
		return true
	case ev.Key() == tcell.KeyRune && ev.Rune() == 'd':
		return p.List.HandleKey(tcell.NewEventKey(tcell.KeyDelete, 0, tcell.ModNone))
	case ev.Key() == tcell.KeyRight && ev.Modifiers() == tcell.ModShift:
		p.moveRight()
		return true
	case ev.Key() == tcell.KeyLeft && ev.Modifiers() == tcell.ModShift:
		p.moveLeft()
		return true
	}
	return p.List.HandleKey(ev)
}

// KeyBindings merges extra hints into the list's existing bindings.
func (p *columnProxy) KeyBindings() []oat.KeyBinding {
	base := p.List.KeyBindings()
	extra := []oat.KeyBinding{
		{Key: tcell.KeyRune, Rune: 'n', Label: "n", Description: "New card"},
		{Key: tcell.KeyRune, Rune: 'd', Label: "d", Description: "Delete card"},
		{Key: tcell.KeyRight, Mod: tcell.ModShift, Label: "⇧→", Description: "Move right"},
		{Key: tcell.KeyLeft, Mod: tcell.ModShift, Label: "⇧←", Description: "Move left"},
		// Display-only hints: Left/Right column switching is handled by canvas
		// as a fallback when the focused component does not consume the key.
		{Key: tcell.KeyLeft, Label: "←/→", Description: "Switch column"},
	}
	return append(extra, base...)
}

func main() {
	a := &App{}
	a.build()
	if err := a.canvas.Run(); err != nil {
		log.Fatal(err)
	}
}
