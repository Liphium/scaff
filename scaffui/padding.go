package scaffui

import "github.com/Liphium/scaff/smath"

type Padding struct {
	Top    int
	Bottom int
	Right  int
	Left   int
}

// NewPadding creates padding with explicit top, right, bottom, and left values.
func NewPadding(top, right, bottom, left int) Padding {
	return Padding{
		Top:    top,
		Bottom: bottom,
		Right:  right,
		Left:   left,
	}
}

// Pad creates uniform padding on all sides.
func Pad(value int) Padding {
	return NewPadding(value, value, value, value)
}

// PadTop creates padding only on top side.
func PadTop(value int) Padding {
	return NewPadding(value, 0, 0, 0)
}

// PadBottom creates padding only on bottom side.
func PadBottom(value int) Padding {
	return NewPadding(0, 0, value, 0)
}

// PadRight creates padding only on right side.
func PadRight(value int) Padding {
	return NewPadding(0, value, 0, 0)
}

// PadLeft creates padding only on left side.
func PadLeft(value int) Padding {
	return NewPadding(0, 0, 0, value)
}

// PadHorizontal creates padding on left and right sides.
func PadHorizontal(value int) Padding {
	return NewPadding(0, value, 0, value)
}

// PadVertical creates padding on top and bottom sides.
func PadVertical(value int) Padding {
	return NewPadding(value, 0, value, 0)
}

// ToVecTopLeft converts top and left padding into position vector.
func (p Padding) ToVecTopLeft() smath.Vec {
	return smath.Vec{X: float64(p.Left), Y: float64(p.Top)}
}
