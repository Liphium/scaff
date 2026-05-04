package smath

import "time"

// NewTimeframe creates a new Timeframe with the given start time and duration.
func NewTimeframe(start time.Time, duration time.Duration) Timeframe {
	return Timeframe{
		duration: duration,
		end:      start.Add(duration),
	}
}

// Timeframe represents a duration with a start and end time.
type Timeframe struct {
	duration  time.Duration
	end       time.Time
	backwards bool
}

// MakeBackwards sets the timeframe to be backwards (end < start).
func (f Timeframe) MakeBackwards() Timeframe {
	f.backwards = true
	return f
}

// MakeForwards sets the timeframe to be forwards (end > start).
func (f Timeframe) MakeForwards() Timeframe {
	f.backwards = false
	return f
}

// IsBackwards returns true if the timeframe is backwards.
func (f Timeframe) IsBackwards() bool {
	return f.backwards
}

// AddDelay adds a delay to the end of the timeframe.
func (f Timeframe) AddDelay(delay time.Duration) Timeframe {
	f.end = f.end.Add(delay)
	return f
}

// Remaining returns the remaining duration until the end of the timeframe.
//
// Also clamps to 0 if the current time is after the end of the timeframe. Or to the duration.
func (f Timeframe) Remaining(now time.Time) time.Duration {
	rem := f.end.Sub(now)
	if f.backwards {
		rem = f.duration - rem
	}
	if rem < 0 {
		return 0
	}
	if rem > f.duration {
		return f.duration
	}
	return rem
}

// Over returns true if the current time is after the end of the timeframe.
func (f Timeframe) Over(now time.Time) bool {
	return now.After(f.end) || now.Equal(f.end)
}

// Started returns true if the current time is after the start of the timeframe.
func (f Timeframe) Started(now time.Time) bool {
	start := f.end.Add(-f.duration)
	return now.After(start) || now.Equal(start)
}
