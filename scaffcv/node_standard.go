package scaffcv

import (
	"math"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/paint"
	"github.com/Liphium/scaff/scath"
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

	props := create.DefaultProps
	return func(context *BuildContext) Node {
		node.tracker = scaff.NewTracker(context.BuildContext, func() {
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

type StandardMethods[P scaff.ChildProps[NodeBuilder]] struct {
	// Should return your current position. This is used by various nodes to perform culling.
	Position func(node *StandardNode[P]) scath.Vec

	// Should return your current size. This is used by various nodes to perform culling.
	Size func(node *StandardNode[P]) scath.Vec

	// Called when the node is loaded.
	OnLoad func(node *StandardNode[P], parent Node)

	// Called when props of the node change.
	OnPropsChanged func(node *StandardNode[P])

	// Called when the node is unloaded.
	OnUnload func(node *StandardNode[P])

	// Called on every Ebiten tick.
	OnUpdate func(node *StandardNode[P], c *scaff.Context) error

	// Called for every event that comes through from scaff.
	OnHandleEvent func(node *StandardNode[P], c *scaff.Context, event scaff.Event) error

	// Should draw the node using the painter.
	OnDraw func(node *StandardNode[P], c *scaff.Context, painter paint.Painter)
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

func (s *StandardNode[P]) Context() *BuildContext {
	return s.context
}

func (s *StandardNode[P]) Position() scath.Vec {
	if s.methods.Position != nil {
		return s.methods.Position(s)
	}

	// As default: Calculate minimum X and Y (we want to make sure all children are in our bounding box, size will be adjusted properly)
	minX, minY := 0.0, 0.0
	for _, child := range s.children {
		pos := child.Position()
		minX = math.Min(pos.X, minX)
		minY = math.Min(pos.Y, minY)
	}
	return scath.Vec{X: minX, Y: minY}
}

func (s *StandardNode[P]) Size() scath.Vec {
	if s.methods.Size != nil {
		return s.methods.Size(s)
	}

	// As default: Calculate minimum X,Y + maximum X,Y and get the difference. This will make sure the bounding box of Position + Size contains all children.
	minX, maxX, minY, maxY := 0.0, 0.0, 0.0, 0.0
	for _, child := range s.children {
		pos := child.Position()
		minX = math.Min(pos.X, minX)
		maxX = math.Max(pos.X, maxX)
		minY = math.Min(pos.Y, minY)
		maxY = math.Max(pos.Y, maxY)
	}
	return scath.Vec{X: maxX - minX, Y: maxY - minY}
}

func (s *StandardNode[P]) Load(parent Node) {
	if s.methods.OnLoad != nil {
		s.methods.OnLoad(s, parent)
	}

	// Set parent + load children by calling PropsChanged
	s.parent = parent
	s.PropsChanged(s.props)
}

func (s *StandardNode[P]) PropsChanged(new P) {

	// If some children changed, build new ones
	changed := new.GetChanged()
	if changed != nil {
		builders := new.GetBuilders()
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
		new.ClearChanged()
	}

	// Call props changed on the actual methods
	if s.methods.OnPropsChanged != nil {
		s.methods.OnPropsChanged(s)
	}

	s.props = new
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
		s.DrawChildren(c, painter)
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
func (s *StandardNode[P]) DrawChildren(c *scaff.Context, painter paint.Painter) {
	vwOrigin, vwSize := s.Context().Camera().Viewport()

	for _, child := range s.children {

		// Do not draw chlildren that are outside of the camera's view
		topLeft, size := child.Position(), child.Size()
		if child.Size() == scath.Zero {
			continue
		}

		// TODO: This needs to support camera rotation in the future, but let's not worry about that for now
		toCheck := []scath.Vec{
			topLeft,
			topLeft.Add(scath.Vec{X: size.X}), // Top right
			topLeft.Add(scath.Vec{Y: size.Y}), // Bottom left
			topLeft.Add(size),                 // Bottom right
		}
		found := false
		for _, vec := range toCheck {
			if vec.IsWithinRectangle(vwOrigin, vwSize) {
				found = true
			}
		}
		if !found {
			continue
		}

		child.Draw(c, painter)
	}
}
