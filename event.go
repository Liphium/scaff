package scaff

import (
	"github.com/Liphium/scaff/smath"
)

type Event interface {
	EventID() EventId
}

type PositionalEvent interface {
	Event
	Position() smath.Vec
}

// All types of events in cgui (don't work yet)
const (
	EventIdDown    EventId = "scaff::down"
	EventIdRelease EventId = "scaff::release"
	EventIdMove    EventId = "scaff::move"
	EventIdScroll  EventId = "scaff::scroll"
)

const (
	LeftClick   = 0
	RightClick  = 1
	MiddleClick = 2
)

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
	Button int
}

func (ce DownEvent) EventID() EventId {
	return EventIdDown
}

func (ce DownEvent) Position() smath.Vec {
	return smath.Vec{X: float64(ce.X), Y: float64(ce.Y)}
}

// Emitted when a user releases a mouse button.
type ReleaseEvent struct {
	X      int
	Y      int
	Button int
}

func (ce ReleaseEvent) EventID() EventId {
	return EventIdRelease
}

func (ce ReleaseEvent) Position() smath.Vec {
	return smath.Vec{X: float64(ce.X), Y: float64(ce.Y)}
}
