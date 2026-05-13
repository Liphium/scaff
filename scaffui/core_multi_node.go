package scaffui

import (
	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/smath"
)

func UseMultiNode[P any](id string, propsCreator func(t *scaff.Tracker, props *P), create func(core *MultiChildConstruct[P])) NodeBuilder {
	return func() Node {
		node := newMultiChildCoreNode[P]()
		node.id = id

		// Create the props
		var props P
		propsCreator(node.Tracker(), &props)
		node.props = props

		// Build the actual node
		create(&MultiChildConstruct[P]{node: node})

		return node
	}
}

type MultiChildConstruct[P any] struct {
	node *MultiChildNode[P]
}

func (m *MultiChildConstruct[P]) Tracker() *scaff.Tracker {
	return m.node.Tracker()
}

func (m *MultiChildConstruct[P]) Props() P {
	return m.node.props
}

func (m *MultiChildConstruct[P]) Child(builder NodeBuilder) {
	m.node.tracker.Add(NewMountedFromBuilder(builder))
}

func (m *MultiChildConstruct[P]) WantedConstraints(fn func(node *MultiChildNode[P], parent Constraints) Constraints) {
	m.node.onWantedConstraints = fn
}

func (m *MultiChildConstruct[P]) Load(fn func(node *MultiChildNode[P])) {
	m.node.onLoad = fn
}

func (m *MultiChildConstruct[P]) Unload(fn func(node *MultiChildNode[P])) {
	m.node.onUnload = fn
}

func (m *MultiChildConstruct[P]) Layout(fn func(node *MultiChildNode[P]) (Size, error)) {
	m.node.onLayout = fn
}

func (m *MultiChildConstruct[P]) HandleEvent(fn func(node *MultiChildNode[P], c *scaff.LayerContext, event Event) error) {
	m.node.onHandleEvent = fn
}

func (m *MultiChildConstruct[P]) Update(fn func(node *MultiChildNode[P], c *scaff.LayerContext) (bool, error)) {
	m.node.onUpdate = fn
}

func (m *MultiChildConstruct[P]) Draw(fn func(node *MultiChildNode[P], position smath.Vec, renderer Renderer)) {
	m.node.onDraw = fn
}

func newMultiChildCoreNode[P any]() *MultiChildNode[P] {
	n := &MultiChildNode[P]{}
	n.tracker = NewMultiTracker(n)
	return n
}

var _ Node = &MultiChildNode[any]{}
var _ WantsConstraints = &MultiChildNode[any]{}

type MultiChildNode[P any] struct {
	tracker     *MultiTracker
	size        Size
	constraints Constraints

	id                  string
	props               P
	onLoad              func(node *MultiChildNode[P])
	onUnload            func(node *MultiChildNode[P])
	onWantedConstraints func(node *MultiChildNode[P], parent Constraints) Constraints
	onLayout            func(node *MultiChildNode[P]) (Size, error)
	onHandleEvent       func(node *MultiChildNode[P], c *scaff.LayerContext, event Event) error
	onUpdate            func(node *MultiChildNode[P], c *scaff.LayerContext) (bool, error)
	onDraw              func(node *MultiChildNode[P], position smath.Vec, renderer Renderer)
}

func (m *MultiChildNode[P]) ID() string {
	return m.id
}

func (m *MultiChildNode[P]) Props() P {
	return m.props
}

func (m *MultiChildNode[P]) Load() {
	if m.onLoad != nil {
		m.onLoad(m)
	}

	m.tracker.Load()
}

func (m *MultiChildNode[P]) Size() Size {
	return m.size
}

func (m *MultiChildNode[P]) Constraints() Constraints {
	return m.constraints
}

func (m *MultiChildNode[P]) SetConstraints(c Constraints) {
	m.constraints = c
}

func (m *MultiChildNode[P]) WantedConstraints(parent Constraints) Constraints {
	if m.onWantedConstraints == nil {
		return Unconstrained()
	}

	return m.onWantedConstraints(m, parent)
}

func (m *MultiChildNode[P]) Layout() (Size, error) {
	if m.onLayout != nil {
		size, err := m.onLayout(m)
		if err != nil {
			return Size{}, scaff.NewTracedError(m, err)
		}
		m.size = size
		return size, nil
	}

	size := Size{Width: m.constraints.MinWidth, Height: m.constraints.MinHeight}
	for _, child := range m.tracker.Nodes() {
		child.Current().SetConstraints(m.constraints)
		childSize, err := child.Current().Layout()
		if err != nil {
			return Size{}, scaff.NewTracedError(m, err)
		}

		size.Width = max(size.Width, childSize.Width)
		size.Height = max(size.Height, childSize.Height)
	}

	m.size = size
	return size, nil
}

func (m *MultiChildNode[P]) HandleEvent(c *scaff.LayerContext, event Event) *scaff.TracedError {
	if m.onHandleEvent != nil {
		if err := m.onHandleEvent(m, c, event); err != nil {
			return scaff.NewTracedError(m, err)
		}
	}

	for _, child := range m.tracker.Nodes() {
		if err := child.Current().HandleEvent(c, event); err != nil {
			return scaff.NewTracedError(m, err)
		}
	}

	return nil
}

func (m *MultiChildNode[P]) Tracker() *scaff.Tracker {
	return m.tracker.Tracker()
}

func (m *MultiChildNode[P]) Update(c *scaff.LayerContext) (bool, *scaff.TracedError) {
	if m.onUpdate != nil {
		relayout, err := m.onUpdate(m, c)
		if err != nil {
			return false, scaff.NewTracedError(m, err)
		}

		childRelayout, updateErr := m.tracker.Update(c)
		if updateErr != nil {
			return false, scaff.NewTracedError(m, updateErr)
		}

		return relayout || childRelayout, nil
	}

	return m.tracker.Update(c)
}

func (m *MultiChildNode[P]) Unload() {
	if m.onUnload != nil {
		m.onUnload(m)
	}

	if m.tracker != nil {
		m.tracker.Unload()
		m.tracker = nil
	}
}

func (m *MultiChildNode[P]) Draw(position smath.Vec, renderer Renderer) {
	if m.onDraw != nil {
		m.onDraw(m, position, renderer)
		return
	}

	for _, child := range m.tracker.Nodes() {
		child.Current().Draw(position, renderer)
	}
}

func (m *MultiChildNode[P]) Children() []*MountedNode {
	if m.tracker == nil {
		return nil
	}

	return m.tracker.Nodes()
}
