package scaff

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Type check to make sure the scene interface is properly implemented
var _ Scene = &Root{}

type Root struct {
	ID              string                          // ID of the scene (MUST BE SET)
	TransitionProps func(bool) TransitionProperties // Set the transition properties for this scene

	sceneRoot   Node
	prevCursorX int
	prevCursorY int
}

// Pass the ID to the GetID function for the Scene interface
func (r *Root) GetId() string {
	return r.ID
}

// Transfers Load from the Scene interface to the scene root.
func (r *Root) Load() {
	if r.sceneRoot == nil {
		return
	}

	// Scene root has no parent, so nil for parent
	r.sceneRoot.Load(nil)
}

// Transfers Unload from the Scene interface to the scene root.
func (r *Root) Unload() {
	if r.sceneRoot == nil {
		return
	}
	r.sceneRoot.Unload()
}

// Forward transition properties from the embedded function
func (r *Root) Transition(in bool) TransitionProperties {
	if r.TransitionProps == nil {
		return NoTransition()
	}
	return r.TransitionProps(in)
}

// Forward update + handle some events that are only available in update
func (r *Root) Update(c SceneContext) error {
	if r.sceneRoot == nil {
		return nil
	}

	ctx := r.buildContext(c)
	x, y := ebiten.CursorPosition()

	// For mouse events, it's important that release is checked before pressed since the user may have released the button and pressed it again between the frame so both events could be emitted. In such a case, the release event should always be first.
	buttons := []ebiten.MouseButton{
		ebiten.MouseButton0,
		ebiten.MouseButton1,
		ebiten.MouseButton2,
		ebiten.MouseButton3,
		ebiten.MouseButton4,
	}

	// Check all of the common mouse buttons for release events
	for _, button := range buttons {
		if inpututil.IsMouseButtonJustReleased(button) {
			if err := r.sceneRoot.HandleEvent(ctx, ReleaseEvent{
				X:      x,
				Y:      y,
				Button: button,
			}); err != nil {
				return err
			}
		}
	}

	// Check all of the common mouse buttons for press events
	for _, button := range buttons {
		if inpututil.IsMouseButtonJustPressed(button) {
			if err := r.sceneRoot.HandleEvent(ctx, ReleaseEvent{
				X:      x,
				Y:      y,
				Button: button,
			}); err != nil {
				return err
			}
		}
	}

	// Check for any delta in scroll
	scrollX, scrollY := ebiten.Wheel()
	if scrollX != 0 || scrollY != 0 {
		r.sceneRoot.HandleEvent(ctx, ScrollEvent{
			X:       x,
			Y:       y,
			ScrollX: scrollX,
			ScrollY: scrollY,
		})
	}

	// Let the actual scene root update itself
	return r.sceneRoot.Update(ctx)
}

// Pass Draw to the layers in the correct order
func (r *Root) Draw(c SceneContext, screen *ebiten.Image) {
	if r.sceneRoot == nil {
		return
	}

	ctx := r.buildContext(c)
	x, y := ebiten.CursorPosition()

	deltaX := x - r.prevCursorX
	deltaY := y - r.prevCursorY
	r.prevCursorX = x
	r.prevCursorY = y

	if deltaX != 0 || deltaY != 0 {
		// Move events are emitted every frame to make sure dragging is smooth
		r.sceneRoot.HandleEvent(r.buildContext(c), MoveEvent{
			X:      x,
			Y:      y,
			DeltaX: deltaX,
			DeltaY: deltaY,
		})
	}

	r.sceneRoot.Draw(ctx, screen)
}

func (r *Root) buildContext(c SceneContext) *Context {
	return &Context{
		Now:             c.Now,
		TransitionFrame: c.TransitionFrame,
		Width:           c.Width,
		Height:          c.Height,
		events:          []EventId{},
	}
}
