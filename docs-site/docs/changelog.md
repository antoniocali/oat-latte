---
sidebar_position: 99
---

# Changelog

All notable changes to the oat-latte framework are listed here, newest first.

---

## v0.2.8

**Cross-axis alignment · `RoundedCorner` theme flag · `callerStyle` pattern**

### Added

- **`oat.HAlign`** — horizontal-axis alignment type for widgets inside a `VBox`:
  - `HAlignFill` (0, default) — fill the full allocated width; **zero-value, all existing code unchanged**.
  - `HAlignLeft` — shrink to natural width, pin to the left.
  - `HAlignCenter` — shrink to natural width, centre horizontally.
  - `HAlignRight` — shrink to natural width, pin to the right.

- **`oat.VAlign`** — vertical-axis alignment type for widgets inside an `HBox`:
  - `VAlignFill` (0, default) — fill the full allocated height; **zero-value, all existing code unchanged**.
  - `VAlignTop` — shrink to natural height, pin to the top.
  - `VAlignMiddle` — shrink to natural height, centre vertically.
  - `VAlignBottom` — shrink to natural height, pin to the bottom.

- **`oat.AlignProvider`** interface — satisfied automatically by any component that embeds `BaseComponent`:
  ```go
  type AlignProvider interface {
      GetHAlign() HAlign
      GetVAlign() VAlign
  }
  ```

- **`BaseComponent.HAlign` / `BaseComponent.VAlign`** fields — carry the alignment preference for a widget. Zero values mean "use parent box default" (fill/stretch).

- **`BaseComponent.GetHAlign()` / `GetVAlign()`** — implement `AlignProvider`.

- **Concrete `WithHAlign` / `WithVAlign` fluent builders on every built-in widget** — `Text`, `Title`, `Button`, `CheckBox`, `EditText`, `List`, `ComponentList`, `Label`, `ProgressBar`, `StatusBar`, `Divider`. Each returns the concrete widget type so alignment can be set inline without a type assertion:

  ```go
  saveBtn   := widget.NewButton("Save",   fn).WithHAlign(oat.HAlignRight)
  cancelBtn := widget.NewButton("Cancel", fn).WithHAlign(oat.HAlignLeft)
  lbl       := widget.NewText("note").WithVAlign(oat.VAlignBottom)
  ```

  Each method is variadic (`a ...oat.HAlign`) — calling with no argument resets to `HAlignFill` / `VAlignFill`. Internally they delegate to `BaseComponent.HAlign` / `BaseComponent.VAlign`; custom widgets that embed `BaseComponent` but have not yet added their own builders can set the field directly: `myWidget.BaseComponent.HAlign = oat.HAlignRight`. `NotificationManager` and `Dialog` are intentionally excluded — they are never placed directly in a `VBox`/`HBox`.

- **`VBox.WithHAlign(a ...oat.HAlign) *VBox`** — sets a box-wide default horizontal alignment for all children that do not declare their own.

- **`HBox.WithVAlign(a ...oat.VAlign) *HBox`** — sets a box-wide default vertical alignment for all children that do not declare their own.

- **`layout.AlignChild`** — a wrapper component that attaches explicit alignment overrides to any child. Takes precedence over both the child's own `BaseComponent` alignment and the box-wide default.

  ```go
  // Constructor:
  func NewAlignChild(child oat.Component, h oat.HAlign, v oat.VAlign) *AlignChild

  // Usage:
  vbox := layout.NewVBox(
      layout.NewAlignChild(saveBtn,   oat.HAlignRight, oat.VAlignFill),
      layout.NewAlignChild(cancelBtn, oat.HAlignLeft,  oat.VAlignFill),
  )
  ```

  `AlignChild` implements `oat.Layout` (via `Children()`), so theme propagation and focus collection recurse into the wrapped component automatically. It is **not** a `FlexSpacer` — wrap it with `AddFlexChild` or `NewFlexChild` when flex behaviour is also needed.

