package scaff

import (
	"maps"
	"sync"
	"sync/atomic"
)

// Internal change interface for all changes that can happen to the internal map of the store.
type change interface{}

// Update operation on the map, can also be used as add.
type update[K comparable, E any] struct {
	Key     K
	Entity  E
	Version int32
}

// Delete operation on the map.
type deletion[K comparable, E any] struct {
	Key     K
	version int32
}

// Store is a versioned, thread-safe and really fast map. Read below to learn about it's internal architecture.
//
// The render thread referred to below here is just the thread where you'll always be updating the Store from. This can not be called from multiple goroutines.
//
// Idea: Have one map that is for general consumption and always up-to-date. All operations will be executed on that one, locked with mutex as usual and stuff. At the same time keep a difference + version number to synchronize the render thread in a basically non-blocking way.
type Store[K comparable, V any] struct {
	// Mirror for outside (we lock this with a mutex as it can be accessed from multiple places)
	mu      sync.Mutex
	current map[K]V

	// State for the render thread
	version  atomic.Int32
	pending  chan change
	rendered map[K]V
}

// Create a new Store.
//
// buffer is how much space the versioning queue will have, if the buffer is exhausted all basic operations like add will start blocking. Define this based on how often you calculate the difference and how often things change between those times.
func NewStore[K comparable, V any](buffer int) *Store[K, V] {
	return &Store[K, V]{
		mu:       sync.Mutex{},
		current:  map[K]V{},
		version:  atomic.Int32{},
		pending:  make(chan change, buffer),
		rendered: map[K]V{},
	}
}

func (e *Store[K, V]) Set(key K, value V) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.current[key] = value

	e.pending <- update[K, V]{
		Key:     key,
		Entity:  value,
		Version: e.version.Load(),
	}
}

// GetCopy returns a copy of the entire internal map.
func (s *Store[K, V]) GetCopy() map[K]V {
	s.mu.Lock()
	defer s.mu.Unlock()

	return maps.Clone(s.current)
}

func (e *Store[K, V]) Clear() {
	// I'm not sure if we actually need the mutex here, thorough analysis would be needed, but should not be a huge performance bottleneck
	e.mu.Lock()
	defer e.mu.Unlock()

	e.version.Add(1)
}

// TODO: Functions for safe iteration: Provide editing functionality directly in the functions to make sure no-one gets the idea to like deadlock if we just return the map + lock mutex or sth.

// THIS METHOD IS NOT CONCURRENCY-SAFE. DO ONLY CALL FROM ONE GOROUTINE REGULARLY TO BE SAFE.
//
// Should diff + deduplicate.
func (e *Store[K, V]) Update() {

}
