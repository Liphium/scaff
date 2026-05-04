package basenode

import (
	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/optional"
	"github.com/Liphium/scaff/scaffui"
	"github.com/Liphium/scaff/smath"
)

type InputProps struct {
	child            optional.O[scaffui.NodeBuilder]
	onDown           optional.O[func(button int) bool]
	onDownOutside    optional.O[func(button int) bool]
	onRelease        optional.O[func(button int) bool]
	onReleaseOutside optional.O[func(button int) bool]
	onMove           optional.O[func(deltaX, deltaY int) bool]
	onMoveOutside    optional.O[func(deltaX, deltaY int) bool]
	onScroll         optional.O[func(scrollX, scrollY float64) bool]
}

func (o *InputProps) Child(builder scaffui.NodeBuilder) {
	o.child.SetValue(builder)
}

func (o *InputProps) OnDown(fn func(button int) bool) {
	o.onDown.SetValue(fn)
}

func (o *InputProps) OnDownOutside(fn func(button int) bool) {
	o.onDownOutside.SetValue(fn)
}

func (o *InputProps) OnRelease(fn func(button int) bool) {
	o.onRelease.SetValue(fn)
}

func (o *InputProps) OnReleaseOutside(fn func(button int) bool) {
	o.onReleaseOutside.SetValue(fn)
}

func (o *InputProps) OnMove(fn func(deltaX, deltaY int) bool) {
	o.onMove.SetValue(fn)
}

func (o *InputProps) OnMoveOutside(fn func(deltaX, deltaY int) bool) {
	o.onMoveOutside.SetValue(fn)
}

func (o *InputProps) OnScroll(fn func(scrollX, scrollY float64) bool) {
	o.onScroll.SetValue(fn)
}

func Input(create func(t *scaffui.Tracker, props *InputProps)) scaffui.NodeBuilder {
	return scaffui.UseSingleNode("input", create, func(core *scaffui.SingleChildConstruct[InputProps]) {

		lastPosition := smath.Vec{X: 0, Y: 0}

		if child, ok := core.Props().child.Value(); ok {
			core.Child(child)
		}

		core.HandleEvent(func(node *scaffui.SingleChildNode[InputProps], c *scaff.LayerContext, event scaffui.Event) error {
			if c.IsHandled(event.EventID()) {
				return nil
			}

			isInside := false
			if posEvent, ok := event.(scaffui.PositionalEvent); ok {
				isInside = scaffui.IsWithin(lastPosition, node.Size(), posEvent.Position())
			} else {
				return nil
			}

			switch ev := event.(type) {
			case scaffui.DownEvent:
				if isInside {
					if fn, ok := core.Props().onDown.Value(); ok {
						if fn(ev.Button) {
							c.Handled(event.EventID())
						}
					}
				} else {
					if fn, ok := core.Props().onDownOutside.Value(); ok {
						if fn(ev.Button) {
							c.Handled(event.EventID())
						}
					}
				}
			case scaffui.ReleaseEvent:
				if isInside {
					if fn, ok := core.Props().onRelease.Value(); ok {
						if fn(ev.Button) {
							c.Handled(event.EventID())
						}
					}
				} else {
					if fn, ok := core.Props().onReleaseOutside.Value(); ok {
						if fn(ev.Button) {
							c.Handled(event.EventID())
						}
					}
				}
			case scaffui.MoveEvent:
				if isInside {
					if fn, ok := core.Props().onMove.Value(); ok {
						if fn(ev.DeltaX, ev.DeltaY) {
							c.Handled(event.EventID())
						}
					}
				} else {
					if fn, ok := core.Props().onMoveOutside.Value(); ok {
						if fn(ev.DeltaX, ev.DeltaY) {
							c.Handled(event.EventID())
						}
					}
				}
			case scaffui.ScrollEvent:
				if isInside {
					if fn, ok := core.Props().onScroll.Value(); ok {
						if fn(ev.ScrollX, ev.ScrollY) {
							c.Handled(event.EventID())
						}
					}
				}
			}

			return nil
		})

		core.Draw(func(node *scaffui.SingleChildNode[InputProps], position smath.Vec, renderer scaffui.Renderer) {
			lastPosition = position
			node.DrawChild(position, renderer)
		})
	})
}
