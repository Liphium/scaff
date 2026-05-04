package smath_test

import (
	"testing"
	"time"

	"github.com/Liphium/scaff/smath"
	"github.com/stretchr/testify/assert"
)

func TestTimeframe(t *testing.T) {
	now := time.Now()
	f := smath.NewTimeframe(now, 1000)

	t.Run("works as expected", func(t *testing.T) {
		assert.Equal(t, 1000*time.Nanosecond, f.Remaining(now))
		assert.Equal(t, 500*time.Nanosecond, f.Remaining(now.Add(500)))
		assert.Equal(t, true, f.Started(now.Add(500)))
		assert.Equal(t, 0*time.Nanosecond, f.Remaining(now.Add(1000)))
		assert.Equal(t, true, f.Over(now.Add(1000)))
	})

	t.Run("bounds work", func(t *testing.T) {
		assert.Equal(t, 1000*time.Nanosecond, f.Remaining(now.Add(-1)))
		assert.Equal(t, false, f.Over(now.Add(-1)))
		assert.Equal(t, false, f.Started(now.Add(-1)))
		assert.Equal(t, 0*time.Nanosecond, f.Remaining(now.Add(1001)))
		assert.Equal(t, true, f.Over(now.Add(1001)))
	})
}

func TestTimeframeBackwards(t *testing.T) {
	now := time.Now()
	f := smath.NewTimeframe(now, 1000).MakeBackwards()

	t.Run("works as expected", func(t *testing.T) {
		assert.Equal(t, 0*time.Nanosecond, f.Remaining(now))
		assert.Equal(t, 500*time.Nanosecond, f.Remaining(now.Add(500)))
		assert.Equal(t, 1000*time.Nanosecond, f.Remaining(now.Add(1000)))
		assert.Equal(t, true, f.Over(now.Add(1000)))
	})

	t.Run("bounds work", func(t *testing.T) {
		assert.Equal(t, 0*time.Nanosecond, f.Remaining(now.Add(-1)))
		assert.Equal(t, false, f.Over(now.Add(-1)))
		assert.Equal(t, 1000*time.Nanosecond, f.Remaining(now.Add(1001)))
		assert.Equal(t, true, f.Over(now.Add(1001)))
	})
}

func TestTimeframeAddDelay(t *testing.T) {
	now := time.Now()
	f := smath.NewTimeframe(now, 1000).AddDelay(300 * time.Nanosecond)

	t.Run("works as expected", func(t *testing.T) {
		assert.Equal(t, 1000*time.Nanosecond, f.Remaining(now))
		assert.Equal(t, false, f.Started(now))
		assert.Equal(t, true, f.Started(now.Add(300)))
		assert.Equal(t, 1000*time.Nanosecond, f.Remaining(now.Add(300)))
		assert.Equal(t, 500*time.Nanosecond, f.Remaining(now.Add(800)))
		assert.Equal(t, 0*time.Nanosecond, f.Remaining(now.Add(1300)))
		assert.Equal(t, true, f.Over(now.Add(1300)))
		assert.Equal(t, true, f.Started(now.Add(1300)))
	})

	t.Run("works with backwards", func(t *testing.T) {
		f = f.MakeBackwards()
		assert.Equal(t, 0*time.Nanosecond, f.Remaining(now))
		assert.Equal(t, 0*time.Nanosecond, f.Remaining(now.Add(300)))
		assert.Equal(t, 500*time.Nanosecond, f.Remaining(now.Add(800)))
		assert.Equal(t, 1000*time.Nanosecond, f.Remaining(now.Add(1300)))
	})
}
