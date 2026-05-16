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
	child optional.O[scaffui.NodeBuilder]

	// When a mouse button is pressed.
	onDown optional.O[func(handled, inside bool, event scaffui.DownEvent) bool]

	// When a mouse button is released.
	onRelease optional.O[func(handled, inside bool, event scaffui.ReleaseEvent) bool]

	// When the mouse is moved.
	onMove optional.O[func(handled, inside bool, event scaffui.MoveEvent) bool]

	// When scrolling with the mouse or potentially differnet methods when no mouse is available.
	onScroll optional.O[func(handled, inside bool, event scaffui.ScrollEvent) bool]
}

func (o *InputProps) Child(builder scaffui.NodeBuilder) {
	o.child.SetValue(builder)
}

func (o *InputProps) OnDown(fn func(handled, inside bool, event scaffui.DownEvent) bool) {
	o.onDown.SetValue(fn)
}

func (o *InputProps) OnRelease(fn func(handled, inside bool, event scaffui.ReleaseEvent) bool) {
	o.onRelease.SetValue(fn)
}

func (o *InputProps) OnMove(fn func(handled, inside bool, event scaffui.MoveEvent) bool) {
	o.onMove.SetValue(fn)
}

func (o *InputProps) OnScroll(fn func(handled, inside bool, event scaffui.ScrollEvent) bool) {
	o.onScroll.SetValue(fn)
}

// Create a new input node exposing a better interface to handle all kinds of input events coming down from scaffui.
func Input(create func(t *scaff.Tracker, props *InputProps)) scaffui.NodeBuilder {
	return scaffui.CreateSingleNode("input", create, func(core *scaffui.SingleChildProps[InputProps]) {

		lastPosition := scath.Vec{X: 0, Y: 0}

		if child, ok := core.Props().child.Value(); ok {
			core.Child(child)
		}

		core.HandleEvent(func(node *scaffui.SingleChildNode[InputProps], c *scaff.LayerContext, event scaffui.Event) error {
			handled := c.IsHandled(event.EventID())

			// If it is a positional event, check if the event was done within the current bounds
			isInside := false
			if posEvent, ok := event.(scaffui.PositionalEvent); ok {
				isInside = scaffui.IsWithin(lastPosition, node.Size(), posEvent.Position())
			} else {
				return nil
			}

			switch ev := event.(type) {
			case scaffui.DownEvent:
				if fn, ok := core.Props().onDown.Value(); ok {
					if fn(handled, isInside, ev) {
						c.Handled(event.EventID())
					}
				}
			case scaffui.ReleaseEvent:
				if fn, ok := core.Props().onRelease.Value(); ok {
					if fn(handled, isInside, ev) {
						c.Handled(event.EventID())
					}
				}
			case scaffui.MoveEvent:
				if fn, ok := core.Props().onMove.Value(); ok {
					if fn(handled, isInside, ev) {
						c.Handled(event.EventID())
					}
				}
			case scaffui.ScrollEvent:
				if fn, ok := core.Props().onScroll.Value(); ok {
					if fn(handled, isInside, ev) {
						c.Handled(event.EventID())
					}
				}
			}

			return nil
		})

		core.Draw(func(node *scaffui.SingleChildNode[InputProps], position scath.Vec, renderer paint.Painter) {
			lastPosition = position
			node.DrawChild(position, renderer)
		})
	})
}
