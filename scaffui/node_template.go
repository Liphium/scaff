package scaffui

import (
	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/optional"
)

// Make sure we actually implement the things
var _ scaff.ChildProps[NodeBuilder] = AcceptNoChild{}
var _ scaff.ChildProps[NodeBuilder] = AcceptChild{}
var _ scaff.ChildProps[NodeBuilder] = AcceptChildren{}

// A struct that can be embedded as a pointer to implement a Child(builder NodeBuilder) function on props for a node that should have no child.
//
// This is recommended for use with a single child node that has no child as it saves a lot of code.
type AcceptNoChild struct{}

func (anc AcceptNoChild) GetBuilders() []NodeBuilder {
	return nil
}

// A struct that can be embedded as a pointer to implement a Child(builder NodeBuilder) function on props for a node with a single child.
//
// This is recommended for use with a single child node as it saves a lot of code.
type AcceptChild struct {
	child optional.O[NodeBuilder]
}

func (ac *AcceptChild) Child(builder NodeBuilder) {
	ac.child.SetValue(builder)
}

func (ac AcceptChild) GetBuilders() []NodeBuilder {
	if child, ok := ac.child.Value(); ok {
		return []NodeBuilder{child}
	}
	return nil
}

// A struct that can be embedded as a pointer to implement a Child(builder NodeBuilder) function on props for a node with a multiple children.
//
// This is recommended for use with a multi child node as it saves a lot of code.
type AcceptChildren struct {
	children []NodeBuilder
}

func (ac *AcceptChildren) Child(builder NodeBuilder) {
	ac.children = append(ac.children, builder)
}

func (ac AcceptChildren) GetBuilders() []NodeBuilder {
	return ac.children
}
