package main

import (
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

	tree := scaff.NewSceneTree("scaffcv-text", paint.NewAssetManager(nil))

	tree.Mount(func(t *scaff.Tracker, props *scaff.RootProps) {
		props.Child(scaffcv.Canvas(func(t *scaff.Tracker, props *scaffcv.CanvasProps) {
			props.CameraPosition(0, 0)

			props.Child(cvnode.Text("Hello, world!", scath.Vec{X: 0, Y: 0}))
		}))
	})

	g := scaff.NewGame()
	g.Goto(tree)
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
