package scaffui

import "github.com/Liphium/scaff/smath"

type Size struct {
	Width  float64
	Height float64
}

// SubtractPadding shrinks size by horizontal and vertical padding totals.
func (s Size) SubtractPadding(padding Padding) Size {
	return Size{
		Width:  max(0, s.Width-(padding.Left+padding.Right)),
		Height: max(0, s.Height-(padding.Top+padding.Bottom)),
	}
}

// AddPadding grows size by horizontal and vertical padding totals.
func (s Size) AddPadding(padding Padding) Size {
	return Size{
		Width:  s.Width + padding.Left + padding.Right,
		Height: s.Height + padding.Top + padding.Bottom,
	}
}

func IsWithin(position smath.Vec, size Size, toCheck smath.Vec) bool {
	return toCheck.X >= position.X && toCheck.X <= position.X+size.Width &&
		toCheck.Y >= position.Y && toCheck.Y <= position.Y+size.Height
}
