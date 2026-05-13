package main

import (
	"image/color"
	"log"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/scaffui"
	"github.com/Liphium/scaff/scaffui/basenode"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(900, 600)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	hovered := scaffui.NewSignal(false)

	layer := scaffui.NewInterfaceLayer(basenode.Root(func(t *scaffui.Tracker, props *basenode.RootProps) {
		props.Child(basenode.Align(func(t *scaffui.Tracker, props *basenode.AlignProps) {
			props.Horizontal(basenode.HorizontalAlignmentCenter)
			props.Vertical(basenode.VerticalAlignmentCenter)

			props.Child(basenode.Flex(func(t *scaffui.Tracker, props *basenode.FlexProps) {
				props.Child(basenode.Input(func(t *scaffui.Tracker, props *basenode.InputProps) {
					props.OnMove(func(handled, inside bool, event scaffui.MoveEvent) bool {
						hovered.Set(inside)
						return false
					})

					props.Child(basenode.Rectangle(func(t *scaffui.Tracker, props *basenode.RectangleProps) {
						props.WantedConstraints(scaffui.Tight(100, 100))
						props.BorderRadius(8)
						if hovered.Track(t) {
							props.FillColor(rainbow.Track(t))
						} else {
							props.FillColor(color.RGBA{255, 255, 255, 255})
						}
					}))
				}))
			}))
		}))
	}))

	var scaffUiSample scaff.Scene = &scaff.LayeredScene{
		ID: "scaffui-sample",
		WorldLayers: []scaff.WorldLayer{
			&rainbowLayer{},
		},
		UILayers: []scaff.UILayer{
			layer,
		},
	}

	g := scaff.NewGame()
	g.Goto(scaffUiSample)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
