package scaff_test

import (
	"testing"
	"time"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/optional"
	"github.com/Liphium/scaff/smath"
	"github.com/stretchr/testify/assert"
)

func TestTransitioningGeneral(t *testing.T) {
	now := time.Now()
	loads := []string{}
	unloads := []string{}
	updates := []string{}

	scene1 := createMockScene("test1", nil, &loads, &unloads)
	scene2 := createMockScene("test2", nil, &loads, &unloads)
	scene1.TransitionFunc = func(in bool) scaff.TransitionProperties {
		return scaff.TransitionProperties{
			Isolated: false,
			Duration: 100 * time.Nanosecond,
		}
	}
	scene2.TransitionFunc = func(in bool) scaff.TransitionProperties {
		return scaff.TransitionProperties{
			Isolated: false,
			Duration: 100 * time.Nanosecond,
		}
	}
	state := scaff.NewTransitioningState(now, scene1)
	state.Update(now, func(ts *TestScene, f smath.Timeframe) error {
		assert.Equal(t, "test1", ts.ID)
		assert.Equal(t, 100*time.Nanosecond, f.Remaining(now))
		assert.Equal(t, false, f.IsBackwards())
		return nil
	})

	t.Run("basic transitions work", func(t *testing.T) {

		// Add new scene and make sure it's actually getting updated together with the other one
		state.Set(now.Add(100), optional.With(scene2))
		expectedBackwards := map[string]bool{
			"test1": true,
			"test2": false,
		}
		state.Update(now.Add(100), func(ts *TestScene, f smath.Timeframe) error {
			updates = append(updates, ts.GetId())
			assert.Equal(t, expectedBackwards[ts.GetId()], f.IsBackwards())
			return nil
		})
		assert.Equal(t, []string{"test1", "test2"}, updates)

		// Other transition should be over now
		state.Update(now.Add(200), func(ts *TestScene, f smath.Timeframe) error {
			assert.Equal(t, "test2", ts.ID)
			assert.Equal(t, 0*time.Nanosecond, f.Remaining(now.Add(200)))
			assert.Equal(t, false, f.IsBackwards())
			return nil
		})

		// Make sure loads and unloads worked properly
		assert.Equal(t, []string{"test1", "test2"}, loads)
		assert.Equal(t, []string{"test1"}, unloads)
	})

	t.Run("transitions back to back", func(t *testing.T) {

		// Insert scene 1 again and make sure it's actually in there
		state.Set(now.Add(200), optional.With(scene1))
		updates = []string{}
		state.Update(now.Add(200), func(ts *TestScene, f smath.Timeframe) error {
			updates = append(updates, ts.GetId())
			return nil
		})
		assert.Equal(t, []string{"test2", "test1"}, updates)

		// Make sure the thing is gone at 300
		state.Update(now.Add(300), func(ts *TestScene, f smath.Timeframe) error {
			assert.Equal(t, "test1", ts.ID)
			assert.Equal(t, 0*time.Nanosecond, f.Remaining(now.Add(300)))
			return nil
		})

		// Make sure loads and unloads worked properly
		assert.Equal(t, []string{"test1", "test2", "test1"}, loads)
		assert.Equal(t, []string{"test1", "test2"}, unloads)
	})
}

func TestTransitionCollision(t *testing.T) {
	now := time.Now()
	loads := []string{}
	unloads := []string{}

	scene1 := createMockScene("test1", nil, &loads, &unloads)
	scene2 := createMockScene("test2", nil, &loads, &unloads)
	scene3 := createMockScene("test3", nil, &loads, &unloads)

	state := scaff.NewTransitioningState(now, scene1)
	state.Set(now, optional.With(scene2))
	state.Set(now, optional.With(scene3))

	// All scenes should be completely gone
	state.Update(now, func(ts *TestScene, f smath.Timeframe) error {
		assert.Equal(t, "test3", ts.ID)
		return nil
	})

	// Make sure loads and unloads worked properly
	assert.Equal(t, []string{"test1", "test2", "test3"}, loads)
	assert.Equal(t, []string{"test1", "test2"}, unloads)
}