- **`latte.Theme.RoundedCorner bool`** — new field on `Theme`. All five built-in themes set `RoundedCorner: true`. When `ApplyTheme` is called, `Button` and `Border` widgets automatically adopt arc corners (`╭─╮/╰─╯`) without any per-widget configuration.

- **`(Theme).WithRoundedCorner(bool) Theme`** — builder method to set `RoundedCorner` on a derived theme. Returns a new `Theme` by value; the originals are never mutated.

### Changed (non-breaking)

- `VBox.Render` — resolves each child's effective `HAlign` before assigning its region. When `HAlignFill`, behaviour is identical to before (unchanged code path). For non-fill values the child is re-measured and positioned within its row.
- `HBox.Render` — same for `VAlign`.

- **`Border.WithRoundedCorner`** — completely reworked. Previously mutated `Style.Border` at call time and panicked on incompatible styles. Now:
  - Stores the intent in an internal `roundedCorner bool` field; **does not mutate `Style.Border`**.
  - The effective corner shape is resolved at **render time**: `BorderSingle ↔ BorderRounded` is toggled based on the field; incompatible styles (`BorderDouble`, `BorderThick`, `BorderDashed`) are **silently left unchanged** — **no panic**.
  - Once called, the explicit preference overrides `theme.RoundedCorner` on all future `ApplyTheme` calls. `WithRoundedCorner(false)` opts out even when `theme.RoundedCorner` is `true`.
  - `ApplyTheme` syncs `roundedCorner` from `t.RoundedCorner` **only** when `WithRoundedCorner` has never been called on this instance (bidirectional — both enabling and disabling work correctly on theme switch).

- **`Button.WithRoundedCorner`** — same rework as `Border`:
  - Stores intent in `roundedCorner bool` / `roundedCornerSet bool` fields; does not mutate `Style.Border`.
  - Incompatible styles are **silently ignored** — **no panic**.
  - `ApplyTheme` syncs from `t.RoundedCorner` when `!roundedCornerSet`.
  - `fmt` import removed from `widget/button.go` (was only used for the now-deleted panic message).

- **`callerStyle` pattern** — all built-in widgets that expose `WithStyle` now store the caller's original intent in a `callerStyle` field and use **that** as the `Merge` base in `ApplyTheme`. This prevents stale colours from a previous theme accumulating in `w.Style` and blocking a new theme from fully replacing them on `SetTheme` calls after the first.

### Resolution order (per child in each box)

1. `AlignChild` wrapper — highest priority.
2. `AlignProvider` on the child itself (e.g. `BaseComponent.HAlign`/`VAlign` set to a non-fill value).
3. Box-wide default (`WithHAlign` / `WithVAlign`).
4. `HAlignFill` / `VAlignFill` — full-stretch fallback; identical to the pre-v0.2.8 behaviour.

---

## v0.2.7

_(v0.2.6 was skipped — these changes were shipped directly as v0.2.7.)_

**`ComponentList` widget · `ThemeLight` colour fixes · theme-aware accent colours · `Buffer` background inheritance**

### Fixed

- **`ThemeLight` colour palette** — several tokens produced illegible or visually incorrect output:
  - `ListSelected` had no background colour → added `LightBgSelected = RGB(219, 228, 253)` (light cobalt tint).
  - `Header` / `Footer` used the canvas background (`LightBg`) → added `LightBgHeader = RGB(232, 229, 224)` (slightly darker warm gray) so header and footer regions are visually distinct.
  - `LightBgScrim` was too close to the canvas colour, making dialogs nearly invisible → darkened from `RGB(224, 220, 215)` to `RGB(180, 175, 168)`.
  - `LightMuted` failed WCAG AA contrast against the light background → darkened from `#8e8e9a` to `RGB(110, 110, 125)` (`#6e6e7d`).
