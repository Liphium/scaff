package main

import (
	"image/color"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/paint"
	"github.com/Liphium/scaff/scaffcv"
	"github.com/Liphium/scaff/scaffcv/cvnode"
	"github.com/Liphium/scaff/scath"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(900, 600)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	cameraPosition := scaff.NewSignal(scath.Zero)
	tree := scaff.NewSceneTree("scaffcv-text", paint.NewAssetManager(nil))

	tree.Mount(func(t *scaff.Tracker, props *scaff.RootProps) {
		props.Child(0, scaffcv.Canvas(func(t *scaff.Tracker, props *scaffcv.CanvasProps) {
			t.Effect(func() {
				props.Position = cameraPosition.Track(t)
			})

			props.Child(cvnode.Stack(func(t *scaff.Tracker, props *cvnode.StackProps) {
				props.Child(0, cvnode.CameraMovement(cameraPosition, nil))

				props.Child(1, cvnode.Rectangle(func(t *scaff.Tracker, props *cvnode.RectangleProps) {
					props.Position = scath.Vec{X: -100, Y: -100}
					props.Size = scath.Vec{X: 200, Y: 200}
					props.FillColor = color.RGBA{90, 0, 0, 255}
				}))

				props.Child(2, cvnode.Text(func(t *scaff.Tracker, props *cvnode.TextProps) {
					props.Text = "Hello, world!"
					props.Position = scath.Vec{X: 0, Y: 0}
				}))
			}))
		}))
	})

	g := scaff.NewGame()
	g.Goto(tree)
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
