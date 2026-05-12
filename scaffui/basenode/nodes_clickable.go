package basenode

import (
	"github.com/Liphium/scaff/optional"
	"github.com/Liphium/scaff/scaffui"
)

type ClickableProps struct {
	child   optional.O[scaffui.NodeBuilder]
	onClick optional.O[func(button int) bool]
}

func (cp *ClickableProps) Child(builder scaffui.NodeBuilder) {
	cp.child.SetValue(builder)
}

func (cp *ClickableProps) OnClick(fn func(button int) bool) {
	cp.onClick.SetValue(fn)
}

func Clickable(create func(t *scaffui.Tracker, props *ClickableProps)) scaffui.NodeBuilder {
	return scaffui.UseSingleNode("clickable", create, func(core *scaffui.SingleChildConstruct[ClickableProps]) {
		pressed := make(map[int]bool)

		core.Child(Input(func(t *scaffui.Tracker, ip *InputProps) {
			if child, ok := core.Props().child.Value(); ok {
				ip.Child(child)
			}

			ip.OnDown(func(handled, inside bool, event scaffui.DownEvent) bool {
				if inside {
					pressed[event.Button] = true
				}
				return inside
			})

			ip.OnRelease(func(handled, inside bool, event scaffui.ReleaseEvent) bool {
				wasPressed := pressed[event.Button]
				pressed[event.Button] = false

				// If the event was not handled before and the button was pressed before + released inside of this element, a click has been detected
				if wasPressed && inside && !handled {
					if fn, ok := core.Props().onClick.Value(); ok {
						return fn(event.Button)
					}
					return true
				}

				return false
			})
		}))
	})
}
