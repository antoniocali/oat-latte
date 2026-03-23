---
sidebar_position: 10
title: Dialog
description: Modal overlay widget that centres itself on screen with a scrim backdrop and confined focus.
---

# Dialog

`Dialog` is a modal overlay that renders a bordered, centred panel on top of the rest of the UI. While a dialog is visible:

- keyboard focus is confined to its subtree (Tab / arrow keys cannot reach the background)
- a full-screen scrim is painted behind the panel to dim the background

```go
import "github.com/antoniocali/oat-latte/widget"
```

## Constructor

```go
dlg := widget.NewDialog("Title")
```

No style argument is needed — colours and border are filled in by the active theme via `ApplyTheme`. Use `WithStyle` to override specific fields.

## Builder methods

| Method | Description |
|---|---|
| `WithChild(c oat.Component) *Dialog` | Sets the body component rendered inside the border |
| `WithID(id string) *Dialog` | Sets the widget ID for canvas lookup |
| `WithTitle(title string) *Dialog` | Overrides the title text |
| `WithStyle(s latte.Style) *Dialog` | Overrides visual style (theme acts as base; explicit fields take precedence) |
| `WithMaxSize(w, h int) *Dialog` | Fixed width × height in terminal cells |
| `WithSize(w, h DialogSize) *Dialog` | Flexible width × height using `DialogFixed` or `DialogPercent` |

## Sizing

Two size helpers are available:

```go
// Exact number of terminal cells.
widget.DialogFixed(60)

// Percentage of the available terminal dimension (1–100).
widget.DialogPercent(70)
```

When the terminal is resized, percent-based dialogs adapt automatically each render pass.

`WithMaxSize` is a shorthand for `WithSize(DialogFixed(w), DialogFixed(h))`. When both are called the last call wins. Default size is 60 × 20 cells.

## Showing and hiding

```go
app.ShowDialog(dlg)   // push modal overlay

// inside a button callback:
app.HideDialog()      // pop topmost overlay, restore body focus
```

## Typical usage

### Fixed-size confirmation

```go
func showConfirm(app *oat.Canvas, msg string, onYes func()) {
    cancelBtn := widget.NewButton("Cancel", func() { app.HideDialog() })
    okBtn     := widget.NewButton("OK",     func() { onYes(); app.HideDialog() })

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

### Percent-based form dialog

```go
func showForm(app *oat.Canvas) {
    titleIn := widget.NewEditText().
        WithStyle(latte.Style{Border: latte.BorderExplicitNone}).
        WithHint("Title").
        WithPlaceholder("Enter title…")

    bodyIn := widget.NewMultiLineEditText().
        WithStyle(latte.Style{Border: latte.BorderExplicitNone}).
        WithHint("Body").
        WithPlaceholder("Write here…")

    cancelBtn := widget.NewButton("Cancel", func() { app.HideDialog() })
    saveBtn   := widget.NewButton("Save",   func() { /* save logic */ app.HideDialog() })

    btnRow := layout.NewHBox()
    btnRow.AddChild(layout.NewHFill())
    btnRow.AddChild(cancelBtn)
    btnRow.AddChild(layout.NewHFill().WithMaxSize(2))
    btnRow.AddChild(saveBtn)

    form := layout.NewVBox(titleIn)
    form.AddFlexChild(bodyIn, 1)
    form.AddChild(layout.NewVFill().WithMaxSize(1))
    form.AddChild(btnRow)

    app.ShowDialog(
        widget.NewDialog("New Item").
            WithChild(layout.NewPaddingUniform(form, 1)).
            WithSize(widget.DialogPercent(60), widget.DialogPercent(70)),
    )
}
```

## Theme tokens

| Token | Applied to |
|---|---|
| `Dialog` | Border and background of the dialog panel |
| `DialogTitle` | Title text rendered into the top border rule |
| `Scrim` | Full-screen backdrop painted behind the dialog |

`ApplyTheme` uses the **theme-as-base** merge pattern — fields you set explicitly via `WithStyle` are preserved:

```go
func (d *Dialog) ApplyTheme(t latte.Theme) {
    d.Style = t.Dialog.Merge(d.Style)   // theme is base; caller-set fields win
    d.titleStyle = t.DialogTitle
    d.scrimStyle = t.Scrim
}
```
