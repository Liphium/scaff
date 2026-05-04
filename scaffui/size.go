package scaffui

import "github.com/Liphium/scaff/smath"

type Size struct {
	Width  int
	Height int
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
	return toCheck.X >= position.X && toCheck.X <= position.X+float64(size.Width) &&
		toCheck.Y >= position.Y && toCheck.Y <= position.Y+float64(size.Height)
}
