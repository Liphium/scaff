package scath

type Padding struct {
	Top    float64
	Bottom float64
	Right  float64
	Left   float64
}

// NewPadding creates padding with explicit top, right, bottom, and left values.
func NewPadding(top, right, bottom, left float64) Padding {
	return Padding{
		Top:    top,
		Bottom: bottom,
		Right:  right,
		Left:   left,
	}
}

// Pad creates uniform padding on all sides.
func Pad(value float64) Padding {
	return NewPadding(value, value, value, value)
}

// PadTop creates padding only on top side.
func PadTop(value float64) Padding {
	return NewPadding(value, 0, 0, 0)
}

// PadBottom creates padding only on bottom side.
func PadBottom(value float64) Padding {
	return NewPadding(0, 0, value, 0)
}

// PadRight creates padding only on right side.
func PadRight(value float64) Padding {
	return NewPadding(0, value, 0, 0)
}

// PadLeft creates padding only on left side.
func PadLeft(value float64) Padding {
	return NewPadding(0, 0, 0, value)
}

// PadHorizontal creates padding on left and right sides.
func PadHorizontal(value float64) Padding {
	return NewPadding(0, value, 0, value)
}

// PadVertical creates padding on top and bottom sides.
func PadVertical(value float64) Padding {
	return NewPadding(value, 0, value, 0)
}

// ToVecTopLeft converts top and left padding into position vector.
func (p Padding) ToVecTopLeft() Vec {
	return Vec{X: p.Left, Y: p.Top}
}
