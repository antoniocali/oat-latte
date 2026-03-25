---
sidebar_position: 99
---

# Changelog

All notable changes to oat-latte are listed here, newest first.

---

## v0.2.0

**`WithNotificationManager` canvas option · delete shortcut · tutorial expansion**

### Added

- `oat.WithNotificationManager(nm)` — new canvas option that wires the notification re-render channel **and** mounts the manager as a persistent overlay in a single call, replacing the previous two-step `SetNotifyChannel` + `ShowPersistentOverlay` pattern.
- `oat.NotificationOverlay` interface — introduced to avoid a circular import between the `oat` and `widget` packages.
- `'d'` key as a delete shortcut in all three example apps (`tasklist`, `notes`, `kanban`), forwarding a synthetic `tcell.KeyDelete` event to `List.HandleKey`.
- `make run-tasklist` target in the `Makefile`.
- Tutorial Step 7: documents the synthetic-event technique for custom key aliases.
- `CONTRIBUTING.md` with setup instructions, code conventions, and PR guidelines.

### Changed

- All example apps updated to the new `WithNotificationManager` API.
- `README.md` updated with the `run-tasklist` make target and a link to `CONTRIBUTING.md`.
- `docs/installation.md` updated to list the tasklist example app.

### Removed

- `app.NotifyChannel()` public method (breaking change — use `WithNotificationManager` instead).

---

## v0.1.1

**`Border.WithRoundedCorner`**

### Added

- `(*Border).WithRoundedCorner(bool)` — switches the border style to `BorderRounded` (`╭─╮│╰─╯`) when `true`, restores `BorderSingle` when `false`. Panics with an explicit message if called with `true` on `BorderDouble`, `BorderThick`, or `BorderDashed` (Unicode provides arc corner codepoints only for light-weight strokes).

### Changed

- `AGENTS.md`, `docs/layout.md`, and `docs/agents.md` updated with the `WithRoundedCorner` API, a compatibility table, and the panic contract.
- Navbar version badge border styling removed; the version link now renders as a plain navbar link.

---

## v0.1.0

**Initial release**

### Added

- Two-pass layout engine (`Measure` / `Render`) built on [tcell](https://github.com/gdamore/tcell).
- Core interfaces: `Component`, `Layout`, `Focusable`, `FocusGuard`.
- `BaseComponent` and `FocusBehavior` embeds for custom widgets.
- Layout containers: `VBox`, `HBox`, `Grid`, `Stack`, `Border`, `Padding`, `VFill`, `HFill`.
- Widget library: `Text`, `Title`, `Button`, `CheckBox`, `EditText` (single- and multi-line), `List`, `Label`, `ProgressBar`, `StatusBar`, `NotificationManager`, `Dialog`.
- Style system with `Style.Merge`, border sentinels (`BorderNone`, `BorderExplicitNone`, `BorderSingle`, `BorderRounded`, `BorderDouble`, `BorderThick`, `BorderDashed`), and true-color support (`latte.RGB`, `latte.Hex`).
- Five built-in themes: `ThemeDefault`, `ThemeDark`, `ThemeLight`, `ThemeDracula`, `ThemeNord`.
- `Canvas` with `WithTheme`, `WithHeader`, `WithBody`, `WithAutoStatusBar`, `WithPrimary`, `WithGlobalKeyBinding`, `SetTheme`, `ShowDialog`, `HideDialog`, `FocusByRef`, `InvalidateLayout`.
- Cooperative focus system with Tab / Shift+Tab cycling, `FocusGuard` for context-sensitive subtrees, and programmatic focus via `FocusByRef`.
- Three example apps: `tasklist`, `notes`, `kanban`.
- Full documentation site at [oat-latte.dev](https://oat-latte.dev).
