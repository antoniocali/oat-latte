// Command notes is a Note Manager TUI built with oat-latte.
//
// Layout:
//
//	╭─────────────────────────────────────────────────────────────╮
//	│  ◆ oat-latte  ·  Notes                       Tab · Esc·Quit │
//	╰─────────────────────────────────────────────────────────────╯
//	╭─ Notes (4) ──╮  ╭─ Editor ────────────────────────────────╮
//	│ > My note    │  │  Title                                   │
//	│   Shopping   │  │  My note title                           │
//	│   Meeting    │  │  Description                             │
//	│   Ideas      │  │  Note body text here…                    │
//	│              │  │                                          │
//	│              │  │  Tags                                    │
//	╰──────────────╯  │  backend, api                            │
//	                  ╰──────────────────────────────────────────╯
//	  [↑] Up  [↓] Down  [n] New  [e] Edit  [Del] Delete  [Tab] Next
//
// Shortcuts:
//   - n         Create new note (from list panel)
//   - e         Jump to editor / title field (from list or editor panel)
//   - Del       Delete note (with confirmation)
//   - ^S        Save note (from any editor field)
//   - ^G        Discard changes (from any editor field)
//   - Tab       Cycle focus between panels
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

var nextNoteID = 1

type Note struct {
	ID    int
	Title string
	Body  string
	Tags  []string
}

func newNote(title, body string, tags []string) Note {
	n := Note{ID: nextNoteID, Title: title, Body: body, Tags: tags}
	nextNoteID++
	return n
}

func seedNotes() []Note {
	return []Note{
		newNote(
			"Project kickoff",
			"Discuss scope, timeline and responsibilities with the whole team.\nAgenda: introductions, roadmap walk-through, Q&A.",
			[]string{"work", "meeting"},
		),
		newNote(
			"Grocery list",
			"Milk\nEggs\nBread\nCheese\nApples\nPasta",
			[]string{"personal", "shopping"},
		),
		newNote(
			"Go concurrency notes",
			"Use errgroup for coordinated goroutines.\nAvoid goroutine leaks — always select on ctx.Done().\nPrefer buffered channels when sender must not block.",
			[]string{"go", "programming"},
		),
		newNote(
			"Book recommendations",
			"- The Pragmatic Programmer\n- Clean Code\n- Designing Data-Intensive Applications\n- Site Reliability Engineering",
			[]string{"books", "personal"},
		),
		newNote(
			"Sprint retrospective",
			"What went well:\n- CI pipeline stabilised\n- Good PR review turnaround\n\nWhat to improve:\n- More thorough manual QA before release\n- Reduce WIP limit",
			[]string{"work", "agile"},
		),
	}
}

// ── app ───────────────────────────────────────────────────────────────────────

type App struct {
	canvas    *oat.Canvas
	statusBar *widget.StatusBar
	notifs    *widget.NotificationManager

	notes         []Note
	activeNoteID  int // ID of note currently loaded in editor, 0 = none
	unsavedChange bool
	editorMode    bool // true while the right panel is being edited

	// left panel
	noteList  *noteListProxy
	rawList   *widget.List
	leftPanel *layout.Border

	// right editor panel
	editorPanel  *layout.Border
	editorShim   *editorFocusShim
	titleInput   *widget.EditText
	titleGuard   *editorInputGuard
	tagsLabel    *widget.Label
	tagsInput    *widget.EditText // comma-separated tags
	tagsSwitcher *tagsSwitcher
	bodyInput    *widget.EditText
	bodyGuard    *editorInputGuard
}

func (a *App) noteByID(id int) (Note, int, bool) {
	for i, n := range a.notes {
		if n.ID == id {
			return n, i, true
		}
	}
	return Note{}, -1, false
}

func (a *App) listItems() []widget.ListItem {
	items := make([]widget.ListItem, len(a.notes))
	for i, n := range a.notes {
		label := n.Title
		if len(n.Tags) > 0 {
			label += "  [" + strings.Join(n.Tags, ", ") + "]"
		}
		items[i] = widget.ListItem{Label: label, Value: n.ID}
	}
	return items
}

func (a *App) refreshList() {
	a.rawList.SetItems(a.listItems())
	if a.leftPanel != nil {
		a.leftPanel.WithTitle(fmt.Sprintf("Notes (%d)", len(a.notes)))
	}
}

