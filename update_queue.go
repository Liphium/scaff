package scaff

import "sync"

type UpdateQueue struct {
	mu        sync.Mutex
	callbacks map[*Tracker]map[int]func()
}

func NewUpdateQueue() *UpdateQueue {
	return &UpdateQueue{
		callbacks: make(map[*Tracker]map[int]func()),
	}
}

// Push adds a callback to the queue, it's indexed by the tracker pointer and effect ID to avoid duplicates
func (u *UpdateQueue) Push(tracker *Tracker, effect int, callback func()) {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.callbacks == nil {
		u.callbacks = make(map[*Tracker]map[int]func())
	}

	trackerCallbacks, ok := u.callbacks[tracker]
	if !ok {
		trackerCallbacks = make(map[int]func())
		u.callbacks[tracker] = trackerCallbacks
	}

	trackerCallbacks[effect] = callback
}

// Update runs all the callbacks in the queue and then clears it
func (u *UpdateQueue) Update() {
	u.mu.Lock()
	callbacks := u.callbacks
	u.callbacks = make(map[*Tracker]map[int]func())
	u.mu.Unlock()

	for tracker, trackerCallbacks := range callbacks {
		// If -1 is present, it means everything should be updated for this tracker.
		// We prioritize it and skip other individual effects for this tracker.
		if cb, ok := trackerCallbacks[-1]; ok {
			cb()
			tracker.runChangeHandlers()
			continue
		}

		for _, cb := range trackerCallbacks {
			cb()
		}
		tracker.runChangeHandlers()
	}
}

// Clear removes all callbacks from the queue without executing them
func (u *UpdateQueue) Clear() {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.callbacks = make(map[*Tracker]map[int]func())
}
