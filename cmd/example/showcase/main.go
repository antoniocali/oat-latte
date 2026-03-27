// Command showcase is a kitchen-sink TUI built with oat-latte.
// It does not model any real domain — its only purpose is to exercise every
// widget, layout container, style option, alignment and anchor the framework
// exposes, so that all features are visible and verifiable in one place.
//
// Layout (three-column body):
//
//	header: title centred + theme name right-aligned
//
//	┌── Left column ──────────────────────────────┐
//	│  Buttons (various HAlign)                   │
//	│  CheckBoxes                                 │
//	│  Dividers (H + V)                           │
//	│  ProgressBars (every anchor)                │
//	└─────────────────────────────────────────────┘
//	┌── Middle column ────────────────────────────┐
//	│  List  (top half, flex 1)                   │
//	│  ComponentList (bottom half, flex 1)        │
//	└─────────────────────────────────────────────┘
//	┌── Right column ─────────────────────────────┐
//	│  EditText  (borderless, with hints)         │
//	│  MultiLineEditText (flex, borderless)       │
//	│  Text (VAlign variants)                     │
//	│  AlignChild showcase                        │
//	│  Labels (tag chips)                         │
//	└─────────────────────────────────────────────┘
//
//	footer: StatusBar (auto key-hints)
//
// Global keys:
//   - ^T   cycle through all 5 built-in themes
//   - ^D   open a modal Dialog (shows AlignChild + Padding + Divider inside)
//   - ^N   push a random notification (cycles through all four kinds)
//   - Esc  quit
package main

import (
	"fmt"
	"log"
	"math"
	"time"

	oat "github.com/antoniocali/oat-latte"
	"github.com/antoniocali/oat-latte/latte"
	"github.com/antoniocali/oat-latte/layout"
	"github.com/antoniocali/oat-latte/widget"
	"github.com/gdamore/tcell/v2"
)

// ── app ───────────────────────────────────────────────────────────────────────

type App struct {
	canvas         *oat.Canvas
	notifs         *widget.NotificationManager
	themeIdx       int
	themeName      *widget.Text // updated whenever theme changes
	progress       [3]*widget.ProgressBar
	progVal        float64
	noBorders      bool
	roundedCorners bool
}

// applyDerivedTheme rebuilds the active theme from themes[a.themeIdx] and
// applies the current noBorders and roundedCorners flags on top, then calls
// SetTheme and updates the theme name label.  Call this whenever either flag
// or the base theme index changes.
func (a *App) applyDerivedTheme() {
	t := themes[a.themeIdx]
	name := themeNames[a.themeIdx]

	if a.noBorders {
		t = t.
			WithPanel(latte.Style{Border: latte.BorderExplicitNone}).
			WithInput(latte.Style{Border: latte.BorderExplicitNone}).
			WithButton(latte.Style{Border: latte.BorderExplicitNone}).
			WithDialog(latte.Style{Border: latte.BorderExplicitNone})
		name += " (borderless)"
	}
	t = t.WithRoundedCorner(a.roundedCorners)

	a.canvas.SetTheme(t)
	a.themeName.SetText("theme: " + name)
}

// ── header ───────────────────────────────────────────────────────────────────

func (a *App) buildHeader() oat.Component {
	title := widget.NewTitle("oat-latte  ·  showcase").
		WithStyle(latte.Style{Bold: true})

	a.themeName = widget.NewText("theme: dark").
		WithHAlign(oat.HAlignRight)

	hdr := layout.NewHBox(
		layout.NewFlexChild(title, 1),
		a.themeName,
	)
	return layout.NewPadding(hdr, oat.Insets{Top: 0, Bottom: 0, Left: 1, Right: 1})
}

// ── left column ──────────────────────────────────────────────────────────────