func (a *App) setEditorMode(on bool) {
	if a.editorMode == on {
		return
	}
	a.editorMode = on
	if a.canvas != nil {
		a.canvas.InvalidateLayout()
	}
}

// loadNote populates the right-panel editor with the given note.
func (a *App) loadNote(note Note) {
	a.activeNoteID = note.ID
	a.unsavedChange = false
	a.setEditorMode(false)
	a.titleInput.SetText(note.Title)
	a.bodyInput.SetText(note.Body)
	a.tagsInput.SetText(strings.Join(note.Tags, ", "))
	a.tagsLabel.SetLabels(note.Tags)
}

func (a *App) clearEditor() {
	a.activeNoteID = 0
	a.unsavedChange = false
	a.setEditorMode(false)
	a.titleInput.SetText("")
	a.bodyInput.SetText("")
	a.tagsInput.SetText("")
	a.tagsLabel.SetLabels(nil)
}

func (a *App) saveActiveNote() {
	if a.activeNoteID == 0 {
		return
	}
	title := strings.TrimSpace(a.titleInput.GetText())
	if title == "" {
		a.notifs.Push("Title cannot be empty", widget.NotificationKindError, 2*time.Second)
		return
	}
	body := a.bodyInput.GetText()
	tagsRaw := strings.Split(a.tagsInput.GetText(), ",")
	var tags []string
	for _, t := range tagsRaw {
		t = strings.TrimSpace(t)
		if t != "" {
			tags = append(tags, t)
		}
	}
	_, idx, ok := a.noteByID(a.activeNoteID)
	if !ok {
		return
	}
	a.notes[idx].Title = title
	a.notes[idx].Body = body
	a.notes[idx].Tags = tags
	a.unsavedChange = false
	a.setEditorMode(false)
	a.tagsLabel.SetLabels(tags)
	a.refreshList()
	a.canvas.FocusByRef(a.editorShim)
	a.notifs.Push("Note saved", widget.NotificationKindSuccess, 2*time.Second)
}

func (a *App) discardChanges() {
	if a.activeNoteID == 0 {
		return
	}
	note, _, ok := a.noteByID(a.activeNoteID)
	if ok {
		a.loadNote(note)
	}
	a.canvas.FocusByRef(a.editorShim)
	a.notifs.Push("Changes discarded", widget.NotificationKindWarning, 2*time.Second)
}

// ── new note dialog ───────────────────────────────────────────────────────────

func (a *App) showNewNoteDialog() {
	titleIn := widget.NewEditText().
		WithStyle(latte.Style{Border: latte.BorderExplicitNone}).
		WithHint("Title").
		WithPlaceholder("Note title…")

	descIn := widget.NewMultiLineEditText().
		WithStyle(latte.Style{Border: latte.BorderExplicitNone}).
		WithHint("Description").
		WithPlaceholder("Write your note here…")

	tagsIn := widget.NewEditText().
		WithStyle(latte.Style{Border: latte.BorderExplicitNone}).
		WithHint("Tags").
		WithPlaceholder("tag1, tag2, tag3…")

	doCreate := func() {
		title := strings.TrimSpace(titleIn.GetText())
		if title == "" {
			a.notifs.Push("Title cannot be empty", widget.NotificationKindError, 2*time.Second)
			return
		}
		body := descIn.GetText()
		tagsRaw := strings.Split(tagsIn.GetText(), ",")
		var tags []string
		for _, t := range tagsRaw {
			t = strings.TrimSpace(t)
			if t != "" {
				tags = append(tags, t)
			}
		}
		note := newNote(title, body, tags)
		a.notes = append(a.notes, note)
		a.refreshList()
		a.loadNote(note)
		a.canvas.HideDialog()
		a.notifs.Push("Note created: "+title, widget.NotificationKindSuccess, 2*time.Second)
	}

	doCancel := func() { a.canvas.HideDialog() }

	cancelBtn := widget.NewButton("Cancel", doCancel)
	createBtn := widget.NewButton("Create", doCreate)

	// ^S / ^G wire-up on all three fields.
	titleIn.WithOnSave(func(_ string) { doCreate() })
	titleIn.WithOnCancel(doCancel)
	descIn.WithOnSave(func(_ string) { doCreate() })
	descIn.WithOnCancel(doCancel)
	tagsIn.WithOnSave(func(_ string) { doCreate() })
	tagsIn.WithOnCancel(doCancel)

	btnRow := layout.NewHBox()
	btnRow.AddChild(layout.NewHFill())
	btnRow.AddChild(cancelBtn)
	btnRow.AddChild(layout.NewHFill().WithMaxSize(2))
	btnRow.AddChild(createBtn)

	dialogVBox := layout.NewVBox(titleIn)
	dialogVBox.AddFlexChild(descIn, 1)
	dialogVBox.AddChild(tagsIn)
	dialogVBox.AddChild(layout.NewVFill().WithMaxSize(1))
	dialogVBox.AddChild(btnRow)

	dlgBody := layout.NewPaddingUniform(dialogVBox, 1)

	dlg := widget.NewDialog("New Note").
		WithChild(dlgBody).
		WithSize(widget.DialogPercent(50), widget.DialogPercent(60))

	a.canvas.ShowDialog(dlg)
}

