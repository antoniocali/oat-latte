---
sidebar_position: 8
title: NotificationManager
description: Toast notification banners with auto-dismiss.
---

# NotificationManager

`widget.NotificationManager` displays transient toast banners anchored to the bottom-right corner of the screen. Banners can auto-dismiss after a duration or stay until dismissed manually.

`NotificationManager` has no style constructor argument and no `WithStyle` builder — its appearance is entirely controlled by the four `Notification*` theme tokens.

## Setup

Mount the `NotificationManager` as a **persistent overlay** after constructing the canvas. Use `ShowPersistentOverlay` — it renders on top of everything else and is never dismissed by Esc.

```go
notifs := widget.NewNotificationManager()

app := oat.NewCanvas(
    oat.WithTheme(latte.ThemeDark),
    oat.WithBody(body),
)

// Connect the timer channel so expiring notifications trigger re-renders.
notifs.SetNotifyChannel(app.NotifyChannel())

// Mount permanently on top of the UI (never dismissed by Esc).
app.ShowPersistentOverlay(notifs)
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

`Push`, `Pop`, and `PopAll` are safe to call from background goroutines. The `NotifyChannel` connection ensures the UI re-renders when a banner expires, even if the user is not pressing any keys.

:::warning Persistent overlay, not a dialog
`NotificationManager` must be mounted with `ShowPersistentOverlay`, not `ShowDialog`. Unlike a dialog it never blocks input and is never dismissed by Esc — it always renders on top.
:::
