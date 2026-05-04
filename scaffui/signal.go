package scaffui

import (
	"maps"
	"sync"
)

// Create a new signal
func NewSignal[T any](def T) *Signal[T] {
	return &Signal[T]{
		mu:        &sync.Mutex{},
		value:     def,
		listeners: make(map[uint64]func(T)),
		nextID:    1,
	}
}

type Signal[T any] struct {
	mu        *sync.Mutex
	value     T
	listeners map[uint64]func(T)
	nextID    uint64
}

// Get the value of the signal (IF YOU WANT THE UI TO UPDATE USE TRACK WITH THE CORRECT TRACKER)
func (s *Signal[T]) Value() T {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.value
}

// Set updates the value and notifies all current listeners; listeners run synchronously and can block callers.
func (s *Signal[T]) Set(value T) {
	s.mu.Lock()
	s.value = value
	listeners := maps.Clone(s.listeners)
	s.mu.Unlock()

	for _, listener := range listeners {
		listener(value)
	}
}

// Refresh re-emits the current value to all listeners without changing it; listeners run synchronously.
func (s *Signal[T]) Refresh() {
	s.mu.Lock()
	value := s.value
	listeners := maps.Clone(s.listeners)
	s.mu.Unlock()

	for _, listener := range listeners {
		listener(value)
	}
}

// AddListener registers a listener and immediately pushes the current value; it returns an unsubscribe function.
func (s *Signal[T]) AddListener(listener func(T)) func() {
	if listener == nil {
		return func() {}
	}

	s.mu.Lock()
	if s.listeners == nil {
		s.listeners = make(map[uint64]func(T))
	}
	id := s.nextID
	s.nextID++
	s.listeners[id] = listener
	value := s.value
	s.mu.Unlock()

	listener(value)

	return func() {
		s.removeListener(id)
	}
}

// removeListener unregisters a listener id if present; it is a no-op for unknown ids.
func (s *Signal[T]) removeListener(id uint64) {
	s.mu.Lock()
	delete(s.listeners, id)
	s.mu.Unlock()
}

// currentValue returns the latest value snapshot; it is internal and intentionally avoids exposing a public getter.
func (s *Signal[T]) currentValue() T {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.value
}

// Convenience wrapper over cgui.TrackValue (tracking is implemented by every Node object by default, you can just pass it in)
func (s *Signal[T]) Track(tracking Tracking) T {
	return TrackValue(tracking.Tracker(), s)
}
