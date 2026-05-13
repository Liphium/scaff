package scaffui

import "github.com/Liphium/scaff"

type MultiTracker struct {
	current Node
	tracker *scaff.Tracker
	nodes   []*MountedNode
}

func NewMultiTracker(current Node) *MultiTracker {
	return &MultiTracker{
		current: current,
		tracker: scaff.NewTracker(),
	}
}

// Get the tracker used by the MultiTracker. Should be returned as the tracker for the current Node.
func (mt *MultiTracker) Tracker() *scaff.Tracker {
	return mt.tracker
}

// Add a new node to the MultiTracker
func (mt *MultiTracker) Add(node *MountedNode) {
	if node == nil {
		return
	}

	mt.nodes = append(mt.nodes, node)
}

func (mt *MultiTracker) Nodes() []*MountedNode {
	return mt.nodes
}

// Load all nodes mounted in the tracker.
func (m *MultiTracker) Load() {
	for _, child := range m.nodes {
		child.Load()
	}
}

func (m *MultiTracker) Update(c *scaff.LayerContext) (relayout bool, err *scaff.TracedError) {
	changed := false
	for _, node := range m.nodes {
		change, err := node.Update(c)
		if err != nil {
			return false, err
		}

		// If there was a change, set changed to true (or ignore if previously true)
		changed = changed || change
	}

	// When one of the children changed, re-layout the entire box, if own size changed, mark as dirty
	if changed {
		currentSize := m.current.Size()
		newSize, err := m.current.Layout()
		if err != nil {
			return false, scaff.NewTracedError(m.current, err)
		}

		// When size changed, indicate relayout to parent
		if currentSize != newSize {
			return true, nil
		}
	}

	return false, nil
}

// Unload all nodes mounted in the tracker.
func (m *MultiTracker) Unload() {
	for _, child := range m.nodes {
		child.Unload()
	}

	m.tracker.Clear()
	m.tracker = nil // Cut tracker off from tree for GC
}
