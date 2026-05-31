package main

import (
	"log"
	"time"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/optional"
	"github.com/Liphium/scaff/scath"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func main() {
	ebiten.SetWindowSize(900, 600)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	images := map[int]*ebiten.Image{}
	timeMachine := scaff.NewStateMachine(scaff.StateMachineCreate[int64, int]{
		Transition: optional.With(scaff.TransitionProperties{
			Isolated: false,
			Duration: 500 * time.Millisecond,
		}),
		Default: scaff.DefaultState[int64, int](0),
		States: []*scaff.State[int64, int]{
			scaff.NewState(1, func(c int64) bool {
				return c%4000 >= 2000
			}),
		},
	})

	tree := scaff.NewSceneTree("state_machine_scene", nil)
	tree.Mount(func(t *scaff.Tracker, props *scaff.RootProps) {
		props.Child(scaff.SingleNode(scaff.SingleNodeCreate[*scaff.AcceptNoChild]{
			ID:           "state_machine",
			DefaultProps: &scaff.AcceptNoChild{},
			Create: func(props *scaff.SingleChildProps[*scaff.AcceptNoChild]) {
				props.OnUpdate = func(node *scaff.SingleChildNode[*scaff.AcceptNoChild], c *scaff.Context) error {
					timeMachine.Update(c.Now(), c.Now().UnixMilli())
					return nil
				}

				props.OnDraw = func(node *scaff.SingleChildNode[*scaff.AcceptNoChild], c *scaff.Context, screen *ebiten.Image) {
					timeMachine.Draw(c.Now(), func(state int, frame scath.Timeframe) {
						text := "Scrolling text"
						if state == 1 {
							text = "is kinda cool"
						}

						if images[state] == nil || images[state].Bounds() != screen.Bounds() {
							images[state] = ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
						}

						bounds := screen.Bounds()
						x := bounds.Min.X + (bounds.Dx()-7*len(text))/2
						y := bounds.Min.Y + bounds.Dy()/2 - 8

						// Add a little bit of offset based on the transition direction
						if frame.IsBackwards() {
							y += frame.LerpInt(c.Now(), -50, 0)
						} else {
							y += frame.LerpInt(c.Now(), 50, 0)
						}

						// Draw the text at the proper location to the text image
						images[state].Clear()
						ebitenutil.DebugPrintAt(images[state], text, x, y)

						// Draw the text image with a change in opacity for a fade effect
						op := &ebiten.DrawImageOptions{}
						op.Blend = ebiten.BlendLighter
						op.ColorScale.ScaleAlpha(float32(frame.LerpFloat(c.Now(), 0, 1)))
						screen.DrawImage(images[state], op)
					})
				}
			},
		}))
	})

	g := scaff.NewGame()
	g.Goto(tree)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
