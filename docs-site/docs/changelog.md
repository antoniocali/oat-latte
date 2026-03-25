---
sidebar_position: 99
---

# Changelog

All notable changes to the oat-latte framework are listed here, newest first.

---

## v0.2.1

**`oat.Anchor` · `ProgressBar.WithPercentage` · `Border.WithTitle` anchor**

### Added

- `oat.Anchor` iota type (`AnchorLeft`, `AnchorCenter`, `AnchorRight`) — general-purpose horizontal-position enum for layout and widget APIs.
- `(*ProgressBar).WithPercentage(show bool, anchor ...oat.Anchor)` — controls whether the `XX%` label is rendered and where it appears (left, center, or right of the bar fill). `WithShowPercent(bool)` is kept as a deprecated backward-compatible alias.
- `(*Border).WithTitle(title string, anchor ...oat.Anchor)` — variadic anchor parameter to position the title in the top border rule. Omitting the anchor preserves left-alignment; all existing call sites compile and behave identically.
- `Buffer.DrawBorderTitle` now accepts a final `anchor Anchor` parameter. Internal callers that do not expose anchor control pass `oat.AnchorLeft` explicitly.

---

## v0.2.0

**`WithNotificationManager` canvas option · `oat.NotificationOverlay` interface**

### Added

- `oat.WithNotificationManager(nm)` — canvas option that wires the notification re-render channel and mounts the manager as a persistent overlay in a single call, replacing the previous two-step `SetNotifyChannel` + `ShowPersistentOverlay` pattern.
- `oat.NotificationOverlay` interface — introduced to decouple the `oat` and `widget` packages and avoid a circular import.

### Removed

- `app.NotifyChannel()` public method — use `oat.WithNotificationManager` instead.

---

## v0.1.1

**`Border.WithRoundedCorner`**

### Added

- `(*Border).WithRoundedCorner(bool)` — switches the border style to `BorderRounded` (`╭─╮│╰─╯`) when `true`, restores `BorderSingle` when `false`. Panics with an explicit message if called with `true` on `BorderDouble`, `BorderThick`, or `BorderDashed` (Unicode provides arc corner codepoints only for light-weight strokes).

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
