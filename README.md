# oat-latte

![oat-latte](static/logo.png)

A component-based TUI framework for Go.

oat-latte provides a two-pass layout engine (Measure/Render), a cooperative focus system, a composable style and theme system, and a library of ready-made widgets — all built on [tcell](https://github.com/gdamore/tcell).

## Install

```sh
go get github.com/antoniocali/oat-latte
```

Requires Go 1.21+ and a true-color terminal for the built-in themes.

## Quick start

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
    input := widget.NewEditText(latte.Style{}).WithPlaceholder("Type something…")
    btn   := widget.NewButton("OK", latte.Style{}, func() { /* handle */ })

    body := layout.NewVBox(input, btn)

    app := oat.NewCanvas(
        oat.WithTheme(latte.ThemeDark),
        oat.WithBody(body),
        oat.WithPrimary(input),
    )

    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

## Examples

```sh
make run-notes    # Notes manager
make run-kanban   # Kanban board
```

## Documentation

Full documentation is at **[antoniocali.github.io/oat-latte](https://antoniocali.github.io/oat-latte/)**.

## License

MIT
