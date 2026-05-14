package scaff

import (
	"slices"
	"time"

	"github.com/Liphium/scaff/smath"
	"github.com/hajimehoshi/ebiten/v2"
)

type Node interface {
	Tracking
	Identifiable

	// Should return your own parent
	Parent() Node

	// Should return your own children
	Children() []Node

	// Called when the Node is initialized in the tree
	Load(parent Node)

	// Called when the Node is removed from the tree
	Unload()

	// Called on every physics tick (like 60 times a second, depending on what ebitens tick rate is)
	Update(c *Context) *TracedError

	// Handle events from the system (you do not have to handle any, but should always push them along to children at least)
	HandleEvent(c *Context, event Event) *TracedError

	// Draw the thing onto the screen
	Draw(c *Context, image *ebiten.Image)
}

type Context struct {
	Now             time.Time       // The current time.
	TransitionFrame smath.Timeframe // The timeframe for the transition.
	Width           int             // Width of the game
	Height          int             // Height of the game

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

// An interface for props of a node that has children. This helps nodes like single child not rebuild as often.
type ChildProps interface {
	GetBuilders() []NodeBuilder
}
