package scaffui

import (
	"errors"
	"math"
)

// Infinite represents an unbounded maximum size.
const Infinite = -1

type Constraints struct {
	MinWidth  int
	MaxWidth  int
	MinHeight int
	MaxHeight int
}

func (c Constraints) RealMaxWidth() int {
	if c.MaxWidth == Infinite {
		return math.MaxInt32
	}
	return c.MaxWidth
}

func (c Constraints) RealMaxHeight() int {
	if c.MaxHeight == Infinite {
		return math.MaxInt32
	}
	return c.MaxHeight
}

func (c Constraints) IsTight() bool {
	return c.MinHeight == c.MaxHeight && c.MinWidth == c.MaxWidth
}

func (c Constraints) DoesWidthFit(w int) bool {
	return c.MinWidth <= w && w <= c.RealMaxWidth()
}

func (c Constraints) DoesHeightFit(h int) bool {
	return c.MinHeight <= h && h <= c.RealMaxHeight()
}

func (c Constraints) Fits(c2 Constraints) bool {
	widthFits := c.MinWidth <= c2.RealMaxWidth() && c2.MinWidth <= c.RealMaxWidth()
	heightFits := c.MinHeight <= c2.RealMaxHeight() && c2.MinHeight <= c.RealMaxHeight()
	return widthFits && heightFits
}

// Find the min size of constraints.
func (c Constraints) Min(horizontal bool) int {
	if horizontal {
		return c.MinWidth
	}
	return c.MinHeight
}

// Find the max size of constraints.
func (c Constraints) Max(horizontal bool) int {
	if horizontal {
		return c.MaxWidth
	}
	return c.MaxHeight
}

// Find the max size of constraints.
func (c Constraints) RealMax(horizontal bool) int {
	if horizontal {
		return c.RealMaxWidth()
	}
	return c.RealMaxHeight()
}

// SubtractPadding shrinks constraints by horizontal and vertical padding totals.
func (c Constraints) SubtractPadding(padding Padding) Constraints {
	horizontal := padding.Left + padding.Right
	vertical := padding.Top + padding.Bottom

	minWidth := max(0, c.MinWidth-horizontal)

	maxWidth := c.MaxWidth
	if maxWidth != Infinite {
		maxWidth = max(0, maxWidth-horizontal)
	}

	minHeight := max(0, c.MinHeight-vertical)

	maxHeight := c.MaxHeight
	if maxHeight != Infinite {
		maxHeight = max(0, maxHeight-vertical)
	}

	return NewConstraints(
		minWidth,
		maxWidth,
		minHeight,
		maxHeight,
	)
}

func (c Constraints) TakeMaxWithin(c2 Constraints) (Size, error) {
	size := Size{
		Width:  min(c.RealMaxWidth(), c2.RealMaxWidth()),
		Height: min(c.RealMaxHeight(), c2.RealMaxHeight()),
	}

	if !c.DoesWidthFit(size.Width) || !c2.DoesWidthFit(size.Width) || !c.DoesHeightFit(size.Height) || !c2.DoesHeightFit(size.Height) {
		return Size{}, errors.New("constraints could not fit")
	}

	return size, nil
}

// NewConstraints creates normalized width and height constraints.
func NewConstraints(minWidth, maxWidth, minHeight, maxHeight int) Constraints {
	minWidth, maxWidth = normalizeConstraintRange(minWidth, maxWidth)
	minHeight, maxHeight = normalizeConstraintRange(minHeight, maxHeight)

	return Constraints{
		MinWidth:  minWidth,
		MaxWidth:  maxWidth,
		MinHeight: minHeight,
		MaxHeight: maxHeight,
	}
}

// Unconstrained returns constraints with no upper bound.
func Unconstrained() Constraints {
	return NewConstraints(0, Infinite, 0, Infinite)
}

// Tight returns constraints with an exact width and height.
func Tight(width, height int) Constraints {
	return NewConstraints(width, width, height, height)
}

// TightFor returns tight constraints for provided axes.
func TightFor(width, height int) Constraints {
	minWidth, maxWidth := 0, Infinite
	minHeight, maxHeight := 0, Infinite

	if width != Infinite {
		minWidth, maxWidth = width, width
	}

	if height != Infinite {
		minHeight, maxHeight = height, height
	}

	return NewConstraints(minWidth, maxWidth, minHeight, maxHeight)
}

// TightForFinite returns tight constraints only for finite axes.
func TightForFinite(width, height int) Constraints {
	minWidth, maxWidth := 0, Infinite
	minHeight, maxHeight := 0, Infinite

	if width >= 0 {
		minWidth, maxWidth = width, width
	}

	if height >= 0 {
		minHeight, maxHeight = height, height
	}

	return NewConstraints(minWidth, maxWidth, minHeight, maxHeight)
}

// Loose returns constraints that are only bounded by maxima.
func Loose(maxWidth, maxHeight int) Constraints {
	return NewConstraints(0, maxWidth, 0, maxHeight)
}

// Expand returns constraints that fill specified finite axes.
func Expand(width, height int) Constraints {
	minWidth, maxWidth := 0, Infinite
	minHeight, maxHeight := 0, Infinite

	if width != Infinite {
		minWidth, maxWidth = width, width
	}

	if height != Infinite {
		minHeight, maxHeight = height, height
	}

	return NewConstraints(minWidth, maxWidth, minHeight, maxHeight)
}

func normalizeConstraintRange(min, max int) (int, int) {
	if min == Infinite {
		min = 0
	}

	if min < Infinite {
		min = 0
	}

	if max < Infinite {
		max = Infinite
	}

	if max != Infinite && min > max {
		max = min
	}

	return min, max
}
