package scaff

import (
	"errors"
	"time"

	"github.com/Liphium/scaff/optional"
	"github.com/Liphium/scaff/smath"
)

type TransitionCapable interface {
	// Transition is called when the state is transitioned in or out to configure the transition properly.
	Transition(in bool) TransitionProperties

	// Load is called when the state is loaded.
	Load()

	// Unload is called when the state is unloaded.
	Unload()
}

type TransitionProperties struct {
	// If the transition is isolated, the next / previous state will be transitioned in / out after.
	//
	// If the duration is not the same, the transition currently coming in always wins since it is drawn on top and the other state will just be dropped.
	//
	// Isolation will always be prioritied over non-isolation.
	Isolated bool

	// How long the transition should take. 0 = instant
	Duration time.Duration
}

// ToTimeframe converts the transition properties to a smath.Timeframe.
func (tp TransitionProperties) ToTimeframe(now time.Time, outTransition optional.O[TransitionProperties]) smath.Timeframe {
	tf := smath.NewTimeframe(now, tp.Duration)
	if val, ok := outTransition.Value(); tp.Isolated && ok && val.Isolated {
		tf = tf.AddDelay(val.Duration)
	}
	return tf
}

func NoTransition() TransitionProperties {
	return TransitionProperties{
		Isolated: false,
		Duration: 0,
	}
}

// TransitioningState is a struct storing one object that can transition.
//
// NOT GOROUTINE SAFE
type TransitioningState[T TransitionCapable] struct {
	from      optional.O[T]
	fromFrame optional.O[smath.Timeframe]
	to        optional.O[T]
	toFrame   optional.O[smath.Timeframe]
}

// NewTransitioningState creates a new TransitioningState with the given state.
//
// The in transition will be initialized to start immediately.
func NewTransitioningState[T TransitionCapable](now time.Time, state T) *TransitioningState[T] {
	ts := &TransitioningState[T]{
		from:      optional.With(state),
		fromFrame: optional.With(state.Transition(true).ToTimeframe(now, optional.None[TransitionProperties]())),
	}
	if val, ok := ts.from.Value(); ok {
		val.Load()
	}
	return ts
}

// Set the state in the TransitioningState. You can also use nil to transition the thing out.
//
// This will not handle setting the same thing twice, please deal with that at a higher level since comparisons in here are expensive.
func (ts *TransitioningState[T]) Set(now time.Time, state optional.O[T]) {

	// If there already is a transition, just skip and set the thing straight away
	if to, ok := ts.to.Value(); ok {
		if from, ok := ts.from.Value(); ok {
			from.Unload()
		}
		to.Unload()
		ts.to = optional.None[T]()
		ts.toFrame = optional.None[smath.Timeframe]()

		if s, ok := state.Value(); ok {
			s.Load()
		}
		ts.from = state
		ts.fromFrame = optional.With(NoTransition().ToTimeframe(now, optional.None[TransitionProperties]()))
		return
	}

	// 1. Let from transition out in case
	outTime := optional.None[TransitionProperties]()
	if from, ok := ts.from.Value(); ok {
		props := from.Transition(false)
		outTime.SetValue(props)
		ts.fromFrame = optional.With(props.ToTimeframe(now, optional.None[TransitionProperties]()).MakeBackwards())
	}

	// 2. Set the new state
	ts.to = state
	if s, ok := state.Value(); ok {
		s.Load()
		ts.toFrame = optional.With(s.Transition(true).ToTimeframe(now, outTime))
	} else {
		ts.toFrame = optional.With(NoTransition().ToTimeframe(now, outTime))
	}
}

func (ts *TransitioningState[T]) SetEmpty(now time.Time) {
	ts.Set(now, optional.None[T]())
}

// Update should be called to both render and update the state of the transitions.
func (ts *TransitioningState[T]) Update(now time.Time, update func(T, smath.Timeframe) error) error {

	// If both transitions are over or not there, switch "to" to "from"
	toFrame, toOk := ts.toFrame.Value()
	fromFrame, fromOk := ts.fromFrame.Value()
	toEnded := !toOk || toFrame.Over(now)
	fromEnded := !fromOk || fromFrame.Over(now)
	if toEnded && fromEnded && toOk {
		if from, ok := ts.from.Value(); ok {
			from.Unload()
		}
		ts.from = ts.to
		ts.fromFrame = ts.toFrame
		ts.to = optional.None[T]()
		ts.toFrame = optional.None[smath.Timeframe]()
	}

	// Render "from" if there
	if from, ok := ts.from.Value(); ok {
		fromFrame, ok := ts.fromFrame.Value()
		if !ok {
			return errors.New("from frame not set")
		}
		if err := update(from, fromFrame); err != nil {
			return err
		}
	}

	// Render "to" if there
	if to, ok := ts.to.Value(); ok {
		toFrame, ok := ts.toFrame.Value()
		if !ok {
			return errors.New("to frame not set")
		}
		if err := update(to, toFrame); err != nil {
			return err
		}
	}
	return nil
}

// GetCurrent returns the current state, immediately returns new state, even when the transition is not over.
func (ts *TransitioningState[T]) GetCurrent() optional.O[T] {
	if to, ok := ts.to.Value(); ok {
		return optional.With(to)
	}
	return ts.from
}
