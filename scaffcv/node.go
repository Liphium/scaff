package scaffcv

import (
	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/engine"
	"github.com/Liphium/scaff/scath"
)

type BuildContext struct {
	cam *Camera
	*scaff.BuildContext
}

func (bc BuildContext) Camera() *Camera {
	return bc.cam
}

type NodeBuilder func(context *BuildContext) Node

type Node interface {
	scaff.Tracking
	scaff.Identifiable
	scaff.Loadable[Node]

	// Should return your current position in world coordinates (this can be used to determine if you should be rendered or not)
	Position() scath.Vec

	// Should return your current size in the world (this can be used to determine if you should be rendered or not)
	Size() scath.Vec

	// Should return your own parent
	Parent() Node

	// Should return your own children
	Children() []Node

	// Called on every physics tick (like 60 times a second, depending on what ebitens tick rate is)
	Update(c *scaff.Context) scaff.TracedError

	// Handle events from the system (you do not have to handle any, but should always push them along to children at least)
	HandleEvent(c *scaff.Context, event scaff.Event) scaff.TracedError

	// Draw the thing onto the screen (world coordinates)
	Draw(c *scaff.Context, painter engine.Painter)
}
