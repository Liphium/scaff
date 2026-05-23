package main

import (
	"embed"
	"image/color"
	"log"
	"time"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/paint"
	"github.com/Liphium/scaff/scaffui"
	"github.com/Liphium/scaff/scaffui/uinode"
	"github.com/Liphium/scaff/scath"
	"github.com/hajimehoshi/ebiten/v2"
)

// Embed the assets file system for getting the images.
//
//go:embed assets/*
var assetsFS embed.FS

func main() {
	ebiten.SetWindowSize(900, 600)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	scaling := scaff.NewSignal(false)
	scaleFactor := scaff.NewSignal(float64(1))

	assetManager := paint.NewAssetManager(assetsFS)
	tree := scaff.NewSceneTree("scaffui-image-sample", assetManager)

	tree.Mount(func(t *scaff.Tracker, props *scaff.RootProps) {
		props.Child(scaff.UseNode("scaling", func(props *scaff.SingleChildProps[any]) {
			props.Draw(func(node *scaff.SingleChildNode[any], c *scaff.Context, image *ebiten.Image) {
				const cycleDuration = 5 * time.Second

				cyclePosition := float64(c.Now.UnixNano()%int64(cycleDuration)) / float64(cycleDuration)
				cyclePosition *= 2
				if cyclePosition > 1 {
					cyclePosition = (2 - cyclePosition) / 2
				} else {
					cyclePosition /= 2
				}
				scaleFactor.Set(0.5 + cyclePosition*4)
			})
		}))

		props.Child(scaffui.Viewport(func(t *scaff.Tracker, props *scaffui.ViewportProps) {
			props.Child(uinode.Stack(func(t *scaff.Tracker, props *uinode.StackProps) {
				props.Child(uinode.Align(func(t *scaff.Tracker, props *uinode.AlignProps) {
					props.Horizontal(uinode.HorizontalAlignmentCenter)
					props.Vertical(uinode.VerticalAlignmentCenter)

					props.Child(uinode.Image(func(t *scaff.Tracker, props *uinode.ImageProps) {
						props.Path("assets/icon.png")
						if scaling.Track(t) {
							props.Constraints(scath.Tight(150*scaleFactor.Track(t), 150*scaleFactor.Track(t)))
						} else {
							props.Constraints(scath.Tight(100, 100))
						}
						props.Filter(ebiten.FilterPixelated)
					}))
				}))

				props.Child(uinode.Align(func(t *scaff.Tracker, props *uinode.AlignProps) {
					props.Horizontal(uinode.HorizontalAlignmentCenter)
					props.Vertical(uinode.VerticalAlignmentTop)

					props.Child(uinode.Text(func(t *scaff.Tracker, props *uinode.TextProps) {
						props.Text("ScaffUI Image Sample")
						props.FontSize(24)
						props.Wrapping(true)
						props.Color(color.White)
					}))
				}))

				props.Child(uinode.Align(func(t *scaff.Tracker, props *uinode.AlignProps) {
					props.Horizontal(uinode.HorizontalAlignmentCenter)
					props.Vertical(uinode.VerticalAlignmentBottom)

					props.Child(uinode.Clickable(func(t *scaff.Tracker, props *uinode.ClickableProps) {
						props.OnClick(func(button ebiten.MouseButton) bool {
							scaling.Set(!scaling.Value())
							return true
						})

						props.Child(uinode.Rectangle(func(t *scaff.Tracker, props *uinode.RectangleProps) {
							props.WantedConstraints(scath.Tight(100, 20))
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
	})

	g := scaff.NewGame()
	g.Goto(tree)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
