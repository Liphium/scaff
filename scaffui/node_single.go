package scaffui

import (
	"github.com/Liphium/scaff/paint"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/scath"
)

// CreateSingleNode lets you create a node with a single child. Simply implement the ChildProps interface on the props you want to have for your node.
//
// id should be a unique id for the node, but also probably be readable as it shows up in error messages.
//
// propsCreator should be the function passed in by users of your node (as in it should probably be an argument of the function creating your node).
//
// create is the function actually specifying your node. You can overwrite all of the functions of the node interface there, with some exceptions that we implement for you.
func CreateSingleNode[P scaff.ChildProps](id string, propsCreator func(t *scaff.Tracker, props *P), create func(core *SingleChildProps[P])) NodeBuilder {
	node := &SingleChildNode[P]{
		id:          id,
		singleProps: &SingleChildProps[P]{},
	}
	create(node.singleProps)

	return func() Node {
		node.tracker = NewSingleTracker(node)

		// Fill the props
		var props P
		propsCreator(node.Tracker(), &props)
		node.props = props

		return node
	}
}

type SingleChildProps[P any] struct {
	onLoad              func(node *SingleChildNode[P])
	onUnload            func(node *SingleChildNode[P])
	onWantedConstraints func(node *SingleChildNode[P], parent Constraints) Constraints
	onLayout            func(node *SingleChildNode[P]) (scath.Vec, error)
	onHandleEvent       func(node *SingleChildNode[P], c *scaff.Context, event scaff.Event) error
	onUpdate            func(node *SingleChildNode[P], c *scaff.Context) (bool, error)
	onDraw              func(node *SingleChildNode[P], position scath.Vec, renderer paint.Painter)
}

func (s *SingleChildProps[P]) WantedConstraints(fn func(node *SingleChildNode[P], parent Constraints) Constraints) {
	s.onWantedConstraints = fn
}

func (s *SingleChildProps[P]) Load(fn func(node *SingleChildNode[P])) {
	s.onLoad = fn
}

func (s *SingleChildProps[P]) Unload(fn func(node *SingleChildNode[P])) {
	s.onUnload = fn
}

func (s *SingleChildProps[P]) Layout(fn func(node *SingleChildNode[P]) (scath.Vec, error)) {
	s.onLayout = fn
}

func (s *SingleChildProps[P]) HandleEvent(fn func(node *SingleChildNode[P], c *scaff.Context, event scaff.Event) error) {
	s.onHandleEvent = fn
}

func (s *SingleChildProps[P]) Update(fn func(node *SingleChildNode[P], c *scaff.Context) (bool, error)) {
	s.onUpdate = fn
}

func (s *SingleChildProps[P]) Draw(fn func(node *SingleChildNode[P], position scath.Vec, renderer paint.Painter)) {
	s.onDraw = fn
}

var _ Node = &SingleChildNode[any]{}
var _ WantsConstraints = &SingleChildNode[any]{}

type SingleChildNode[P any] struct {
	tracker     *SingleTracker
	size        scath.Vec
	constraints Constraints

	id          string
	props       P
	singleProps *SingleChildProps[P]
}

func (s *SingleChildNode[P]) ID() string {
	return s.id
}

func (s *SingleChildNode[P]) Props() P {
	return s.props
}

func (s *SingleChildNode[P]) Load(parent Node) {
	if s.singleProps.onLoad != nil {
		s.singleProps.onLoad(s)
	}

	s.tracker.Load(parent)
}

func (s *SingleChildNode[P]) Size() scath.Vec {
	return s.size
}

func (s *SingleChildNode[P]) Constraints() Constraints {
	return s.constraints
}

func (s *SingleChildNode[P]) SetConstraints(c Constraints) {
	s.constraints = c
}

func (s *SingleChildNode[P]) WantedConstraints(parent Constraints) Constraints {
	if s.singleProps.onWantedConstraints == nil {
		return Unconstrained()
	}

	return s.singleProps.onWantedConstraints(s, parent)
}

func (s *SingleChildNode[P]) Layout() (scath.Vec, error) {
	size, err := s.layout()
	if err != nil {
		return scath.Vec{}, err
	}

	s.size = size
	return size, nil
}

func (s *SingleChildNode[P]) layout() (scath.Vec, error) {
	if s.singleProps.onLayout != nil {
		size, err := s.singleProps.onLayout(s)
		if err != nil {
			return size, scaff.NewTracedError(s, err)
		}
		s.size = size
		return size, nil
	}

	// As a default just take the size of the child
	size := scath.Vec{X: s.constraints.MinX, Y: s.constraints.MinY}

	if child, ok := s.tracker.Node(); ok {
		child.Current().SetConstraints(s.constraints)
		childSize, err := child.Current().Layout()
		if err != nil {
			return scath.Vec{}, scaff.NewTracedError(s, err)
		}

		size.X = max(size.X, childSize.X)
		size.Y = max(size.Y, childSize.Y)
	}

	return size, nil
}

func (s *SingleChildNode[P]) HandleEvent(c *scaff.Context, event scaff.Event) *scaff.TracedError {
	if s.singleProps.onHandleEvent != nil {
		if err := s.singleProps.onHandleEvent(s, c, event); err != nil {
			return scaff.NewTracedError(s, err)
		}
	}

	return s.HandleEventChild(c, event)
}

func (s *SingleChildNode[P]) Tracker() *scaff.Tracker {
	return s.tracker.Tracker()
}

func (s *SingleChildNode[P]) Update(c *scaff.Context) (UpdateResult, *scaff.TracedError) {
	relayout := false
	var err error
	if s.singleProps.onUpdate != nil {
		relayout, err = s.singleProps.onUpdate(s, c)
		if err != nil {
			return NoUpdate(), scaff.NewTracedError(s, err)
		}
	}

	// We should still update all the children after our own update
	result, updateErr := s.tracker.Update(s, c)
	if updateErr != nil {
		return NoUpdate(), scaff.NewTracedError(s, updateErr)
	}
	if relayout {
		result.Stack(SizeChanged())
	}
	return result, nil
}

func (s *SingleChildNode[P]) Unload() {
	if s.singleProps.onUnload != nil {
		s.singleProps.onUnload(s)
	}

	if s.tracker != nil {
		s.tracker.Unload()
		s.tracker = nil
	}
}

func (s *SingleChildNode[P]) Draw(position scath.Vec, renderer paint.Painter) {
	if s.singleProps.onDraw != nil {
		s.singleProps.onDraw(s, position, renderer)
		return
	}

	s.DrawChild(position, renderer)
}

func (s *SingleChildNode[P]) Child() (*MountedNode, bool) {
	return s.tracker.Node()
}

func (s *SingleChildNode[P]) LayoutChild(constraints Constraints) (scath.Vec, error) {
	child, ok := s.tracker.Node()
	if !ok {
		return scath.Vec{}, nil
	}

	child.Current().SetConstraints(constraints)
	return child.Current().Layout()
}

func (s *SingleChildNode[P]) HandleEventChild(c *scaff.Context, event scaff.Event) *scaff.TracedError {
	child, ok := s.tracker.Node()
	if !ok {
		return nil
	}

	// scaff.Event should always be passed to the child as well so it doesn't get missed (even if already handled)
	return child.Current().HandleEvent(c, event)
}

// Draw the child of the node somewhere
func (s *SingleChildNode[P]) DrawChild(position scath.Vec, renderer paint.Painter) {
	child, ok := s.tracker.Node()
	if !ok {
		return
	}

	child.Current().Draw(position, renderer)
}
