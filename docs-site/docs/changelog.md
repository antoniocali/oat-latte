---
sidebar_position: 99
---

# Changelog

All notable changes to the oat-latte framework are listed here, newest first.

---

## v0.2.3

**`layout.FlexChild` · `Button.WithRoundedCorner` · `Label.WithHighlight` · `Dialog` percent-height fix · `Button` render fix · Label `FillBG` fix**

_(v0.2.2 was skipped — these changes were shipped directly as v0.2.3.)_

### Added

- `layout.NewFlexChild(child oat.Component, weight ...int) *FlexChild` — wraps any `Component` as a flex slot so it can be passed to the variadic `NewVBox` / `NewHBox` constructors without a separate `AddFlexChild` call. Weight defaults to `1`; minimum effective weight is `1`. `FlexChild` implements `oat.Layout` via `Children()` so theme propagation and focus collection recurse into the wrapped component automatically.
- `(*Button).WithRoundedCorner(bool)` — draws arc corners (`╭╮╰╯`) on the button border when `true`. Panics at render time if the effective border style is `BorderDouble`, `BorderThick`, or `BorderDashed`. `WithStyle` also panics immediately at construction time for those styles.
- `(*Button).WithStyle(latte.Style)` — new builder; validates that the border style is compatible (`BorderNone`, `BorderExplicitNone`, `BorderSingle`, or `BorderRounded`) and panics otherwise.
- `(*Label).WithHighlight(bool)` — controls whether chip badges render with their background colour fill. `false` keeps the foreground colour and text attributes but strips the background, useful for minimal or transparent UIs. Default is `true` (existing behaviour).

### Changed

- All built-in themes: `Button` token now carries `Border: BorderSingle` so buttons always render with a visible border regardless of focus state. `ButtonFocus` now carries only `Reverse: true` and an accent `BorderFG` — the border shape is no longer placed in `ButtonFocus`. This makes `Button.Measure` return a stable `Height: 3` (border + label + border) at all times, fixing a layout instability where the button would grow from 1 to 3 rows upon receiving focus.
- `Button.Measure` and `Button.Render` now derive border presence from `b.Style` (the unfocused base style) rather than `EffectiveStyle(IsFocused())`. `EffectiveStyle` is still used for colour and attribute rendering. This is the correct separation: shape is stable, colour changes on focus.
- `cmd/example/tasklist`: `showNewDialog` updated to `WithMaxSize(50, 13)` (was 11) to accommodate the stable `Height: 3` that buttons now always report (border + label + border) when `Border: BorderSingle` is set by the theme.

### Fixed

- `Dialog.Measure` and `Dialog.Render` no longer shrink the dialog panel to its content height. Both now derive the panel dimensions solely from `maxDimensions(region)`, so percent-based dialogs (`DialogPercent`) correctly fill the requested fraction of the terminal each frame and adapt when the terminal is resized.
- `Button.Render` inner label is now drawn into `sub.Sub({X:1, Y:1, ...})` (a clipped sub-buffer of the button's allocated region) rather than `buf.Sub({X: region.X+1, ...})`. The old code bypassed the clipping fence of the enclosing `Border` and could paint the label one character outside the button's allocated area.
- `Label.Render` was calling `sub.FillBG(latte.Style{BG: l.Style.BG})`, which leaked the tag background colour across the entire row width — bleeding into separators and trailing space beyond the last chip. Changed to `sub.FillBG(latte.Style{})` so the parent background shows through outside chip boundaries.

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
