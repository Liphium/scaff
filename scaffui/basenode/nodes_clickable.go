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

			ip.OnDown(func(button int) bool {
				pressed[button] = true
				return true
			})

			ip.OnRelease(func(button int) bool {
				wasPressed := pressed[button]
				pressed[button] = false

				if wasPressed {
					if fn, ok := core.Props().onClick.Value(); ok {
						return fn(button)
					}
				}
				return false
			})

			ip.OnReleaseOutside(func(button int) bool {
				pressed[button] = false
				return false
			})
		}))
	})
}
