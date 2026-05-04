package scaff

import (
	"slices"
	"time"

	"github.com/Liphium/scaff/smath"
	"github.com/hajimehoshi/ebiten/v2"
)

type EventId string

const (
	// Any handling of mouse input
	EventMouseInput EventId = "mouse-input"
)

// Type check to make sure the scene interface is properly implemented
var _ Scene = &LayeredScene{}

type WorldLayer interface {
	Update(c *LayerContext) error
	Draw(c *LayerContext, camera *Camera, screen *ebiten.Image)
	Load()
	Unload()
}

type UILayer interface {
	Update(c *LayerContext) error
	Draw(c *LayerContext, screen *ebiten.Image)
	Load()
	Unload()
}

type LayerContext struct {
	Now             time.Time       // The current time.
	TransitionFrame smath.Timeframe // The timeframe for the transition.
	Width           int             // Width of the game
	Height          int             // Height of the game

	// Any events that have already been handled.
	//
	// This is for telling bottom layers that mouse input has already been handled and stuff.
	events []EventId
}

// Check if any kind of event has already been handled.
func (lc *LayerContext) IsHandled(event EventId) bool {
	return slices.Contains(lc.events, event)
}

// Mark a type of event as handled.
func (lc *LayerContext) Handled(event EventId) {
	lc.events = append(lc.events, event)
}

type LayeredScene struct {
	ID              string                          // ID of the scene (MUST BE SET)
	WorldLayers     []WorldLayer                    // All world layers in the scene
	UILayers        []UILayer                       // All UI layers in the scene
	StartX, StartY  float64                         // Start coordinates for the camera
	TransitionProps func(bool) TransitionProperties // Set the transition properties for this scene

	cam *Camera
}

// Pass the ID to the GetID function for the Scene interface
func (ls *LayeredScene) GetId() string {
	return ls.ID
}

// Transfers Load from the Scene interface to the layers.
func (ls *LayeredScene) Load() {
	ls.cam = NewCamera(ls.StartX, ls.StartY, 100, 100)

	for _, layer := range ls.WorldLayers {
		layer.Load()
	}
	for _, layer := range ls.UILayers {
		layer.Load()
	}
}

// Transfers Unload from the Scene interface to the layers.
func (ls *LayeredScene) Unload() {
	for _, layer := range ls.WorldLayers {
		layer.Unload()
	}
	for _, layer := range ls.UILayers {
		layer.Unload()
	}
}

// Forward transition properties from the embedded function
func (ls *LayeredScene) Transition(in bool) TransitionProperties {
	if ls.TransitionProps == nil {
		return NoTransition()
	}
	return ls.TransitionProps(in)
}

func (ls *LayeredScene) buildContext(c SceneContext) *LayerContext {
	return &LayerContext{
		Now:             c.Now,
		TransitionFrame: c.TransitionFrame,
		Width:           c.Width,
		Height:          c.Height,
		events:          []EventId{},
	}
}

// Forward update + do own logic for scene layouting and stuff
func (ls *LayeredScene) Update(c SceneContext) error {
	ctx := ls.buildContext(c)

	// Handle UI layers first (are above world layers), do it backwards to have the top layer first
	for _, layer := range slices.Backward(ls.UILayers) {
		if err := layer.Update(ctx); err != nil {
			return err
		}
	}

	// Handle world layers second (backwards for top layers first)
	for _, layer := range slices.Backward(ls.WorldLayers) {
		if err := layer.Update(ctx); err != nil {
			return err
		}
	}
	return nil
}

// Pass Draw to the layers in the correct order
func (ls *LayeredScene) Draw(c SceneContext, screen *ebiten.Image) {
	ctx := ls.buildContext(c)

	// Update size for the camera (might change at any time)
	ls.cam.SetSize(float64(ctx.Width), float64(ctx.Height))

	// Handle world layers first (below UI layers)
	for _, layer := range ls.WorldLayers {
		layer.Draw(ctx, ls.cam, screen)
	}

	// Draw UI layers first (are above world layers)
	for _, layer := range ls.UILayers {
		layer.Draw(ctx, screen)
	}
}
