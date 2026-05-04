package basenode

import (
	"github.com/Liphium/scaff"
)

var log = scaff.NewLogger("basenode")

// Deletes an element without preserving order.
func deleteUnordered[T any](s []T, i int) []T {
	last := len(s) - 1
	s[i] = s[last]
	var zero T
	s[last] = zero // avoid holding refs for GC
	return s[:last]
}

type LayoutDirection uint

const (
	LayoutTopToBottom LayoutDirection = 0
	LayoutBottomToTop LayoutDirection = 1
	LayoutLeftToRight LayoutDirection = 2
	LayoutRightToLeft LayoutDirection = 3
)

type HorizontalAlignment uint

const (
	HorizontalAlignmentRight  HorizontalAlignment = 0
	HorizontalAlignmentCenter HorizontalAlignment = 1
	HorizontalAlignmentLeft   HorizontalAlignment = 2
)

type VerticalAlignment uint

const (
	VerticalAlignmentTop    VerticalAlignment = 0
	VerticalAlignmentCenter VerticalAlignment = 1
	VerticalAlignmentBottom VerticalAlignment = 2
)
