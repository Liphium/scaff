package scaffui

import (
	"github.com/Liphium/scaff"
)

type Node interface {
	scaff.Node

	// Should return the current size of the node
	Size() Size

	// Should return the current constraints (that were last set)
	Constraints() Constraints

	// Set the constraints of the node (should be used for layouting)
	SetConstraints(c Constraints)

	// Should layout the node and return the size within the (previously set) constraints
	Layout() (Size, error)
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