- **All four true-color themes (`ThemeDark`, `ThemeLight`, `ThemeDracula`, `ThemeNord`)** — `Panel.BG` was `ColorDefault` (unfilled), causing panel backgrounds to appear as terminal-black instead of the theme canvas colour. Each theme now has an explicit `Panel.BG` matching its canvas background.
- **`StatusBar` — hardcoded `ColorBrightCyan` for key-hint brackets** → `ApplyTheme` now reads `t.Accent.FG` and stores it as `accentColor`; bracket labels are rendered in the active theme's accent colour. Falls back to `ColorBrightCyan` for `ThemeDefault` (ANSI-16).
- **`ComponentList` — hardcoded `ColorBrightCyan` for the `>` cursor glyph** → same fix; cursor colour is now driven by `t.Accent.FG` via `accentColor`.
- **`Buffer` background inheritance** — the root cause of all "terminal-black bleed-through" bugs on light (and custom) themes. tcell has no transparency: every `SetContent` call with `tcell.ColorDefault` paints terminal black over any background the parent already drew. Fixed by adding a `bg latte.Color` field to `Buffer`:
  - `Fill` / `FillBG` record `style.BG` into `b.bg` when it is a concrete colour (the canvas pre-fill establishes the theme background for the whole tree).
  - `Sub` propagates `bg` to child buffers, so every widget in the tree inherits the canvas background colour automatically.
  - New `resolveStyle`: substitutes `b.bg` for `ColorDefault` in `style.BG` before passing to tcell.
  - New `resolveBorderStyle`: same substitution for `style.BorderBG`.
  - `SetCell`, `DrawText`, `DrawTextAligned` all call `resolveStyle`.
  - `DrawBorderTitle` calls `resolveBorderStyle` for border runes and `resolveStyle` for title text.
  - Custom widgets that draw with `BG == ColorDefault` automatically inherit the canvas background — no code changes required. Only override `BG` when a widget intentionally wants a specific background colour.

### Changed

- `latte/colors.go` — two new `Light*` palette constants (`LightBgSelected`, `LightBgHeader`) and updated values for `LightBgScrim` and `LightMuted`.

---

## v0.2.6

**`widget.ComponentList` — component-row list**

### Added

- `widget.ComponentList` — a vertically scrollable, keyboard-navigable list whose rows are arbitrary `Component` values instead of plain label strings. This allows each row to contain rich layouts such as `HBox(Text, Flex(Text), Text)`.

  - `widget.ComponentListItem{Component oat.Component, Value interface{}}` — the item type. `Component` is the widget rendered for the row; `Value` is an opaque identifier the caller can use to correlate a row with application data (e.g. a record ID).
  - `widget.NewComponentList(items []ComponentListItem) *ComponentList` — constructor.
  - Builder options mirror `List` exactly: `WithStyle`, `WithID`, `WithSelectedStyle`, `WithHighlight`, `WithCursor`, `WithOnSelect`, `WithOnDelete`, `WithOnCursorChange`.
  - `SetItems`, `SelectedItem() (ComponentListItem, bool)`, `SelectedIndex() int`.
  - Row heights are **variable**: each row's component is measured via `Measure` (unconstrained on Y) to determine how many terminal lines it occupies. Scroll is tracked by item index so the viewport always shows complete rows.
  - Implements `oat.Layout` (`Children()` / `AddChild`) so theme propagation and the focus collector recurse into every row component automatically.
  - Implements `oat.ValueGetter`: `Canvas.GetValue(id)` returns the `Value` field of the currently selected item.
  - Theme tokens: `t.Text` (base), `t.ListSelected` (highlight), `t.FocusBorder` (focused border colour) — identical to `List`.

```go
makeRow := func(name, role, status string, id int) widget.ComponentListItem {
    row := layout.NewHBox(
        widget.NewText(name),
        layout.NewFlexChild(widget.NewText(role), 1),
        widget.NewText(status),
    )
    return widget.ComponentListItem{Component: row, Value: id}
}

list := widget.NewComponentList([]widget.ComponentListItem{
    makeRow("Alice",   "Backend engineer",  "active",   1),
    makeRow("Bob",     "Frontend engineer", "inactive", 2),
    makeRow("Charlie", "DevOps",            "active",   3),
}).
    WithID("people-list").
    WithOnSelect(func(idx int, item widget.ComponentListItem) {
        id := item.Value.(int)
        // open record
    })
```

---

## v0.2.5

