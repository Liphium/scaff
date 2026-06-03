package main

import (
	"embed"
	"image/color"
	"time"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/paint"
	"github.com/Liphium/scaff/scaffui"
	"github.com/Liphium/scaff/scaffui/uinode"
	"github.com/Liphium/scaff/scath"
	sutil "github.com/Liphium/scaff/util"
	"github.com/hajimehoshi/ebiten/v2"
)

// Embed the assets file system for getting the images.
//
//go:embed assets/*
var assetsFS embed.FS

var log = sutil.NewLogger("app")

func main() {
	ebiten.SetWindowSize(900, 600)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	scaling := scaff.NewSignal(false)
	scaleFactor := scaff.NewSignal(float64(1))

	assetManager := paint.NewAssetManager(assetsFS)
	tree := scaff.NewSceneTree("scaffui-image-sample", assetManager)

	tree.Mount(func(t *scaff.Tracker, props *scaff.RootProps) {
		props.Child(0, scaff.Standard(scaff.StandardCreate[scaff.AcceptNoChild]{
			ID: "scaling",
			Create: func(props *scaff.StandardMethods[scaff.AcceptNoChild]) {
				props.OnDraw = func(node *scaff.StandardNode[scaff.AcceptNoChild], c *scaff.Context, image *ebiten.Image) {
					const cycleDuration = 5 * time.Second

					cyclePosition := float64(c.Now().UnixNano()%int64(cycleDuration)) / float64(cycleDuration)
					cyclePosition *= 2
					if cyclePosition > 1 {
						cyclePosition = (2 - cyclePosition) / 2
					} else {
						cyclePosition /= 2
					}
					scaleFactor.Set(0.5 + cyclePosition*4)
				}
			},
		}))

		props.Child(1, scaffui.Viewport(func(t *scaff.Tracker, props *scaffui.ViewportProps) {
			props.Child(uinode.Stack(func(t *scaff.Tracker, props *uinode.StackProps) {
				props.Child(0, uinode.Align(func(t *scaff.Tracker, props *uinode.AlignProps) {
					props.HorizontalAligment.SetValue(uinode.HorizontalAlignmentCenter)
					props.VerticalAlignment.SetValue(uinode.VerticalAlignmentCenter)

					props.Child(uinode.Image(func(t *scaff.Tracker, props *uinode.ImageProps) {
						props.Path.SetValue("assets/icon.png")
						props.FilterMode.SetValue(ebiten.FilterPixelated)

						t.Effect(func() {
							if scaling.Track(t) {
								props.Constraints.SetValue(scath.Tight(150*scaleFactor.Track(t), 150*scaleFactor.Track(t)))
							} else {
								props.Constraints.SetValue(scath.Tight(100, 100))
							}
						})
					}))
				}))

				props.Child(1, uinode.Align(func(t *scaff.Tracker, props *uinode.AlignProps) {
					props.HorizontalAligment.SetValue(uinode.HorizontalAlignmentCenter)
					props.VerticalAlignment.SetValue(uinode.VerticalAlignmentTop)

					props.Child(uinode.Text(func(t *scaff.Tracker, props *uinode.TextProps) {
						props.Text = "ScaffUI Image Sample"
						props.FontSize = 24
						props.Wrapping = true
						props.Color = color.White
					}))
				}))

				props.Child(2, uinode.Align(func(t *scaff.Tracker, props *uinode.AlignProps) {
					props.HorizontalAligment.SetValue(uinode.HorizontalAlignmentCenter)
					props.VerticalAlignment.SetValue(uinode.VerticalAlignmentBottom)

					props.Child(uinode.Padding(func(t *scaff.Tracker, props *uinode.PaddingProps) {
						props.Padding.SetValue(scath.PadBottom(8))

						props.Child(uinode.Clickable(func(t *scaff.Tracker, props *uinode.ClickableProps) {
							props.OnClick.SetValue(func(button ebiten.MouseButton) bool {
								scaling.Set(!scaling.Value())
								return true
							})

							props.Child(uinode.Rectangle(func(t *scaff.Tracker, props *uinode.RectangleProps) {
								props.FillColor = color.RGBA{30, 30, 30, 255}
								props.BorderRadius = 12
								props.Padding = scath.Pad(12)
								props.StrokeThickness = 4
								props.StrokeColor = color.RGBA{50, 50, 50, 255}

								props.Child(uinode.Text(func(t *scaff.Tracker, props *uinode.TextProps) {
									props.FontSize = 24
									props.Wrapping = true

									t.Effect(func() {
										if scaling.Track(t) {
											props.Text = "Scale animation: ON"
										} else {
											props.Text = "Scale animation: OFF"
										}
									})
								}))
							}))
						}))
					}))
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
