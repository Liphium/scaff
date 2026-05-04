package scaffui

import (
	"image/color"

	"github.com/Liphium/scaff/smath"
	"github.com/hajimehoshi/ebiten/v2"
)

// This is so we can actually parse commands to JSON. Why? For live-reloading, but that's a future ambition.
type RenderCommand interface {
	ID() string
}

type RectangleCommand struct {
	Position     smath.Vec  `json:"position"`
	Size         Size       `json:"size"`
	FillColor    color.RGBA `json:"fillColor"`
	BorderRadius int        `json:"borderRadius"`
}

func (RectangleCommand) ID() string {
	return "rectangle"
}

type RectangleStrokeCommand struct {
	Position     smath.Vec  `json:"position"`
	Size         Size       `json:"size"`
	Color        color.RGBA `json:"color"`
	BorderRadius int        `json:"borderRadius"`
	Thickness    int        `json:"thickness"`
}

func (RectangleStrokeCommand) ID() string {
	return "rectangle-stroke"
}

type ImageCommand struct {
	Path       string        `json:"path"`
	Position   smath.Vec     `json:"position"`
	Size       Size          `json:"size"`
	FilterMode ebiten.Filter `json:"filterMode"`
}

func (ImageCommand) ID() string {
	return "image"
}

type TextCommand struct {
	Font     string    `json:"font"`
	Text     string    `json:"text"`
	Position smath.Vec `json:"position"`
	FontSize int       `json:"fontSize"`
	Weight   int       `json:"weight"`
}

func (TextCommand) ID() string {
	return "text"
}

type Renderer interface {
	// Clear the canvas for a new frame
	Clear()

	// Should draw one render command on top of everything else that has already been drawn.
	DrawOne(command RenderCommand)

	// Should draw all of the commands in order, the first index gets drawn first, etc.
	Draw(commands []RenderCommand)
}
