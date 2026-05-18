package basenode

import (
	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/optional"
	"github.com/Liphium/scaff/scaffui"
	"github.com/hajimehoshi/ebiten/v2"
)

type ClickableProps struct {
	onClick optional.O[func(button ebiten.MouseButton) bool]
	*scaffui.AcceptChild
}

func (cp *ClickableProps) OnClick(fn func(button ebiten.MouseButton) bool) {
	cp.onClick.SetValue(fn)
}

func Clickable(create func(t *scaff.Tracker, props *ClickableProps)) scaffui.NodeBuilder {
	pressed := make(map[ebiten.MouseButton]bool)

	// Create an input node for actually listening to the events
	return Input(func(t *scaff.Tracker, input *InputProps) {
		props := &ClickableProps{
			AcceptChild: &scaffui.AcceptChild{},
		}
		create(t, props)

		// Pass the child to the input node
		if builder, ok := props.GetChild().Value(); ok {
			input.Child(builder)
		}

		input.OnDown(func(handled, inside bool, event scaff.DownEvent) bool {
			log.Debug("clicky")
			if inside {
				pressed[event.Button] = true
			}
			return inside
		})

		input.OnRelease(func(handled, inside bool, event scaff.ReleaseEvent) bool {
			wasPressed := pressed[event.Button]
			pressed[event.Button] = false

			// If the event was not handled before and the button was pressed before + released inside of this element, a click has been detected
			if wasPressed && inside && !handled {
				if fn, ok := props.onClick.Value(); ok {
					return fn(event.Button)
				}
				return true
			}

			return false
		})
	})
}
