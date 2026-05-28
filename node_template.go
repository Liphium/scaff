package scaff

import (
	"github.com/Liphium/scaff/optional"
)

type AcceptNoChild = AcceptNoChildTemplate[NodeBuilder]
type AcceptChild = AcceptChildTemplate[NodeBuilder]
type AcceptChildren = AcceptChildrenTemplate[NodeBuilder]

// Make sure we actually implement the things
var _ ChildProps[NodeBuilder] = AcceptNoChildTemplate[NodeBuilder]{}
var _ ChildProps[NodeBuilder] = AcceptChildTemplate[NodeBuilder]{}
var _ ChildProps[NodeBuilder] = AcceptChildrenTemplate[NodeBuilder]{}

// A struct that can be embedded as a pointer to implement a Child(builder NodeBuilder) function on props for a node that should have no child.
//
// This is recommended for use with a single child node that has no child as it saves a lot of code.
type AcceptNoChildTemplate[B any] struct{}

func (anc AcceptNoChildTemplate[B]) GetBuilders() []B {
	return nil
}

// A struct that can be embedded as a pointer to implement a Child(builder NodeBuilder) function on props for a node with a single child.
//
// This is recommended for use with a single child node as it saves a lot of code.
type AcceptChildTemplate[B any] struct {
	child optional.O[B]
}

func (ac *AcceptChildTemplate[B]) Child(builder B) {
	ac.child.SetValue(builder)
}

func (ac *AcceptChildTemplate[B]) GetChild() optional.O[B] {
	return ac.child
}

func (ac AcceptChildTemplate[B]) GetBuilders() []B {
	if child, ok := ac.child.Value(); ok {
		return []B{child}
	}
	return nil
}

// A struct that can be embedded as a pointer to implement a Child(builder NodeBuilder) function on props for a node with a multiple children.
//
// This is recommended for use with a multi child node as it saves a lot of code.
type AcceptChildrenTemplate[B any] struct {
	children []B
}

func (ac *AcceptChildrenTemplate[B]) Child(builder B) {
	ac.children = append(ac.children, builder)
}

func (ac *AcceptChildrenTemplate[B]) GetChildren() []B {
	return ac.children
}

func (ac AcceptChildrenTemplate[B]) GetBuilders() []B {
	return ac.children
}
