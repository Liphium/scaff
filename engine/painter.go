package engine

import "github.com/hajimehoshi/ebiten/v2"

type Painter interface {
	// Get the current transform
	Transform() Transform

	// Set a transform
	SetTransform(transform Transform)

	// Clear the canvas for a new frame
	Clear()

	// Should draw one render command on top of everything else that has already been drawn.
	Paint(command RenderCommand)

	// Draw an image directly onto the screen
	DrawRaw(image *ebiten.Image, op *ebiten.DrawImageOptions)

	// Should draw all of the commands in order, the first index gets drawn first, etc.
	PaintMulti(commands []RenderCommand)
}
