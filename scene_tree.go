package scaff

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// NewSceneTree makes a new scene tree based on an identifier. This tree supports a builder pattern, just use .Mount(<node>) to mount a node into it or add a transition with SetTransitionProps.
func NewSceneTree(id string) *SceneTree {
	return &SceneTree{
		id: id,
	}
}

// Type check to make sure the scene interface is properly implemented
var _ Scene = &SceneTree{}

type SceneTree struct {
	id              string                          // ID of the scene (MUST BE SET)
	transitionProps func(bool) TransitionProperties // Set the transition properties for this scene
	sceneRoot       Node                            // Node at the root of the scene tree

	prevCursorX int
	prevCursorY int
}

// Mount a node as the root of the scene tree.
func (st *SceneTree) Mount(builder NodeBuilder) *SceneTree {

	// Create a single child node that essentially just exists to refresh the builder passed in
	node := &SingleChildNode[int8]{
		id:      "root",
		tracker: NewTracker(),
	}
	node.builder = builder // Node will automatically be built on load

	// Mount the node inside of a node that can refresh
	st.sceneRoot = node
	return st
}

// Set what kind of transition properties the scene should use (the parameter in the function is if the transition is reverse or not)
func (st *SceneTree) SetTransitionProps(handler func(bool) TransitionProperties) *SceneTree {
	st.transitionProps = handler
	return st
}

// Pass the ID to the GetID function for the Scene interface
func (st *SceneTree) GetId() string {
	return st.id
}

// Transfers Load from the Scene interface to the scene root.
func (st *SceneTree) Load() {
	if st.sceneRoot == nil {
		return
	}

	// Scene root has no parent, so nil for parent
	st.sceneRoot.Load(nil)
}

// Transfers Unload from the Scene interface to the scene root.
func (st *SceneTree) Unload() {
	if st.sceneRoot == nil {
		return
	}
	st.sceneRoot.Unload()
}

// Forward transition properties from the embedded function
func (st *SceneTree) Transition(in bool) TransitionProperties {
	if st.transitionProps == nil {
		return NoTransition()
	}
	return st.transitionProps(in)
}

// Forward update + handle some events that are only available in update
func (st *SceneTree) Update(c SceneContext) error {
	if st.sceneRoot == nil {
		return nil
	}

	ctx := st.buildContext(c)
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
			if err := st.sceneRoot.HandleEvent(ctx, ReleaseEvent{
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
			if err := st.sceneRoot.HandleEvent(ctx, ReleaseEvent{
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
		st.sceneRoot.HandleEvent(ctx, ScrollEvent{
			X:       x,
			Y:       y,
			ScrollX: scrollX,
			ScrollY: scrollY,
		})
	}

	// Let the actual scene root update itself
	return st.sceneRoot.Update(ctx)
}

// Pass Draw to the layers in the correct order
func (st *SceneTree) Draw(c SceneContext, screen *ebiten.Image) {
	if st.sceneRoot == nil {
		return
	}

	ctx := st.buildContext(c)
	x, y := ebiten.CursorPosition()

	deltaX := x - st.prevCursorX
	deltaY := y - st.prevCursorY
	st.prevCursorX = x
	st.prevCursorY = y

	if deltaX != 0 || deltaY != 0 {
		// Move events are emitted every frame to make sure dragging is smooth
		st.sceneRoot.HandleEvent(st.buildContext(c), MoveEvent{
			X:      x,
			Y:      y,
			DeltaX: deltaX,
			DeltaY: deltaY,
		})
	}

	st.sceneRoot.Draw(ctx, screen)
}

func (st *SceneTree) buildContext(c SceneContext) *Context {
	return &Context{
		Now:             c.Now,
		TransitionFrame: c.TransitionFrame,
		Width:           c.Width,
		Height:          c.Height,
		events:          []EventId{},
	}
}
