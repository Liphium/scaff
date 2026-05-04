package smath_test

import (
	"testing"
	"time"

	"github.com/Liphium/scaff/smath"
	"github.com/stretchr/testify/assert"
)

func TestLerpFloat(t *testing.T) {
	now := time.Now()
	f := smath.NewTimeframe(now, 1000)

	t.Run("normal values work", func(t *testing.T) {
		assert.Equal(t, 0.0, f.LerpFloat(now, 0.0, 1.0))
		assert.Equal(t, 0.5, f.LerpFloat(now.Add(500), 0.0, 1.0))
		assert.Equal(t, 1.0, f.LerpFloat(now.Add(1000), 0.0, 1.0))
	})

	t.Run("bounds are kept", func(t *testing.T) {
		assert.Equal(t, 0.0, f.LerpFloat(now.Add(-1), 0.0, 1.0))
		assert.Equal(t, 1.0, f.LerpFloat(now.Add(1001), 0.0, 1.0))
	})
}

func TestLerpInt(t *testing.T) {
	now := time.Now()
	f := smath.NewTimeframe(now, 1000)

	t.Run("normal values work", func(t *testing.T) {
		assert.Equal(t, 1, f.LerpInt(now, 1, 100))
		assert.Equal(t, 50, f.LerpInt(now.Add(500), 1, 100))
		assert.Equal(t, 100, f.LerpInt(now.Add(1000), 1, 100))
	})

	t.Run("bounds are kept", func(t *testing.T) {
		assert.Equal(t, 1, f.LerpInt(now.Add(-1), 1, 100))
		assert.Equal(t, 100, f.LerpInt(now.Add(1001), 1, 100))
	})
}

func TestLerpUint(t *testing.T) {
	now := time.Now()
	f := smath.NewTimeframe(now, 1000)

	t.Run("normal values work", func(t *testing.T) {
		assert.Equal(t, uint(1), f.LerpUint(now, uint(1), uint(100)))
		assert.Equal(t, uint(50), f.LerpUint(now.Add(500), uint(1), uint(100)))
		assert.Equal(t, uint(100), f.LerpUint(now.Add(1000), uint(1), uint(100)))
	})

	t.Run("bounds are kept", func(t *testing.T) {
		assert.Equal(t, uint(1), f.LerpUint(now.Add(-1), uint(1), uint(100)))
		assert.Equal(t, uint(100), f.LerpUint(now.Add(1001), uint(1), uint(100)))
	})
}

func TestLerpTypewriter(t *testing.T) {
	now := time.Now()
	f := smath.NewTimeframe(now, 1000)

	t.Run("normal values work", func(t *testing.T) {
		assert.Equal(t, "", f.LerpTypewriter(now, "ab"))
		assert.Equal(t, "a", f.LerpTypewriter(now.Add(500), "ab"))
		assert.Equal(t, "ab", f.LerpTypewriter(now.Add(1000), "ab"))
	})

	t.Run("bounds are kept", func(t *testing.T) {
		assert.Equal(t, "", f.LerpTypewriter(now.Add(-1), "ab"))
		assert.Equal(t, "ab", f.LerpTypewriter(now.Add(1001), "ab"))
	})
}

func TestLerpTypewriterBetween(t *testing.T) {
	now := time.Now()
	f := smath.NewTimeframe(now, 1000)

	t.Run("normal values work", func(t *testing.T) {
		assert.Equal(t, "a", f.LerpTypewriterBetween(now, "a", "b"))
		assert.Equal(t, "", f.LerpTypewriterBetween(now.Add(500), "a", "b"))
		assert.Equal(t, "b", f.LerpTypewriterBetween(now.Add(1000), "a", "b"))
	})

	t.Run("bounds are kept", func(t *testing.T) {
		assert.Equal(t, "a", f.LerpTypewriterBetween(now.Add(-1), "a", "b"))
		assert.Equal(t, "b", f.LerpTypewriterBetween(now.Add(1001), "a", "b"))
	})
}