func (a *App) buildLeft() oat.Component {
	// ── buttons with every HAlign ─────────────────────────────────────────────
	btnFill := widget.NewButton("HAlignFill (default)", func() {
		a.notifs.Push("HAlignFill button pressed", widget.NotificationKindInfo, 2*time.Second)
	})
	btnLeft := widget.NewButton("HAlignLeft", func() {
		a.notifs.Push("HAlignLeft button pressed", widget.NotificationKindInfo, 2*time.Second)
	}).WithHAlign(oat.HAlignLeft)

	btnCenter := widget.NewButton("HAlignCenter", func() {
		a.notifs.Push("HAlignCenter button pressed", widget.NotificationKindSuccess, 2*time.Second)
	}).WithHAlign(oat.HAlignCenter)

	btnRight := widget.NewButton("HAlignRight", func() {
		a.notifs.Push("HAlignRight button pressed", widget.NotificationKindSuccess, 2*time.Second)
	}).WithHAlign(oat.HAlignRight)

	// WithRoundedCorner(true) on a dashed border is a silent no-op —
	// incompatible border styles keep their square corners without panicking.
	btnRoundedDashed := widget.NewButton("Rounded (dashed, no-op)", func() {
		a.notifs.Push("Dashed stays square — no panic!", widget.NotificationKindInfo, 2*time.Second)
	}).
		WithStyle(latte.Style{Border: latte.BorderDashed}).
		WithRoundedCorner(true).
		WithHAlign(oat.HAlignCenter)

	btnSection := layout.NewVBox(
		widget.NewText("─── Buttons ───────────────────").
			WithStyle(latte.Style{Bold: true}),
		btnFill,
		btnLeft,
		btnCenter,
		btnRight,
		btnRoundedDashed,
	)

	// ── checkboxes ────────────────────────────────────────────────────────────
	cb1 := widget.NewCheckBox("Enable notifications").
		WithOnToggle(func(on bool) {
			msg := "Notifications disabled"
			if on {
				msg = "Notifications enabled"
			}
			a.notifs.Push(msg, widget.NotificationKindInfo, 2*time.Second)
		})
	cb1.SetChecked(true)

	cb2 := widget.NewCheckBox("No borders").
		WithOnToggle(func(on bool) {
			a.noBorders = on
			a.applyDerivedTheme()
		})

	cb3 := widget.NewCheckBox("Rounded corners").
		WithOnToggle(func(on bool) {
			a.roundedCorners = on
			a.applyDerivedTheme()
		})
	cb3.SetChecked(true) // all built-in themes default to RoundedCorner: true

	cbSection := layout.NewVBox(
		widget.NewText("─── CheckBoxes ─────────────────").
			WithStyle(latte.Style{Bold: true}),
		cb1,
		cb2,
		cb3,
	)

	// ── dividers ──────────────────────────────────────────────────────────────
	// Horizontal dividers: full, 50% centred, 8-cell fixed left
	hdFull := widget.NewHDivider()
	hdHalf := widget.NewHDivider().
		WithMaxSize(widget.DividerPercent(50), oat.AnchorCenter)
	hdFixed := widget.NewHDivider().
		WithRune('═').
		WithMaxSize(widget.DividerFixed(12), oat.AnchorLeft)
	hdRight := widget.NewHDivider().
		WithRune('─').
		WithMaxSize(widget.DividerPercent(60), oat.AnchorRight)

	// A short HBox that embeds vertical dividers between fixed-width text items.
	// Wrapped in a Border so it gets an explicit 3-row height (top border +
	// 1 content row + bottom border), making the │ characters visible.
	// The text items use %-6s format so they report a fixed natural width when
	// measured; without this they would each claim the full available width and
	// leave the VDividers with zero columns.
	vdRow := layout.NewBorder(layout.NewHBox(
		widget.NewText(fmt.Sprintf("%-6s", "left")),
		widget.NewVDivider(),
		widget.NewText(fmt.Sprintf("%-8s", "middle")),
		widget.NewVDivider().
			WithMaxSizeV(widget.DividerFixed(1), oat.VAnchorMiddle),
		widget.NewText(fmt.Sprintf("%-6s", "right")),
	))

	divSection := layout.NewVBox(
		widget.NewText("─── Dividers ───────────────────").
			WithStyle(latte.Style{Bold: true}),
		widget.NewText("full:"),
		hdFull,
		widget.NewText("50% centred:"),
		hdHalf,
		widget.NewText("12-cell fixed double (left):"),
		hdFixed,
		widget.NewText("60% right:"),
		hdRight,
		widget.NewText("vertical (in HBox):"),
		vdRow,
	)

	// ── progress bars ─────────────────────────────────────────────────────────
	a.progress[0] = widget.NewProgressBar().
		WithPercentage(true, oat.AnchorLeft)
	a.progress[1] = widget.NewProgressBar().
		WithPercentage(true, oat.AnchorCenter)
	a.progress[2] = widget.NewProgressBar().
		WithPercentage(true, oat.AnchorRight)

	a.progress[0].SetValue(0.30)
	a.progress[1].SetValue(0.60)
	a.progress[2].SetValue(0.85)

	pbSection := layout.NewVBox(
		widget.NewText("─── ProgressBars ───────────────").
			WithStyle(latte.Style{Bold: true}),
		widget.NewText("label AnchorLeft:"),
		a.progress[0],
		widget.NewText("label AnchorCenter:"),
		a.progress[1],
		widget.NewText("label AnchorRight:"),
		a.progress[2],
	)

	col := layout.NewVBox(
		btnSection,
		layout.NewVFill().WithMaxSize(1),
		cbSection,
		layout.NewVFill().WithMaxSize(1),
		divSection,
		layout.NewVFill().WithMaxSize(1),
		pbSection,
	)
	return layout.NewBorder(
		layout.NewPadding(col, oat.Insets{Top: 0, Right: 1, Bottom: 0, Left: 1}),
	).WithTitle("Widgets", oat.AnchorLeft).WithRoundedCorner(true)
}

