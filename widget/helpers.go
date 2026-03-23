package widget

// helpers.go — shared utilities for the widget package.

import (
	oat "github.com/antoniocali/oat-latte"
	"github.com/antoniocali/oat-latte/latte"
)

// toOatInsets converts a latte.Insets to an oat.Insets.
// The two types are structurally identical; this helper avoids a circular import.
func toOatInsets(i latte.Insets) oat.Insets {
	return oat.Insets{Top: i.Top, Right: i.Right, Bottom: i.Bottom, Left: i.Left}
}

// clamp returns v clamped to [lo, hi].
func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
