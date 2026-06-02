package scaffui

import (
	"fmt"

	"github.com/Liphium/scaff/paint"

	"github.com/Liphium/scaff"

	"github.com/Liphium/scaff/scath"
)

// Struct for defining a standard node.
//
// ID should be a unique id for the node, but also probably be readable as it shows up in error messages.
//
// DefaultProps are the props as they are by default. Make sure to also specify any embedded struct pointers.
//
// PropsCreator should be the function passed in by users of your node (as in it should probably be an argument of the function creating your node).
//
// Create is the function actually specifying your node. You can overwrite all of the functions of the node interface there, with some exceptions that we implement for you.
type StandardCreate[P scaff.ChildProps[NodeBuilder]] struct {
	ID           string
	DefaultProps P
	PropsCreator func(t *scaff.Tracker, props *P)
	Create       func(props *StandardMethods[P])
}

// Standard lets you create a node with multiple children. Simply implement the ChildProps interface on the props you want to have for your node.
func Standard[P scaff.ChildProps[NodeBuilder]](create StandardCreate[P]) NodeBuilder {
	node := &StandardNode[P]{
		id:         create.ID,
		multiProps: &StandardMethods[P]{},
	}
	if create.Create != nil {
		create.Create(node.multiProps)
	}

	return func(context *scaff.BuildContext) Node {
		node.tracker = scaff.NewTracker(context, node.PropsChanged)
		node.context = context

		// Create the props (if desired)
		props := create.DefaultProps
		if create.PropsCreator != nil {
			create.PropsCreator(node.Tracker(), &props)
		}
		node.props = props

		return node
	}
}

type StandardMethods[P scaff.ChildProps[NodeBuilder]] struct {
	OnLoad              func(node *StandardNode[P])
	OnPropsChanged      func(node *StandardNode[P])
	OnUnload            func(node *StandardNode[P])
	OnWantedConstraints func(node *StandardNode[P], parent scath.Constraints) scath.Constraints
	OnLayout            func(node *StandardNode[P]) (scath.Vec, error)
	OnHandleEvent       func(node *StandardNode[P], c *scaff.Context, event scaff.Event) error
	OnDraw              func(node *StandardNode[P], position scath.Vec, renderer paint.Painter)
}

var _ Node = &StandardNode[*AcceptChildren]{}

type StandardNode[P scaff.ChildProps[NodeBuilder]] struct {
	id       string
	parent   Node
	children []Node

	context     *scaff.BuildContext
	tracker     *scaff.Tracker
	size        scath.Vec
	constraints scath.Constraints

	props      P
	dirty      bool
	multiProps *StandardMethods[P]
}

func (s *StandardNode[P]) ID() string {
	return s.id
}

func (s *StandardNode[P]) Parent() Node {
	return s.parent
}

func (s *StandardNode[P]) Children() []Node {
	return s.children
}

func (s *StandardNode[P]) Props() P {
	return s.props
}

func (s *StandardNode[P]) Context() *scaff.BuildContext {
	return s.context
}

func (s *StandardNode[P]) Load(parent Node) {

	// Set parent + build the children
	s.parent = parent
	s.PropsChanged()

	if s.multiProps.OnLoad != nil {
		s.multiProps.OnLoad(s)
	}
}

func (s *StandardNode[P]) PropsChanged() {
	s.dirty = true

	changed := s.props.GetChanged()
	if changed == nil {
		return
	}

	builders := s.props.GetBuilders()
	for _, i := range changed {
		if len(s.children) <= int(i) {
			log.Error("index out of bounds for props update", "i", i, "children", len(s.children))
			continue
		}

		s.children[i].Unload()
		s.children[i] = builders[i](s.context)
		s.children[i].Load(s)
	}
}

func (s *StandardNode[P]) Size() scath.Vec {
	return s.size
}

func (s *StandardNode[P]) Constraints() scath.Constraints {
	return s.constraints
}

func (s *StandardNode[P]) SetConstraints(c scath.Constraints) {
	s.constraints = c
}