// ── middle column ─────────────────────────────────────────────────────────────

func (a *App) buildMiddle() oat.Component {
	// ── plain List ────────────────────────────────────────────────────────────
	listItems := []widget.ListItem{
		{Label: "Apples", Value: 1},
		{Label: "Bananas", Value: 2},
		{Label: "Cherries", Value: 3},
		{Label: "Dates", Value: 4},
		{Label: "Elderberries", Value: 5},
		{Label: "Figs", Value: 6},
		{Label: "Grapes", Value: 7},
	}
	list := widget.NewList(listItems).
		WithID("fruit-list").
		WithOnSelect(func(_ int, item widget.ListItem) {
			a.notifs.Push(
				fmt.Sprintf("Selected: %s", item.Label),
				widget.NotificationKindSuccess, 2*time.Second,
			)
		})

	listPanel := layout.NewBorder(list).
		WithTitle("List", oat.AnchorCenter).
		WithRoundedCorner(true)

	// ── ComponentList with styled rows ────────────────────────────────────────
	type entry struct{ name, kind, level string }
	entries := []entry{
		{"TCP/IP", "network", "info"},
		{"HTTP/2", "protocol", "ok"},
		{"TLS 1.3", "security", "ok"},
		{"DNS", "network", "warn"},
		{"WebSocket", "protocol", "info"},
		{"gRPC", "protocol", "ok"},
		{"QUIC", "network", "warn"},
	}

	levelStyle := func(l string) latte.Style {
		switch l {
		case "ok":
			return latte.Style{FG: latte.ColorGreen, Bold: true}
		case "warn":
			return latte.Style{FG: latte.ColorYellow, Bold: true}
		default:
			return latte.Style{FG: latte.ColorBrightCyan}
		}
	}

	compItems := make([]widget.ComponentListItem, len(entries))
	for i, e := range entries {
		row := layout.NewHBox(
			widget.NewText(fmt.Sprintf("%-12s", e.name)),
			layout.NewFlexChild(
				widget.NewText(e.kind).WithStyle(latte.Style{FG: latte.ColorBrightBlack}),
				1,
			),
			widget.NewText(e.level).WithStyle(levelStyle(e.level)),
		)
		compItems[i] = widget.ComponentListItem{Component: row, Value: i}
	}

	compList := widget.NewComponentList(compItems).
		WithID("proto-list").
		WithOnSelect(func(_ int, item widget.ComponentListItem) {
			a.notifs.Push(
				fmt.Sprintf("Protocol #%d selected", item.Value.(int)),
				widget.NotificationKindInfo, 2*time.Second,
			)
		})

	compPanel := layout.NewBorder(compList).
		WithTitle("ComponentList", oat.AnchorCenter).
		WithRoundedCorner(true)

	col := layout.NewVBox()
	col.AddFlexChild(listPanel, 1)
	col.AddFlexChild(compPanel, 1)
	return col
}

// ── right column ─────────────────────────────────────────────────────────────

