package scaff

import "github.com/hajimehoshi/ebiten/v2"

// Struct for defining a new standard node.
//
// ID should be a unique id for the node, but also probably be readable as it shows up in error messages.
//
// DefaultProps are the props as they are by default. Make sure to also specify any embedded struct pointers.
//
// PropsCreator should be the function passed in by users of your node (as in it should probably be an argument of the function creating your node).
//
// Create is the function actually specifying your node. You can overwrite all of the functions of the node interface there, with some exceptions that we implement for you.
type StandardCreate[P ChildProps[NodeBuilder]] struct {
	ID           string
	DefaultProps P
	PropsCreator func(t *Tracker, props *P)
	Create       func(methods *StandardMethods[P])
}

// Standard lets you create a node with multiple children. Simply implement the ChildProps interface on the props you want to have for your node.
func Standard[P ChildProps[NodeBuilder]](create StandardCreate[P]) NodeBuilder {
	node := &StandardNode[P]{
		id:      create.ID,
		methods: &StandardMethods[P]{},
	}
	if create.Create != nil {
		create.Create(node.methods)
	}

	props := create.DefaultProps
	return func(context *BuildContext) Node {
		node.tracker = NewTracker(context, func() {
			node.PropsChanged(props)
		})
		node.context = context

		// Fill the props
		if create.PropsCreator != nil {
			create.PropsCreator(node.Tracker(), &props)
		}
		node.props = props

		return node
	}
}

type StandardMethods[P ChildProps[NodeBuilder]] struct {
	OnLoad         func(node *StandardNode[P], parent Node)
	OnUnload       func(node *StandardNode[P])
	OnPropsChanged func(node *StandardNode[P])
	OnUpdate       func(node *StandardNode[P], c *Context) error
	OnHandleEvent  func(node *StandardNode[P], c *Context, event Event) error
	OnDraw         func(node *StandardNode[P], c *Context, image *ebiten.Image)
}

// Just for making sure we implement the Node interface
var _ Node = &StandardNode[*AcceptChildren]{}

type StandardNode[P ChildProps[NodeBuilder]] struct {
	parent   Node
	children []Node

	tracker *Tracker
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

func (s *StandardNode[P]) Context() *BuildContext {
	return s.context
}

func (s *StandardNode[P]) Load(parent Node) {
	if s.methods.OnLoad != nil {
		s.methods.OnLoad(s, parent)
	}

	// Set parent + build the children
	s.parent = parent
	s.PropsChanged(s.props)
}

func (s *StandardNode[P]) PropsChanged(new P) {
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

func (s *StandardNode[P]) HandleEvent(c *Context, event Event) TracedError {

	// First handle event on this node
	if s.methods.OnHandleEvent != nil {
		if err := s.methods.OnHandleEvent(s, c, event); err != nil {
			return NewTracedError(s, err)
		}
	}

	// Then pass event to children
	return s.HandleEventChild(c, event)
}

func (s *StandardNode[P]) Tracker() *Tracker {
	return s.tracker
}

func (s *StandardNode[P]) Update(c *Context) TracedError {

	// First call the update handler on the props for this node
	if s.methods.OnUpdate != nil {
		if err := s.methods.OnUpdate(s, c); err != nil {
			return NewTracedError(s, err)
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

func (s *StandardNode[P]) Draw(c *Context, image *ebiten.Image) {
	if s.methods.OnDraw != nil {
		s.methods.OnDraw(s, c, image)
	} else {
		// Default implementation: just draw children
		s.DrawChild(c, image)
	}
}

func (s *StandardNode[P]) Parent() Node {
	return s.parent
}

func (s *StandardNode[P]) Children() []Node {
	return s.children
}

func (s *StandardNode[P]) HandleEventChild(c *Context, event Event) TracedError {
	// Event should always be passed to the children as well so it doesn't get missed (even if already handled)
	for _, child := range s.children {
		if err := child.HandleEvent(c, event); err != nil {
			return err
		}
	}

	return nil
}

// Draw the children of the node
func (s *StandardNode[P]) DrawChild(c *Context, image *ebiten.Image) {
	for _, child := range s.children {
		child.Draw(c, image)
	}
}
