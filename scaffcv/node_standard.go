package scaffcv

import (
	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/paint"
)

// Struct for defining a new standard node.
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
		id:      create.ID,
		methods: &StandardMethods[P]{},
	}
	if create.Create != nil {
		create.Create(node.methods)
	}

	return func(context *BuildContext) Node {
		node.tracker = scaff.NewTracker(context.BuildContext, node.PropsChanged)
		node.context = context

		// Fill the props
		props := create.DefaultProps
		if create.PropsCreator != nil {
			create.PropsCreator(node.Tracker(), &props)
		}
		node.props = props

		return node
	}
}

type StandardMethods[P scaff.ChildProps[NodeBuilder]] struct {
	OnLoad         func(node *StandardNode[P], parent Node)
	OnPropsChanged func(node *StandardNode[P])
	OnUnload       func(node *StandardNode[P])
	OnUpdate       func(node *StandardNode[P], c *scaff.Context) error
	OnHandleEvent  func(node *StandardNode[P], c *scaff.Context, event scaff.Event) error
	OnDraw         func(node *StandardNode[P], c *scaff.Context, image paint.Painter)
}

// Just for making sure we implement the Node interface
var _ Node = &StandardNode[*AcceptChildren]{}

type StandardNode[P scaff.ChildProps[NodeBuilder]] struct {
	parent   Node
	children []Node
	builders []NodeBuilder

	tracker *scaff.Tracker
	context *BuildContext

	id      string
	props   P
	methods *StandardMethods[P]
}

func (s *StandardNode[P]) ID() string {
	return s.id
}

func (s *StandardNode[P]) Props() P {
	return s.props
}

func (s *StandardNode[P]) Load(parent Node) {
	if s.methods.OnLoad != nil {
		s.methods.OnLoad(s, parent)
	}

	// Set parent + load children by calling PropsChanged
	s.parent = parent
	s.PropsChanged()
}

func (s *StandardNode[P]) PropsChanged() {

	// If some children changed, build new ones
	changed := s.props.GetChanged()
	if changed != nil {
		builders := s.props.GetBuilders()
		if s.children == nil {
			s.children = make([]Node, len(builders))
		}

		for _, i := range changed {
			if len(builders) < int(i) {
				log.Error("index out of bounds for props update", "i", i, "children", len(builders))
				continue
			}

			if s.children[i] != nil {
				s.children[i].Unload()
			}
			s.children[i] = builders[i](s.context)
			s.children[i].Load(s)
		}
		s.props.ClearChanged()
	}

	// Call props changed on the actual methods
	if s.methods.OnPropsChanged != nil {
		s.methods.OnPropsChanged(s)
	}
}

func (s *StandardNode[P]) HandleEvent(c *scaff.Context, event scaff.Event) scaff.TracedError {

	// First handle event on this node
	if s.methods.OnHandleEvent != nil {
		if err := s.methods.OnHandleEvent(s, c, event); err != nil {
			return scaff.NewTracedError(s, err)
		}
	}

	// Then pass event to children
	return s.HandleEventChildren(c, event)
}

func (s *StandardNode[P]) Tracker() *scaff.Tracker {
	return s.tracker
}

func (s *StandardNode[P]) Update(c *scaff.Context) scaff.TracedError {

	// First call the update handler on the props for this node
	if s.methods.OnUpdate != nil {
		if err := s.methods.OnUpdate(s, c); err != nil {
			return scaff.NewTracedError(s, err)
		}
	}

	// Forward the update to the children
	for _, child := range s.children {
		if err := child.Update(c); err != nil {
			return err
		}
	}

	return nil
}

func (s *StandardNode[P]) Unload() {
	if s.methods.OnUnload != nil {
		s.methods.OnUnload(s)
	}

	// Unload the children properly
	for _, child := range s.children {
		child.Unload()
	}

	s.tracker.Clear()
	s.tracker = nil // Cut tracker off from tree for GC
}

func (s *StandardNode[P]) Draw(c *scaff.Context, painter paint.Painter) {
	if s.methods.OnDraw != nil {
		s.methods.OnDraw(s, c, painter)
	} else {

		// Default implementation: just draw children
		s.DrawChild(c, painter)
	}
}

func (s *StandardNode[P]) Parent() Node {
	return s.parent
}

func (s *StandardNode[P]) Children() []Node {
	return s.children
}

func (s *StandardNode[P]) HandleEventChildren(c *scaff.Context, event scaff.Event) scaff.TracedError {
	// Event should always be passed to the children as well so it doesn't get missed (even if already handled)
	for _, child := range s.children {
		if err := child.HandleEvent(c, event); err != nil {
			return err
		}
	}

	return nil
}

// Draw the children of the node
func (s *StandardNode[P]) DrawChild(c *scaff.Context, painter paint.Painter) {
	for _, child := range s.children {
		child.Draw(c, painter)
	}
}
