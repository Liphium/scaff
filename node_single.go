package scaff

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// UseNode creates a node without a child. You can use this if you simply want to create a node that does something very custom outside of scaff. This node still gives you access to all all functions on the Node interface through the props, but you can choose which ones you actually want to implement.
//
// If you want to have one or multiple children for this node, CreateSingleNode or CreateMultiNode might be a better fit for your use case.
func UseNode[P any](id string, create func(props *SingleChildProps[P])) NodeBuilder {
	node := &SingleChildNode[P]{
		id:      id,
		tracker: NewTracker(),
	}

	// Create the actual node
	singleProps := &SingleChildProps[P]{}
	create(singleProps)
	node.singleProps = singleProps

	return func() Node {
		return node
	}
}

// CreateSingleNode lets you create a node with a single child. Simply implement the ChildProps interface on the props you want to have for your node.
//
// id should be a unique id for the node, but also probably be readable as it shows up in error messages.
//
// propsCreator should be the function passed in by users of your node (as in it should probably be an argument of the function creating your node).
//
// create is the function actually specifying your node. You can overwrite all of the functions of the node interface there, with some exceptions that we implement for you.
func CreateSingleNode[P ChildProps](id string, propsCreator func(t *Tracker, props *P), create func(props *SingleChildProps[P])) NodeBuilder {
	node := &SingleChildNode[P]{
		id: id,
	}

	// Create the actual node
	singleProps := &SingleChildProps[P]{}
	create(singleProps)
	node.singleProps = singleProps

	return func() Node {
		node.tracker = NewTracker()

		// Fill the props
		var props P
		propsCreator(node.Tracker(), &props)
		node.props = props
		if len(props.GetBuilders()) > 1 {
			log.Error("node can not have multiple children", "id", id, "children", len(props.GetBuilders()))
		}
		if len(props.GetBuilders()) == 1 {
			node.builder = props.GetBuilders()[0]
			node.current = node.builder()
		}

		return node
	}
}

type SingleChildProps[P any] struct {
	onLoad        func(node *SingleChildNode[P], parent Node)
	onUnload      func(node *SingleChildNode[P])
	onUpdate      func(node *SingleChildNode[P], c *Context) error
	onHandleEvent func(node *SingleChildNode[P], c *Context, event Event) error
	onDraw        func(node *SingleChildNode[P], c *Context, image *ebiten.Image)
}

func (s *SingleChildProps[P]) Load(fn func(node *SingleChildNode[P], parent Node)) {
	s.onLoad = fn
}

func (s *SingleChildProps[P]) Unload(fn func(node *SingleChildNode[P])) {
	s.onUnload = fn
}

func (s *SingleChildProps[P]) Update(fn func(node *SingleChildNode[P], c *Context) error) {
	s.onUpdate = fn
}

func (s *SingleChildProps[P]) HandleEvent(fn func(node *SingleChildNode[P], c *Context, event Event) error) {
	s.onHandleEvent = fn
}

func (s *SingleChildProps[P]) Draw(fn func(node *SingleChildNode[P], c *Context, image *ebiten.Image)) {
	s.onDraw = fn
}

var _ Node = &MultiChildNode[any]{}

type SingleChildNode[P any] struct {
	parent  Node
	current Node
	builder NodeBuilder
	tracker *Tracker

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
	s.parent = parent

	// If there is a builder, actually build and load the child
	if s.builder != nil {
		s.current.Load(s)
	}

	if s.singleProps.onLoad != nil {
		s.singleProps.onLoad(s, parent)
	}
}

func (s *SingleChildNode[P]) HandleEvent(c *Context, event Event) *TracedError {

	// First handle event on this node
	if s.singleProps.onHandleEvent != nil {
		if err := s.singleProps.onHandleEvent(s, c, event); err != nil {
			return NewTracedError(s, err)
		}
	}

	// Then pass event to child
	return s.HandleEventChild(c, event)
}

func (s *SingleChildNode[P]) Tracker() *Tracker {
	return s.tracker
}

func (s *SingleChildNode[P]) Update(c *Context) *TracedError {

	// First call the update handler on the props for this node
	if s.singleProps.onUpdate != nil {
		if err := s.singleProps.onUpdate(s, c); err != nil {
			return NewTracedError(s, err)
		}
	}

	// If there is no node, nothing to check
	if s.current == nil {
		return nil
	}

	// If dirty, rebuild
	if s.current.Tracker().Changed() {
		s.current.Unload()
		s.current = s.builder()
		s.current.Load(s)
	}

	// Forward the update to the child
	return s.current.Update(c)
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

func (s *SingleChildNode[P]) Draw(c *Context, image *ebiten.Image) {
	if s.singleProps.onDraw != nil {
		s.singleProps.onDraw(s, c, image)
	} else {

		// Default implementation: just draw child
		s.DrawChild(c, image)
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

func (s *SingleChildNode[P]) HandleEventChild(c *Context, event Event) *TracedError {
	if s.current == nil {
		return nil
	}

	// Event should always be passed to the child as well so it doesn't get missed (even if already handled)
	return s.current.HandleEvent(c, event)
}

// Draw the child of the node
func (s *SingleChildNode[P]) DrawChild(c *Context, image *ebiten.Image) {
	if s.current != nil {
		s.current.Draw(c, image)
	}
}