**Named true-color palette · `widget.Divider` · `oat.VAnchor`**

### Added

- `latte/colors.go` — ~120 named true-color `Color` constants grouped into utility scales and design-system palettes. All values are `latte.Color` variables (via `latte.RGB`), compatible with any API that accepts a `latte.Color`.

  **Utility scales** (Tailwind-style, shade numbers 50–950): `Slate`, `Zinc`, `Stone`, `Sky`, `Blue`, `Indigo`, `Cornflower`, `Cyan`, `Teal`, `Emerald`, `Green`, `Lime`, `Yellow`, `Amber`, `Orange`, `Red`, `Rose`, `Pink`, `Violet`, `Purple`, `Fuchsia`.

  **Design-system palettes** — named constants extracted from every hex literal that appears in the four true-color built-in themes:
  - `Dark*` — `DarkBg`, `DarkBgElevated`, `DarkBgScrim`, `DarkBorder`, `DarkMuted`, `DarkAccent`, `DarkSuccess`, `DarkWarning`, `DarkError`
  - `Light*` — `LightBg`, `LightBgElevated`, `LightBgScrim`, `LightBorder`, `LightMuted`, `LightText`, `LightAccent`, `LightSuccess`, `LightWarning`, `LightError`
  - `Dracula*` — `DraculaBg`, `DraculaBgElevated`, `DraculaBgScrim`, `DraculaFg`, `DraculaComment`, `DraculaSelection`, `DraculaPurple`, `DraculaCyan`, `DraculaGreen`, `DraculaOrange`, `DraculaRed`, `DraculaYellow`, `DraculaPink`
  - `Nord0`–`Nord15` + `NordBg`, `NordBgElevated`, `NordBgScrim`

- `widget.Divider` — axis-agnostic rule widget for placing horizontal (`─`) or vertical (`│`) separators between layout children.
  - `widget.NewHDivider()` / `widget.NewVDivider()` — convenience constructors.
  - `widget.NewDivider(axis widget.Axis)` — explicit-axis constructor (`widget.AxisHorizontal` / `widget.AxisVertical`).
  - `DividerSize` — controls how much of the allocated space the visible rule occupies: `widget.DividerFill` (default, full span), `widget.DividerFixed(n)` (exactly `n` cells), `widget.DividerPercent(p)` (1–100% of the allocated length).
  - `(*Divider).WithRune(r rune)` — override the line character (e.g. `'═'` for a double rule).
  - `(*Divider).WithMaxSize(size DividerSize, anchor ...oat.Anchor)` — for `AxisHorizontal` dividers: limits width and positions the rule horizontally.
  - `(*Divider).WithMaxSizeV(size DividerSize, anchor ...oat.VAnchor)` — for `AxisVertical` dividers: limits height and positions the rule vertically.
  - `(*Divider).WithStyle(s latte.Style)` — override the display style.
  - `ApplyTheme` maps the `Muted` theme token onto the divider; override with `WithStyle`.

- `oat.VAnchor` — new vertical-axis positioning type (`VAnchorTop`, `VAnchorMiddle`, `VAnchorBottom`), the V-axis counterpart to `oat.Anchor`. The two types are kept separate so the compiler enforces correct axis usage — APIs that accept H-axis placement cannot accidentally receive a `VAnchor` and vice versa.

### Changed

- `ThemeDark`, `ThemeLight`, `ThemeDracula`, `ThemeNord` — all raw `Hex("...")` literals replaced with the new named constants. No visual change; pure readability / maintainability improvement.

```go
// Named constants instead of raw hex strings.
myTheme := latte.ThemeDark.
    WithAccent(latte.Style{FG: latte.Pink500}).
    WithFocusBorder(latte.Pink500).
    WithName("dark-pink")

// Horizontal rule — place in a VBox between items
hd := widget.NewHDivider()
hd := widget.NewHDivider().WithRune('═')
hd := widget.NewHDivider().WithMaxSize(widget.DividerPercent(60), oat.AnchorCenter)

// Vertical rule — place in an HBox between items
vd := widget.NewVDivider()
vd := widget.NewVDivider().WithMaxSizeV(widget.DividerFixed(8), oat.VAnchorMiddle)
```

