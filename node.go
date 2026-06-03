package scaff

import (
	"github.com/Liphium/scaff/paint"
	"github.com/hajimehoshi/ebiten/v2"
)

// TODO: Only rebuild props of nodes when their tracker is dirty: Check in the node builder, if not own thing dirty, just return old thingy (nothing changed anyway).

// All the context for nodes being created within ScaffUI
type BuildContext struct {
	assetManager *paint.AssetManager
	updateQueue  *UpdateQueue
}

func (bc BuildContext) UpdateQueue() *UpdateQueue {
	return bc.updateQueue
}

func (bc BuildContext) AssetManager() *paint.AssetManager {
	return bc.assetManager
}

func (bc BuildContext) CopyWithNewUpdateQueue() *BuildContext {
	return &BuildContext{
		updateQueue:  NewUpdateQueue(),
		assetManager: bc.assetManager,
	}
}

type NodeBuilder func(*BuildContext) Node

type Loadable[T any] interface {

	// Called when the Node is initialized in the tree
	Load(parent T)

	// Called when the Node is removed from the tree
	Unload()
}

type Node interface {
	Tracking
	Identifiable
	Loadable[Node]

	// Should return your own parent
	Parent() Node

	// Should return your own children
	Children() []Node

	// Called on every physics tick (like 60 times a second, depending on what ebitens tick rate is)
	Update(c *Context) TracedError

	// Handle events from the system (you do not have to handle any, but should always push them along to children at least)
	HandleEvent(c *Context, event Event) TracedError

	// Draw the thing onto the screen
	Draw(c *Context, image *ebiten.Image)
}

// An interface for props of a node that has children. This helps nodes like single child not rebuild as often.
type ChildProps[B any] interface {
	// GetBuilders gets all the builders for the children of this node
	GetBuilders() []B

	// GetChanged gets all the changed indices since the last time ClearChanged was called.
	GetChanged() []uint

	// ClearChanged clears all the changed indices
	ClearChanged()
}
