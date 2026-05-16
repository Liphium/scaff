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

	TransitionCapable
}

type Context struct {
	Now             time.Time       // The current time.
	Focused         bool            // If a scene is in the front of the scene stack, it is focused
	TransitionFrame scath.Timeframe // The timeframe for the transition.
	Width           float64         // Width of the game
	Height          float64         // Height of the game

	// Any events that have already been handled.
	//
	// This is for telling other nodes that, an event with some id has already been handled and stuff (useful for clicks and such).
	events []EventId
}

// Check if any kind of event has already been handled.
func (c *Context) IsHandled(event EventId) bool {
	return slices.Contains(c.events, event)
}

// Mark a type of event as handled.
func (c *Context) Handled(event EventId) {
	c.events = append(c.events, event)
}
