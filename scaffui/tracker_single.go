package scaffui

import "github.com/Liphium/scaff"

type SingleTracker struct {
	current Node
	tracker *Tracker
	node    *MountedNode
}

func NewSingleTracker(current Node) *SingleTracker {
	return &SingleTracker{
		current: current,
		tracker: NewTracker(),
		node:    nil,
	}
}

// Get tracker used by SingleTracker. Should be returned as tracker for current Node.
func (s *SingleTracker) Tracker() *Tracker {
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
func (s *SingleTracker) Load() {
	if s.node == nil {
		return
	}
	s.node.Load()
}

// Get the current node stored in the tracker
func (s *SingleTracker) Node() (*MountedNode, bool) {
	return s.node, s.node != nil
}

func (s *SingleTracker) Update(c *scaff.LayerContext) (relayout bool, err *Error) {

	// If there is no node, nothing to check
	if s.node == nil {
		return false, nil
	}

	changed, err := s.node.Update(c)
	if err != nil {
		return false, err
	}

	// When child changed, re-layout entire box. If own size changed, mark relayout for parent.
	if changed {
		currentSize := s.current.Size()
		newSize, err := s.current.Layout()
		if err != nil {
			return false, NewError(s.current, err)
		}

		if currentSize != newSize {
			return true, nil
		}
	}

	return false, nil
}

// Unload node mounted in tracker.
func (s *SingleTracker) Unload() {
	if s.node != nil {
		s.node.Unload()
	}

	s.tracker.Clear()
	s.tracker = nil // Cut tracker off from tree for GC
}
