package uinode

import (
	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/scaffui"
	"github.com/hajimehoshi/ebiten/v2"
)

type ClickableProps struct {
	OnClick        func(button ebiten.MouseButton) bool
	OnClickOutside func() bool
	Cursor         ebiten.CursorShapeType
	*scaffui.AcceptChild
}

func Clickable(create func(t *scaff.Tracker, props *ClickableProps)) scaffui.NodeBuilder {
	pressed := make(map[ebiten.MouseButton]bool)

	// Create an input node for actually listening to the events
	return Input(func(t *scaff.Tracker, input *InputProps) {
		props := &ClickableProps{
			Cursor:      ebiten.CursorShapePointer,
			AcceptChild: input.AcceptChild,
		}
		create(t, props)

		set := false
		input.OnMove = func(handled, inside bool, event scaff.MoveEvent) bool {
			if inside {
				ebiten.SetCursorShape(props.Cursor)
				set = true
			} else if set {
				ebiten.SetCursorShape(ebiten.CursorShapeDefault)
				set = false
			}
			return false
		}

		input.OnDown = func(handled, inside bool, event scaff.PressEvent) bool {
			if inside {
				pressed[event.Button] = true
			} else {
				if props.OnClickOutside != nil {
					return props.OnClickOutside()
				}
			}
			return inside
		}

		input.OnRelease = func(handled, inside bool, event scaff.ReleaseEvent) bool {
			wasPressed := pressed[event.Button]
			pressed[event.Button] = false

			// If the event was not handled before and the button was pressed before + released inside of this element, a click has been detected
			if wasPressed && inside && !handled {
				if props.OnClick != nil {
					return props.OnClick(event.Button)
				}
				return true
			}

			return false
		}
	})
}
