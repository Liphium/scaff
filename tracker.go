package scaff

import (
	"slices"
	"sync"
)

type Tracking interface {
	// Should return the tracker for the current Node
	Tracker() *Tracker
}

var _ Tracking = &Tracker{}

type effect struct {
	handler      func()
	dependencies []any
}

// TODO: We need a method for actually, after the props are created, running all the effects cause otherwise those parts of the props won't be set, an ideal thing to do additionally would be to when the effects are ran, track which signals were added in each effect and then create a mapping between signal -> effect index, that way we could only re-run the effects that matter, for others we just re-run all effects obv.
type Tracker struct {
	mu               sync.Mutex
	runMu            sync.Mutex
	context          *BuildContext
	onChangeHandlers []func()

	removal       map[any]func()
	effectsToCall map[any][]int

	currentEffect      int // -1 for no effect, index for other effect
	effects            []effect
	effectDependencies []any
}

// Implement Tracking interface
func (t *Tracker) Tracker() *Tracker {
	return t
}

func NewTracker(context *BuildContext, onChange func()) *Tracker {
	return &Tracker{
		context:          context,
		removal:          make(map[any]func()),
		effectsToCall:    make(map[any][]int),
		currentEffect:    -1,
		onChangeHandlers: []func(){onChange},
	}
}

// Clear removes all tracked signals and is safe to call multiple times.
func (t *Tracker) Clear() {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, remove := range t.removal {
		remove()
	}
	t.removal = make(map[any]func())
}

// On change will be called when anything subscribed to the tracker changes, unlike effect, it does not track dependencies, this should be used for components that share a tracker across multiple nodes.
//
// Also runs the handler once after being added.
func (t *Tracker) OnChange(t2 *Tracker, handler func()) {
	t.mu.Lock()
	t.onChangeHandlers = append(t.onChangeHandlers, func() {
		handler()
		if t != t2 {
			t2.runChangeHandlers()
		}
	})
	t.mu.Unlock()

	handler()
}

// Run all change handlers
func (t *Tracker) runChangeHandlers() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.onChangeHandlers != nil {
		for _, handler := range t.onChangeHandlers {
			handler()
		}
	}
}

// Effect adds a new handler being called when the thing changes, use this to actually react to changes of signals (runs first immediately after being called).
//
// RECURSIVE EFFECTS ARE NOT ALLOWED.
func (t *Tracker) Effect(handler func()) {
	t.mu.Lock()
	nested := t.currentEffect != -1
	t.mu.Unlock()

	if nested {
		panic("recursive effects are not allowed")
	}

	t.runMu.Lock()
	defer t.runMu.Unlock()

	t.mu.Lock()
	i := len(t.effects)
	t.effects = append(t.effects, effect{
		handler: handler,
	})
	t.mu.Unlock()

	t.runEffectLocked(i)
}

func (t *Tracker) runEffect(i int) {
	t.runMu.Lock()
	defer t.runMu.Unlock()
	t.runEffectLocked(i)
}

func (t *Tracker) runEffectLocked(i int) {
	t.mu.Lock()
	effect := &t.effects[i]

	t.currentEffect = int(i)
	t.effectDependencies = nil // Clear dependencies to start fresh for this effect
	t.mu.Unlock()

	effect.handler()

	t.mu.Lock()
	defer t.mu.Unlock()

	// Set all the new dependencies
	oldDependencies := effect.dependencies
	effect.dependencies = t.effectDependencies
	t.effectDependencies = nil

	// Remove all old dependencies that aren't there anymore from the effectsToCall list
	for _, dep := range oldDependencies {
		if t.effectsToCall[dep] != nil && !slices.ContainsFunc(effect.dependencies, func(dep2 any) bool {
			return dep == dep2
		}) {
			toCall := t.effectsToCall[dep]
			toCall = slices.DeleteFunc(toCall, func(effect int) bool {
				return effect == i
			})

			if len(toCall) == 0 {
				// When nothing is left, remove
				if rem, ok := t.removal[dep]; ok {
					rem()
				}
				delete(t.effectsToCall, dep)
				delete(t.removal, dep)
			} else {
				// Otherwise simply insert again since we're just not a dependency anymore, but other effects still are
				t.effectsToCall[dep] = toCall
			}
		}
	}

	// Add all new ones (that aren't already added), to the effectsToCall list
	for _, dep := range effect.dependencies {
		if !slices.ContainsFunc(oldDependencies, func(dep2 any) bool {
			return dep == dep2
		}) {
			toCall := t.effectsToCall[dep]
			if slices.Contains(toCall, -1) {
				continue
			}

			if toCall == nil {
				toCall = append(toCall, i)
			} else {
				// Add if not contained already
				if !slices.ContainsFunc(toCall, func(effect int) bool {
					return effect == i
				}) {
					toCall = append(toCall, i)
				}
			}

			t.effectsToCall[dep] = toCall
		}
	}

	t.currentEffect = -1
}

// Update calls all the effects on the tracker to synchronize everything
func (t *Tracker) update() {
	t.mu.Lock()
	count := len(t.effects)
	t.mu.Unlock()

	for i := range count {
		t.runEffect(i)
	}
}

// TrackValue ensures tracker is subscribed to signal and returns the current value; closed or nil inputs return zero value.
func TrackValue[T any](tracker *Tracker, signal *Signal[T]) T {
	if tracker == nil || signal == nil {
		var zero T
		return zero
	}

	tracker.mu.Lock()
	if tracker.currentEffect != -1 {
		tracker.effectDependencies = append(tracker.effectDependencies, signal)
	} else {
		// When we are outside of an effect, we want to ensure that if this signal changes, we update everything.
		tracker.effectsToCall[signal] = []int{-1}
	}

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

			tracker.mu.Lock()
			defer tracker.mu.Unlock()

			if tracker.context != nil && tracker.context.updateQueue != nil {
				if effects, ok := tracker.effectsToCall[signal]; ok {
					if slices.Contains(effects, -1) {
						tracker.context.updateQueue.Push(tracker, -1, tracker.update)
					} else {
						for _, effectID := range effects {
							tracker.context.updateQueue.Push(tracker, effectID, func() {
								tracker.runEffect(effectID)
							})
						}
					}
				} else {
					// Fallback to update everything
					tracker.context.updateQueue.Push(tracker, -1, tracker.update)
				}
			}
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