func (a *App) buildRight() oat.Component {
	// ── EditText fields (borderless, hints) ───────────────────────────────────
	emailIn := widget.NewEditText().
		WithStyle(latte.Style{Border: latte.BorderExplicitNone}).
		WithHint("Email").
		WithPlaceholder("user@example.com")

	urlIn := widget.NewEditText().
		WithStyle(latte.Style{Border: latte.BorderExplicitNone}).
		WithHint("URL").
		WithPlaceholder("https://…")

	passIn := widget.NewEditText().
		WithStyle(latte.Style{Border: latte.BorderExplicitNone}).
		WithHint("Token (max 32 chars)").
		WithPlaceholder("••••••••••••••••").
		WithMaxLength(32)

	formPanel := layout.NewBorder(
		layout.NewPadding(
			layout.NewVBox(emailIn, urlIn, passIn),
			oat.Insets{Top: 0, Right: 1, Bottom: 1, Left: 1},
		),
	).WithTitle("Form fields", oat.AnchorLeft).WithRoundedCorner(true)

	// ── MultiLine EditText ────────────────────────────────────────────────────
	noteIn := widget.NewMultiLineEditText().
		WithStyle(latte.Style{Border: latte.BorderExplicitNone}).
		WithHint("Freeform notes").
		WithPlaceholder("Type anything here…")

	notePanel := layout.NewBorder(
		noteIn,
	).WithTitle("MultiLine EditText", oat.AnchorCenter).WithRoundedCorner(true)

	// ── Text with VAlign variants inside HBox ─────────────────────────────────
	// Each cell is a bordered box positioned differently in the HBox slot.
	// VAlign must be set on the direct HBox child — not on the Text inside the
	// Border — because HBox only applies the alignment offset to its own direct
	// children. We use AlignChild as the direct HBox child so the VAlign
	// is resolved at the correct layer.
	makeAlignBox := func(label, align string, va oat.VAlign) oat.Component {
		inner := widget.NewText(label).
			WithStyle(latte.Style{Bold: true})
		bordered := layout.NewBorder(inner).
			WithTitle(align, oat.AnchorCenter)
		return layout.NewAlignChild(bordered, oat.HAlignFill, va)
	}

	alignRow := layout.NewHBox(
		layout.NewFlexChild(makeAlignBox("top", "VAlignTop", oat.VAlignTop), 1),
		layout.NewFlexChild(makeAlignBox("mid", "VAlignMiddle", oat.VAlignMiddle), 1),
		layout.NewFlexChild(makeAlignBox("bot", "VAlignBottom", oat.VAlignBottom), 1),
	)

	// ── AlignChild showcase ───────────────────────────────────────────────────
	// Three buttons inside a VBox, each wrapped by AlignChild.
	acLeft := layout.NewAlignChild(
		widget.NewButton("AlignChild Left", func() {
			a.notifs.Push("AlignChild Left", widget.NotificationKindInfo, 2*time.Second)
		}),
		oat.HAlignLeft, oat.VAlignFill,
	)
	acCenter := layout.NewAlignChild(
		widget.NewButton("AlignChild Center", func() {
			a.notifs.Push("AlignChild Center", widget.NotificationKindInfo, 2*time.Second)
		}),
		oat.HAlignCenter, oat.VAlignFill,
	)
	acRight := layout.NewAlignChild(
		widget.NewButton("AlignChild Right", func() {
			a.notifs.Push("AlignChild Right", widget.NotificationKindInfo, 2*time.Second)
		}),
		oat.HAlignRight, oat.VAlignFill,
	)

	acSection := layout.NewBorder(
		layout.NewVBox(
			widget.NewText("AlignChild wrappers:"),
			acLeft,
			acCenter,
			acRight,
		),
	).WithTitle("AlignChild", oat.AnchorRight).WithRoundedCorner(true)

	// ── labels (tag chips) ────────────────────────────────────────────────────
	lblHighlight := widget.NewLabel([]string{"go", "tui", "oat-latte"})
	lblNoHighlight := widget.NewLabel([]string{"no-bg", "plain", "tags"}).
		WithHighlight(false)
	lblRight := widget.NewLabel([]string{"right", "aligned"}).
		WithHAlign(oat.HAlignRight)

	lblSection := layout.NewBorder(
		layout.NewPadding(
			layout.NewVBox(
				widget.NewText("with highlight:"),
				lblHighlight,
				widget.NewText("no highlight:"),
				lblNoHighlight,
				widget.NewText("right-aligned:"),
				lblRight,
			),
			oat.Insets{Top: 0, Right: 1, Bottom: 1, Left: 1},
		),
	).WithTitle("Labels (tag chips)", oat.AnchorLeft).WithRoundedCorner(true)

	col := layout.NewVBox(
		formPanel,
		layout.NewFlexChild(notePanel, 1),
		layout.NewFlexChild(
			layout.NewBorder(alignRow).
				WithTitle("VAlign in HBox", oat.AnchorCenter).
				WithRoundedCorner(true),
			1,
		),
		acSection,
		lblSection,
	)
	return col
}

// ── dialog ───────────────────────────────────────────────────────────────────

