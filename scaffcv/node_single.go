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
		if node.singleProps.OnPropsChanged != nil {
			node.singleProps.OnPropsChanged(node)
		}

		return node
	}
}

type SingleChildProps[P any] struct {
	OnLoad        func(node *SingleChildNode[P], parent Node)
	OnPropsChanged func(node *SingleChildNode[P])
	OnUnload      func(node *SingleChildNode[P])
	OnUpdate      func(node *SingleChildNode[P], c *scaff.Context) error
	OnHandleEvent func(node *SingleChildNode[P], c *scaff.Context, event scaff.Event) error
	OnDraw        func(node *SingleChildNode[P], c *scaff.Context, painter paint.Painter)
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

	if s.singleProps.OnLoad != nil {
		s.singleProps.OnLoad(s, parent)
	}
}

func (s *SingleChildNode[P]) HandleEvent(c *scaff.Context, event scaff.Event) scaff.TracedError {

	// First handle event on this node
	if s.singleProps.OnHandleEvent != nil {
		if err := s.singleProps.OnHandleEvent(s, c, event); err != nil {
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
	if s.singleProps.OnUpdate != nil {
		if err := s.singleProps.OnUpdate(s, c); err != nil {
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
	if s.singleProps.OnUnload != nil {
		s.singleProps.OnUnload(s)
	}

	// Unload the child properly
	if s.current != nil {
		s.current.Unload()
	}

	s.tracker.Clear()
	s.tracker = nil // Cut tracker off from tree for GC
}

func (s *SingleChildNode[P]) Draw(c *scaff.Context, painter paint.Painter) {
	if s.singleProps.OnDraw != nil {
		s.singleProps.OnDraw(s, c, painter)
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
