---
sidebar_position: 8
title: NotificationManager
description: Toast notification banners with auto-dismiss.
---

# NotificationManager

`widget.NotificationManager` displays transient toast banners anchored to the bottom-right corner of the screen. Banners can auto-dismiss after a duration or stay until dismissed manually.

`NotificationManager` has no style constructor argument and no `WithStyle` builder — its appearance is entirely controlled by the four `Notification*` theme tokens.

## Setup

Pass the `NotificationManager` to `oat.WithNotificationManager` when constructing the canvas. This wires the timer channel and mounts the manager as a persistent overlay automatically.

```go
notifs := widget.NewNotificationManager()

app := oat.NewCanvas(
    oat.WithTheme(latte.ThemeDark),
    oat.WithBody(body),
    oat.WithNotificationManager(notifs),  // wires channel + mounts as persistent overlay
)
```

## Pushing notifications

```go
// Auto-dismiss after 2 seconds.
notifs.Push("Note saved", widget.NotificationKindSuccess, 2*time.Second)

// Auto-dismiss after 4 seconds.
notifs.Push("Welcome!", widget.NotificationKindInfo, 4*time.Second)

// Stays until manually cleared (duration 0 = no auto-dismiss).
notifs.Push("Build failed", widget.NotificationKindError, 0)
```

## Dismissing manually

```go
notifs.Pop()    // remove the most recently added notification (LIFO)
notifs.PopAll() // clear all active notifications
```

## Notification kinds

| Kind | Theme token | Use for |
|---|---|---|
| `NotificationKindInfo` | `NotificationInfo` | General messages |
| `NotificationKindSuccess` | `NotificationSuccess` | Completed actions |
| `NotificationKindWarning` | `NotificationWarning` | Reversible actions, cautions |
| `NotificationKindError` | `NotificationError` | Errors, failed operations |

## Thread safety

`Push`, `Pop`, and `PopAll` are safe to call from background goroutines. The `WithNotificationManager` option wires the re-render channel automatically so the UI re-renders when a banner expires, even if the user is not pressing any keys.

:::warning Persistent overlay, not a dialog
`NotificationManager` must be wired via `oat.WithNotificationManager`, not `ShowDialog`. It never blocks input and is never dismissed by Esc — it always renders on top.
:::
