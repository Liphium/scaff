package scaff

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type inputHandler struct {
	prevCursorX int
	prevCursorY int
}

var buttons = []ebiten.MouseButton{
	ebiten.MouseButton0,
	ebiten.MouseButton1,
	ebiten.MouseButton2,
	ebiten.MouseButton3,
	ebiten.MouseButton4,
}

// getInputEventsUpdate returns (for an update) all input events that should be emitted.
func (i *inputHandler) getInputEventsUpdate() []Event {
	var events []Event
	x, y := ebiten.CursorPosition()

	// Press events
	for _, button := range buttons {
		if inpututil.IsMouseButtonJustPressed(button) {
			events = append(events, PressEvent{
				X:      x,
				Y:      y,
				Button: button,
			})
		}
	}

	// Release events
	for _, button := range buttons {
		if inpututil.IsMouseButtonJustReleased(button) {
			events = append(events, ReleaseEvent{
				X:      x,
				Y:      y,
				Button: button,
			})
		}
	}

	// Scroll event (in case there is one)
	scrollX, scrollY := ebiten.Wheel()
	if scrollX != 0 || scrollY != 0 {
		events = append(events, ScrollEvent{
			X:       x,
			Y:       y,
			ScrollX: scrollX,
			ScrollY: scrollY,
		})
	}

	// Handle released keys before (as a key might be pressed and released within the same tick)
	for _, key := range inpututil.AppendJustReleasedKeys([]ebiten.Key{}) {
		events = append(events, KeyReleaseEvent{
			Key: key,
		})
	}

	// Handle pressed keys
	for _, key := range inpututil.AppendJustPressedKeys([]ebiten.Key{}) {
		events = append(events, KeyPressEvent{
			Key: key,
		})
	}

	return events
}

// getInputEventsDraw returns (for a draw call) all input events that should be emitted.
func (i *inputHandler) getInputEventsDraw() []Event {
	var events []Event
	x, y := ebiten.CursorPosition()

	deltaX := x - i.prevCursorX
	deltaY := y - i.prevCursorY
	i.prevCursorX = x
	i.prevCursorY = y

	// Append movement event in case there is one
	if deltaX != 0 || deltaY != 0 {
		events = append(events, MoveEvent{
			X:      x,
			Y:      y,
			DeltaX: deltaX,
			DeltaY: deltaY,
		})
	}

	return events
}
