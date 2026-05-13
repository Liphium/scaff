package main

import (
	"embed"
	"image/color"
	"log"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/scaffui"
	"github.com/Liphium/scaff/scaffui/basenode"
	"github.com/hajimehoshi/ebiten/v2"
)

// Embed the assets file system for getting the images.
//
//go:embed assets/*
var assetsFS embed.FS

func main() {
	ebiten.SetWindowSize(900, 600)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	scaling := scaffui.NewSignal(false)

	layer := scaffui.NewInterfaceLayer(assetsFS, basenode.Root(func(t *scaffui.Tracker, props *basenode.RootProps) {
		props.Child(basenode.Stack(func(t *scaffui.Tracker, props *basenode.StackProps) {
			props.Child(basenode.Align(func(t *scaffui.Tracker, props *basenode.AlignProps) {
				props.Horizontal(basenode.HorizontalAlignmentCenter)
				props.Vertical(basenode.VerticalAlignmentCenter)

				props.Child(basenode.Image(func(t *scaffui.Tracker, props *basenode.ImageProps) {
					props.Path("assets/icon.png")
					if scaling.Track(t) {
						props.Constraints(scaffui.Tight(int(150*scaleFactor.Track(t)), int(150*scaleFactor.Track(t))))
					} else {
						props.Constraints(scaffui.Tight(100, 100))
					}
					props.Filter(ebiten.FilterPixelated)
				}))
			}))

			props.Child(basenode.Align(func(t *scaffui.Tracker, props *basenode.AlignProps) {
				props.Horizontal(basenode.HorizontalAlignmentCenter)
				props.Vertical(basenode.VerticalAlignmentBottom)

				props.Child(basenode.Clickable(func(t *scaffui.Tracker, props *basenode.ClickableProps) {
					props.OnClick(func(button int) bool {
						scaling.Set(!scaling.Value())
						return true
					})

					props.Child(basenode.Rectangle(func(t *scaffui.Tracker, props *basenode.RectangleProps) {
						props.WantedConstraints(scaffui.Tight(100, 20))
						if scaling.Track(t) {
							props.FillColor(color.RGBA{0, 255, 0, 255})
						} else {
							props.FillColor(color.RGBA{255, 255, 255, 255})
						}
					}))
				}))
			}))
		}))
	}))

	var scaffUiSample scaff.Scene = &scaff.LayeredScene{
		ID: "scaffui-image-sample",
		WorldLayers: []scaff.WorldLayer{
			&scaleLayer{},
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
