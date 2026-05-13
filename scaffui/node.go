package scaffui

import (
	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/smath"
)

type Node interface {
	scaff.Tracking
	scaff.Identifiable

	// Should return the current size of the node
	Size() Size

	// Should return the current constraints (that were last set)
	Constraints() Constraints

	// Set the constraints of the node (should be used for layouting)
	SetConstraints(c Constraints)

	// Should layout the node and return the size within the (previously set) constraints
	Layout() (Size, error)

	// Called when the Node is initialized in the tree
	Load()

	// Called when the Node is removed from the tree
	Unload()

	// Called on every tick, use to handle state updates, etc.
	Update(c *scaff.LayerContext) (sizeChange bool, err *scaff.TracedError)

	// Draw the thing onto the screen at a specified position (next step is getting this to work)
	Draw(position smath.Vec, renderer Renderer)

	// Handle events from cgui (you do not have to handle any, but should always push them along to children at least)
	HandleEvent(c *scaff.LayerContext, event Event) *scaff.TracedError
}

type WantsConstraints interface {
	// Should return the constraints wanted by this Node
	WantedConstraints(parent Constraints) Constraints
}

// Get the constraints (default is unconstrained)
func WantedConstraints(node Node, parent Constraints) Constraints {
	wanted := Unconstrained()
	if constrained, ok := node.(WantsConstraints); ok {
		wanted = constrained.WantedConstraints(parent)
	}
	return wanted
}