// ── delete confirmation ───────────────────────────────────────────────────────

func (a *App) showDeleteDialog(note Note) {
	var dlg *widget.Dialog

	cancelBtn := widget.NewButton("Cancel", func() {
		a.canvas.HideDialog()
	})
	deleteBtn := widget.NewButton("Delete", func() {
		_, idx, ok := a.noteByID(note.ID)
		if ok {
			a.notes = append(a.notes[:idx], a.notes[idx+1:]...)
		}
		if a.activeNoteID == note.ID {
			a.clearEditor()
		}
		a.refreshList()
		a.canvas.HideDialog()
		a.notifs.Push("Note deleted", widget.NotificationKindWarning, 2*time.Second)
	})

	msg := widget.NewText(fmt.Sprintf("Delete \"%s\"?", note.Title))
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
		WithMaxSize(52, 9)

	a.canvas.ShowDialog(dlg)
}

// ── build ─────────────────────────────────────────────────────────────────────

func (a *App) build() {
	a.statusBar = widget.NewStatusBar()
	a.notifs = widget.NewNotificationManager()
	a.notes = seedNotes()

	// ── left panel ───────────────────────────────────────────────────────────

	rawList := widget.NewList(a.listItems()).WithID("note-list")
	rawList.WithOnCursorChange(func(_ int, item widget.ListItem) {
		if id, ok := item.Value.(int); ok {
			if n, _, ok := a.noteByID(id); ok {
				a.loadNote(n)
			}
		}
	})
	rawList.WithOnSelect(func(_ int, item widget.ListItem) {
		if id, ok := item.Value.(int); ok {
			if n, _, ok := a.noteByID(id); ok {
				a.loadNote(n)
			}
		}
	})
	rawList.WithOnDelete(func(_ int, item widget.ListItem) {
		if id, ok := item.Value.(int); ok {
			if n, _, ok := a.noteByID(id); ok {
				a.showDeleteDialog(n)
			}
		}
	})
	a.rawList = rawList

	proxy := &noteListProxy{List: rawList, app: a}
	a.noteList = proxy

	leftVBox := layout.NewVBox()
	leftVBox.AddFlexChild(proxy, 1)

	leftPanel := layout.NewBorder(leftVBox).
		WithTitle(fmt.Sprintf("Notes (%d)", len(a.notes)))
	a.leftPanel = leftPanel

	// ── right editor panel ───────────────────────────────────────────────────────

	a.titleInput = widget.NewEditText().
		WithStyle(latte.Style{Border: latte.BorderExplicitNone}).
		WithID("note-title-input").
		WithHint("Title").
		WithPlaceholder("Note title…")
	a.titleInput.WithOnChange(func(_ string) { a.unsavedChange = true })
	a.titleGuard = &editorInputGuard{EditText: a.titleInput, app: a}

	a.tagsLabel = widget.NewLabel(nil)

	a.tagsInput = widget.NewEditText().
		WithStyle(latte.Style{Border: latte.BorderExplicitNone}).
		WithID("note-tags-input").
		WithHint("Tags").
		WithPlaceholder("tag1, tag2, tag3…")
	a.tagsInput.WithOnChange(func(text string) {
		a.unsavedChange = true
		parts := strings.Split(text, ",")
		var tags []string
		for _, t := range parts {
			t = strings.TrimSpace(t)
			if t != "" {
				tags = append(tags, t)
			}
		}
		a.tagsLabel.SetLabels(tags)
	})

	a.bodyInput = widget.NewMultiLineEditText().
		WithStyle(latte.Style{Border: latte.BorderExplicitNone}).
		WithID("note-body-input").
		WithHint("Description").
		WithPlaceholder("Write your note here…")
	a.bodyInput.WithOnChange(func(_ string) { a.unsavedChange = true })
	a.bodyGuard = &editorInputGuard{EditText: a.bodyInput, app: a}

	// Wire ^S / ^G on all inputs
	a.titleInput.WithOnSave(func(_ string) { a.saveActiveNote() })
	a.titleInput.WithOnCancel(func() { a.discardChanges() })
	a.bodyInput.WithOnSave(func(_ string) { a.saveActiveNote() })
	a.bodyInput.WithOnCancel(func() { a.discardChanges() })
	a.tagsInput.WithOnSave(func(_ string) { a.saveActiveNote() })
	a.tagsInput.WithOnCancel(func() { a.discardChanges() })

	a.editorShim = newEditorFocusShim(a)

	a.tagsSwitcher = &tagsSwitcher{app: a}
	a.tagsSwitcher.EnsureID()

	editorVBox := layout.NewVBox(
		a.editorShim,
		a.titleGuard,
	)
	editorVBox.AddFlexChild(a.bodyGuard, 1)
	editorVBox.AddChild(a.tagsSwitcher)

	a.editorPanel = layout.NewBorder(editorVBox).
		WithStyle(latte.Style{Padding: latte.Insets{Bottom: 1}}).
		WithTitle("Editor")

	// ── body ─────────────────────────────────────────────────────────────────

	body := layout.NewHBox()
	body.AddFlexChild(leftPanel, 1)
	body.AddFlexChild(a.editorPanel, 3)

	// ── header ───────────────────────────────────────────────────────────────

	titleWidget := widget.NewText("  ◆ oat-latte  ·  Notes")

	// ── canvas ───────────────────────────────────────────────────────────────

	a.canvas = oat.NewCanvas(
		oat.WithTheme(latte.ThemeDark),
		oat.WithHeader(titleWidget),
		oat.WithBody(body),
		oat.WithAutoStatusBar(a.statusBar),
		oat.WithPrimary(proxy),
	)

	a.notifs.SetNotifyChannel(a.canvas.NotifyChannel())
	a.canvas.ShowPersistentOverlay(a.notifs)

	// load first note into editor
	if len(a.notes) > 0 {
		a.loadNote(a.notes[0])
	}

	a.notifs.Push("Welcome to Notes! Press [n] to create a note.", widget.NotificationKindInfo, 4*time.Second)
}

