package scaffcv

import (
	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/paint"
)

// Struct for defining a new single child node.
//
// ID should be a unique id for the node, but also probably be readable as it shows up in error messages.
//
// DefaultProps are the props as they are by default. Make sure to also specify any embedded struct pointers.
//
// PropsCreator should be the function passed in by users of your node (as in it should probably be an argument of the function creating your node).
//
// Create is the function actually specifying your node. You can overwrite all of the functions of the node interface there, with some exceptions that we implement for you.
type SingleNodeCreate[P scaff.ChildProps[NodeBuilder]] struct {
	ID           string
	DefaultProps P
	PropsCreator func(t *scaff.Tracker, props *P)
	Create       func(props *SingleChildProps[P])
}

// SingleNode lets you create a node with a single child. Simply implement the ChildProps interface on the props you want to have for your node.
func SingleNode[P scaff.ChildProps[NodeBuilder]](create SingleNodeCreate[P]) NodeBuilder {

	// Create the actual node
	node := &SingleChildNode[P]{
		id:          create.ID,
		singleProps: &SingleChildProps[P]{},
	}
	if create.Create != nil {
		create.Create(node.singleProps)
	}

	return func(context *BuildContext) Node {
		node.tracker = scaff.NewTracker()
		node.context = context

		// Fill the props
		props := create.DefaultProps
		if create.PropsCreator != nil {
			create.PropsCreator(node.Tracker(), &props)

			// Build the children, in case there are any
			if builders := props.GetBuilders(); builders != nil {
				if len(builders) > 1 {
					log.Error("node can not have multiple children", "id", create.ID, "children", len(props.GetBuilders()))
				}
				if len(builders) == 1 {
					node.builder = props.GetBuilders()[0]
				}
			}
		}
		node.props = props

		// Run the state change hook
		if node.singleProps.onPropsChange != nil {
			node.singleProps.onPropsChange(node)
		}

		return node
	}
}

type SingleChildProps[P any] struct {
	onLoad        func(node *SingleChildNode[P], parent Node)
	onPropsChange func(node *SingleChildNode[P])
	onUnload      func(node *SingleChildNode[P])
	onUpdate      func(node *SingleChildNode[P], c *scaff.Context) error
	onHandleEvent func(node *SingleChildNode[P], c *scaff.Context, event scaff.Event) error
	onDraw        func(node *SingleChildNode[P], c *scaff.Context, painter paint.Painter)
}

func (s *SingleChildProps[P]) Load(fn func(node *SingleChildNode[P], parent Node)) {
	s.onLoad = fn
}

func (s *SingleChildProps[P]) PropsChanged(fn func(node *SingleChildNode[P])) {
	s.onPropsChange = fn
}

func (s *SingleChildProps[P]) Unload(fn func(node *SingleChildNode[P])) {
	s.onUnload = fn
}

func (s *SingleChildProps[P]) Update(fn func(node *SingleChildNode[P], c *scaff.Context) error) {
	s.onUpdate = fn
}

func (s *SingleChildProps[P]) HandleEvent(fn func(node *SingleChildNode[P], c *scaff.Context, event scaff.Event) error) {
	s.onHandleEvent = fn
}

func (s *SingleChildProps[P]) Draw(fn func(node *SingleChildNode[P], c *scaff.Context, painter paint.Painter)) {
	s.onDraw = fn
}

var _ Node = &SingleChildNode[any]{}

type SingleChildNode[P any] struct {
	parent  Node
	current Node
	builder NodeBuilder

	// Things internal to the node
	context *BuildContext
	tracker *scaff.Tracker

	// Configuration of the node
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

func (s *SingleChildNode[P]) Context() *BuildContext {
	return s.context
}

func (s *SingleChildNode[P]) Load(parent Node) {
	s.parent = parent

	// If there is a builder, load the child
	if s.builder != nil {
		s.current = s.builder(s.context)
		s.current.Load(s)
	}

	if s.singleProps.onLoad != nil {
		s.singleProps.onLoad(s, parent)
	}
}

func (s *SingleChildNode[P]) HandleEvent(c *scaff.Context, event scaff.Event) scaff.TracedError {

	// First handle event on this node
	if s.singleProps.onHandleEvent != nil {
		if err := s.singleProps.onHandleEvent(s, c, event); err != nil {
			return scaff.NewTracedError(s, err)
		}
	}

	// Then pass event to child
	return s.HandleEventChild(c, event)
}

func (s *SingleChildNode[P]) Tracker() *scaff.Tracker {
	return s.tracker
}

func (s *SingleChildNode[P]) Update(c *scaff.Context) scaff.TracedError {

	// First call the update handler on the props for this node
	if s.singleProps.onUpdate != nil {
		if err := s.singleProps.onUpdate(s, c); err != nil {
			return scaff.NewTracedError(s, err)
		}
	}

	// If there is no node, nothing to check
	if s.current == nil {
		return nil
	}

	// Forward the update to the child
	if err := s.current.Update(c); err != nil {
		return err
	}

	// If dirty, rebuild
	if s.current.Tracker().Changed() {
		s.current.Unload()
		s.current = s.builder(s.context)
		s.current.Load(s)
	}

	return nil
}

func (s *SingleChildNode[P]) Unload() {
	if s.singleProps.onUnload != nil {
		s.singleProps.onUnload(s)
	}

	// Unload the child properly
	if s.current != nil {
		s.current.Unload()
	}

	s.tracker.Clear()
	s.tracker = nil // Cut tracker off from tree for GC
}

func (s *SingleChildNode[P]) Draw(c *scaff.Context, painter paint.Painter) {
	if s.singleProps.onDraw != nil {
		s.singleProps.onDraw(s, c, painter)
	} else {

		// Default implementation: just draw child
		s.DrawChild(c, painter)
	}
}

func (s *SingleChildNode[P]) Parent() Node {
	return s.parent
}

func (s *SingleChildNode[P]) Children() []Node {
	if s.current == nil {
		return []Node{}
	}
	return []Node{s.current}
}

func (s *SingleChildNode[P]) HandleEventChild(c *scaff.Context, event scaff.Event) scaff.TracedError {
	if s.current == nil {
		return nil
	}

	// Event should always be passed to the child as well so it doesn't get missed (even if already handled)
	return s.current.HandleEvent(c, event)
}

// Draw the child of the node
func (s *SingleChildNode[P]) DrawChild(c *scaff.Context, painter paint.Painter) {
	if s.current != nil {
		s.current.Draw(c, painter)
	}
}
