package basenode

import (
	"github.com/Liphium/scaff/paint"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/optional"
	"github.com/Liphium/scaff/scaffui"
	"github.com/Liphium/scaff/scath"
)

// Props for creating a new Input node. All of the listeners should return wether or not the event was handled, meaning no other UI components should handle the event.
type InputProps struct {

	// When a mouse button is pressed.
	onDown optional.O[func(handled, inside bool, event scaff.DownEvent) bool]

	// When a mouse button is released.
	onRelease optional.O[func(handled, inside bool, event scaff.ReleaseEvent) bool]

	// When the mouse is moved.
	onMove optional.O[func(handled, inside bool, event scaff.MoveEvent) bool]

	// When scrolling with the mouse or potentially differnet methods when no mouse is available.
	onScroll optional.O[func(handled, inside bool, event scaff.ScrollEvent) bool]

	*scaffui.AcceptChild
}

func (o *InputProps) OnDown(fn func(handled, inside bool, event scaff.DownEvent) bool) {
	o.onDown.SetValue(fn)
}

func (o *InputProps) OnRelease(fn func(handled, inside bool, event scaff.ReleaseEvent) bool) {
	o.onRelease.SetValue(fn)
}

func (o *InputProps) OnMove(fn func(handled, inside bool, event scaff.MoveEvent) bool) {
	o.onMove.SetValue(fn)
}

func (o *InputProps) OnScroll(fn func(handled, inside bool, event scaff.ScrollEvent) bool) {
	o.onScroll.SetValue(fn)
}

// Create a new input node exposing a better interface to handle all kinds of input events coming down from scaffui.
func Input(create func(t *scaff.Tracker, props *InputProps)) scaffui.NodeBuilder {
	return scaffui.CreateSingleNode(scaffui.SingleNodeCreate[InputProps]{
		ID: "input",
		DefaultProps: InputProps{
			AcceptChild: &scaffui.AcceptChild{},
		},
		PropsCreator: create,
		Create: func(props *scaffui.SingleChildProps[InputProps]) {
			lastPosition := scath.Vec{X: 0, Y: 0}

			props.HandleEvent(func(node *scaffui.SingleChildNode[InputProps], c *scaff.Context, event scaff.Event) error {
				handled := c.IsHandled(event.EventID())

				// If it is a positional event, check if the event was done within the current bounds
				isInside := false
				if posEvent, ok := event.(scaff.PositionalEvent); ok {
					isInside = posEvent.Position().IsWithinRectangle(lastPosition, node.Size())
				} else {
					return nil
				}

				switch ev := event.(type) {
				case scaff.DownEvent:
					if fn, ok := node.Props().onDown.Value(); ok {
						if fn(handled, isInside, ev) {
							c.Handled(event.EventID())
						}
					}
				case scaff.ReleaseEvent:
					if fn, ok := node.Props().onRelease.Value(); ok {
						if fn(handled, isInside, ev) {
							c.Handled(event.EventID())
						}
					}
				case scaff.MoveEvent:
					if fn, ok := node.Props().onMove.Value(); ok {
						if fn(handled, isInside, ev) {
							c.Handled(event.EventID())
						}
					}
				case scaff.ScrollEvent:
					if fn, ok := node.Props().onScroll.Value(); ok {
						if fn(handled, isInside, ev) {
							c.Handled(event.EventID())
						}
					}
				}

				return nil
			})

			props.Draw(func(node *scaffui.SingleChildNode[InputProps], position scath.Vec, renderer paint.Painter) {
				lastPosition = position
				node.DrawChild(position, renderer)
			})
		},
	})
}
