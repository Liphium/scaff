package scaffui

import (
	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/smath"
)

func UseSingleNode[P any](id string, propsCreator func(t *Tracker, props *P), create func(core *SingleChildConstruct[P])) NodeBuilder {
	return func() Node {
		node := newSingleChildCoreNode[P]()
		node.id = id

		// Fill the props
		var props P
		propsCreator(node.Tracker(), &props)
		node.props = props

		// Create the actual node
		create(&SingleChildConstruct[P]{node: node})

		return node
	}
}

type SingleChildConstruct[P any] struct {
	node *SingleChildNode[P]
}

func (s *SingleChildConstruct[P]) Tracker() *Tracker {
	return s.node.Tracker()
}

func (s *SingleChildConstruct[P]) Props() P {
	return s.node.props
}

func (s *SingleChildConstruct[P]) Child(builder NodeBuilder) {
	s.node.tracker.SetNode(NewMountedFromBuilder(builder))
}

func (s *SingleChildConstruct[P]) WantedConstraints(fn func(node *SingleChildNode[P], parent Constraints) Constraints) {
	s.node.onWantedConstraints = fn
}

func (s *SingleChildConstruct[P]) Load(fn func(node *SingleChildNode[P])) {
	s.node.onLoad = fn
}

func (s *SingleChildConstruct[P]) Unload(fn func(node *SingleChildNode[P])) {
	s.node.onUnload = fn
}

func (s *SingleChildConstruct[P]) Layout(fn func(node *SingleChildNode[P]) (Size, error)) {
	s.node.onLayout = fn
}

func (s *SingleChildConstruct[P]) HandleEvent(fn func(node *SingleChildNode[P], c *scaff.LayerContext, event Event) error) {
	s.node.onHandleEvent = fn
}

func (s *SingleChildConstruct[P]) Update(fn func(node *SingleChildNode[P], c *scaff.LayerContext) (bool, error)) {
	s.node.onUpdate = fn
}

func (s *SingleChildConstruct[P]) Draw(fn func(node *SingleChildNode[P], position smath.Vec, renderer Renderer)) {
	s.node.onDraw = fn
}

func newSingleChildCoreNode[P any]() *SingleChildNode[P] {
	n := &SingleChildNode[P]{}
	n.tracker = NewSingleTracker(n)
	return n
}

var _ Node = &SingleChildNode[any]{}
var _ WantsConstraints = &SingleChildNode[any]{}

type SingleChildNode[P any] struct {
	tracker     *SingleTracker
	size        Size
	constraints Constraints

	id                  string
	props               P
	onLoad              func(node *SingleChildNode[P])
	onUnload            func(node *SingleChildNode[P])
	onWantedConstraints func(node *SingleChildNode[P], parent Constraints) Constraints
	onLayout            func(node *SingleChildNode[P]) (Size, error)
	onHandleEvent       func(node *SingleChildNode[P], c *scaff.LayerContext, event Event) error
	onUpdate            func(node *SingleChildNode[P], c *scaff.LayerContext) (bool, error)
	onDraw              func(node *SingleChildNode[P], position smath.Vec, renderer Renderer)
}

func (s *SingleChildNode[P]) ID() string {
	return s.id
}

func (s *SingleChildNode[P]) Props() P {
	return s.props
}

func (s *SingleChildNode[P]) Load() {
	if s.onLoad != nil {
		s.onLoad(s)
	}

	s.tracker.Load()
}

func (s *SingleChildNode[P]) Size() Size {
	return s.size
}

func (s *SingleChildNode[P]) Constraints() Constraints {
	return s.constraints
}

func (s *SingleChildNode[P]) SetConstraints(c Constraints) {
	s.constraints = c
}

func (s *SingleChildNode[P]) WantedConstraints(parent Constraints) Constraints {
	if s.onWantedConstraints == nil {
		return Unconstrained()
	}

	return s.onWantedConstraints(s, parent)
}

func (s *SingleChildNode[P]) Layout() (Size, error) {
	size, err := s.layout()
	if err != nil {
		return Size{}, err
	}

	s.size = size
	return size, nil
}

func (s *SingleChildNode[P]) layout() (Size, error) {
	if s.onLayout != nil {
		size, err := s.onLayout(s)
		if err != nil {
			return size, NewError(s, err)
		}
		s.size = size
		return size, nil
	}

	// As a default just take the size of the child
	size := Size{Width: s.constraints.MinWidth, Height: s.constraints.MinHeight}

	if child, ok := s.tracker.Node(); ok {
		child.Current().SetConstraints(s.constraints)
		childSize, err := child.Current().Layout()
		if err != nil {
			return Size{}, NewError(s, err)
		}

		size.Width = max(size.Width, childSize.Width)
		size.Height = max(size.Height, childSize.Height)
	}

	return size, nil
}

func (s *SingleChildNode[P]) HandleEvent(c *scaff.LayerContext, event Event) *Error {
	if s.onHandleEvent != nil {
		if err := s.onHandleEvent(s, c, event); err != nil {
			return NewError(s, err)
		}
	}

	return s.HandleEventChild(c, event)
}

func (s *SingleChildNode[P]) Tracker() *Tracker {
	return s.tracker.Tracker()
}

func (s *SingleChildNode[P]) Update(c *scaff.LayerContext) (bool, *Error) {
	relayout := false
	var err error
	if s.onUpdate != nil {
		relayout, err = s.onUpdate(s, c)
		if err != nil {
			return false, NewError(s, err)
		}
	}

	// We should still update all the children after our own update
	childRelayout, updateErr := s.tracker.Update(c)
	if updateErr != nil {
		updateErr.add(s)
		return false, updateErr
	}
	return relayout || childRelayout, nil
}

func (s *SingleChildNode[P]) Unload() {
	if s.onUnload != nil {
		s.onUnload(s)
	}

	if s.tracker != nil {
		s.tracker.Unload()
		s.tracker = nil
	}
}

func (s *SingleChildNode[P]) Draw(position smath.Vec, renderer Renderer) {
	if s.onDraw != nil {
		s.onDraw(s, position, renderer)
		return
	}

	s.DrawChild(position, renderer)
}

func (s *SingleChildNode[P]) Child() (*MountedNode, bool) {
	return s.tracker.Node()
}

func (s *SingleChildNode[P]) LayoutChild(constraints Constraints) (Size, error) {
	child, ok := s.tracker.Node()
	if !ok {
		return Size{}, nil
	}

	child.Current().SetConstraints(constraints)
	return child.Current().Layout()
}

func (s *SingleChildNode[P]) HandleEventChild(c *scaff.LayerContext, event Event) *Error {
	child, ok := s.tracker.Node()
	if !ok {
		return nil
	}

	// Event should always be passed to the child as well so it doesn't get missed (even if already handled)
	return child.Current().HandleEvent(c, event)
}

// Draw the child of the node somewhere
func (s *SingleChildNode[P]) DrawChild(position smath.Vec, renderer Renderer) {
	child, ok := s.tracker.Node()
	if !ok {
		return
	}

	child.Current().Draw(position, renderer)
}
