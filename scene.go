package scaff

import (
	"slices"
	"time"

	"github.com/Liphium/scaff/scath"
	"github.com/hajimehoshi/ebiten/v2"
)

type Scene interface {
	// GetId returns the unique identifier of the scene.
	GetId() string

	// Update updates the scene with the given context.
	Update(c *Context) error

	// Draw draws the scene onto the screen.
	Draw(c *Context, screen *ebiten.Image)

	// HandleEvent should handle an event (like input events + resizing, etc.)
	HandleEvent(c *Context, e Event) error

	TransitionCapable
}

type Context struct {
	now             time.Time       // The current time.
	focused         bool            // If a scene is in the front of the scene stack, it is focused
	transitionFrame scath.Timeframe // The timeframe for the transition.
	width           float64         // Width of the game
	height          float64         // Height of the game

	*eventContext
}

// The current time (use for all things to make testing easier)
func (c Context) Now() time.Time {
	return c.now
}

// If a scene is in the front of the scene stack, it is focused
func (c Context) Focused() bool {
	return c.focused
}

// The timeframe for the transition
func (c Context) TransitionFrame() scath.Timeframe {
	return c.transitionFrame
}

// Width of the game
func (c Context) Width() float64 {
	return c.width
}

// Height of the game
func (c Context) Height() float64 {
	return c.height
}

type eventContext struct {
	// Any events that have already been handled.
	//
	// This is for telling other nodes that, an event with some id has already been handled and stuff (useful for clicks and such).
	events []EventId
}

// Check if any kind of event has already been handled.
func (c *eventContext) IsHandled(event EventId) bool {
	return slices.Contains(c.events, event)
}

// Mark a type of event as handled.
func (c *eventContext) Handled(event EventId) {
	c.events = append(c.events, event)
}
