package scaff

import (
	"github.com/Liphium/scaff/optional"
)

// A struct that can be embedded as a pointer to implement a Child(builder NodeBuilder) function on props for a node with no children.
//
// This is recommended for use with a single child node as it saves a lot of code and implements all of the things required for children properly according to best practices.
type AcceptNoChild = AcceptNoChildTemplate[NodeBuilder]

// A struct that can be embedded as a pointer to implement a Child(builder NodeBuilder) function on props for a node with a single child.
//
// This is recommended for use with a single child node as it saves a lot of code and implements all of the things required for children properly according to best practices.
type AcceptChild = AcceptChildTemplate[NodeBuilder]

// A struct that can be embedded as a pointer to implement a Child(builder NodeBuilder) function on props for a node with a multiple children.
//
// This is recommended for use with a multi child node as it saves a lot of code and implements all of the things required for children properly according to best practices.
type AcceptChildren = AcceptChildrenTemplate[NodeBuilder]

// Make sure we actually implement the things
var _ ChildProps[NodeBuilder] = &AcceptNoChildTemplate[NodeBuilder]{}
var _ ChildProps[NodeBuilder] = &AcceptChildTemplate[NodeBuilder]{}
var _ ChildProps[NodeBuilder] = &AcceptChildrenTemplate[NodeBuilder]{}

// Template struct for type aliases that can be embedded as a pointer to implement a Child(builder NodeBuilder) function on props for a node that should have no child.
type AcceptNoChildTemplate[B any] struct{}

func (anc AcceptNoChildTemplate[B]) GetBuilders() []B {
	return nil
}

func (anc AcceptNoChildTemplate[B]) GetChanged() []uint {
	return nil
}

func (anc AcceptNoChildTemplate[B]) ClearChanged() {}

// Template struct for type aliases that can be embedded as a pointer to implement a Child(builder NodeBuilder) function on props for a node with a single child.
type AcceptChildTemplate[B any] struct {
	changed bool
	child   optional.O[B]
}

// Child adds a child to the node. There can only be one child.
//
// There are internal reasons for this to be a method. Trackers just need to work with this properly.
func (ac *AcceptChildTemplate[B]) Child(builder B) {
	ac.child.SetValue(builder)
	ac.changed = true
}

func (ac AcceptChildTemplate[B]) GetBuilders() []B {
	if child, ok := ac.child.Value(); ok {
		return []B{child}
	}
	return nil
}

func (ac AcceptChildTemplate[B]) GetChanged() []uint {
	if ac.changed {
		return []uint{0}
	}
	return nil
}

func (ac *AcceptChildTemplate[B]) ClearChanged() {
	ac.changed = false
}

// Template struct for type aliases that can be embedded as a pointer to implement a Child(i uint, builder NodeBuilder) function on props for a node with a multiple children.
//
// There are internal reasons for this to be a method. Trackers just need to work with this properly.
type AcceptChildrenTemplate[B any] struct {
	changed  []uint
	children []B
}

// Child adds a new child to the node. i is the index of the child. Make sure it is actually valid and stuff.
//
// There are internal reasons for this to be a method. Trackers just need to work with this properly.
func (ac *AcceptChildrenTemplate[B]) Child(i uint, builder B) {
	ac.changed = append(ac.changed, i)
	ac.children = append(ac.children, builder)
}

func (ac AcceptChildrenTemplate[B]) GetBuilders() []B {
	return ac.children
}

func (ac AcceptChildrenTemplate[B]) GetChanged() []uint {
	return ac.changed
}

func (ac *AcceptChildrenTemplate[B]) ClearChanged() {
	ac.changed = nil
}
