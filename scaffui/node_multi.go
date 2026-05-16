package scaffui

import (
	"github.com/Liphium/scaff/paint"

	"github.com/Liphium/scaff"

	"github.com/Liphium/scaff/scath"
)

// CreateMultiNode lets you create a node with multiple children. Simply implement the ChildProps interface on the props you want to have for your node.
//
// id should be a unique id for the node, but also probably be readable as it shows up in error messages.
//
// propsCreator should be the function passed in by users of your node (as in it should probably be an argument of the function creating your node).
//
// create is the function actually specifying your node. You can overwrite all of the functions of the node interface there, with some exceptions that we implement for you.
func CreateMultiNode[P scaff.ChildProps](id string, propsCreator func(t *scaff.Tracker, props *P), create func(core *MultiChildProps[P])) NodeBuilder {
	node := &MultiChildNode[P]{
		id:         id,
		multiProps: &MultiChildProps[P]{},
	}
	create(node.multiProps)

	return func() Node {
		node.tracker = NewMultiTracker(node)

		// Create the props
		var props P
		propsCreator(node.Tracker(), &props)
		node.props = props

		return node
	}
}

type MultiChildProps[P any] struct {
	onLoad              func(node *MultiChildNode[P])
	onUnload            func(node *MultiChildNode[P])
	onWantedConstraints func(node *MultiChildNode[P], parent Constraints) Constraints
	onLayout            func(node *MultiChildNode[P]) (scath.Vec, error)
	onHandleEvent       func(node *MultiChildNode[P], c *scaff.Context, event scaff.Event) error
	onUpdate            func(node *MultiChildNode[P], c *scaff.Context) (bool, error)
	onDraw              func(node *MultiChildNode[P], position scath.Vec, renderer paint.Painter)
}

func (m *MultiChildProps[P]) WantedConstraints(fn func(node *MultiChildNode[P], parent Constraints) Constraints) {
	m.onWantedConstraints = fn
}

func (m *MultiChildProps[P]) Load(fn func(node *MultiChildNode[P])) {
	m.onLoad = fn
}

func (m *MultiChildProps[P]) Unload(fn func(node *MultiChildNode[P])) {
	m.onUnload = fn
}

func (m *MultiChildProps[P]) Layout(fn func(node *MultiChildNode[P]) (scath.Vec, error)) {
	m.onLayout = fn
}

func (m *MultiChildProps[P]) HandleEvent(fn func(node *MultiChildNode[P], c *scaff.Context, event scaff.Event) error) {
	m.onHandleEvent = fn
}

func (m *MultiChildProps[P]) Update(fn func(node *MultiChildNode[P], c *scaff.Context) (bool, error)) {
	m.onUpdate = fn
}

func (m *MultiChildProps[P]) Draw(fn func(node *MultiChildNode[P], position scath.Vec, renderer paint.Painter)) {
	m.onDraw = fn
}

var _ Node = &MultiChildNode[any]{}
var _ WantsConstraints = &MultiChildNode[any]{}

type MultiChildNode[P any] struct {
	tracker     *MultiTracker
	size        scath.Vec
	constraints Constraints

	id         string
	props      P
	multiProps *MultiChildProps[P]
}

func (m *MultiChildNode[P]) ID() string {
	return m.id
}

func (m *MultiChildNode[P]) Props() P {
	return m.props
}

func (m *MultiChildNode[P]) Load(parent Node) {
	if m.multiProps.onLoad != nil {
		m.multiProps.onLoad(m)
	}

	m.tracker.Load(parent)
}

func (m *MultiChildNode[P]) Size() scath.Vec {
	return m.size
}

func (m *MultiChildNode[P]) Constraints() Constraints {
	return m.constraints
}

func (m *MultiChildNode[P]) SetConstraints(c Constraints) {
	m.constraints = c
}

func (m *MultiChildNode[P]) WantedConstraints(parent Constraints) Constraints {
	if m.multiProps.onWantedConstraints == nil {
		return Unconstrained()
	}

	return m.multiProps.onWantedConstraints(m, parent)
}

func (m *MultiChildNode[P]) Layout() (scath.Vec, error) {
	if m.multiProps.onLayout != nil {
		size, err := m.multiProps.onLayout(m)
		if err != nil {
			return scath.Vec{}, scaff.NewTracedError(m, err)
		}
		m.size = size
		return size, nil
	}

	size := scath.Vec{X: m.constraints.MinX, Y: m.constraints.MinY}
	for _, child := range m.tracker.Nodes() {
		child.Current().SetConstraints(m.constraints)
		childSize, err := child.Current().Layout()
		if err != nil {
			return scath.Vec{}, scaff.NewTracedError(m, err)
		}

		size.X = max(size.X, childSize.X)
		size.Y = max(size.Y, childSize.Y)
	}

	m.size = size
	return size, nil
}

func (m *MultiChildNode[P]) HandleEvent(c *scaff.Context, event scaff.Event) *scaff.TracedError {
	if m.multiProps.onHandleEvent != nil {
		if err := m.multiProps.onHandleEvent(m, c, event); err != nil {
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

func (m *MultiChildNode[P]) Update(c *scaff.Context) (UpdateResult, *scaff.TracedError) {
	if m.multiProps.onUpdate != nil {
		changed, err := m.multiProps.onUpdate(m, c)
		if err != nil {
			return NoUpdate(), scaff.NewTracedError(m, err)
		}

		result, updateErr := m.tracker.Update(m, c)
		if updateErr != nil {
			return NoUpdate(), scaff.NewTracedError(m, updateErr)
		}

		// If the parent has a relayout to be done, it's stronger than the children because they don't know about it
		if changed {
			result.Stack(SizeChanged())
		}

		return result, nil
	}

	return m.tracker.Update(m, c)
}

func (m *MultiChildNode[P]) Unload() {
	if m.multiProps.onUnload != nil {
		m.multiProps.onUnload(m)
	}

	if m.tracker != nil {
		m.tracker.Unload()
		m.tracker = nil
	}
}

func (m *MultiChildNode[P]) Draw(position scath.Vec, renderer paint.Painter) {
	if m.multiProps.onDraw != nil {
		m.multiProps.onDraw(m, position, renderer)
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