func TestTransitioningIsolated(t *testing.T) {
	now := time.Now()
	loads := []string{}
	unloads := []string{}

	scene1 := createMockScene("test1", nil, &loads, &unloads)
	scene2 := createMockScene("test2", nil, &loads, &unloads)
	scene1.TransitionFunc = func(in bool) scaff.TransitionProperties {
		return scaff.TransitionProperties{
			Isolated: true,
			Duration: 100 * time.Nanosecond,
		}
	}
	scene2.TransitionFunc = func(in bool) scaff.TransitionProperties {
		return scaff.TransitionProperties{
			Isolated: true,
			Duration: 100 * time.Nanosecond,
		}
	}
	state := scaff.NewTransitioningState(now, scene1)
	state.Update(now, func(ts *TestScene, f smath.Timeframe) error {
		assert.Equal(t, "test1", ts.ID)
		assert.Equal(t, 100*time.Nanosecond, f.Remaining(now))
		assert.Equal(t, false, f.IsBackwards())
		return nil
	})

	// At 100, scene1 should start to transition out
	state.Set(now.Add(100), optional.With(scene2))
	expectedRemaining := map[string]time.Duration{
		"test1": 0 * time.Nanosecond,
		"test2": 100 * time.Nanosecond,
	}
	updates := []string{}
	state.Update(now.Add(100), func(ts *TestScene, f smath.Timeframe) error {
		updates = append(updates, ts.GetId())
		assert.Equal(t, expectedRemaining[ts.GetId()], f.Remaining(now.Add(100)))
		return nil
	})
	assert.Equal(t, []string{"test1", "test2"}, updates)

	// At 200, test1 should have transitioned out, but still kept around because the other scene is not transitioned in yet, that should start now
	expectedRemaining = map[string]time.Duration{
		"test1": 100 * time.Nanosecond,
		"test2": 100 * time.Nanosecond,
	}
	updates = []string{}
	state.Update(now.Add(200), func(ts *TestScene, f smath.Timeframe) error {
		updates = append(updates, ts.GetId())
		assert.Equal(t, expectedRemaining[ts.GetId()], f.Remaining(now.Add(200)))
		return nil
	})
	assert.Equal(t, []string{"test1", "test2"}, updates)

	// At 300, test1 should be fully transitioned in and still be there
	state.Update(now.Add(300), func(ts *TestScene, f smath.Timeframe) error {
		assert.Equal(t, "test2", ts.GetId())
		assert.Equal(t, 0*time.Nanosecond, f.Remaining(now.Add(300)))
		return nil
	})

	// Make sure loads and unloads worked properly
	assert.Equal(t, []string{"test1", "test2"}, loads)
	assert.Equal(t, []string{"test1"}, unloads)
}

func TestTransitioningEmpty(t *testing.T) {
	now := time.Now()
	loads := []string{}
	unloads := []string{}

	scene1 := createMockScene("test1", nil, &loads, &unloads)
	scene1.TransitionFunc = func(in bool) scaff.TransitionProperties {
		return scaff.TransitionProperties{
			Isolated: true,
			Duration: 100 * time.Nanosecond,
		}
	}
	state := scaff.NewTransitioningState(now, scene1)
	state.Update(now, func(ts *TestScene, f smath.Timeframe) error {
		assert.Equal(t, "test1", ts.ID)
		assert.Equal(t, 100*time.Nanosecond, f.Remaining(now))
		assert.Equal(t, false, f.IsBackwards())
		return nil
	})

	t.Run("proper out when empty", func(t *testing.T) {
		// Set empty and make sure test1 starts transitioning out
		state.Set(now.Add(200), optional.None[*TestScene]())
		state.Update(now.Add(200), func(ts *TestScene, f smath.Timeframe) error {
			assert.Equal(t, "test1", ts.ID)
			assert.Equal(t, 0*time.Nanosecond, f.Remaining(now.Add(200)))
			assert.Equal(t, 100*time.Nanosecond, f.Remaining(now.Add(300)))
			assert.Equal(t, true, f.IsBackwards())
			return nil
		})
	})
}
