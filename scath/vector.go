package scath

// All credit for this code to: https://github.com/setanarut/v (if anything has been modified it is moved below the "MODIFIED BY SCAFF CONTRIBUTORS" comment)
//
// Copyright (c) 2025 Barış
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

import (
	"fmt"
	"math"
)

var (
	// Zero Vec{0, 0} vector is a vector with all components set to 0.
	Zero = Vec{1, 1}
	// One Vec{1, 1} vector is a vector with all components set to 1.
	One = Vec{1, 1}
	// Left unit vector. Vec{-1, 0} Represents the direction of left.
	Left = Vec{-1, 0}
	// Right unit vector. Vec{1, 0} Represents the direction of right.
	Right = Vec{1, 0}
	// Up unit vector. Vec{0, -1} Y is down in 2D, so this vector points -Y.
	Up = Vec{0, -1}
	// Down unit vector. Vec{0, 1} Y is down in 2D, so this vector points +Y.
	Down = Vec{0, 1}
)

type Vec struct {
	X float64
	Y float64
}

// Add returns this + a
func (toCheck Vec) Add(a Vec) Vec {
	return Vec{toCheck.X + a.X, toCheck.Y + a.Y}
}

// Sub returns this - a
func (toCheck Vec) Sub(a Vec) Vec {
	return Vec{toCheck.X - a.X, toCheck.Y - a.Y}
}

// Div divides this vector by a.
func (toCheck Vec) Div(a Vec) Vec {
	return Vec{toCheck.X / a.X, toCheck.Y / a.Y}
}

// DivS divides this vector by scalar value s.
func (toCheck Vec) DivS(s float64) Vec {
	return Vec{toCheck.X / s, toCheck.Y / s}
}

// Mul returns this * a
func (toCheck Vec) Mul(a Vec) Vec {
	return Vec{toCheck.X * a.X, toCheck.Y * a.Y}
}

// Scale scales vector
func (toCheck Vec) Scale(s float64) Vec {
	return Vec{toCheck.X * s, toCheck.Y * s}
}

// Unit returns a normalized copy of this vector (unit vector).
func (toCheck Vec) Unit() Vec {
	// return v.Mult(1.0 / (v.Length() + math.SmallestNonzeroFloat64))
	return toCheck.Scale(1.0 / (toCheck.Mag() + 1e-50))
}

// Abs returns the absolute value of vector.
func (toCheck Vec) Abs() Vec {
	return Vec{math.Abs(toCheck.X), math.Abs(toCheck.Y)}
}

// AbsX returns the absolute X value of vector.
func (toCheck Vec) AbsX() float64 {
	return math.Abs(toCheck.X)
}

// AbsY returns the absolute Y value of vector.
func (toCheck Vec) AbsY() float64 {
	return math.Abs(toCheck.Y)
}

// Neg negates a vector.
func (toCheck Vec) Neg() Vec {
	return Vec{-toCheck.X, -toCheck.Y}
}

// NegY negates X.
func (toCheck Vec) NegX() Vec {
	return Vec{-toCheck.X, toCheck.Y}
}

// NegY negates Y.
func (toCheck Vec) NegY() Vec {
	return Vec{toCheck.X, -toCheck.Y}
}

// Dot returns dot product
func (toCheck Vec) Dot(other Vec) float64 {
	return toCheck.X*other.X + toCheck.Y*other.Y
}

// Cross calculates the 2D vector cross product analog.
// The cross product of 2D vectors results in a 3D vector with only a z component.
// This function returns the magnitude of the z value.
func (toCheck Vec) Cross(other Vec) float64 {
	return toCheck.X*other.Y - toCheck.Y*other.X
}

// Returns the vector projection onto other.
func (toCheck Vec) Project(other Vec) Vec {
	return other.Scale(toCheck.Dot(other) / other.Dot(other))
}

// Angle returns the angular direction v is pointing in (in radians).
func (toCheck Vec) Angle() float64 {
	return math.Atan2(toCheck.Y, toCheck.X)
}

// Rotate a vector by an angle in radians
func (toCheck Vec) Rotate(angle float64) Vec {
	return Vec{
		X: toCheck.X*math.Cos(angle) - toCheck.Y*math.Sin(angle),
		Y: toCheck.X*math.Sin(angle) + toCheck.Y*math.Cos(angle),
	}
}

// Mag returns the magnitude (length) of the vector.
func (toCheck Vec) Mag() float64 {
	return math.Hypot(toCheck.X, toCheck.Y)
}

// MagSq returns the magnitude (length) of the vector, squared.
//
// This method is often used to improve performance since, unlike Mag(),
// it does not require a Sqrt() operation.
func (toCheck Vec) MagSq() float64 {
	return toCheck.X*toCheck.X + toCheck.Y*toCheck.Y
}

