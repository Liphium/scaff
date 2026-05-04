package scaffui

import (
	"sync"
	"sync/atomic"
)

var _ Tracking = &Tracker{}

type Tracker struct {
	mu      *sync.Mutex
	changed atomic.Bool
	removal map[any]func()
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

// SetChanged marks the tracker as changed
func (t *Tracker) SetChanged() {
	t.changed.CompareAndSwap(false, true)
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