func (a *App) showShowcaseDialog() {
	// A dialog that itself shows Padding, AlignChild, Dividers, and styled Text.
	closeBtn := widget.NewButton("Close", func() { a.canvas.HideDialog() })

	titleText := widget.NewText("This dialog uses:").
		WithStyle(latte.Style{Bold: true})

	bullets := layout.NewVBox(
		widget.NewText("  • layout.NewPadding — left+right insets"),
		widget.NewText("  • widget.NewHDivider — full-width rule"),
		widget.NewText("  • layout.AlignChild  — centre + right buttons"),
		widget.NewText("  • latte.Style        — custom FG colours"),
		widget.NewText("  • Border.WithTitle   — AnchorCenter title"),
	)

	div := widget.NewHDivider()

	// Two info labels aligned differently
	infoLeft := widget.NewText("left info").
		WithStyle(latte.Style{FG: latte.ColorBrightCyan})
	infoRight := widget.NewText("right info").
		WithStyle(latte.Style{FG: latte.ColorYellow}).
		WithHAlign(oat.HAlignRight)

	// Close button centred via AlignChild
	closeCentered := layout.NewAlignChild(closeBtn, oat.HAlignCenter, oat.VAlignFill)

	body := layout.NewPadding(
		layout.NewVBox(
			titleText,
			layout.NewVFill().WithMaxSize(1),
			bullets,
			layout.NewVFill().WithMaxSize(1),
			div,
			layout.NewVFill().WithMaxSize(1),
			infoLeft,
			infoRight,
			layout.NewVFill().WithMaxSize(1),
			closeCentered,
		),
		oat.Insets{Top: 1, Right: 2, Bottom: 1, Left: 2},
	)

	a.canvas.ShowDialog(
		widget.NewDialog("Feature showcase dialog").
			WithChild(body).
			WithSize(widget.DialogPercent(55), widget.DialogPercent(65)),
	)
}

// ── build ─────────────────────────────────────────────────────────────────────

var themes = []latte.Theme{
	latte.ThemeDark,
	latte.ThemeLight,
	latte.ThemeDracula,
	latte.ThemeNord,
	latte.ThemeDefault,
}

var themeNames = []string{"dark", "light", "dracula", "nord", "default"}

var notifKinds = []widget.NotificationKind{
	widget.NotificationKindInfo,
	widget.NotificationKindSuccess,
	widget.NotificationKindWarning,
	widget.NotificationKindError,
}
var notifMessages = []string{
	"Info notification",
	"Success notification",
	"Warning notification",
	"Error notification",
}

func (a *App) build() {
	a.notifs = widget.NewNotificationManager()

	header := a.buildHeader()
	left := a.buildLeft()
	middle := a.buildMiddle()
	right := a.buildRight()

	body := layout.NewHBox()
	body.AddFlexChild(left, 3)
	body.AddFlexChild(middle, 3)
	body.AddFlexChild(right, 4)

	statusBar := widget.NewStatusBar()

	notifIdx := 0

	a.canvas = oat.NewCanvas(
		oat.WithTheme(themes[a.themeIdx]),
		oat.WithHeader(header),
		oat.WithBody(body),
		oat.WithAutoStatusBar(statusBar),
		oat.WithNotificationManager(a.notifs),
		oat.WithGlobalKeyBinding(
			oat.KeyBinding{
				Key:         tcell.KeyCtrlT,
				Mod:         tcell.ModCtrl,
				Label:       "^T",
				Description: "Cycle theme",
				Handler: func() {
					a.themeIdx = (a.themeIdx + 1) % len(themes)
					a.applyDerivedTheme()
				},
			},
			oat.KeyBinding{
				Key:         tcell.KeyCtrlD,
				Mod:         tcell.ModCtrl,
				Label:       "^D",
				Description: "Open dialog",
				Handler:     func() { a.showShowcaseDialog() },
			},
			oat.KeyBinding{
				Key:         tcell.KeyCtrlN,
				Mod:         tcell.ModCtrl,
				Label:       "^N",
				Description: "Notification",
				Handler: func() {
					k := notifKinds[notifIdx%len(notifKinds)]
					m := notifMessages[notifIdx%len(notifMessages)]
					notifIdx++
					a.notifs.Push(m, k, 3*time.Second)

					// Animate progress bars forward a bit
					a.progVal = math.Mod(a.progVal+0.1, 1.1)
					if a.progVal > 1.0 {
						a.progVal = 0
					}
					for _, pb := range a.progress {
						pb.SetValue(a.progVal)
					}
				},
			},
		),
	)

	a.notifs.Push("^T theme  ^D dialog  ^N notification  Tab focus", widget.NotificationKindInfo, 4*time.Second)
}

func main() {
	a := &App{roundedCorners: true}
	a.build()
	if err := a.canvas.Run(); err != nil {
		log.Fatal(err)
	}
}