// ── noteListProxy ─────────────────────────────────────────────────────────────

// noteListProxy wraps *widget.List and adds 'n' and 'e' shortcuts.
type noteListProxy struct {
	*widget.List
	app *App
}

func (p *noteListProxy) HandleKey(ev *oat.KeyEvent) bool {
	if ev.Key() == tcell.KeyRune {
		switch ev.Rune() {
		case 'n':
			p.app.showNewNoteDialog()
			return true
		case 'e':
			p.app.setEditorMode(true)
			p.app.canvas.FocusByRef(p.app.titleGuard)
			return true
		}
	}
	return p.List.HandleKey(ev)
}

func (p *noteListProxy) IsFocusable() bool {
	return !p.app.editorMode
}

func (p *noteListProxy) KeyBindings() []oat.KeyBinding {
	base := p.List.KeyBindings()
	extra := []oat.KeyBinding{
		{Key: tcell.KeyRune, Rune: 'n', Label: "n", Description: "New note"},
		{Key: tcell.KeyRune, Rune: 'e', Label: "e", Description: "Edit note"},
	}
	return append(extra, base...)
}

// ── editorFocusShim ───────────────────────────────────────────────────────────

// editorFocusShim is a zero-size focusable shim placed at the top of the editor
// panel.  When focused it shows the 'e' / Tab shortcuts; pressing 'e' or Enter
// moves focus to the titleInput so the user can start typing immediately.
type editorFocusShim struct {
	oat.BaseComponent
	oat.FocusBehavior
	app *App
}

