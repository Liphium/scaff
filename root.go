package scaff

import (
	"slices"
	"time"

	"github.com/Liphium/scaff/smath"
	"github.com/hajimehoshi/ebiten/v2"
)

type EventId string

// Type check to make sure the scene interface is properly implemented
var _ Scene = &Root{}

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
func (lc *Context) IsHandled(event EventId) bool {
	return slices.Contains(lc.events, event)
}

// Mark a type of event as handled.
func (lc *Context) Handled(event EventId) {
	lc.events = append(lc.events, event)
}

type Root struct {
	ID              string                          // ID of the scene (MUST BE SET)
	TransitionProps func(bool) TransitionProperties // Set the transition properties for this scene
}

// Pass the ID to the GetID function for the Scene interface
func (ls *Root) GetId() string {
	return ls.ID
}

// Transfers Load from the Scene interface to the layers.
func (ls *Root) Load() {

}

// Transfers Unload from the Scene interface to the layers.
func (ls *Root) Unload() {
}

// Forward transition properties from the embedded function
func (ls *Root) Transition(in bool) TransitionProperties {
	if ls.TransitionProps == nil {
		return NoTransition()
	}
	return ls.TransitionProps(in)
}

func (ls *Root) buildContext(c SceneContext) *Context {
	return &Context{
		Now:             c.Now,
		TransitionFrame: c.TransitionFrame,
		Width:           c.Width,
		Height:          c.Height,
		events:          []EventId{},
	}
}

// Forward update + do own logic for scene layouting and stuff
func (ls *Root) Update(c SceneContext) error {
	ls.buildContext(c)
	return nil
}

// Pass Draw to the layers in the correct order
func (ls *Root) Draw(c SceneContext, screen *ebiten.Image) {
	ls.buildContext(c)
}