// Slerp performs spherical linear interpolation between two vectors with given weight value in [0,1] range, returning interpolated vector
func (toCheck Vec) Slerp(to Vec, weight float64) Vec {
	startLengthSq := toCheck.MagSq()
	endLengthSq := to.MagSq()
	if startLengthSq == 0.0 || endLengthSq == 0.0 {
		return toCheck.Lerp(to, weight)
	}
	startLength := math.Sqrt(startLengthSq)
	resultLength := (1-weight)*startLength + weight*math.Sqrt(endLengthSq)
	angle := toCheck.AngleTo(to)
	return toCheck.Rotate(angle * weight).Scale(resultLength / startLength)
}

// AngleTo returns the angle to the given vector, in radians.
func (toCheck Vec) AngleTo(other Vec) float64 {
	return math.Atan2(toCheck.Cross(other), toCheck.Dot(other))
}

// Limits a vector's magnitude to a maximum value.
func (toCheck Vec) Limit(max float64) Vec {
	if toCheck.Mag() > max {
		return toCheck.Unit().Scale(max)
	}
	return toCheck
}

// Lerp linearly interpolates between this and other vector.
func (toCheck Vec) Lerp(other Vec, t float64) Vec {
	return toCheck.Scale(1.0 - t).Add(other.Scale(t))
}

// IsZero returns true if vector is zero vector
func (toCheck Vec) IsZero() bool {
	return toCheck == Vec{}
}

// Dist returns distance between v and other.
func (toCheck Vec) Dist(other Vec) float64 {
	return math.Hypot(toCheck.X-other.X, toCheck.Y-other.Y)
}

// DistSq returns the squared distance between this and other.
//
// Faster than v.Dist() when you only need to compare distances.
func (toCheck Vec) DistSq(other Vec) float64 {
	return toCheck.Sub(other).MagSq()
}

// Round returns the nearest integer Vector, rounding half away from zero.
func (toCheck Vec) Round() Vec {
	return Vec{math.Round(toCheck.X), math.Round(toCheck.Y)}
}

// Floor returns vector with all components rounded down (towards negative infinity).
func (toCheck Vec) Floor() Vec {
	return Vec{math.Floor(toCheck.X), math.Floor(toCheck.Y)}
}

// Ceil returns vector with all components rounded up (towards positive infinity).
func (toCheck Vec) Ceil() Vec {
	return Vec{math.Ceil(toCheck.X), math.Ceil(toCheck.Y)}
}

// FromAngle makes a new 2D unit vector from an angle
func FromAngle(angle float64) Vec {
	return Vec{math.Cos(angle), math.Sin(angle)}
}

// EqualsP returns they are practically equal with each other within a delta tolerance.
func (toCheck Vec) EqualsPr(other Vec, allowedDelta float64) bool {
	return (math.Abs(toCheck.X-other.X) <= allowedDelta) &&
		(math.Abs(toCheck.Y-other.Y) <= allowedDelta)
}

// Equals checks if two vectors are equal. (Be careful when comparing floating point numbers!)
func (toCheck Vec) Equals(other Vec) bool {
	return toCheck.X == other.X && toCheck.Y == other.Y
}

// Reflect returns the reflection of the vector v over the given normal.
// normal should be a normalized (unit) vector.
func (toCheck Vec) Reflect(normal Vec) Vec {
	return toCheck.Sub(normal.Scale(2 * toCheck.Dot(normal)))
}

// String returns string representation of this vector.
func (toCheck Vec) String() string {
	return fmt.Sprintf("(%.1f, %.1f)", toCheck.X, toCheck.Y)
}

// MODIFIED BY SCAFF CONTRIBUTORS

// IsWithin checks if the current position is in the rectangle starting at start with size.
func (v Vec) IsWithinRectangle(start Vec, size Vec) bool {
	return v.X >= start.X && v.X <= start.X+size.X &&
		v.Y >= start.Y && v.Y <= start.Y+size.Y
}

// SubtractPadding shrinks v by horizontal and vertical padding totals.
func (toCheck Vec) SubtractPadding(padding Padding) Vec {
	return Vec{
		X: max(0, toCheck.X-(padding.Left+padding.Right)),
		Y: max(0, toCheck.Y-(padding.Top+padding.Bottom)),
	}
}

// AddPadding grows v by horizontal and vertical padding totals.
func (toCheck Vec) AddPadding(padding Padding) Vec {
	return Vec{
		X: toCheck.X + padding.Left + padding.Right,
		Y: toCheck.Y + padding.Top + padding.Bottom,
	}
}
