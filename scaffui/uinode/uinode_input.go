package uinode

import (
	"github.com/Liphium/scaff/paint"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/scaffui"
	"github.com/Liphium/scaff/scath"
)

// Props for creating a new Input node. All of the listeners should return wether or not the event was handled, meaning no other UI components should handle the event.
type InputProps struct {

	// When a mouse button is pressed.
	OnDown func(handled, inside bool, event scaff.PressEvent) bool

	// When a mouse button is released.
	OnRelease func(handled, inside bool, event scaff.ReleaseEvent) bool

	// When the mouse is moved.
	OnMove func(handled, inside bool, event scaff.MoveEvent) bool

	// When scrolling with the mouse or potentially differnet methods when no mouse is available.
	OnScroll func(handled, inside bool, event scaff.ScrollEvent) bool

	// When a key is pressed.
	OnKeyPress func(handled bool, event scaff.KeyPressEvent) bool

	*scaffui.AcceptChild
}

// Create a new input node exposing a better interface to handle all kinds of input events coming down from scaffui.
func Input(create func(t *scaff.Tracker, props *InputProps)) scaffui.NodeBuilder {
	return scaffui.Standard(scaffui.StandardCreate[InputProps]{
		ID: "input",
		DefaultProps: InputProps{
			AcceptChild: &scaffui.AcceptChild{},
		},
		PropsCreator: create,
		Create: func(props *scaffui.StandardMethods[InputProps]) {
			lastPosition := scath.Vec{X: 0, Y: 0}

			props.OnHandleEvent = func(node *scaffui.StandardNode[InputProps], c *scaff.Context, event scaff.Event) error {
				handled := c.IsHandled(event.EventID())

				// If it is a positional event, check if the event was done within the current bounds
				isInside := false
				if posEvent, ok := event.(scaff.PositionalEvent); ok {
					isInside = posEvent.Position().IsWithinRectangle(lastPosition, node.Size())
				}

				switch ev := event.(type) {
				case scaff.PressEvent:
					if node.Props().OnDown != nil {
						if node.Props().OnDown(handled, isInside, ev) {
							c.Handled(event.EventID())
						}
					}
				case scaff.ReleaseEvent:
					if node.Props().OnRelease != nil {
						if node.Props().OnRelease(handled, isInside, ev) {
							c.Handled(event.EventID())
						}
					}
				case scaff.MoveEvent:
					if node.Props().OnMove != nil {
						if node.Props().OnMove(handled, isInside, ev) {
							c.Handled(event.EventID())
						}
					}
				case scaff.ScrollEvent:
					if node.Props().OnScroll != nil {
						if node.Props().OnScroll(handled, isInside, ev) {
							c.Handled(event.EventID())
						}
					}
				case scaff.KeyPressEvent:
					if node.Props().OnKeyPress != nil {
						if node.Props().OnKeyPress(handled, ev) {
							c.Handled(event.EventID())
						}
					}
				}

				return nil
			}

			props.OnDraw = func(node *scaffui.StandardNode[InputProps], position scath.Vec, renderer paint.Painter) {
				lastPosition = position
				if child, ok := node.Child(); ok {
					child.Draw(position, renderer)
				}
			}
		},
	})
}
