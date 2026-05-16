package scath

import (
	"math"
	"time"
)

// Linear interpolation between two floating-point numbers.
func (f Timeframe) LerpFloat(now time.Time, a, b float64) float64 {
	t := float64(float64(f.Remaining(now)) / float64(f.duration))
	return b + (a-b)*t
}

// Linear interpolation between two integers.
func (f Timeframe) LerpInt(now time.Time, a, b int) int {
	return int(math.RoundToEven(f.LerpFloat(now, float64(a), float64(b))))
}

// Linear interpolation between two unsigned integers.
func (f Timeframe) LerpUint(now time.Time, a, b uint) uint {
	return uint(math.RoundToEven(f.LerpFloat(now, float64(a), float64(b))))
}

// LerpTypewriter writes a string with a typewriter effect, constructing it over time.
func (f Timeframe) LerpTypewriter(now time.Time, a string) string {
	progress := f.LerpInt(now, 0, len(a))
	return a[:progress]
}

// LerpTypewriterBetween writes a string with a typewriter effect between two strings, deleting the first and then writing the second.
func (f Timeframe) LerpTypewriterBetween(now time.Time, a, b string) string {
	progress := f.LerpInt(now, -len(a), len(b))
	if progress <= 0 {
		return a[:-progress]
	}
	return b[:progress]
}
