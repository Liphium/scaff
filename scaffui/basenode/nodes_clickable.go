package basenode

import (
	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/optional"
	"github.com/Liphium/scaff/scaffui"
	"github.com/hajimehoshi/ebiten/v2"
)

type ClickableProps struct {
	child   optional.O[scaffui.NodeBuilder]
	onClick optional.O[func(button ebiten.MouseButton) bool]
}

func (cp *ClickableProps) Child(builder scaffui.NodeBuilder) {
	cp.child.SetValue(builder)
}

func (cp *ClickableProps) OnClick(fn func(button ebiten.MouseButton) bool) {
	cp.onClick.SetValue(fn)
}

func Clickable(create func(t *scaff.Tracker, props *ClickableProps)) scaffui.NodeBuilder {
	return scaffui.CreateSingleNode("clickable", create, func(core *scaffui.SingleChildProps[ClickableProps]) {
		pressed := make(map[ebiten.MouseButton]bool)

		core.Child(Input(func(t *scaff.Tracker, ip *InputProps) {
			if child, ok := core.Props().child.Value(); ok {
				ip.Child(child)
			}

			ip.OnDown(func(handled, inside bool, event scaff.DownEvent) bool {
				if inside {
					pressed[event.Button] = true
				}
				return inside
			})

			ip.OnRelease(func(handled, inside bool, event scaff.ReleaseEvent) bool {
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
