package uinode

import (
	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/optional"
	"github.com/Liphium/scaff/scaffui"
	"github.com/hajimehoshi/ebiten/v2"
)

type ClickableProps struct {
	OnClick optional.O[func(button ebiten.MouseButton) bool]
	cursor  ebiten.CursorShapeType
	*scaffui.AcceptChild
}

func Clickable(create func(t *scaff.Tracker, props *ClickableProps)) scaffui.NodeBuilder {
	pressed := make(map[ebiten.MouseButton]bool)

	// Create an input node for actually listening to the events
	return Input(func(t *scaff.Tracker, input *InputProps) {
		props := &ClickableProps{
			cursor:      ebiten.CursorShapePointer,
			AcceptChild: &scaffui.AcceptChild{},
		}
		create(t, props)

		// Pass the child to the input node
		if builder, ok := props.GetChild().Value(); ok {
			input.Child(builder)
		}

		input.OnMove = func(handled, inside bool, event scaff.MoveEvent) bool {
			if inside {
				ebiten.SetCursorShape(props.cursor)
			} else {
				ebiten.SetCursorShape(ebiten.CursorShapeDefault)
			}
			return false
		}

		input.OnDown = func(handled, inside bool, event scaff.PressEvent) bool {
			if inside {
				pressed[event.Button] = true
			}
			return inside
		}

		input.OnRelease = func(handled, inside bool, event scaff.ReleaseEvent) bool {
			wasPressed := pressed[event.Button]
			pressed[event.Button] = false

			// If the event was not handled before and the button was pressed before + released inside of this element, a click has been detected
			if wasPressed && inside && !handled {
				if fn, ok := props.OnClick.Value(); ok {
					return fn(event.Button)
				}
				return true
			}

			return false
		}
	})
}
