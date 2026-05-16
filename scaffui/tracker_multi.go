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
func (m *MultiTracker) Load(parent Node) {
	for _, child := range m.nodes {
		child.Load(parent)
	}
}

func (m *MultiTracker) Update(parent Node, c *scaff.Context) (UpdateResult, *scaff.TracedError) {
	result := NoUpdate()
	for _, node := range m.nodes {
		change, err := node.Update(parent, c)
		if err != nil {
			return NoUpdate(), err
		}

		// Stack the updates on top of each other (will mark changed in case was not before)
		result.Stack(change)
	}

	// When one of the children changed, re-layout the entire box, if own size changed, mark as dirty
	if result.SizeChanged {
		currentSize := m.current.Size()
		newSize, err := m.current.Layout()
		if err != nil {
			return NoUpdate(), scaff.NewTracedError(m.current, err)
		}

		// When size changed, indicate relayout to parent
		if currentSize != newSize {
			return SizeChanged(), nil
		}
		return LayoutChanged(), nil
	}

	return result, nil
}

// Unload all nodes mounted in the tracker.
func (m *MultiTracker) Unload() {
	for _, child := range m.nodes {
		child.Unload()
	}

	m.tracker.Clear()
	m.tracker = nil // Cut tracker off from tree for GC
}
