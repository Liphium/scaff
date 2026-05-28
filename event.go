package scaff

import (
	"fmt"

	"github.com/Liphium/scaff/scath"
	"github.com/hajimehoshi/ebiten/v2"
)

type EventId string

type Event interface {
	EventID() EventId
}

type PositionalEvent interface {
	Event
	Position() scath.Vec
}

// All types of events in Scaff
const (
	EventIdMove       EventId = "scaff::move"
	EventIdScroll     EventId = "scaff::scroll"
	EventIdSizeChange EventId = "scaff::size-change"
)

// This event id has the button in it to make sure it can be marked as handled separately from events for other buttons
func EventIdPress(button ebiten.MouseButton) EventId {
	return EventId(fmt.Sprintf("scaff::down::%d", button))
}

// This event id has the button in it to make sure it can be marked as handled separately from events for other buttons
func EventIdRelease(button ebiten.MouseButton) EventId {
	return EventId(fmt.Sprintf("scaff::release::%d", button))
}

// This event id has the key in it to make sure it can be handled separately from other events.
func EventIdKeyPress(key ebiten.Key) EventId {
	return EventId(fmt.Sprintf("scaff::key-press::%d", key))
}

// This event id has the key in it to make sure it can be handled separately from other events.
func EventIdKeyRelease(key ebiten.Key) EventId {
	return EventId(fmt.Sprintf("scaff::key-release::%d", key))
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

func (me MoveEvent) Position() scath.Vec {
	return scath.Vec{X: float64(me.X), Y: float64(me.Y)}
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

func (se ScrollEvent) Position() scath.Vec {
	return scath.Vec{X: float64(se.X), Y: float64(se.Y)}
}

// Emitted when the user presses a mouse button down.
type PressEvent struct {
	X      int
	Y      int
	Button ebiten.MouseButton
}

func (de PressEvent) EventID() EventId {
	return EventIdPress(de.Button)
}

func (de PressEvent) Position() scath.Vec {
	return scath.Vec{X: float64(de.X), Y: float64(de.Y)}
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

func (re ReleaseEvent) Position() scath.Vec {
	return scath.Vec{X: float64(re.X), Y: float64(re.Y)}
}

// Emitted when the size of the application changes.
//
// This does not contain the new size as the context already does anyway.
type SizeChangeEvent struct{}

func (sce SizeChangeEvent) EventID() EventId {
	return EventIdSizeChange
}

// Emitted when a user presses a key.
type KeyPressEvent struct {
	Key ebiten.Key
}

func (kp KeyPressEvent) EventID() EventId {
	return EventIdKeyPress(kp.Key)
}

// Emitted when a user releases a key.
type KeyReleaseEvent struct {
	Key ebiten.Key
}

func (kr KeyReleaseEvent) EventID() EventId {
	return EventIdKeyRelease(kr.Key)
}
