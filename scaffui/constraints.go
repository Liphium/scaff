package scaffui

import (
	"github.com/Liphium/scaff/scath"
	"errors"
	"math"
)

// Infinite represents an unbounded maximum size.
const Infinite float64 = -1

type Constraints struct {
	MinX  float64
	MaxX  float64
	MinY float64
	MaxY float64
}

func (c Constraints) RealMaxX() float64 {
	if c.MaxX == Infinite {
		return math.MaxFloat64
	}
	return c.MaxX
}

func (c Constraints) RealMaxY() float64 {
	if c.MaxY == Infinite {
		return math.MaxFloat64
	}
	return c.MaxY
}

func (c Constraints) IsTight() bool {
	return c.MinY == c.MaxY && c.MinX == c.MaxX
}

func (c Constraints) DoesXFit(w float64) bool {
	return c.MinX <= w && w <= c.RealMaxX()
}

func (c Constraints) DoesYFit(h float64) bool {
	return c.MinY <= h && h <= c.RealMaxY()
}

func (c Constraints) Fits(c2 Constraints) bool {
	widthFits := c.MinX <= c2.RealMaxX() && c2.MinX <= c.RealMaxX()
	heightFits := c.MinY <= c2.RealMaxY() && c2.MinY <= c.RealMaxY()
	return widthFits && heightFits
}

// Find the min size of constraints.
func (c Constraints) Min(horizontal bool) float64 {
	if horizontal {
		return c.MinX
	}
	return c.MinY
}

// Find the max size of constraints.
func (c Constraints) Max(horizontal bool) float64 {
	if horizontal {
		return c.MaxX
	}
	return c.MaxY
}

// Find the max size of constraints.
func (c Constraints) RealMax(horizontal bool) float64 {
	if horizontal {
		return c.RealMaxX()
	}
	return c.RealMaxY()
}

// SubtractPadding shrinks constraints by horizontal and vertical padding totals.
func (c Constraints) SubtractPadding(padding scath.Padding) Constraints {
	horizontal := padding.Left + padding.Right
	vertical := padding.Top + padding.Bottom

	minX := max(0, c.MinX-horizontal)

	maxX := c.MaxX
	if maxX != Infinite {
		maxX = max(0, maxX-horizontal)
	}

	minY := max(0, c.MinY-vertical)

	maxY := c.MaxY
	if maxY != Infinite {
		maxY = max(0, maxY-vertical)
	}

	return NewConstraints(
		minX,
		maxX,
		minY,
		maxY,
	)
}

func (c Constraints) TakeMaxWithin(c2 Constraints) (scath.Vec, error) {
	size := scath.Vec{
		X:  min(c.RealMaxX(), c2.RealMaxX()),
		Y: min(c.RealMaxY(), c2.RealMaxY()),
	}

	if !c.DoesXFit(size.X) || !c2.DoesXFit(size.X) || !c.DoesYFit(size.Y) || !c2.DoesYFit(size.Y) {
		return scath.Vec{}, errors.New("constraints could not fit")
	}

	return size, nil
}

// NewConstraints creates normalized width and height constraints.
func NewConstraints(minX, maxX, minY, maxY float64) Constraints {
	minX, maxX = normalizeConstraintRange(minX, maxX)
	minY, maxY = normalizeConstraintRange(minY, maxY)

	return Constraints{
		MinX:  minX,
		MaxX:  maxX,
		MinY: minY,
		MaxY: maxY,
	}
}

// Unconstrained returns constraints with no upper bound.
func Unconstrained() Constraints {
	return NewConstraints(0, Infinite, 0, Infinite)
}

// Tight returns constraints with an exact width and height.
func Tight(width, height float64) Constraints {
	return NewConstraints(width, width, height, height)
}

// TightFor returns tight constraints for provided axes.
func TightFor(width, height float64) Constraints {
	minX, maxX := 0.0, Infinite
	minY, maxY := 0.0, Infinite

	if width != Infinite {
		minX, maxX = width, width
	}

	if height != Infinite {
		minY, maxY = height, height
	}

	return NewConstraints(minX, maxX, minY, maxY)
}

// TightForFinite returns tight constraints only for finite axes.
func TightForFinite(width, height float64) Constraints {
	minX, maxX := 0.0, Infinite
	minY, maxY := 0.0, Infinite

	if width >= 0 {
		minX, maxX = width, width
	}

	if height >= 0 {
		minY, maxY = height, height
	}

	return NewConstraints(minX, maxX, minY, maxY)
}

// Loose returns constraints that are only bounded by maxima.
func Loose(maxX, maxY float64) Constraints {
	return NewConstraints(0, maxX, 0, maxY)
}

// Expand returns constraints that fill specified finite axes.
func Expand(width, height float64) Constraints {
	minX, maxX := 0.0, Infinite
	minY, maxY := 0.0, Infinite

	if width != Infinite {
		minX, maxX = width, width
	}

	if height != Infinite {
		minY, maxY = height, height
	}

	return NewConstraints(minX, maxX, minY, maxY)
}

func normalizeConstraintRange(min, max float64) (float64, float64) {
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
