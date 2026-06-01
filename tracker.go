package scaff

import (
	"sync"
	"sync/atomic"
)

type Tracking interface {
	// Should return the tracker for the current Node
	Tracker() *Tracker
}

var _ Tracking = &Tracker{}

// TODO: We need a method for actually, after the props are created, running all the effects cause otherwise those parts of the props won't be set, an ideal thing to do additionally would be to when the effects are ran, track which signals were added in each effect and then create a mapping between signal -> effect index, that way we could only re-run the effects that matter, for others we just re-run all effects obv.
type Tracker struct {
	mu      *sync.Mutex
	changed atomic.Bool
	removal map[any]func()
	effects []func()
}

// Implement Tracking interface
func (t *Tracker) Tracker() *Tracker {
	return t
}

func NewTracker() *Tracker {
	return &Tracker{
		mu:      &sync.Mutex{},
		changed: atomic.Bool{},
		removal: make(map[any]func()),
	}
}

// Changed reports whether any tracked signal emitted after the initial immediate push.
func (t *Tracker) Changed() bool {
	return t.changed.Load()
}

// SetChanged marks the tracker as changed (returns if anything changed)
func (t *Tracker) SetChanged() bool {
	return t.changed.CompareAndSwap(false, true)
}

// SetUnchanged marks the tracker as unchanged (returns if anything changed)
func (t *Tracker) SetUnchanged() bool {
	return t.changed.CompareAndSwap(true, false)
}

// Clear removes all tracked signals and is safe to call multiple times.
func (t *Tracker) Clear() {
	t.mu.Lock()
	for _, remove := range t.removal {
		remove()
	}
	t.removal = make(map[any]func())
	t.mu.Unlock()
}

// Effect adds a new handler being called when the thing changes, use this to actually react to changes of signals
func (t *Tracker) Effect(handler func()) {
	t.mu.Lock()
	t.effects = append(t.effects, handler)
	t.mu.Unlock()
}

// Update calls all the effects on the tracker to synchronize everything
func (t *Tracker) Update() {
	t.mu.Lock()
	effectsCopy := make([]func(), len(t.effects))
	copy(effectsCopy, t.effects)
	t.mu.Unlock()

	for _, effect := range t.effects {
		effect()
	}
}

// TrackValue ensures tracker is subscribed to signal and returns the current value; closed or nil inputs return zero value.
func TrackValue[T any](tracker *Tracker, signal *Signal[T]) T {
	if tracker == nil || signal == nil {
		var zero T
		return zero
	}

	tracker.mu.Lock()
	if tracker.removal == nil {
		tracker.removal = make(map[any]func())
	}
	_, exists := tracker.removal[signal]
	tracker.mu.Unlock()

	if !exists {
		initial := true
		remove := signal.AddListener(func(T) {
			if initial {
				initial = false
				return
			}
			tracker.SetChanged()
		})

		tracker.mu.Lock()
		if _, alreadyExists := tracker.removal[signal]; alreadyExists {
			tracker.mu.Unlock()
			remove()
		} else {
			tracker.removal[signal] = remove
			tracker.mu.Unlock()
		}
	}

	return signal.currentValue()
}
