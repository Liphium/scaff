package scaff

import (
	"time"

	"github.com/Liphium/scaff/smath"
	"github.com/hajimehoshi/ebiten/v2"
)

type Scene interface {
	// GetId returns the unique identifier of the scene.
	GetId() string

	// Update updates the scene with the given context.
	Update(c SceneContext) error

	// Draw draws the scene onto the screen.
	Draw(c SceneContext, screen *ebiten.Image)

	TransitionCapable
}

type SceneContext struct {
	Now             time.Time       // The current time.
	Focused         bool            // If a scene is in the front of the scene stack, it is focused
	Width           int             // Width of the game
	Height          int             // Height of the game
	TransitionFrame smath.Timeframe // The timeframe for the transition.
}
