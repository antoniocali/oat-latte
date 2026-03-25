# Contributing to oat-latte

Thank you for your interest in contributing. This document covers how to set up the project, the conventions used, and how to submit changes.

## Prerequisites

- Go 1.21 or later
- A true-color terminal (for running examples)
- `make` (optional, for convenience targets)

## Getting started

```sh
git clone https://github.com/antoniocali/oat-latte.git
cd oat-latte
go mod download
```

Verify everything builds and passes vet:

```sh
go build ./...
go vet ./...
```

Run one of the example apps to check the UI works:

```sh
make run-tasklist
```

## Project structure

| Path | Contents |
|---|---|
| `*.go` (root) | Core interfaces: `Component`, `Layout`, `Focusable`, `Canvas`, `Buffer`, geometry types |
| `latte/` | `Style`, `Color`, `BorderStyle`, `Theme`, built-in themes |
| `layout/` | Layout containers: `VBox`, `HBox`, `Grid`, `Border`, `Padding`, fill spacers |
| `widget/` | Ready-made widgets: `Text`, `Button`, `EditText`, `List`, `ProgressBar`, etc. |
| `cmd/example/` | Example applications (`tasklist`, `notes`, `kanban`) |
| `docs-site/` | Docusaurus documentation site |

## Making changes

### Code conventions

- Follow standard Go formatting (`gofmt`).
- Every new `Component` must implement the two-pass contract: **Measure** before **Render**, never storing `Buffer` or `Region` between frames.
- Embed `oat.BaseComponent` in custom widgets. Call `e.EnsureID()` in constructors.
- Use `Style.Merge` in `ApplyTheme` implementations — never direct struct assignment. See `AGENTS.md` for the full rationale.
- New interactive widgets must implement `Focusable` and embed `oat.FocusBehavior`.

### Before opening a pull request

```sh
go build ./...
go vet ./...
```

Both must exit cleanly (no output, status 0). There is no automated test suite yet; manual verification with the example apps is expected.

### Commit messages

Use the imperative mood and a short subject line (≤ 72 chars). Prefix with a type when appropriate:

```
feat: add RadioGroup widget
fix: correct border rendering when width < 4
docs: update focus system section in AGENTS.md
ci: cache npm dependencies in deploy workflow
```

### Pull requests

- Target the `main` branch.
- Keep PRs focused — one logical change per PR.
- Describe *what* changed and *why* in the PR body.
- Reference any related issues with `Closes #N`.

## Documentation

The docs site lives in `docs-site/` and is built with [Docusaurus 3](https://docusaurus.io). To preview locally:

```sh
cd docs-site
npm install
npm start
```

The `AGENTS.md` file at the repo root is the authoritative reference for AI coding agents and doubles as detailed internal documentation. Keep it in sync with any API changes.

## License

By contributing you agree that your contributions will be licensed under the [MIT License](LICENSE).
