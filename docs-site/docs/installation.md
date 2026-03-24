---
sidebar_position: 2
title: Installation
description: Add oat-latte to your Go module.
---

# Installation

## Requirements

- Go 1.21 or later
- A terminal that supports at least 256 colors (true-color recommended for the built-in themes)

## Add the module

```bash
go get github.com/antoniocali/oat-latte
```

## Packages

| Import path | What it contains |
|---|---|
| `github.com/antoniocali/oat-latte` | `Canvas`, `Buffer`, `FocusManager`, core interfaces, geometry types |
| `github.com/antoniocali/oat-latte/latte` | `Style`, `Color`, `BorderStyle`, `Theme`, built-in themes |
| `github.com/antoniocali/oat-latte/layout` | `VBox`, `HBox`, `Grid`, `Border`, `Padding`, `VFill`, `HFill`, `Dialog` |
| `github.com/antoniocali/oat-latte/widget` | All built-in widgets |

## Minimal working program

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
    msg  := widget.NewText("Hello from oat-latte!")
    body := layout.NewVBox(msg)

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

```bash
go run .
```

Press `Esc` or `Ctrl+C` to quit.

## Example apps

The repository ships two full example applications:

```bash
go run github.com/antoniocali/oat-latte/cmd/example/notes
go run github.com/antoniocali/oat-latte/cmd/example/kanban
```

Or, if you have cloned the repo:

```bash
make run-notes
make run-kanban
```
