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

var _ scaff.Scene = &StateMachineScene{}

type StateMachineScene struct {
	images      map[int]*ebiten.Image
	timeMachine *scaff.StateMachine[int64, int]
}

func (t *StateMachineScene) Load() {
	t.images = map[int]*ebiten.Image{}
	t.timeMachine = scaff.NewStateMachine(scaff.StateMachineCreate[int64, int]{
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
}

func (t *StateMachineScene) Unload() {}

func (t *StateMachineScene) GetId() string {
	return "state_machine_scene"
}

func (t *StateMachineScene) Update(c scaff.SceneContext) error {
	t.timeMachine.Update(c.Now, c.Now.UnixMilli())
	return nil
}

func (t *StateMachineScene) Draw(c scaff.SceneContext, screen *ebiten.Image) {
	t.timeMachine.Draw(c.Now, func(state int, frame scath.Timeframe) {
		text := "Scrolling text"
		if state == 1 {
			text = "is kinda cool"
		}

		if t.images[state] == nil || t.images[state].Bounds() != screen.Bounds() {
			t.images[state] = ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
		}

		bounds := screen.Bounds()
		x := bounds.Min.X + (bounds.Dx()-7*len(text))/2
		y := bounds.Min.Y + bounds.Dy()/2 - 8

		// Add a little bit of offset based on the transition direction
		if frame.IsBackwards() {
			y += frame.LerpInt(c.Now, -50, 0)
		} else {
			y += frame.LerpInt(c.Now, 50, 0)
		}

		// Draw the text at the proper location to the text image
		t.images[state].Clear()
		ebitenutil.DebugPrintAt(t.images[state], text, x, y)

		// Draw the text image with a change in opacity for a fade effect
		op := &ebiten.DrawImageOptions{}
		op.Blend = ebiten.BlendLighter
		op.ColorScale.ScaleAlpha(float32(frame.LerpFloat(c.Now, 0, 1)))
		screen.DrawImage(t.images[state], op)
	})
}

func (t *StateMachineScene) Transition(in bool) scaff.TransitionProperties {
	return scaff.NoTransition()
}

func main() {
	ebiten.SetWindowSize(900, 600)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g := scaff.NewGame()
	g.Goto(&StateMachineScene{})
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
