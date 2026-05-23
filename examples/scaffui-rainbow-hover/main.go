package main

import (
	"image/color"
	"log"
	"math"
	"time"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/scaffui"
	"github.com/Liphium/scaff/scaffui/uinode"
	"github.com/Liphium/scaff/scath"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(900, 600)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	hovered := scaff.NewSignal(false)
	rainbow := scaff.NewSignal(color.RGBA{255, 255, 255, 255})

	tree := scaff.NewSceneTree("scaffui-rainbow-hover", nil)

	tree.Mount(func(t *scaff.Tracker, props *scaff.RootProps) {
		props.Child(scaff.UseNode("rainbow", func(props *scaff.SingleChildProps[any]) {
			props.Draw(func(node *scaff.SingleChildNode[any], c *scaff.Context, image *ebiten.Image) {
				const cycleDuration = 10 * time.Second

				cyclePosition := float64(c.Now.UnixNano()%int64(cycleDuration)) / float64(cycleDuration)
				rainbow.Set(hsvToRGBA(cyclePosition, 1, 1))
			})
		}))

		props.Child(scaffui.Viewport(func(t *scaff.Tracker, props *scaffui.ViewportProps) {
			props.Child(uinode.Align(func(t *scaff.Tracker, props *uinode.AlignProps) {
				props.Horizontal(uinode.HorizontalAlignmentCenter)
				props.Vertical(uinode.VerticalAlignmentCenter)

				props.Child(uinode.Flex(func(t *scaff.Tracker, props *uinode.FlexProps) {
					props.Child(uinode.Input(func(t *scaff.Tracker, props *uinode.InputProps) {
						props.OnMove(func(handled, inside bool, event scaff.MoveEvent) bool {
							hovered.Set(inside)
							return false
						})

						props.Child(uinode.Rectangle(func(t *scaff.Tracker, props *uinode.RectangleProps) {
							props.WantedConstraints(scath.Tight(100, 100))
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
	})

	g := scaff.NewGame()
	g.Goto(tree)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func hsvToRGBA(h, s, v float64) color.RGBA {
	h = math.Mod(h, 1)
	if h < 0 {
		h += 1
	}

	segment := h * 6
	chroma := v * s
	x := chroma * (1 - math.Abs(math.Mod(segment, 2)-1))
	m := v - chroma

	var r, g, b float64
	switch {
	case segment < 1:
		r, g, b = chroma, x, 0
	case segment < 2:
		r, g, b = x, chroma, 0
	case segment < 3:
		r, g, b = 0, chroma, x
	case segment < 4:
		r, g, b = 0, x, chroma
	case segment < 5:
		r, g, b = x, 0, chroma
	default:
		r, g, b = chroma, 0, x
	}

	return color.RGBA{
		R: uint8((r + m) * 255),
		G: uint8((g + m) * 255),
		B: uint8((b + m) * 255),
		A: 255,
	}
}
