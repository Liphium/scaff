package scaff

type NodeBuilder func() Node

type SingleTracker struct {
	parent  Node
	current Node
	builder NodeBuilder
	tracker *Tracker
}

// Create a new single tracker. If the builder is nil, it can never change from that state again (cause that wouldn't make sense anyway).
func NewSingleTracker(parent Node, builder NodeBuilder) *SingleTracker {
	if parent == nil {
		log.Error("parent of node can not be nil for single tracker")
		return nil
	}

	// If there is no builder, just return an empty single tracker, this will basically be doing nothing, but that's fine since no child
	if builder == nil {
		return &SingleTracker{
			current: nil,
			tracker: NewTracker(),
		}
	}

	// Build the current scene
	current := builder()
	current.Load(parent)

	return &SingleTracker{
		parent:  parent,
		current: current,
		builder: builder,
		tracker: NewTracker(),
	}
}

// Get tracker used by SingleTracker. Should be returned as tracker for current Node.
func (s *SingleTracker) Tracker() *Tracker {
	return s.tracker
}

// Load the node mounted in the tracker.
func (s *SingleTracker) Load() {
	if s.current == nil {
		return
	}
	s.current.Load(s.parent)
}

func (s *SingleTracker) Update(c *Context) *TracedError {

	// If there is no node, nothing to check
	if s.current == nil {
		return nil
	}

	// If dirty, rebuild
	if s.tracker.Changed() {
		s.current.Unload()
		s.current = s.builder()
		s.current.Load(s.parent)
	}

	// Forward the update to the child
	return s.current.Update(c)
}

// Unload node mounted in tracker.
func (s *SingleTracker) Unload() {
	if s.current != nil {
		s.current.Unload()
	}

	s.tracker.Clear()
	s.tracker = nil // Cut tracker off from tree for GC
}
