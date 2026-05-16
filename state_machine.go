package scaff

import (
	"time"

	"github.com/Liphium/scaff/optional"
	"github.com/Liphium/scaff/scath"
)

// Type check to make sure State confirms to TransitionCapable
var _ TransitionCapable = &State[any, any]{}

type State[O any, S any] struct {
	id      S
	checker func(obj O) bool // To check if the state is the current one

	// Passed by the parent
	transition TransitionProperties
}

func (s *State[O, S]) Transition(in bool) TransitionProperties {
	return s.transition
}

// These functions are just required to implement the interface, they don't really matter
func (s *State[O, S]) Load()   {}
func (s *State[O, S]) Unload() {}

// Create a new state. Will be activated if checker returned true and no state upper in the ladder has been activated.
func NewState[O any, S any](state S, checker func(O) bool) *State[O, S] {
	return &State[O, S]{
		id:         state,
		checker:    checker,
		transition: NoTransition(),
	}
}

// ONLY USE THIS FOR THE DEFAULT STATE. Sort of like the else branch of an if statement, this state is active if no other state's checker returns true.
func DefaultState[O any, S any](state S) *State[O, S] {
	return &State[O, S]{
		id:         state,
		checker:    nil,
		transition: NoTransition(),
	}
}

type StateMachine[O any, S comparable] struct {
	transition TransitionProperties
	states     []*State[O, S]
	def        *State[O, S]

	// Internal state
	state *TransitioningState[*State[O, S]]
}

// A struct to create a new StateMachine.
type StateMachineCreate[O any, S comparable] struct {
	// What transition to use between states. If not set, no transition will be used.
	Transition optional.O[TransitionProperties]

	// All the states in StateMachine
	States []*State[O, S]

	// The default state (at the start, should not be in States)
	Default *State[O, S]
}

// Create a new StateMachine.
//
// Some things to keep in mind:
// - The first states are the highest priority (just like in a if statement).
// - def is the default state meaning when no other state works for the conditions, it will be used.
func NewStateMachine[O any, S comparable](create StateMachineCreate[O, S]) *StateMachine[O, S] {

	// Update the transition in all states
	for _, state := range create.States {
		state.transition = create.Transition.Or(NoTransition())
	}
	create.Default.transition = create.Transition.Or(NoTransition())

	return &StateMachine[O, S]{
		transition: create.Transition.Or(NoTransition()),
		states:     create.States,
		def:        create.Default,
	}
}

// Call to update the state in the StateMachine.
func (sm *StateMachine[O, S]) Update(now time.Time, obj O) {

	// Find the current state (iterate backwards cause the first ones are supposed to have the highest priority)
	current := sm.def
	for _, state := range sm.states {
		if state.checker(obj) {
			current = state
			break
		}
	}

	// Create transitioning state if not there yet
	if sm.state == nil {
		sm.state = NewTransitioningState(now, current)
	}

	// Change to the new state
	if curr, ok := sm.state.GetCurrent().Value(); ok && curr.id != current.id {
		sm.state.Set(now, optional.With(current))
	} else if !ok {
		// Set anyway, even though this should not happen
		sm.state.Set(now, optional.With(current))
	}

	// Update the transitioning state properly (to make sure it switches its internal logic)
	sm.state.Update(now, func(s *State[O, S], t scath.Timeframe) error {
		return nil
	})
}

// Use this to draw the actual object that is supposed to be drawn from the state in the StateMachine.
func (sm *StateMachine[O, S]) Draw(now time.Time, drawFunc func(state S, frame scath.Timeframe)) {
	sm.state.Update(now, func(s *State[O, S], frame scath.Timeframe) error {
		drawFunc(s.id, frame)
		return nil
	})
}
