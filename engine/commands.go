package engine

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
	BorderRadius float64     `json:"borderRadius"`
}

func (Rectangle) ID() string {
	return "rectangle"
}

type RectangleStroke struct {
	Position     scath.Vec   `json:"position"`
	Size         scath.Vec   `json:"size"`
	Color        color.Color `json:"color"`
	BorderRadius float64     `json:"borderRadius"`
	Thickness    float64     `json:"thickness"`
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

type Line struct {
	Start     scath.Vec   `json:"start"`
	End       scath.Vec   `json:"end"`
	Color     color.Color `json:"color"`
	Thickness float64     `json:"thickness"`
}

func (Line) ID() string {
	return "line"
}
