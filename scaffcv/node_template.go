package scaffcv

import "github.com/Liphium/scaff"

// Wrap templates for use in scaffcv

type AcceptNoChild = scaff.AcceptNoChildTemplate[NodeBuilder]
type AcceptChild = scaff.AcceptChildTemplate[NodeBuilder]
type AcceptChildren = scaff.AcceptChildrenTemplate[NodeBuilder]

// Wrap helper functions for use in scaffcv

func EmptyNoChild() AcceptNoChild {
	return scaff.EmptyNoChildTemplate[NodeBuilder]()
}

func EmptyChild() *AcceptChild {
	return scaff.EmptyChildTemplate[NodeBuilder]()
}

func WithChild(builder NodeBuilder) *AcceptChild {
	return scaff.WithChildTemplate(builder)
}

func EmptyChildren() *AcceptChildren {
	return scaff.EmptyChildrenTemplate[NodeBuilder]()
}

func WithChildren(builders []NodeBuilder) *AcceptChildren {
	return scaff.WithChildrenTemplate(builders)
}
