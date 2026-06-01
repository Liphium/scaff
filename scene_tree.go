package scaff

import (
	"github.com/Liphium/scaff/paint"
	"github.com/hajimehoshi/ebiten/v2"
)

// NewSceneTree makes a new scene tree based on an identifier. This tree supports a builder pattern, just use .Mount(<node>) to mount a node into it or add a transition with SetTransitionProps.
func NewSceneTree(id string, assetManager *paint.AssetManager) *SceneTree {
	return &SceneTree{
		id:           id,
		assetManager: assetManager,
	}
}

// Type check to make sure the scene interface is properly implemented
var _ Scene = &SceneTree{}

type SceneTree struct {
	id              string                          // ID of the scene (MUST BE SET)
	assetManager    *paint.AssetManager             // Asset manager for the scene
	transitionProps func(bool) TransitionProperties // Set the transition properties for this scene
	sceneRoot       Node                            // Node at the root of the scene tree
}

// Mount a node as the root of the scene tree.
func (st *SceneTree) Mount(create func(t *Tracker, props *RootProps)) *SceneTree {

	// Create a single child node that essentially just exists to refresh the builder passed in
	node := &SingleChildNode[AcceptNoChild]{
		id:      "root",
		tracker: NewTracker(),
		context: &BuildContext{
			assetManager: st.assetManager,
		},
		singleProps: &SingleChildProps[AcceptNoChild]{},
	}
	node.builder = root(create) // Node will automatically be built on load

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
func (st *SceneTree) Update(c *Context) error {
	if st.sceneRoot == nil {
		return nil
	}

	// Let the actual scene root update itself
	return st.sceneRoot.Update(c)
}

// Pass events to the root node
func (st *SceneTree) HandleEvent(c *Context, e Event) error {
	return st.sceneRoot.HandleEvent(c, e)
}

// Pass Draw to the layers in the correct order
func (st *SceneTree) Draw(c *Context, screen *ebiten.Image) {
	if st.sceneRoot == nil {
		return
	}

	st.sceneRoot.Draw(c, screen)
}
