package scaffui

import "github.com/Liphium/scaff"

func NewMountedFromBuilder(builder NodeBuilder) *MountedNode {
	return &MountedNode{
		current:   builder(),
		construct: builder,
	}
}

type NodeBuilder func() Node

type MountedNode struct {
	current   Node
	construct NodeBuilder
}

func (w *MountedNode) Current() Node {
	return w.current
}

func (w *MountedNode) Load(parent Node) {
	w.current.Load(parent)
}

func (w *MountedNode) Unload() {
	w.current.Unload()
	w.construct = nil
	w.current = nil
}

// Should be called for an update from the parent, the boolean indicates whether a re-layout should be done (forwards errors from the update of the child)
func (w *MountedNode) Update(parent Node, c *scaff.Context) (UpdateResult, *scaff.TracedError) {
	result, err := w.current.Update(c)
	if err != nil {
		return result, err
	}

	if w.current.Tracker() != nil && w.current.Tracker().Changed() {
		w.current.Unload()
		w.current = nil // For GC to absolutely know that this Node is no longer needed
		w.Rebuild()
		w.current.Load(parent)

		// Technically our size didn't change, but we need our parent to check that for us
		return SizeChanged(), nil
	}
	return result, nil
}

func (w *MountedNode) Rebuild() {
	w.current = w.construct()
}
