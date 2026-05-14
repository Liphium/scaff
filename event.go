package scaff

import (
	"fmt"

	"github.com/Liphium/scaff/smath"
	"github.com/hajimehoshi/ebiten/v2"
)

type EventId string

type Event interface {
	EventID() EventId
}

type PositionalEvent interface {
	Event
	Position() smath.Vec
}

// All types of events in cgui (don't work yet)
const (
	EventIdMove   EventId = "scaff::move"
	EventIdScroll EventId = "scaff::scroll"
)

// This event id has the button in it to make sure it can be marked as handled separately from events for other buttons
func EventIdDown(button ebiten.MouseButton) EventId {
	return EventId(fmt.Sprintf("scaff::down::%d", button))
}

// This event id has the button in it to make sure it can be marked as handled separately from events for other buttons
func EventIdRelease(button ebiten.MouseButton) EventId {
	return EventId(fmt.Sprintf("scaff::release::%d", button))
}

type MoveEvent struct {
	X      int
	Y      int
	DeltaX int
	DeltaY int
}

func (me MoveEvent) EventID() EventId {
	return EventIdMove
}

func (me MoveEvent) Position() smath.Vec {
	return smath.Vec{X: float64(me.X), Y: float64(me.Y)}
}

// Emitted when there is a new scroll delta. This can be either the mouse, a touchpad or potentially also swiping around on mobile.
type ScrollEvent struct {
	X       int
	Y       int
	ScrollX float64
	ScrollY float64
}

func (se ScrollEvent) EventID() EventId {
	return EventIdScroll
}

func (se ScrollEvent) Position() smath.Vec {
	return smath.Vec{X: float64(se.X), Y: float64(se.Y)}
}

// Emitted when the user presses a mouse button down.
type DownEvent struct {
	X      int
	Y      int
	Button ebiten.MouseButton
}

func (de DownEvent) EventID() EventId {
	return EventIdDown(de.Button)
}

func (de DownEvent) Position() smath.Vec {
	return smath.Vec{X: float64(de.X), Y: float64(de.Y)}
}

// Emitted when a user releases a mouse button.
type ReleaseEvent struct {
	X      int
	Y      int
	Button ebiten.MouseButton
}

func (re ReleaseEvent) EventID() EventId {
	return EventIdRelease(re.Button)
}

func (re ReleaseEvent) Position() smath.Vec {
	return smath.Vec{X: float64(re.X), Y: float64(re.Y)}
}