func newEditorFocusShim(a *App) *editorFocusShim {
	s := &editorFocusShim{app: a}
	s.EnsureID()
	return s
}

func (s *editorFocusShim) Measure(_ oat.Constraint) oat.Size  { return oat.Size{} }
func (s *editorFocusShim) Render(_ *oat.Buffer, _ oat.Region) {}

// IsFocusable implements oat.FocusGuard: the shim is only in the focus cycle
// when NOT in editor mode (it acts as the right-panel representative in view mode).
func (s *editorFocusShim) IsFocusable() bool { return !s.app.editorMode }

func (s *editorFocusShim) HandleKey(ev *oat.KeyEvent) bool {
	if ev.Key() == tcell.KeyRune && ev.Rune() == 'e' {
		s.app.setEditorMode(true)
		s.app.canvas.FocusByRef(s.app.titleGuard)
		return true
	}
	if ev.Key() == tcell.KeyEnter {
		s.app.setEditorMode(true)
		s.app.canvas.FocusByRef(s.app.titleGuard)
		return true
	}
	return false
}

func (s *editorFocusShim) KeyBindings() []oat.KeyBinding {
	return []oat.KeyBinding{
		{Key: tcell.KeyRune, Rune: 'e', Label: "e", Description: "Edit note"},
		{Key: tcell.KeyEnter, Label: "Enter", Description: "Edit note"},
	}
}

// ── tagsSwitcher ──────────────────────────────────────────────────────────────

// tagsSwitcher is a thin component that renders tagsLabel when the app is in
// view mode, and tagsInput (an editable EditText) when in editor mode.
// Both children always exist in memory; the switcher simply delegates
// Measure/Render to the active one, keeping the focus tree stable.
type tagsSwitcher struct {
	oat.BaseComponent
	oat.FocusBehavior
	app *App
}

// IsFocusable implements oat.FocusGuard: the switcher is in the focus cycle
// only when in editor mode.
func (s *tagsSwitcher) IsFocusable() bool { return s.app.editorMode }

func (s *tagsSwitcher) active() oat.Component {
	if s.app.editorMode {
		return s.app.tagsInput
	}
	return s.app.tagsLabel
}

func (s *tagsSwitcher) Measure(c oat.Constraint) oat.Size {
	return s.active().Measure(c)
}

func (s *tagsSwitcher) Render(buf *oat.Buffer, region oat.Region) {
	s.active().Render(buf, region)
}

// SetFocused / IsFocused / HandleKey / KeyBindings forward to tagsInput only
// when in editor mode; otherwise the switcher is inert.

func (s *tagsSwitcher) SetFocused(focused bool) {
	if s.app.editorMode {
		s.app.tagsInput.SetFocused(focused)
	} else {
		s.FocusBehavior.SetFocused(focused)
	}
}

func (s *tagsSwitcher) IsFocused() bool {
	if s.app.editorMode {
		return s.app.tagsInput.IsFocused()
	}
	return s.FocusBehavior.IsFocused()
}

func (s *tagsSwitcher) HandleKey(ev *oat.KeyEvent) bool {
	if s.app.editorMode {
		return s.app.tagsInput.HandleKey(ev)
	}
	return false
}

func (s *tagsSwitcher) KeyBindings() []oat.KeyBinding {
	if s.app.editorMode {
		return s.app.tagsInput.KeyBindings()
	}
	return nil
}

func main() {
	a := &App{}
	a.build()
	if err := a.canvas.Run(); err != nil {
		log.Fatal(err)
	}
}

// ── editorInputGuard ──────────────────────────────────────────────────────────

// editorInputGuard wraps a *widget.EditText and implements oat.FocusGuard so
// the input is only included in the Tab cycle when the app is in editor mode.
type editorInputGuard struct {
	*widget.EditText
	app *App
}

// IsFocusable implements oat.FocusGuard.
func (g *editorInputGuard) IsFocusable() bool { return g.app.editorMode }
