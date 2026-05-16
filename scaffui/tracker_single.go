package scaffui

import "github.com/Liphium/scaff"

type SingleTracker struct {
	current Node
	tracker *scaff.Tracker
	node    *MountedNode
}

func NewSingleTracker(current Node) *SingleTracker {
	return &SingleTracker{
		current: current,
		tracker: scaff.NewTracker(),
		node:    nil,
	}
}

// Get tracker used by SingleTracker. Should be returned as tracker for current Node.
func (s *SingleTracker) Tracker() *scaff.Tracker {
	return s.tracker
}

// Set the node mounted in the tracker.
func (s *SingleTracker) SetNode(node *MountedNode) {
	if node.current == nil {
		return
	}
	s.node = node
}

// Load the node mounted in the tracker.
func (s *SingleTracker) Load(parent Node) {
	if s.node == nil {
		return
	}
	s.node.Load(parent)
}

// Get the current node stored in the tracker
func (s *SingleTracker) Node() (*MountedNode, bool) {
	return s.node, s.node != nil
}

func (s *SingleTracker) Update(parent Node, c *scaff.Context) (UpdateResult, *scaff.TracedError) {

	// If there is no node, nothing to check
	if s.node == nil {
		return NoUpdate(), nil
	}

	result, err := s.node.Update(parent, c)
	if err != nil {
		return NoUpdate(), err
	}

	// When child changed, re-layout entire box. If own size changed, mark relayout for parent.
	if result.SizeChanged {
		currentSize := s.current.Size()
		newSize, err := s.current.Layout()
		if err != nil {
			return NoUpdate(), scaff.NewTracedError(s.current, err)
		}

		if currentSize != newSize {
			return SizeChanged(), nil
		}
		return LayoutChanged(), nil
	}

	return NoUpdate(), nil
}

// Unload node mounted in tracker.
func (s *SingleTracker) Unload() {
	if s.node != nil {
		s.node.Unload()
	}

	s.tracker.Clear()
	s.tracker = nil // Cut tracker off from tree for GC
}
