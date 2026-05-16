package scaff

import "github.com/hajimehoshi/ebiten/v2"

// CreateMultiNode lets you create a node with multiple children. Simply implement the ChildProps interface on the props you want to have for your node.
//
// id should be a unique id for the node, but also probably be readable as it shows up in error messages.
//
// propsCreator should be the function passed in by users of your node (as in it should probably be an argument of the function creating your node).
//
// create is the function actually specifying your node. You can overwrite all of the functions of the node interface there, with some exceptions that we implement for you.
func CreateMultiNode[P ChildProps[NodeBuilder]](id string, propsCreator func(t *Tracker, props *P), create func(props *MultiChildProps[P])) NodeBuilder {
	node := &MultiChildNode[P]{
		id:         id,
		multiProps: &MultiChildProps[P]{},
	}
	create(node.multiProps)

	return func() Node {
		node.tracker = NewTracker()

		// Fill the props
		var props P
		propsCreator(node.Tracker(), &props)
		node.props = props

		node.builders = props.GetBuilders()
		node.children = make([]Node, len(node.builders))
		for i, builder := range node.builders {
			node.children[i] = builder()
		}

		return node
	}
}

type MultiChildProps[P any] struct {
	onLoad        func(node *MultiChildNode[P], parent Node)
	onUnload      func(node *MultiChildNode[P])
	onUpdate      func(node *MultiChildNode[P], c *Context) error
	onHandleEvent func(node *MultiChildNode[P], c *Context, event Event) error
	onDraw        func(node *MultiChildNode[P], c *Context, image *ebiten.Image)
}

func (s *MultiChildProps[P]) Load(fn func(node *MultiChildNode[P], parent Node)) {
	s.onLoad = fn
}

func (s *MultiChildProps[P]) Unload(fn func(node *MultiChildNode[P])) {
	s.onUnload = fn
}

func (s *MultiChildProps[P]) Update(fn func(node *MultiChildNode[P], c *Context) error) {
	s.onUpdate = fn
}

func (s *MultiChildProps[P]) HandleEvent(fn func(node *MultiChildNode[P], c *Context, event Event) error) {
	s.onHandleEvent = fn
}

func (s *MultiChildProps[P]) Draw(fn func(node *MultiChildNode[P], c *Context, image *ebiten.Image)) {
	s.onDraw = fn
}

// Just for making sure we implement the Node interface
var _ Node = &MultiChildNode[any]{}

type MultiChildNode[P any] struct {
	parent   Node
	children []Node
	builders []NodeBuilder
	tracker  *Tracker

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

	// Actually load the children
	for _, child := range s.children {
		child.Load(s)
	}

	if s.multiProps.onLoad != nil {
		s.multiProps.onLoad(s, parent)
	}
}

func (s *MultiChildNode[P]) HandleEvent(c *Context, event Event) *TracedError {

	// First handle event on this node
	if s.multiProps.onHandleEvent != nil {
		if err := s.multiProps.onHandleEvent(s, c, event); err != nil {
			return NewTracedError(s, err)
		}
	}

	// Then pass event to children
	return s.HandleEventChild(c, event)
}

func (s *MultiChildNode[P]) Tracker() *Tracker {
	return s.tracker
}

func (s *MultiChildNode[P]) Update(c *Context) *TracedError {

	// First call the update handler on the props for this node
	if s.multiProps.onUpdate != nil {
		if err := s.multiProps.onUpdate(s, c); err != nil {
			return NewTracedError(s, err)
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
			s.children[i] = s.builders[i]()
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

func (s *MultiChildNode[P]) Draw(c *Context, image *ebiten.Image) {
	if s.multiProps.onDraw != nil {
		s.multiProps.onDraw(s, c, image)
	} else {

		// Default implementation: just draw children
		s.DrawChild(c, image)
	}
}

func (s *MultiChildNode[P]) Parent() Node {
	return s.parent
}

func (s *MultiChildNode[P]) Children() []Node {
	return s.children
}

func (s *MultiChildNode[P]) HandleEventChild(c *Context, event Event) *TracedError {
	// Event should always be passed to the children as well so it doesn't get missed (even if already handled)
	for _, child := range s.children {
		if err := child.HandleEvent(c, event); err != nil {
			return err
		}
	}

	return nil
}

// Draw the children of the node
func (s *MultiChildNode[P]) DrawChild(c *Context, image *ebiten.Image) {
	for _, child := range s.children {
		child.Draw(c, image)
	}
}
