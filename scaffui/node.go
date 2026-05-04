package scaffui

import (
	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/smath"
)

// Layout algorithm basic idea (copy what Flutter is doing kinda)
//
// 1. Constraints pass: Constraints go down
// 2. Size pass: Sizes go up

type Node interface {
	Tracking

	// Should return the name or something of the node (doesn't need to be unique, is used in the error path especially)
	ID() string

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
	Update(c *scaff.LayerContext) (sizeChange bool, err *Error)

	// Draw the thing onto the screen at a specified position (next step is getting this to work)
	Draw(position smath.Vec, renderer Renderer)

	// Handle events from cgui (you do not have to handle any, but should always push them along to children at least)
	HandleEvent(c *scaff.LayerContext, event Event) *Error
}

type Tracking interface {
	// Should return the tracker for the current Node
	Tracker() *Tracker
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
