package scaffui

import (
	"github.com/Liphium/scaff/engine"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/scath"
)

type NodeBuilder func(*scaff.BuildContext) Node

type UpdateResult struct {
	SizeChanged     bool // The size of the current node changed due to updates, re-layouting is needed
	AnythingChanged bool // If any of the children changed, but no more re-layouting is needed
}

// Stack a UpdateResult with a different one, will keep everything true that was before and sort of combine every boolean using or.
func (ur *UpdateResult) Stack(ur2 UpdateResult) {
	ur.AnythingChanged = ur.AnythingChanged || ur2.AnythingChanged
	ur.SizeChanged = ur.SizeChanged || ur2.SizeChanged
}

// When no update happened, return this UpdateResult
func NoUpdate() UpdateResult {
	return UpdateResult{
		SizeChanged:     false,
		AnythingChanged: false,
	}
}

// When your size changed, return this to signal a change in size
func SizeChanged() UpdateResult {
	return UpdateResult{
		SizeChanged:     true,
		AnythingChanged: true,
	}
}

// Wen you did re-layouting, but your size didn't change
func LayoutChanged() UpdateResult {
	return UpdateResult{
		SizeChanged:     false,
		AnythingChanged: true,
	}
}

type Node interface {
	scaff.Tracking
	scaff.Identifiable
	scaff.Loadable[Node]

	// Get the parent node of the current Node
	Parent() Node

	// Get all children nodes of the current Node
	Children() []Node

	// Should return the current size of the node
	Size() scath.Vec

	// Should return the constraints wanted by this Node
	WantedConstraints(parent scath.Constraints) scath.Constraints

	// Should return the current constraints (that were last set)
	Constraints() scath.Constraints

	// Set the constraints of the node (should be used for layouting)
	SetConstraints(c scath.Constraints)

	// Should layout the node and return the size within the (previously set) constraints
	Layout() (scath.Vec, error)

	// Called to see if a draw is needed based on changed props, etc.
	Update() (UpdateResult, scaff.TracedError)

	// Draw the thing onto the screen at a specified position (next step is getting this to work)
	Draw(position scath.Vec, painter engine.Painter)

	// Handle events from scaff (you do not have to handle any, but should always push them along to children at least)
	HandleEvent(c *scaff.Context, event scaff.Event) scaff.TracedError
}