---

## v0.2.4

**`Theme` fluent builder methods**

### Added

- `(Theme).WithName(string) Theme` — returns a copy of the theme with a new name, useful when naming a derived theme.
- `(Theme).WithFocusBorder(Color) Theme` — replaces the `FocusBorder` colour (a plain `Color`, not a `Style`).
- One `With<Token>(latte.Style) Theme` method for every `Style`-typed field on `Theme`: `WithCanvas`, `WithText`, `WithMuted`, `WithAccent`, `WithSuccess`, `WithWarning`, `WithError`, `WithPanel`, `WithPanelTitle`, `WithInput`, `WithInputFocus`, `WithListSelected`, `WithButton`, `WithButtonFocus`, `WithCheckBox`, `WithCheckBoxFocus`, `WithHeader`, `WithFooter`, `WithDialog`, `WithDialogTitle`, `WithScrim`, `WithTag`, `WithNotificationInfo`, `WithNotificationSuccess`, `WithNotificationWarning`, `WithNotificationError`.

All methods return `Theme` by value — built-in theme variables (`ThemeDark`, `ThemeNord`, etc.) are never mutated. Style-typed methods use `Style.Merge` internally so only the non-zero fields of the supplied `Style` are applied; the rest of the token is preserved.

```go
// Nord but with no borders anywhere
borderless := latte.ThemeNord.
    WithPanel(latte.Style{Border: latte.BorderExplicitNone}).
    WithInput(latte.Style{Border: latte.BorderExplicitNone}).
    WithButton(latte.Style{Border: latte.BorderExplicitNone}).
    WithDialog(latte.Style{Border: latte.BorderExplicitNone}).
    WithName("nord-borderless")

// Dark theme with a custom accent and focus colour
pink := latte.ThemeDark.
    WithAccent(latte.Style{FG: latte.Hex("#ff69b4")}).
    WithFocusBorder(latte.Hex("#ff69b4")).
    WithName("dark-pink")

app.SetTheme(borderless)
```

---

## v0.2.3

**`layout.FlexChild` · `Button.WithRoundedCorner` · `Label.WithHighlight` · `Dialog` percent-height fix · `Button` render fix · Label `FillBG` fix**

_(v0.2.2 was skipped — these changes were shipped directly as v0.2.3.)_

### Added

- `layout.NewFlexChild(child oat.Component, weight ...int) *FlexChild` — wraps any `Component` as a flex slot so it can be passed to the variadic `NewVBox` / `NewHBox` constructors without a separate `AddFlexChild` call. Weight defaults to `1`; minimum effective weight is `1`. `FlexChild` implements `oat.Layout` via `Children()` so theme propagation and focus collection recurse into the wrapped component automatically.
- `(*Button).WithRoundedCorner(bool)` — draws arc corners (`╭╮╰╯`) on the button border when `true`. Incompatible border styles (`BorderDouble`, `BorderThick`, `BorderDashed`) are silently ignored. *(Panic behaviour was removed in v0.2.8.)*
- `(*Button).WithStyle(latte.Style)` — new builder; accepts any border style.
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

- `(*Border).WithRoundedCorner(bool)` — switches the effective border corner style to arc (`╭─╮│╰─╯`) when `true`. Stores intent in an internal field; does not mutate `Style.Border`. Incompatible styles (`BorderDouble`, `BorderThick`, `BorderDashed`) are silently ignored. *(Original panic behaviour was removed in v0.2.8.)*

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
- `Canvas` with `WithTheme`, `WithHeader`, `WithBody`, `WithAutoStatusBar`, `WithPrimary`, `WithGlobalKeyBinding`, `SetTheme`, `GetTheme`, `ShowDialog`, `HideDialog`, `FocusByRef`, `InvalidateLayout`.
- Cooperative focus system with Tab / Shift+Tab cycling, `FocusGuard` for context-sensitive subtrees, and programmatic focus via `FocusByRef`.
- Three example apps: `tasklist`, `notes`, `kanban`.
