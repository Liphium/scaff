package scaffcv

import (
	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/paint"
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

	return func(context *BuildContext) Node {
		node.tracker = scaff.NewTracker()
		node.context = context

		// Fill the props
		props := create.DefaultProps
		if create.PropsCreator != nil {
			create.PropsCreator(node.Tracker(), &props)

			node.builders = props.GetBuilders()
		}
		node.props = props

		// Call props change hook (in case registered)
		if node.multiProps.onPropsChange != nil {
			node.multiProps.onPropsChange(node)
		}

		return node
	}
}

type MultiChildProps[P any] struct {
	onLoad        func(node *MultiChildNode[P], parent Node)
	onPropsChange func(node *MultiChildNode[P])
	onUnload      func(node *MultiChildNode[P])
	onUpdate      func(node *MultiChildNode[P], c *scaff.Context) error
	onHandleEvent func(node *MultiChildNode[P], c *scaff.Context, event scaff.Event) error
	onDraw        func(node *MultiChildNode[P], c *scaff.Context, image paint.Painter)
}

func (s *MultiChildProps[P]) Load(fn func(node *MultiChildNode[P], parent Node)) {
	s.onLoad = fn
}

func (s *MultiChildProps[P]) PropsChanged(fn func(node *MultiChildNode[P])) {
	s.onPropsChange = fn
}

func (s *MultiChildProps[P]) Unload(fn func(node *MultiChildNode[P])) {
	s.onUnload = fn
}

func (s *MultiChildProps[P]) Update(fn func(node *MultiChildNode[P], c *scaff.Context) error) {
	s.onUpdate = fn
}

func (s *MultiChildProps[P]) HandleEvent(fn func(node *MultiChildNode[P], c *scaff.Context, event scaff.Event) error) {
	s.onHandleEvent = fn
}

func (s *MultiChildProps[P]) Draw(fn func(node *MultiChildNode[P], c *scaff.Context, painter paint.Painter)) {
	s.onDraw = fn
}

// Just for making sure we implement the Node interface
var _ Node = &MultiChildNode[any]{}

type MultiChildNode[P any] struct {
	parent   Node
	children []Node
	builders []NodeBuilder

	tracker *scaff.Tracker
	context *BuildContext

	id         string
	props      P
	multiProps *MultiChildProps[P]
}

func (s *MultiChildNode[P]) ID() string {
	return s.id
}

func (s *MultiChildNode[P]) Props() P {
	return s.props
}

func (s *MultiChildNode[P]) Load(parent Node) {
	s.parent = parent

	// Actually load the children and build them
	if s.builders != nil {
		s.children = make([]Node, len(s.builders))
		for i, builder := range s.builders {
			s.children[i] = builder(s.context)
			s.children[i].Load(s)
		}
	}

	if s.multiProps.onLoad != nil {
		s.multiProps.onLoad(s, parent)
	}
}

func (s *MultiChildNode[P]) HandleEvent(c *scaff.Context, event scaff.Event) scaff.TracedError {

	// First handle event on this node
	if s.multiProps.onHandleEvent != nil {
		if err := s.multiProps.onHandleEvent(s, c, event); err != nil {
			return scaff.NewTracedError(s, err)
		}
	}

	// Then pass event to children
	return s.HandleEventChild(c, event)
}

func (s *MultiChildNode[P]) Tracker() *scaff.Tracker {
	return s.tracker
}

func (s *MultiChildNode[P]) Update(c *scaff.Context) scaff.TracedError {

	// First call the update handler on the props for this node
	if s.multiProps.onUpdate != nil {
		if err := s.multiProps.onUpdate(s, c); err != nil {
			return scaff.NewTracedError(s, err)
		}
	}

	// Forward the update to the children
	for _, child := range s.children {
		if err := child.Update(c); err != nil {
			return err
		}
	}

	// If any of the children are dirty, rebuild them
	for i, child := range s.children {
		if child.Tracker().Changed() {
			child.Unload()
			s.children[i] = s.builders[i](s.context)
			s.children[i].Load(s)
		}
	}

	return nil
}

func (s *MultiChildNode[P]) Unload() {
	if s.multiProps.onUnload != nil {
		s.multiProps.onUnload(s)
	}

	// Unload the children properly
	for _, child := range s.children {
		child.Unload()
	}

	s.tracker.Clear()
	s.tracker = nil // Cut tracker off from tree for GC
}

func (s *MultiChildNode[P]) Draw(c *scaff.Context, painter paint.Painter) {
	if s.multiProps.onDraw != nil {
		s.multiProps.onDraw(s, c, painter)
	} else {

		// Default implementation: just draw children
		s.DrawChild(c, painter)
	}
}

func (s *MultiChildNode[P]) Parent() Node {
	return s.parent
}

func (s *MultiChildNode[P]) Children() []Node {
	return s.children
}

func (s *MultiChildNode[P]) HandleEventChild(c *scaff.Context, event scaff.Event) scaff.TracedError {
	// Event should always be passed to the children as well so it doesn't get missed (even if already handled)
	for _, child := range s.children {
		if err := child.HandleEvent(c, event); err != nil {
			return err
		}
	}

	return nil
}

// Draw the children of the node
func (s *MultiChildNode[P]) DrawChild(c *scaff.Context, painter paint.Painter) {
	for _, child := range s.children {
		child.Draw(c, painter)
	}
}
