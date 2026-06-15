package scaffui

import (
	"fmt"

	"github.com/Liphium/scaff/engine"

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
	Create       func(methods *StandardMethods[P])
}

// Standard lets you create a node with multiple children. Simply implement the ChildProps interface on the props you want to have for your node.
func Standard[P scaff.ChildProps[NodeBuilder]](create StandardCreate[P]) NodeBuilder {
	node := &StandardNode[P]{
		id:      create.ID,
		methods: &StandardMethods[P]{},
	}
	if create.Create != nil {
		create.Create(node.methods)
	}

	props := create.DefaultProps
	return func(context *scaff.BuildContext) Node {
		node.tracker = scaff.NewTracker(context, func() {
			node.PropsChanged(props)
		})
		node.context = context

		// Create the props (if desired)
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
	OnDraw              func(node *StandardNode[P], position scath.Vec, renderer engine.Painter)
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

	props   P
	dirty   bool
	methods *StandardMethods[P]
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

// Get the first child, this is just a helper function for nodes that have one child.
func (s *StandardNode[P]) Child() (Node, bool) {
	if len(s.children) > 0 {
		return s.children[0], true
	}
	return nil, false
}

func (s *StandardNode[P]) Props() P {
	return s.props
}

func (s *StandardNode[P]) Context() *scaff.BuildContext {
	return s.context
}

func (s *StandardNode[P]) Load(parent Node) {
	if s.methods.OnLoad != nil {
		s.methods.OnLoad(s)
	}

	// Set parent + build the children
	s.parent = parent
	s.PropsChanged(s.props)
}

func (s *StandardNode[P]) PropsChanged(new P) {
	s.dirty = true
	s.props = new

	// Call props changed on the actual methods
	if s.methods.OnPropsChanged != nil {
		s.methods.OnPropsChanged(s)
	}

	// If some children changed, build new ones
	changed := new.GetChanged()
	if changed != nil {
		builders := new.GetBuilders()

		// Diff size to make sure length is properly done
		if s.children == nil {
			s.children = make([]Node, len(builders))
		} else if len(s.children) > len(builders) {
			for _, node := range s.children[len(builders)-1:] {
				node.Unload()
			}
		} else if len(s.children) < len(builders) {
			s.children = append(s.children, make([]Node, len(builders)-len(s.children))...)
		}

		for _, i := range changed {
			if int(i) >= len(builders) {
				log.Error("index out of bounds for props update", "i", i, "children", len(builders))
				continue
			}

			if s.children[i] != nil {
				s.children[i].Unload()
				s.children[i] = nil
			}
			if builders[i] != nil {
				s.children[i] = builders[i](s.context)
				s.children[i].Load(s)
			}
		}
		new.ClearChanged()
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
	if s.methods.OnWantedConstraints == nil {
		return scath.Unconstrained()
	}

	return s.methods.OnWantedConstraints(s, parent)
}

func (s *StandardNode[P]) Layout() (scath.Vec, error) {

	// If layouting is handled by a method, use that
	if s.methods.OnLayout != nil {
		size, err := s.methods.OnLayout(s)
		if err != nil {
			return scath.Vec{}, scaff.NewTracedError(s, err)
		}

		finalSize, err := s.finalSize(size)
		s.size = finalSize
		return finalSize, err
	}

	return s.LayoutChildren()
}

// Default implementation of layout
func (s *StandardNode[P]) LayoutChildren() (scath.Vec, error) {

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
	if s.methods.OnHandleEvent != nil {
		if err := s.methods.OnHandleEvent(s, c, event); err != nil {
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
	if s.methods.OnUnload != nil {
		s.methods.OnUnload(s)
	}

	if s.tracker != nil {
		s.tracker.Clear()
		s.tracker = nil
	}
}

func (s *StandardNode[P]) Draw(position scath.Vec, renderer engine.Painter) {
	if s.methods.OnDraw != nil {
		s.methods.OnDraw(s, position, renderer)
		return
	}

	for _, child := range s.children {
		child.Draw(position, renderer)
	}
}
