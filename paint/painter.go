package paint

import (
	"image/color"

	"github.com/Liphium/scaff/scath"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// This is so we can actually parse commands to JSON. Why? For live-reloading, but that's a future ambition.
type RenderCommand interface {
	ID() string
}

type Rectangle struct {
	Position     scath.Vec   `json:"position"`
	Size         scath.Vec   `json:"size"`
	FillColor    color.Color `json:"fillColor"`
	BorderRadius int         `json:"borderRadius"`
}

func (Rectangle) ID() string {
	return "rectangle"
}

type RectangleStroke struct {
	Position     scath.Vec   `json:"position"`
	Size         scath.Vec   `json:"size"`
	Color        color.Color `json:"color"`
	BorderRadius int         `json:"borderRadius"`
	Thickness    int         `json:"thickness"`
}

func (RectangleStroke) ID() string {
	return "rectangle-stroke"
}

type Image struct {
	Path       string        `json:"path"`
	Position   scath.Vec     `json:"position"`
	Size       scath.Vec     `json:"size"`
	FilterMode ebiten.Filter `json:"filterMode"`
}

func (Image) ID() string {
	return "image"
}

type Text struct {
	Direction      text.Direction `json:"direction"`
	Font           string         `json:"font"`
	Color          color.Color    `json:"color"`
	FontSize       float64        `json:"fontSize"`
	LineSpacing    float64        `json:"lineSpacing"`
	Text           string         `json:"text"`
	Position       scath.Vec      `json:"position"`
	PrimaryAlign   text.Align     `json:"primaryAlign"`
	SecondaryAlign text.Align     `json:"secondaryAlign"`
}

func (Text) ID() string {
	return "text"
}

type Painter interface {
	// Clear the canvas for a new frame
	Clear()

	// Should draw one render command on top of everything else that has already been drawn.
	Paint(command RenderCommand)

	// Should draw all of the commands in order, the first index gets drawn first, etc.
	PaintMulti(commands []RenderCommand)
}