func (s *StandardNode[P]) WantedConstraints(parent scath.Constraints) scath.Constraints {
	if s.multiProps.OnWantedConstraints == nil {
		return scath.Unconstrained()
	}

	return s.multiProps.OnWantedConstraints(s, parent)
}

func (s *StandardNode[P]) Layout() (scath.Vec, error) {

	// If layouting is handled by a method, use that
	if s.multiProps.OnLayout != nil {
		size, err := s.multiProps.OnLayout(s)
		if err != nil {
			return scath.Vec{}, scaff.NewTracedError(s, err)
		}

		finalSize, err := s.finalSize(size)
		s.size = finalSize
		return finalSize, err
	}

	// Use minimum size or biggest size of children
	size := scath.Vec{X: s.constraints.MinX, Y: s.constraints.MinY}
	for _, child := range s.children {
		child.SetConstraints(s.constraints)
		childSize, err := child.Layout()
		if err != nil {
			return scath.Vec{}, scaff.NewTracedError(s, err)
		}

		size.X = max(size.X, childSize.X)
		size.Y = max(size.Y, childSize.Y)
	}

	finalSize, err := s.finalSize(size)
	s.size = finalSize
	return finalSize, err
}

// TODO: Fix thight constraint handling everywhere to make sure we always check that the size fits within the parent + that when we get get a min width or height but our child is smaller, we still take the minimum height / width
func (s *StandardNode[P]) finalSize(size scath.Vec) (scath.Vec, error) {

	// Still size bigger when the parent gives us larger constraints with minimums
	if size.X < s.constraints.MinX {
		size.X = s.constraints.MinX
	}
	if size.Y < s.constraints.MinY {
		size.Y = s.constraints.MinY
	}

	// Make sure the size is actually correct
	if !size.FitsWithin(s.constraints) {
		return scath.Vec{}, fmt.Errorf("child size is too big for parent constraints w=%v h=%v parent=%v", size.X, size.Y, s.constraints)
	}

	return size, nil
}

func (s *StandardNode[P]) HandleEvent(c *scaff.Context, event scaff.Event) scaff.TracedError {
	if s.multiProps.OnHandleEvent != nil {
		if err := s.multiProps.OnHandleEvent(s, c, event); err != nil {
			return scaff.NewTracedError(s, err)
		}
	}

	for _, child := range s.children {
		if err := child.HandleEvent(c, event); err != nil {
			return scaff.NewTracedError(s, err)
		}
	}

	return nil
}

func (s *StandardNode[P]) Tracker() *scaff.Tracker {
	return s.tracker.Tracker()
}

func (s *StandardNode[P]) Update() (UpdateResult, scaff.TracedError) {
	// If we're dirty, instantly say that the size changed (relayout is needed)
	if s.dirty {
		s.dirty = false
		return SizeChanged(), nil
	}

	// Otherwise, go through all children and stack the result
	result := NoUpdate()
	for _, child := range s.children {
		change, err := child.Update()
		if err != nil {
			return NoUpdate(), scaff.NewTracedError(s, err)
		}

		// Stack the updates on top of each other (will mark changed in case was not before)
		result.Stack(change)
	}

	// When one of the children changed, re-layout the entire box, if own size changed, tell the parent to also re-layout
	if result.SizeChanged {
		currentSize := s.Size()
		newSize, err := s.Layout()
		if err != nil {
			return NoUpdate(), scaff.NewTracedError(s, err)
		}

		// When size changed, indicate relayout to parent
		if currentSize != newSize {
			return SizeChanged(), nil
		}
		return LayoutChanged(), nil
	}
	return result, nil
}

func (s *StandardNode[P]) Unload() {
	if s.multiProps.OnUnload != nil {
		s.multiProps.OnUnload(s)
	}

	if s.tracker != nil {
		s.tracker.Clear()
		s.tracker = nil
	}
}

func (s *StandardNode[P]) Draw(position scath.Vec, renderer paint.Painter) {
	if s.multiProps.OnDraw != nil {
		s.multiProps.OnDraw(s, position, renderer)
		return
	}

	for _, child := range s.children {
		child.Draw(position, renderer)
	}
}
