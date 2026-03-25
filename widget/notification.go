package widget

import (
	"sync"
	"time"

	oat "github.com/antoniocali/oat-latte"
	"github.com/antoniocali/oat-latte/latte"
)

// NotificationKind controls the visual style of a notification banner.
type NotificationKind int

const (
	NotificationKindInfo    NotificationKind = iota // informational (blue/cyan)
	NotificationKindSuccess                         // positive outcome (green)
	NotificationKindWarning                         // caution (yellow/orange)
	NotificationKindError                           // failure / destructive (red)
)

// notificationEntry is one active notification in the manager's queue.
type notificationEntry struct {
	message string
	kind    NotificationKind
	expires time.Time // zero means never expires
}

// NotificationManager renders a stack of toast-style banners anchored to the
// bottom-right corner of whatever region it is given.
//
// Mount it via oat.WithNotificationManager, which wires the Canvas event loop
// and mounts the manager as a persistent overlay in a single step:
//
//	notifs := widget.NewNotificationManager()
//
//	app := oat.NewCanvas(
//	    oat.WithTheme(latte.ThemeDark),
//	    oat.WithBody(body),
//	    oat.WithNotificationManager(notifs),
//	)
//
//	// Show a timed notification:
//	notifs.Push("Task saved!", widget.NotificationKindSuccess, 3*time.Second)
//
// Each notification is a single line of text drawn with the appropriate
// theme colour.  Notifications with a non-zero duration are automatically
// removed from the queue; the manager sends to the Canvas notify channel
// so the screen refreshes without requiring a key press.
//
// NotificationManager is safe for concurrent use.
type NotificationManager struct {
	oat.BaseComponent
	mu      sync.Mutex
	entries []notificationEntry

	// theme styles — populated by ApplyTheme
	styleInfo    latte.Style
	styleSuccess latte.Style
	styleWarning latte.Style
	styleError   latte.Style

	// notifyCh is wired by Canvas via WithNotificationManager; used to trigger
	// re-renders on timer expiry.
	notifyCh chan<- time.Time
}

// NewNotificationManager creates an empty NotificationManager.
func NewNotificationManager() *NotificationManager {
	nm := &NotificationManager{}
	nm.EnsureID()
	return nm
}

// SetNotifyChannel implements oat.NotificationOverlay.  It is called
// automatically by oat.WithNotificationManager during Canvas construction —
// callers should not invoke this method directly.
func (nm *NotificationManager) SetNotifyChannel(ch chan<- time.Time) {
	nm.mu.Lock()
	nm.notifyCh = ch
	nm.mu.Unlock()
}

// Push adds a notification banner.
// If dur > 0 the notification is automatically dismissed after dur elapses.
// If dur == 0 the notification persists until dismissed by PopAll or a
// subsequent call to Pop.
func (nm *NotificationManager) Push(message string, kind NotificationKind, dur time.Duration) {
	nm.mu.Lock()
	entry := notificationEntry{
		message: message,
		kind:    kind,
	}
	if dur > 0 {
		entry.expires = time.Now().Add(dur)
		ch := nm.notifyCh
		nm.mu.Unlock()
		time.AfterFunc(dur, func() {
			nm.expire()
			if ch != nil {
				select {
				case ch <- time.Now():
				default:
				}
			}
		})
		nm.mu.Lock()
	}
	nm.entries = append(nm.entries, entry)
	nm.mu.Unlock()
}

// expire removes all entries whose deadline has passed.
func (nm *NotificationManager) expire() {
	now := time.Now()
	nm.mu.Lock()
	defer nm.mu.Unlock()
	kept := nm.entries[:0]
	for _, e := range nm.entries {
		if e.expires.IsZero() || e.expires.After(now) {
			kept = append(kept, e)
		}
	}
	nm.entries = kept
}

// Pop removes the most recently added notification (LIFO).
func (nm *NotificationManager) Pop() {
	nm.mu.Lock()
	defer nm.mu.Unlock()
	if len(nm.entries) > 0 {
		nm.entries = nm.entries[:len(nm.entries)-1]
	}
}

// PopAll clears all active notifications.
func (nm *NotificationManager) PopAll() {
	nm.mu.Lock()
	nm.entries = nm.entries[:0]
	nm.mu.Unlock()
}

// ApplyTheme maps the four Notification* theme tokens to internal styles.
func (nm *NotificationManager) ApplyTheme(t latte.Theme) {
	nm.styleInfo = t.NotificationInfo
	nm.styleSuccess = t.NotificationSuccess
	nm.styleWarning = t.NotificationWarning
	nm.styleError = t.NotificationError
}

// styleFor returns the latte.Style for the given kind.
func (nm *NotificationManager) styleFor(kind NotificationKind) latte.Style {
	switch kind {
	case NotificationKindSuccess:
		return nm.styleSuccess
	case NotificationKindWarning:
		return nm.styleWarning
	case NotificationKindError:
		return nm.styleError
	default:
		return nm.styleInfo
	}
}

// Measure returns the total height needed for all active notifications (one
// row each) and the width of the widest message + 4 padding chars.
func (nm *NotificationManager) Measure(c oat.Constraint) oat.Size {
	nm.expire()
	nm.mu.Lock()
	n := len(nm.entries)
	maxMsg := 0
	for _, e := range nm.entries {
		if len(e.message) > maxMsg {
			maxMsg = len(e.message)
		}
	}
	nm.mu.Unlock()

	if n == 0 {
		return oat.Size{}
	}
	w := maxMsg + 4 // 2 spaces padding on each side
	h := n
	if c.MaxWidth > 0 && w > c.MaxWidth {
		w = c.MaxWidth
	}
	if c.MaxHeight > 0 && h > c.MaxHeight {
		h = c.MaxHeight
	}
	return oat.Size{Width: w, Height: h}
}

// Render draws each notification as a right-aligned banner anchored to the
// bottom-right of region.  Newer notifications appear at the bottom.
func (nm *NotificationManager) Render(buf *oat.Buffer, region oat.Region) {
	nm.expire()
	nm.mu.Lock()
	entries := make([]notificationEntry, len(nm.entries))
	copy(entries, nm.entries)
	nm.mu.Unlock()

	if len(entries) == 0 {
		return
	}

	// Determine banner width: widest message + 4 (padding).
	maxMsg := 0
	for _, e := range entries {
		if len(e.message) > maxMsg {
			maxMsg = len(e.message)
		}
	}
	bannerW := maxMsg + 4
	if bannerW > region.Width {
		bannerW = region.Width
	}

	// Anchor to bottom-right.
	startY := region.Y + region.Height - len(entries)
	if startY < region.Y {
		startY = region.Y
	}
	startX := region.X + region.Width - bannerW
	if startX < region.X {
		startX = region.X
	}

	for i, e := range entries {
		y := startY + i
		if y >= region.Y+region.Height {
			break
		}
		style := nm.styleFor(e.kind)
		text := "  " + e.message + "  "
		// Pad to banner width.
		for len(text) < bannerW {
			text += " "
		}
		if len(text) > bannerW {
			text = text[:bannerW]
		}
		bannerRegion := oat.Region{X: startX, Y: y, Width: bannerW, Height: 1}
		buf.Sub(bannerRegion).DrawText(0, 0, text, style)
	}
}
