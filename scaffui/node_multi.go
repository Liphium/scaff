package scaffui

import (
	"github.com/Liphium/scaff/paint"

	"github.com/Liphium/scaff"

	"github.com/Liphium/scaff/scath"
)

// Struct for defining a new multi child node.
//
// ID should be a unique id for the node, but also probably be readable as it shows up in error messages.
//
// DefaultProps are the props as they are by default. Make sure to also specify any embedded struct pointers.
//
// PropsCreator should be the function passed in by users of your node (as in it should probably be an argument of the function creating your node).
//
// Create is the function actually specifying your node. You can overwrite all of the functions of the node interface there, with some exceptions that we implement for you.
type MultiNodeCreate[P scaff.ChildProps[NodeBuilder]] struct {
	ID           string
	DefaultProps P
	PropsCreator func(t *scaff.Tracker, props *P)
	Create       func(props *MultiChildProps[P])
}

// MultiNode lets you create a node with multiple children. Simply implement the ChildProps interface on the props you want to have for your node.
func MultiNode[P scaff.ChildProps[NodeBuilder]](create MultiNodeCreate[P]) NodeBuilder {
	node := &MultiChildNode[P]{
		id:         create.ID,
		multiProps: &MultiChildProps[P]{},
	}
	if create.Create != nil {
		create.Create(node.multiProps)
	}

	return func(context *scaff.BuildContext) Node {
		node.tracker = NewMultiTracker(node)
		node.context = context

		// Create the props (if desired)
		props := create.DefaultProps
		if create.PropsCreator != nil {
			create.PropsCreator(node.Tracker(), &props)

			// Add children (if desired)
			if builders := props.GetBuilders(); len(builders) > 0 {
				for _, builder := range builders {
					node.tracker.Add(NewMountedFromBuilder(builder, context))
				}
			}
		}
		node.props = props

		// Call props changed hook when defined
		if node.multiProps.OnPropsChanged != nil {
			node.multiProps.OnPropsChanged(node)
		}

		return node
	}
}

type MultiChildProps[P any] struct {
	OnLoad              func(node *MultiChildNode[P])
	OnPropsChanged      func(node *MultiChildNode[P])
	OnUnload            func(node *MultiChildNode[P])
	OnWantedConstraints func(node *MultiChildNode[P], parent scath.Constraints) scath.Constraints
	OnLayout            func(node *MultiChildNode[P]) (scath.Vec, error)
	OnHandleEvent       func(node *MultiChildNode[P], c *scaff.Context, event scaff.Event) error
	OnUpdate            func(node *MultiChildNode[P], c *scaff.Context) (bool, error)
	OnDraw              func(node *MultiChildNode[P], position scath.Vec, renderer paint.Painter)
}

var _ Node = &MultiChildNode[any]{}
var _ WantsConstraints = &MultiChildNode[any]{}

type MultiChildNode[P any] struct {
	context     *scaff.BuildContext
	tracker     *MultiTracker
	size        scath.Vec
	constraints scath.Constraints

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

func (s *MultiChildNode[P]) Context() *scaff.BuildContext {
	return s.context
}

func (m *MultiChildNode[P]) Load(parent Node) {
	if m.multiProps.OnLoad != nil {
		m.multiProps.OnLoad(m)
	}

	m.tracker.Load(parent)
}

func (m *MultiChildNode[P]) Size() scath.Vec {
	return m.size
}

func (m *MultiChildNode[P]) Constraints() scath.Constraints {
	return m.constraints
}

func (m *MultiChildNode[P]) SetConstraints(c scath.Constraints) {
	m.constraints = c
}

func (m *MultiChildNode[P]) WantedConstraints(parent scath.Constraints) scath.Constraints {
	if m.multiProps.OnWantedConstraints == nil {
		return scath.Unconstrained()
	}

	return m.multiProps.OnWantedConstraints(m, parent)
}

func (m *MultiChildNode[P]) Layout() (scath.Vec, error) {
	if m.multiProps.OnLayout != nil {
		size, err := m.multiProps.OnLayout(m)
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

func (m *MultiChildNode[P]) HandleEvent(c *scaff.Context, event scaff.Event) scaff.TracedError {
	if m.multiProps.OnHandleEvent != nil {
		if err := m.multiProps.OnHandleEvent(m, c, event); err != nil {
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

func (m *MultiChildNode[P]) Update(c *scaff.Context) (UpdateResult, scaff.TracedError) {
	if m.multiProps.OnUpdate != nil {
		changed, err := m.multiProps.OnUpdate(m, c)
		if err != nil {
			return NoUpdate(), scaff.NewTracedError(m, err)
		}

		result, updateErr := m.tracker.Update(m, c, m.context)
		if updateErr != nil {
			return NoUpdate(), scaff.NewTracedError(m, updateErr)
		}

		// If the parent has a relayout to be done, it's stronger than the children because they don't know about it
		if changed {
			result.Stack(SizeChanged())
		}

		return result, nil
	}

	return m.tracker.Update(m, c, m.context)
}

func (m *MultiChildNode[P]) Unload() {
	if m.multiProps.OnUnload != nil {
		m.multiProps.OnUnload(m)
	}

	if m.tracker != nil {
		m.tracker.Unload()
		m.tracker = nil
	}
}

func (m *MultiChildNode[P]) Draw(position scath.Vec, renderer paint.Painter) {
	if m.multiProps.OnDraw != nil {
		m.multiProps.OnDraw(m, position, renderer)
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
