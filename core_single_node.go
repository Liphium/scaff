package scaff

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func UseSingleNode[P any](id string, propsCreator func(t *Tracker, props *P), create func(core *SingleChildProps[P])) NodeBuilder {
	return func() Node {
		node := &SingleChildNode[P]{
			id:      id,
			tracker: NewTracker(),
		}

		// Fill the props
		var props P
		propsCreator(node.Tracker(), &props)
		node.props = props

		// Create the actual node
		singleProps := &SingleChildProps[P]{
			props: props,
		}
		create(singleProps)
		node.singleProps = singleProps

		return node
	}
}

type SingleChildProps[P any] struct {
	props    P
	builder  NodeBuilder
	onLoad   func(node *SingleChildNode[P], parent Node)
	onUnload func(node *SingleChildNode[P])
	onUpdate func(node *SingleChildNode[P], c *Context) error
	onEvent  func(node *SingleChildNode[P], c *Context, event Event) error
	onDraw   func(node *SingleChildNode[P], c *Context, image *ebiten.Image)
}

func (s *SingleChildProps[P]) Props() P {
	return s.props
}

func (s *SingleChildProps[P]) Child(builder NodeBuilder) {
	s.builder = builder
}

var _ Node = &SingleChildNode[any]{}

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

	// If there is a builder, actually build and load the child
	if s.singleProps.builder != nil {
		s.current = s.builder()
		s.current.Load(s)
	}

	if s.singleProps.onLoad != nil {
		s.singleProps.onLoad(s, parent)
	}
}

func (s *SingleChildNode[P]) HandleEvent(c *Context, event Event) *TracedError {

	// First handle event on this node
	if s.singleProps.onEvent != nil {
		if err := s.singleProps.onEvent(s, c, event); err != nil {
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
	if s.tracker.Changed() {
		s.current.Unload()
		s.current = s.builder()
		s.current.Load(s.parent)
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
